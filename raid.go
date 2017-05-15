package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func MountRaidDrives(drives []EbsVol, volId int, logger *log.Logger) error {
	logger.Printf("Mounting raid drives")
	raidLevel := drives[0].RaidLevel
	mountPath := drives[0].MountPath

	if raidLevel != 0 && raidLevel != 1 {
		return fmt.Errorf("Valid raid levels are 0 and 1")
	}
	logger.Printf("Checking if drives exist")

	driveNames := []string{}
	for _, drive := range drives {
		var attempts int
		for driveExists := false; driveExists == true; driveExists = DoesDriveExist(drive.AttachedName, logger) {
			time.Sleep(time.Duration(1 * time.Second))
			attempts++
			if attempts >= statAttempts {
				logger.Printf("Exceeded max (%d) stat attempts waiting for drive %s to exist", statAttempts, drive.AttachedName)
				return fmt.Errorf("Stat failed")
			}
		}
		driveNames = append(driveNames, drive.AttachedName)
	}

	raidDriveName := "/dev/md" + strconv.Itoa(volId)

	cmd := "mdadm"
	args := []string{
		"--create",
		raidDriveName,
		"--level=" + strconv.Itoa(raidLevel),
		"--name=KRAKEN",
		"--raid-devices=" + strconv.Itoa(len(driveNames)),
	}
	args = append(args, driveNames...)
	logger.Printf("Executing: %s %s\n", cmd, args)
	if _, err := ExecuteCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}

	return MountSingleDrive(raidDriveName, mountPath, drives[0].FsType, logger)
}
