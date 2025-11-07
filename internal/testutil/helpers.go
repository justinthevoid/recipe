package testutil

import (
	"io"
	"os"
	"testing"
)

// CopyFile copies a file for testing purposes
func CopyFile(t *testing.T, src, dst string) {
	t.Helper()

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		t.Fatalf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		t.Fatalf("failed to copy file: %v", err)
	}
}

// CreateTempFile creates a temporary file with given content
func CreateTempFile(t *testing.T, content []byte) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "test-*.dat")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

// ValidateRecipeRanges validates that all recipe parameters are within expected ranges
func ValidateRecipeRanges(t *testing.T, recipe interface{}) {
	t.Helper()
	// This is a placeholder - each format test file implements its own validation
	// based on the specific format's constraints
}
