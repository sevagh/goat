package filesystem

import (
	"fmt"
	"os"
	"strconv"
)

//CreateRaidArray runs the appropriate mdadm command for the given list of EbsVol that should be raided together.
func CreateRaidArray(driveNames []string, volName string, raidLevel int) (string, error) {
	var raidDriveName string
	var err error

	if raidDriveName, err = randRaidDriveNamePicker(); err != nil {
		return "", fmt.Errorf("Couldn't select unused RAID drive name: %v", err)
	}

	cmd := "mdadm"

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

	if _, err := Command(cmd, args, ""); err != nil {
		return "", fmt.Errorf("Error when executing mdadm command: %v", err)
	}

	return raidDriveName, nil
}

//PersistMdadm dumps the current mdadm config to /etc/mdadm.conf
func PersistMdadm() error {
	cmd := "mdadm"

	args := []string{
		"--verbose",
		"--detail",
		"--scan",
	}

	var out CommandOut
	var err error
	if out, err = Command(cmd, args, ""); err != nil {
		return fmt.Errorf("Error when executing mdadm command: %v", err)
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

func randRaidDriveNamePicker() (string, error) {
	ctr := 0
	deviceName := "/dev/md"
	runes := []rune("0123456789")
	for {
		if ctr >= len(runes) {
			return "", fmt.Errorf("Ran out of raid drive names")
		}
		if !DoesDriveExist(deviceName + string(runes[ctr])) {
			break
		}
		ctr++
	}
	return deviceName + string(runes[ctr]), nil
}
