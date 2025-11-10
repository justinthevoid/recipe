package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Color styles for different file formats
var (
	styleNP3    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))   // Blue
	styleXMP    = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))  // Orange
	styleLRT    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))   // Green
	styleCOS    = lipgloss.NewStyle().Foreground(lipgloss.Color("135"))  // Purple (Capture One brand color)
	styleCOSP   = lipgloss.NewStyle().Foreground(lipgloss.Color("135"))  // Purple (Capture One bundle)
	styleDir    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))  // Gray
	styleCursor = lipgloss.NewStyle().Reverse(true)                      // Reverse colors
	styleBorder = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))  // Gray border
	styleHeader = lipgloss.NewStyle().Bold(true)                         // Bold header
	styleStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))  // Light gray status
)

// renderFileList renders the main file browser view
func renderFileList(m model) string {
	var b strings.Builder

	// Header: Current directory
	header := styleHeader.Render(fmt.Sprintf("Current: %s", m.currentDir))
	b.WriteString(header)
	b.WriteString("\n")

	// Border separator
	border := styleBorder.Render(strings.Repeat("─", min(m.termWidth, 80)))
	b.WriteString(border)
	b.WriteString("\n")

	// Handle empty directory
	if len(m.files) == 0 {
		b.WriteString("\n")
		b.WriteString("  No preset files found in this directory.\n")
		b.WriteString("  (Showing only .np3, .xmp, .lrtemplate, .costyle, and .costylepack files)\n")
		b.WriteString("\n")
	} else {
		// Calculate visible range for scrolling
		maxVisible := m.termHeight - 6 // Reserve space for header, footer, status
		start := 0
		end := len(m.files)

		if len(m.files) > maxVisible && maxVisible > 0 {
			// Scroll to keep cursor visible
			if m.cursor >= maxVisible/2 {
				start = m.cursor - maxVisible/2
			}
			end = start + maxVisible
			if end > len(m.files) {
				end = len(m.files)
				start = max(0, end-maxVisible)
			}
		}

		// Render visible files
		for i := start; i < end; i++ {
			line := renderFileLine(m.files[i], i == m.cursor, m.selected[m.files[i].Path])
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if start > 0 || end < len(m.files) {
			scrollInfo := styleStatus.Render(fmt.Sprintf("  [Showing %d-%d of %d files]", start+1, end, len(m.files)))
			b.WriteString(scrollInfo)
			b.WriteString("\n")
		}
	}

	// Status bar
	b.WriteString("\n")
	selectedCount := len(m.selected)
	if selectedCount > 0 {
		statusLine := fmt.Sprintf("[%d files selected] Press 'c' to convert (Story 4-3)", selectedCount)
		b.WriteString(styleStatus.Render(statusLine))
		b.WriteString("\n")
	}

	// Show current file info if cursor is valid
	if m.cursor < len(m.files) && len(m.files) > 0 {
		file := m.files[m.cursor]
		fileInfo := fmt.Sprintf("%s | %s", file.Path, formatSize(file.Size))
		if !file.IsDir {
			fileInfo += fmt.Sprintf(" | Modified: %s", file.ModTime.Format("2006-01-02"))
		}
		b.WriteString(styleStatus.Render(fileInfo))
		b.WriteString("\n")
	}

	// Help hint
	b.WriteString(styleStatus.Render("Press '?' for help | 'q' to quit"))
	b.WriteString("\n")

	return b.String()
}

// renderFileLine renders a single file or directory line
func renderFileLine(file FileInfo, isCursor bool, isSelected bool) string {
	// Cursor indicator
	cursor := " "
	if isCursor {
		cursor = ">"
	}

	// Selection checkbox
	checkbox := " "
	if isSelected {
		checkbox = "✓"
	}

	// Icon
	icon := "📄"
	if file.IsDir {
		icon = "📁"
	}

	// Format badge
	badge := formatBadge(file.Format)

	// Size (only for files)
	size := ""
	if !file.IsDir {
		size = formatSize(file.Size)
	}

	// Build line
	line := fmt.Sprintf("%s %s %s %-40s %3s %8s", cursor, checkbox, icon, file.Name, badge, size)

	// Apply cursor style if needed
	if isCursor {
		return styleCursor.Render(line)
	}

	return line
}

// formatBadge returns a color-coded format badge
func formatBadge(format string) string {
	switch format {
	case "np3":
		return styleNP3.Render("NP3")
	case "xmp":
		return styleXMP.Render("XMP")
	case "lrtemplate":
		return styleLRT.Render("LRT")
	case "costyle":
		return styleCOS.Render("COS")
	case "costylepack":
		return styleCOSP.Render("CPK")
	case "dir":
		return styleDir.Render("DIR")
	default:
		return "   "
	}
}

// formatSize converts bytes to human-readable format
func formatSize(bytes int64) string {
	const kb = 1024
	const mb = kb * 1024

	if bytes < kb {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < mb {
		return fmt.Sprintf("%.1f KB", float64(bytes)/kb)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/mb)
	}
}

// renderHelp renders the help overlay
func renderHelp(m model) string {
	help := `
┌─────────────────── Keyboard Shortcuts ────────────────────┐
│                                                            │
│ Navigation:                                                │
│   ↑/k        Move cursor up                                │
│   ↓/j        Move cursor down                              │
│   Enter      Navigate into directory                       │
│   Backspace  Navigate to parent directory                  │
│   Home       Jump to first item                            │
│   End        Jump to last item                             │
│                                                            │
│ Preview Pane (when terminal ≥ 120 columns):                │
│   Tab        Toggle focus between file list and preview    │
│   j/k        Scroll preview line by line (when focused)    │
│   PgUp/PgDn  Scroll preview page by page (when focused)    │
│                                                            │
│ Selection:                                                 │
│   Space      Toggle selection on current file              │
│   a          Select all files                              │
│   n          Deselect all files                            │
│                                                            │
│ Actions:                                                   │
│   c          Convert selected files (Story 4-3)            │
│   r          Refresh current directory                     │
│   ?          Toggle this help                              │
│   q/Ctrl+C   Quit                                          │
│                                                            │
│ Press ? or Esc to close                                    │
└────────────────────────────────────────────────────────────┘
`
	return help
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
