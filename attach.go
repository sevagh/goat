package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes []EbsVol, logger *log.Logger) (map[int][]string, error) {
	ret := map[int][]string{}
	drivesToMount := map[int][]EbsVol{}
	ctr := 0

	var letterRunes = []rune("bcdefghijklmnopqrstuvwxyz")
	for _, volume := range volumes {
		drivesToMount[volume.VolumeId] = []EbsVol{}
		if volume.AttachedName != "" {
			logger.Printf("%s already attached\n", volume.EbsVolId)
			drivesToMount[volume.VolumeId] = append(drivesToMount[volume.VolumeId], volume)
		} else {
			var deviceName string
			for {
				deviceName = "/dev/xvd" + string(letterRunes[ctr])
				ctr++
				if !DoesDriveExist(deviceName, logger) {
					break
				}
			}
			attachVolIn := &ec2.AttachVolumeInput{
				Device:     &deviceName,
				InstanceId: &ec2Instance.InstanceId,
				VolumeId:   &volume.EbsVolId,
			}
			volAttachments, err := ec2Instance.Ec2Client.AttachVolume(attachVolIn)
			if err != nil {
				return ret, err
			}
			logger.Println(volAttachments)
			volume.AttachedName = deviceName

			drivesToMount[volume.VolumeId] = append(drivesToMount[volume.VolumeId], volume)
		}
	}

	for volId, volumes := range drivesToMount {
		//check for volume mismatch
		ret[volId] = []string{}
		volSize := volumes[0].VolumeSize
		if len(volumes) == 1 && volSize == 1 {
			ret[volId] = append(ret[volId], volumes[0].AttachedName)
			continue
		} else {
			for _, vol := range volumes[1:] {
				if volSize != vol.VolumeSize {
					return ret, fmt.Errorf("Mismatched VolumeSize tags among disks of same volume")
				}
				ret[volId] = append(ret[volId], vol.AttachedName)
			}
			if len(volumes) != volSize {
				return ret, fmt.Errorf("Found %d volumes, expected %d from VolumeSize tag", len(volumes), volSize)
			}
		}
	}
	for _, attachedDrives := range ret {
		logger.Printf("Attached drives: %s\n", attachedDrives)
	}
	return ret, nil
}
