package batch

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// BatchManifest represents the complete report of a batch processing run.
type BatchManifest struct {
	Version string       `json:"version"`
	Summary BatchSummary `json:"summary"`
	Files   []FileResult `json:"files"`
}

// BatchSummary contains aggregate statistics about the batch run.
type BatchSummary struct {
	TotalProcessed int       `json:"total_processed"`
	SuccessCount   int       `json:"success_count"`
	FailureCount   int       `json:"failure_count"`
	SkippedCount   int       `json:"skipped_count"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Duration       string    `json:"duration"` // Human readable duration
	NP3Hash        string    `json:"np3_hash,omitempty"`
}

// FileResult represents the outcome of processing a single file.
type FileResult struct {
	InputPath    string `json:"input_path"`
	Status       string `json:"status"` // "success" or "error"
	OutputPath   string `json:"output_path,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

// WriteManifest writes the manifest to the specified path using an atomic write pattern.
// It writes to a temporary file first, then renames it to the target path.
func WriteManifest(path string, manifest *BatchManifest) error {
	tmpPath := path + ".tmp"

	// Create/Overwrite temp file
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp manifest: %w", err)
	}
	// Ensure we close the file, handling potential errors if not already closed
	defer func() {
		file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}

	// Close before rename to ensure flush (defer runs after rename otherwise)
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close temp manifest: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename manifest: %w", err)
	}

	return nil
}
