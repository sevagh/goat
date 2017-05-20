package main

import (
	"fmt"
	"strings"
)

func RandDriveNamePicker() (string, error) {
	ctr := 0
	deviceName := "/dev/xvd"
	runes := []rune("bcdefghijklmnopqrstuvwxyz")
	for {
		if ctr >= len(runes) {
			return "", fmt.Errorf("Ran out of drive names")
		}
		if !DoesDriveExist(deviceName + string(runes[ctr])) {
			break
		}
		ctr++
	}
	return deviceName + string(runes[ctr]), nil
}

func RandRaidDriveNamePicker() (string, error) {
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

func DrawAsciiBanner(headLine string) string {
	return fmt.Sprintf("\n%[1]s\n# %[2]s #\n%[1]s\n",
		strings.Repeat("#", len(headLine)+4),
		headLine)
}
