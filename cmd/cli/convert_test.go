package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		input  string
		format string
		want   string
	}{
		{"portrait.xmp", "np3", "portrait.np3"},
		{"portrait.np3", "xmp", "portrait.xmp"},
		{"file.with.dots.xmp", "np3", "file.with.dots.np3"},
		{`C:\Windows\path\file.xmp`, "np3", `C:\Windows\path\file.np3`},
	}

	for _, tt := range tests {
		got := generateOutputPath(tt.input, tt.format)
		if got != tt.want {
			t.Errorf("generateOutputPath(%q, %q) = %q, want %q", tt.input, tt.format, got, tt.want)
		}
	}
}

func TestCheckOutputExists(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")

	// Create an existing file
	if err := os.WriteFile(existingFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		path      string
		overwrite bool
		wantErr   bool
	}{
		{"existing file, no overwrite", existingFile, false, true},
		{"existing file, with overwrite", existingFile, true, false},
		{"non-existent file, no overwrite", nonExistentFile, false, false},
		{"non-existent file, with overwrite", nonExistentFile, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkOutputExists(tt.path, tt.overwrite)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkOutputExists(%q, %v) error = %v, wantErr %v", tt.path, tt.overwrite, err, tt.wantErr)
			}
		})
	}
}

func TestEnsureOutputDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"existing directory", filepath.Join(tmpDir, "file.txt"), false},
		{"nested new directory", filepath.Join(tmpDir, "new", "nested", "dirs", "file.txt"), false},
		{"current directory", "file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ensureOutputDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureOutputDir(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}

			// Verify directory was created (except for current directory test)
			if tt.name != "current directory" {
				dir := filepath.Dir(tt.path)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Errorf("ensureOutputDir(%q) did not create directory %q", tt.path, dir)
				}
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes int
		want  string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1234, "1.2 KB"},
		{10240, "10.0 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{10485760, "10.0 MB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.bytes)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{1 * time.Millisecond, "1ms"},
		{15 * time.Millisecond, "15ms"},
		{999 * time.Millisecond, "999ms"},
		{1 * time.Second, "1.00s"},
		{1500 * time.Millisecond, "1.50s"},
		{2 * time.Second, "2.00s"},
		{2*time.Second + 345*time.Millisecond, "2.34s"}, // Float precision
	}

	for _, tt := range tests {
		got := formatDuration(tt.duration)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
		}
	}
}
