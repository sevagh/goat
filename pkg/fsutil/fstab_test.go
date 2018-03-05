package fsutil

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestAppendToFstab(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	if err := AppendToFstab("test_label", "ext999", "/dummy/path", true); err != nil {
		t.Errorf("Error: %v", err)
	}

	bufString := buf.String()
	if !strings.Contains(bufString, "FSTAB: would have appended: LABEL=test_label /dummy/path ext999 defaults 0 1") {
		t.Errorf("logged wrong thing. Actual: %s", bufString)
	}
}
