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
	if err := checkFilesystem(driveName, desiredFs, logger); err != nil {
		return err
	}
	cmd := "mount"
	args := []string{
		driveName,
		mountPath,
	}
	logger.Printf("Executing: %s %s\n", cmd, args)
	if _, err := ExecuteCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}

	return nil
}

func checkFilesystem(driveName string, desiredFs string, logger *log.Logger) error {
	logger.Printf("Checking filesystem on %s", driveName)
        cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		driveName,
	}
	fsOut, err := ExecuteCommand(cmd, args, logger);
	if err != nil {
		logger.Printf("%v", err)
		return err
	}
	if fsOut.Status == 2 {
		return nil
	}
	switch fsOut.Stdout {
	case desiredFs:
		return nil
	default:
		return fmt.Errorf("Desired fs: %s, actual fs: %s", desiredFs, fsOut.Stdout)
	}
}
