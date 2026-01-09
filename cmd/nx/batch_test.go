package main

import (
	"testing"
)

func TestNewBatchCmd(t *testing.T) {
	cmd := newBatchCmd()
	if cmd.Use != "batch" {
		t.Errorf("expected Use to be 'batch', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}
