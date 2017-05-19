package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes []EbsVol) (map[int][]EbsVol, error) {
	drivesToMount := map[int][]EbsVol{}
	var deviceName string
	var err error

	log.Printf("Now attaching EBS volumes")
	for _, volume := range volumes {
		if volume.AttachedName != "" {
			log.Printf("%s already attached\n", volume.EbsVolId)
			drivesToMount[volume.VolumeId] = append(drivesToMount[volume.VolumeId], volume)
		} else {
			log.Printf("Picking a drive that doesn't exist")
			if deviceName, err = RandDriveNamePicker(); err != nil {
				return drivesToMount, err
			}
			log.Printf("Executing AWS SDK attach command on attached volume %s", deviceName)
			attachVolIn := &ec2.AttachVolumeInput{
				Device:     &deviceName,
				InstanceId: &ec2Instance.InstanceId,
				VolumeId:   &volume.EbsVolId,
				DryRun:     &DryRun,
			}
			volAttachments, err := ec2Instance.Ec2Client.AttachVolume(attachVolIn)
			if err != nil {
				return drivesToMount, err
			}
			log.Println(volAttachments)
			volume.AttachedName = deviceName

			drivesToMount[volume.VolumeId] = append(drivesToMount[volume.VolumeId], volume)
		}
	}

	return sanityCheckVols(drivesToMount)
}

func sanityCheckVols(drivesToMount map[int][]EbsVol) (map[int][]EbsVol, error) {
	for _, volumes := range drivesToMount {
		//check for volume mismatch
		volSize := volumes[0].VolumeSize
		mountPath := volumes[0].MountPath
		if len(volumes) == 1 && volSize == 1 {
			continue
		} else {
			for _, vol := range volumes[1:] {
				if volSize != vol.VolumeSize {
					return drivesToMount, fmt.Errorf("Mismatched VolumeSize tags among disks of same volume")
				}
				if mountPath != vol.MountPath {
					return drivesToMount, fmt.Errorf("Mismatched MountPath tags among disks of same volume")
				}
			}
			if len(volumes) != volSize {
				return drivesToMount, fmt.Errorf("Found %d volumes, expected %d from VolumeSize tag", len(volumes), volSize)
			}

			if !DoDrivesExist(volumes) {
				return drivesToMount, fmt.Errorf("Attached drives can't be stat")
			}
		}
	}
	return drivesToMount, nil
}
