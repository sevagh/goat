package execute

import (
	"testing"
)

func TestExecuteCommandReturnCodeSuccess(t *testing.T) {
	if out, err := ExecuteCommand("true", []string{}); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if !(out.Status == 0) {
			t.Fatalf("Expected exit code: 0")
		}
	}
}

func TestExecuteCommandReturnCodeFailure(t *testing.T) {
	if out, err := ExecuteCommand("false", []string{}); err == nil {
		t.Fatalf("Expected an error")
	} else {
		if !(out.Status == 1) {
			t.Fatal("Expected exit code: 1")
		}
	}
}

func TestExecuteCommandStdout(t *testing.T) {
	if out, err := ExecuteCommand("bash", []string{"-c", "echo hello"}); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if !(out.Stderr == "") {
			t.Fatal("Expected stdout to be empty")
		}
		if !(out.Stdout == "hello\n") {
			t.Fatalf("Expected hello to be printed to stdout. Actual: %s", out.Stdout)
		}
	}
}

func TestExecuteCommandStderr(t *testing.T) {
	if out, err := ExecuteCommand("bash", []string{"-c", "echo hello 1>&2"}); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if !(out.Stdout == "") {
			t.Fatal("Expected stdout to be empty")
		}
		if !(out.Stderr == "hello\n") {
			t.Fatalf("Expected hello to be printed to stderr. Actual: %s", out.Stderr)
		}
	}
}

func TestExecuteCommandFailedToStart(t *testing.T) {
	if out, err := ExecuteCommand("i_shouldnt_exist", []string{}); err == nil {
		t.Fatal("Expected error")
	} else {
		if !(out.Status == 0) {
			t.Fatal("Return code should be uninitialized i.e. 0")
		}
		if !(out.Stderr == "" && out.Stdout == "") {
			t.Fatal("Both stdout and stderr should be empty")
		}
	}
}
