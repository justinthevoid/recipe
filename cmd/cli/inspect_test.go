package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInspectCommand_ValidFile verifies inspect works with valid files
func TestInspectCommand_ValidFile(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		format string
	}{
		{"NP3 file", "testdata/np3/Classic Chrome.np3", "np3"},
		{"XMP file", "testdata/xmp/AFGA APX 100.xmp", "xmp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run inspect command
			cmd := rootCmd
			cmd.SetArgs([]string{"inspect", tt.file})

			err := cmd.Execute()
			if err != nil {
				t.Errorf("inspect command failed: %v", err)
			}
		})
	}
}

// TestInspectCommand_OutputFlag verifies --output flag
func TestInspectCommand_OutputFlag(t *testing.T) {
	outputFile := "tmp_test_output.json"
	defer os.Remove(outputFile)

	// Run inspect with --output flag
	cmd := rootCmd
	cmd.SetArgs([]string{"inspect", "testdata/np3/Classic Chrome.np3", "--output", outputFile})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("inspect command failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify file has content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Output file is empty")
	}

	// Verify JSON structure
	jsonStr := string(content)
	if !strings.Contains(jsonStr, `"metadata"`) {
		t.Error("Output missing 'metadata' section")
	}
	if !strings.Contains(jsonStr, `"parameters"`) {
		t.Error("Output missing 'parameters' section")
	}
}

// TestInspectCommand_OutputDirectoryCreation verifies parent directory creation
func TestInspectCommand_OutputDirectoryCreation(t *testing.T) {
	outputPath := filepath.Join("tmp_test_dir", "subdir", "output.json")
	defer os.RemoveAll("tmp_test_dir")

	// Run inspect with nested output path
	cmd := rootCmd
	cmd.SetArgs([]string{"inspect", "testdata/np3/Classic Chrome.np3", "--output", outputPath})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("inspect command failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created (parent directory creation failed)")
	}
}

// TestInspectCommand_FileNotFound verifies error for missing file
func TestInspectCommand_FileNotFound(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"inspect", "nonexistent.xmp"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "failed to read file") {
		t.Errorf("Error message should mention file read failure, got: %v", errMsg)
	}
}

// TestInspectCommand_UnknownFormat verifies error for unknown format
func TestInspectCommand_UnknownFormat(t *testing.T) {
	// Create temp file with unknown extension
	tmpFile := "tmp_test.txt"
	os.WriteFile(tmpFile, []byte("test"), 0644)
	defer os.Remove(tmpFile)

	cmd := rootCmd
	cmd.SetArgs([]string{"inspect", tmpFile})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for unknown format, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "unable to detect format") {
		t.Errorf("Error should mention format detection failure, got: %v", errMsg)
	}
	if !strings.Contains(errMsg, "Supported formats") {
		t.Errorf("Error should list supported formats, got: %v", errMsg)
	}
}

// TestInspectCommand_InvalidFileContent verifies parse error handling
func TestInspectCommand_InvalidFileContent(t *testing.T) {
	// Create temp file with invalid NP3 content
	tmpFile := "tmp_invalid.np3"
	os.WriteFile(tmpFile, []byte("invalid content"), 0644)
	defer os.Remove(tmpFile)

	cmd := rootCmd
	cmd.SetArgs([]string{"inspect", tmpFile})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid file content, got nil")
	}

	// Error should be a ConversionError with operation=parse
	errMsg := err.Error()
	if !strings.Contains(errMsg, "parse") {
		t.Errorf("Error should mention parse operation, got: %v", errMsg)
	}
}
