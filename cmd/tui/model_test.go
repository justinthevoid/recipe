package main

import (
	"testing"
)

// mockKeyMsg is a simple mock for testing key handling
type mockKeyMsg struct {
	str string
}

func (m mockKeyMsg) String() string {
	return m.str
}

func TestInitialModel(t *testing.T) {
	m := initialModel()

	if m.currentDir == "" {
		t.Error("currentDir should not be empty")
	}

	if m.selected == nil {
		t.Error("selected map should be initialized")
	}

	if m.termWidth != 80 || m.termHeight != 24 {
		t.Errorf("default terminal size should be 80x24, got %dx%d", m.termWidth, m.termHeight)
	}

	if m.showHelp {
		t.Error("help should not be shown initially")
	}

	if m.cursor != 0 {
		t.Error("cursor should start at 0")
	}
}

func TestNavigationKeys(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
		{Name: "file3.lrtemplate", Path: "/test/file3.lrtemplate", IsDir: false, Format: "lrtemplate"},
	}

	// Test down navigation
	newM, _ := m.handleKeyPress(mockKeyMsg{"down"})
	m = newM.(model)
	if m.cursor != 1 {
		t.Errorf("cursor should be at 1 after down, got %d", m.cursor)
	}

	// Test up navigation
	newM, _ = m.handleKeyPress(mockKeyMsg{"up"})
	m = newM.(model)
	if m.cursor != 0 {
		t.Errorf("cursor should be at 0 after up, got %d", m.cursor)
	}

	// Test boundary - up from 0 should stay at 0
	newM, _ = m.handleKeyPress(mockKeyMsg{"up"})
	m = newM.(model)
	if m.cursor != 0 {
		t.Error("cursor should stay at 0 when at top")
	}

	// Test end key
	newM, _ = m.handleKeyPress(mockKeyMsg{"end"})
	m = newM.(model)
	if m.cursor != 2 {
		t.Errorf("cursor should be at last item (2), got %d", m.cursor)
	}

	// Test down from bottom should stay at bottom
	newM, _ = m.handleKeyPress(mockKeyMsg{"down"})
	m = newM.(model)
	if m.cursor != 2 {
		t.Error("cursor should stay at bottom when at last item")
	}

	// Test home key
	newM, _ = m.handleKeyPress(mockKeyMsg{"home"})
	m = newM.(model)
	if m.cursor != 0 {
		t.Errorf("cursor should be at 0 after home, got %d", m.cursor)
	}
}

func TestFileSelection(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
		{Name: "subdir", Path: "/test/subdir", IsDir: true, Format: "dir"},
	}
	m.cursor = 0

	// Test toggle selection on file
	newM, _ := m.toggleSelection()
	m = newM.(model)
	if !m.selected["/test/file1.xmp"] {
		t.Error("file1.xmp should be selected")
	}
	if len(m.selected) != 1 {
		t.Errorf("expected 1 selected file, got %d", len(m.selected))
	}

	// Test toggle again (deselect)
	newM, _ = m.toggleSelection()
	m = newM.(model)
	if m.selected["/test/file1.xmp"] {
		t.Error("file1.xmp should be deselected")
	}
	if len(m.selected) != 0 {
		t.Errorf("expected 0 selected files, got %d", len(m.selected))
	}

	// Test can't select directories
	m.cursor = 2
	newM, _ = m.toggleSelection()
	m = newM.(model)
	if len(m.selected) != 0 {
		t.Error("directories should not be selectable")
	}
}

func TestSelectAll(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
		{Name: "subdir", Path: "/test/subdir", IsDir: true, Format: "dir"},
	}

	// Select all
	newM, _ := m.selectAll()
	m = newM.(model)

	if len(m.selected) != 2 {
		t.Errorf("expected 2 selected files (excluding directory), got %d", len(m.selected))
	}

	if !m.selected["/test/file1.xmp"] || !m.selected["/test/file2.np3"] {
		t.Error("all files should be selected")
	}

	if m.selected["/test/subdir"] {
		t.Error("directories should not be selected by selectAll")
	}
}

func TestDeselectAll(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
	}
	m.selected["/test/file1.xmp"] = true
	m.selected["/test/file2.np3"] = true

	// Deselect all
	newM, _ := m.deselectAll()
	m = newM.(model)

	if len(m.selected) != 0 {
		t.Errorf("expected 0 selected files, got %d", len(m.selected))
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		isDir    bool
		want     string
	}{
		{"NP3 file", "test.np3", false, "np3"},
		{"XMP file", "test.xmp", false, "xmp"},
		{"lrtemplate file", "test.lrtemplate", false, "lrtemplate"},
		{"Directory", "testdir", true, "dir"},
		{"Unknown file", "test.txt", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectFormat(tt.filename, tt.isDir)
			if got != tt.want {
				t.Errorf("detectFormat(%q, %v) = %q, want %q", tt.filename, tt.isDir, got, tt.want)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"Bytes", 512, "512 B"},
		{"Kilobytes", 2048, "2.0 KB"},
		{"Megabytes", 2097152, "2.0 MB"},
		{"Large file", 5242880, "5.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSize(tt.bytes)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatBadge(t *testing.T) {
	// Test that formatBadge returns non-empty strings for valid formats
	badges := []string{
		formatBadge("np3"),
		formatBadge("xmp"),
		formatBadge("lrtemplate"),
		formatBadge("dir"),
	}

	for i, badge := range badges {
		if badge == "" {
			t.Errorf("badge %d should not be empty", i)
		}
	}

	// Test unknown format returns empty badge
	unknownBadge := formatBadge("unknown")
	if unknownBadge != "   " {
		t.Errorf("unknown format should return empty badge, got %q", unknownBadge)
	}
}

func TestHelpToggle(t *testing.T) {
	m := initialModel()

	// Show help
	newM, _ := m.handleKeyPress(mockKeyMsg{"?"})
	m = newM.(model)
	if !m.showHelp {
		t.Error("help should be visible after pressing '?'")
	}

	// Hide help with '?'
	newM, _ = m.handleKeyPress(mockKeyMsg{"?"})
	m = newM.(model)
	if m.showHelp {
		t.Error("help should be hidden after pressing '?' again")
	}

	// Show help again
	newM, _ = m.handleKeyPress(mockKeyMsg{"?"})
	m = newM.(model)
	if !m.showHelp {
		t.Error("help should be visible")
	}

	// Hide help with Escape
	newM, _ = m.handleKeyPress(mockKeyMsg{"esc"})
	m = newM.(model)
	if m.showHelp {
		t.Error("help should be hidden after pressing Escape")
	}
}

// mockWindowSizeMsg is a simple mock for window resize
type mockWindowSizeMsg struct {
	Width  int
	Height int
}

func (m mockWindowSizeMsg) GetWidth() int {
	return m.Width
}

func (m mockWindowSizeMsg) GetHeight() int {
	return m.Height
}

func TestWindowResize(t *testing.T) {
	m := initialModel()

	// Simulate window resize
	newM, _ := m.Update(mockWindowSizeMsg{Width: 120, Height: 40})
	m = newM.(model)

	if m.termWidth != 120 || m.termHeight != 40 {
		t.Errorf("terminal size should be 120x40, got %dx%d", m.termWidth, m.termHeight)
	}
}

func TestFilesLoadedMsg(t *testing.T) {
	m := initialModel()

	// Simulate files loaded
	files := []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
		{Name: "file2.np3", Path: "/test/file2.np3", IsDir: false, Format: "np3"},
	}

	newM, _ := m.Update(filesLoadedMsg{files: files})
	m = newM.(model)

	if len(m.files) != 2 {
		t.Errorf("expected 2 files, got %d", len(m.files))
	}

	if m.files[0].Name != "file1.xmp" {
		t.Errorf("first file should be file1.xmp, got %s", m.files[0].Name)
	}
}

func TestNavigateIntoFile(t *testing.T) {
	m := initialModel()
	m.files = []FileInfo{
		{Name: "file1.xmp", Path: "/test/file1.xmp", IsDir: false, Format: "xmp"},
	}
	m.cursor = 0

	// Try to navigate into a file (should do nothing)
	newM, _ := m.navigateInto()
	m = newM.(model)

	// Directory should not change when trying to enter a file
	if m.currentDir != initialModel().currentDir {
		t.Error("currentDir should not change when trying to navigate into a file")
	}
}

func TestMinMaxHelpers(t *testing.T) {
	if min(5, 10) != 5 {
		t.Error("min(5, 10) should be 5")
	}

	if min(10, 5) != 5 {
		t.Error("min(10, 5) should be 5")
	}

	if max(5, 10) != 10 {
		t.Error("max(5, 10) should be 10")
	}

	if max(10, 5) != 10 {
		t.Error("max(10, 5) should be 10")
	}
}
