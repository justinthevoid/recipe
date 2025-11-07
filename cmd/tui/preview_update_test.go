package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestUpdatePreviewWithCache tests that cached content is used immediately
func TestUpdatePreviewWithCache(t *testing.T) {
	m := model{
		files: []FileInfo{
			{Path: "test1.xmp", Name: "test1.xmp", IsDir: false},
			{Path: "test2.xmp", Name: "test2.xmp", IsDir: false},
		},
		cursor:       0,
		showPreview:  true,
		previewCache: make(map[string]string),
	}

	// Add cached content
	cachedContent := "Cached parameters"
	m.previewCache["test1.xmp"] = cachedContent

	// Call updatePreview
	m, cmd := m.updatePreview()

	// Should use cached content
	if m.previewContent != cachedContent {
		t.Errorf("expected cached content %q, got %q", cachedContent, m.previewContent)
	}

	// Should not trigger async load (cmd should be nil)
	if cmd != nil {
		t.Error("expected no command when using cache, but got command")
	}

	// Should reset scroll offset
	if m.scrollOffset != 0 {
		t.Errorf("expected scroll offset to reset to 0, got %d", m.scrollOffset)
	}
}

// TestUpdatePreviewWithoutCache tests that async load is triggered without cache
func TestUpdatePreviewWithoutCache(t *testing.T) {
	m := model{
		files: []FileInfo{
			{Path: "test1.xmp", Name: "test1.xmp", IsDir: false},
		},
		cursor:       0,
		showPreview:  true,
		previewCache: make(map[string]string),
	}

	// Call updatePreview without cached content
	m, cmd := m.updatePreview()

	// Should set loading state
	if !m.previewLoading {
		t.Error("expected previewLoading to be true")
	}

	// Should trigger async load command
	if cmd == nil {
		t.Error("expected command to trigger async load, but got nil")
	}

	// Should set preview file
	if m.previewFile != "test1.xmp" {
		t.Errorf("expected previewFile %q, got %q", "test1.xmp", m.previewFile)
	}

	// Should reset scroll offset
	if m.scrollOffset != 0 {
		t.Errorf("expected scroll offset to reset to 0, got %d", m.scrollOffset)
	}
}

// TestUpdatePreviewForDirectory tests that directories show special message
func TestUpdatePreviewForDirectory(t *testing.T) {
	m := model{
		files: []FileInfo{
			{Path: "subdir", Name: "subdir", IsDir: true},
		},
		cursor:       0,
		showPreview:  true,
		previewCache: make(map[string]string),
	}

	// Call updatePreview on directory
	m, cmd := m.updatePreview()

	// Should show directory message
	if m.previewContent != "  (Directories cannot be previewed)" {
		t.Errorf("expected directory message, got %q", m.previewContent)
	}

	// Should not trigger async load
	if cmd != nil {
		t.Error("expected no command for directory, but got command")
	}
}

// TestUpdatePreviewWhenDisabled tests that preview doesn't update when disabled
func TestUpdatePreviewWhenDisabled(t *testing.T) {
	m := model{
		files: []FileInfo{
			{Path: "test1.xmp", Name: "test1.xmp", IsDir: false},
		},
		cursor:       0,
		showPreview:  false, // Disabled
		previewCache: make(map[string]string),
	}

	// Call updatePreview when preview is disabled
	m, cmd := m.updatePreview()

	// Should clear preview
	if m.previewContent != "" {
		t.Errorf("expected empty content when preview disabled, got %q", m.previewContent)
	}

	// Should not trigger async load
	if cmd != nil {
		t.Error("expected no command when preview disabled, but got command")
	}
}

// TestPreviewLoadedMsg tests handling of async preview results
func TestPreviewLoadedMsg(t *testing.T) {
	m := model{
		previewFile:    "test.xmp",
		previewLoading: true,
		previewCache:   make(map[string]string),
	}

	// Simulate successful async load
	msg := previewLoadedMsg{
		filePath: "test.xmp",
		content:  "Exposure: +0.5\nContrast: +10",
		err:      nil,
	}

	updatedModel, _ := m.Update(msg)
	m = updatedModel.(model)

	// Should clear loading state
	if m.previewLoading {
		t.Error("expected previewLoading to be false after load")
	}

	// Should cache the content
	if cachedContent, ok := m.previewCache["test.xmp"]; !ok {
		t.Error("expected content to be cached")
	} else if cachedContent != msg.content {
		t.Errorf("expected cached content %q, got %q", msg.content, cachedContent)
	}

	// Should update display
	if m.previewContent != msg.content {
		t.Errorf("expected preview content %q, got %q", msg.content, m.previewContent)
	}
}

// TestPreviewLoadedMsgWithError tests error handling in async load
func TestPreviewLoadedMsgWithError(t *testing.T) {
	m := model{
		previewFile:    "test.xmp",
		previewLoading: true,
		previewCache:   make(map[string]string),
	}

	// Simulate failed async load
	msg := previewLoadedMsg{
		filePath: "test.xmp",
		content:  "",
		err:      fmt.Errorf("parse error"),
	}

	updatedModel, _ := m.Update(msg)
	m = updatedModel.(model)

	// Should clear loading state
	if m.previewLoading {
		t.Error("expected previewLoading to be false after error")
	}

	// Should show error message
	if !strings.Contains(m.previewContent, "Error loading preview") {
		t.Errorf("expected error message in preview, got %q", m.previewContent)
	}
}

// TestTabTogglesPreviewFocus tests Tab key toggles focus
func TestTabTogglesPreviewFocus(t *testing.T) {
	m := model{
		showPreview:    true,
		previewFocused: false,
	}

	// Press Tab
	newM, _ := m.handleKeyPress(mockKeyMsg{"tab"})
	m = newM.(model)

	// Should toggle focus to true
	if !m.previewFocused {
		t.Error("expected previewFocused to be true after first Tab")
	}

	// Press Tab again
	newM, _ = m.handleKeyPress(mockKeyMsg{"tab"})
	m = newM.(model)

	// Should toggle focus back to false
	if m.previewFocused {
		t.Error("expected previewFocused to be false after second Tab")
	}
}

// TestPageUpDownScrolling tests PageUp/PageDown keys
func TestPageUpDownScrolling(t *testing.T) {
	// Create content with many lines
	lines := []string{}
	for i := 0; i < 50; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i))
	}
	content := strings.Join(lines, "\n")

	m := model{
		showPreview:    true,
		previewFocused: true,
		previewContent: content,
		viewportHeight: 10,
		scrollOffset:   0,
	}

	// Press PageDown
	newM, _ := m.handleKeyPress(mockKeyMsg{"pgdn"})
	m = newM.(model)

	// Should scroll down by viewport height
	if m.scrollOffset != 10 {
		t.Errorf("expected scrollOffset 10 after PageDown, got %d", m.scrollOffset)
	}

	// Press PageUp
	newM, _ = m.handleKeyPress(mockKeyMsg{"pgup"})
	m = newM.(model)

	// Should scroll back to top
	if m.scrollOffset != 0 {
		t.Errorf("expected scrollOffset 0 after PageUp, got %d", m.scrollOffset)
	}
}

// TestSplitLines tests line splitting helper
func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single line", "Line 1", 1},
		{"Multiple lines", "Line 1\nLine 2\nLine 3", 3},
		{"Trailing newline", "Line 1\nLine 2\n", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := splitLines(tt.content)
			if len(lines) != tt.expected {
				t.Errorf("expected %d lines, got %d", tt.expected, len(lines))
			}
		})
	}
}
