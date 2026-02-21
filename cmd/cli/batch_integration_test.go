package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBatch_100Files tests batch conversion of 100 files for performance (AC-3)
func TestBatch_100Files(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	// Create 100 test files
	tmpDir := t.TempDir()
	for i := 0; i < 100; i++ {
		filename := fmt.Sprintf("file%d.xmp", i)
		createValidXMPFile(t, filepath.Join(tmpDir, filename))
	}

	// Run batch conversion
	start := time.Now()
	cmd := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--overwrite")
	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	// Assertions
	if err != nil {
		t.Fatalf("batch conversion failed: %v\nOutput: %s", err, output)
	}

	// Verify performance (<2s target)
	if elapsed > 2*time.Second {
		t.Errorf("batch took %v, want <2s", elapsed)
	}
	t.Logf("Performance: Converted 100 files in %v (target: <2s)", elapsed)

	// Verify all files converted
	files, _ := filepath.Glob(filepath.Join(tmpDir, "*.np3"))
	if len(files) != 100 {
		t.Errorf("got %d output files, want 100", len(files))
	}

	// Verify summary message
	outputStr := string(output)
	if !strings.Contains(outputStr, "100 success") {
		t.Errorf("success summary not found in output: %s", outputStr)
	}
}

// TestBatch_ErrorHandling tests continue-on-error and fail-fast modes (AC-5)
func TestBatch_ErrorHandling(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()

	// Create mix of valid and invalid files
	createValidXMPFile(t, filepath.Join(tmpDir, "valid1.xmp"))
	createValidXMPFile(t, filepath.Join(tmpDir, "valid2.xmp"))
	createInvalidFile(t, filepath.Join(tmpDir, "invalid.xmp"))

	t.Run("continue-on-error", func(t *testing.T) {
		// Default behavior: continue on error
		cmd := exec.Command("./recipe-test.exe", "batch",
			filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--overwrite")
		output, err := cmd.CombinedOutput()

		// Should exit with error code but process all files
		if err == nil {
			t.Error("expected error exit code, got success")
		}

		outputStr := string(output)
		// Verify all 3 files were attempted
		if !strings.Contains(outputStr, "3 files") {
			t.Errorf("expected 3 files processed, output: %s", outputStr)
		}
		// Verify 2 successes
		if !strings.Contains(outputStr, "2 success") {
			t.Errorf("expected 2 successes, output: %s", outputStr)
		}
		// Verify 1 error
		if !strings.Contains(outputStr, "1 error") {
			t.Errorf("expected 1 error, output: %s", outputStr)
		}
	})

	// Clean up output files for fail-fast test
	os.Remove(filepath.Join(tmpDir, "valid1.np3"))
	os.Remove(filepath.Join(tmpDir, "valid2.np3"))

	t.Run("fail-fast", func(t *testing.T) {
		cmd := exec.Command("./recipe-test.exe", "batch",
			filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--fail-fast", "--overwrite")
		output, err := cmd.CombinedOutput()

		// Should exit with error
		if err == nil {
			t.Error("expected error exit code, got success")
		}

		outputStr := string(output)
		t.Logf("fail-fast output: %s", outputStr)
		// May not process all files if stops on first error
		// Just verify it exited with error
	})
}

// TestBatch_OverwriteProtection tests overwrite behavior (AC-9)
func TestBatch_OverwriteProtection(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()

	// Create test files
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("test%d.xmp", i)
		createValidXMPFile(t, filepath.Join(tmpDir, filename))
	}

	// First conversion - should succeed
	cmd1 := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--overwrite")
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		t.Fatalf("first batch failed: %v\nOutput: %s", err, output1)
	}

	// Verify 5 successes
	if !strings.Contains(string(output1), "5 success") {
		t.Errorf("first batch: expected 5 successes, output: %s", output1)
	}

	// Second conversion without --overwrite - should error (files exist)
	cmd2 := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3")
	output2, err := cmd2.CombinedOutput()

	// Should fail because all files already exist
	if err == nil {
		t.Error("second batch should have failed due to existing files without --overwrite")
	}

	// Verify all files had errors (overwrite protection)
	outputStr2 := string(output2)
	if !strings.Contains(outputStr2, "5 errors") {
		t.Errorf("second batch: expected 5 errors (overwrite protection), output: %s", outputStr2)
	}

	// Third conversion with --overwrite - should overwrite all
	cmd3 := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--overwrite")
	output3, err := cmd3.CombinedOutput()
	if err != nil {
		t.Fatalf("third batch failed: %v\nOutput: %s", err, output3)
	}

	// Verify 5 successes again
	if !strings.Contains(string(output3), "5 success") {
		t.Errorf("third batch: expected 5 successes, output: %s", output3)
	}
}

// TestBatch_CustomOutputDirectory tests --output-dir flag (AC-6)
func TestBatch_CustomOutputDirectory(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "converted")

	// Create test files
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("test%d.xmp", i)
		createValidXMPFile(t, filepath.Join(tmpDir, filename))
	}

	// Run batch with custom output directory
	cmd := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3",
		"--output-dir", outputDir, "--overwrite")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("batch failed: %v\nOutput: %s", err, output)
	}

	// Verify output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("output directory not created: %s", outputDir)
	}

	// Verify files are in output directory
	files, _ := filepath.Glob(filepath.Join(outputDir, "*.np3"))
	if len(files) != 5 {
		t.Errorf("got %d output files in %s, want 5", len(files), outputDir)
	}

	// Verify original directory has no output files
	originalFiles, _ := filepath.Glob(filepath.Join(tmpDir, "*.np3"))
	if len(originalFiles) > 0 {
		t.Errorf("found %d output files in original directory, want 0", len(originalFiles))
	}
}

// TestBatch_JSONOutput tests --json flag output format (AC-7)
func TestBatch_JSONOutput(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()

	// Create test files
	for i := 0; i < 3; i++ {
		filename := fmt.Sprintf("test%d.xmp", i)
		createValidXMPFile(t, filepath.Join(tmpDir, filename))
	}

	// Run batch with JSON output
	cmd := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "*.xmp"), "--to", "np3", "--json", "--overwrite")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("batch failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Verify JSON structure (AC-4: Check for required batch fields)
	if !strings.Contains(outputStr, "\"batch\"") {
		t.Error("JSON output missing batch field")
	}
	if !strings.Contains(outputStr, "\"total\"") {
		t.Error("JSON output missing total field")
	}
	if !strings.Contains(outputStr, "\"success_count\"") {
		t.Error("JSON output missing success_count field")
	}
	if !strings.Contains(outputStr, "\"error_count\"") {
		t.Error("JSON output missing error_count field")
	}
	if !strings.Contains(outputStr, "\"duration_ms\"") {
		t.Error("JSON output missing duration_ms field")
	}
	if !strings.Contains(outputStr, "\"results\"") {
		t.Error("JSON output missing results array")
	}

	// Verify no progress messages in JSON mode
	if strings.Contains(outputStr, "Processing") {
		t.Error("JSON output should not contain progress messages")
	}
}

// TestBatch_GlobPatternValidation tests glob pattern error handling (AC-8)
func TestBatch_GlobPatternValidation(t *testing.T) {
	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "recipe-test.exe", ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}
	defer os.Remove("recipe-test.exe")

	tmpDir := t.TempDir()

	// Test: No matches
	cmd := exec.Command("./recipe-test.exe", "batch",
		filepath.Join(tmpDir, "nonexistent*.xmp"), "--to", "np3")
	output, err := cmd.CombinedOutput()

	// Should fail
	if err == nil {
		t.Error("expected error for no matches, got success")
	}

	// Verify error message mentions pattern
	outputStr := string(output)
	if !strings.Contains(outputStr, "no files match") {
		t.Errorf("error message should mention no files match, got: %s", outputStr)
	}
}
