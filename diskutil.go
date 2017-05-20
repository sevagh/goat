package main

import (
	log "github.com/sirupsen/logrus"
)

func PrepAndMountDrives(volName string, vols []EbsVol, dryRun bool) {
	driveLogger := log.WithFields(log.Fields{"vol_name": volName, "vols": vols})
	var driveName string
	if len(vols) == 1 {
		driveLogger.Info("Single drive, no RAID")
		driveName = vols[0].AttachedName
	} else {
		driveLogger.Info("Creating RAID array")
		driveName = CreateRaidArray(vols, volName, dryRun)
	}

	mountPath := vols[0].MountPath
	desiredFs := vols[0].FsType

	driveLogger.Info("Checking for existing filesystem")
	if err := CheckFilesystem(driveName, desiredFs, volName, dryRun); err != nil {
		driveLogger.Fatalf("Checking for existing filesystem: %v", err)
	}

	driveLogger.Info("Checking if something already mounted at %s", mountPath)
	if isMounted, err := IsMountpointAlreadyMounted(mountPath); err != nil {
		driveLogger.Fatalf("Checking mount point for existing mounts: %v", err)
	} else {
		if isMounted {
			driveLogger.Fatalf("Something already mounted at %s", mountPath)
		}
	}

	if err := Mkdir(mountPath); err != nil {
		driveLogger.Fatalf("Couldn't mkdir: %v", err)
	}

	driveLogger.Info("Appending fstab entry")
	if err := AppendToFstab(PREFIX+"-"+volName, desiredFs, mountPath, dryRun); err != nil {
		driveLogger.Fatalf("Couldn't append to fstab: %v", err)
	}

	driveLogger.Info("Now mounting")
	if err := Mount(mountPath); err != nil {
		driveLogger.Fatalf("Couldn't mount: %v", err)
	}
}
