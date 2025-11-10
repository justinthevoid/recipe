package main

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// triggerConversion initiates the conversion process
func (m model) triggerConversion() (tea.Model, tea.Cmd) {
	// Check if any files are selected
	if len(m.selected) == 0 {
		return m, nil // Do nothing if no files selected
	}

	// Show format selection prompt
	m.showFormatPrompt = true
	m.formatMenuCursor = 0
	m.outputDir = m.currentDir // Default to current directory

	return m, nil
}

// selectFormat handles format menu selection
func (m model) selectFormat(choice string) (tea.Model, tea.Cmd) {
	switch choice {
	case "1":
		m.targetFormat = "np3"
	case "2":
		m.targetFormat = "xmp"
	case "3":
		m.targetFormat = "lrtemplate"
	default:
		return m, nil
	}

	// Move to validation screen (Story 4-4)
	m.showFormatPrompt = false
	newModel, cmd := m.triggerValidation()
	return newModel, cmd
}

// startConversion begins the batch conversion process
func (m model) startConversion() (tea.Model, tea.Cmd) {
	// Hide confirmation screen
	m.showConfirmation = false

	// Prepare for conversion
	m.converting = true
	m.currentFile = 0
	m.completedFiles = 0
	m.errorCount = 0
	m.warningCount = 0
	m.startTime = timeNow() // Use timeNow() for testability
	m.results = make([]ConversionResult, 0)
	m.completedList = make([]string, 0)
	m.cancelling = false

	// Get list of selected files
	selectedFiles := make([]FileInfo, 0, len(m.selected))
	for _, file := range m.files {
		if m.selected[file.Path] {
			selectedFiles = append(selectedFiles, file)
		}
	}

	m.totalFiles = len(selectedFiles)

	// Start conversion in background and tick for progress updates
	return m, tea.Batch(
		convertBatchCmd(selectedFiles, m.targetFormat, m.outputDir, m.cancelChan),
		tickCmd(),
	)
}

// formatMenuOptions returns the format selection menu text
func formatMenuOptions() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("Select Target Format:\n")
	b.WriteString("\n")
	b.WriteString("  1. NP3 (Nikon NX Studio)\n")
	b.WriteString("  2. XMP (Adobe Lightroom)\n")
	b.WriteString("  3. lrtemplate (Lightroom Template)\n")
	b.WriteString("\n")
	b.WriteString("Enter choice (1-3) or Esc to cancel:")
	return b.String()
}

// confirmationScreen returns the confirmation screen text
func confirmationScreen(m model) string {
	var b strings.Builder

	// Count source formats
	formatCounts := make(map[string]int)
	for path := range m.selected {
		for _, file := range m.files {
			if file.Path == path {
				formatCounts[file.Format]++
				break
			}
		}
	}

	// Build format breakdown string
	formatBreakdown := make([]string, 0, len(formatCounts))
	for format, count := range formatCounts {
		formatName := strings.ToUpper(format)
		if format == "lrtemplate" {
			formatName = "lrtemplate"
		}
		formatBreakdown = append(formatBreakdown, fmt.Sprintf("%s (%d)", formatName, count))
	}

	targetFormatName := strings.ToUpper(m.targetFormat)
	if m.targetFormat == "lrtemplate" {
		targetFormatName = "lrtemplate (Lightroom Template)"
	} else if m.targetFormat == "np3" {
		targetFormatName = "NP3 (Nikon NX Studio)"
	} else if m.targetFormat == "xmp" {
		targetFormatName = "XMP (Adobe Lightroom)"
	}

	b.WriteString("\n")
	b.WriteString("Confirm Batch Conversion:\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Files:          %d selected\n", len(m.selected)))
	b.WriteString(fmt.Sprintf("Source formats: %s\n", strings.Join(formatBreakdown, ", ")))
	b.WriteString(fmt.Sprintf("Target format:  %s\n", targetFormatName))
	b.WriteString(fmt.Sprintf("Output dir:     %s\n", m.outputDir))
	b.WriteString("\n")
	b.WriteString("Proceed? [y/n]:")

	return b.String()
}

// getSelectedFiles returns a slice of selected file paths
func (m model) getSelectedFiles() []string {
	files := make([]string, 0, len(m.selected))
	for path := range m.selected {
		files = append(files, filepath.Base(path))
	}
	return files
}
