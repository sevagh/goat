package main

import (
	"fmt"
	"log"
)

func MountSingleVolume(drive EbsVol) error {
	log.Printf("Mounting single drive: %s", drive.AttachedName)
	return MountSingleDrive(drive.AttachedName, drive.MountPath, drive.FsType, drive.VolumeName)
}

func MountSingleDrive(driveName string, mountPath string, desiredFs string, label string) error {
	if err := checkAndCreateFilesystem(driveName, desiredFs, label); err != nil {
		return err
	}
	if err := mkdir(mountPath); err != nil {
		return err
	}
	cmd := "mount"
	args := []string{
		driveName,
		mountPath,
	}
	log.Printf("Executing: %s %s", cmd, args)
	if _, err := ExecuteCommand(cmd, args); err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

func mkdir(mountPath string) error {
	cmd := "mkdir"
	args := []string{
		"-p",
		mountPath,
	}
	log.Printf("Executing: %s %s", cmd, args)
	if _, err := ExecuteCommand(cmd, args); err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func checkAndCreateFilesystem(driveName string, desiredFs string, label string) error {
	log.Printf("Checking filesystem on %s", driveName)
	cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		driveName,
	}
	fsOut, err := ExecuteCommand(cmd, args)
	if DryRun {
		return nil
	}
	if err != nil {
		if fsOut.Status == 2 {
			log.Printf("Creating fs %s on %s", desiredFs, driveName)
			cmd = "mkfs." + desiredFs
			argsCreateFs := []string{
				driveName,
				"-L",
				"KRAKEN-"+label,
			}
			if _, err := ExecuteCommand(cmd, argsCreateFs); err != nil {
				return err
			}
			return nil
		} else {
			log.Printf("%v", err)
			return err
		}
	}
	switch fsOut.Stdout {
	case desiredFs + "\n":
		return nil
	default:
		return fmt.Errorf("Desired fs: %s, actual fs: %s", desiredFs, fsOut.Stdout)
	}
}
