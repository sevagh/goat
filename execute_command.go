package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

type CommandOut struct {
	Stdout string
	Stderr string
	Status int
}

func ExecuteCommand(commandString string, args []string) (CommandOut, error) {
	out := CommandOut{}
	cmd := exec.Command(commandString, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Start(); err != nil {
		return out, fmt.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				out.Stdout = cmdOut.String()
				out.Stderr = cmdErr.String()
				out.Status = status.ExitStatus()
				return out, fmt.Errorf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			out.Stdout = cmdOut.String()
			out.Stderr = cmdErr.String()
			return out, fmt.Errorf("cmd.Wait: %v", err)
		}
	}
	out.Stdout = cmdOut.String()
	out.Stderr = cmdErr.String()
	out.Status = 0
	return out, nil
}
