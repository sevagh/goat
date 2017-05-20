package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CreateRaidArray(drives []EbsVol, volName string, dryRun bool) string {
	raidLogger := log.WithFields(log.Fields{"vol_name": volName, "drives": drives})

	raidLogger.Info("Mounting raid drives")
	raidLevel := drives[0].RaidLevel

	var raidDriveName string
	var err error
	if raidDriveName, err = RandRaidDriveNamePicker(); err != nil {
		raidLogger.Fatalf("Couldn't select unused raid drive name: %v", err)
	}

	cmd := "mdadm"

	driveNames := []string{}
	for _, drive := range drives {
		driveNames = append(driveNames, drive.AttachedName)
	}

	args := []string{
		"--create",
		raidDriveName,
		"--level=" + strconv.Itoa(raidLevel),
		"--name=\"" + PREFIX + "-" + volName + "\"",
		"--raid-devices=" + strconv.Itoa(len(driveNames)),
	}
	args = append(args, driveNames...)
	log.Info("Creating RAID drive: %s %s", cmd, args)
	if !dryRun {
		if _, err := ExecuteCommand(cmd, args); err != nil {
			raidLogger.Fatalf("Error when executing mdadm command: %v", err)
		}
	}

	return raidDriveName
}
