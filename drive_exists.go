package main

import (
	"time"
)

const statAttempts = 5

//DoesDriveExistWithTimeout makes 5 attempts, 1 second sleep between each, to stat a drive to check for its existence
func DoesDriveExistWithTimeout(driveName string) bool {
	var attempts int
	for !DoesDriveExist(driveName) {
		time.Sleep(time.Duration(1 * time.Second))
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

//DoesLabelExist checks /dev/disk/by-label to search for a labelled filesystem
func DoesLabelExist(label string) bool {
	if _, err := ExecuteCommand("ls", []string{"/dev/disk/by-label/" + label}); err != nil {
		return false
	}
	return true
}

//DoesRaidDriveExist uses mdadm to check if a given RAID device exists
func DoesRaidDriveExist(raidDriveName string) bool {
	if _, err := ExecuteCommand("mdadm", []string{raidDriveName}); err != nil {
		return false
	}
	return true
}
