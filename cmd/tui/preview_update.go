package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// previewLoadedMsg is sent when preview content has been loaded
type previewLoadedMsg struct {
	filePath string
	content  string
	err      error
}

// loadPreviewCmd loads preview content for a file asynchronously
func loadPreviewCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		// Extract parameters from file
		recipe, err := extractParameters(filePath)
		if err != nil {
			return previewLoadedMsg{
				filePath: filePath,
				content:  "",
				err:      err,
			}
		}

		// Format parameters for display
		content := formatParameters(recipe)

		return previewLoadedMsg{
			filePath: filePath,
			content:  content,
			err:      nil,
		}
	}
}

// updatePreview updates the preview pane for the currently selected file
func (m model) updatePreview() (model, tea.Cmd) {
	// Check if preview is enabled and we have files
	if !m.showPreview || len(m.files) == 0 || m.cursor >= len(m.files) {
		m.previewContent = ""
		m.previewFile = ""
		return m, nil
	}

	// Get currently selected file
	currentFile := m.files[m.cursor]

	// Skip directories
	if currentFile.IsDir {
		m.previewContent = "  (Directories cannot be previewed)"
		m.previewFile = ""
		return m, nil
	}

	// Check cache first
	if cachedContent, ok := m.previewCache[currentFile.Path]; ok {
		m.previewContent = cachedContent
		m.previewFile = currentFile.Path
		m.previewLoading = false
		m.scrollOffset = 0 // Reset scroll on file change
		return m, nil
	}

	// Start loading preview
	m.previewLoading = true
	m.previewFile = currentFile.Path
	m.scrollOffset = 0 // Reset scroll on file change

	return m, loadPreviewCmd(currentFile.Path)
}

// splitLines splits content into lines for scroll calculations
func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	return strings.Split(content, "\n")
}
