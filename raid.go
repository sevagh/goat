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
	log.Printf("Checking if drives exist")

	var raidDriveName string
	var err error
	if raidDriveName, err = RandRaidDriveNamePicker(); err != nil {
		return err
	}

	cmd := "mdadm"

	argsExist := []string{
		raidDriveName,
	}

	log.Printf("Checking if %s exists in mdadm", raidDriveName)
	_, err = ExecuteCommand(cmd, argsExist)

	driveNames := []string{}
	for _, drive := range drives {
		driveNames = append(driveNames, drive.AttachedName)
	}

	if DryRun || err != nil {
		log.Printf("Raid drive doesn't exist, creating")
		args := []string{
			"--create",
			raidDriveName,
			"--level=" + strconv.Itoa(raidLevel),
			"--name=KRAKEN-" + volName,
			"--raid-devices=" + strconv.Itoa(len(driveNames)),
		}
		args = append(args, driveNames...)
		log.Printf("Executing: %s %s\n", cmd, args)
		if _, err := ExecuteCommand(cmd, args); err != nil {
			log.Printf("%v", err)
			return err
		}
	}

	return MountSingleDrive(raidDriveName, mountPath, drives[0].FsType, volName)
}
