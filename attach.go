package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes []EbsVol, dryRun bool) (map[int][]EbsVol, error) {
	drivesToMount := map[int][]EbsVol{}
	ctr := 0

	log.Printf("Now attaching EBS volumes")
	var letterRunes = []rune("bcdefghijklmnopqrstuvwxyz")
	for _, volume := range volumes {
		if volume.AttachedName != "" {
			log.Printf("%s already attached\n", volume.EbsVolId)
			drivesToMount[volume.VolumeId] = append(drivesToMount[volume.VolumeId], volume)
		} else {
			log.Printf("Picking a drive that doesn't exist")
			var deviceName string
			for {
				if ctr >= len(letterRunes) {
					return drivesToMount, fmt.Errorf("Ran out of drive letter names")
				}
				deviceName = "/dev/xvd" + string(letterRunes[ctr])
				ctr++
				if !DoesDriveExist(deviceName) {
					break
				}
			}
			log.Printf("Executing AWS SDK attach command on attached volume %s", deviceName)
			attachVolIn := &ec2.AttachVolumeInput{
				Device:     &deviceName,
				InstanceId: &ec2Instance.InstanceId,
				VolumeId:   &volume.EbsVolId,
				DryRun:     &dryRun,
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
		}
	}
	return drivesToMount, nil
}
