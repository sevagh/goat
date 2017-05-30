package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//EbsVol is a struct defining the discovered EBS volumes and its metadata parsed from the tags
type EbsVol struct {
	EbsVolID     string
	VolumeName   string
	RaidLevel    int
	VolumeSize   int
	AttachedName string
	MountPath    string
	FsType       string
	Touched      bool
}

//MapEbsVolumes discovers and creates a {'VolumeName':[]EbsVol} map for all the required EBS volumes given an EC2Instance struct
func MapEbsVolumes(ec2Instance *EC2Instance) map[string][]EbsVol {
	drivesToMount := map[string][]EbsVol{}

	log.Info("Searching for EBS volumes with previously established EC2 client")

	volumes, err := findEbsVolumes(ec2Instance)
	if err != nil {
		log.Fatal("Error when searching for EBS volumes")
	}

	log.Info("Classifying EBS volumes based on tags")
	for _, volume := range volumes {
		drivesToMount[volume.VolumeName] = append(drivesToMount[volume.VolumeName], volume)
	}

	toDelete := []string{}

	for volName, volumes := range drivesToMount {
		volGroupLogger := log.WithFields(log.Fields{"vol_name": volName})

		//check for volume mismatch
		volSize := volumes[0].VolumeSize
		mountPath := volumes[0].MountPath
		fsType := volumes[0].FsType
		raidLevel := volumes[0].RaidLevel
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

	return drivesToMount
}

func findEbsVolumes(ec2Instance *EC2Instance) ([]EbsVol, error) {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:" + PREFIX + "-IN:Prefix"),
				Values: []*string{
					aws.String(ec2Instance.Prefix),
				},
			},
			&ec2.Filter{
				Name: aws.String("tag:" + PREFIX + "-IN:NodeId"),
				Values: []*string{
					aws.String(ec2Instance.NodeID),
				},
			},
			&ec2.Filter{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(ec2Instance.Az),
				},
			},
		},
	}

	volumes := []EbsVol{}

	result, err := ec2Instance.EC2Client.DescribeVolumes(params)
	if err != nil {
		return volumes, err
	}

	for _, volume := range result.Volumes {
		ebsVolume := EbsVol{
			EbsVolID: *volume.VolumeId,
		}
		if len(volume.Attachments) > 0 {
			for _, attachment := range volume.Attachments {
				if *attachment.InstanceId != ec2Instance.InstanceID {
					return volumes, fmt.Errorf("Volume %s attached to different instance-id: %s", *volume.VolumeId, *attachment.InstanceId)
				}
				ebsVolume.AttachedName = *attachment.Device
			}
		} else {
			ebsVolume.AttachedName = ""
		}
		tagCtr := 0
		ebsVolume.Touched = false
		for _, tag := range volume.Tags {
			switch *tag.Key {
			case PREFIX + "-IN:VolumeName":
				ebsVolume.VolumeName = *tag.Value
				tagCtr++
			case PREFIX + "-IN:RaidLevel":
				if ebsVolume.RaidLevel, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse RaidLevel tag as int: %v", err)
				}
				tagCtr++
			case PREFIX + "-IN:VolumeSize":
				if ebsVolume.VolumeSize, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse VolumeSize tag as int: %v", err)
				}
				tagCtr++
			case PREFIX + "-IN:MountPath":
				ebsVolume.MountPath = *tag.Value
				tagCtr++
			case PREFIX + "-IN:FsType":
				ebsVolume.FsType = *tag.Value
				tagCtr++
			case PREFIX + "-IN:NodeId": //do nothing
				tagCtr++
			case PREFIX + "-IN:Prefix": //do nothing
				tagCtr++
			case PREFIX + "-OUT:Touched":
				ebsVolume.Touched = true
			default:
			}
		}

		if tagCtr != 7 {
			return volumes, fmt.Errorf("Missing required KRK-IN tags VolumeName, RaidLevel, MountPath, VolumeSize, NodeId, Prefix, FsType")
		}
		volumes = append(volumes, ebsVolume)
	}
	return volumes, nil
}
