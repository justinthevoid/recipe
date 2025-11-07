package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFilesCmd(t *testing.T) {
	// Create temp directory with test files
	tmpDir, err := os.MkdirTemp("", "tui-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"test1.xmp",
		"test2.np3",
		"test3.lrtemplate",
		"test4.txt", // Should be filtered out
	}

	for _, name := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Create a subdirectory
	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// Load files
	cmd := loadFilesCmd(tmpDir)
	msg := cmd()

	// Check the message
	switch msg := msg.(type) {
	case filesLoadedMsg:
		if msg.err != nil {
			t.Errorf("loadFiles returned error: %v", msg.err)
		}

		// Should have 3 preset files + 1 directory = 4 items
		if len(msg.files) != 4 {
			t.Errorf("expected 4 files (3 presets + 1 dir), got %d", len(msg.files))
		}

		// Verify txt file was filtered out
		for _, file := range msg.files {
			if file.Name == "test4.txt" {
				t.Error("txt file should have been filtered out")
			}
		}

		// Verify directory is included
		foundDir := false
		for _, file := range msg.files {
			if file.Name == "subdir" {
				foundDir = true
				if !file.IsDir {
					t.Error("subdir should be marked as directory")
				}
				if file.Format != "dir" {
					t.Errorf("subdir format should be 'dir', got '%s'", file.Format)
				}
			}
		}
		if !foundDir {
			t.Error("directory should be included in file list")
		}

	default:
		t.Errorf("expected filesLoadedMsg, got %T", msg)
	}
}

func TestLoadFilesCmdError(t *testing.T) {
	// Try to load from non-existent directory
	cmd := loadFilesCmd("/nonexistent/path/that/doesnt/exist")
	msg := cmd()

	switch msg := msg.(type) {
	case filesLoadedMsg:
		if msg.err == nil {
			t.Error("expected error for non-existent directory")
		}
	default:
		t.Errorf("expected filesLoadedMsg, got %T", msg)
	}
}

func TestDetectFormatComprehensive(t *testing.T) {
	tests := []struct {
		filename string
		isDir    bool
		expected string
	}{
		{"test.np3", false, "np3"},
		{"test.NP3", false, ""}, // Case sensitive
		{"test.xmp", false, "xmp"},
		{"test.XMP", false, ""}, // Case sensitive
		{"test.lrtemplate", false, "lrtemplate"},
		{"test.LRTEMPLATE", false, ""}, // Case sensitive
		{"directory", true, "dir"},
		{"file.txt", false, ""},
		{"file.jpg", false, ""},
		{"no-extension", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := detectFormat(tt.filename, tt.isDir)
			if got != tt.expected {
				t.Errorf("detectFormat(%q, %v) = %q, want %q", tt.filename, tt.isDir, got, tt.expected)
			}
		})
	}
}
