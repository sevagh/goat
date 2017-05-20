package main

import (
	"fmt"
)

func CheckFilesystem(driveName string, desiredFs string, label string, dryRun bool) error {
	cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		driveName,
	}
	fsOut, err := ExecuteCommand(cmd, args)
	if dryRun {
		return nil
	}
	if err != nil {
		if fsOut.Status == 2 {
			//go ahead and create filesystem
			return nil
		} else {
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

func CreateFilesystem(driveName string, desiredFs string, label string) error {
	cmd := "mkfs." + desiredFs
	args := []string{
		driveName,
		"-L",
		PREFIX + "-" + label,
	}
	if _, err := ExecuteCommand(cmd, args); err != nil {
		return err
	}
	return nil
}
