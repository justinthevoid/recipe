package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestTriggerValidation tests the validation screen trigger
func TestTriggerValidation(t *testing.T) {
	m := initialModel()
	m.targetFormat = "xmp"
	m.outputDir = "."
	m.selected = map[string]bool{
		"file1.np3": true,
		"file2.xmp": true,
	}
	m.files = []FileInfo{
		{Name: "file1.np3", Path: "file1.np3", Size: 1024, Format: "np3"},
		{Name: "file2.xmp", Path: "file2.xmp", Size: 2048, Format: "xmp"},
	}

	newModel, _ := m.triggerValidation()

	if !newModel.showValidation {
		t.Error("Expected showValidation to be true")
	}
	if len(newModel.validationFiles) != 2 {
		t.Errorf("Expected 2 validation files, got %d", len(newModel.validationFiles))
	}
}

// TestValidateOutputDirectory tests directory validation
func TestValidateOutputDirectory(t *testing.T) {
	// Test with current directory (should pass)
	cwd, _ := os.Getwd()
	err := validateOutputDirectory(cwd)
	if err != nil {
		t.Errorf("Expected current directory to be valid, got error: %v", err)
	}

	// Test with non-existent directory
	err = validateOutputDirectory("/nonexistent/directory/path")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

// TestDetectOverwrites tests overwrite detection
func TestDetectOverwrites(t *testing.T) {
	// Create a temp directory with test files
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.xmp")
	os.WriteFile(testFile, []byte("test"), 0644)

	files := []ValidationFile{
		{Name: "test.np3", Path: filepath.Join(tmpDir, "test.np3"), Size: 1024, SourceFormat: "np3"},
	}

	overwrites := detectOverwrites(files, tmpDir, "xmp")

	if len(overwrites) != 1 {
		t.Errorf("Expected 1 overwrite, got %d", len(overwrites))
	}
	if len(overwrites) > 0 && overwrites[0].File != "test.np3" {
		t.Errorf("Expected overwrite for test.np3, got %s", overwrites[0].File)
	}
}

// TestChangeExtension tests file extension changing
func TestChangeExtension(t *testing.T) {
	tests := []struct {
		filename string
		format   string
		expected string
	}{
		{"file.np3", "xmp", "file.xmp"},
		{"file.xmp", "np3", "file.np3"},
		{"file.lrtemplate", "xmp", "file.xmp"},
		{"file.name.np3", "lrtemplate", "file.name.lrtemplate"},
	}

	for _, tt := range tests {
		result := changeExtension(tt.filename, tt.format)
		if result != tt.expected {
			t.Errorf("changeExtension(%q, %q) = %q, want %q", tt.filename, tt.format, result, tt.expected)
		}
	}
}

// TestCalculateConversionPlan tests conversion plan calculation
func TestCalculateConversionPlan(t *testing.T) {
	files := []ValidationFile{
		{Name: "file1.np3", Size: 1024, SourceFormat: "np3", TargetFormat: "xmp"},
		{Name: "file2.xmp", Size: 2048, SourceFormat: "xmp", TargetFormat: "xmp"},
		{Name: "file3.lrtemplate", Size: 512, SourceFormat: "lrtemplate", TargetFormat: "xmp"},
	}

	plan := calculateConversionPlan(files, "xmp")

	if plan.FileCount != 3 {
		t.Errorf("Expected FileCount 3, got %d", plan.FileCount)
	}
	if plan.TotalInputSize != 3584 { // 1024 + 2048 + 512
		t.Errorf("Expected TotalInputSize 3584, got %d", plan.TotalInputSize)
	}
	if plan.CrossFormatCount != 2 {
		t.Errorf("Expected CrossFormatCount 2, got %d", plan.CrossFormatCount)
	}
	if plan.SameFormatCount != 1 {
		t.Errorf("Expected SameFormatCount 1, got %d", plan.SameFormatCount)
	}
	if plan.EstimatedTime == 0 {
		t.Error("Expected EstimatedTime to be non-zero")
	}
}

// TestRenderValidationScreen tests validation screen rendering
func TestRenderValidationScreen(t *testing.T) {
	m := initialModel()
	m.showValidation = true
	m.targetFormat = "xmp"
	m.outputDir = "/tmp/test"
	m.validationFiles = []ValidationFile{
		{Name: "file1.np3", Size: 1024, SourceFormat: "np3", TargetFormat: "xmp", HasWarnings: false},
	}
	m.validationPlan = ConversionPlan{
		FileCount:           1,
		TotalInputSize:      1024,
		EstimatedOutputSize: 1075,
		CrossFormatCount:    1,
		SameFormatCount:     0,
	}
	m.validationPassed = true

	content := renderValidationScreen(m)

	if content == "" {
		t.Error("Expected non-empty validation screen content")
	}
	if !contains(content, "Conversion Validation") {
		t.Error("Expected validation screen to contain title")
	}
	if !contains(content, "file1.np3") {
		t.Error("Expected validation screen to contain file name")
	}
}

// TestDirectoryValidationDisplay tests directory issue rendering
func TestDirectoryValidationDisplay(t *testing.T) {
	// Test missing directory
	m := initialModel()
	m.showValidation = true
	m.targetFormat = "xmp"
	m.outputDir = "/nonexistent/directory"
	m.directoryIssue = "missing"
	m.validationPassed = false
	m.validationFiles = []ValidationFile{
		{Name: "file1.np3", Size: 1024, SourceFormat: "np3", TargetFormat: "xmp"},
	}
	m.validationPlan = ConversionPlan{FileCount: 1}

	content := renderValidationScreen(m)

	if !contains(content, "Directory Issues") {
		t.Error("Expected 'Directory Issues' in validation screen with directory issue")
	}
	if !contains(content, "Output directory does not exist") {
		t.Error("Expected missing directory message")
	}

	// Test permission issue
	m.directoryIssue = "permission"
	content = renderValidationScreen(m)

	if !contains(content, "No write permission") {
		t.Error("Expected permission error message")
	}
}

// TestOverwriteWarningsDisplay tests overwrite warnings rendering
func TestOverwriteWarningsDisplay(t *testing.T) {
	m := initialModel()
	m.showValidation = true
	m.targetFormat = "xmp"
	m.outputDir = "."
	m.validationPassed = true
	m.validationFiles = []ValidationFile{
		{Name: "file1.np3", Size: 1024, SourceFormat: "np3", TargetFormat: "xmp"},
	}
	m.validationPlan = ConversionPlan{FileCount: 1}
	m.overwriteFiles = []OverwriteInfo{
		{File: "file1.np3", ExistingSize: 2048, NewSize: 1024},
		{File: "file2.np3", ExistingSize: 3072, NewSize: 1536},
	}

	content := renderValidationScreen(m)

	if !contains(content, "Overwrite Warnings") {
		t.Error("Expected 'Overwrite Warnings' in validation screen with overwrites")
	}
	if !contains(content, "file1.np3") {
		t.Error("Expected overwrite file name in display")
	}
}

// TestValidationToConfirmationFlow tests the validation → confirmation flow (AC-6)
func TestValidationToConfirmationFlow(t *testing.T) {
	m := initialModel()
	m.showValidation = true
	m.validationPassed = true
	m.targetFormat = "xmp"
	m.outputDir = "."
	m.selected["/test/file1.np3"] = true
	m.files = []FileInfo{
		{Name: "file1.np3", Path: "/test/file1.np3", Format: "np3", Size: 1024},
	}
	m.validationFiles = []ValidationFile{
		{Name: "file1.np3", Path: "/test/file1.np3", Size: 1024, SourceFormat: "np3", TargetFormat: "xmp"},
	}

	// Press 'c' to move from validation to confirmation
	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"c"})
	m = newModel.(model)

	if m.showValidation {
		t.Error("Expected showValidation to be false after 'c' key")
	}
	if !m.showConfirmation {
		t.Error("Expected showConfirmation to be true after 'c' key")
	}

	// Press 'y' to start conversion
	newModel, cmd := m.handleKeyPress(mockKeyMsgBatch{"y"})
	m = newModel.(model)

	if !m.converting {
		t.Error("Expected converting to be true after 'y' key")
	}
	if cmd == nil {
		t.Error("Expected non-nil command after starting conversion")
	}
}

// TestValidationBlocksWithDirectoryIssue tests that validation blocks confirmation when directory is invalid (AC-5)
func TestValidationBlocksWithDirectoryIssue(t *testing.T) {
	m := initialModel()
	m.showValidation = true
	m.validationPassed = false // Blocked by directory issue
	m.directoryIssue = "missing"
	m.targetFormat = "xmp"
	m.outputDir = "/nonexistent"

	// Press 'c' - should NOT move to confirmation
	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"c"})
	m = newModel.(model)

	if !m.showValidation {
		t.Error("Expected to remain on validation screen when validation failed")
	}
	if m.showConfirmation {
		t.Error("Expected showConfirmation to remain false when validation failed")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
