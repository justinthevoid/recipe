package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// openSettingsEditor opens the settings editor from validation screen
func (m model) openSettingsEditor() (model, tea.Cmd) {
	m.showSettingsEditor = true
	m.editorCursor = 0

	// Copy current settings to temporary fields
	m.editedTargetFormat = m.targetFormat
	m.editedOutputDir = m.outputDir
	m.editedFileSelection = make(map[string]bool)
	for k, v := range m.selected {
		m.editedFileSelection[k] = v
	}

	return m, nil
}

// closeSettingsEditor closes the settings editor without saving
func (m model) closeSettingsEditor() (model, tea.Cmd) {
	m.showSettingsEditor = false
	m.editorCursor = 0
	m.editedTargetFormat = ""
	m.editedOutputDir = ""
	m.editedFileSelection = nil

	return m, nil
}

// saveSettings saves edited settings and re-runs validation
func (m model) saveSettings() (model, tea.Cmd) {
	// Apply edited settings
	m.targetFormat = m.editedTargetFormat
	m.outputDir = m.editedOutputDir
	m.selected = make(map[string]bool)
	for k, v := range m.editedFileSelection {
		m.selected[k] = v
	}

	// Close editor
	m.showSettingsEditor = false
	m.editorCursor = 0
	m.editedTargetFormat = ""
	m.editedOutputDir = ""
	m.editedFileSelection = nil

	// Re-run validation with new settings
	m = m.performValidation()

	return m, nil
}

// renderSettingsEditor renders the settings editor screen
func renderSettingsEditor(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("╔════════════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║ Edit Conversion Settings                                               ║\n")
	b.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	b.WriteString("║                                                                        ║\n")

	// Target Format (field 0)
	cursor0 := " "
	if m.editorCursor == 0 {
		cursor0 = ">"
	}
	b.WriteString("║ Target Format:                                                         ║\n")

	formats := []struct {
		key  string
		name string
	}{
		{"np3", "NP3 (Nikon NX Studio)"},
		{"xmp", "XMP (Adobe Lightroom)"},
		{"lrtemplate", "lrtemplate (Lightroom Template)"},
		{"costyle", "Costyle (Capture One)"},
		{"costylepack", "Costylepack (Capture One Bundle)"},
	}

	for _, format := range formats {
		selected := " "
		if m.editedTargetFormat == format.key {
			selected = "●"
		}
		line := fmt.Sprintf("%s  %s %s", cursor0, selected, format.name)
		b.WriteString(fmt.Sprintf("║ %-70s ║\n", line))
	}

	b.WriteString("║                                                                        ║\n")

	// Output Directory (field 1)
	cursor1 := " "
	if m.editorCursor == 1 {
		cursor1 = ">"
	}
	b.WriteString("║ Output Directory:                                                      ║\n")
	dirLine := fmt.Sprintf("%s  [%-65s]", cursor1, truncateString(m.editedOutputDir, 65))
	b.WriteString(fmt.Sprintf("║ %-70s ║\n", dirLine))

	b.WriteString("║                                                                        ║\n")

	// Files to Convert (field 2)
	cursor2 := " "
	if m.editorCursor == 2 {
		cursor2 = ">"
	}
	b.WriteString("║ Files to Convert: (space to toggle)                                   ║\n")

	// Show first 8 files
	fileIndex := 0
	for _, file := range m.files {
		if _, exists := m.editedFileSelection[file.Path]; exists {
			checked := "☐"
			if m.editedFileSelection[file.Path] {
				checked = "✓"
			}
			line := fmt.Sprintf("%s  %s %s", cursor2, checked, truncateString(file.Name, 60))
			b.WriteString(fmt.Sprintf("║ %-70s ║\n", line))
			fileIndex++
			if fileIndex >= 8 {
				break
			}
		}
	}

	b.WriteString("║                                                                        ║\n")

	// Footer with shortcuts
	b.WriteString("║ ↑/↓ Navigate fields  |  Space: Toggle (format/file)                  ║\n")
	b.WriteString("║ Enter: Save and return  |  Esc: Cancel changes                        ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════════════════╝\n")

	return b.String()
}
