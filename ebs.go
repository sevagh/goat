package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/sevagh/goat/filesystem"
)

//GoatEbs runs Goat for your EBS volumes - attach, mount, mkfs, etc.
func GoatEbs(debug bool, tagPrefix string) {
	log.Printf("WELCOME TO GOAT")
	log.Printf("1: COLLECTING EC2 INFO")
	ec2Instance := GetEC2InstanceData(tagPrefix)

	log.Printf("2: COLLECTING EBS INFO")
	ec2Instance.FindEbsVolumes(tagPrefix)

	log.Printf("3: ATTACHING EBS VOLS")
	ec2Instance.AttachEbsVolumes()

	log.Printf("4: MOUNTING ATTACHED VOLS")

	if len(ec2Instance.Vols) == 0 {
		log.Warn("Empty vols, nothing to do")
		os.Exit(0)
	}

	for volName, vols := range ec2Instance.Vols {
		prepAndMountDrives(volName, vols)
	}
}

func prepAndMountDrives(volName string, vols []EbsVol) {
	driveLogger := log.WithFields(log.Fields{"vol_name": volName, "vols": vols})

	mountPath := vols[0].MountPath
	desiredFs := vols[0].FsType
	raidLevel := vols[0].RaidLevel

	if volName == "" {
		driveLogger.Info("No volume name given, not performing further actions")
		return
	}

	if filesystem.DoesDriveExist("/dev/disk/by-label/GOAT-" + volName) {
		driveLogger.Info("Label already exists, jumping to mount phase")
	} else {
		var driveName string
		var err error
		if len(vols) == 1 {
			driveLogger.Info("Single drive, no RAID")
			driveName, err = filesystem.GetActualBlockDeviceName(vols[0].AttachedName)
			if err != nil {
				driveLogger.Fatalf("Block device is not available %s : %v", vols[0].AttachedName, err)
			}
			driveLogger.Infof("Using %s as local block device", driveName)
		} else {
			if raidLevel == -1 {
				driveLogger.Info("Raid level not provided, not performing further actions")
				return
			}
			driveLogger.Info("Creating RAID array")
			driveNames := []string{}
			for _, vol := range vols {
				driveName, err = filesystem.GetActualBlockDeviceName(vol.AttachedName)
				if err != nil {
					driveLogger.Fatalf("Block device is not available %s : %v", vol.AttachedName, err)
				}
				driveLogger.Infof("Using %s as local block device", driveName)
				driveNames = append(driveNames, driveName)
			}
			if driveName, err = filesystem.CreateRaidArray(driveNames, volName, raidLevel); err != nil {
				driveLogger.Fatalf("Error when creating reaid array: %v", err)
			}
		}

		if desiredFs == "" {
			driveLogger.Info("Desired filesystem not provided, not performing further actions")
			return
		}

		driveLogger.Info("Checking for existing filesystem")

		if err := filesystem.CheckFilesystem(driveName, desiredFs, volName); err != nil {
			driveLogger.Fatalf("Checking for existing filesystem: %v", err)
		}
		if err := filesystem.CreateFilesystem(driveName, desiredFs, volName); err != nil {
			driveLogger.Fatalf("Error when creating filesystem: %v", err)
		}
	}

	if mountPath == "" {
		driveLogger.Info("Mount point not provided, not performing further actions")
		return
	}

	driveLogger.Info("Checking if something already mounted at %s", mountPath)
	if isMounted, err := filesystem.IsMounted(mountPath); err != nil {
		driveLogger.Fatalf("Error when checking mount point for existing mounts: %v", err)
	} else {
		if isMounted {
			driveLogger.Fatalf("Something already mounted at %s", mountPath)
		}
	}

	if err := os.MkdirAll(mountPath, 0777); err != nil {
		driveLogger.Fatalf("Couldn't mkdir: %v", err)
	}

	driveLogger.Info("Appending fstab entry")
	if err := filesystem.AppendToFstab("GOAT-"+volName, desiredFs, mountPath); err != nil {
		driveLogger.Fatalf("Couldn't append to fstab: %v", err)
	}

	driveLogger.Info("Now mounting")
	if err := syscall.Mount("", mountPath, "", 0, ""); err != nil {
		driveLogger.Fatalf("Couldn't mount: %v", err)
	}

	if len(vols) > 1 {
		driveLogger.Info("Now persisting mdadm conf")
		if err := filesystem.PersistMdadm(); err != nil {
			driveLogger.Fatalf("Couldn't persist mdadm conf: %v", err)
		}
	}
}

//EbsVol is a struct defining the discovered EBS volumes and its metadata parsed from the tags
type EbsVol struct {
	EbsVolID     string
	VolumeName   string
	RaidLevel    int
	VolumeSize   int
	AttachedName string
	MountPath    string
	FsType       string
}

//FindEbsVolumes discovers and creates a {'VolumeName':[]EbsVol} map for all the required EBS volumes given an EC2Instance struct
func (e *EC2Instance) FindEbsVolumes(tagPrefix string) {
	drivesToMount := map[string][]EbsVol{}

	log.Info("Searching for EBS volumes")

	volumes, err := e.findEbsVolumes(tagPrefix)
	if err != nil {
		log.Fatalf("Error when searching for EBS volumes: %v", err)
	}

	log.Info("Classifying EBS volumes based on tags")
	for _, volume := range volumes {
		drivesToMount[volume.VolumeName] = append(drivesToMount[volume.VolumeName], volume)
	}

	for volName, volumes := range drivesToMount {
		volGroupLogger := log.WithFields(log.Fields{"vol_name": volName})

		//check for volume mismatch
		volSize := volumes[0].VolumeSize
		mountPath := volumes[0].MountPath
		fsType := volumes[0].FsType
		raidLevel := volumes[0].RaidLevel
		if volSize != -1 {
			if len(volumes) != volSize {
				volGroupLogger.Fatalf("Found %d volumes, expected %d from VolumeSize tag", len(volumes), volSize)
			}
			for _, vol := range volumes[1:] {
				volLogger := log.WithFields(log.Fields{"vol_id": vol.EbsVolID, "vol_name": vol.VolumeName})
				if volSize != vol.VolumeSize || mountPath != vol.MountPath || fsType != vol.FsType || raidLevel != vol.RaidLevel {
					volLogger.Fatal("Mismatched tags among disks of same volume")
				}
			}
		}
	}

	e.Vols = drivesToMount
}

func (e *EC2Instance) findEbsVolumes(tagPrefix string) ([]EbsVol, error) {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:" + tagPrefix + ":Prefix"),
				Values: []*string{
					aws.String(e.Prefix),
				},
			},
			{
				Name: aws.String("tag:" + tagPrefix + ":NodeId"),
				Values: []*string{
					aws.String(e.NodeID),
				},
			},
			{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(e.Az),
				},
			},
		},
	}

	volumes := []EbsVol{}

	result, err := e.EC2Client.DescribeVolumes(params)
	if err != nil {
		return volumes, err
	}

	for _, volume := range result.Volumes {
		ebsVolume := EbsVol{
			EbsVolID:   *volume.VolumeId,
			VolumeName: "",
			RaidLevel:  -1,
			VolumeSize: -1,
			MountPath:  "",
			FsType:     "",
		}
		if len(volume.Attachments) > 0 {
			for _, attachment := range volume.Attachments {
				if *attachment.InstanceId != e.InstanceID {
					return volumes, fmt.Errorf("Volume %s attached to different instance-id: %s", *volume.VolumeId, *attachment.InstanceId)
				}
				ebsVolume.AttachedName = *attachment.Device
			}
		} else {
			ebsVolume.AttachedName = ""
		}
		for _, tag := range volume.Tags {
			switch *tag.Key {
			case tagPrefix + ":VolumeName":
				ebsVolume.VolumeName = *tag.Value
			case tagPrefix + ":RaidLevel":
				if ebsVolume.RaidLevel, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse RaidLevel tag as int: %v", err)
				}
			case tagPrefix + ":VolumeSize":
				if ebsVolume.VolumeSize, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse VolumeSize tag as int: %v", err)
				}
			case tagPrefix + ":MountPath":
				ebsVolume.MountPath = *tag.Value
			case tagPrefix + ":FsType":
				ebsVolume.FsType = *tag.Value
			default:
			}
		}
		volumes = append(volumes, ebsVolume)
	}
	return volumes, nil
}
