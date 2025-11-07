package main

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

// KeyMsgInterface for testing
type KeyMsgInterface interface {
	String() string
}

// handleKeyPress processes keyboard input
func (m model) handleKeyPress(msg KeyMsgInterface) (tea.Model, tea.Cmd) {
	// Handle any-key dismissal of summary screen
	if m.showSummary {
		m.showSummary = false
		m.results = nil
		m.completedList = nil
		// Clear selection after conversion
		m.selected = make(map[string]bool)
		return m, nil
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "ctrl+c":
		// Handle cancellation during conversion
		if m.converting && !m.cancelling {
			m.cancelling = true
			// Send cancel signal
			select {
			case m.cancelChan <- true:
			default:
			}
			return m, nil
		}
		return m, tea.Quit

	case "up", "k":
		// Settings editor navigation
		if m.showSettingsEditor {
			if m.editorCursor > 0 {
				m.editorCursor--
			}
			return m, nil
		}

		// If preview is focused, scroll preview pane instead of moving cursor
		if m.showPreview && m.previewFocused && msg.String() == "k" {
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
			return m, nil
		}

		if m.cursor > 0 {
			m.cursor--
			// Update preview for new file
			return m.updatePreview()
		}
		return m, nil

	case "down", "j":
		// Settings editor navigation
		if m.showSettingsEditor {
			if m.editorCursor < 2 { // 0=format, 1=directory, 2=files
				m.editorCursor++
			}
			return m, nil
		}

		// If preview is focused, scroll preview pane instead of moving cursor
		if m.showPreview && m.previewFocused && msg.String() == "j" {
			// Calculate max scroll
			lines := len(splitLines(m.previewContent))
			maxScroll := lines - m.viewportHeight
			if maxScroll < 0 {
				maxScroll = 0
			}
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
			return m, nil
		}

		if m.cursor < len(m.files)-1 {
			m.cursor++
			// Update preview for new file
			return m.updatePreview()
		}
		return m, nil

	case "home":
		m.cursor = 0
		return m.updatePreview()

	case "end":
		if len(m.files) > 0 {
			m.cursor = len(m.files) - 1
		}
		return m.updatePreview()

	case "tab":
		// Toggle preview focus (for keyboard scrolling)
		if m.showPreview {
			m.previewFocused = !m.previewFocused
		}
		return m, nil

	case "pgup":
		// Page up in preview pane
		if m.showPreview && m.previewFocused {
			pageSize := m.viewportHeight
			if m.scrollOffset > pageSize {
				m.scrollOffset -= pageSize
			} else {
				m.scrollOffset = 0
			}
		}
		return m, nil

	case "pgdn":
		// Page down in preview pane
		if m.showPreview && m.previewFocused {
			lines := len(splitLines(m.previewContent))
			maxScroll := lines - m.viewportHeight
			if maxScroll < 0 {
				maxScroll = 0
			}
			pageSize := m.viewportHeight
			if m.scrollOffset+pageSize < maxScroll {
				m.scrollOffset += pageSize
			} else {
				m.scrollOffset = maxScroll
			}
		}
		return m, nil

	case "enter":
		// Settings editor: save and return
		if m.showSettingsEditor {
			return m.saveSettings()
		}
		return m.navigateInto()

	case "backspace", "left":
		return m.navigateUp()

	case " ", "space":
		// Settings editor: toggle format or file selection
		if m.showSettingsEditor {
			if m.editorCursor == 0 {
				// Cycle through formats
				formats := []string{"np3", "xmp", "lrtemplate"}
				for i, f := range formats {
					if m.editedTargetFormat == f {
						m.editedTargetFormat = formats[(i+1)%len(formats)]
						break
					}
				}
			} else if m.editorCursor == 2 {
				// Toggle file selection (simplified - toggle first selected file)
				for path := range m.editedFileSelection {
					m.editedFileSelection[path] = !m.editedFileSelection[path]
					break
				}
			}
			return m, nil
		}
		return m.toggleSelection()

	case "a":
		return m.selectAll()

	case "n":
		// Cancel conversion if in confirmation mode, otherwise deselect all
		if m.showConfirmation {
			m.showConfirmation = false
			return m, nil
		}
		return m.deselectAll()

	case "r":
		return m, loadFilesCmd(m.currentDir)

	case "?":
		m.showHelp = !m.showHelp
		return m, nil

	case "esc":
		if m.showHelp {
			m.showHelp = false
		} else if m.showSettingsEditor {
			// Cancel settings editor without saving (Story 4-4)
			return m.closeSettingsEditor()
		} else if m.showFormatPrompt {
			// Cancel format selection
			m.showFormatPrompt = false
			m.formatMenuCursor = 0
		} else if m.showValidation {
			// Cancel validation (Story 4-4)
			m.showValidation = false
			m.validationFiles = nil
			m.validationWarnings = nil
		} else if m.showConfirmation {
			// Cancel conversion
			m.showConfirmation = false
		}
		return m, nil

	case "c":
		// Story 4-4: 'c' on validation screen moves to confirmation
		if m.showValidation && m.validationPassed {
			m.showValidation = false
			m.showConfirmation = true
			return m, nil
		}
		// Trigger batch conversion if files are selected
		return m.triggerConversion()

	case "e":
		// Story 4-4: 'e' opens settings editor on validation screen
		if m.showValidation {
			return m.openSettingsEditor()
		}
		return m, nil

	case "1", "2", "3":
		// Handle format menu selection
		if m.showFormatPrompt {
			return m.selectFormat(msg.String())
		}
		return m, nil

	case "y":
		// Confirm conversion
		if m.showConfirmation {
			return m.startConversion()
		}
		return m, nil
	}

	return m, nil
}

// navigateInto navigates into a directory
func (m model) navigateInto() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.files) || len(m.files) == 0 {
		return m, nil
	}

	selected := m.files[m.cursor]
	if !selected.IsDir {
		return m, nil // Can't navigate into file
	}

	m.currentDir = selected.Path
	m.cursor = 0
	m.previewContent = "" // Clear preview when changing directories
	m.previewFile = ""
	return m, loadFilesCmd(m.currentDir)
}

// navigateUp navigates to parent directory
func (m model) navigateUp() (tea.Model, tea.Cmd) {
	parent := filepath.Dir(m.currentDir)
	if parent == m.currentDir {
		return m, nil // Already at root
	}

	m.currentDir = parent
	m.cursor = 0
	m.previewContent = "" // Clear preview when changing directories
	m.previewFile = ""
	return m, loadFilesCmd(m.currentDir)
}

// toggleSelection toggles selection on current file
func (m model) toggleSelection() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.files) || len(m.files) == 0 {
		return m, nil
	}

	file := m.files[m.cursor]
	if file.IsDir {
		return m, nil // Can't select directories
	}

	if m.selected[file.Path] {
		delete(m.selected, file.Path)
	} else {
		m.selected[file.Path] = true
	}

	return m, nil
}

// selectAll selects all files in current directory
func (m model) selectAll() (tea.Model, tea.Cmd) {
	for _, file := range m.files {
		if !file.IsDir {
			m.selected[file.Path] = true
		}
	}
	return m, nil
}

// deselectAll clears all selections
func (m model) deselectAll() (tea.Model, tea.Cmd) {
	m.selected = make(map[string]bool)
	return m, nil
}
