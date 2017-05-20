package main

import (
	"path/filepath"
	"strings"
)

func Mount(mountPath string) error {
	cmd := "mount"
	args := []string{
		mountPath,
	}

	if _, err := ExecuteCommand(cmd, args); err != nil {
		return err
	}

	return nil
}

func IsMountpointAlreadyMounted(mountPoint string) (bool, error) {
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
