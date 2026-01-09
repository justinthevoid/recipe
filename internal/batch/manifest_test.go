package batch

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestManifestSerialization(t *testing.T) {
	startTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 1, 8, 10, 5, 0, 0, time.UTC)

	manifest := BatchManifest{
		Version: "1.0",
		Summary: BatchSummary{
			TotalProcessed: 2,
			SuccessCount:   1,
			SkippedCount:   0,
			FailureCount:   1,
			StartTime:      startTime,
			EndTime:        endTime,
			Duration:       "5m0s",
			NP3Hash:        "abcdef123",
		},
		Files: []FileResult{
			{
				InputPath:  "test/input/a.nef",
				Status:     "success",
				OutputPath: "test/output/a.nksc",
			},
			{
				InputPath:    "test/input/b.nef",
				Status:       "error",
				ErrorMessage: "dummy error",
			},
		},
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	jsonStr := string(data)

	// Check for JSON tags (snake_case expectations)
	checks := []string{
		`"version":"1.0"`,
		`"total_processed":2`,
		`"success_count":1`,
		`"skipped_count":0`,
		`"input_path":"test/input/a.nef"`,
		`"error":"dummy error"`,
		`"duration":"5m0s"`,
		`"np3_hash":"abcdef123"`,
	}

	for _, check := range checks {
		if !strings.Contains(jsonStr, check) {
			t.Errorf("JSON output missing expected string: %s. Got: %s", check, jsonStr)
		}
	}
}

func TestWriteManifest_Atomic(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := tmpDir + "/manifest.json"

	manifest := &BatchManifest{
		Version: "1.0",
		Summary: BatchSummary{
			TotalProcessed: 1,
			SuccessCount:   1,
		},
	}

	if err := WriteManifest(manifestPath, manifest); err != nil {
		t.Fatalf("WriteManifest failed: %v", err)
	}

	// Verify file exists and has content
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("Failed to read manifest file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"version": "1.0"`) {
		t.Errorf("Manifest content missing version. Got: %s", content)
	}
}
