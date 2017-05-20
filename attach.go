package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

func AttachEbsVolumes(ec2Instance Ec2Instance, volumes map[string][]EbsVol, dryRun bool) map[string][]EbsVol {
	var deviceName string
	var err error

	localVolumes := map[string][]EbsVol{}

	for key, volumes_ := range volumes {
		localVolumes[key] = []EbsVol{}
		for _, volume := range volumes_ {
			volLogger := log.WithFields(log.Fields{"vol_id": volume.EbsVolId, "vol_name": volume.VolumeName})
			if volume.AttachedName == "" {
				volLogger.Info("Volume is unattached, picking drive name")
				if deviceName, err = RandDriveNamePicker(); err != nil {
					volLogger.Fatal("Couldn't find an unused drive name")
				}
				attachVolIn := &ec2.AttachVolumeInput{
					Device:     &deviceName,
					InstanceId: &ec2Instance.InstanceId,
					VolumeId:   &volume.EbsVolId,
					DryRun:     &dryRun,
				}
				volLogger.Info("Executing AWS SDK attach command")
				volAttachments, err := ec2Instance.Ec2Client.AttachVolume(attachVolIn)
				if err != nil {
					volLogger.Fatalf("Couldn't attach: %v", err)
				}
				volLogger.Info(volAttachments)
				volume.AttachedName = deviceName

				if !dryRun && !DoesDriveExistWithTimeout(deviceName) {
					volLogger.Fatalf("Drive %s doesn't exist after attaching - checked with stat %d times", deviceName, statAttempts)
				}
				localVolumes[key] = append(localVolumes[key], volume)
			}

		}
	}
	return localVolumes
}
