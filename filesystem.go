package main

import (
	"log"
	"fmt"
)

func MountSingleVolume(drive EbsVol, logger *log.Logger) error {
	logger.Printf("Mounting single drive: %s", drive.AttachedName)
	return MountSingleDrive(drive.AttachedName, drive.MountPath, drive.FsType, logger)
}


func MountSingleDrive(driveName string, mountPath string, desiredFs string, logger *log.Logger) error {
	if err := checkAndCreateFilesystem(driveName, desiredFs, logger); err != nil {
		return err
	}
	if err := mkdir(mountPath, logger); err != nil {
		return err
	}
	cmd := "mount"
	args := []string{
		driveName,
		mountPath,
	}
	logger.Printf("Executing: %s %s", cmd, args)
	if _, err := ExecuteCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}

	return nil
}

func mkdir(mountPath string, logger *log.Logger) error {
	cmd := "mkdir"
	args := []string{
		"-p",
		mountPath,
	}
	logger.Printf("Executing: %s %s", cmd, args)
	if _, err := ExecuteCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}
	return nil
}

func checkAndCreateFilesystem(driveName string, desiredFs string, logger *log.Logger) error {
	logger.Printf("Checking filesystem on %s", driveName)
        cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		driveName,
	}
	fsOut, err := ExecuteCommand(cmd, args, logger)
	if err != nil {
		if fsOut.Status == 2 {
			logger.Printf("Creating fs %s on %s", desiredFs, driveName)
			cmd = "mkfs"
			argsCreateFs := []string{
				"-t",
				desiredFs,
				driveName,
			}
			if _, err := ExecuteCommand(cmd, argsCreateFs, logger); err != nil {
				return err
			}
			return nil
		} else {
		    logger.Printf("%v", err)
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
