package main

import (
	"testing"
	"strings"
	"bytes"
	log "github.com/sirupsen/logrus"
)

func TestCheckFilesystem(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	if err := CheckFilesystem("/dev/dummy", "ext999", "dummy_label", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if ! strings.Contains(bufString, "FILESYSTEM: would have executed blkid [-o value -s TYPE /dev/dummy]") {
	    t.Errorf("printed wrong thing to stderr. Actual: %s", bufString)
	}
}
