package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestWorkerPoolProcessesAllFiles verifies all files are processed exactly once
func TestWorkerPoolProcessesAllFiles(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	testFiles := createTestFiles(t, tmpDir, 10)

	flags := BatchFlags{
		To:              "np3",
		From:            "xmp",
		Parallel:        4,
		ContinueOnError: true,
		Overwrite:       true,
	}

	// Process batch
	result, err := processBatch(testFiles, flags)
	if err != nil {
		t.Fatalf("processBatch failed: %v", err)
	}

	// Verify all files processed
	if result.Total != 10 {
		t.Errorf("got %d files processed, want 10", result.Total)
	}

	// Verify results match files
	if len(result.Results) != 10 {
		t.Errorf("got %d results, want 10", len(result.Results))
	}
}

// TestGlobPatternExpansion tests glob pattern expansion (AC-1, AC-8)
func TestGlobPatternExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFiles(t, tmpDir, 5)

	tests := []struct {
		name        string
		pattern     string
		expectFiles int
		expectError bool
	}{
		{
			name:        "wildcard pattern",
			pattern:     filepath.Join(tmpDir, "*.xmp"),
			expectFiles: 5,
			expectError: false,
		},
		{
			name:        "no matches",
			pattern:     filepath.Join(tmpDir, "*.nonexistent"),
			expectFiles: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := filepath.Glob(tt.pattern)
			if err != nil {
				if !tt.expectError {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			if len(files) == 0 && tt.expectFiles > 0 {
				if !tt.expectError {
					t.Errorf("got 0 files, want %d", tt.expectFiles)
				}
			}

			if len(files) != tt.expectFiles {
				t.Errorf("got %d files, want %d", len(files), tt.expectFiles)
			}
		})
	}
}

// TestContinueOnError verifies continue-on-error behavior (AC-5)
func TestContinueOnError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mix of valid and invalid files
	validFile := filepath.Join(tmpDir, "valid.xmp")
	invalidFile := filepath.Join(tmpDir, "invalid.xmp")

	createValidXMPFile(t, validFile)
	createInvalidFile(t, invalidFile)

	files := []string{validFile, invalidFile}
	flags := BatchFlags{
		To:              "np3",
		Parallel:        2,
		ContinueOnError: true,
		Overwrite:       true,
	}

	result, err := processBatch(files, flags)
	if err != nil {
		t.Fatalf("processBatch failed: %v", err)
	}

	// Should process both files
	if result.Total != 2 {
		t.Errorf("got %d files processed, want 2", result.Total)
	}

	// Should have 1 success and 1 error
	if result.SuccessCount != 1 {
		t.Errorf("got %d successes, want 1", result.SuccessCount)
	}
	if result.ErrorCount != 1 {
		t.Errorf("got %d errors, want 1", result.ErrorCount)
	}
}

// TestOverwriteProtection tests overwrite protection in batch mode (AC-9)
func TestOverwriteProtection(t *testing.T) {
	tmpDir := t.TempDir()
	testFiles := createTestFiles(t, tmpDir, 3)

	// First conversion - should succeed
	flags1 := BatchFlags{
		To:        "np3",
		Parallel:  2,
		Overwrite: true,
	}

	result1, err := processBatch(testFiles, flags1)
	if err != nil {
		t.Fatalf("first processBatch failed: %v", err)
	}

	if result1.SuccessCount != 3 {
		t.Errorf("first batch: got %d successes, want 3", result1.SuccessCount)
	}

	// Second conversion without overwrite - should skip all
	flags2 := BatchFlags{
		To:        "np3",
		Parallel:  2,
		Overwrite: false,
	}

	result2, err := processBatch(testFiles, flags2)
	if err != nil {
		t.Fatalf("second processBatch failed: %v", err)
	}

	// Skipped files are now counted in error_count (not success)
	if result2.ErrorCount != 3 {
		t.Errorf("second batch: got %d errors (skipped), want 3", result2.ErrorCount)
	}

	// Third conversion with overwrite - should succeed again
	flags3 := BatchFlags{
		To:        "np3",
		Parallel:  2,
		Overwrite: true,
	}

	result3, err := processBatch(testFiles, flags3)
	if err != nil {
		t.Fatalf("third processBatch failed: %v", err)
	}

	if result3.SuccessCount != 3 {
		t.Errorf("third batch: got %d successes, want 3", result3.SuccessCount)
	}
}

// TestCustomOutputDirectory tests output-dir flag (AC-6)
func TestCustomOutputDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testFiles := createTestFiles(t, tmpDir, 3)
	outputDir := filepath.Join(tmpDir, "converted")

	flags := BatchFlags{
		To:        "np3",
		OutputDir: outputDir,
		Parallel:  2,
		Overwrite: true,
	}

	result, err := processBatch(testFiles, flags)
	if err != nil {
		t.Fatalf("processBatch failed: %v", err)
	}

	if result.SuccessCount != 3 {
		t.Errorf("got %d successes, want 3", result.SuccessCount)
	}

	// Verify output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("output directory not created: %s", outputDir)
	}

	// Verify files are in output directory
	for _, r := range result.Results {
		if !strings.HasPrefix(r.Output, outputDir) {
			t.Errorf("output file %s not in output directory %s", r.Output, outputDir)
		}
	}
}

// TestParallelWorkerCount tests --parallel flag (AC-2)
func TestParallelWorkerCount(t *testing.T) {
	tmpDir := t.TempDir()
	testFiles := createTestFiles(t, tmpDir, 20)

	tests := []struct {
		name    string
		workers int
	}{
		{"single worker", 1},
		{"4 workers", 4},
		{"8 workers", 8},
		{"default (NumCPU)", runtime.NumCPU()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := BatchFlags{
				To:        "np3",
				Parallel:  tt.workers,
				Overwrite: true,
			}

			start := time.Now()
			result, err := processBatch(testFiles, flags)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("processBatch failed: %v", err)
			}

			if result.Total != 20 {
				t.Errorf("got %d files processed, want 20", result.Total)
			}

			t.Logf("%s: processed 20 files in %v", tt.name, elapsed)
		})
	}
}

// TestResultAggregation tests batch result aggregation (AC-7)
func TestResultAggregation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mix of scenarios
	validFiles := createTestFiles(t, tmpDir, 5)
	invalidFile := filepath.Join(tmpDir, "invalid.xmp")
	createInvalidFile(t, invalidFile)

	allFiles := append(validFiles, invalidFile)

	flags := BatchFlags{
		To:              "np3",
		Parallel:        4,
		ContinueOnError: true,
		Overwrite:       true,
	}

	result, err := processBatch(allFiles, flags)
	if err != nil {
		t.Fatalf("processBatch failed: %v", err)
	}

	// Verify counts
	expectedTotal := 6
	expectedSuccess := 5
	expectedError := 1

	if result.Total != expectedTotal {
		t.Errorf("Total: got %d, want %d", result.Total, expectedTotal)
	}
	if result.SuccessCount != expectedSuccess {
		t.Errorf("SuccessCount: got %d, want %d", result.SuccessCount, expectedSuccess)
	}
	if result.ErrorCount != expectedError {
		t.Errorf("ErrorCount: got %d, want %d", result.ErrorCount, expectedError)
	}

	// Verify individual results present
	if len(result.Results) != expectedTotal {
		t.Errorf("Results length: got %d, want %d", len(result.Results), expectedTotal)
	}

	// Verify counts match individual results
	actualSuccess := 0
	actualError := 0
	for _, r := range result.Results {
		if r.Success {
			actualSuccess++
		}
		if !r.Success {
			actualError++
		}
	}

	if actualSuccess != expectedSuccess {
		t.Errorf("Counted success: got %d, want %d", actualSuccess, expectedSuccess)
	}
	if actualError != expectedError {
		t.Errorf("Counted errors: got %d, want %d", actualError, expectedError)
	}
}

// Helper functions

func createTestFiles(t *testing.T, dir string, count int) []string {
	t.Helper()
	files := make([]string, count)

	for i := 0; i < count; i++ {
		// Use sprintf for proper filename formatting
		filename := fmt.Sprintf("test%d.xmp", i)
		path := filepath.Join(dir, filename)
		createValidXMPFile(t, path)
		files[i] = path
	}

	return files
}

func createValidXMPFile(t *testing.T, path string) {
	t.Helper()
	// Minimal valid XMP content
	content := `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description rdf:about="">
      <crs:Exposure2012>+0.50</crs:Exposure2012>
      <crs:Contrast2012>+10</crs:Contrast2012>
      <crs:Highlights2012>-20</crs:Highlights2012>
      <crs:Shadows2012>+15</crs:Shadows2012>
      <crs:Whites2012>+5</crs:Whites2012>
      <crs:Blacks2012>-5</crs:Blacks2012>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", path, err)
	}
}

func createInvalidFile(t *testing.T, path string) {
	t.Helper()
	content := "invalid content"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create invalid file %s: %v", path, err)
	}
}
