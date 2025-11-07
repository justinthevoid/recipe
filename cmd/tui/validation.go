package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// triggerValidation initiates the validation screen after format selection
func (m model) triggerValidation() (model, tea.Cmd) {
	// Prepare validation screen
	m.showValidation = true
	m.validationPassed = false
	m.fileListScrollOffset = 0

	// Run validation checks
	m = m.performValidation()

	return m, nil
}

// performValidation runs all validation checks on selected files
func (m model) performValidation() model {
	// Build ValidationFile list from selected files
	m.validationFiles = make([]ValidationFile, 0, len(m.selected))

	for _, file := range m.files {
		if m.selected[file.Path] {
			vf := ValidationFile{
				Name:         file.Name,
				Path:         file.Path,
				Size:         file.Size,
				SourceFormat: file.Format,
				TargetFormat: m.targetFormat,
				HasWarnings:  false,
				Warnings:     []string{},
			}
			m.validationFiles = append(m.validationFiles, vf)
		}
	}

	// Detect warnings for each file
	m.validationWarnings = detectBatchWarnings(m.validationFiles, m.targetFormat)

	// Update ValidationFile.HasWarnings and Warnings fields
	for i := range m.validationFiles {
		for _, warning := range m.validationWarnings {
			if warning.File == m.validationFiles[i].Name {
				m.validationFiles[i].HasWarnings = true
				m.validationFiles[i].Warnings = warning.Parameters
				break
			}
		}
	}

	// Calculate conversion plan
	m.validationPlan = calculateConversionPlan(m.validationFiles, m.targetFormat)

	// Validate output directory
	dirErr := validateOutputDirectory(m.outputDir)
	if dirErr != nil {
		// Directory validation failed
		m.validationPassed = false
		m.showDirectoryPrompt = true
		if strings.Contains(dirErr.Error(), "does not exist") {
			m.directoryIssue = "missing"
		} else {
			m.directoryIssue = "permission"
		}
	} else {
		m.showDirectoryPrompt = false
		m.directoryIssue = ""
	}

	// Detect overwrites
	m.overwriteFiles = detectOverwrites(m.validationFiles, m.outputDir, m.targetFormat)

	// Validation passes if directory is OK (warnings are non-blocking)
	if m.directoryIssue == "" {
		m.validationPassed = true
	}

	return m
}

// detectBatchWarnings scans all files for unmappable parameters
func detectBatchWarnings(files []ValidationFile, targetFormat string) []Warning {
	warnings := make([]Warning, 0)

	for _, file := range files {
		// Parse file to extract parameters
		recipe, err := extractParameters(file.Path)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Detect unmappable parameters
		unmappable := detectUnmappableParams(recipe, file.SourceFormat, targetFormat)

		if len(unmappable) > 0 {
			// Classify severity
			severity := "minor"
			if len(unmappable) >= 3 {
				severity = "significant"
			}

			description := fmt.Sprintf("%d unmappable parameter(s)", len(unmappable))
			if severity == "significant" {
				description += " (significant data loss possible)"
			}

			warnings = append(warnings, Warning{
				File:           file.Name,
				ParameterCount: len(unmappable),
				Parameters:     unmappable,
				Severity:       severity,
				Description:    description,
			})
		}
	}

	return warnings
}

// detectUnmappableParams identifies parameters that can't be mapped to target format
func detectUnmappableParams(recipe interface{}, sourceFormat, targetFormat string) []string {
	// For now, return empty list - will be implemented based on actual parameter mapping
	// This is a placeholder that will be enhanced with real mapping rules
	unmappable := []string{}

	// Example logic (simplified):
	// - NP3 → XMP: Some NP3-specific adjustments don't map
	// - lrtemplate → NP3: Some Lightroom features don't map

	// TODO: Implement real detection based on UniversalRecipe fields and target format

	return unmappable
}

// validateOutputDirectory checks if output directory exists and is writable
func validateOutputDirectory(path string) error {
	// Check if directory exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("cannot access directory: %w", err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	// Check write permissions by creating a temp file
	testFile := filepath.Join(path, ".recipe-test-"+time.Now().Format("20060102150405"))
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("no write permission: %s", path)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

// detectOverwrites checks for existing files that will be overwritten
func detectOverwrites(files []ValidationFile, outputDir, targetFormat string) []OverwriteInfo {
	overwrites := make([]OverwriteInfo, 0)

	for _, file := range files {
		// Construct output filename
		outputName := changeExtension(file.Name, targetFormat)
		outputPath := filepath.Join(outputDir, outputName)

		// Check if file exists
		if info, err := os.Stat(outputPath); err == nil {
			overwrites = append(overwrites, OverwriteInfo{
				File:         file.Name,
				ExistingSize: info.Size(),
				NewSize:      estimateOutputSize(file, targetFormat),
			})
		}
	}

	return overwrites
}

// changeExtension changes file extension to match target format
func changeExtension(filename, targetFormat string) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	switch targetFormat {
	case "np3":
		return nameWithoutExt + ".np3"
	case "xmp":
		return nameWithoutExt + ".xmp"
	case "lrtemplate":
		return nameWithoutExt + ".lrtemplate"
	default:
		return filename
	}
}

// estimateOutputSize estimates output file size based on format
func estimateOutputSize(file ValidationFile, targetFormat string) int64 {
	// Rough estimates based on format characteristics
	switch targetFormat {
	case "np3":
		// NP3 is binary, usually compact
		return file.Size
	case "xmp":
		// XMP is XML, usually slightly larger due to metadata
		return int64(float64(file.Size) * 1.1)
	case "lrtemplate":
		// lrtemplate is Lua, similar size to XMP
		return int64(float64(file.Size) * 1.05)
	default:
		return file.Size
	}
}

// calculateConversionPlan computes batch statistics and estimates
func calculateConversionPlan(files []ValidationFile, targetFormat string) ConversionPlan {
	plan := ConversionPlan{
		FileCount:      len(files),
		CrossFormatCount: 0,
		SameFormatCount:  0,
	}

	// Calculate sizes and format counts
	for _, file := range files {
		plan.TotalInputSize += file.Size

		if file.SourceFormat == targetFormat {
			plan.SameFormatCount++
		} else {
			plan.CrossFormatCount++
		}
	}

	// Estimate output size (average 1.05x input for format overhead)
	plan.EstimatedOutputSize = int64(float64(plan.TotalInputSize) * 1.05)

	// Estimate time (avg 100ms per file)
	plan.EstimatedTime = time.Duration(len(files)) * 100 * time.Millisecond

	// Get available disk space (cross-platform)
	plan.AvailableDiskSpace = getAvailableDiskSpace()

	return plan
}

// renderValidationScreen renders the validation screen layout
func renderValidationScreen(m model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("╔════════════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║ Conversion Validation                                                  ║\n")
	b.WriteString("╠════════════════════════════════════════════════════════════════════════╣\n")
	b.WriteString("║                                                                        ║\n")

	// Batch Details
	b.WriteString("║ Batch Details:                                                         ║\n")
	b.WriteString(fmt.Sprintf("║   Files to convert:  %-50d║\n", m.validationPlan.FileCount))

	// Source formats
	formatCounts := make(map[string]int)
	for _, f := range m.validationFiles {
		formatCounts[f.SourceFormat]++
	}
	formatStr := ""
	for format, count := range formatCounts {
		if formatStr != "" {
			formatStr += ", "
		}
		formatStr += fmt.Sprintf("%s (%d)", strings.ToUpper(format), count)
	}
	b.WriteString(fmt.Sprintf("║   Source formats:    %-50s║\n", formatStr))

	// Target format
	targetName := strings.ToUpper(m.targetFormat)
	if m.targetFormat == "np3" {
		targetName = "NP3 (Nikon NX Studio)"
	} else if m.targetFormat == "xmp" {
		targetName = "XMP (Adobe Lightroom)"
	} else if m.targetFormat == "lrtemplate" {
		targetName = "lrtemplate (Lightroom Template)"
	}
	b.WriteString(fmt.Sprintf("║   Target format:     %-50s║\n", targetName))
	b.WriteString(fmt.Sprintf("║   Output directory:  %-50s║\n", truncateString(m.outputDir, 50)))

	b.WriteString("║                                                                        ║\n")

	// File List (first 10, scrollable)
	b.WriteString("║ Files:                                                                 ║\n")
	visibleCount := min(10, len(m.validationFiles))
	for i := m.fileListScrollOffset; i < m.fileListScrollOffset+visibleCount && i < len(m.validationFiles); i++ {
		file := m.validationFiles[i]
		sizeStr := formatFileSize(file.Size)
		arrow := fmt.Sprintf("%s → %s", strings.ToUpper(file.SourceFormat), strings.ToUpper(file.TargetFormat))
		warning := ""
		if file.HasWarnings {
			warning = " ⚠️"
		}
		line := fmt.Sprintf("  %2d. %-20s %7s  %-15s%s", i+1, truncateString(file.Name, 20), sizeStr, arrow, warning)
		b.WriteString(fmt.Sprintf("║ %-70s ║\n", line))
	}

	if len(m.validationFiles) > 10 {
		remaining := len(m.validationFiles) - visibleCount
		b.WriteString(fmt.Sprintf("║   (... %d more files, ↓ to scroll)                                      ║\n", remaining))
	}

	b.WriteString("║                                                                        ║\n")

	// Warnings
	if len(m.validationWarnings) > 0 {
		b.WriteString("║ Warnings:                                                              ║\n")
		for _, w := range m.validationWarnings {
			icon := "⚠️"
			if w.Severity == "significant" {
				icon = "🔶"
			}
			b.WriteString(fmt.Sprintf("║   %s  %-60s ║\n", icon, truncateString(w.File+": "+w.Description, 60)))
			// Show first 2 parameters
			for j, param := range w.Parameters {
				if j >= 2 {
					break
				}
				b.WriteString(fmt.Sprintf("║      → %-63s ║\n", truncateString(param, 63)))
			}
		}
	} else {
		b.WriteString("║ ✓ No warnings - All conversions are lossless                          ║\n")
	}

	b.WriteString("║                                                                        ║\n")

	// Conversion Plan Summary
	b.WriteString("║ Conversion Plan:                                                       ║\n")
	b.WriteString(fmt.Sprintf("║   Total input size:    %-47s║\n", formatFileSize(m.validationPlan.TotalInputSize)))
	b.WriteString(fmt.Sprintf("║   Estimated output:    %-47s║\n", formatFileSize(m.validationPlan.EstimatedOutputSize)))
	b.WriteString(fmt.Sprintf("║   Estimated time:      %-47s║\n", m.validationPlan.EstimatedTime.String()))
	diskSpaceStr := formatFileSize(m.validationPlan.AvailableDiskSpace)
	if m.validationPlan.AvailableDiskSpace < m.validationPlan.EstimatedOutputSize {
		diskSpaceStr += " ⚠️  INSUFFICIENT"
	}
	b.WriteString(fmt.Sprintf("║   Available space:     %-47s║\n", diskSpaceStr))

	b.WriteString("║                                                                        ║\n")

	// Directory Issues (AC-5)
	if m.directoryIssue != "" {
		b.WriteString("║ Directory Issues:                                                      ║\n")
		if m.directoryIssue == "missing" {
			b.WriteString("║   ⚠️  Output directory does not exist                                 ║\n")
			b.WriteString(fmt.Sprintf("║      Path: %-60s║\n", truncateString(m.outputDir, 60)))
			b.WriteString("║      Press 'e' to edit settings and create/choose directory           ║\n")
		} else if m.directoryIssue == "permission" {
			b.WriteString("║   ⚠️  No write permission to output directory                         ║\n")
			b.WriteString(fmt.Sprintf("║      Path: %-60s║\n", truncateString(m.outputDir, 60)))
			b.WriteString("║      Press 'e' to choose a different directory                         ║\n")
		}
		b.WriteString("║                                                                        ║\n")
	}

	// Overwrite Warnings (AC-5)
	if len(m.overwriteFiles) > 0 {
		b.WriteString("║ Overwrite Warnings:                                                    ║\n")
		showCount := min(3, len(m.overwriteFiles))
		for i := 0; i < showCount; i++ {
			ow := m.overwriteFiles[i]
			existingSize := formatFileSize(ow.ExistingSize)
			newSize := formatFileSize(ow.NewSize)
			line := fmt.Sprintf("   🔄  %s (%s → %s)", truncateString(ow.File, 35), existingSize, newSize)
			b.WriteString(fmt.Sprintf("║ %-70s ║\n", line))
		}
		if len(m.overwriteFiles) > 3 {
			b.WriteString(fmt.Sprintf("║   (... %d more files will be overwritten)                              ║\n", len(m.overwriteFiles)-3))
		}
		b.WriteString("║                                                                        ║\n")
	}

	// Footer with shortcuts
	if m.validationPassed {
		b.WriteString("║ Press 'c' to confirm conversion                                        ║\n")
	} else {
		b.WriteString("║ ⚠️  Validation failed - Press 'e' to fix issues                       ║\n")
	}
	b.WriteString("║ Press 'e' to edit settings                                             ║\n")
	b.WriteString("║ Press Esc to cancel                                                    ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════════════════╝\n")

	return b.String()
}

// Helper functions
// Note: getAvailableDiskSpace is implemented in diskspace_unix.go and diskspace_windows.go
