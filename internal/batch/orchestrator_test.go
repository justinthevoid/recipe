package batch_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

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

	result, err := orch.ProcessBatch(context.Background())
	// Since the output directory is read-only/invalid, WriteManifest will also fail.
	// We accept this specific error but verify we still got the file processing results.
	if err != nil {
		t.Logf("Got expected manifest write error: %v", err)
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

func TestProcessBatch_GeneratesManifest(t *testing.T) {
	inDir := t.TempDir()
	nefPath := filepath.Join(inDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := t.TempDir()

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  inDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{RawBytes: []byte("np3data")},
	})

	_, err := orch.ProcessBatch(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Check manifest
	manifestPath := filepath.Join(outDir, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatal("manifest.json not created")
	}
}

func TestOrchestrator_IdempotencyIntegration(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// 1. Setup Input: 2 files
	file1 := filepath.Join(inDir, "f1.nef")
	file2 := filepath.Join(inDir, "f2.nef")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Dummy NP3
	recipe := &np3.Metadata{
		RawBytes: []byte("dummy_recipe"),
	}

	cfg := batch.Config{
		InputPath:  inDir,
		OutputPath: outDir,
		Recipe:     recipe,
	}

	orch := batch.NewOrchestrator(cfg)

	// 2. First Run: Process all
	res1, err := orch.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.Processed != 2 {
		t.Errorf("Run 1: expected 2 processed, got %d", res1.Processed)
	}

	// 3. Second Run: No changes -> Skip all
	// Re-create orchestrator to simulate fresh run (reloads manifest)
	orch2 := batch.NewOrchestrator(cfg)
	res2, err := orch2.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.Skipped != 2 {
		t.Errorf("Run 2: expected 2 skipped, got %d", res2.Skipped)
	}
	if res2.Processed != 0 {
		t.Errorf("Run 2: expected 0 processed, got %d", res2.Processed)
	}

	// 4. Modify one file (change content -> hash changes)
	newTime := time.Now().Add(1 * time.Hour)
	if err := os.WriteFile(file1, []byte("content1_modified"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file1, newTime, newTime); err != nil {
		t.Fatal(err)
	}

	orch3 := batch.NewOrchestrator(cfg)
	res3, err := orch3.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("Run 3 failed: %v", err)
	}
	// Run 3 checks - f1 should process (modified), f2 should skip (unchanged)
	_ = res3
	if res3.Processed != 1 {
		t.Errorf("Run 3: expected 1 processed, got %d. Errors: %v", res3.Processed, res3.Errors)
	}
	if res3.Skipped != 1 {
		t.Errorf("Run 3: expected 1 skipped, got %d", res3.Skipped)
	}

	// 5. Force run -> Process all
	cfg.Force = true
	orch4 := batch.NewOrchestrator(cfg)
	res4, err := orch4.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("Run 4 failed: %v", err)
	}
	if res4.Processed != 2 {
		t.Errorf("Run 4: expected 2 processed (force), got %d", res4.Processed)
	}
}

func TestProcessBatch_DryRun(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// Create a valid file
	nefPath := filepath.Join(inDir, "test.nef")
	if err := os.WriteFile(nefPath, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  inDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{RawBytes: []byte("np3data")},
		DryRun:     true,
	})

	// Run ProcessBatch
	result, err := orch.ProcessBatch(context.Background())
	if err != nil {
		t.Fatalf("ProcessBatch failed in dry-run: %v", err)
	}

	// 1. Check no files written
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) > 0 {
		t.Errorf("DryRun should not write files, found: %v", entries)
	}

	// 2. Check result status
	// We expect files to be Skipped
	if result.Skipped != 1 {
		t.Errorf("DryRun should count files as skipped, got Skipped=%d", result.Skipped)
	}
}

func TestProcessBatch_Strict_AbortsOnError(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	// 1. Bad file (permissions) - fail immediately
	badFile := filepath.Join(inDir, "bad.nef")
	if err := os.WriteFile(badFile, []byte("bad"), 0000); err != nil {
		t.Fatal(err)
	}
	os.Chmod(badFile, 0000)

	orch := batch.NewOrchestrator(batch.Config{
		InputPath:  inDir,
		OutputPath: outDir,
		Recipe:     &np3.Metadata{RawBytes: []byte("np3data")},
		Strict:     true,
		Workers:    1,
	})

	_, err := orch.ProcessBatch(context.Background())
	if err == nil {
		t.Fatal("Strict mode should return error on failure")
	}

	if !strings.Contains(err.Error(), "strict mode abort") && !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
