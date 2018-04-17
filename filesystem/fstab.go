package filesystem

import (
	"fmt"
	"os"
)

//AppendToFstab appends a line to fstab with the given parameters.
func AppendToFstab(label string, fs string, mountPoint string) error {
	fstabEntry := fmt.Sprintf("LABEL=%s %s %s defaults 0 1\n", label, mountPoint, fs)

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
