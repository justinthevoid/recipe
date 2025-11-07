package main

import (
	"testing"
)

// mockKeyMsg for testing
type mockKeyMsgBatch struct {
	str string
}

func (m mockKeyMsgBatch) String() string {
	return m.str
}

// TestConversionTrigger tests the conversion trigger with 'c' key
func TestConversionTrigger(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", Format: "xmp"},
		{Name: "file2.lrtemplate", Path: "/test/file2.lrtemplate", Format: "lrtemplate"},
	}
	m.selected["/test/file1.xmp"] = true
	m.selected["/test/file2.lrtemplate"] = true

	// Press 'c' key
	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"c"})
	m = newModel.(model)

	if !m.showFormatPrompt {
		t.Error("Expected showFormatPrompt to be true after 'c' key press")
	}
}

// TestConversionTriggerNoSelection tests that 'c' key does nothing when no files are selected
func TestConversionTriggerNoSelection(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", Format: "xmp"},
	}

	// Press 'c' key with no selection
	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"c"})
	m = newModel.(model)

	if m.showFormatPrompt {
		t.Error("Expected showFormatPrompt to remain false when no files selected")
	}
}

// TestFormatSelection tests format menu navigation
func TestFormatSelection(t *testing.T) {
	m := initialModel()
	m.showFormatPrompt = true
	m.outputDir = "." // Set output directory for validation
	m.selected["/test/file1.np3"] = true
	m.files = []FileInfo{
		{Name: "file1.np3", Path: "/test/file1.np3", Format: "np3", Size: 1024},
	}

	tests := []struct {
		key            string
		expectedFormat string
	}{
		{"1", "np3"},
		{"2", "xmp"},
		{"3", "lrtemplate"},
	}

	for _, tt := range tests {
		m.showFormatPrompt = true
		m.showValidation = false // Reset validation state
		newModel, _ := m.handleKeyPress(mockKeyMsgBatch{tt.key})
		m = newModel.(model)

		if m.targetFormat != tt.expectedFormat {
			t.Errorf("Expected targetFormat %s, got %s", tt.expectedFormat, m.targetFormat)
		}

		// Story 4-4: After format selection, should show validation screen
		if !m.showValidation {
			t.Error("Expected showValidation to be true after format selection")
		}

		// Should not go directly to confirmation anymore
		if m.showConfirmation {
			t.Error("Expected showConfirmation to be false (validation comes first)")
		}
	}
}

// TestConfirmationScreen tests confirmation screen logic
func TestConfirmationScreen(t *testing.T) {
	m := initialModel()
	m.showConfirmation = true
	m.targetFormat = "xmp"
	m.selected["/test/file1.np3"] = true
	m.files = []FileInfo{
		{Name: "file1.np3", Path: "/test/file1.np3", Format: "np3"},
	}

	// Test 'y' (confirm)
	newModel, cmd := m.handleKeyPress(mockKeyMsgBatch{"y"})
	m = newModel.(model)

	if !m.converting {
		t.Error("Expected converting to be true after confirmation")
	}

	if cmd == nil {
		t.Error("Expected non-nil command after confirmation")
	}
}

// TestConfirmationCancel tests cancelling confirmation
func TestConfirmationCancel(t *testing.T) {
	m := initialModel()
	m.showConfirmation = true

	// Test 'n' (cancel)
	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"n"})
	m = newModel.(model)

	if m.showConfirmation {
		t.Error("Expected showConfirmation to be false after cancel")
	}
}

// TestEscapeHandling tests Esc key in various screens
func TestEscapeHandling(t *testing.T) {
	// Test Esc in format prompt
	m := initialModel()
	m.showFormatPrompt = true

	newModel, _ := m.handleKeyPress(mockKeyMsgBatch{"esc"})
	m = newModel.(model)

	if m.showFormatPrompt {
		t.Error("Expected showFormatPrompt to be false after Esc")
	}

	// Test Esc in confirmation screen
	m.showConfirmation = true
	newModel, _ = m.handleKeyPress(mockKeyMsgBatch{"esc"})
	m = newModel.(model)

	if m.showConfirmation {
		t.Error("Expected showConfirmation to be false after Esc")
	}
}
