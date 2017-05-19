package main

import (
	"log"
	"time"
)

const statAttempts = 5

func DoDrivesExist(drives []EbsVol) bool {
	for _, drive := range drives {
		log.Printf("Checking if drive %s exists", drive.AttachedName)
		var attempts int
		for !DoesDriveExist(drive.AttachedName) {
			time.Sleep(time.Duration(1 * time.Second))
			attempts++
			if attempts >= statAttempts {
				log.Printf("Exceeded max (%d) stat attempts waiting for drive %s to exist", statAttempts, drive.AttachedName)
				return false
			}
		}
	}
	return true
}

func DoesDriveExist(driveName string) bool {
	log.Printf("Checking if device %s exists", driveName)
	if _, err := ExecuteCommand("stat", []string{driveName}); err != nil {
		log.Printf("%s doesn't exist", driveName)
		return false
	}
	log.Printf("%s exists", driveName)
	return true
}

func DoesLabelExist(label string) bool {
	log.Println("Checking if label exists")
	if _, err := ExecuteCommand("ls", []string{"/dev/disk/by-label" + label}); err != nil {
		log.Printf("%s doesn't exist", label)
		return false
	}
	log.Printf("%s exists", label)
	return true
}
