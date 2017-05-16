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

func DoesDriveExist(driveName string) bool {
	log.Printf("Checking if device %s exists", driveName)
	if _, err := ExecuteCommand("stat", []string{driveName}); err != nil {
		log.Printf("%s doesn't exist", driveName)
		return false
	}
	log.Printf("%s exists", driveName)
	return true
}

func ExecuteCommand(commandString string, args []string) (CommandOut, error) {
	out := CommandOut{}
	cmd := exec.Command(commandString, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	dryRunPrefix := ""
	if DryRun {
		dryRunPrefix = "[DRY]: "
	}

	log.Printf("%sCmd args: %s", dryRunPrefix, cmd.Args)

	if DryRun {
		out.Status = 0
		return out, nil
	}

	if err := cmd.Start(); err != nil {
		return out, fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
				out.Stdout = cmdOut.String()
				out.Stderr = cmdErr.String()
				out.Status = status.ExitStatus()
				return out, fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
			out.Stdout = cmdOut.String()
			out.Stderr = cmdErr.String()
			return out, fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	log.Printf("OUT: %s, ERR: %s", cmdOut.String(), cmdErr.String())
	out.Stdout = cmdOut.String()
	out.Stderr = cmdErr.String()
	out.Status = 0
	return out, nil
}
