package main

import (
	"time"
)

const statAttempts = 10

//DoesDriveExistWithTimeout makes 10 attempts, 2 second sleep between each, to stat a drive to check for its existence
func DoesDriveExistWithTimeout(driveName string) bool {
	var attempts int
	for !DoesDriveExist(driveName) {
		time.Sleep(time.Duration(2 * time.Second))
		attempts++
		if attempts >= statAttempts {
			return false
		}
	}
	return true
}

//DoesDriveExist does a single stat call to check if a drive exists
func DoesDriveExist(driveName string) bool {
	if _, err := ExecuteCommand("stat", []string{driveName}); err != nil {
		return false
	}
	return true
}
