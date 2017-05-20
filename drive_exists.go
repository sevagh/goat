package main

import (
	"time"
)

const statAttempts = 5

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

func DoesDriveExist(driveName string) bool {
	if _, err := ExecuteCommand("stat", []string{driveName}); err != nil {
		return false
	}
	return true
}

func DoesLabelExist(label string) bool {
	if _, err := ExecuteCommand("ls", []string{"/dev/disk/by-label/" + label}); err != nil {
		return false
	}
	return true
}

func DoesRaidDriveExist(raidDriveName string) bool {
	if _, err := ExecuteCommand("mdadm", []string{raidDriveName}); err != nil {
		return false
	}
	return true
}
