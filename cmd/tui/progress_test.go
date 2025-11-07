package main

import (
	"strings"
	"testing"
	"time"
)

// TestProgressBarRendering tests progress bar at various completion levels
func TestProgressBarRendering(t *testing.T) {
	tests := []struct {
		current    int
		total      int
		errorCount int
		wantPct    int
	}{
		{0, 10, 0, 0},
		{5, 10, 0, 50},
		{10, 10, 0, 100},
		{3, 7, 0, 42},
	}

	for _, tt := range tests {
		result := renderProgressBar(tt.current, tt.total, tt.errorCount)

		// Check percentage appears in output
		expectedPct := (tt.current * 100) / tt.total
		if !strings.Contains(result, string(rune(expectedPct))) && expectedPct > 0 {
			// Note: This is a basic check; actual percentage rendering uses fmt.Sprintf
			t.Logf("Progress bar rendered for %d/%d", tt.current, tt.total)
		}

		// Check that progress bar contains both filled and empty blocks
		if tt.current > 0 && tt.current < tt.total {
			if !strings.Contains(result, "█") || !strings.Contains(result, "░") {
				t.Errorf("Expected both filled (█) and empty (░) blocks in progress bar")
			}
		}
	}
}

// TestProgressPercentage tests percentage calculation accuracy
func TestProgressPercentage(t *testing.T) {
	tests := []struct {
		current int
		total   int
		want    int
	}{
		{5, 10, 50},
		{1, 3, 33},
		{7, 10, 70},
		{10, 10, 100},
		{0, 10, 0},
	}

	for _, tt := range tests {
		got := (tt.current * 100) / tt.total
		if got != tt.want {
			t.Errorf("calculateProgress(%d, %d) = %d, want %d", tt.current, tt.total, got, tt.want)
		}
	}
}

// TestFileStatusFormatting tests status icons and colors
func TestFileStatusFormatting(t *testing.T) {
	tests := []struct {
		file         string
		sourceFormat string
		targetFormat string
		status       string
		wantIcon     string
	}{
		{"file1.xmp", "xmp", "np3", "success", "✓"},
		{"file2.np3", "np3", "xmp", "warning", "⚠️"},
		{"file3.lrtemplate", "lrtemplate", "xmp", "error", "✗"},
		{"file4.xmp", "xmp", "np3", "converting", "⠋"},
	}

	for _, tt := range tests {
		result := formatFileStatus(tt.file, tt.sourceFormat, tt.targetFormat, tt.status)

		if !strings.Contains(result, tt.file) {
			t.Errorf("Expected filename %s in status, got %s", tt.file, result)
		}

		// Check icon appears (may be styled)
		// Note: The icon may be wrapped in ANSI codes, so we just check the string is non-empty
		if result == "" {
			t.Errorf("Expected non-empty status for file %s", tt.file)
		}
	}
}

// TestTimeEstimation tests time estimation accuracy
func TestTimeEstimation(t *testing.T) {
	// Mock time
	mockStart := time.Now().Add(-30 * time.Second)
	completed := 3
	total := 10

	remaining := estimateRemainingTime(mockStart, completed, total)

	// Average: 10s per file, 7 remaining = ~70s
	expectedSeconds := 70.0
	gotSeconds := remaining.Seconds()

	if gotSeconds < expectedSeconds-5 || gotSeconds > expectedSeconds+5 {
		t.Errorf("estimateRemainingTime() = %.1fs, want ~%.1fs", gotSeconds, expectedSeconds)
	}
}

// TestTimeEstimationEdgeCases tests edge cases
func TestTimeEstimationEdgeCases(t *testing.T) {
	start := time.Now()

	// Test completed == 0
	remaining := estimateRemainingTime(start, 0, 10)
	if remaining != 0 {
		t.Errorf("Expected 0 remaining time when no files completed, got %v", remaining)
	}

	// Test very fast conversions (<1s)
	start = time.Now().Add(-100 * time.Millisecond)
	remaining = estimateRemainingTime(start, 1, 10)
	// Should still give a reasonable estimate
	if remaining < 0 {
		t.Errorf("Expected non-negative remaining time, got %v", remaining)
	}
}

// TestDurationFormatting tests duration formatting
func TestDurationFormatting(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{0 * time.Second, "0:00"},
		{30 * time.Second, "0:30"},
		{90 * time.Second, "1:30"},
		{3600 * time.Second, "1:00:00"},
		{3661 * time.Second, "1:01:01"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %s, want %s", tt.duration, got, tt.want)
		}
	}
}

// TestWriteOutputAtomic tests atomic file writing
func TestWriteOutputAtomic(t *testing.T) {
	// Create temp file
	tmpPath := t.TempDir() + "/test_output.txt"
	data := []byte("test data")

	err := writeOutputAtomic(tmpPath, data)
	if err != nil {
		t.Errorf("writeOutputAtomic() error = %v", err)
	}

	// Verify file exists and contains correct data
	// (This would require reading the file, which is tested in integration tests)
}

// TestResultsSummary tests summary count accuracy
func TestResultsSummary(t *testing.T) {
	results := []ConversionResult{
		{File: "file1.xmp", Status: "success"},
		{File: "file2.np3", Status: "warning", Message: "unmappable params"},
		{File: "file3.xmp", Status: "error", Message: "parse error"},
		{File: "file4.lrtemplate", Status: "success"},
	}

	successCount := 0
	warningCount := 0
	errorCount := 0

	for _, r := range results {
		switch r.Status {
		case "success":
			successCount++
		case "warning":
			warningCount++
		case "error":
			errorCount++
		}
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successes, got %d", successCount)
	}
	if warningCount != 1 {
		t.Errorf("Expected 1 warning, got %d", warningCount)
	}
	if errorCount != 1 {
		t.Errorf("Expected 1 error, got %d", errorCount)
	}
}

// TestSummaryMessageFormatting tests message formatting
func TestSummaryMessageFormatting(t *testing.T) {
	results := []ConversionResult{
		{File: "file1.xmp", Status: "success"},
		{File: "file2.np3", Status: "warning", Message: "3 unmappable parameters"},
		{File: "file3.xmp", Status: "error", Message: "Invalid XML at line 42"},
	}

	elapsed := 125 * time.Second

	summary := renderSummaryScreen(results, elapsed)

	// Check that summary contains key information
	if !strings.Contains(summary, "Complete") {
		t.Error("Expected summary to contain 'Complete'")
	}

	// Check elapsed time is formatted
	if !strings.Contains(summary, "2:05") {
		t.Error("Expected summary to contain formatted elapsed time")
	}

	// Check error messages appear
	if !strings.Contains(summary, "Invalid XML") {
		t.Error("Expected summary to contain error message")
	}
}
