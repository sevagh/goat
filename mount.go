package main

import (
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func Mount(mountPath string, dryRun bool) error {
	cmd := "mount"
	args := []string{
		mountPath,
	}

	if dryRun {
		log.WithFields(log.Fields{"mount_path": mountPath}).Infof("MOUNT: Would have executed: %s %s", cmd, args)
		return nil
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
