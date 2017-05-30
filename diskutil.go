package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

//PrepAndMountDrives prepares the filesystem, RAIDs (if necessary) and mounts a given list of EbsVol (can be size 1 for non-RAID)
func PrepAndMountDrives(volName string, vols []EbsVol, ec2Instance EC2Instance, dryRun bool) {
	driveLogger := log.WithFields(log.Fields{"vol_name": volName, "vols": vols})

	mountPath := vols[0].MountPath
	desiredFs := vols[0].FsType

	if DoesLabelExist(PREFIX + "-" + volName) {
		driveLogger.Info("Label already exists, jumping to mount phase")
	} else {
		var driveName string
		if len(vols) == 1 {
			driveLogger.Info("Single drive, no RAID")
			driveName = vols[0].AttachedName
		} else {
			driveLogger.Info("Creating RAID array")
			driveName = CreateRaidArray(vols, volName, dryRun)
		}

		driveLogger.Info("Checking for existing filesystem")

		if err := CheckFilesystem(driveName, desiredFs, volName, dryRun); err != nil {
			driveLogger.Fatalf("Checking for existing filesystem: %v", err)
		}
		if err := CreateFilesystem(driveName, desiredFs, volName, dryRun); err != nil {
			driveLogger.Fatalf("Error when creating filesystem: %v", err)
		}
	}

	driveLogger.Info("Checking if something already mounted at %s", mountPath)
	if isMounted, err := IsMountpointAlreadyMounted(mountPath); err != nil {
		driveLogger.Fatalf("Error when checking mount point for existing mounts: %v", err)
	} else {
		if isMounted {
			driveLogger.Fatalf("Something already mounted at %s", mountPath)
		}
	}

	if !dryRun {
		if err := os.MkdirAll(mountPath, 0777); err != nil {
			driveLogger.Fatalf("Couldn't mkdir: %v", err)
		}
	}

	driveLogger.Info("Appending fstab entry")
	if err := AppendToFstab(PREFIX+"-"+volName, desiredFs, mountPath, dryRun); err != nil {
		driveLogger.Fatalf("Couldn't append to fstab: %v", err)
	}

	driveLogger.Info("Now mounting")
	if err := Mount(mountPath, dryRun); err != nil {
		driveLogger.Fatalf("Couldn't mount: %v", err)
	}
}
