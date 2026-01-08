package batch_test

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/justin/recipe/internal/batch"
)

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
