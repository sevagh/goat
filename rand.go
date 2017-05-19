package main

import (
	"fmt"
)

func RandDriveNamePicker() (string, error) {
	ctr := 0
	deviceName := "/dev/xvd"
	runes := []rune("bcdefghijklmnopqrstuvwxyz")
	if DryRun {
		return deviceName + string(runes[0]), nil
	}
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
	if DryRun {
		return deviceName + string(runes[0]), nil
	}
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
