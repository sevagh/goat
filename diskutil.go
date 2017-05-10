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

func MountSingleDrive(driveName string, mountPath string, logger *log.Logger) error {
	logger.Printf("Mounting single drive: %s", driveName)
	return nil
}

func MountRaidDrives(driveNames []string, mountPath string, raidLevel int, logger *log.Logger) error {
	logger.Printf("Mounting raid drives: %s", driveNames)
	if raidLevel != 0 && raidLevel != 1 {
		return fmt.Errorf("Valid raid levels are 0 and 1")
	}
	logger.Printf("Checking if drives exist")
	for _, driveName := range driveNames {
		var attempts int
		for err := fmt.Errorf("dummy_error"); err != nil; err = executeCommand("stat", []string{driveName}, logger) {
			time.Sleep(time.Duration(1 * time.Second))
			attempts++
			if attempts >= statAttempts {
				logger.Fatalf("Exceeded max (%d) stat attempts waiting for drive %s to exist", statAttempts, driveName)
				return fmt.Errorf("Stat failed")
			}
		}
	}

	cmd := "mdadm"
	args := []string{
		"--create",
		"/dev/md0",
		"--level=" + strconv.Itoa(raidLevel),
		"--name=KRAKEN",
		"--raid-devices=" + strconv.Itoa(len(driveNames)),
	}
	args = append(args, driveNames...)
	logger.Printf("Executing: %s %s\n", cmd, args)
	if err := executeCommand(cmd, args, logger); err != nil {
		logger.Fatalf("%v", err)
		return err
	}

	return nil
}

func mountSingleDrive(driveName string, mountPath string) error {
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
				logger.Fatalf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
				return fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			logger.Fatalf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
			return fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
	return nil
}
