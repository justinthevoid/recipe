package main

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Attempt to create the root command.
	// This function doesn't exist yet, so compilation will fail (Red phase).
	cmd := newRootCmd()

	if cmd == nil {
		t.Fatal("Root command must not be nil")
	}

	if cmd.Use != "recipe-nx" {
		t.Errorf("Expected Use='recipe-nx', got '%s'", cmd.Use)
	}

	// Verify command hierarchy: recipe-nx -> batch -> apply
	batchCmd, _, err := cmd.Find([]string{"batch"})
	if err != nil || batchCmd.Use != "batch" {
		t.Fatal("Start-up failed: 'batch' command not found")
	}

	applyCmd, _, err := batchCmd.Find([]string{"apply"})
	if err != nil || applyCmd.Use != "apply" {
		t.Fatal("Hierarchy error: 'apply' subcommand not found under 'batch'")
	}

	// Verify log-level flag
	flag := cmd.PersistentFlags().Lookup("log-level")
	if flag == nil {
		t.Error("Flag 'log-level' not found on root command")
	}
}
