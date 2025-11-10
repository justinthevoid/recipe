package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestIntegrationFullWorkflow tests the complete TUI workflow end-to-end
func TestIntegrationFullWorkflow(t *testing.T) {
	// Initialize model
	m := initialModel()

	// Test 1: Terminal resize triggers preview mode
	t.Run("Terminal resize enables preview mode", func(t *testing.T) {
		msg := mockWindowSizeMsg{Width: 120, Height: 30}
		newModel, _ := m.Update(msg)
		m = newModel.(model)

		if !m.showPreview {
			t.Error("expected showPreview to be true for 120 column terminal")
		}
		if m.viewportHeight == 0 {
			t.Error("expected viewportHeight to be calculated")
		}
	})

	// Test 2: Load files from examples directory
	t.Run("Load files from examples directory", func(t *testing.T) {
		examplesDir := filepath.Join("..", "..", "examples", "lrtemplate", "015. PRESETPRO - Emulation K")
		if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
			t.Skip("examples directory not found, skipping integration test")
		}

		m.currentDir = examplesDir
		msg := loadFilesCmd(examplesDir)()
		newModel, _ := m.Update(msg)
		m = newModel.(model)

		if len(m.files) == 0 {
			t.Error("expected to load files from examples directory")
		}

		// Verify we have different file formats
		formats := make(map[string]bool)
		for _, file := range m.files {
			formats[file.Format] = true
		}
		if len(formats) == 0 {
			t.Error("expected to detect file formats")
		}
		t.Logf("Found %d files with formats: %v", len(m.files), formats)
	})

	// Test 3: Navigate through files and trigger preview updates
	t.Run("Navigate through files with preview updates", func(t *testing.T) {
		if len(m.files) == 0 {
			t.Skip("no files loaded, skipping")
		}

		// Move down through files
		for i := 0; i < min(3, len(m.files)-1); i++ {
			newM, cmd := m.handleKeyPress(mockKeyMsg{"j"})
			m = newM.(model)

			// Preview should update
			if cmd == nil && m.showPreview {
				t.Logf("Moved to file %d: %s", m.cursor, m.files[m.cursor].Name)
			}
		}

		// Move back up
		newM, _ := m.handleKeyPress(mockKeyMsg{"k"})
		m = newM.(model)

		if m.cursor < 0 || m.cursor >= len(m.files) {
			t.Error("cursor out of bounds after navigation")
		}
	})

	// Test 4: Test preview loading with real files
	t.Run("Load preview for real preset files", func(t *testing.T) {
		if len(m.files) == 0 {
			t.Skip("no files loaded, skipping")
		}

		// Find a preset file (not directory)
		presetIdx := -1
		for i, file := range m.files {
			if !file.IsDir && (file.Format == "xmp" || file.Format == "np3" || file.Format == "lrtemplate") {
				presetIdx = i
				break
			}
		}

		if presetIdx == -1 {
			t.Skip("no preset files found")
		}

		m.cursor = presetIdx
		file := m.files[presetIdx]
		t.Logf("Testing preview loading for: %s (%s)", file.Name, file.Format)

		// Trigger preview update
		m, cmd := m.updatePreview()

		// Should trigger async load command
		if cmd == nil {
			t.Error("expected command to load preview")
		}

		// Simulate async load completing
		if cmd != nil {
			result := cmd()
			if previewMsg, ok := result.(previewLoadedMsg); ok {
				newModel, _ := m.Update(previewMsg)
				m = newModel.(model)

				// Verify preview was loaded
				if previewMsg.err != nil {
					t.Logf("Preview load error (may be expected): %v", previewMsg.err)
				} else {
					if m.previewContent == "" {
						t.Error("expected preview content to be populated")
					}
					if !strings.Contains(m.previewContent, "Basic Adjustments:") &&
						!strings.Contains(m.previewContent, "Color:") &&
						!strings.Contains(m.previewContent, "No adjustments") {
						t.Logf("Preview content: %s", m.previewContent)
					}
				}

				// Verify caching
				if _, ok := m.previewCache[file.Path]; !ok {
					t.Error("expected content to be cached")
				}
			}
		}
	})

	// Test 5: Test Tab key to toggle preview focus
	t.Run("Toggle preview focus with Tab", func(t *testing.T) {
		if !m.showPreview {
			t.Skip("preview not enabled, skipping")
		}

		initialFocus := m.previewFocused

		// Press Tab
		newM, _ := m.handleKeyPress(mockKeyMsg{"tab"})
		m = newM.(model)

		if m.previewFocused == initialFocus {
			t.Error("expected preview focus to toggle")
		}

		// Press Tab again
		newM, _ = m.handleKeyPress(mockKeyMsg{"tab"})
		m = newM.(model)

		if m.previewFocused != initialFocus {
			t.Error("expected preview focus to toggle back")
		}
	})

	// Test 6: Test scrolling when preview is focused
	t.Run("Scroll preview pane when focused", func(t *testing.T) {
		if !m.showPreview {
			t.Skip("preview not enabled, skipping")
		}

		// Create multi-line content
		lines := []string{}
		for i := 0; i < 30; i++ {
			lines = append(lines, "Line "+string(rune('A'+i%26)))
		}
		m.previewContent = strings.Join(lines, "\n")
		m.previewFocused = true
		m.viewportHeight = 10
		m.scrollOffset = 0

		// Scroll down with j
		newM, _ := m.handleKeyPress(mockKeyMsg{"j"})
		m = newM.(model)

		if m.scrollOffset != 1 {
			t.Errorf("expected scrollOffset to be 1, got %d", m.scrollOffset)
		}

		// Page down
		newM, _ = m.handleKeyPress(mockKeyMsg{"pgdn"})
		m = newM.(model)

		if m.scrollOffset != 11 {
			t.Errorf("expected scrollOffset to be 11 after PageDown, got %d", m.scrollOffset)
		}

		// Page up
		newM, _ = m.handleKeyPress(mockKeyMsg{"pgup"})
		m = newM.(model)

		if m.scrollOffset != 1 {
			t.Errorf("expected scrollOffset to be 1 after PageUp, got %d", m.scrollOffset)
		}
	})

	// Test 7: Test metadata rendering
	t.Run("Render metadata for different file types", func(t *testing.T) {
		if len(m.files) == 0 {
			t.Skip("no files loaded, skipping")
		}

		for _, file := range m.files {
			metadata := renderMetadata(file)

			// Should contain basic metadata
			if file.IsDir {
				if !strings.Contains(metadata, "Type: Directory") {
					t.Errorf("expected directory metadata for %s", file.Name)
				}
			} else {
				if !strings.Contains(metadata, "Format:") {
					t.Errorf("expected format in metadata for %s", file.Name)
				}
				if !strings.Contains(metadata, "Size:") {
					t.Errorf("expected size in metadata for %s", file.Name)
				}
			}

			if !strings.Contains(metadata, "Modified:") {
				t.Errorf("expected modified date in metadata for %s", file.Name)
			}
		}
	})

	// Test 8: Test directory navigation with preview clearing
	t.Run("Navigate directories clears preview", func(t *testing.T) {
		// Set some preview content
		m.previewContent = "Test content"
		m.previewFile = "test.xmp"

		// Navigate up
		originalDir := m.currentDir
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			newM, _ := m.navigateUp()
			m = newM.(model)

			// Preview should be cleared
			if m.previewContent != "" {
				t.Error("expected preview content to be cleared after directory navigation")
			}
			if m.previewFile != "" {
				t.Error("expected preview file to be cleared after directory navigation")
			}

			// Navigate back
			m.currentDir = originalDir
			m.Update(loadFilesCmd(originalDir)())
		}
	})

	// Test 9: Test cache performance
	t.Run("Cache improves preview performance", func(t *testing.T) {
		if len(m.files) == 0 {
			t.Skip("no files loaded, skipping")
		}

		// Find a preset file
		presetIdx := -1
		for i, file := range m.files {
			if !file.IsDir && file.Format == "xmp" {
				presetIdx = i
				break
			}
		}

		if presetIdx == -1 {
			t.Skip("no XMP file found")
		}

		m.cursor = presetIdx
		file := m.files[presetIdx]

		// First load (cache miss)
		m.previewCache = make(map[string]string)
		start := time.Now()
		m, cmd := m.updatePreview()
		firstLoadTime := time.Since(start)

		// Execute async command if returned
		if cmd != nil {
			result := cmd()
			if previewMsg, ok := result.(previewLoadedMsg); ok {
				newModel, _ := m.Update(previewMsg)
				m = newModel.(model)
			}
		}

		// Second load (cache hit)
		start = time.Now()
		m, cmd = m.updatePreview()
		secondLoadTime := time.Since(start)

		// Cache hit should return immediately (no command)
		if cmd != nil {
			t.Error("expected no command for cached preview")
		}

		t.Logf("First load: %v, Second load: %v", firstLoadTime, secondLoadTime)

		// Verify cache was used
		if _, ok := m.previewCache[file.Path]; !ok {
			t.Error("expected content to be in cache")
		}
	})
}

// TestIntegrationFormatDetection tests format detection for all supported types
func TestIntegrationFormatDetection(t *testing.T) {
	formats := []string{".xmp", ".np3", ".lrtemplate"}

	for _, ext := range formats {
		t.Run("Format detection for "+ext, func(t *testing.T) {
			file := FileInfo{
				Name:    "test" + ext,
				Path:    "/test/test" + ext,
				Format:  strings.TrimPrefix(ext, "."),
				Size:    1024,
				ModTime: time.Now(),
			}

			metadata := renderMetadata(file)

			expectedFormat := strings.ToUpper(strings.TrimPrefix(ext, "."))
			if !strings.Contains(metadata, "Format: "+expectedFormat) {
				t.Errorf("expected format %s in metadata, got: %s", expectedFormat, metadata)
			}
		})
	}
}

// TestIntegrationEdgeCases tests edge cases and error handling
func TestIntegrationEdgeCases(t *testing.T) {
	t.Run("Empty file list", func(t *testing.T) {
		m := initialModel()
		m.files = []FileInfo{}
		m.showPreview = true

		// Try to update preview with no files
		m, cmd := m.updatePreview()

		if cmd != nil {
			t.Error("expected no command for empty file list")
		}
		if m.previewContent != "" {
			t.Error("expected empty preview content")
		}
	})

	t.Run("Preview when disabled", func(t *testing.T) {
		m := initialModel()
		m.files = []FileInfo{{Name: "test.xmp", Path: "/test.xmp", Format: "xmp"}}
		m.cursor = 0
		m.showPreview = false // Disabled

		// Try to update preview when disabled
		m, cmd := m.updatePreview()

		if cmd != nil {
			t.Error("expected no command when preview disabled")
		}
	})

	t.Run("Cursor out of bounds", func(t *testing.T) {
		m := initialModel()
		m.files = []FileInfo{{Name: "test.xmp", Path: "/test.xmp", Format: "xmp"}}
		m.cursor = 10 // Out of bounds
		m.showPreview = true

		// Try to update preview with invalid cursor
		m, cmd := m.updatePreview()

		if cmd != nil {
			t.Error("expected no command for out-of-bounds cursor")
		}
	})
}
