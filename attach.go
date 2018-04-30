package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"

	"github.com/sevagh/goat/filesystem"
)

//AttachEbsVolumes attaches the given map of {'VolumeName':[]EbsVol} with the EC2 client in the provided ec2Instance
func (e *EC2Instance) AttachEbsVolumes() map[string][]EbsVol {
	var deviceName string
	var err error

	localVolumes := map[string][]EbsVol{}

	for key, volumes := range e.Vols {
		localVolumes[key] = []EbsVol{}
		for _, volume := range volumes {
			volLogger := log.WithFields(log.Fields{"vol_id": volume.EbsVolID, "vol_name": volume.VolumeName})
			if volume.AttachedName == "" {
				volLogger.Info("Volume is unattached, picking drive name")
				if deviceName, err = randDriveNamePicker(); err != nil {
					volLogger.Fatal("Couldn't find an unused drive name")
				}
				attachVolIn := &ec2.AttachVolumeInput{
					Device:     &deviceName,
					InstanceId: &e.InstanceID,
					VolumeId:   &volume.EbsVolID,
				}
				volLogger.Info("Executing AWS SDK attach command")
				volAttachments, err := e.EC2Client.AttachVolume(attachVolIn)
				if err != nil {
					volLogger.Fatalf("Couldn't attach: %v", err)
				}
				volLogger.Info(volAttachments)
				volume.AttachedName = deviceName

				if !filesystem.DoesDriveExistWithTimeout(deviceName, 10) {
					volLogger.Fatalf("Drive %s doesn't exist after attaching - checked with stat 10 times", deviceName)
				}
				localVolumes[key] = append(localVolumes[key], volume)
			}

		}
	}
	return localVolumes
}

//AttachEnis attaches the given array of Eni Ids with the EC2 client in the provided ec2Instance
func (e *EC2Instance) AttachEnis() {
	for eniIdx, eni := range e.Enis {
		eniLogger := log.WithFields(log.Fields{"eni_id": eni})

		deviceIdx := int64(eniIdx + 1)
		attachEniIn := &ec2.AttachNetworkInterfaceInput{
			NetworkInterfaceId: &eni,
			InstanceId:         &e.InstanceID,
			DeviceIndex:        &deviceIdx,
		}

		eniLogger.Info("Executing AWS SDK attach command")
		_, err := e.EC2Client.AttachNetworkInterface(attachEniIn)
		if err != nil {
			eniLogger.Fatalf("Couldn't attach: %v", err)
		}
	}
}

func randDriveNamePicker() (string, error) {
	ctr := 0
	deviceName := "/dev/xvd"
	runes := []rune("bcdefghijklmnopqrstuvwxyz")
	for {
		if ctr >= len(runes) {
			return "", fmt.Errorf("Ran out of drive names")
		}
		if !filesystem.DoesDriveExist(deviceName + string(runes[ctr])) {
			break
		}
		ctr++
	}
	return deviceName + string(runes[ctr]), nil
}
