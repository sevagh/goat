package fsutil

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/sevagh/goat/execute"
)

//CheckFilesystem checks for a filesystem on a given drive using blkid. It returns ok if there is no filesystem or the filesystem is the correct type. Error if there's a different filesystem
func CheckFilesystem(driveName string, desiredFs string, label string, dryRun bool) error {
	cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		driveName,
	}

	if dryRun {
		log.WithFields(log.Fields{"drive_name": driveName, "fs": desiredFs, "label": label}).Infof("FILESYSTEM: would have executed %s %s", cmd, args)
		return nil
	}

	fsOut, err := execute.Command(cmd, args)
	if err != nil {
		if fsOut.Status == 2 {
			//go ahead and create filesystem
			return nil
		}
		return err
	}
	switch fsOut.Stdout {
	case desiredFs + "\n":
		return nil
	default:
		return fmt.Errorf("Desired fs: %s, actual fs: %s", desiredFs, fsOut.Stdout)
	}
}

//CreateFilesystem executes mkfs.<desired_filesystem> on the requested drive. It can do a dryRun where it logs its command without running it
func CreateFilesystem(driveName string, desiredFs string, label string, dryRun bool) error {
	cmd := "mkfs." + desiredFs
	args := []string{
		driveName,
		"-L",
		"GOAT-" + label,
	}

	if dryRun {
		log.WithFields(log.Fields{"drive_name": driveName, "fs": desiredFs, "label": label}).Infof("FILESYSTEM: would have executed %s %s", cmd, args)
		return nil
	}

	if _, err := execute.Command(cmd, args); err != nil {
		return err
	}
	return nil
}
