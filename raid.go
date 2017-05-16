package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func MountRaidDrives(drives []EbsVol, volId int) error {
	log.Printf("Mounting raid drives")
	raidLevel := drives[0].RaidLevel
	mountPath := drives[0].MountPath

	if raidLevel != 0 && raidLevel != 1 {
		return fmt.Errorf("Valid raid levels are 0 and 1")
	}
	log.Printf("Checking if drives exist")

	driveNames := []string{}
	for _, drive := range drives {
		log.Printf("Checking if drive %s exists", drive.AttachedName)
		var attempts int
		for !DoesDriveExist(drive.AttachedName) {
			time.Sleep(time.Duration(1 * time.Second))
			attempts++
			if attempts >= statAttempts {
				log.Printf("Exceeded max (%d) stat attempts waiting for drive %s to exist", statAttempts, drive.AttachedName)
				return fmt.Errorf("Stat failed")
			}
		}
		driveNames = append(driveNames, drive.AttachedName)
	}

	raidDriveName := "/dev/md" + strconv.Itoa(volId)

	cmd := "mdadm"

	argsExist := []string{
		raidDriveName,
	}

	log.Printf("Checking if %s exists in mdadm", raidDriveName)
	_, err := ExecuteCommand(cmd, argsExist);
	if DryRun || err != nil {
		log.Printf("Raid drive doesn't exist, creating")
		args := []string{
			"--create",
			raidDriveName,
			"--level=" + strconv.Itoa(raidLevel),
			"--name=KRAKEN" + strconv.Itoa(volId),
			"--raid-devices=" + strconv.Itoa(len(driveNames)),
		}
		args = append(args, driveNames...)
		log.Printf("Executing: %s %s\n", cmd, args)
		if _, err := ExecuteCommand(cmd, args); err != nil {
			log.Printf("%v", err)
			return err
		}
	}

	return MountSingleDrive(raidDriveName, mountPath, drives[0].FsType)
}
