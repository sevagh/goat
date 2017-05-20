package main

import (
	"fmt"
	"log"
	"strconv"
)

func MountRaidDrives(drives []EbsVol, volName string) error {
	log.Printf("Mounting raid drives")
	raidLevel := drives[0].RaidLevel
	mountPath := drives[0].MountPath

	if raidLevel != 0 && raidLevel != 1 {
		return fmt.Errorf("Valid raid levels are 0 and 1")
	}

	var raidDriveName string
	var err error
	if raidDriveName, err = RandRaidDriveNamePicker(); err != nil {
		return err
	}

	cmd := "mdadm"

	argsExist := []string{
		raidDriveName,
	}

	_, err = ExecuteCommand(cmd, argsExist)

	driveNames := []string{}
	for _, drive := range drives {
		driveNames = append(driveNames, drive.AttachedName)
	}

	if DryRun || err != nil {
		args := []string{
			"--create",
			raidDriveName,
			"--level=" + strconv.Itoa(raidLevel),
			"--name=\"" + PREFIX + "-" + volName + "\"",
			"--raid-devices=" + strconv.Itoa(len(driveNames)),
		}
		args = append(args, driveNames...)
		log.Printf("Creating RAID drive: %s %s", cmd, args)
		if _, err := ExecuteCommand(cmd, args); err != nil {
			return err
		}
	}

	return MountSingleDrive(raidDriveName, mountPath, drives[0].FsType, volName)
}
