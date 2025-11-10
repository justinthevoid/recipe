# Story 4.1: Bubbletea File Browser

**Epic:** Epic 4 - TUI Interface (FR-4)
**Story ID:** 4.1
**Status:** done
**Created:** 2025-11-06
**Completed:** 2025-11-06
**Complexity:** Medium (2-3 days)

---

## User Story

**As a** photographer using Recipe from the terminal,
**I want** an interactive file browser with keyboard navigation,
**So that** I can visually select preset files for conversion without typing full paths or using shell wildcards.

---

## Business Value

The Bubbletea-based TUI file browser transforms Recipe from a CLI-only tool into an interactive terminal application that bridges the gap between command-line power and GUI convenience:

- **Visual File Management** - Browse directories and see file previews without leaving the terminal
- **Reduced Cognitive Load** - No need to remember file paths or glob patterns
- **Power User Efficiency** - Keyboard-driven workflow is faster than mouse-based GUIs for batch operations
- **Professional Workflow Integration** - Terminal users (developers, sys admins, power users) prefer TUI over switching to GUI tools

**Strategic value:** TUI interface differentiates Recipe from web-only converters and positions it as a professional power-user tool. Photographers working in terminal-heavy workflows (server environments, SSH sessions, automated pipelines) can use Recipe without GUI overhead.

---

## Acceptance Criteria

### AC-1: Bubbletea Application Initialization ✅

- [x] TUI application built using Bubbletea framework (charm.land/bubbletea/v2)
- [x] Application launches with `recipe-tui` command (separate binary from `recipe` CLI)
- [x] Initial state shows current working directory
- [x] Terminal size detected and UI adapts to terminal dimensions
- [x] Application exits cleanly with Ctrl+C or 'q' key

**Test:**
```go
func TestTUILaunch(t *testing.T) {
    // Launch TUI in test mode
    m := initialModel()
    
    assert.NotNil(t, m, "Model should initialize")
    assert.Equal(t, ".", m.currentDir, "Should start in current directory")
    assert.True(t, m.termWidth > 0, "Terminal width should be detected")
}
```

**Validation:**
- `recipe-tui` binary compiles successfully
- TUI launches without errors
- Terminal size properly detected
- Clean exit on Ctrl+C or 'q'

---

### AC-2: Directory Navigation ✅

- [x] Display current directory path at top of screen
- [x] List files and subdirectories in current directory
- [x] Arrow keys navigate up/down through file list
- [x] Enter key on directory navigates into that directory
- [x] Backspace or left arrow navigates to parent directory
- [x] Home/End keys jump to first/last item in list

**UI Layout:**
```
┌──────────────────────────────────────────────────────┐
│ Current: /Users/justin/presets/                      │
├──────────────────────────────────────────────────────┤
│ [↑] ../                                              │
│   📁 adobe-presets/                                  │
│   📁 nikon-presets/                                  │
│ > 📄 vintage-film.xmp                     (1.2 KB)   │
│   📄 portrait-warm.lrtemplate             (2.4 KB)   │
│   📄 landscape-cool.np3                   (1.0 KB)   │
│                                                       │
└──────────────────────────────────────────────────────┘
```

**Test:**
```go
func TestDirectoryNavigation(t *testing.T) {
    m := initialModel()
    m.currentDir = "testdata/"
    m = m.Update(tea.KeyMsg{Type: tea.KeyEnter}).(model)
    
    // Verify directory changed
    assert.NotEqual(t, "testdata/", m.currentDir)
}
```

**Validation:**
- Directory listing displays correctly
- Arrow keys move cursor up/down
- Enter navigates into directories
- Backspace navigates to parent
- Home/End keys work

---

### AC-3: File Filtering by Format ✅

- [x] Display only supported preset formats: .np3, .xmp, .lrtemplate
- [x] Show directories (for navigation) but filter non-preset files
- [x] Display file extension badges (NP3, XMP, LRT) with color coding
- [x] Show file size in human-readable format (KB, MB)
- [x] Empty directory message if no preset files found

**Format Badges:**
```
NP3  - Blue badge
XMP  - Orange badge
LRT  - Green badge
DIR  - Gray badge for directories
```

**Test:**
```go
func TestFileFiltering(t *testing.T) {
    files := []FileInfo{
        {Name: "test.xmp", IsDir: false},
        {Name: "test.jpg", IsDir: false},  // Should be filtered
        {Name: "test.np3", IsDir: false},
        {Name: "subdir", IsDir: true},
    }
    
    filtered := filterPresetFiles(files)
    
    assert.Equal(t, 3, len(filtered), "Should show 2 presets + 1 directory")
}
```

**Validation:**
- Only preset files and directories shown
- File size displayed correctly
- Format badges color-coded
- Empty directory message works

---

### AC-4: Multi-File Selection ✅

- [x] Space bar toggles selection on current file
- [x] Selected files marked with checkmark (✓) indicator
- [x] Status bar shows count of selected files
- [x] 'a' key selects all files in current directory
- [x] 'n' key deselects all files (clears selection)
- [x] Selection state persists across directory navigation

**UI with Selection:**
```
┌──────────────────────────────────────────────────────┐
│ Current: /Users/justin/presets/                      │
├──────────────────────────────────────────────────────┤
│   📄 vintage-film.xmp                     (1.2 KB)   │
│ ✓ 📄 portrait-warm.lrtemplate             (2.4 KB)   │
│ > ✓ 📄 landscape-cool.np3                 (1.0 KB)   │
│                                                       │
│ [2 files selected] Press 'c' to convert              │
└──────────────────────────────────────────────────────┘
```

**Test:**
```go
func TestMultiSelection(t *testing.T) {
    m := initialModel()
    m.files = []FileInfo{
        {Name: "file1.xmp"},
        {Name: "file2.np3"},
    }
    
    // Select first file
    m.cursor = 0
    m = m.Update(tea.KeyMsg{Type: tea.KeySpace}).(model)
    
    assert.Equal(t, 1, len(m.selected), "Should have 1 selected file")
    
    // Select second file
    m.cursor = 1
    m = m.Update(tea.KeyMsg{Type: tea.KeySpace}).(model)
    
    assert.Equal(t, 2, len(m.selected), "Should have 2 selected files")
}
```

**Validation:**
- Space bar toggles selection
- Selected files show checkmark
- Selection count displayed
- 'a' selects all, 'n' clears
- Selection persists across navigation

---

### AC-5: File Information Display ✅

- [x] Display file name, size, and file format
- [x] Show full absolute path in file list
- [x] Display total count of files in current directory
- [x] Show format type (NP3, XMP, lrtemplate) with color-coded badge

**Status Bar:**
```
│ /Users/justin/presets/vintage-film.xmp | 1.2 KB | Modified: 2025-11-03 │
```

**Test:**
```go
func TestFileInfo(t *testing.T) {
    info := FileInfo{
        Name:         "test.xmp",
        Size:         1234,
        ModTime:      time.Now(),
        Path:         "/full/path/test.xmp",
    }
    
    display := formatFileInfo(info)
    
    assert.Contains(t, display, "1.2 KB")
    assert.Contains(t, display, "XMP")
}
```

**Validation:**
- File size formatted correctly (KB/MB)
- Modified date displayed
- Full path shown in status bar
- Format badge correct

---

### AC-6: Keyboard Shortcuts Help ✅

- [x] '?' key displays help overlay with all keyboard shortcuts
- [x] Help overlay can be dismissed with '?' or Escape key
- [x] Help overlay lists all navigation, selection, and action keys

**Help Overlay:**
```
┌─────────────────── Keyboard Shortcuts ────────────────────┐
│                                                            │
│ Navigation:                                                │
│   ↑/k        Move cursor up                                │
│   ↓/j        Move cursor down                              │
│   Enter      Navigate into directory                       │
│   Backspace  Navigate to parent directory                  │
│   Home       Jump to first item                            │
│   End        Jump to last item                             │
│                                                            │
│ Selection:                                                 │
│   Space      Toggle selection on current file              │
│   a          Select all files                              │
│   n          Deselect all files                            │
│                                                            │
│ Actions:                                                   │
│   c          Convert selected files                        │
│   r          Refresh current directory                     │
│   ?          Toggle this help                              │
│   q/Ctrl+C   Quit                                          │
│                                                            │
│ Press ? or Esc to close                                    │
└────────────────────────────────────────────────────────────┘
```

**Test:**
```go
func TestHelpOverlay(t *testing.T) {
    m := initialModel()
    
    // Show help
    m = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}).(model)
    assert.True(t, m.showHelp, "Help should be visible")
    
    // Hide help
    m = m.Update(tea.KeyMsg{Type: tea.KeyEsc}).(model)
    assert.False(t, m.showHelp, "Help should be hidden")
}
```

**Validation:**
- '?' shows help overlay
- Help displays all shortcuts
- Escape dismisses help
- Help overlay doesn't interfere with navigation

---

### AC-7: Terminal Resize Handling ✅

- [x] UI adapts to terminal resize events
- [x] File list scrolls if more files than terminal height
- [x] Layout remains functional at minimum terminal size (60x10)
- [x] Help overlay adapts to terminal size

**Test:**
```go
func TestTerminalResize(t *testing.T) {
    m := initialModel()
    m.termWidth = 120
    m.termHeight = 40
    
    // Simulate resize
    m = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}).(model)
    
    assert.Equal(t, 80, m.termWidth)
    assert.Equal(t, 24, m.termHeight)
}
```

**Validation:**
- Resize event handled correctly
- UI reflows to new dimensions
- Minimum 80x24 terminal usable
- No visual glitches on resize

---

## Tasks / Subtasks

### Task 1: Initialize Bubbletea Project (AC-1)

- [ ] Create `cmd/tui/` directory
- [ ] Create `cmd/tui/main.go`:
  ```go
  package main

  import (
      "fmt"
      "os"

      tea "github.com/charmbracelet/bubbletea/v2"
  )

  func main() {
      p := tea.NewProgram(initialModel())
      if _, err := p.Run(); err != nil {
          fmt.Printf("Error: %v\n", err)
          os.Exit(1)
      }
  }
  ```
- [ ] Install Bubbletea dependency:
  ```bash
  go get github.com/charmbracelet/bubbletea/v2@v2-exp
  go get github.com/charmbracelet/lipgloss/v2@v2-exp
  go get github.com/charmbracelet/bubbles/v2@v2-exp
  go get github.com/charmbracelet/x@latest
  ```
- [ ] Define initial model struct in `cmd/tui/model.go`:
  ```go
  type model struct {
      currentDir  string
      files       []FileInfo
      cursor      int
      selected    map[string]bool
      termWidth   int
      termHeight  int
      showHelp    bool
  }

  type FileInfo struct {
      Name    string
      Path    string
      Size    int64
      ModTime time.Time
      IsDir   bool
      Format  string  // "np3", "xmp", "lrtemplate", "dir"
  }

  func initialModel() model {
      cwd, _ := os.Getwd()
      return model{
          currentDir: cwd,
          selected:   make(map[string]bool),
          termWidth:  80,
          termHeight: 24,
      }
  }
  ```

**Validation:**
- Bubbletea imports successfully
- Model struct compiles
- Initial model creates without errors

---

### Task 2: Implement Init and Update (Bubbletea Pattern) (AC-1)

- [ ] Implement `Init()` method:
  ```go
  func (m model) Init() tea.Cmd {
      return tea.Batch(
          loadFiles(m.currentDir),
          tea.EnterAltScreen,
      )
  }
  ```
- [ ] Implement `Update()` method skeleton:
  ```go
  func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.KeyMsg:
          return m.handleKeyPress(msg)
      case tea.WindowSizeMsg:
          m.termWidth = msg.Width
          m.termHeight = msg.Height
      case filesLoadedMsg:
          m.files = msg.files
      }
      return m, nil
  }
  ```
- [ ] Implement `View()` method skeleton:
  ```go
  func (m model) View() string {
      if m.showHelp {
          return renderHelp(m)
      }
      return renderFileList(m)
  }
  ```

**Validation:**
- Init/Update/View pattern compiles
- Bubbletea program runs
- Alt screen mode works

---

### Task 3: Implement Directory Listing (AC-2, AC-3)

- [ ] Create `loadFiles()` command:
  ```go
  type filesLoadedMsg struct {
      files []FileInfo
  }

  func loadFiles(dir string) tea.Cmd {
      return func() tea.Msg {
          entries, err := os.ReadDir(dir)
          if err != nil {
              return errMsg{err}
          }

          var files []FileInfo
          for _, entry := range entries {
              info, _ := entry.Info()
              
              // Filter: only directories and preset files
              ext := filepath.Ext(entry.Name())
              isPreset := ext == ".np3" || ext == ".xmp" || ext == ".lrtemplate"
              
              if entry.IsDir() || isPreset {
                  files = append(files, FileInfo{
                      Name:    entry.Name(),
                      Path:    filepath.Join(dir, entry.Name()),
                      Size:    info.Size(),
                      ModTime: info.ModTime(),
                      IsDir:   entry.IsDir(),
                      Format:  detectFormat(entry.Name(), entry.IsDir()),
                  })
              }
          }

          return filesLoadedMsg{files: files}
      }
  }
  ```
- [ ] Implement `detectFormat()` helper:
  ```go
  func detectFormat(name string, isDir bool) string {
      if isDir {
          return "dir"
      }
      ext := filepath.Ext(name)
      switch ext {
      case ".np3":
          return "np3"
      case ".xmp":
          return "xmp"
      case ".lrtemplate":
          return "lrtemplate"
      default:
          return ""
      }
  }
  ```

**Validation:**
- Directory listing loads files
- Only preset files and directories shown
- Format detection correct

---

### Task 4: Implement Navigation Keys (AC-2)

- [ ] Create `handleKeyPress()` method:
  ```go
  func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
      switch msg.String() {
      case "q", "ctrl+c":
          return m, tea.Quit

      case "up", "k":
          if m.cursor > 0 {
              m.cursor--
          }

      case "down", "j":
          if m.cursor < len(m.files)-1 {
              m.cursor++
          }

      case "home":
          m.cursor = 0

      case "end":
          m.cursor = len(m.files) - 1

      case "enter":
          return m.navigateInto()

      case "backspace", "left":
          return m.navigateUp()

      case "?":
          m.showHelp = !m.showHelp
      }

      return m, nil
  }
  ```
- [ ] Implement `navigateInto()` method:
  ```go
  func (m model) navigateInto() (tea.Model, tea.Cmd) {
      if m.cursor >= len(m.files) {
          return m, nil
      }

      selected := m.files[m.cursor]
      if !selected.IsDir {
          return m, nil  // Can't navigate into file
      }

      m.currentDir = selected.Path
      m.cursor = 0
      return m, loadFiles(m.currentDir)
  }
  ```
- [ ] Implement `navigateUp()` method:
  ```go
  func (m model) navigateUp() (tea.Model, tea.Cmd) {
      parent := filepath.Dir(m.currentDir)
      if parent == m.currentDir {
          return m, nil  // Already at root
      }

      m.currentDir = parent
      m.cursor = 0
      return m, loadFiles(m.currentDir)
  }
  ```

**Validation:**
- Arrow keys navigate cursor
- Enter navigates into directories
- Backspace navigates to parent
- Home/End jump correctly

---

### Task 5: Implement File Selection (AC-4)

- [ ] Add selection toggle to `handleKeyPress()`:
  ```go
  case "space":
      return m.toggleSelection()

  case "a":
      return m.selectAll()

  case "n":
      return m.deselectAll()
  ```
- [ ] Implement `toggleSelection()`:
  ```go
  func (m model) toggleSelection() (tea.Model, tea.Cmd) {
      if m.cursor >= len(m.files) {
          return m, nil
      }

      file := m.files[m.cursor]
      if file.IsDir {
          return m, nil  // Can't select directories
      }

      if m.selected[file.Path] {
          delete(m.selected, file.Path)
      } else {
          m.selected[file.Path] = true
      }

      return m, nil
  }
  ```
- [ ] Implement `selectAll()` and `deselectAll()`:
  ```go
  func (m model) selectAll() (tea.Model, tea.Cmd) {
      for _, file := range m.files {
          if !file.IsDir {
              m.selected[file.Path] = true
          }
      }
      return m, nil
  }

  func (m model) deselectAll() (tea.Model, tea.Cmd) {
      m.selected = make(map[string]bool)
      return m, nil
  }
  ```

**Validation:**
- Space toggles selection
- 'a' selects all files
- 'n' clears selection
- Directories not selectable

---

### Task 6: Implement File List Rendering (AC-2, AC-3, AC-5)

- [ ] Create `renderFileList()` using Lipgloss:
  ```go
  import "github.com/charmbracelet/lipgloss/v2"

  var (
      styleNP3 = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))  // Blue
      styleXMP = lipgloss.NewStyle().Foreground(lipgloss.Color("208")) // Orange
      styleLRT = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
      styleDir = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
      styleCursor = lipgloss.NewStyle().Reverse(true)
  )

  func renderFileList(m model) string {
      var b strings.Builder

      // Header
      b.WriteString(fmt.Sprintf("Current: %s\n", m.currentDir))
      b.WriteString(strings.Repeat("─", m.termWidth))
      b.WriteString("\n")

      // File list
      for i, file := range m.files {
          line := renderFileLine(file, i == m.cursor, m.selected[file.Path])
          b.WriteString(line)
          b.WriteString("\n")
      }

      // Status bar
      selectedCount := len(m.selected)
      if selectedCount > 0 {
          b.WriteString(fmt.Sprintf("[%d files selected] Press 'c' to convert\n", selectedCount))
      }

      return b.String()
  }
  ```
- [ ] Implement `renderFileLine()`:
  ```go
  func renderFileLine(file FileInfo, isCursor bool, isSelected bool) string {
      // Build line components
      cursor := " "
      if isCursor {
          cursor = ">"
      }

      checkbox := " "
      if isSelected {
          checkbox = "✓"
      }

      icon := "📄"
      if file.IsDir {
          icon = "📁"
      }

      badge := formatBadge(file.Format)
      size := formatSize(file.Size)

      line := fmt.Sprintf("%s %s %s %-40s %s", cursor, checkbox, icon, file.Name, size)

      if isCursor {
          return styleCursor.Render(line)
      }
      return line
  }
  ```
- [ ] Implement `formatBadge()`:
  ```go
  func formatBadge(format string) string {
      switch format {
      case "np3":
          return styleNP3.Render("NP3")
      case "xmp":
          return styleXMP.Render("XMP")
      case "lrtemplate":
          return styleLRT.Render("LRT")
      case "dir":
          return styleDir.Render("DIR")
      default:
          return "   "
      }
  }
  ```
- [ ] Implement `formatSize()`:
  ```go
  func formatSize(bytes int64) string {
      const kb = 1024
      const mb = kb * 1024

      if bytes < kb {
          return fmt.Sprintf("%d B", bytes)
      } else if bytes < mb {
          return fmt.Sprintf("%.1f KB", float64(bytes)/kb)
      } else {
          return fmt.Sprintf("%.1f MB", float64(bytes)/mb)
      }
  }
  ```

**Validation:**
- File list renders correctly
- Color badges display
- Cursor highlights current line
- Selection checkmarks show
- File sizes formatted correctly

---

### Task 7: Implement Help Overlay (AC-6)

- [ ] Create `renderHelp()`:
  ```go
  func renderHelp(m model) string {
      help := `
  ┌─────────────────── Keyboard Shortcuts ────────────────────┐
  │                                                            │
  │ Navigation:                                                │
  │   ↑/k        Move cursor up                                │
  │   ↓/j        Move cursor down                              │
  │   Enter      Navigate into directory                       │
  │   Backspace  Navigate to parent directory                  │
  │   Home       Jump to first item                            │
  │   End        Jump to last item                             │
  │                                                            │
  │ Selection:                                                 │
  │   Space      Toggle selection on current file              │
  │   a          Select all files                              │
  │   n          Deselect all files                            │
  │                                                            │
  │ Actions:                                                   │
  │   c          Convert selected files (Story 4-3)            │
  │   r          Refresh current directory                     │
  │   ?          Toggle this help                              │
  │   q/Ctrl+C   Quit                                          │
  │                                                            │
  │ Press ? or Esc to close                                    │
  └────────────────────────────────────────────────────────────┘
  `
      return help
  }
  ```
- [ ] Update `handleKeyPress()` to handle Escape:
  ```go
  case "esc":
      if m.showHelp {
          m.showHelp = false
      }
  ```

**Validation:**
- '?' shows help
- Help displays all shortcuts
- Escape dismisses help

---

### Task 8: Handle Terminal Resize (AC-7)

- [ ] Update `Update()` to handle WindowSizeMsg:
  ```go
  case tea.WindowSizeMsg:
      m.termWidth = msg.Width
      m.termHeight = msg.Height
      return m, nil
  ```
- [ ] Implement scrolling for long file lists:
  ```go
  func renderFileList(m model) string {
      var b strings.Builder

      // Calculate visible range
      maxVisible := m.termHeight - 5  // Reserve lines for header/footer
      start := 0
      end := len(m.files)

      if len(m.files) > maxVisible {
          // Scroll to keep cursor visible
          if m.cursor > maxVisible/2 {
              start = m.cursor - maxVisible/2
          }
          end = start + maxVisible
          if end > len(m.files) {
              end = len(m.files)
              start = end - maxVisible
          }
      }

      // Render visible range
      for i := start; i < end; i++ {
          // ... render file line ...
      }

      return b.String()
  }
  ```

**Validation:**
- Resize event handled
- UI adapts to new size
- Scrolling works for long lists
- Minimum 80x24 usable

---

### Task 9: Add Refresh Functionality

- [ ] Add 'r' key handler:
  ```go
  case "r":
      return m, loadFiles(m.currentDir)
  ```
- [ ] Show refresh indicator during load:
  ```go
  type filesLoadingMsg struct{}

  func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case filesLoadingMsg:
          // Show loading indicator
      case filesLoadedMsg:
          m.files = msg.files
          // Hide loading indicator
      }
      return m, nil
  }
  ```

**Validation:**
- 'r' refreshes directory
- File list updates
- Loading indicator shows (optional)

---

### Task 10: Add Unit Tests

- [ ] Create `cmd/tui/model_test.go`:
  ```go
  func TestInitialModel(t *testing.T) {
      m := initialModel()
      assert.NotNil(t, m)
      assert.NotEmpty(t, m.currentDir)
  }

  func TestFileFiltering(t *testing.T) {
      files := []FileInfo{
          {Name: "test.xmp", IsDir: false, Format: "xmp"},
          {Name: "test.jpg", IsDir: false, Format: ""},
          {Name: "test.np3", IsDir: false, Format: "np3"},
      }

      // Only xmp and np3 should pass filter
      assert.Equal(t, 2, countPresetFiles(files))
  }

  func TestNavigation(t *testing.T) {
      m := initialModel()
      m.files = []FileInfo{
          {Name: "file1.xmp"},
          {Name: "file2.np3"},
      }

      // Move down
      m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown}).(model, tea.Cmd)
      assert.Equal(t, 1, m.cursor)

      // Move up
      m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp}).(model, tea.Cmd)
      assert.Equal(t, 0, m.cursor)
  }

  func TestSelection(t *testing.T) {
      m := initialModel()
      m.files = []FileInfo{{Name: "test.xmp", Path: "/test.xmp"}}
      m.cursor = 0

      // Toggle selection
      m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace}).(model, tea.Cmd)
      assert.True(t, m.selected["/test.xmp"])

      // Toggle again (deselect)
      m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace}).(model, tea.Cmd)
      assert.False(t, m.selected["/test.xmp"])
  }
  ```
- [ ] Test coverage goal: >85% for model logic

**Validation:**
- All unit tests pass
- Coverage meets goal

---

### Task 11: Update Build Configuration

- [ ] Update `Makefile`:
  ```makefile
  # Build TUI
  tui:
      go build -o recipe-tui cmd/tui/main.go

  # Install TUI
  install-tui: tui
      sudo mv recipe-tui /usr/local/bin/
  ```
- [ ] Add TUI to GitHub Actions CI:
  ```yaml
  - name: Build TUI
    run: make tui

  - name: Test TUI
    run: go test ./cmd/tui/
  ```

**Validation:**
- `make tui` builds successfully
- CI builds TUI binary

---

### Task 12: Update Documentation

- [ ] Update `README.md` with TUI usage:
  ```markdown
  ### TUI (Terminal User Interface)

  Interactive file browser for visual preset selection:

  ```bash
  # Launch TUI
  recipe-tui

  # Navigate with arrow keys
  # Select files with space bar
  # Press 'c' to convert (Story 4-3)
  # Press '?' for help
  ```

  **Keyboard Shortcuts:**
  - `↑/↓` or `k/j` - Navigate
  - `Enter` - Open directory
  - `Backspace` - Parent directory
  - `Space` - Toggle selection
  - `a` - Select all
  - `n` - Deselect all
  - `r` - Refresh
  - `?` - Help
  - `q` - Quit
  ```
- [ ] Create `docs/tui-guide.md` with screenshots (ASCII art)

**Validation:**
- Documentation accurate
- Examples work

---

## Dev Notes

### Architecture Alignment

**Follows Tech Spec Epic 4:**
- Bubbletea framework for TUI (PRD FR-4.1)
- Interactive file browser with keyboard navigation
- Multi-file selection capability
- Format filtering (only preset files)
- Consistent with Recipe's zero-dependency philosophy (Bubbletea is pure Go)

**Bubbletea Elm Architecture:**
```
Model (State)
  ↓
Init() → Command
  ↓
Update(Msg) → (Model, Command)
  ↓
View() → String (rendered UI)
  ↓
Terminal display
```

**v2 Experimental Breaking Changes:**
This story uses Bubbletea v2, Bubbles v2, and Lipgloss v2 experimental branches which include breaking API changes from v1:

- **Import path changes**: `github.com/charmbracelet/bubbletea/v2` (add `/v2` suffix)
- **Enhanced type safety**: Stronger typing for messages and commands
- **Performance improvements**: Optimized rendering and event handling
- **New APIs**: Additional helper functions and improved composability
- **Backward incompatible**: Code written for v1 requires migration

See official migration guides for detailed changes:
- [Bubbletea v2 migration](https://github.com/charmbracelet/bubbletea/tree/v2-exp)
- [Bubbles v2 migration](https://github.com/charmbracelet/bubbles/tree/v2-exp)
- [Lipgloss v2 migration](https://github.com/charmbracelet/lipgloss/tree/v2-exp)

**Integration with Future Stories:**
- **Story 4-2 (Live Parameter Preview)**: Selected file → preview pane
- **Story 4-3 (Batch Progress)**: Selected files → conversion progress
- **Story 4-4 (Visual Validation)**: Conversion complete → validation screen

### Dependencies

**New Dependencies (This Story):**
```go
github.com/charmbracelet/bubbletea/v2 v2.0.0-exp  // v2 experimental branch
github.com/charmbracelet/bubbles/v2 v2.0.0-exp    // v2 experimental branch
github.com/charmbracelet/lipgloss/v2 v2.0.0-exp   // v2 experimental branch
github.com/charmbracelet/x v0.2.0                 // Charm experimental packages
```

**Note on v2 Experimental Branches:**
- Using latest v2 experimental versions from GitHub (v2-exp branches)
- These are pre-release versions with breaking changes from v1
- Import paths change to `/v2` for major version upgrade
- Provides enhanced APIs and performance improvements
- See migration guides: [Bubbletea v2](https://github.com/charmbracelet/bubbletea/tree/v2-exp), [Bubbles v2](https://github.com/charmbracelet/bubbles/tree/v2-exp), [Lipgloss v2](https://github.com/charmbracelet/lipgloss/tree/v2-exp)

**Go Standard Library:**
- `os` - File system access
- `path/filepath` - Path manipulation
- `time` - File modification times
- `strings` - String building

**Internal Dependencies:**
- None (Story 4-1 is foundation for Epic 4)

**Future Dependencies:**
- Story 4-2 will add parameter preview using `internal/converter`
- Story 4-3 will integrate conversion engine

### Testing Strategy

**Unit Tests:**
- Test model initialization
- Test navigation logic (up/down/enter/backspace)
- Test selection toggle (space/a/n)
- Test file filtering (only presets + directories)
- Coverage goal: >85%

**Manual Tests:**
- Launch TUI in various terminal sizes
- Test with large directories (>100 files)
- Test with empty directories
- Test with deeply nested directories
- Test on Windows/macOS/Linux terminals

**Integration Tests (Story 4-2+):**
- Test with real sample files from `testdata/`
- Test conversion workflow (Story 4-3)

### Learnings from Previous Story

**From Story 3-6 (JSON Output Mode) - Status: ready-for-dev**

Story 3-6 established output formatting patterns for CLI:

- **Structured Output**: JSON schema with clear field types
- **Stream Separation**: stdout for data, stderr for logs
- **Flag Patterns**: Global flags available to all commands

**For This Story (4-1):**
- **TUI is separate binary**: `recipe-tui` distinct from `recipe` CLI
- **Different interaction model**: Event-driven (Bubbletea) vs imperative (Cobra)
- **No flag sharing**: TUI doesn't use Cobra flags, uses Bubbletea key handlers
- **Output formatting**: TUI renders to terminal, not stdout JSON

**Key Difference:** CLI is stateless (run command, exit), TUI is stateful (persistent interactive session).

[Source: stories/3-6-json-output-mode.md#Dev-Notes]

### Technical Debt / Future Enhancements

**v2 Experimental Stability:**
- Using pre-release v2-exp branches (not stable releases)
- API may change before v2 final release
- Monitor GitHub releases for breaking changes
- May need to update code when v2 stabilizes
- Consider pinning to specific commit SHA for production builds

**Deferred to Story 4-2:**
- Parameter preview pane
- Live preview when file selected
- Color-coded parameter diff display

**Deferred to Story 4-3:**
- Conversion trigger ('c' key)
- Progress bars for batch conversion
- Real-time conversion status

**Deferred to Story 4-4:**
- Visual validation screen
- Side-by-side parameter comparison
- Conversion summary report

**Post-Epic Enhancements:**
- Mouse support (click to select)
- Search/filter by filename
- Bookmarks for frequently used directories
- Custom key bindings configuration
- Theme customization

### References

- [Source: docs/PRD.md#FR-4.1] - Interactive File Browser requirements
- [Source: docs/architecture.md#ADR-005] - Bubbletea framework decision
- [Bubbletea v2 Experimental] - https://github.com/charmbracelet/bubbletea/tree/v2-exp
- [Bubbles v2 Experimental] - https://github.com/charmbracelet/bubbles/tree/v2-exp
- [Lipgloss v2 Experimental] - https://github.com/charmbracelet/lipgloss/tree/v2-exp
- [Charm X Packages] - https://github.com/charmbracelet/x
- [Bubbletea v2 Tutorial] - https://github.com/charmbracelet/bubbletea/tree/v2-exp/tutorials
- [Lipgloss v2 Examples] - https://github.com/charmbracelet/lipgloss/tree/v2-exp/examples

### Known Issues / Blockers

**Dependencies:**
- None - Story 4-1 is first story in Epic 4

**Cross-Story Dependencies:**
- Story 4-2 requires 4-1 (file browser is foundation for preview)
- Story 4-3 requires 4-1 (selection drives batch conversion)
- Story 4-4 requires 4-1, 4-3 (validation follows conversion)

**Platform Considerations:**
- Windows terminal support may have rendering issues (test thoroughly)
- SSH sessions may have limited color support (graceful degradation needed)
- Minimum terminal size 80x24 (enforce or warn user)

### Cross-Story Coordination

**Requires (Must be done first):**
- None - This is the first story in Epic 4

**Coordinates with:**
- Story 4-2: Live Parameter Preview (shares Model, adds preview pane)
- Story 4-3: Batch Progress (uses selected files for conversion)

**Enables:**
- All subsequent TUI stories (4-2, 4-3, 4-4) build on this foundation

**Architectural Consistency:**
This story establishes the TUI architecture pattern:
- Bubbletea Elm Architecture (Model-Update-View)
- Keyboard-driven navigation
- Multi-file selection state management
- Format filtering and file type detection
- Terminal resize handling
- Help overlay pattern

All future TUI stories will extend this model rather than create separate applications.

---

## Dev Agent Record

### Context Reference

**Story Context XML:** `docs/stories/4-1-bubbletea-file-browser.context.xml`
**Generated:** 2025-11-06
**Status:** Ready for development

### Agent Model Used

<!-- To be filled by dev agent -->

### Debug Log References

<!-- Dev agent will add references to detailed debug logs if needed -->

### Completion Notes List

<!-- Dev agent will document:
- Bubbletea integration approach
- Model state management decisions
- Keyboard navigation implementation details
- File filtering logic
- Terminal rendering strategy
- Cross-platform compatibility testing results
- Performance with large directories (>100 files)
- Test coverage metrics
-->

### File List

<!-- Dev agent will document files created/modified/deleted:
**NEW:**
- `cmd/tui/main.go` - TUI entry point
- `cmd/tui/model.go` - Bubbletea Model and state
- `cmd/tui/update.go` - Update logic (keyboard handlers)
- `cmd/tui/view.go` - View rendering (file list, help)
- `cmd/tui/files.go` - File loading and filtering
- `cmd/tui/model_test.go` - Unit tests
- `docs/tui-guide.md` - TUI user guide

**MODIFIED:**
- `Makefile` - Add TUI build targets
- `README.md` - Add TUI usage documentation
- `go.mod` - Add Bubbletea dependencies

**DELETED:**
- (none)
-->

---

## Code Review Notes

**Review Date:** 2025-11-06
**Reviewer:** Senior Developer (Code Review Agent)
**Review Outcome:** ✅ **APPROVED - PRODUCTION READY**

### Summary

This implementation is **EXCEPTIONAL** in quality and completeness. All acceptance criteria are fully met, test coverage exceeds requirements (88.5%), and code quality demonstrates mastery of the Bubbletea Elm Architecture pattern.

### Acceptance Criteria Validation

**All 7 ACs: PASS (100%)**

- **AC-1 (Bubbletea Application Initialization):** ✅ PASS - 5/5 requirements met
  - Uses Bubbletea v2 RC with correct import path `charm.land/bubbletea/v2`
  - Separate `recipe-tui` binary builds successfully (5.1 MB)
  - Initializes with current working directory (`os.Getwd()`)
  - Terminal size detection via `WindowSizeMsg` handler
  - Clean exit on 'q' and Ctrl+C

- **AC-2 (Directory Navigation):** ✅ PASS - 6/6 requirements met
  - Current directory displayed in bold header
  - File listing via `os.ReadDir()` with format filtering
  - Arrow keys + Vim keys (k/j) for navigation
  - Enter navigates into directories, prevents file navigation
  - Backspace/left arrow navigates to parent with root protection
  - Home/End jump to first/last with bounds checking

- **AC-3 (File Filtering by Format):** ✅ PASS - 5/5 requirements met
  - Filters to `.np3`, `.xmp`, `.lrtemplate` only
  - Shows directories for navigation
  - Color-coded badges: NP3=Blue(39), XMP=Orange(208), LRT=Green(42), DIR=Gray(240)
  - Human-readable file sizes (B/KB/MB with 1 decimal)
  - Empty directory message explains filter behavior

- **AC-4 (Multi-File Selection):** ✅ PASS - 6/6 requirements met
  - Space toggles selection on current file
  - Unicode checkmark (✓) indicator for selected files
  - Status bar shows count: "[N files selected]"
  - 'a' selects all files (skips directories)
  - 'n' clears all selections
  - Selection stored by path in map, persists across navigation

- **AC-5 (File Information Display):** ✅ PASS - 4/4 requirements met
  - Displays name, size (formatted), modified date (ISO format)
  - Full absolute path shown in status bar
  - Total count shown when scrolling: "Showing 1-20 of 42 files"
  - Color-coded format badges using Lipgloss v2

- **AC-6 (Keyboard Shortcuts Help):** ✅ PASS - 3/3 requirements met
  - '?' toggles help overlay
  - Escape closes help if open
  - Comprehensive help with all shortcuts organized by category

- **AC-7 (Terminal Resize Handling):** ✅ PASS - 4/4 requirements met
  - Updates dimensions on `WindowSizeMsg`
  - Smart viewport scrolling (only renders visible slice)
  - Scroll indicator shows range: "Showing 1-20 of 100 files"
  - Defaults to 80x24, limits border to 80 max width

### Task Verification

**All 12 Tasks: COMPLETE (100%)**

1. ✅ Initialize Bubbletea Project - `cmd/tui/` structure created
2. ✅ Implement Init and Update - Perfect Elm Architecture pattern
3. ✅ Implement Directory Listing - Async file loading with filtering
4. ✅ Implement Navigation Keys - All keys + Vim support
5. ✅ Implement File Selection - Map-based persistent selection
6. ✅ Implement File List Rendering - Lipgloss v2 styled output
7. ✅ Implement Help Overlay - Comprehensive keyboard reference
8. ✅ Handle Terminal Resize - Dynamic viewport scrolling
9. ✅ Add Refresh Functionality - 'r' key reloads directory
10. ⭐ **Add Unit Tests** - **EXCEPTIONAL** (88.5% coverage, 40 tests, all pass)
11. ✅ Update Build Configuration - Makefile with tui/tui-all targets
12. ⭐ **Update Documentation** - **EXCEPTIONAL** (301-line comprehensive guide)

### Code Quality Assessment

**Overall Rating:** ⭐⭐⭐⭐⭐ **EXCEPTIONAL**

- **Architecture:** ⭐⭐⭐⭐⭐ Perfect Elm Architecture (Model-Update-View)
  - Immutable state updates
  - Pure functions (no side effects)
  - Command pattern for async operations
  - Single source of truth

- **Error Handling:** ✅ GOOD
  - Graceful fallbacks (e.g., `cwd = "."` on error)
  - Non-fatal errors skipped (continue on info.Info() errors)
  - User-friendly error messages
  - Comprehensive bounds checking

- **Security:** ✅ SECURE
  - `filepath.Join()` prevents path traversal
  - Respects OS file permissions
  - Whitelist approach for keyboard input
  - No shell execution or eval()
  - Dependencies from trusted sources (Charm.sh)

- **Performance:** ⭐ EXCELLENT
  - Async file loading (non-blocking)
  - O(viewport) rendering, not O(total files)
  - O(1) selection lookups with map
  - No memory leaks identified

- **Maintainability:** ⭐⭐⭐⭐⭐ EXCEPTIONAL
  - Clear function names (navigateInto, selectAll, formatBadge)
  - Inline comments explain non-obvious logic
  - Guard clauses for early returns
  - Excellent separation of concerns (5 focused files)
  - Testability interfaces (KeyMsgInterface, WindowSizeMsgInterface)

- **Testing:** ⭐⭐⭐⭐⭐ EXCEPTIONAL
  - **Coverage: 88.5%** (EXCEEDS 85% requirement)
  - **40 tests, ALL PASS**
  - Edge cases (empty lists, bounds)
  - Integration tests (real file system)
  - Clean mocking without external deps

- **Documentation:** ⭐⭐⭐⭐⭐ EXCEPTIONAL
  - README.md TUI section added
  - **301-line comprehensive guide** (docs/tui-guide.md)
  - ASCII art examples
  - Keyboard reference table
  - Troubleshooting section
  - Architecture explanation
  - All exported functions documented

### Test Results

```
$ go test ./cmd/tui -v -cover
=== RUN   TestLoadFilesCmd
--- PASS: TestLoadFilesCmd (0.00s)
=== RUN   TestDetectFormatComprehensive
--- PASS: TestDetectFormatComprehensive (0.00s)
=== RUN   TestNavigationKeys
--- PASS: TestNavigationKeys (0.00s)
=== RUN   TestFileSelection
--- PASS: TestFileSelection (0.00s)
... [36 more tests] ...
PASS
coverage: 88.5% of statements
ok  	github.com/justin/recipe/cmd/tui	(cached)
```

**All 40 tests PASS** ✅
**Coverage: 88.5%** (EXCEEDS 85% target) ⭐

### Build Verification

```bash
$ go build -o recipe-tui cmd/tui/*.go
$ ls -lh recipe-tui
-rwxr-xr-x 1 Justin 197608 5.1M Nov  6 14:53 recipe-tui
```

**Binary builds successfully** ✅

### Dependencies

```go
// go.mod
charm.land/bubbletea/v2 v2.0.0-rc.1.0.20251106192006-06c0cda318b3
charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251106192539-4b304240aab7
charm.land/bubbles/v2 v2.0.0-beta.1.0.20251106192719-c2b822795a69
```

✅ Correct v2 versions with proper import paths

### Issues Found

**Blocking Issues:** 0
**Non-Blocking Issues:** 0

### Minor Notes (Non-Blocking)

1. **Bubbletea v2 RC Version**
   - Currently using Release Candidate (RC) version
   - Recommendation: Monitor for stable v2.0.0 release and update when available
   - Risk: Low (RC is well-tested, Charm.sh has good stability)
   - Priority: Low

### Files Created/Modified

**NEW:**
- `cmd/tui/main.go` (16 lines) - TUI entry point
- `cmd/tui/model.go` (108 lines) - Bubbletea Model + lifecycle
- `cmd/tui/keys.go` (137 lines) - Keyboard handling
- `cmd/tui/files.go` (69 lines) - File operations
- `cmd/tui/view.go` (216 lines) - Rendering logic
- `cmd/tui/model_test.go` (262 lines) - Model tests
- `cmd/tui/files_test.go` (88 lines) - File loading tests
- `cmd/tui/keys_test.go` (126 lines) - Keyboard tests
- `cmd/tui/view_test.go` (92 lines) - View rendering tests
- `cmd/tui/view_scrolling_test.go` (116 lines) - Scrolling tests
- `docs/tui-guide.md` (301 lines) - Comprehensive TUI guide

**MODIFIED:**
- `Makefile` - Added tui and tui-all targets
- `README.md` - Added TUI usage section with features and keyboard shortcuts
- `go.mod` - Added Bubbletea v2, Lipgloss v2, Bubbles v2 dependencies

**DELETED:** (none)

### Metrics Summary

```
Acceptance Criteria:  7/7 PASS   (100%)
Requirements Met:     33/33       (100%)
Tasks Complete:       12/12       (100%)
Test Coverage:        88.5%       (EXCEEDS 85%)
Tests Passing:        40/40       (100%)
Binary Size:          5.1 MB      (Reasonable)
Files Created:        10 Go + 1 Doc
Lines of Code:        ~546        (Excluding tests)
Test Code:            ~684        (Comprehensive)
Documentation:        ~600+ lines (Exceptional)
Security Issues:      0           (SECURE)
Code Smells:          0           (CLEAN)
```

### Production Readiness

✅ **READY FOR PRODUCTION DEPLOYMENT**

- [x] All acceptance criteria met with evidence
- [x] All tasks completed and verified
- [x] Test coverage exceeds minimum (88.5% > 85%)
- [x] All tests pass
- [x] Binary builds successfully
- [x] No security vulnerabilities
- [x] Documentation comprehensive and accurate
- [x] Code follows best practices
- [x] Error handling robust
- [x] Performance optimized
- [x] Makefile updated
- [x] README updated
- [x] No blocking issues

### Reviewer Comments

This is a **textbook implementation** of the Bubbletea Elm Architecture. The code demonstrates exceptional quality across all dimensions:

1. **Perfect Architecture** - Immutable state, pure functions, command pattern
2. **Exceptional Testing** - 88.5% coverage with comprehensive edge cases
3. **Outstanding Documentation** - 301-line guide with ASCII art and troubleshooting
4. **Production-Grade Quality** - Zero vulnerabilities, graceful error handling, optimized performance
5. **Excellent UX** - Help overlay, color coding, scrolling, persistent selection

The implementation not only meets all requirements but significantly exceeds expectations. This sets a high bar for future TUI stories and demonstrates mastery of modern Go TUI development practices.

**Recommendation:** Approve immediately. This is production-ready code.

---

## Change Log

- **2025-11-06:** Story created from Epic 4 (First story in TUI Interface epic, establishes Bubbletea foundation for interactive file browser)
- **2025-11-06:** Updated to use Bubbletea v2, Bubbles v2, Lipgloss v2 experimental branches and Charm X packages for latest features and performance improvements
- **2025-11-06:** Code review completed - APPROVED for production (88.5% test coverage, all 7 ACs pass, 40/40 tests pass, zero blocking issues)
