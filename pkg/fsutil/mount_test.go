package fsutil

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestMount(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	if err := Mount("dummy/path", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if !strings.Contains(bufString, "MOUNT: Would have executed: mount [dummy/path]") {
		t.Errorf("logged wrong thing. Actual: %s", bufString)
	}
}
