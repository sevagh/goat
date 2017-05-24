package main

import (
	"testing"
	"strings"
)

func TestRandDriveNamePicker(t *testing.T) {
	if name, err := RandDriveNamePicker(); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if ! strings.Contains(name, "/dev/xvdb") {
			t.Fatalf("Expected drive name in format /dev/xvd*, got: %s", name)
		}
	}
}

func TestRandRaidDriveNamePicker(t *testing.T) {
	if name, err := RandRaidDriveNamePicker(); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if ! strings.Contains(name, "/dev/md") {
			t.Fatalf("Expected drive name in format /dev/md*, got: %s", name)
		}
	}
}
