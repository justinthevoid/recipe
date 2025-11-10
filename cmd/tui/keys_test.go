package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNavigateUp(t *testing.T) {
	// Create a nested directory structure for testing
	tmpDir, err := os.MkdirTemp("", "tui-nav-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	m := initialModel()
	m.currentDir = subdir

	// Navigate up
	newM, _ := m.navigateUp()
	m = newM.(model)

	// Should navigate to parent
	if m.currentDir == subdir {
		t.Error("should navigate to parent directory")
	}

	// Cursor should reset
	if m.cursor != 0 {
		t.Error("cursor should reset to 0 after navigation")
	}
}

func TestNavigateUpAtRoot(t *testing.T) {
	m := initialModel()

	// Get a root path (platform specific)
	root := filepath.VolumeName(m.currentDir) + "/"
	if root == "/" && os.PathSeparator == '\\' {
		root = "C:\\"
	} else if root == "C:\\" && os.PathSeparator == '/' {
		root = "/"
	}

	// Try to navigate up from close to root
	parent := filepath.Dir(m.currentDir)
	m.currentDir = parent

	newM, _ := m.navigateUp()
	m = newM.(model)

	// This test just ensures it doesn't panic
	// and returns a valid model
	if m.currentDir == "" {
		t.Error("currentDir should not be empty after navigation")
	}
}

func TestNavigateIntoDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tui-nav-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subdir := filepath.Join(tmpDir, "testdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	m := initialModel()
	m.currentDir = tmpDir
	m.files = []FileInfo{
		{Name: "testdir", Path: subdir, IsDir: true, Format: "dir"},
	}
	m.cursor = 0

	// Navigate into directory
	newM, _ := m.navigateInto()
	m = newM.(model)

	if m.currentDir != subdir {
		t.Errorf("should navigate into directory, expected %s, got %s", subdir, m.currentDir)
	}

	if m.cursor != 0 {
		t.Error("cursor should reset after navigation")
	}
}

func TestKeyPressQuit(t *testing.T) {
	m := initialModel()

	// Test 'q' key
	_, cmd := m.handleKeyPress(mockKeyMsg{"q"})
	if cmd == nil {
		t.Error("'q' should return tea.Quit command")
	}

	// Test 'ctrl+c'
	_, cmd = m.handleKeyPress(mockKeyMsg{"ctrl+c"})
	if cmd == nil {
		t.Error("'ctrl+c' should return tea.Quit command")
	}
}

func TestKeyPressRefresh(t *testing.T) {
	m := initialModel()

	// Test 'r' key
	_, cmd := m.handleKeyPress(mockKeyMsg{"r"})
	if cmd == nil {
		t.Error("'r' should return loadFiles command")
	}
}

func TestKeyPressVimKeys(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
	}

	// Test 'j' (vim down)
	newM, _ := m.handleKeyPress(mockKeyMsg{"j"})
	m = newM.(model)
	if m.cursor != 1 {
		t.Errorf("'j' should move cursor down, got %d", m.cursor)
	}

	// Test 'k' (vim up)
	newM, _ = m.handleKeyPress(mockKeyMsg{"k"})
	m = newM.(model)
	if m.cursor != 0 {
		t.Errorf("'k' should move cursor up, got %d", m.cursor)
	}
}

func TestToggleSelectionEmptyFiles(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{} // Empty file list

	// Try to toggle selection with empty file list
	newM, _ := m.toggleSelection()
	m = newM.(model)

	// Should handle gracefully
	if len(m.selected) != 0 {
		t.Error("no files should be selected from empty list")
	}
}

func TestNavigateIntoEmptyFiles(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{} // Empty file list

	oldDir := m.currentDir

	// Try to navigate with empty file list
	newM, _ := m.navigateInto()
	m = newM.(model)

	// Should stay in same directory
	if m.currentDir != oldDir {
		t.Error("should not navigate with empty file list")
	}
}

func TestSelectAllNoFiles(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{} // Empty

	newM, _ := m.selectAll()
	m = newM.(model)

	if len(m.selected) != 0 {
		t.Error("no files should be selected from empty list")
	}
}
