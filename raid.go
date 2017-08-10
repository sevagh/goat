package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"

	"github.com/sevagh/goat/execute"
)

//CreateRaidArray runs the appropriate mdadm command for the given list of EbsVol that should be raided together. It takes dryRun as a boolean, where it tells you which mdadm it would have run
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

	nameString := "--name='GOAT-" + volName + "'"

	var args []string
	args = []string{
		"--create",
		raidDriveName,
		"--level=" + strconv.Itoa(raidLevel),
		nameString,
		"--raid-devices=" + strconv.Itoa(len(driveNames)),
	}

	args = append(args, driveNames...)
	raidLogger.Infof("RAID: Creating RAID drive: %s %s", cmd, args)
	if dryRun {
		return raidDriveName
	}
	if _, err := execute.ExecuteCommand(cmd, args); err != nil {
		raidLogger.Fatalf("Error when executing mdadm command: %v", err)
	}

	return raidDriveName
}

//PersistMdadm dumps the current mdadm config to /etc/mdadm.conf
func PersistMdadm() error {
	cmd := "mdadm"

	args := []string{
		"--verbose",
		"--detail",
		"--scan",
	}

	log.Infof("Persisting mdadm settings: %s %s", cmd, args)

	var out execute.CommandOut
	var err error
	if out, err = execute.ExecuteCommand(cmd, args); err != nil {
		log.Fatalf("Error when executing mdadm command: %v", err)
	}

	f, err := os.OpenFile("/etc/mdadm.conf", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(out.Stdout); err != nil {
		return err
	}
	return nil
}
