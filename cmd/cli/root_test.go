package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Create a new root command instance for testing
	cmd := rootCmd
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Convert photo presets") {
		t.Error("Help text should contain 'Convert photo presets'")
	}
	if !strings.Contains(output, "Supported formats") {
		t.Error("Help text should contain 'Supported formats'")
	}
	if !strings.Contains(output, "Examples") {
		t.Error("Help text should contain 'Examples'")
	}
}

func TestVersionFlag(t *testing.T) {
	// Test that Version field is set correctly in root command
	if rootCmd.Version != "Recipe CLI v0.1.0" {
		t.Errorf("Expected version 'Recipe CLI v0.1.0', got: %s", rootCmd.Version)
	}

	// Test version string variable
	if version != "Recipe CLI v0.1.0" {
		t.Errorf("Expected version variable 'Recipe CLI v0.1.0', got: %s", version)
	}
}

func TestGlobalFlags(t *testing.T) {
	// Create a new root command instance for testing
	cmd := rootCmd
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "--verbose") {
		t.Error("Help text should contain --verbose flag")
	}
	if !strings.Contains(output, "--json") {
		t.Error("Help text should contain --json flag")
	}
}

func TestConvertCommandListed(t *testing.T) {
	// Create a new root command instance for testing
	cmd := rootCmd
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "convert") {
		t.Error("Help text should list convert command")
	}
}
