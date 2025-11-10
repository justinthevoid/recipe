package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestConvertJSONOutput tests single file conversion with JSON output (AC-2)
func TestConvertJSONOutput(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Use existing test file - use absolute path from project root
	inputFile := filepath.Join("..", "..", "testdata", "xmp", "sample.xmp")

	cmd := exec.Command("./recipe-test.exe", "convert", inputFile,
		"--to", "np3", "--json", "--overwrite")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("convert with --json failed: %v\nOutput: %s", err, output)
	}

	// Parse JSON output
	var result ConversionResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Verify required fields (AC-2)
	if result.Input == "" {
		t.Error("JSON output missing input field")
	}
	if result.Output == "" {
		t.Error("JSON output missing output field")
	}
	if result.SourceFormat == "" {
		t.Error("JSON output missing source_format field")
	}
	if result.TargetFormat != "np3" {
		t.Errorf("target_format = %s, want np3", result.TargetFormat)
	}
	if !result.Success {
		t.Errorf("success = false, want true (error: %s)", result.Error)
	}
	if result.DurationMs <= 0 {
		t.Errorf("duration_ms = %d, want > 0", result.DurationMs)
	}
	if result.FileSizeBytes <= 0 {
		t.Errorf("file_size_bytes = %d, want > 0", result.FileSizeBytes)
	}
}

// TestConvertJSONOutput_FailedConversion tests JSON output for errors (AC-3)
func TestConvertJSONOutput_FailedConversion(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Create a corrupted file
	tmpDir := t.TempDir()
	corruptedFile := filepath.Join(tmpDir, "corrupted.xmp")

	// Write invalid XML
	if err := writeTestFile(corruptedFile, []byte("not valid xml")); err != nil {
		t.Fatalf("failed to create corrupted file: %v", err)
	}

	cmd := exec.Command("./recipe-test.exe", "convert", corruptedFile,
		"--to", "np3", "--json")

	// Capture stdout only (JSON goes there)
	var stdout strings.Builder
	cmd.Stdout = &stdout
	err := cmd.Run()

	// Should exit with error code (AC-3)
	if err == nil {
		t.Error("expected non-zero exit code for failed conversion")
	}

	// But should still output valid JSON to stdout (AC-3)
	var result ConversionResult
	if err := json.Unmarshal([]byte(stdout.String()), &result); err != nil {
		t.Fatalf("failed to parse JSON output for error case: %v\nOutput: %s", err, stdout.String())
	}

	// Verify error fields (AC-3)
	if result.Success {
		t.Error("success should be false for failed conversion")
	}
	if result.Error == "" {
		t.Error("error field should not be empty for failed conversion")
	}
}

// TestBatchJSONOutput tests batch conversion with JSON output (AC-4)
func TestBatchJSONOutput(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Create test files
	tmpDir := t.TempDir()
	createTestFiles(t, tmpDir, 3)

	cmd := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("batch with --json failed: %v\nOutput: %s", err, output)
	}

	// Parse JSON output
	var result BatchResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("failed to parse batch JSON output: %v\nOutput: %s", err, output)
	}

	// Verify batch fields (AC-4)
	if !result.Batch {
		t.Error("batch field should be true")
	}
	if result.Total != 3 {
		t.Errorf("total = %d, want 3", result.Total)
	}
	if result.SuccessCount != 3 {
		t.Errorf("success_count = %d, want 3", result.SuccessCount)
	}
	if result.ErrorCount != 0 {
		t.Errorf("error_count = %d, want 0", result.ErrorCount)
	}
	if result.DurationMs <= 0 {
		t.Errorf("duration_ms = %d, want > 0", result.DurationMs)
	}
	if len(result.Results) != 3 {
		t.Errorf("len(results) = %d, want 3", len(result.Results))
	}

	// Verify each result follows single conversion schema (AC-4)
	for i, r := range result.Results {
		if r.Input == "" {
			t.Errorf("results[%d].input is empty", i)
		}
		if r.Output == "" {
			t.Errorf("results[%d].output is empty", i)
		}
		if r.SourceFormat == "" {
			t.Errorf("results[%d].source_format is empty", i)
		}
		if r.TargetFormat != "np3" {
			t.Errorf("results[%d].target_format = %s, want np3", i, r.TargetFormat)
		}
		if !r.Success {
			t.Errorf("results[%d].success = false (error: %s)", i, r.Error)
		}
	}
}

// TestJSONWithVerbose tests JSON output with verbose flag (AC-7)
func TestJSONWithVerbose(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	inputFile := filepath.Join("..", "..", "testdata", "xmp", "sample.xmp")

	cmd := exec.Command("./recipe-test.exe", "convert", inputFile,
		"--to", "np3", "--json", "--verbose", "--overwrite")
	cmd.Dir = "."

	// Capture stdout and stderr separately
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start command: %v", err)
	}

	// Read outputs
	stdoutData := make([]byte, 10240)
	stderrData := make([]byte, 10240)

	n1, _ := stdout.Read(stdoutData)
	n2, _ := stderr.Read(stderrData)

	cmd.Wait()

	stdoutStr := string(stdoutData[:n1])
	stderrStr := string(stderrData[:n2])

	// Verify stdout contains ONLY valid JSON (AC-7)
	var result ConversionResult
	if err := json.Unmarshal([]byte(stdoutStr), &result); err != nil {
		t.Errorf("stdout should contain valid JSON: %v\nStdout: %s", err, stdoutStr)
	}

	// Verify stdout does NOT contain log messages (AC-7)
	if strings.Contains(stdoutStr, "DEBUG") || strings.Contains(stdoutStr, "INFO") {
		t.Errorf("stdout should not contain log messages, got: %s", stdoutStr)
	}

	// Verify stderr contains verbose logs (AC-7)
	if !strings.Contains(stderrStr, "DEBUG") && !strings.Contains(stderrStr, "INFO") {
		t.Errorf("stderr should contain verbose logs, got: %s", stderrStr)
	}
}

// Helper function for writing test files
func writeTestFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}

// TestJSONFieldNamingConvention tests snake_case field names (AC-5)
func TestJSONFieldNamingConvention(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	inputFile := filepath.Join("..", "..", "testdata", "xmp", "sample.xmp")

	cmd := exec.Command("./recipe-test.exe", "convert", inputFile,
		"--to", "np3", "--json", "--overwrite")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("convert with --json failed: %v", err)
	}

	outputStr := string(output)

	// Verify snake_case field names (AC-5)
	expectedFields := []string{
		"\"input\"",
		"\"output\"",
		"\"source_format\"",
		"\"target_format\"",
		"\"success\"",
		"\"duration_ms\"",
		"\"file_size_bytes\"",
	}

	for _, field := range expectedFields {
		if !strings.Contains(outputStr, field) {
			t.Errorf("JSON should contain snake_case field %s", field)
		}
	}

	// Verify NO camelCase field names
	unwantedFields := []string{
		"\"sourceFormat\"",
		"\"targetFormat\"",
		"\"durationMs\"",
		"\"fileSizeBytes\"",
	}

	for _, field := range unwantedFields {
		if strings.Contains(outputStr, field) {
			t.Errorf("JSON should NOT contain camelCase field %s", field)
		}
	}
}
