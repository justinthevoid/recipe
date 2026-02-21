package batch

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestManifestSerialization(t *testing.T) {
	startTime := time.Date(2026, 1, 8, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 1, 8, 10, 5, 0, 0, time.UTC)

	manifest := BatchManifest{
		Version: ManifestVersion,
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
				InputPath:   "test/input/a.nef",
				Status:      StatusSuccess,
				OutputPath:  "test/output/a.nksc",
				Size:        1024,
				ModTime:     startTime,
				PayloadHash: "nef_hash_1",
				NP3Hash:     "np3_hash_v1",
			},
			{
				InputPath:    "test/input/b.nef",
				Status:       StatusError,
				ErrorMessage: "dummy error",
			},
		},
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	var unmarshaledManifest BatchManifest
	if err := json.Unmarshal(data, &unmarshaledManifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	if unmarshaledManifest.Version != ManifestVersion {
		t.Errorf("Expected version %s, got %s", ManifestVersion, unmarshaledManifest.Version)
	}
	if unmarshaledManifest.Summary.TotalProcessed != 2 {
		t.Errorf("Expected TotalProcessed 2, got %d", unmarshaledManifest.Summary.TotalProcessed)
	}
	if unmarshaledManifest.Summary.SuccessCount != 1 {
		t.Errorf("Expected SuccessCount 1, got %d", unmarshaledManifest.Summary.SuccessCount)
	}
	if len(unmarshaledManifest.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(unmarshaledManifest.Files))
	}

	file0 := unmarshaledManifest.Files[0]
	if file0.InputPath != "test/input/a.nef" {
		t.Errorf("Expected input path 'test/input/a.nef', got %s", file0.InputPath)
	}
	if file0.Status != StatusSuccess {
		t.Errorf("Expected status %s, got %s", StatusSuccess, file0.Status)
	}

	file1 := unmarshaledManifest.Files[1]
	if file1.Status != StatusError {
		t.Errorf("Expected status %s, got %s", StatusError, file1.Status)
	}
}

func TestWriteManifest_Atomic(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := tmpDir + "/manifest.json"

	manifest := &BatchManifest{
		Version: ManifestVersion,
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

	var readManifest BatchManifest
	if err := json.Unmarshal(data, &readManifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest from file: %v", err)
	}
	if readManifest.Version != ManifestVersion {
		t.Errorf("Expected version %s, got %s", ManifestVersion, readManifest.Version)
	}
}

func TestReadManifest(t *testing.T) {
	manifestJSON := fmt.Sprintf(`{
		"version": "%s",
		"summary": {
			"total_processed": 1,
			"np3_hash": "abc"
		},
		"files": [
			{
				"input_path": "a.nef",
				"status": "success",
				"file_size": 123,
				"mod_time": "2026-01-08T10:00:00Z"
			}
		]
	}`, ManifestVersion)

	tmpPath := t.TempDir() + "/read_manifest.json"
	if err := os.WriteFile(tmpPath, []byte(manifestJSON), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := ReadManifest(tmpPath)
	if err != nil {
		t.Fatalf("ReadManifest failed: %v", err)
	}

	if m.Version != ManifestVersion {
		t.Errorf("Expected version %s, got %s", ManifestVersion, m.Version)
	}
	if len(m.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(m.Files))
	}
	if m.Files[0].InputPath != "a.nef" {
		t.Errorf("Expected input path a.nef, got %s", m.Files[0].InputPath)
	}
	if m.Files[0].Size != 123 {
		t.Errorf("Expected size 123, got %d", m.Files[0].Size)
	}
}
