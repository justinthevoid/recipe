package main

import (
	"strings"
	"testing"
)

func TestRenderFileList(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", Size: 1234, IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", Size: 2048, IsDir: false, Format: "np3"},
	}
	m.cursor = 0

	output := renderFileList(m)

	if !strings.Contains(output, "file1.xmp") {
		t.Error("output should contain file1.xmp")
	}

	if !strings.Contains(output, "file2.np3") {
		t.Error("output should contain file2.np3")
	}

	if !strings.Contains(output, "Press '?' for help") {
		t.Error("output should contain help hint")
	}
}

func TestRenderEmptyDirectory(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{} // Empty directory

	output := renderFileList(m)

	if !strings.Contains(output, "No preset files found") {
		t.Error("output should show empty directory message")
	}
}

func TestRenderFileListWithSelection(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", Size: 1234, IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", Size: 2048, IsDir: false, Format: "np3"},
	}
	m.selected["/test/file1.xmp"] = true
	m.selected["/test/file2.np3"] = true

	output := renderFileList(m)

	if !strings.Contains(output, "[2 files selected]") {
		t.Error("output should show selection count")
	}
}

func TestRenderFileLine(t *testing.T) {
	file := FileInfo{
		Name:   "test.xmp",
		Path:   "/test/test.xmp",
		Size:   1234,
		IsDir:  false,
		Format: "xmp",
	}

	// Test without cursor
	line := renderFileLine(file, false, false)
	if !strings.Contains(line, "test.xmp") {
		t.Error("line should contain filename")
	}
	if !strings.Contains(line, "1.2 KB") {
		t.Error("line should contain formatted size")
	}

	// Test with cursor
	lineCursor := renderFileLine(file, true, false)
	if lineCursor == line {
		t.Error("cursor line should be different from regular line (styled)")
	}

	// Test with selection
	lineSelected := renderFileLine(file, false, true)
	if !strings.Contains(lineSelected, "✓") {
		t.Error("selected line should contain checkmark")
	}
}

func TestRenderHelp(t *testing.T) {
	m := initialModel()
	m.showHelp = true

	output := renderHelp(m)

	requiredStrings := []string{
		"Keyboard Shortcuts",
		"Navigation:",
		"Selection:",
		"Actions:",
		"Press ? or Esc to close",
	}

	for _, str := range requiredStrings {
		if !strings.Contains(output, str) {
			t.Errorf("help should contain '%s'", str)
		}
	}
}

func TestViewMethodReturnsView(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", Size: 1234, IsDir: false, Format: "xmp"},
	}

	view := m.View()

	// Just verify that View returns a non-nil view
	// We can't easily inspect the content in v2 API, but we can verify it doesn't panic
	if view.Content == nil {
		t.Error("View should return non-nil Content")
	}
}
