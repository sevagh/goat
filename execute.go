package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

type CommandOut struct {
	Stdout    string
	Stderr    string
	Status int
}

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

func ExecuteCommand(commandString string, args []string, logger *log.Logger) (CommandOut, error) {
	out := CommandOut{}
	cmd := exec.Command(commandString, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	logger.Printf("Cmd args: %s", cmd.Args)

	if err := cmd.Start(); err != nil {
		return out, fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
				out.Stdout = cmdOut.String()
				out.Stderr = cmdErr.String()
				out.Status = status.ExitStatus()
				return out, fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
			out.Stdout = cmdOut.String()
			out.Stderr = cmdErr.String()
			return out, fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	logger.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
	out.Stdout = cmdOut.String()
	out.Stderr = cmdErr.String()
	out.Status = 0
	return out, nil
}
