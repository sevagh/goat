package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes map[string][]EbsVol) error {
	var deviceName string
	var err error

	log.Printf("Now attaching EBS volumes")
	for _, volumes_ := range volumes {
		for _, volume := range volumes_ {
			if volume.AttachedName != "" {
				log.Printf("%s already attached\n", volume.EbsVolId)
			} else {
				log.Printf("Picking a drive that doesn't exist")
				if deviceName, err = RandDriveNamePicker(); err != nil {
					return err
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
					return err
				}
				log.Println(volAttachments)
				volume.AttachedName = deviceName
			}
			if !DoDrivesExist(volumes_) {
				return fmt.Errorf("Attached drives can't be stat")
			}
		}
	}
	return nil
}
