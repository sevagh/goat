package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	isMounted, err := isMountpointAlreadyMounted(mountPath)
	if err != nil {
		return err
	}
	if isMounted {
		return fmt.Errorf("Something already mounted at %s", mountPath)
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

	log.Printf("Appending fstab entry")
	if err := appendFstabEntry("KRAKEN-"+label, desiredFs, mountPath); err != nil {
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
				"KRAKEN-" + label,
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

func appendFstabEntry(label string, fs string, mountPoint string) error {
	fstabEntry := fmt.Sprintf("LABEL=%s %s %s defaults 0 1\n", label, mountPoint, fs)
	log.Printf("Appending to fstab: %s", fstabEntry)
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
