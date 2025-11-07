package main

import (
	"strings"
	"testing"
)

func TestRenderFileListScrolling(t *testing.T) {
	m := initialModel()
	m.termHeight = 10 // Small terminal

	// Create many files to trigger scrolling
	for i := 0; i < 50; i++ {
		m.files = append(m.files, FileInfo{
			Name:   "file" + string(rune(i)) + ".xmp",
			Path:   "/test/file.xmp",
			Size:   1234,
			IsDir:  false,
			Format: "xmp",
		})
	}

	// Set cursor to middle
	m.cursor = 25

	output := renderFileList(m)

	// Should show scroll indicator
	if !strings.Contains(output, "Showing") {
		t.Error("should show scroll indicator for large file list")
	}
}

func TestRenderFileListWithDirectories(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "dir1", Path: "/test/dir1", IsDir: true, Format: "dir"},
		{Name: "file1.xmp", Path: "/test/file1.xmp", Size: 1234, IsDir: false, Format: "xmp"},
	}
	m.cursor = 0

	output := renderFileList(m)

	// Directories should be shown
	if !strings.Contains(output, "dir1") {
		t.Error("output should contain directory")
	}

	// Directory badge
	if !strings.Contains(output, "📁") {
		t.Error("output should contain directory icon")
	}
}

func TestFileLineWithDirectory(t *testing.T) {
	dir := FileInfo{
		Name:   "testdir",
		Path:   "/test/testdir",
		IsDir:  true,
		Format: "dir",
	}

	line := renderFileLine(dir, false, false)

	// Directory should not show size
	if !strings.Contains(line, "testdir") {
		t.Error("line should contain directory name")
	}
}

func TestCursorBoundsChecking(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
	}

	// Simulate filesLoadedMsg with cursor out of bounds
	m.cursor = 10 // Way out of bounds

	newM, _ := m.Update(filesLoadedMsg{files: m.files})
	m = newM.(model)

	// Cursor should be reset
	if m.cursor >= len(m.files) {
		t.Error("cursor should be reset when out of bounds")
	}
}

func TestInitReturnsCommand(t *testing.T) {
	m := initialModel()

	cmd := m.Init()

	if cmd == nil {
		t.Error("Init should return a command")
	}
}

func TestMultipleKeysInSequence(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
		{Name: "file3.lrtemplate", Path: "/test/file3.lrtemplate", IsDir: false, Format: "lrtemplate"},
	}

	// Down, space (select), down, space (select), home
	newM, _ := m.handleKeyPress(mockKeyMsg{"down"})
	m = newM.(model)
	if m.cursor != 1 {
		t.Errorf("cursor should be at 1, got %d", m.cursor)
	}

	newM, _ = m.handleKeyPress(mockKeyMsg{" "})
	m = newM.(model)
	if len(m.selected) != 1 {
		t.Error("should have 1 selected file")
	}

	newM, _ = m.handleKeyPress(mockKeyMsg{"down"})
	m = newM.(model)
	if m.cursor != 2 {
		t.Errorf("cursor should be at 2, got %d", m.cursor)
	}

	newM, _ = m.handleKeyPress(mockKeyMsg{"space"})
	m = newM.(model)
	if len(m.selected) != 2 {
		t.Error("should have 2 selected files")
	}

	newM, _ = m.handleKeyPress(mockKeyMsg{"home"})
	m = newM.(model)
	if m.cursor != 0 {
		t.Error("cursor should be at 0 after home")
	}

	// Selection should persist
	if len(m.selected) != 2 {
		t.Error("selection should persist across navigation")
	}
}

func TestRenderFileLineWithLongName(t *testing.T) {
	file := FileInfo{
		Name:   "this-is-a-very-long-filename-that-might-cause-layout-issues-in-the-tui.xmp",
		Path:   "/test/long.xmp",
		Size:   999999,
		IsDir:  false,
		Format: "xmp",
	}

	// Should not panic with long filename
	line := renderFileLine(file, false, false)

	if !strings.Contains(line, "this-is-a-very-long-filename") {
		t.Error("line should contain filename (even if truncated)")
	}
}
