package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Styles for preview pane
var (
	stylePreviewHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))  // Blue
	stylePreviewLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))            // Light gray
	stylePreviewValue  = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))            // White
	stylePreviewError  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))            // Red
	stylePreviewWarn   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))            // Orange/Yellow
	styleDivider       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Gray divider
)

// renderSplitView renders the split-pane layout with file list on left and preview on right
func renderSplitView(m model) string {
	// Calculate pane widths (50/50 split with 1 column for divider)
	leftWidth := (m.termWidth - 1) / 2
	rightWidth := m.termWidth - leftWidth - 1

	// Render left pane (file list)
	leftPane := renderFileListPane(m, leftWidth)

	// Render right pane (preview)
	rightPane := renderPreviewPane(m, rightWidth)

	// Join panes horizontally with vertical divider
	divider := renderVerticalDivider(m.termHeight)

	// Combine all parts
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, divider, rightPane)
}

// renderFileListPane renders the left pane with file list (reuses existing logic)
func renderFileListPane(m model, width int) string {
	var b strings.Builder

	// Header: Current directory
	header := styleHeader.Render(fmt.Sprintf("Current: %s", truncateString(m.currentDir, width-10)))
	b.WriteString(header)
	b.WriteString("\n")

	// Border separator
	border := styleBorder.Render(strings.Repeat("─", width))
	b.WriteString(border)
	b.WriteString("\n")

	// Handle empty directory
	if len(m.files) == 0 {
		b.WriteString("\n")
		b.WriteString("  No preset files\n")
		b.WriteString("  found.\n")
		b.WriteString("\n")
	} else {
		// Calculate visible range for scrolling
		maxVisible := m.termHeight - 8 // Reserve more space for split view
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

		// Render visible files (compact format for split view)
		for i := start; i < end; i++ {
			line := renderFileLineCompact(m.files[i], i == m.cursor, m.selected[m.files[i].Path], width)
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if start > 0 || end < len(m.files) {
			scrollInfo := styleStatus.Render(fmt.Sprintf("  [%d-%d of %d]", start+1, end, len(m.files)))
			b.WriteString(scrollInfo)
			b.WriteString("\n")
		}
	}

	// Status bar
	b.WriteString("\n")
	selectedCount := len(m.selected)
	if selectedCount > 0 {
		statusLine := fmt.Sprintf("[%d selected]", selectedCount)
		b.WriteString(styleStatus.Render(statusLine))
		b.WriteString("\n")
	}

	return lipgloss.NewStyle().Width(width).Height(m.termHeight).Render(b.String())
}

// renderFileLineCompact renders a compact file line for split view
func renderFileLineCompact(file FileInfo, isCursor bool, isSelected bool, width int) string {
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

	// Truncate filename to fit width
	maxNameLen := width - 15 // Reserve space for cursor, checkbox, icon, badge
	name := truncateString(file.Name, maxNameLen)

	// Build line
	line := fmt.Sprintf("%s %s %s %-*s %s", cursor, checkbox, icon, maxNameLen, name, badge)

	// Apply cursor style if needed
	if isCursor {
		return styleCursor.Render(line)
	}

	return line
}

// renderPreviewPane renders the right pane with parameter preview
func renderPreviewPane(m model, width int) string {
	var b strings.Builder

	// Header with filename
	if m.cursor < len(m.files) && len(m.files) > 0 {
		file := m.files[m.cursor]
		header := stylePreviewHeader.Render(fmt.Sprintf("Preview: %s", truncateString(file.Name, width-10)))
		b.WriteString(header)
	} else {
		header := stylePreviewHeader.Render("Preview")
		b.WriteString(header)
	}
	b.WriteString("\n")

	// Metadata (format, size, modified date)
	if m.cursor < len(m.files) && len(m.files) > 0 {
		file := m.files[m.cursor]
		metadata := renderMetadata(file)
		b.WriteString(metadata)
		b.WriteString("\n")
	}

	// Border separator
	border := styleBorder.Render(strings.Repeat("─", width))
	b.WriteString(border)
	b.WriteString("\n")

	// Content area
	if m.previewLoading {
		// Show loading indicator
		b.WriteString("\n")
		b.WriteString(stylePreviewLabel.Render("  Loading..."))
		b.WriteString("\n")
	} else if m.previewContent != "" {
		// Show preview content with scrolling
		lines := strings.Split(m.previewContent, "\n")
		totalLines := len(lines)

		// Calculate visible slice
		start := m.scrollOffset
		end := start + m.viewportHeight
		if end > totalLines {
			end = totalLines
		}
		if start < 0 {
			start = 0
		}

		// Render visible lines
		for i := start; i < end; i++ {
			if i < len(lines) {
				b.WriteString(lines[i])
				b.WriteString("\n")
			}
		}

		// Show scroll indicator if needed
		if totalLines > m.viewportHeight {
			b.WriteString("\n")
			scrollInfo := styleStatus.Render(fmt.Sprintf("  [Lines %d-%d of %d] j/k to scroll", start+1, end, totalLines))
			b.WriteString(scrollInfo)
			b.WriteString("\n")
		}

		// Footer note about unmapped parameters (AC-6)
		if m.cursor < len(m.files) && len(m.files) > 0 {
			file := m.files[m.cursor]
			if !file.IsDir {
				b.WriteString("\n")
				note := styleStatus.Render("  Note: Format-specific features may not be displayed")
				b.WriteString(note)
			}
		}
	} else {
		// No content to show
		b.WriteString("\n")
		b.WriteString(stylePreviewLabel.Render("  Select a file to preview"))
		b.WriteString("\n")
	}

	return lipgloss.NewStyle().Width(width).Height(m.termHeight).Render(b.String())
}

// renderVerticalDivider renders a vertical line divider
func renderVerticalDivider(height int) string {
	var b strings.Builder
	for i := 0; i < height; i++ {
		b.WriteString(styleDivider.Render("│"))
		b.WriteString("\n")
	}
	return b.String()
}

// renderMetadata formats file metadata (format, size, modified date)
func renderMetadata(file FileInfo) string {
	var parts []string

	// Format
	if file.IsDir {
		parts = append(parts, "Type: Directory")
	} else {
		formatName := strings.ToUpper(file.Format)
		if formatName == "" {
			formatName = "UNKNOWN"
		}
		parts = append(parts, fmt.Sprintf("Format: %s", formatName))

		// File size (human-readable)
		size := formatFileSize(file.Size)
		parts = append(parts, fmt.Sprintf("Size: %s", size))
	}

	// Modified date
	modTime := file.ModTime.Format("2006-01-02 15:04")
	parts = append(parts, fmt.Sprintf("Modified: %s", modTime))

	// Join with separator and apply dimmed style
	metadata := strings.Join(parts, " • ")
	return styleStatus.Render(fmt.Sprintf("  %s", metadata))
}

// formatFileSize formats bytes into human-readable size
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// truncateString truncates a string to maxLen with ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
