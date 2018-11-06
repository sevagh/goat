package main

import (
	"testing"
)

func TestRunCli(t *testing.T) {
	c := runCli()
	if c.Name != "go-ebsnvme" {
		t.Fatalf("Expected c.Name to be go-ebsnvme, got '%v'", c.Name)
	}
}
