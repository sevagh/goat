package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func MountSingleDrive(driveName string, mountPath string, desiredFs string, label string) error {
	log.Printf("Checking for existing filesystem and/or creating it")
	if err := checkAndCreateFilesystem(driveName, desiredFs, label); err != nil {
		return err
	}
	log.Printf("Checking if something already mounted at %s", mountPath)
	isMounted, err := isMountpointAlreadyMounted(mountPath)
	if err != nil {
		return err
	}
	if isMounted {
		return fmt.Errorf("Something already mounted at %s", mountPath)
	}
	if err := mkdir(mountPath); err != nil {
		return err
	}

	log.Printf("Appending fstab entry")
	if err := appendFstabEntry(PREFIX+"-"+label, desiredFs, mountPath); err != nil {
		return err
	}

	cmd := "mount"
	args := []string{
		mountPath,
	}

	log.Printf("Running mount command")
	if _, err := ExecuteCommand(cmd, args); err != nil {
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
	if _, err := ExecuteCommand(cmd, args); err != nil {
		return err
	}
	return nil
}

func checkAndCreateFilesystem(driveName string, desiredFs string, label string) error {
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
			cmd = "mkfs." + desiredFs
			argsCreateFs := []string{
				driveName,
				"-L",
				PREFIX + "-" + label,
			}
			if _, err := ExecuteCommand(cmd, argsCreateFs); err != nil {
				return err
			}
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

func appendFstabEntry(label string, fs string, mountPoint string) error {
	fstabEntry := fmt.Sprintf("LABEL=%s %s %s defaults 0 1\n", label, mountPoint, fs)
	if DryRun {
		return nil
	}
	f, err := os.OpenFile("/etc/fstab", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(fstabEntry); err != nil {
		return err
	}
	return nil
}

func isMountpointAlreadyMounted(mountPoint string) (bool, error) {
	if mountOut, err := ExecuteCommand("mount", []string{}); err != nil {
		return true, err
	} else {
		for _, line := range strings.Split(mountOut.Stdout, "\n") {
			for _, word := range strings.Split(line, " ") {
				if filepath.Clean(word) == filepath.Clean(mountPoint) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
