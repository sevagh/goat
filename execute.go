package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

const statAttempts = 5

func DoesDriveExist(driveName string, logger *log.Logger) bool {
	logger.Printf("Checking if device %s exists", driveName)
	if _, err := ExecuteCommand("stat", []string{driveName}, logger); err != nil {
		logger.Printf("%s doesn't exist", driveName)
		return false
	}
	logger.Printf("%s exists", driveName)
	return true
}

func ExecuteCommand(commandString string, args []string, logger *log.Logger) (string, error) {
	cmd := exec.Command(commandString, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	logger.Printf("Cmd args: %s", cmd.Args)

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
				return "", fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
			return "", fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
	return cmdOut.String(), nil
}
