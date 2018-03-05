package fsutil

import (
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"

	"github.com/sevagh/goat/pkg/execute"
)

//Mount calls mount with no parameters. It relies on there being a correct fstab entry on the provided mountpoint. In the case of a dryRun it doesn't actually execute it, just logs what it would have executed
func Mount(mountPath string, dryRun bool) error {
	cmd := "mount"
	args := []string{
		mountPath,
	}

	if dryRun {
		log.WithFields(log.Fields{"mount_path": mountPath}).Infof("MOUNT: Would have executed: %s %s", cmd, args)
		return nil
	}

	if _, err := execute.Command(cmd, args); err != nil {
		return err
	}

	return nil
}

//IsMountpointAlreadyMounted checks if a mountPoint appears in the output of the mount command. If yes, it returns false. This is to protect from multiple mounts.
func IsMountpointAlreadyMounted(mountPoint string) (bool, error) {
	var mountOut execute.CommandOut
	var err error
	if mountOut, err = execute.Command("mount", []string{}); err != nil {
		return true, err
	}
	for _, line := range strings.Split(mountOut.Stdout, "\n") {
		for _, word := range strings.Split(line, " ") {
			if filepath.Clean(word) == filepath.Clean(mountPoint) {
				return true, nil
			}
		}
	}
	return false, nil
}
