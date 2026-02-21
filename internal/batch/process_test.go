package batch_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/batch"
	"github.com/justin/recipe/internal/formats/np3"
)

func TestProcessBatch_Parallel(t *testing.T) {
	tmpDir := t.TempDir()
	outDir := t.TempDir()

	// Create 10 files
	count := 10
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("file%d.nef", i)
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  tmpDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{Label: "ParallelTest"},
		Workers:    4,
	})

	result, err := orch.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("ProcessBatch failed: %v", err)
	}

	if result.Processed != count {
		t.Errorf("Expected %d processed, got %d", count, result.Processed)
	}
}

func TestProcessBatch_PartialFailure(t *testing.T) {
	tmpDir := t.TempDir()
	outDir := t.TempDir()

	files := map[string]os.FileMode{
		"valid1.nef": 0644,
		"valid2.nrw": 0644,
		"bad.nef":    0000,
	}

	for name, perm := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("data"), perm); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		os.Chmod(path, perm)
	}

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  tmpDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{Label: "TestRecipe"},
	})

	result, err := orch.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("ProcessBatch failed with fatal error: %v", err)
	}

	if result.TotalFiles != 3 {
		t.Errorf("expected 3 total files, got %d", result.TotalFiles)
	}
	if result.Processed != 2 {
		t.Errorf("expected 2 processed, got %d", result.Processed)
	}
	if result.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", result.Failed)
	}

	foundBad := false
	for _, err := range result.Errors {
		if strings.Contains(err.Error(), "bad.nef") {
			foundBad = true
			break
		}
	}
	if !foundBad {
		t.Errorf("expected error for bad.nef, got errors: %v", result.Errors)
	}

	if _, err := os.Stat(filepath.Join(outDir, "valid1.nef")); err != nil {
		t.Errorf("missing valid1.nef in output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "valid1.nef.nksc")); err != nil {
		t.Errorf("missing valid1.nef.nksc sidecar: %v", err)
	}
}
