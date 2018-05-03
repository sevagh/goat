package filesystem

import (
	"time"

	"github.com/sevagh/goat/execute"
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
	if _, err := execute.Command("stat", []string{driveName}, ""); err != nil {
		return false
	}
	return true
}
