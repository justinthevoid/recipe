package batch_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/justin/recipe/internal/batch"
)

func TestCopyFile_Metadata(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Create source file with specific content and permissions
	content := []byte("hello world")
	err := os.WriteFile(srcPath, content, 0600) // rw-------
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Set specific modtime
	mtime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(srcPath, mtime, mtime); err != nil {
		t.Fatalf("failed to set mtime: %v", err)
	}

	// Execute CopyFile
	err = batch.CopyFile(context.Background(), srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify content
	gotContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read dest: %v", err)
	}
	if string(gotContent) != string(content) {
		t.Errorf("content mismatch")
	}

	// Verify metadata
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Fatalf("failed to stat dest: %v", err)
	}

	// Check Mode (permissions)
	// Note: on some systems permissions might be masked, but in temp dir it should work for owner bits.
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}

	// Check ModTime
	if !info.ModTime().Equal(mtime) {
		t.Errorf("expected mtime %v, got %v", mtime, info.ModTime())
	}
}

func TestCalculateFileHash(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "hash.txt")
	content := []byte("hello")
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := batch.CalculateFileHash(tmpFile)
	if err != nil {
		t.Fatalf("CalculateFileHash failed: %v", err)
	}

	// Calculate manually for verification
	hasher := sha256.New()
	hasher.Write(content)
	expected := hex.EncodeToString(hasher.Sum(nil))

	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}
