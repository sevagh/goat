package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

func MountSingleDrive(driveName string, mountPath string, logger *log.Logger) error {
	logger.Printf("Mounting single drive: %s", driveName)
	return nil
}

func MountRaidDrives(driveNames []string, mountPath string, raidLevel int, logger *log.Logger) error {
	logger.Printf("Mounting raid drives: %s", driveNames)
	if raidLevel != 0 && raidLevel != 1 {
		return fmt.Errorf("Valid raid levels are 0 and 1")
	}
	driveString := strings.Join(driveNames, " ")
	cmd := "mdadm"
	args := []string{
		"--create",
		"/dev/md0",
		"--level=" + string(raidLevel),
		"--name=KRAKEN",
		"--raid-devices=" + string(len(driveNames)),
		driveString,
	}
	logger.Printf("Executing: %s %s\n", cmd, args)
	if err := executeCommand(cmd, args); err != nil {
		logger.Fatalf("%v", err)
		return err
	}

	return nil
}

func mountSingleDrive(driveName string, mountPath string) error {
	return nil
}

func executeCommand(commandString string, args []string) error {
	cmd := exec.Command(commandString, args...)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			return fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	return nil
}
