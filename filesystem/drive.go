package filesystem

import (
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/mvisonneau/go-ebsnvme/pkg/ebsnvme"
)

//DoesDriveExistWithTimeout makes 10 attempts, 2 second sleep between each, to stat a drive to check for its existence
func DoesDriveExistWithTimeout(driveName string, maxAttempts int) bool {
	var attempts int
	for !DoesDriveExist(driveName) {
		time.Sleep(time.Duration(2 * time.Second))
		attempts++
		if attempts >= maxAttempts {
			return false
		}
	}
	return true
}

//DoesDriveExist does a single stat call to check if a drive exists
func DoesDriveExist(driveName string) bool {
	if _, err := os.Stat(driveName); os.IsNotExist(err) {
		for _, file := range listNvmeBlockDevices() {
			if d, _ := ebsnvme.ScanDevice(file); d.Name == driveName {
				return true
			}
		}
		return false
	}
	return true
}

// GetLocalBlockDeviceName returns the actual name of the block device seen
// within the instance (useful for nitros)
func GetActualBlockDeviceName(name string) (string, error) {
	for _, device := range listNvmeBlockDevices() {
		if d, _ := ebsnvme.ScanDevice(device); d.Name == name {
			return device, nil
		}
	}
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return "", err
	}
	return name, nil
}

func listNvmeBlockDevices() (devices []string) {
	re := regexp.MustCompile("(^\\/dev\\/nvme[0-9]+n1$)")
	f, _ := filepath.Glob("/dev/nvme*")
	for _, d := range f {
		if re.Match([]byte(d)) {
			devices = append(devices, d)
		}
	}
	return
}
