package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// ConversionResult represents the result of a single file conversion
type ConversionResult struct {
	Input         string   `json:"input"`
	Output        string   `json:"output"`
	SourceFormat  string   `json:"source_format"`
	TargetFormat  string   `json:"target_format"`
	Success       bool     `json:"success"`
	DurationMs    int64    `json:"duration_ms"`
	FileSizeBytes int64    `json:"file_size_bytes,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
	Error         string   `json:"error,omitempty"`
}

// BatchResult aggregates results from all file conversions
type BatchResult struct {
	Batch        bool               `json:"batch"`
	Total        int                `json:"total"`
	SuccessCount int                `json:"success_count"`
	ErrorCount   int                `json:"error_count"`
	DurationMs   int64              `json:"duration_ms"`
	Results      []ConversionResult `json:"results"`
}

// isJSONMode checks if the --json flag is enabled
func isJSONMode(cmd *cobra.Command) bool {
	json, _ := cmd.Flags().GetBool("json")
	return json
}

// outputConversionResult outputs a single conversion result
func outputConversionResult(result ConversionResult, jsonMode bool) {
	if jsonMode {
		// JSON mode: output to stdout only
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			// This should never happen with our simple structs, but handle it gracefully
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		// Human-readable mode
		if result.Success {
			// Success messages go to stdout
			fileSize := formatBytes(int(result.FileSizeBytes))
			duration := formatMilliseconds(result.DurationMs)
			fmt.Fprintf(os.Stdout, "✓ Converted %s → %s (%s, %s)\n",
				result.Input, result.Output, fileSize, duration)

			// Warnings go to stderr
			if len(result.Warnings) > 0 {
				fmt.Fprintf(os.Stderr, "⚠ %d warning(s):\n", len(result.Warnings))
				for _, warning := range result.Warnings {
					fmt.Fprintf(os.Stderr, "  - %s\n", warning)
				}
			}
		} else {
			// Error messages go to stderr
			fmt.Fprintf(os.Stderr, "✗ Error converting %s: %s\n", result.Input, result.Error)
		}
	}
}

// outputBatchResult outputs batch conversion results
func outputBatchResult(result BatchResult, jsonMode bool) {
	if jsonMode {
		// JSON mode: output to stdout only
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			// This should never happen with our simple structs, but handle it gracefully
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		// Clear progress line if present
		if result.Total > 0 {
			fmt.Fprint(os.Stderr, "\r\033[K") // Clear line
		}

		// Summary message to stdout
		fmt.Fprintf(os.Stdout, "✓ Converted %d files: %d success, %d errors (%s total)\n",
			result.Total, result.SuccessCount, result.ErrorCount,
			formatMilliseconds(result.DurationMs))

		// Error details to stderr
		if result.ErrorCount > 0 {
			fmt.Fprintln(os.Stderr, "\nErrors:")
			for _, r := range result.Results {
				if !r.Success {
					fmt.Fprintf(os.Stderr, "  - %s: %s\n", r.Input, r.Error)
				}
			}
		}
	}
}

// formatMilliseconds formats milliseconds to human-readable duration
func formatMilliseconds(ms int64) string {
	d := time.Duration(ms) * time.Millisecond
	if d < time.Second {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
