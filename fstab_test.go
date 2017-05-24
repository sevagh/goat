package main

import (
	"testing"
	"github.com/sevagh/stdcap"
	"strings"
)

func TestAppendToFstab(t *testing.T) {
	sc := stdcap.StdoutCapture()

	out := sc.Capture(func() {
		if err := AppendToFstab("test_label", "ext999", "/dummy/path", true); err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	if strings.Contains(out, "FSTAB: would have appended: %sLABEL=test_label /dummy/path ext999 defaults 0 1") {
	    t.Errorf("printed wrong thing to stdout")
	}
}
