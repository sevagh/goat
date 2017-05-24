package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

//AppendToFstab appends a line to fstab with the given parameters. It takes dryRun as a param, where it says what it would have appended without actually modifying '/etc/fstab'
func AppendToFstab(label string, fs string, mountPoint string, dryRun bool) error {
	fstabEntry := fmt.Sprintf("LABEL=%s %s %s defaults 0 1\n", label, mountPoint, fs)
	if dryRun {
		log.WithFields(log.Fields{"label": label, "fs": fs, "mount_point": mountPoint}).Infof("FSTAB: would have appended: %s", fstabEntry)
		return nil
	}

	f, err := os.OpenFile("/etc/fstab", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(fstabEntry); err != nil {
		return err
	}
	return nil
}
