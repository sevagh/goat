package filesystem

import (
	"path/filepath"
	"strings"

	"github.com/sevagh/goat/execute"
)

//Mount calls mount with no parameters. It relies on there being a correct fstab entry on the provided mountpoint.
func Mount(mountPath string) error {
	cmd := "mount"
	args := []string{
		mountPath,
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
