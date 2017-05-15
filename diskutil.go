package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

const statAttempts = 5

func MountSingleDrive(drive EbsVol, logger *log.Logger) error {
	logger.Printf("Mounting single drive: %s", drive.AttachedName)
	return mountSingleDrive(drive.AttachedName, drive.MountPath, logger)
}

func DoesDriveExist(driveName string, logger *log.Logger) bool {
	logger.Printf("Checking if device %s exists", driveName)
	if err := executeCommand("stat", []string{driveName}, logger); err != nil {
		logger.Printf("%s doesn't exist", driveName)
		return false
	}
	logger.Printf("%s exists", driveName)
	return true
}

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
	if err := executeCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}

	return mountSingleDrive(raidDriveName, mountPath, logger)
}

func mountSingleDrive(driveName string, mountPath string, logger *log.Logger) error {
	cmd := "mount"
	args := []string{
		driveName,
		mountPath,
	}
	logger.Printf("Executing: %s %s\n", cmd, args)
	if err := executeCommand(cmd, args, logger); err != nil {
		logger.Printf("%v", err)
		return err
	}

	return nil
}

func executeCommand(commandString string, args []string, logger *log.Logger) error {
	cmd := exec.Command(commandString, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	logger.Printf("Cmd args: %s", cmd.Args)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
				return fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
			return fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
	return nil
}
