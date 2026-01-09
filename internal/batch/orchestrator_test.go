package batch_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/justin/recipe/internal/apperr"
	"github.com/justin/recipe/internal/batch"
	"github.com/justin/recipe/internal/formats/np3"
)

func TestOrchestrator_ProcessBatch_ErrorWrapping(t *testing.T) {
	// Setup input with one file
	inDir := t.TempDir()
	nefPath := filepath.Join(inDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(t.TempDir(), "readonly")
	if err := os.Mkdir(outDir, 0444); err != nil { // Read-only directory
		t.Fatal(err)
	}
	outDir = filepath.Join(outDir, "nested") // Should fail to create

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  inDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{},
	})

	result, err := orch.ProcessBatch()
	if err != nil {
		t.Fatalf("Unexpected ProcessBatch error: %v", err)
	}

	if result.Failed != 1 {
		t.Errorf("Expected 1 failure, got %d", result.Failed)
	}

	if len(result.Errors) == 0 {
		t.Fatal("Expected errors in result")
	}

	var appErr *apperr.Error
	if !errors.As(result.Errors[0], &appErr) {
		t.Errorf("Error is not apperr.Error: %T %v", result.Errors[0], result.Errors[0])
	} else {
		// verify context
		if appErr.File != "test.nef" {
			t.Errorf("Expected context file 'test.nef', got '%s'", appErr.File)
		}
	}
}

func TestOrchestrator_FindFiles(t *testing.T) {
	// Setup temp directory with nested files
	tmpDir := t.TempDir()

	filesToCreate := []string{
		"test1.NEF",
		"sub/test2.nrw", // Case insensitive check
		"sub/ignore.txt",
		"TEST3.nef",
	}

	for _, relPath := range filesToCreate {
		path := filepath.Join(tmpDir, relPath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte("fake content"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	// Create Orchestrator
	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  tmpDir,
		OutputPath: t.TempDir(),
	})

	// Execute FindFiles
	files, err := orch.FindFiles()
	if err != nil {
		t.Fatalf("FindFiles returned error: %v", err)
	}

	// Verify results
	expectedFiles := []string{
		"TEST3.nef",
		"sub/test2.nrw",
		"test1.NEF",
	}

	var foundRelPaths []string
	for _, f := range files {
		rel, err := filepath.Rel(tmpDir, f)
		if err != nil {
			t.Fatalf("failed to get relative path: %v", err)
		}
		foundRelPaths = append(foundRelPaths, rel)
	}
	sort.Strings(foundRelPaths)
	sort.Strings(expectedFiles)

	if !reflect.DeepEqual(expectedFiles, foundRelPaths) {
		t.Errorf("FindFiles() mismatch\nExpected: %v\nGot: %v", expectedFiles, foundRelPaths)
	}
}
