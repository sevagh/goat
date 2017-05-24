package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func CreateRaidArray(drives []EbsVol, volName string, dryRun bool) string {
	raidLogger := log.WithFields(log.Fields{"vol_name": volName, "drives": drives})

	raidLogger.Info("Mounting raid drives")

	var raidDriveName string
	var err error
	raidLogger.Info("Searching for unused RAID drive name")
	if raidDriveName, err = RandRaidDriveNamePicker(); err != nil {
		raidLogger.Fatalf("Couldn't select unused RAID drive name: %v", err)
	}

	raidLevel := drives[0].RaidLevel
	cmd := "mdadm"

	driveNames := []string{}
	for _, drive := range drives {
		driveNames = append(driveNames, drive.AttachedName)
	}

	args := []string{
		"--create",
		raidDriveName,
		"--level=" + strconv.Itoa(raidLevel),
		"--name='" + PREFIX + "-" + volName + "'",
		"--raid-devices=" + strconv.Itoa(len(driveNames)),
	}
	args = append(args, driveNames...)
	raidLogger.Infof("RAID: Creating RAID drive: %s %s", cmd, args)
	if dryRun {
		return raidDriveName
	}
	if _, err := ExecuteCommand(cmd, args); err != nil {
		raidLogger.Fatalf("Error when executing mdadm command: %v", err)
	}

	args = []string{
		"--verbose",
		"--detail",
		"--scan",
	}

	raidLogger.Infof("Persisting mdadm settings: %s %s", cmd, args)
	if out, err := ExecuteCommand(cmd, args); err != nil {
		raidLogger.Fatalf("Error when executing mdadm command: %v", err)
	} else {
		if err := appendToMdadmConf(out.Stdout); err != nil {
			raidLogger.Fatalf("Error when persisting mdadm settings to /etc/mdadm.conf: %v", err)
		}
	}

	return raidDriveName
}

func appendToMdadmConf(content string) error {
	f, err := os.OpenFile("/etc/mdadm.conf", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}
	return nil
}
