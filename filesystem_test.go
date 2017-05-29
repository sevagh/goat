package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestCheckFilesystem(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	if err := CheckFilesystem("/dev/dummy", "ext999", "dummy_label", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if !strings.Contains(bufString, "FILESYSTEM: would have executed blkid [-o value -s TYPE /dev/dummy]") {
		t.Errorf("logged wrong thing. Actual: %s", bufString)
	}
}

func TestCreateFilesystem(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	if err := CreateFilesystem("/dev/dummy", "ext999", "dummy_label", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if !strings.Contains(bufString, "FILESYSTEM: would have executed mkfs.ext999 [/dev/dummy -L EWIZ-dummy_label]") {
		t.Errorf("logged wrong thing. Actual: %s", bufString)
	}
}
