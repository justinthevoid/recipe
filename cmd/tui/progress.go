package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/justin/recipe/internal/converter"
)

// Color styles for progress display
var (
	styleBlue   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))   // In-progress
	styleGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))   // Success
	styleYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))  // Warning
	styleRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))  // Error
	styleGray   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))  // Dimmed
)

// Message types
type tickMsg time.Time
type conversionCompleteMsg struct {
	result ConversionResult
}
type batchCompleteMsg struct{}

// timeNow is a variable for testing
var timeNow = time.Now

// tickCmd returns a command that sends tick messages every 100ms
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// convertBatchCmd converts files in parallel using worker pool pattern
func convertBatchCmd(files []FileInfo, targetFormat string, outputDir string, cancelChan chan bool) tea.Cmd {
	return func() tea.Msg {
		numWorkers := runtime.NumCPU()
		if numWorkers > len(files) {
			numWorkers = len(files)
		}

		jobs := make(chan FileInfo, len(files))
		results := make(chan ConversionResult, len(files))

		// Context for cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Listen for cancel signal
		go func() {
			<-cancelChan
			cancel()
		}()

		// Start workers
		var wg sync.WaitGroup
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for file := range jobs {
					// Check for cancellation
					select {
					case <-ctx.Done():
						return
					default:
					}

					result := convertFile(file, targetFormat, outputDir)
					results <- result
				}
			}()
		}

		// Send jobs
		for _, file := range files {
			jobs <- file
		}
		close(jobs)

		// Wait for workers and close results
		go func() {
			wg.Wait()
			close(results)
		}()

		// Send results back to UI
		// Note: This is a simplified version. In production, we'd send
		// individual results as they complete for real-time updates
		for range results {
			// Results are collected but batch completes as a whole for now
		}

		return batchCompleteMsg{}
	}
}

// convertFile converts a single file
func convertFile(file FileInfo, targetFormat string, outputDir string) ConversionResult {
	// Read input file
	input, err := os.ReadFile(file.Path)
	if err != nil {
		slog.Error("Failed to read file", "file", file.Path, "error", err)
		return ConversionResult{
			File:         file.Name,
			Status:       "error",
			Message:      fmt.Sprintf("Read error: %v", err),
			SourceFormat: file.Format,
			TargetFormat: targetFormat,
		}
	}

	// Convert using converter from Epic 1
	output, err := converter.Convert(input, file.Format, targetFormat)
	if err != nil {
		slog.Error("Conversion failed", "file", file.Path, "error", err)
		return ConversionResult{
			File:         file.Name,
			Status:       "error",
			Message:      fmt.Sprintf("Conversion error: %v", err),
			SourceFormat: file.Format,
			TargetFormat: targetFormat,
		}
	}

	// Determine output filename
	baseName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
	var outputExt string
	switch targetFormat {
	case "np3":
		outputExt = ".np3"
	case "xmp":
		outputExt = ".xmp"
	case "lrtemplate":
		outputExt = ".lrtemplate"
	}
	outputPath := filepath.Join(outputDir, baseName+outputExt)

	// Write output atomically
	if err := writeOutputAtomic(outputPath, output); err != nil {
		slog.Error("Failed to write output", "file", outputPath, "error", err)
		return ConversionResult{
			File:         file.Name,
			Status:       "error",
			Message:      fmt.Sprintf("Write error: %v", err),
			SourceFormat: file.Format,
			TargetFormat: targetFormat,
		}
	}

	// Success
	return ConversionResult{
		File:         file.Name,
		Status:       "success",
		Message:      "",
		SourceFormat: file.Format,
		TargetFormat: targetFormat,
	}
}

// writeOutputAtomic writes a file atomically using temp file + rename
func writeOutputAtomic(path string, data []byte) error {
	// Write to temp file
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath) // Cleanup on failure
		return err
	}

	return nil
}

// renderProgressBar renders the progress bar
func renderProgressBar(current, total, errorCount int) string {
	if total == 0 {
		return ""
	}

	percentage := (current * 100) / total
	barWidth := 40

	// Calculate filled blocks
	filled := (current * barWidth) / total
	if filled > barWidth {
		filled = barWidth
	}

	// Build bar
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Color based on status
	var styledBar string
	if errorCount > 0 {
		styledBar = styleRed.Render(bar)
	} else if current == total {
		styledBar = styleGreen.Render(bar)
	} else {
		styledBar = styleBlue.Render(bar)
	}

	// Format display
	display := fmt.Sprintf("Overall Progress: %d%% (%d/%d)", percentage, current, total)
	if errorCount > 0 {
		display += fmt.Sprintf(", %d errors", errorCount)
	}

	return display + "\n" + styledBar
}

// formatFileStatus formats per-file status line
func formatFileStatus(file, sourceFormat, targetFormat, status string) string {
	// Icon and color based on status
	var icon string
	var style lipgloss.Style
	switch status {
	case "converting":
		icon = "⠋" // Spinner frame
		style = styleBlue
	case "success":
		icon = "✓"
		style = styleGreen
	case "warning":
		icon = "⚠️"
		style = styleYellow
	case "error":
		icon = "✗"
		style = styleRed
	default:
		icon = " "
		style = lipgloss.NewStyle()
	}

	// Format conversion
	conversion := fmt.Sprintf("%s → %s", strings.ToUpper(sourceFormat), strings.ToUpper(targetFormat))

	// Build status line
	statusText := ""
	switch status {
	case "converting":
		statusText = "Converting..."
	case "success":
		statusText = "Success"
	case "warning":
		statusText = "Warning"
	case "error":
		statusText = "Error"
	}

	line := fmt.Sprintf("%s %s (%s) - %s", icon, file, conversion, statusText)
	return style.Render(line)
}

// estimateRemainingTime calculates remaining time estimate
func estimateRemainingTime(start time.Time, completed, total int) time.Duration {
	if completed == 0 {
		return 0
	}

	elapsed := time.Since(start)
	avgPerFile := elapsed / time.Duration(completed)
	remaining := avgPerFile * time.Duration(total-completed)

	return remaining
}

// formatDuration formats a duration as MM:SS or HH:MM:SS
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// renderConversionScreen renders the batch conversion progress screen
func renderConversionScreen(m model) string {
	var b strings.Builder

	// Title
	title := fmt.Sprintf("Converting %d files to %s...", m.totalFiles, strings.ToUpper(m.targetFormat))
	b.WriteString(styleHeader.Render(title))
	b.WriteString("\n")
	b.WriteString(styleBorder.Render(strings.Repeat("─", min(m.termWidth, 80))))
	b.WriteString("\n\n")

	// Progress bar
	progressBar := renderProgressBar(m.currentFile, m.totalFiles, m.errorCount)
	b.WriteString(progressBar)
	b.WriteString("\n\n")

	// Current file status
	if m.currentFile < m.totalFiles && m.currentFileName != "" {
		fileNum := m.currentFile + 1
		status := formatFileStatus(m.currentFileName, "", m.targetFormat, m.currentStatus)
		b.WriteString(fmt.Sprintf("File %d of %d:\n", fileNum, m.totalFiles))
		b.WriteString(status)
		b.WriteString("\n\n")
	}

	// Time estimation
	elapsed := formatDuration(m.elapsedTime)
	remaining := formatDuration(m.estimatedRemaining)
	timeInfo := fmt.Sprintf("Elapsed: %s | Remaining: ~%s", elapsed, remaining)
	b.WriteString(styleGray.Render(timeInfo))
	b.WriteString("\n\n")

	// Completed files (last 5)
	if len(m.completedList) > 0 {
		b.WriteString("Completed:\n")
		for _, file := range m.completedList {
			b.WriteString(fmt.Sprintf("  %s %s\n", styleGreen.Render("✓"), file))
		}
		b.WriteString("\n")
	}

	// Cancel hint
	if !m.cancelling {
		b.WriteString(styleGray.Render("Press Ctrl+C to cancel"))
	} else {
		b.WriteString(styleYellow.Render("Cancelling..."))
	}

	return b.String()
}

// renderSummaryScreen renders the final results summary
func renderSummaryScreen(results []ConversionResult, elapsed time.Duration) string {
	var b strings.Builder

	// Title
	b.WriteString(styleHeader.Render("Batch Conversion Complete!"))
	b.WriteString("\n")
	b.WriteString(styleBorder.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	// Count results
	successCount := 0
	warningCount := 0
	errorCount := 0
	cancelledCount := 0

	for _, r := range results {
		switch r.Status {
		case "success":
			successCount++
		case "warning":
			warningCount++
		case "error":
			errorCount++
		case "cancelled":
			cancelledCount++
		}
	}

	// Display counts
	b.WriteString("Results:\n")
	if successCount > 0 {
		b.WriteString(fmt.Sprintf("  %s %d succeeded\n", styleGreen.Render("✓"), successCount))
	}
	if warningCount > 0 {
		b.WriteString(fmt.Sprintf("  %s %d warnings\n", styleYellow.Render("⚠️"), warningCount))
	}
	if errorCount > 0 {
		b.WriteString(fmt.Sprintf("  %s %d errors\n", styleRed.Render("✗"), errorCount))
	}
	if cancelledCount > 0 {
		b.WriteString(fmt.Sprintf("  %s %d cancelled\n", styleGray.Render("⊘"), cancelledCount))
	}
	b.WriteString("\n")

	// List warnings
	warnings := make([]ConversionResult, 0)
	for _, r := range results {
		if r.Status == "warning" {
			warnings = append(warnings, r)
		}
	}
	if len(warnings) > 0 {
		b.WriteString("Warnings:\n")
		for _, w := range warnings {
			b.WriteString(fmt.Sprintf("  %s %s\n", styleYellow.Render("⚠️"), w.File))
			if w.Message != "" {
				b.WriteString(fmt.Sprintf("     → %s\n", w.Message))
			}
		}
		b.WriteString("\n")
	}

	// List errors
	errors := make([]ConversionResult, 0)
	for _, r := range results {
		if r.Status == "error" {
			errors = append(errors, r)
		}
	}
	if len(errors) > 0 {
		b.WriteString("Errors:\n")
		for _, e := range errors {
			b.WriteString(fmt.Sprintf("  %s %s\n", styleRed.Render("✗"), e.File))
			if e.Message != "" {
				b.WriteString(fmt.Sprintf("     → %s\n", e.Message))
			}
		}
		b.WriteString("\n")
	}

	// Total time
	b.WriteString(fmt.Sprintf("Total time: %s\n", formatDuration(elapsed)))
	b.WriteString("\n")

	// Footer
	b.WriteString(styleGray.Render("Press any key to continue"))

	return b.String()
}
