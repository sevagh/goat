package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"

	"github.com/sevagh/goat/driveutil"
)

//AttachEbsVolumes attaches the given map of {'VolumeName':[]EbsVol} with the EC2 client in the provided ec2Instance
func AttachEbsVolumes(ec2Instance EC2Instance, volumes map[string][]EbsVol, dryRun bool) map[string][]EbsVol {
	var deviceName string
	var err error

	localVolumes := map[string][]EbsVol{}

	for key, volumes := range volumes {
		localVolumes[key] = []EbsVol{}
		for _, volume := range volumes {
			volLogger := log.WithFields(log.Fields{"vol_id": volume.EbsVolID, "vol_name": volume.VolumeName})
			if volume.AttachedName == "" {
				volLogger.Info("Volume is unattached, picking drive name")
				if deviceName, err = driveutil.RandDriveNamePicker(); err != nil {
					volLogger.Fatal("Couldn't find an unused drive name")
				}
				attachVolIn := &ec2.AttachVolumeInput{
					Device:     &deviceName,
					InstanceId: &ec2Instance.InstanceID,
					VolumeId:   &volume.EbsVolID,
					DryRun:     &dryRun,
				}
				volLogger.Info("Executing AWS SDK attach command")
				volAttachments, err := ec2Instance.EC2Client.AttachVolume(attachVolIn)
				if err != nil {
					volLogger.Fatalf("Couldn't attach: %v", err)
				}
				volLogger.Info(volAttachments)
				volume.AttachedName = deviceName

				if !dryRun && !driveutil.DoesDriveExistWithTimeout(deviceName, 10) {
					volLogger.Fatalf("Drive %s doesn't exist after attaching - checked with stat 10 times", deviceName)
				}
				localVolumes[key] = append(localVolumes[key], volume)
			}

		}
	}
	return localVolumes
}
