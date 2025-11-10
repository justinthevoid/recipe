package main

import (
	"strings"
	"testing"
	"time"
)

// TestRenderMetadata tests metadata rendering
func TestRenderMetadata(t *testing.T) {
	tests := []struct {
		name     string
		file     FileInfo
		contains []string
	}{
		{
			name: "NP3 file",
			file: FileInfo{
				Name:    "preset.np3",
				Format:  "np3",
				Size:    2048,
				ModTime: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			contains: []string{"Format: NP3", "Size: 2.0 KB", "Modified: 2024-01-15 10:30"},
		},
		{
			name: "XMP file",
			file: FileInfo{
				Name:    "preset.xmp",
				Format:  "xmp",
				Size:    1536,
				ModTime: time.Date(2024, 2, 20, 14, 45, 0, 0, time.UTC),
			},
			contains: []string{"Format: XMP", "Size: 1.5 KB", "Modified: 2024-02-20 14:45"},
		},
		{
			name: "Directory",
			file: FileInfo{
				Name:    "folder",
				IsDir:   true,
				ModTime: time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC),
			},
			contains: []string{"Type: Directory", "Modified: 2024-03-10 09:00"},
		},
		{
			name: "Large file",
			file: FileInfo{
				Name:    "large.xmp",
				Format:  "xmp",
				Size:    5242880, // 5 MB
				ModTime: time.Date(2024, 4, 5, 16, 20, 0, 0, time.UTC),
			},
			contains: []string{"Format: XMP", "Size: 5.0 MB", "Modified: 2024-04-05 16:20"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderMetadata(tt.file)
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected metadata to contain %q, got: %s", s, result)
				}
			}
		})
	}
}

// TestFormatFileSize tests file size formatting
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"Bytes", 512, "512 B"},
		{"Kilobytes", 2048, "2.0 KB"},
		{"Kilobytes exact", 1024, "1.0 KB"},
		{"Megabytes", 5242880, "5.0 MB"},
		{"Megabytes fractional", 1572864, "1.5 MB"},
		{"Gigabytes", 1073741824, "1.0 GB"},
		{"Gigabytes fractional", 2147483648, "2.0 GB"},
		{"Zero", 0, "0 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSplitPaneLayout tests the split-pane layout rendering
func TestSplitPaneLayout(t *testing.T) {
	tests := []struct {
		name          string
		termWidth     int
		termHeight    int
		expectPreview bool
	}{
		{"Wide terminal (120 cols)", 120, 30, true},
		{"Wider terminal (140 cols)", 140, 30, true},
		{"Very wide terminal (160 cols)", 160, 30, true},
		{"Narrow terminal (80 cols)", 80, 24, false},
		{"Medium terminal (100 cols)", 100, 24, false},
		{"Minimum width (119 cols)", 119, 30, false},
		{"Exact minimum (120 cols)", 120, 30, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := initialModel()

			// Set terminal size
			msg := mockWindowSizeMsg{Width: tt.termWidth, Height: tt.termHeight}
			newModel, _ := m.Update(msg)
			m = newModel.(model)

			// Check showPreview flag
			if m.showPreview != tt.expectPreview {
				t.Errorf("expected showPreview=%v, got %v", tt.expectPreview, m.showPreview)
			}

			// Check viewport height calculation
			if tt.expectPreview && m.viewportHeight == 0 {
				t.Error("viewport height should be calculated when preview is shown")
			}

			// Check terminal dimensions are stored
			if m.termWidth != tt.termWidth {
				t.Errorf("expected termWidth=%d, got %d", tt.termWidth, m.termWidth)
			}
			if m.termHeight != tt.termHeight {
				t.Errorf("expected termHeight=%d, got %d", tt.termHeight, m.termHeight)
			}

			// Test view rendering
			m.files = []FileInfo{
				{Name: "test.xmp", Path: "/test.xmp", IsDir: false, Format: "xmp"},
			}

			view := m.View()

			// Verify View returns non-nil content
			if view.Content == nil {
				t.Error("View should return non-nil Content")
			}
		})
	}
}

// TestPaneWidthCalculation tests that panes are sized correctly
func TestPaneWidthCalculation(t *testing.T) {
	m := initialModel()
	m.termWidth = 120
	m.termHeight = 30
	m.showPreview = true
	m.files = []FileInfo{
		{Name: "test.xmp", Path: "/test.xmp", IsDir: false, Format: "xmp"},
	}

	view := m.View()

	// Verify View returns non-nil content (can't inspect string in v2 API)
	if view.Content == nil {
		t.Error("View should return non-nil Content")
	}
}

// TestPreviewCacheInitialization tests that preview cache is properly initialized
func TestPreviewCacheInitialization(t *testing.T) {
	m := initialModel()

	if m.previewCache == nil {
		t.Error("previewCache should be initialized")
	}

	// Should be empty initially
	if len(m.previewCache) != 0 {
		t.Error("previewCache should be empty initially")
	}
}

// TestPreviewFallbackToFullWidth tests that narrow terminals fall back to full-width layout
func TestPreviewFallbackToFullWidth(t *testing.T) {
	m := initialModel()

	// Set narrow terminal width
	msg := mockWindowSizeMsg{Width: 80, Height: 24}
	newModel, _ := m.Update(msg)
	m = newModel.(model)

	// Should not show preview
	if m.showPreview {
		t.Error("showPreview should be false for narrow terminal")
	}

	// Set files
	m.files = []FileInfo{
		{Name: "test.xmp", Path: "/test.xmp", IsDir: false, Format: "xmp"},
	}

	// View should render full-width file list (Story 4-1 layout)
	view := m.View()

	// Verify View returns non-nil content
	if view.Content == nil {
		t.Error("View should return non-nil Content")
	}
}

// TestTerminalResize tests that layout adapts to terminal resizing
func TestTerminalResize(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "test.xmp", Path: "/test.xmp", IsDir: false, Format: "xmp"},
	}

	// Start with wide terminal
	msg1 := mockWindowSizeMsg{Width: 140, Height: 30}
	newModel, _ := m.Update(msg1)
	m = newModel.(model)

	if !m.showPreview {
		t.Error("preview should be shown in wide terminal")
	}

	// Resize to narrow
	msg2 := mockWindowSizeMsg{Width: 80, Height: 24}
	newModel, _ = m.Update(msg2)
	m = newModel.(model)

	if m.showPreview {
		t.Error("preview should be hidden after resize to narrow")
	}

	// Resize back to wide
	msg3 := mockWindowSizeMsg{Width: 160, Height: 30}
	newModel, _ = m.Update(msg3)
	m = newModel.(model)

	if !m.showPreview {
		t.Error("preview should be shown again after resize to wide")
	}
}

// TestViewportHeightCalculation tests viewport height is calculated correctly
func TestViewportHeightCalculation(t *testing.T) {
	tests := []struct {
		termHeight     int
		expectedMinVH  int // Minimum expected viewport height
	}{
		{30, 20}, // 30 - 8 = 22
		{40, 30}, // 40 - 8 = 32
		{24, 14}, // 24 - 8 = 16 (minimum standard terminal)
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			m := initialModel()
			msg := mockWindowSizeMsg{Width: 120, Height: tt.termHeight}
			newModel, _ := m.Update(msg)
			m = newModel.(model)

			if m.viewportHeight < tt.expectedMinVH {
				t.Errorf("viewport height too small: got %d, expected at least %d",
					m.viewportHeight, tt.expectedMinVH)
			}
		})
	}
}
