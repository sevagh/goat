package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes map[string][]EbsVol) (map[string][]EbsVol, error) {
	var deviceName string
	var err error

	localVolumes := map[string][]EbsVol{}

	log.Printf("Now attaching EBS volumes")
	for key, volumes_ := range volumes {
		localVolumes[key] = []EbsVol{}
		for _, volume := range volumes_ {
			if volume.AttachedName != "" {
				log.Printf("%s already attached\n", volume.EbsVolId)
			} else {
				log.Printf("Picking a drive that doesn't exist")
				if deviceName, err = RandDriveNamePicker(); err != nil {
					return localVolumes, err
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
					return localVolumes, err
				}
				log.Println(volAttachments)
				volume.AttachedName = deviceName

				if !DoesDriveExistWithTimeout(deviceName) {
					return localVolumes, fmt.Errorf("Drive %s doesn't exist", deviceName)
				}

				localVolumes[key] = append(localVolumes[key], volume)
			}

		}
	}
	log.Printf("VOL MAP AFTER ATTACH: %s", localVolumes)
	return localVolumes, nil
}
