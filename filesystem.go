package main

import (
	"log"
	"fmt"
)

func MountSingleVolume(drive EbsVol) error {
	log.Printf("Mounting single drive: %s", drive.AttachedName)
	return MountSingleDrive(drive.AttachedName, drive.MountPath, drive.FsType)
}


func MountSingleDrive(driveName string, mountPath string, desiredFs string) error {
	if err := checkAndCreateFilesystem(driveName, desiredFs); err != nil {
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

func checkAndCreateFilesystem(driveName string, desiredFs string) error {
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
	if err != nil {
		if fsOut.Status == 2 {
			log.Printf("Creating fs %s on %s", desiredFs, driveName)
			cmd = "mkfs"
			argsCreateFs := []string{
				"-t",
				desiredFs,
				driveName,
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
	case desiredFs+"\n":
		return nil
	default:
		return fmt.Errorf("Desired fs: %s, actual fs: %s", desiredFs, fsOut.Stdout)
	}
}
