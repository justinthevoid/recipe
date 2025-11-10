package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestVerboseFlag_ShortFlag tests that -v flag is recognized (AC-1)
func TestVerboseFlag_ShortFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Find a test file
	testFile := findTestFile(t, "../../testdata/xmp")

	// Run with -v flag
	cmd := exec.Command(binPath, "convert", testFile, "--to", "np3", "-v", "--overwrite")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stderr
	stderrBytes, _ := io.ReadAll(stderr)
	stderrStr := string(stderrBytes)

	cmd.Wait()

	// Verify verbose logs appear in stderr (AC-3)
	if !strings.Contains(stderrStr, "level=DEBUG") && !strings.Contains(stderrStr, "level=INFO") {
		t.Errorf("Expected verbose logs with -v flag, got:\n%s", stderrStr)
	}
}

// TestVerboseFlag_LongFlag tests that --verbose flag is recognized (AC-1)
func TestVerboseFlag_LongFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Find a test file
	testFile := findTestFile(t, "../../testdata/xmp")

	// Run with --verbose flag
	cmd := exec.Command(binPath, "convert", testFile, "--to", "np3", "--verbose")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stderr
	stderrBytes, _ := io.ReadAll(stderr)
	stderrStr := string(stderrBytes)

	cmd.Wait()

	// Verify verbose logs appear in stderr (AC-3)
	if !strings.Contains(stderrStr, "level=DEBUG") && !strings.Contains(stderrStr, "level=INFO") {
		t.Errorf("Expected verbose logs with --verbose flag, got:\n%s", stderrStr)
	}
}

// TestVerboseConversion_AllStepsLogged verifies all workflow steps are logged (AC-3)
func TestVerboseConversion_AllStepsLogged(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Find a test file
	testFile := findTestFile(t, "../../testdata/xmp")

	// Run with verbose flag
	cmd := exec.Command(binPath, "convert", testFile, "--to", "np3", "-v", "--overwrite")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stderr
	stderrBytes, _ := io.ReadAll(stderr)
	stderrStr := string(stderrBytes)

	cmd.Wait()

	// Verify expected log messages (AC-3)
	expectedMessages := []string{
		"reading input",      // Reading input file
		"parsing file",       // Parsing file
		"converting formats", // Converting formats
		"generating output",  // Generating output
		"writing output",     // Writing output
		"conversion completed", // Completion with timing
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(stderrStr, msg) {
			t.Errorf("Expected log message %q not found in stderr:\n%s", msg, stderrStr)
		}
	}
}

// TestVerboseConversion_StructuredFields verifies slog structured format (AC-2)
func TestVerboseConversion_StructuredFields(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Find a test file
	testFile := findTestFile(t, "../../testdata/xmp")

	// Run with verbose flag
	cmd := exec.Command(binPath, "convert", testFile, "--to", "np3", "-v", "--overwrite")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stderr
	stderrBytes, _ := io.ReadAll(stderr)
	stderrStr := string(stderrBytes)

	cmd.Wait()

	// Verify structured fields (AC-2)
	requiredFields := []string{
		"level=",      // Log level
		"msg=",        // Message
		"file=",       // File field
		"format=",     // Format field
		"duration_ms=", // Duration field
	}

	for _, field := range requiredFields {
		if !strings.Contains(stderrStr, field) {
			t.Errorf("Expected structured field %q not found in stderr:\n%s", field, stderrStr)
		}
	}
}

// TestNormalMode_NoVerboseLogs verifies normal mode has no debug logs (AC-2)
func TestNormalMode_NoVerboseLogs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Find a test file
	testFile := findTestFile(t, "../../testdata/xmp")

	// Run WITHOUT verbose flag
	cmd := exec.Command(binPath, "convert", testFile, "--to", "np3")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stderr
	stderrBytes, _ := io.ReadAll(stderr)
	stderrStr := string(stderrBytes)

	cmd.Wait()

	// Verify NO debug logs in normal mode (AC-2)
	if strings.Contains(stderrStr, "level=DEBUG") {
		t.Errorf("Debug logs should not appear in normal mode:\n%s", stderrStr)
	}
}

// buildBinary builds the recipe CLI binary for testing
func buildBinary(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "recipe.exe")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	return binPath
}

// findTestFile finds a test file in the specified directory
func findTestFile(t *testing.T, dir string) string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.xmp"))
	if err != nil || len(files) == 0 {
		t.Fatalf("No XMP test files found in %s", dir)
	}

	return files[0]
}
