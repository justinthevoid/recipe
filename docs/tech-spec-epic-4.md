# Epic Technical Specification: TUI Interface

Date: 2025-11-06
Author: Justin
Epic ID: 4
Status: Draft

---

## Overview

Epic 4 delivers an interactive Terminal User Interface (TUI) for Recipe, providing a visual yet terminal-based experience for preset conversion workflows. Built with the Bubbletea framework (Elm-architecture pattern), this TUI bridges the gap between the simple CLI (Epic 3) and the full Web interface (Epic 2), offering power users an efficient, keyboard-driven workflow with real-time parameter preview and visual feedback.

The TUI reuses the proven conversion engine from Epic 1 (`internal/converter`), ensuring 95%+ accuracy and consistency across all Recipe interfaces. Unlike the CLI's fire-and-forget approach, the TUI provides interactive file browsing, live parameter extraction display, batch progress visualization with color-coded status, and confirmation screens before conversion—all while maintaining the terminal's speed and efficiency.

This epic targets photographers who prefer terminal-based workflows but want more visibility than the CLI provides, especially for batch operations where seeing parameter details and conversion progress in real-time improves confidence and reduces errors.

## Objectives and Scope

**In Scope:**
- Bubbletea-based interactive file browser with directory navigation (up/down arrow keys, enter to select)
- Multi-file selection using space bar with visual indicators (checkbox/highlight)
- Live parameter preview pane showing extracted preset parameters as files are selected
- Side-by-side layout: file list (left) | parameter preview (right)
- Real-time batch progress visualization with progress bars, file counts, and color-coded status
- Visual confirmation screen before conversion showing files, target format, and output location
- Keyboard-driven navigation (arrow keys, tab, enter, space, ESC, Ctrl+C)
- Terminal resize responsiveness using Bubbletea's built-in window size handling
- Format filtering in file browser (show only .np3, .xmp, .lrtemplate files)
- Unmappable parameter highlighting in preview (warnings for parameters that won't convert)
- Cancellable operations (Ctrl+C gracefully exits without partial conversions)
- Cross-platform terminal support (Windows Terminal, iTerm2, GNOME Terminal, etc.)

**Out of Scope:**
- Mouse interaction (keyboard-only to maintain terminal efficiency)
- GUI elements or graphics (pure text-based UI using box-drawing characters)
- Real-time file watching / daemon mode (single-session conversion workflow)
- Built-in text editor for preset parameter editing (view-only preview)
- Network-based conversion or remote file access (all processing local)
- Configuration persistence across sessions (stateless, fresh state each run)
- Custom color schemes or theming (use Lipgloss defaults for consistency)
- Interactive search/filter beyond format filtering (navigate via arrows)

**Deferred to Future Epics:**
- Advanced parameter diff view (compare two presets side-by-side)
- Undo/redo for batch operations (stateless design prevents this)
- Saved conversion profiles or templates
- Integration with file watchers for auto-conversion

## System Architecture Alignment

The TUI leverages the existing hub-and-spoke architecture from Epic 1, calling `converter.Convert()` as the single source of truth for all conversion logic. This ensures format consistency and validation rules are identical across Web, CLI, and TUI interfaces.

**Architecture References:**
- **Shared Library:** `internal/converter/converter.go` - Stateless Convert() function used by all interfaces
- **Format Parsers:** `internal/formats/{np3,xmp,lrtemplate}/parse.go` - Already implemented and tested in Epic 1
- **Format Generators:** `internal/formats/{np3,xmp,lrtemplate}/generate.go` - Already implemented and tested in Epic 1
- **TUI Framework:** Bubbletea v2 (github.com/charmbracelet/bubbletea/v2) - Elm-architecture pattern for terminal UIs
- **Styling Library:** Lipgloss v2 (github.com/charmbracelet/lipgloss/v2) - CSS-like styling for terminal output
- **Component Library:** Bubbles v2 (github.com/charmbracelet/bubbles/v2) - Reusable TUI components (list, progress, viewport)
- **Charm Utilities:** Charm X (github.com/charmbracelet/x) - Additional terminal utilities
- **Project Structure:** `cmd/tui/main.go` as entry point, following standard Go project layout
- **⚠️ CRITICAL:** All Charm libraries use v2 experimental branches with breaking changes from v1.x

**Key Constraints from Architecture:**
- Zero OS dependencies in converter package (maintains WASM compatibility from Epic 2)
- All file I/O happens in TUI layer (`cmd/tui/`), not in converter
- Errors wrapped as `ConversionError` type with operation/format context
- Logging via `slog` with structured fields (Go 1.21+ stdlib)
- Table-driven tests with real sample files from `testdata/` (1,501 files)

**Bubbletea Architecture Pattern:**
```
┌─────────────────────────────────────────────────────────────┐
│                    Elm Architecture Flow                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  User Input (Keyboard)                                      │
│         ↓                                                   │
│  Update(model, msg) → new Model + Commands                 │
│         ↓                                                   │
│  View(model) → Render String                               │
│         ↓                                                   │
│  Terminal Display                                           │
│         ↓ (loop)                                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**TUI State Management:**
- Model struct holds all UI state (selected files, current directory, preview data)
- Update function handles all state transitions (file selection, navigation, conversion)
- View function renders current state to terminal (pure function, no side effects)
- Commands for async operations (file I/O, conversion) executed outside Update
- No global state—all data in Model passed through Bubbletea runtime

## Detailed Design

### Services and Modules

| Module                                    | Responsibility                                            | Inputs                                     | Outputs                            | Owner           |
| ----------------------------------------- | --------------------------------------------------------- | ------------------------------------------ | ---------------------------------- | --------------- |
| **cmd/tui/main.go**                       | TUI entry point, Bubbletea program initialization         | Command-line args, keyboard events         | Exit code (0 or 1)                 | Epic 4          |
| **cmd/tui/model.go**                      | Core TUI state (Elm Model), holds all UI data             | Update messages (key presses, file events) | Updated Model struct               | Story 4-1       |
| **cmd/tui/update.go**                     | State transition logic (Elm Update), handles all messages | Model + tea.Msg                            | New Model + tea.Cmd                | All Stories     |
| **cmd/tui/view.go**                       | UI rendering (Elm View), generates terminal output        | Model                                      | Rendered string (terminal display) | All Stories     |
| **cmd/tui/filebrowser.go**                | File browser component, directory navigation              | Directory path, filter criteria            | File list, selection state         | Story 4-1       |
| **cmd/tui/preview.go**                    | Parameter preview component, displays preset details      | Selected file path                         | Formatted parameter list           | Story 4-2       |
| **cmd/tui/progress.go**                   | Batch progress tracking component                         | Conversion events (start, success, error)  | Progress bar, status text          | Story 4-3       |
| **cmd/tui/confirm.go**                    | Confirmation screen component                             | Selected files, target format              | User decision (yes/no/cancel)      | Story 4-4       |
| **cmd/tui/styles.go**                     | Lipgloss styling definitions                              | None                                       | Styled text renderers              | All Stories     |
| **cmd/tui/keys.go**                       | Keyboard binding definitions                              | Key presses                                | Key event handlers                 | All Stories     |
| **internal/converter**                    | **Shared conversion engine (Epic 1)**                     | Byte array, source/target formats          | Converted byte array or error      | Epic 1 (reused) |
| **internal/formats/{np3,xmp,lrtemplate}** | **Format parsers/generators (Epic 1)**                    | Format-specific bytes                      | UniversalRecipe or format bytes    | Epic 1 (reused) |

**Key Design Principles:**
- **Thin TUI Layer:** All business logic in `internal/converter`, TUI only handles UI state and rendering
- **Bubbletea Pattern:** Strict separation of Model (data), Update (logic), View (rendering)
- **Component Composition:** Each UI component (file browser, preview, progress) is self-contained Bubble
- **Stateless Conversion:** Each conversion call is independent, no persistent state between operations
- **Reactive Updates:** UI updates automatically when state changes (Elm architecture guarantee)

### Data Models and Contracts

**TUI State Model (Bubbletea Pattern):**

```go
// cmd/tui/model.go
package main

import (
    "github.com/charmbracelet/bubbles/v2/list"
    "github.com/charmbracelet/bubbles/v2/progress"
    "github.com/charmbracelet/bubbles/v2/viewport"
    tea "github.com/charmbracelet/bubbletea/v2"
)

// Model holds all TUI state
type Model struct {
    // Current screen/view
    CurrentView ViewState  // browse | preview | confirm | progress | error

    // File browser state
    FileBrowser    list.Model      // Bubbles list component for files
    CurrentDir     string          // Current directory path
    SelectedFiles  []FileItem      // Files marked for conversion (space bar)
    FilteredFiles  []FileItem      // Files matching format filter

    // Parameter preview state
    PreviewPane    viewport.Model  // Bubbles viewport for scrollable preview
    PreviewData    *PreviewContent // Parsed parameters from selected file
    PreviewError   error          // Parse error if file invalid

    // Conversion state
    TargetFormat   string          // np3 | xmp | lrtemplate
    OutputDir      string          // Where to save converted files
    ConversionMode ConversionMode  // single | batch

    // Batch progress state
    ProgressBar    progress.Model  // Bubbles progress component
    TotalFiles     int            // Count of files to convert
    ProcessedFiles int            // Count completed (success + error)
    SuccessCount   int            // Successful conversions
    ErrorCount     int            // Failed conversions
    Results        []ConversionResult  // Detailed results per file

    // Window size (for responsive layout)
    WindowWidth    int
    WindowHeight   int

    // Keyboard state
    Keys           KeyMap         // Key bindings
    Help           string         // Context-sensitive help text

    // Error state (for error view)
    LastError      error
}

type ViewState int
const (
    ViewBrowse ViewState = iota
    ViewPreview
    ViewConfirm
    ViewProgress
    ViewError
    ViewComplete
)

type ConversionMode int
const (
    ModeSingle ConversionMode = iota
    ModeBatch
)

// FileItem represents a file in the browser
type FileItem struct {
    Name       string
    Path       string
    Size       int64
    Format     string  // np3 | xmp | lrtemplate | unknown
    Selected   bool    // Marked for batch conversion
    IsDir      bool
}

// PreviewContent holds parsed preset parameters
type PreviewContent struct {
    FileName       string
    DetectedFormat string
    Parameters     map[string]interface{}  // Key-value pairs
    Warnings       []string                 // Unmappable parameters
    ParsedAt       time.Time
}

// ConversionResult (from CLI, reused)
type ConversionResult struct {
    InputFile    string
    OutputFile   string
    SourceFormat string
    TargetFormat string
    Success      bool
    Error        string
    Duration     time.Duration
    FileSize     int64
    Warnings     []string
}

// KeyMap defines keyboard shortcuts
type KeyMap struct {
    Up       key.Binding
    Down     key.Binding
    Select   key.Binding  // Space bar
    Confirm  key.Binding  // Enter
    Cancel   key.Binding  // ESC
    Quit     key.Binding  // q or Ctrl+C
    Help     key.Binding  // ?
}
```

**Message Types (Bubbletea Update Messages):**

```go
// cmd/tui/messages.go

// File browser messages
type dirChangedMsg struct{ path string }
type fileSelectedMsg struct{ file FileItem }
type filesLoadedMsg struct{ files []FileItem }

// Parameter preview messages
type previewLoadedMsg struct{ content *PreviewContent }
type previewErrorMsg struct{ err error }

// Conversion messages
type convertStartMsg struct{ files []FileItem, targetFormat string }
type convertProgressMsg struct{ current, total int, result ConversionResult }
type convertCompleteMsg struct{ results []ConversionResult }
type convertErrorMsg struct{ err error }

// Window resize
type tea.WindowSizeMsg  // Built-in Bubbletea message
```

**TUI does not introduce new conversion data models.** All conversion models defined in Epic 1 and Epic 3 (CLI):

| Model                | Definition                           | Usage in TUI                                        |
| -------------------- | ------------------------------------ | --------------------------------------------------- |
| **UniversalRecipe**  | `internal/model/recipe.go`           | Intermediate format (not directly exposed to user)  |
| **ConversionError**  | `internal/converter/errors.go`       | Wrapped errors displayed in error view              |
| **ConversionResult** | `cmd/cli/types.go` (reused from CLI) | Tracks per-file conversion status in batch progress |

### APIs and Interfaces

**TUI-to-Converter API (Existing from Epic 1):**

```go
// internal/converter/converter.go (already implemented)
func Convert(input []byte, sourceFormat, targetFormat string) ([]byte, error)
```

**New TUI Internal APIs:**

```go
// cmd/tui/filebrowser.go
func NewFileBrowser(initialDir string) (list.Model, error)
// Creates Bubbles list component pre-populated with files

func loadDirectory(path string) ([]FileItem, error)
// Scans directory, filters for preset formats (.np3, .xmp, .lrtemplate)
// Returns: Sorted file list with metadata (size, format detection)

func detectFileFormat(path string) (string, error)
// Returns: "np3" | "xmp" | "lrtemplate" based on extension
// Fallback: Read file header if extension ambiguous

// cmd/tui/preview.go
func loadPreview(filePath string) tea.Cmd
// Async command: Read file → Parse → Return previewLoadedMsg
// Returns: Bubbles command (executed outside Update)

func parsePresetFile(data []byte, format string) (*PreviewContent, error)
// Calls converter to parse file
// Returns: Parameter map + warnings for unmappable fields

func formatPreviewContent(content *PreviewContent) string
// Generates formatted text for viewport display
// Highlights warnings in different color (Lipgloss styling)

// cmd/tui/progress.go
func startConversion(files []FileItem, targetFormat string) tea.Cmd
// Async command: Converts files in parallel (worker pool like CLI)
// Sends convertProgressMsg for each file completion
// Returns: Final convertCompleteMsg when all done

func updateProgress(model *Model, result ConversionResult) *Model
// Updates progress bar, success/error counts
// Returns: Modified model with new state

// cmd/tui/confirm.go
func renderConfirmScreen(files []FileItem, targetFormat string, outputDir string) string
// Generates confirmation screen text showing conversion plan
// Returns: Formatted string for display

// cmd/tui/styles.go (Lipgloss styling)
func titleStyle() lipgloss.Style
func selectedStyle() lipgloss.Style
func errorStyle() lipgloss.Style
func successStyle() lipgloss.Style
func warningStyle() lipgloss.Style
func helpStyle() lipgloss.Style
// Each returns Lipgloss style for consistent UI theming

// cmd/tui/keys.go
func DefaultKeyMap() KeyMap
// Returns: Key bindings configuration
// Example: Up=↑, Down=↓, Select=Space, Confirm=Enter, Quit=q/Ctrl+C
```

**Bubbletea Core Functions (Implemented):**

```go
// cmd/tui/model.go
func New() Model
// Initializes Model with default state
// Returns: Fresh Model ready for Bubbletea.NewProgram()

// cmd/tui/update.go
func (m Model) Init() tea.Cmd
// Bubbletea lifecycle: Run once at startup
// Returns: Initial command (e.g., load current directory files)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
// Bubbletea lifecycle: Called for every message (keyboard, async results)
// Handles ALL state transitions via switch on msg type
// Returns: Updated model + optional command for async operations

// cmd/tui/view.go
func (m Model) View() string
// Bubbletea lifecycle: Renders current model state to string
// Called after every Update()
// Returns: Terminal output (ANSI-colored text using Lipgloss)
```

**Key Interactions:**

1. **File Browser Navigation:**
   ```go
   // User presses ↓
   msg := tea.KeyMsg{Type: tea.KeyDown}
   newModel, _ := model.Update(msg)
   // newModel.FileBrowser cursor moved down
   ```

2. **File Selection (Space Bar):**
   ```go
   // User presses space on file
   msg := tea.KeyMsg{Type: tea.KeySpace}
   newModel, _ := model.Update(msg)
   // Toggles file.Selected = !file.Selected
   // Updates SelectedFiles list
   ```

3. **Parameter Preview (Enter):**
   ```go
   // User presses enter on file
   msg := tea.KeyMsg{Type: tea.KeyEnter}
   newModel, cmd := model.Update(msg)
   // cmd = loadPreview(file.Path) → async file read

   // Later: Preview loaded
   msg := previewLoadedMsg{content: previewData}
   newModel, _ = model.Update(msg)
   // model.PreviewData = previewData
   // model.CurrentView = ViewPreview
   ```

4. **Batch Conversion:**
   ```go
   // User confirms conversion
   msg := tea.KeyMsg{Type: tea.KeyEnter}  // On confirm screen
   newModel, cmd := model.Update(msg)
   // cmd = startConversion(model.SelectedFiles, model.TargetFormat)

   // Progress updates
   msg := convertProgressMsg{current: 5, total: 10, result: {...}}
   newModel, _ = model.Update(msg)
   // Updates progress bar, increments ProcessedFiles

   // Completion
   msg := convertCompleteMsg{results: allResults}
   newModel, _ = model.Update(msg)
   // model.CurrentView = ViewComplete
   ```

**Exit Codes (Same as CLI):**

| Code | Meaning     | When                                           |
| ---- | ----------- | ---------------------------------------------- |
| 0    | Success     | User completed conversions or quit normally    |
| 1    | Error       | Conversion error or unexpected failure         |
| 2    | Usage Error | Invalid args (though TUI has no required args) |

### Workflows and Sequencing

**Single File Conversion Workflow:**

```
User launches TUI: recipe-tui
    ↓
1. Init() → Load current directory files
    ↓
2. View: File Browser (ViewBrowse)
   - Display file list with format badges
   - Show help text: "↑↓ navigate, Enter preview, Space select, q quit"
    ↓
3. User navigates with ↑↓ arrows
   - Update() handles tea.KeyUp/tea.KeyDown
   - FileBrowser cursor position changes
   - View() re-renders with highlighted file
    ↓
4. User presses Enter on file
   - Update() dispatches loadPreview(file.Path) command
   - View: Show "Loading preview..." spinner
    ↓
5. Async: loadPreview() → Read file → Parse parameters
   - Sends previewLoadedMsg back to Update()
    ↓
6. Update() receives previewLoadedMsg
   - model.PreviewData = content
   - model.CurrentView = ViewPreview
    ↓
7. View: Parameter Preview Screen (ViewPreview)
   Left pane: File details (name, size, format)
   Right pane: Scrollable parameter list
     - Exposure: +0.5
     - Contrast: +15
     - Saturation: -10
     ⚠ Warning: Parameter 'Grain' unmappable to NP3
   Help text: "Tab select format, Enter convert, ESC back"
    ↓
8. User presses Tab → Cycle target formats (np3 | xmp | lrtemplate)
   - Update() changes model.TargetFormat
   - View() updates format selector display
    ↓
9. User presses Enter → Confirm conversion
   - Update() sets model.CurrentView = ViewConfirm
    ↓
10. View: Confirmation Screen (ViewConfirm)
    ┌─────────────────────────────────────┐
    │ Confirm Conversion                  │
    ├─────────────────────────────────────┤
    │ File: portrait.xmp                  │
    │ Convert: xmp → np3                  │
    │ Output: portrait.np3                │
    │                                      │
    │ [Y] Confirm  [N] Cancel             │
    └─────────────────────────────────────┘
    ↓
11. User presses Y
    - Update() dispatches startConversion([file], targetFormat) command
    - model.CurrentView = ViewProgress
    ↓
12. View: Progress Screen (ViewProgress)
    Converting...
    [████████████████████████] 100%
    ✓ portrait.np3 (1.2 KB, 15ms)
    ↓
13. Conversion complete
    - Update() receives convertCompleteMsg
    - model.CurrentView = ViewComplete
    ↓
14. View: Success Screen (ViewComplete)
    ✓ Conversion Complete!
    Converted 1 file successfully
    Press any key to return to browser...
    ↓
15. User presses any key → Back to file browser (ViewBrowse)
```

**Batch Conversion Workflow:**

```
User in File Browser (ViewBrowse)
    ↓
1. Navigate to file with ↑↓
    ↓
2. Press Space to select file
   - Update() toggles file.Selected = true
   - FileBrowser shows checkbox ☑ next to file
   - model.SelectedFiles.append(file)
   - Help text updates: "5 files selected"
    ↓
3. Repeat for multiple files (space bar on each)
   - SelectedFiles list grows: [file1, file2, file3, ...]
    ↓
4. Press 'c' key (batch convert shortcut)
   - Update() validates: len(SelectedFiles) > 0
   - model.CurrentView = ViewConfirm
    ↓
5. View: Batch Confirmation Screen (ViewConfirm)
    ┌─────────────────────────────────────┐
    │ Batch Conversion Confirmation       │
    ├─────────────────────────────────────┤
    │ Files to convert: 5                 │
    │   • portrait1.xmp                   │
    │   • portrait2.xmp                   │
    │   • landscape.xmp                   │
    │   • ...and 2 more                   │
    │                                      │
    │ Target format: np3                  │
    │ Output directory: ./                │
    │                                      │
    │ [Y] Confirm  [N] Cancel  [E] Edit   │
    └─────────────────────────────────────┘
    ↓
6. User presses E → Edit settings
   - Tab through: Target format | Output directory
   - Enter to confirm choice
    ↓
7. User presses Y → Confirm batch
   - Update() dispatches startConversion(SelectedFiles, targetFormat) command
   - model.CurrentView = ViewProgress
   - Resets progress: ProcessedFiles=0, SuccessCount=0, ErrorCount=0
    ↓
8. Async: Parallel conversions begin (worker pool like CLI)
   - Each file completion sends convertProgressMsg
    ↓
9. View: Real-time Batch Progress (ViewProgress)
    Converting 5 files...
    [████████████░░░░░░░░░░] 60% (3/5)

    ✓ portrait1.np3 (1.1 KB, 12ms)
    ✓ portrait2.np3 (1.3 KB, 14ms)
    ✓ landscape.np3 (1.2 KB, 13ms)
    ⏳ portrait3.xmp converting...
    ⏳ portrait4.xmp converting...

    Elapsed: 0.5s | Est. remaining: 0.3s
    Ctrl+C to cancel
    ↓
10. Update() receives convertProgressMsg for each file
    - model.ProcessedFiles += 1
    - model.SuccessCount += 1 (or ErrorCount if error)
    - model.Results.append(result)
    - ProgressBar updates: current/total
    - View() re-renders with new progress
    ↓
11. All conversions complete
    - Update() receives convertCompleteMsg{results: [...]}
    - model.CurrentView = ViewComplete
    ↓
12. View: Batch Complete Screen (ViewComplete)
    ✓ Batch Conversion Complete!

    Total: 5 files
    Success: 4 ✓
    Errors: 1 ✗

    ✓ portrait1.np3
    ✓ portrait2.np3
    ✓ landscape.np3
    ✓ portrait4.np3
    ✗ portrait3.xmp - Parse error: Invalid XML

    Press any key to return to browser...
    ↓
13. User presses any key
    - Update() clears SelectedFiles
    - model.CurrentView = ViewBrowse
```

**Error Handling Workflow:**

```
At any point during conversion:
    ↓
If error occurs (file not found, parse failure, write error):
    ↓
1. Update() receives convertErrorMsg{err: error}
   - model.LastError = err
   - model.CurrentView = ViewError
    ↓
2. View: Error Screen (ViewError)
    ┌─────────────────────────────────────┐
    │ ✗ Conversion Error                  │
    ├─────────────────────────────────────┤
    │ Failed to convert portrait.xmp      │
    │                                      │
    │ Error: Parse error: Invalid XML     │
    │ Expected crs:xmpmeta root element   │
    │                                      │
    │ Suggestions:                        │
    │ • Re-export preset from Lightroom   │
    │ • Check file isn't corrupted        │
    │                                      │
    │ [R] Retry  [B] Back to browser      │
    └─────────────────────────────────────┘
    ↓
3. User presses B → Back to browser
   - Update() clears error
   - model.CurrentView = ViewBrowse
```

**Terminal Resize Handling:**

```
Terminal window resized (user drags window edge)
    ↓
1. Bubbletea sends tea.WindowSizeMsg
   - Update() receives msg
   - model.WindowWidth = msg.Width
   - model.WindowHeight = msg.Height
    ↓
2. View() recalculates layout
   - File browser width = WindowWidth / 2
   - Preview pane width = WindowWidth / 2
   - Progress bar width = WindowWidth - 4 (padding)
   - Viewport height = WindowHeight - 10 (header/footer)
    ↓
3. View() re-renders with new dimensions
   - No content lost, just reflowed
   - Scrollable viewports adjust automatically (Bubbles feature)
```

**Cancellation Workflow (Ctrl+C during batch):**

```
User presses Ctrl+C during batch conversion
    ↓
1. Bubbletea sends tea.KeyMsg{Type: tea.KeyCtrlC}
   - Update() checks CurrentView == ViewProgress
    ↓
2. If conversion in progress:
   - Set cancellation flag
   - Wait for current file to finish (atomic operation)
   - Stop starting new conversions
    ↓
3. Update() transitions to ViewComplete (partial results)
   - Shows: "✓ Converted 3/5 files (cancelled by user)"
   - Lists completed files
    ↓
4. User can review partial results
   - No partial file writes (atomic conversion)
   - Can retry remaining files manually
```

**Key State Transitions:**

```
ViewBrowse → ViewPreview (Enter on file)
ViewPreview → ViewConfirm (Enter after selecting format)
ViewConfirm → ViewProgress (Y to confirm)
ViewProgress → ViewComplete (conversion finishes)
ViewComplete → ViewBrowse (any key to continue)

ViewError can be reached from any state
ESC always returns to previous view (or quits from ViewBrowse)
```

## Non-Functional Requirements

### Performance

**UI Responsiveness:**
- **Keystroke Latency:** <16ms from key press to screen update (60 FPS equivalent)
- **File Browser Rendering:** <50ms to render file list with 1,000+ files
- **Preview Load Time:** <100ms to parse and display parameters for typical preset file
- **Screen Transitions:** <30ms for view changes (Browse → Preview → Confirm, etc.)
- **Progress Updates:** Real-time batch progress updates every 100ms (smooth animation)
- **Terminal Resize:** <50ms to recalculate and re-render layout on window size change

**Conversion Performance (Inherited from Epic 1/3):**
- **Single File:** <20ms per conversion (same as CLI target)
- **Batch Processing:** Parallel CPU utilization (use all cores via worker pool)
- **100 File Batch:** <2 seconds total (same as CLI target, progress visible in TUI)
- **Memory Efficiency:** TUI process uses <150MB for 1,000 file batch operations

**File I/O Performance:**
- **Directory Scan:** <200ms to load and filter 10,000 files
- **Format Detection:** <1ms per file (extension-based)
- **Parameter Preview Parsing:** <50ms for typical file (XMP ~50KB, lrtemplate ~100KB)
- **Lazy Loading:** Preview only parses when user selects file (not all files upfront)

**Optimization Strategies:**
- **Lazy Preview Loading:** Only parse file when user navigates to it (avoid parsing all files)
- **Viewport Virtualization:** Bubbles list component virtualizes large file lists (render only visible items)
- **Async Operations:** All file I/O and conversion happens in Bubbletea commands (non-blocking)
- **Worker Pool:** Batch conversions use parallel goroutines (same pattern as CLI)
- **Minimal Re-renders:** Bubbletea only re-renders when state changes (efficient diff algorithm)
- **Terminal Output Buffering:** Lipgloss batches ANSI escape codes to minimize terminal writes

**Performance Degradation Scenarios:**
- **Large Directories (10,000+ files):** Initial directory scan may take 1-2 seconds (show loading spinner)
- **Slow File System:** Network drives or USB may slow preview loading (show "Loading..." indicator)
- **Very Large Files (>1MB):** Parsing may take 200-500ms (acceptable, show spinner)
- **Mitigation:** All slow operations have visual feedback (spinners, progress bars)

### Security

**Local Processing Only:**
- Zero network requests during conversion (same as CLI/Web)
- No telemetry, analytics, or crash reporting
- No auto-update mechanisms (user-controlled updates)
- All processing happens on local machine in terminal session

**Input Validation (Inherited from Epic 1):**
- File size limits: Max 10MB per file (prevent memory exhaustion)
- Format validation: Magic bytes verification before parsing
- Path validation: Prevent directory traversal attacks
- Buffer overflow protection: Go's memory safety model

**File System Access:**
- Read-only access to input files
- Write-only access to specified output paths
- No directory traversal beyond user-specified paths
- Atomic writes: Use temp files, rename on success (no partial writes)
- Fail-safe: Refuse overwrite without user confirmation

**Terminal Security:**
- No execution of shell commands (pure Go file I/O)
- ANSI escape code sanitization: Lipgloss handles all terminal output (no raw escape codes)
- No external binary execution (all conversion in-process)
- Input sanitization: Key presses validated (no terminal injection)

**Dependency Security:**
- Bubbletea framework: Well-maintained, security-audited Go library
- Lipgloss: Pure Go, no C dependencies
- Bubbles: Official Charmbracelet component library
- Zero custom crypto or network code
- Go stdlib only for core conversion (minimal attack surface)

**Error Information Disclosure:**
- Error messages don't reveal system paths (sanitize file paths in display)
- No stack traces in production builds
- Debug mode: Additional details only when explicitly enabled via --debug flag
- Verbose logging: Controlled by user, not automatic

**State Isolation:**
- No persistent state across TUI sessions (stateless design)
- No configuration files that could be tampered with
- Fresh state on every launch
- No clipboard access or external data exfiltration

### Reliability/Availability

**Conversion Accuracy (Inherited from Epic 1):**
- 95%+ parameter fidelity target (validated via round-trip testing)
- Graceful handling of unmappable parameters (user warnings in preview)
- Consistent results across platforms (Windows, macOS, Linux)
- Deterministic output (same input always produces same output)

**Error Handling:**
- Clear error messages for all failure modes (file not found, parse errors, write failures)
- Error view shows actionable suggestions ("Re-export from Lightroom", "Check file permissions")
- No silent failures (all errors displayed to user)
- Retry capability from error screen (user can attempt conversion again)
- Partial batch recovery: Successful conversions saved even if some files fail

**Cross-Platform Terminal Compatibility:**
- **Windows:** Windows Terminal, PowerShell, CMD
- **macOS:** Terminal.app, iTerm2
- **Linux:** GNOME Terminal, Konsole, xterm, Alacritty
- **Terminal Feature Detection:** Bubbletea auto-detects terminal capabilities (color support, size)
- **Fallback:** Graceful degradation if terminal doesn't support full ANSI colors

**Batch Processing Resilience:**
- **Continue-on-error:** Process all files even if some fail (default behavior)
- **Atomic Conversions:** Each file conversion is independent (failure doesn't affect others)
- **Progress Preservation:** UI shows which files succeeded/failed during batch
- **Cancellation Safety:** Ctrl+C during batch doesn't leave partial files (atomic writes)
- **Result Summary:** Detailed report shows success/error counts with file-level details

**File Integrity:**
- Atomic writes: Use temp files, rename on success (no partial conversions)
- Verification: Generated files validated before overwriting existing files
- Original files never modified (read-only access)
- Output validation: Conversion engine guarantees well-formed output or error

**Terminal State Recovery:**
- **Clean Exit:** Restore terminal to normal state on exit (cursor visible, no alt screen artifacts)
- **Crash Handling:** Bubbletea cleanup functions restore terminal even on panic
- **Signal Handling:** Graceful shutdown on SIGINT (Ctrl+C), SIGTERM
- **Alt Screen:** TUI uses alternate screen buffer (main screen preserved when exiting)

**State Consistency:**
- **Elm Architecture Guarantee:** All state transitions through Update() function (no race conditions)
- **No Global State:** All data in Model struct passed through Bubbletea runtime
- **Immutable Updates:** Update() returns new Model (doesn't modify in-place, prevents bugs)
- **Single Source of Truth:** Model holds all UI state (no hidden state in components)

### Observability

**Visual Feedback (TUI-Specific):**
- **Real-time Progress:** Progress bars for batch conversions with percentage, file counts
- **Color-Coded Status:** Green (✓ success), Red (✗ error), Yellow (⚠ warning), Blue (⏳ in progress)
- **File-Level Details:** Each conversion shows filename, size, duration immediately after completion
- **Context-Sensitive Help:** Bottom of screen shows available commands for current view
- **Status Messages:** Top of screen shows current operation ("Converting...", "Loading preview...", etc.)

**Parameter Preview Visibility:**
- **Structured Display:** Parameters grouped by category (Basic, Color, Tone Curve, etc.)
- **Value Formatting:** Human-readable ranges (Exposure: +0.5 EV, Contrast: +15%)
- **Warning Highlighting:** Unmappable parameters shown in yellow with warning icon (⚠)
- **Scrollable View:** Long parameter lists scrollable with ↑↓ arrows
- **Source Format Badge:** Shows detected format (NP3, XMP, lrtemplate) with color coding

**Batch Operation Tracking:**
- **Overall Progress Bar:** Visual bar showing completion percentage [████░░░░] 60%
- **File Counter:** "Processing 45/100 files..."
- **Success/Error Tallies:** "Success: 42 ✓ | Errors: 3 ✗"
- **Per-File Results:** Scrollable list of completed files with status icons
- **Time Estimates:** "Elapsed: 1.2s | Est. remaining: 0.8s"
- **Live Updates:** Progress refreshes every 100ms during batch operations

**Error Reporting:**
- **Error View Screen:** Dedicated screen for displaying errors with context
- **Error Message:** User-friendly description (not technical stack traces)
- **File Context:** Shows which file caused error and at what stage (parse, generate, write)
- **Suggestions:** Actionable next steps ("Check file permissions", "Re-export from Lightroom")
- **Error Log:** Option to save error details to file for debugging (--debug mode)

**Logging Framework (Optional Debug Mode):**
- **Normal Mode:** No logging, clean TUI only
- **Debug Mode (--debug flag):** Structured logging to file via slog
  - Log file: `~/.recipe-tui.log`
  - Log levels: Debug, Info, Warn, Error
  - Structured fields: timestamp, view, operation, file, error

**Debug Mode Logging Example:**
```json
// ~/.recipe-tui.log (JSON structured logging)
{"time":"2025-11-06T10:15:30Z","level":"INFO","msg":"TUI started","view":"browse","dir":"/home/user/presets"}
{"time":"2025-11-06T10:15:35Z","level":"INFO","msg":"File selected","file":"portrait.xmp","format":"xmp"}
{"time":"2025-11-06T10:15:36Z","level":"DEBUG","msg":"Preview loading","file":"portrait.xmp","size":51234}
{"time":"2025-11-06T10:15:36Z","level":"INFO","msg":"Preview loaded","file":"portrait.xmp","params":42,"warnings":1}
{"time":"2025-11-06T10:15:40Z","level":"ERROR","msg":"Conversion failed","file":"corrupted.xmp","error":"parse error: invalid XML"}
```

**UI State Inspection (for development/debugging):**
- **Model Dump:** Press 'd' key (debug builds only) to dump model state to log
- **Terminal Info:** Display terminal capabilities (color support, size) in help screen
- **Version Info:** Show TUI version, Go version, Bubbletea version in about screen

**User Experience Metrics (Visible to User):**
- **Conversion Duration:** Per-file timing displayed immediately (e.g., "portrait.np3 (15ms)")
- **Batch Summary:** Total time for batch operations (e.g., "Converted 100 files in 1.5s")
- **File Size:** Output file size shown after conversion (e.g., "1.2 KB")
- **Warning Count:** Number of unmappable parameters shown in preview

**Accessibility Considerations:**
- **High Contrast Mode:** Lipgloss styles work well with terminal high contrast settings
- **Screen Reader:** Text-based UI more compatible with screen readers than GUI
- **Keyboard-Only:** 100% keyboard-driven (no mouse required)
- **Clear Labels:** All UI elements have descriptive text (no icon-only buttons)

## Dependencies and Integrations

### External Dependencies

**Bubbletea Framework:**
- **Package:** `github.com/charmbracelet/bubbletea/v2`
- **Version:** v2.0.0-exp (v2 experimental branch)
- **Purpose:** Elm-architecture TUI framework (Model-Update-View pattern)
- **Rationale:** Industry standard for Go TUIs, active maintenance, excellent documentation
- **Installation:** `go get github.com/charmbracelet/bubbletea/v2@v2-exp`
- **⚠️ CRITICAL:** v2 experimental version with breaking changes from v1. Import path includes `/v2` suffix. Migration guide: https://github.com/charmbracelet/bubbletea/tree/v2-exp

**Lipgloss Styling Library:**
- **Package:** `github.com/charmbracelet/lipgloss/v2`
- **Version:** v2.0.0-exp (v2 experimental branch)
- **Purpose:** CSS-like styling for terminal output (colors, borders, alignment)
- **Rationale:** Official Charmbracelet styling companion to Bubbletea
- **Installation:** `go get github.com/charmbracelet/lipgloss/v2@v2-exp`
- **⚠️ CRITICAL:** v2 experimental version with breaking changes from v1. Import path includes `/v2` suffix. Migration guide: https://github.com/charmbracelet/lipgloss/tree/v2-exp

**Bubbles Component Library:**
- **Package:** `github.com/charmbracelet/bubbles/v2`
- **Version:** v2.0.0-exp (v2 experimental branch)
- **Purpose:** Reusable TUI components (list, progress bar, viewport, spinner)
- **Components Used:**
  - `list.Model` - File browser
  - `progress.Model` - Batch progress bar
  - `viewport.Model` - Scrollable parameter preview
  - `spinner.Model` - Loading indicators
- **Rationale:** Official Charmbracelet component library, well-tested, composable
- **Installation:** `go get github.com/charmbracelet/bubbles/v2@v2-exp`
- **⚠️ CRITICAL:** v2 experimental version with breaking changes from v1. Import path includes `/v2` suffix. Migration guide: https://github.com/charmbracelet/bubbles/tree/v2-exp

**Charm X Experimental Utilities:**
- **Package:** `github.com/charmbracelet/x`
- **Version:** v0.2.0
- **Purpose:** Additional Charm utilities (ansi, input, term helpers)
- **Installation:** `go get github.com/charmbracelet/x@v0.2.0`

**Go Standard Library (No Version Constraints):**
- `os` - File I/O operations
- `path/filepath` - Cross-platform path handling, directory traversal
- `errors` - Error wrapping and unwrapping
- `log/slog` - Structured logging (Go 1.21+, debug mode only)
- `time` - Duration tracking, performance metrics, timestamps
- `runtime` - CPU count for parallel batch processing
- `sync` - Goroutine synchronization (worker pools for batch conversions)
- `io/ioutil` - File reading (preview loading)

### Internal Dependencies (Shared with All Interfaces)

**Core Conversion Engine:**
- **Package:** `github.com/justin/recipe/internal/converter`
- **API:** `Convert([]byte, string, string) ([]byte, error)`
- **Shared By:** TUI (Epic 4), Web/WASM (Epic 2), CLI (Epic 3)
- **Owner:** Epic 1 (Format Parsers/Generators)

**Format Parsers:**
- `github.com/justin/recipe/internal/formats/np3`
- `github.com/justin/recipe/internal/formats/xmp`
- `github.com/justin/recipe/internal/formats/lrtemplate`

**Data Models:**
- `github.com/justin/recipe/internal/model` (UniversalRecipe struct)

**Error Types:**
- `github.com/justin/recipe/internal/converter` (ConversionError)

### Integration Points

**File System:**
- **Read Operations:** `os.ReadFile()` for input files, directory scanning
- **Write Operations:** `os.WriteFile()` for output files (atomic writes via temp files)
- **Directory Operations:**
  - `os.ReadDir()` for file browser directory listing
  - `os.Stat()` for file metadata (size, permissions)
  - `filepath.Walk()` for recursive directory traversal (if needed)
- **Glob Patterns:** Not used in TUI (user navigates directories interactively)

**Terminal/Shell:**
- **Standard Output:** All TUI rendering to `os.Stdout` (Bubbletea handles this)
- **Standard Input:** Keyboard input via Bubbletea (handles all terminal escape sequences)
- **Terminal Control:**
  - Alt screen buffer (preserve main terminal screen)
  - Cursor visibility control (hide during operation, show on exit)
  - Terminal size detection (dynamic layout based on window size)
- **Exit Codes:** `os.Exit(0)` for success, `os.Exit(1)` for errors (same as CLI)

**Conversion Engine (Critical Integration):**
- **Interface Contract:**
  ```go
  func Convert(input []byte, sourceFormat, targetFormat string) ([]byte, error)
  ```
- **Input:** File bytes read via `os.ReadFile()` in loadPreview() and startConversion()
- **Output:** Converted bytes written via `os.WriteFile()` in conversion commands
- **Error Handling:** Wrap conversion errors with file context for display in error view

**Parallel Processing (Batch Conversions):**
- **Worker Pool Pattern:** Goroutines + channels (same as CLI Epic 3)
- **Concurrency:** `runtime.NumCPU()` workers (default, no user override in TUI)
- **Synchronization:** `sync.WaitGroup` for worker coordination
- **Progress Updates:** Each worker sends convertProgressMsg to Update() via channel

**Bubbletea Runtime Integration:**
- **Message Bus:** All async operations send messages back to Update() via tea.Cmd
- **Commands:** File I/O, conversion, preview loading all execute as Bubbletea commands
- **Lifecycle Hooks:**
  - `Init()` - Initial setup (load current directory)
  - `Update()` - State transitions (handle all messages)
  - `View()` - Rendering (called after every Update())
- **Subscriptions:** Window resize events (tea.WindowSizeMsg) automatically handled

### Build Dependencies

**Go Toolchain:**
- **Minimum Version:** Go 1.24.0+
- **Reason:** Enhanced WASM support (shared with Epic 2), slog stdlib
- **Verification:** `go version` must output go1.24.0 or higher

**Build Tools:**
- **Make:** Optional but recommended for build automation
- **No Additional Tools:** Bubbletea compiles to single binary (no external runtime)

### Development Dependencies

**Testing:**
- `testing` (stdlib) - Go test framework
- `testdata/` directory - 1,501 sample files for validation
- **Bubbletea Testing:** Use `Model.Update()` directly to test state transitions

**Benchmarking:**
- `testing` (stdlib) - Benchmark functions (though TUI benchmarks less critical than CLI)
- `time` (stdlib) - Performance measurement

**Linting/Formatting:**
- `gofmt` (stdlib) - Code formatting
- `go vet` (stdlib) - Static analysis

### Distribution Dependencies

**Cross-Platform Compilation:**
- No runtime dependencies required
- Single binary distribution (static linking)
- Supported platforms: Windows (amd64), macOS (amd64, arm64), Linux (amd64, arm64)
- Terminal requirements: ANSI color support, 80x24 minimum size

**Package Managers (Optional):**
- Homebrew (macOS/Linux) - For `brew install recipe-tui` distribution
- Scoop (Windows) - For `scoop install recipe-tui` distribution
- Neither required for manual installation (download binary from GitHub Releases)

### Dependency Version Management

**go.mod Entry (Expected):**
```go
module github.com/justin/recipe

go 1.25.1

require (
    github.com/charmbracelet/bubbletea/v2 v2.0.0-exp
    github.com/charmbracelet/lipgloss/v2 v2.0.0-exp
    github.com/charmbracelet/bubbles/v2 v2.0.0-exp
    github.com/charmbracelet/x v0.2.0
)
```

**⚠️ CRITICAL VERSION NOTES:**
- All Charm libraries use **v2 experimental branches** with `/v2` import path suffix
- These are pre-release versions with breaking API changes from v1.x
- Code written for v1.x will NOT work with v2 without migration
- Import paths MUST include `/v2`: `import tea "github.com/charmbracelet/bubbletea/v2"`
- Monitor upstream for changes: v2 experimental may have breaking changes before stable v2 release

**Dependency Update Strategy:**
- Pin exact versions in go.mod for reproducible builds
- Monitor Charmbracelet releases for security updates
- Test thoroughly before upgrading major versions
- Use `go mod tidy` to keep dependencies clean

### Zero Additional Dependencies Policy

**Explicitly NOT Using:**
- No database libraries (TUI doesn't persist data)
- No network/HTTP libraries (all processing is local)
- No configuration file parsers (TUI is stateless)
- No external CLI frameworks beyond Bubbletea (no Cobra in TUI, that's CLI-specific)
- No logging frameworks beyond slog (stdlib only)

**Rationale:**
- Minimal attack surface
- Faster builds
- Smaller binary size
- Simpler maintenance
- Consistent with project philosophy (privacy-first, zero-telemetry)

## Acceptance Criteria (Authoritative)

### AC-1: Interactive File Browser with Directory Navigation

**Given** a user launches the TUI (`recipe-tui`)
**When** the TUI starts
**Then** the file browser displays:
- Current directory path at top of screen
- List of files in current directory (filtered to .np3, .xmp, .lrtemplate)
- File metadata (name, size, format badge)
- Cursor on first file
- Help text showing available commands: "↑↓ navigate, Enter preview, Space select, q quit"

**And** when user presses ↑↓ arrow keys
**Then** cursor moves between files with visual highlighting

**And** when user presses Enter on ".." (parent directory)
**Then** TUI navigates to parent directory and refreshes file list

**And** when user presses Enter on a directory
**Then** TUI enters that directory and shows its contents

**Validation:**
- Start TUI in directory with 100+ files of mixed types
- Verify only preset files shown (.np3, .xmp, .lrtemplate)
- Navigate up/down with arrow keys
- Enter subdirectories and return to parent
- Check file count displayed correctly

---

### AC-2: Multi-File Selection with Space Bar

**Given** user is in file browser view
**When** user presses Space bar on a file
**Then** file is marked as selected with visual indicator (☑ checkbox or highlight)
**And** file is added to SelectedFiles list
**And** status line shows "X files selected"

**And** when user presses Space bar again on same file
**Then** selection is toggled off (☐ unchecked)
**And** file removed from SelectedFiles list
**And** status line updates count

**And** when user navigates to different files and selects multiple
**Then** all selections persist (can select 5, 10, 20+ files)
**And** status line shows accurate count

**Validation:**
- Select 10 files using space bar
- Verify all show selection indicator
- Deselect 3 files
- Verify count updates to 7 files selected
- Navigate away and back to verify selections persist

---

### AC-3: Live Parameter Preview Display

**Given** user is in file browser view
**When** user presses Enter on a preset file
**Then** TUI switches to preview view showing:
- **Left pane:** File details (name, size, format, path)
- **Right pane:** Extracted parameters in structured format:
  ```
  Basic Adjustments:
    Exposure: +0.5 EV
    Contrast: +15
    Highlights: -10
    Shadows: +5

  Color Adjustments:
    Saturation: -10
    Vibrance: +5

  Warnings:
    ⚠ Parameter 'Grain' not supported in NP3 format
  ```
- Help text: "↑↓ scroll, Tab select format, Enter convert, ESC back"

**And** when parameter list is longer than screen height
**Then** user can scroll with ↑↓ arrow keys
**And** scroll position indicator shows current position

**And** when preview fails (corrupted file)
**Then** error message displayed: "Failed to parse file: [error details]"
**And** user can press ESC to return to browser

**Validation:**
- Select XMP file with 30+ parameters
- Verify all parameters displayed correctly
- Scroll through full list
- Check warnings appear for unmappable parameters
- Test with corrupted file, verify error handling

---

### AC-4: Target Format Selection

**Given** user is in preview view
**When** user presses Tab key
**Then** target format cycles through options: NP3 → XMP → lrtemplate → NP3 (loop)
**And** selected format is highlighted/indicated clearly
**And** format selector shows all three options with one highlighted

**And** when source format is XMP
**Then** default selected target is NP3 (logical inverse)

**And** when user selects same format as source
**Then** warning shown: "⚠ Source and target are the same format (no conversion)"
**But** user can still proceed if desired (no hard block)

**Validation:**
- Open XMP file preview
- Verify NP3 pre-selected as target
- Press Tab 3 times, verify cycles through all formats
- Select XMP as target (same as source), verify warning appears

---

### AC-5: Batch Progress Visualization

**Given** user has selected 10 files for batch conversion
**And** user confirms conversion on confirmation screen
**When** batch conversion starts
**Then** progress view displays:
- Overall progress bar: [████░░░░] 40% (4/10)
- File-by-file status:
  ```
  ✓ file1.np3 (1.1 KB, 12ms)
  ✓ file2.np3 (1.3 KB, 14ms)
  ✓ file3.np3 (1.2 KB, 13ms)
  ✓ file4.np3 (1.0 KB, 11ms)
  ⏳ file5.xmp converting...
  ⏳ file6.xmp converting...
  ⏳ file7.xmp queued...
  ...
  ```
- Success/error counters: "Success: 4 ✓ | Errors: 0 ✗"
- Time info: "Elapsed: 0.5s | Est. remaining: 0.6s"
- Cancellation hint: "Press Ctrl+C to cancel"

**And** progress bar updates in real-time (every 100ms)
**And** file status updates immediately after each conversion completes
**And** color coding: ✓ green for success, ✗ red for errors, ⏳ blue for in-progress

**Validation:**
- Convert 20 files in batch
- Verify progress bar animates smoothly
- Check all files listed with status icons
- Confirm timing info updates
- Test with mixed success/error scenarios

---

### AC-6: Visual Confirmation Screen

**Given** user has selected target format in preview (single file) OR selected multiple files in browser (batch)
**When** user presses Enter to initiate conversion
**Then** confirmation screen displays:

**For single file:**
```
┌─────────────────────────────────────┐
│ Confirm Conversion                  │
├─────────────────────────────────────┤
│ File: portrait.xmp                  │
│ Convert: xmp → np3                  │
│ Output: portrait.np3                │
│ Location: /home/user/presets/       │
│                                      │
│ [Y] Confirm  [N] Cancel             │
└─────────────────────────────────────┘
```

**For batch:**
```
┌─────────────────────────────────────┐
│ Batch Conversion Confirmation       │
├─────────────────────────────────────┤
│ Files to convert: 10                │
│   • file1.xmp                       │
│   • file2.xmp                       │
│   • file3.xmp                       │
│   ...and 7 more                     │
│                                      │
│ Target format: np3                  │
│ Output directory: ./                │
│                                      │
│ [Y] Confirm  [N] Cancel  [E] Edit   │
└─────────────────────────────────────┘
```

**And** when user presses Y
**Then** conversion starts (transitions to progress view)

**And** when user presses N or ESC
**Then** returns to previous view without conversion (no side effects)

**And** when user presses E (batch only)
**Then** can edit target format or output directory before confirming

**Validation:**
- Test single file confirmation flow
- Test batch confirmation with 5+ files
- Verify all details shown correctly
- Test cancel (ESC) returns without conversion
- Test edit mode (E key) in batch confirmation

---

### AC-7: Keyboard-Driven Navigation

**Given** user is in any TUI view
**Then** ALL operations are accessible via keyboard:

| Key    | Action                                                 |
| ------ | ------------------------------------------------------ |
| ↑      | Move cursor up / Scroll up                             |
| ↓      | Move cursor down / Scroll down                         |
| Enter  | Select file / Confirm action / Open directory          |
| Space  | Toggle file selection (multi-select)                   |
| Tab    | Cycle target format                                    |
| ESC    | Go back to previous view / Cancel                      |
| q      | Quit TUI (with confirmation if conversion in progress) |
| Ctrl+C | Cancel operation / Quit                                |
| ?      | Show help screen                                       |
| y      | Confirm (on confirmation screen)                       |
| n      | Cancel (on confirmation screen)                        |
| e      | Edit (on batch confirmation screen)                    |

**And** no mouse interaction required or supported
**And** help text always visible showing available keys for current view
**And** context-sensitive shortcuts (e.g., 'e' only works on batch confirmation screen)

**Validation:**
- Navigate entire TUI workflow using only keyboard
- Verify all actions accessible
- Check help text updates per view
- Confirm mouse clicks have no effect

---

### AC-8: Terminal Resize Responsiveness

**Given** user has TUI running
**When** user resizes terminal window (drag window edge or fullscreen toggle)
**Then** TUI layout adjusts immediately:
- File browser and preview panes recalculate width (50/50 split maintained)
- Progress bar width adjusts to new terminal width
- Viewport heights adjust to new terminal height
- No content lost or corrupted
- No visual glitches (smooth reflow)

**And** when terminal is very small (<80x24)
**Then** TUI shows warning: "Terminal too small, minimum 80x24 recommended"
**But** TUI continues to function (degraded layout acceptable)

**Validation:**
- Start TUI at default size (80x24)
- Resize to large (200x60)
- Resize to small (70x20)
- Verify layout adjusts smoothly each time
- Check no crashes or visual corruption

---

### AC-9: Error Handling and Recovery

**Given** user attempts conversion that fails (corrupted file, invalid format, write permission error)
**When** error occurs
**Then** TUI displays error view:
```
┌─────────────────────────────────────┐
│ ✗ Conversion Error                  │
├─────────────────────────────────────┤
│ Failed to convert portrait.xmp      │
│                                      │
│ Error: Parse error: Invalid XML     │
│ Expected crs:xmpmeta root element   │
│                                      │
│ Suggestions:                        │
│ • Re-export preset from Lightroom   │
│ • Check file isn't corrupted        │
│ • Verify file has correct extension │
│                                      │
│ [R] Retry  [B] Back to browser      │
└─────────────────────────────────────┘
```

**And** error message is user-friendly (no stack traces in normal mode)
**And** actionable suggestions provided based on error type
**And** user can retry conversion (R key) or go back (B key)

**And** during batch conversion with errors
**Then** batch continues for remaining files (don't stop on first error)
**And** final summary shows: "Success: 8 ✓ | Errors: 2 ✗"
**And** error details available for each failed file

**Validation:**
- Test with corrupted file (truncated)
- Test with invalid format (wrong extension)
- Test with write-protected directory
- Verify error messages are clear
- Confirm batch continues after individual file errors

---

### AC-10: Format Filtering in File Browser

**Given** user is in directory with mixed file types (.jpg, .txt, .xmp, .np3, .lrtemplate, .pdf)
**When** file browser loads
**Then** only preset formats are shown (.np3, .xmp, .lrtemplate)
**And** other file types are hidden from list
**And** directories are always shown (for navigation)

**And** file list shows format badge for each file:
- **[NP3]** for .np3 files
- **[XMP]** for .xmp files
- **[LRT]** for .lrtemplate files

**Validation:**
- Create test directory with 50 mixed files (25 presets, 25 other types)
- Verify only 25 presets + directories shown
- Check format badges appear correctly
- Confirm other file types completely hidden (not just grayed out)

---

### AC-11: Unmappable Parameter Highlighting

**Given** user is viewing parameter preview for a file
**When** file contains parameters that cannot map to target format
**Then** those parameters are highlighted with warning icon (⚠) and yellow color
**And** warning tooltip explains: "This parameter is not supported in [target format]"

**Example:**
```
Color Adjustments:
  Saturation: -10
  Vibrance: +5
  ⚠ Grain: +15  ← Warning: Not supported in NP3 format
```

**And** when user hovers (or views) warning
**Then** explanation shown: "NP3 format does not support Grain parameter. This value will be omitted in conversion."

**And** warning count shown in preview summary: "Parameters: 42 | Warnings: 3"

**Validation:**
- Convert XMP (with Grain) to NP3
- Verify Grain parameter shows warning
- Check warning color/icon visible
- Confirm tooltip/explanation available
- Test multiple unmappable parameters (Dehaze, Texture, etc.)

---

### AC-12: Cancellable Operations (Ctrl+C)

**Given** user is performing batch conversion (progress view)
**When** user presses Ctrl+C
**Then** TUI cancels remaining conversions gracefully:
- Current in-progress conversions allowed to finish (atomic)
- No new conversions started
- Partial results saved (files completed before cancel)
- Cancel confirmation message: "✓ Cancelled. Converted 7/10 files before cancellation."

**And** no partial file writes (atomic conversions prevent corruption)
**And** user returned to file browser view
**And** can review which files converted successfully

**Validation:**
- Start batch of 20 files
- Press Ctrl+C after 5 files complete
- Verify no more conversions start
- Check 5 completed files exist and are valid
- Confirm no partial/corrupted files created

---

### AC-13: Cross-Platform Terminal Compatibility

**Given** Recipe TUI installed on Windows, macOS, and Linux
**When** user runs TUI on each platform
**Then** TUI renders correctly on all platforms:

**Windows:**
- Windows Terminal (ANSI colors work)
- PowerShell (basic colors, may degrade gracefully)
- CMD (minimal colors, text-only fallback)

**macOS:**
- Terminal.app (full color support)
- iTerm2 (full color support)

**Linux:**
- GNOME Terminal (full color support)
- Konsole (full color support)
- xterm (basic colors)
- Alacritty (full color support)

**And** TUI detects terminal capabilities and adjusts:
- Full ANSI color support → Vibrant colors
- Basic color support → Simplified palette
- No color support → Text-only with symbols (no colors)

**Validation:**
- Test TUI on Windows Terminal, PowerShell, CMD
- Test on macOS Terminal.app and iTerm2
- Test on Linux GNOME Terminal and xterm
- Verify no visual corruption on any platform
- Check color degradation graceful when unsupported

---

### AC-14: Performance - UI Responsiveness

**Given** user interacts with TUI
**Then** all operations meet performance targets:

| Operation                         | Target      | Measurement                 |
| --------------------------------- | ----------- | --------------------------- |
| Key press to screen update        | <16ms       | 60 FPS, no lag              |
| File browser render (1,000 files) | <50ms       | Instant display             |
| Parameter preview load            | <100ms      | Fast enough to feel instant |
| Screen transition                 | <30ms       | Smooth view changes         |
| Progress bar update               | Every 100ms | Smooth animation            |
| Terminal resize reflow            | <50ms       | No visual delay             |

**Validation:**
- Load directory with 1,000+ files, verify no lag when navigating
- Open 20 different file previews rapidly, check all <100ms
- Resize terminal rapidly, verify smooth reflow
- Monitor batch progress bar for smooth 10 FPS updates
- Use `time` command or profiler to verify keystroke latency <16ms

## Traceability Mapping

| Acceptance Criteria                     | Tech Spec Section                                                                                               | Component/Module                                                                           | Test Approach                                                                                                                                                   |
| --------------------------------------- | --------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **AC-1: File Browser Navigation**       | Services → cmd/tui/filebrowser.go<br>Data Models → FileItem struct<br>APIs → loadDirectory()                    | `cmd/tui/filebrowser.go`<br>`Bubbles list.Model`<br>`loadDirectory()` function             | Integration test: Load directory with 100+ mixed files, verify only presets shown. Navigate up/down, enter subdirectories, check cursor position updates.       |
| **AC-2: Multi-File Selection**          | Data Models → Model.SelectedFiles<br>Update Logic → Space key handling                                          | `cmd/tui/update.go`<br>`Model.SelectedFiles` array<br>Space key binding                    | Unit test: Simulate space key presses on FileItems, verify Selected flag toggles. Integration test: Select 10 files, verify count displayed correctly.          |
| **AC-3: Live Parameter Preview**        | Services → cmd/tui/preview.go<br>APIs → loadPreview(), parsePresetFile()<br>Workflows → Single File Conversion  | `cmd/tui/preview.go`<br>`viewport.Model` for scrolling<br>`loadPreview()` async command    | Integration test: Open XMP file with 30+ parameters, verify all displayed. Scroll through list, check warnings appear. Test corrupted file error handling.      |
| **AC-4: Target Format Selection**       | Data Models → Model.TargetFormat<br>Update Logic → Tab key cycling                                              | `cmd/tui/update.go`<br>Tab key binding<br>Format cycling logic                             | Unit test: Simulate Tab key presses, verify format cycles NP3→XMP→lrtemplate→NP3. Check default target is inverse of source.                                    |
| **AC-5: Batch Progress Visualization**  | Services → cmd/tui/progress.go<br>Data Models → ProgressBar, Results array<br>Workflows → Batch Conversion Flow | `cmd/tui/progress.go`<br>`Bubbles progress.Model`<br>`updateProgress()` function           | Integration test: Convert 20 files, verify progress bar updates every 100ms. Check color coding (green ✓, red ✗, blue ⏳). Measure time estimates accuracy.      |
| **AC-6: Visual Confirmation**           | Services → cmd/tui/confirm.go<br>Views → renderConfirmScreen()<br>Workflows → Confirmation Screen               | `cmd/tui/confirm.go`<br>`renderConfirmScreen()` function<br>View confirmation state        | Integration test: Trigger single file and batch confirmations, verify all details displayed. Test Y/N/E/ESC keys, confirm correct state transitions.            |
| **AC-7: Keyboard-Driven Navigation**    | Data Models → KeyMap struct<br>Update Logic → Key bindings<br>APIs → DefaultKeyMap()                            | `cmd/tui/keys.go`<br>`KeyMap` struct<br>All Update() key handlers                          | Manual test: Navigate entire TUI using only keyboard. Integration test: Simulate all key combinations, verify correct actions triggered.                        |
| **AC-8: Terminal Resize**               | Update Logic → tea.WindowSizeMsg handling<br>View Logic → Responsive layout calculation                         | `cmd/tui/update.go`<br>`cmd/tui/view.go`<br>Bubbletea window size handling                 | Integration test: Send WindowSizeMsg with different dimensions (80x24, 200x60, 70x20), verify layout recalculates. Check no visual corruption.                  |
| **AC-9: Error Handling**                | Data Models → Model.LastError<br>Views → Error screen<br>Workflows → Error Handling Workflow                    | `cmd/tui/view.go` (error view)<br>`ConversionError` type<br>Error recovery logic           | Negative test: Provide corrupted files, invalid formats, write-protected directories. Verify user-friendly error messages, retry/back options work.             |
| **AC-10: Format Filtering**             | APIs → loadDirectory(), detectFileFormat()<br>Data Models → FilteredFiles array                                 | `cmd/tui/filebrowser.go`<br>`detectFileFormat()` function<br>File filtering logic          | Unit test: Load directory with 50 mixed files, verify only .np3/.xmp/.lrtemplate shown. Check directories always visible. Verify format badges correct.         |
| **AC-11: Unmappable Parameters**        | Services → cmd/tui/preview.go<br>Data Models → PreviewContent.Warnings<br>Styling → warningStyle()              | `cmd/tui/preview.go`<br>`parsePresetFile()` warning detection<br>`Lipgloss warningStyle()` | Integration test: Load XMP with Grain parameter, convert to NP3, verify warning shown in yellow with ⚠ icon. Check tooltip/explanation available.               |
| **AC-12: Cancellable Operations**       | Update Logic → Ctrl+C handling<br>Workflows → Cancellation Workflow<br>Batch Processing → Cancellation flag     | `cmd/tui/update.go`<br>Worker pool cancellation<br>Atomic conversion guarantee             | Integration test: Start 20-file batch, press Ctrl+C after 5 complete. Verify remaining conversions cancelled, completed files valid, no partial writes.         |
| **AC-13: Cross-Platform Compatibility** | NFR → Terminal Compatibility<br>Bubbletea → Terminal capability detection                                       | Bubbletea runtime<br>Lipgloss color detection<br>Cross-platform builds                     | Cross-platform test: Run TUI on Windows Terminal/PowerShell/CMD, macOS Terminal/iTerm2, Linux GNOME/xterm. Verify rendering correct, colors degrade gracefully. |
| **AC-14: Performance Targets**          | NFR → Performance → UI Responsiveness<br>All performance optimizations                                          | All TUI components<br>Lazy loading<br>Async operations                                     | Benchmark test: Measure key press latency (<16ms), preview load (<100ms), directory scan (<50ms for 1,000 files). Profile with pprof if needed.                 |

### Traceability to PRD Requirements

| PRD Requirement                      | Epic 4 Implementation                                                                      | Verification                                                             |
| ------------------------------------ | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------ |
| **FR-4.1: Interactive File Browser** | cmd/tui/filebrowser.go with directory navigation, multi-select, format filtering           | AC-1, AC-2, AC-10 (keyboard navigation, space select, filtering)         |
| **FR-4.2: Live Parameter Preview**   | cmd/tui/preview.go with side-by-side layout, parameter extraction, unmappable highlighting | AC-3, AC-11 (preview display, warnings for unmappable params)            |
| **FR-4.3: Batch Progress**           | cmd/tui/progress.go with progress bars, file counts, color-coded status, estimated time    | AC-5 (real-time updates, color coding, cancellable)                      |
| **FR-4.4: Visual Validation**        | cmd/tui/confirm.go with confirmation screen showing files, format, output location         | AC-6 (user confirmation, can edit settings, cancel without side effects) |
| **NFR-1.3: TUI Performance**         | <16ms keystroke latency, <100ms preview load, <50ms directory scan                         | AC-14 (all performance targets)                                          |
| **NFR-2.1: Zero Data Exfiltration**  | All processing local, no network access in TUI code                                        | Code review (no net/http imports)                                        |
| **NFR-3.1: Conversion Accuracy**     | Uses shared converter.Convert() API (95%+ accuracy from Epic 1)                            | Round-trip tests in internal/converter                                   |
| **Cross-Platform TUI**               | Bubbletea framework with terminal capability detection, Windows/macOS/Linux support        | AC-13 (multiple terminal tests)                                          |

### Test Coverage Matrix

| Component                | Unit Tests                                                             | Integration Tests                                          | Manual Tests                    | Performance Tests           |
| ------------------------ | ---------------------------------------------------------------------- | ---------------------------------------------------------- | ------------------------------- | --------------------------- |
| `cmd/tui/filebrowser.go` | ✓ loadDirectory()<br>✓ detectFileFormat()<br>✓ File filtering          | ✓ Navigate directories<br>✓ Multi-select files             | ✓ Real directory navigation     | ✓ 1,000 file directory scan |
| `cmd/tui/preview.go`     | ✓ parsePresetFile()<br>✓ formatPreviewContent()<br>✓ Warning detection | ✓ Load various file types<br>✓ Scroll long parameter lists | ✓ Visual parameter display      | ✓ Preview load time <100ms  |
| `cmd/tui/progress.go`    | ✓ updateProgress()<br>✓ Progress calculation                           | ✓ Batch conversion tracking<br>✓ Real-time updates         | ✓ Visual progress bar animation | ✓ Update frequency 10 FPS   |
| `cmd/tui/confirm.go`     | ✓ renderConfirmScreen()                                                | ✓ Single/batch confirmation flows                          | ✓ Edit mode (batch)             | -                           |
| `cmd/tui/update.go`      | ✓ Key binding handlers<br>✓ State transitions                          | ✓ Full TUI workflows<br>✓ Error scenarios                  | ✓ Keyboard-only navigation      | ✓ Keystroke latency <16ms   |
| `cmd/tui/view.go`        | ✓ Layout calculations<br>✓ Styling functions                           | ✓ Terminal resize<br>✓ View transitions                    | ✓ Visual rendering quality      | ✓ Render time per frame     |
| **Shared Converter**     | (Epic 1 tests)                                                         | ✓ TUI calls converter.Convert()                            | -                               | ✓ Conversion time <20ms     |

### PRD Success Criteria Mapping

| PRD Success Criteria                       | Epic 4 Contribution                                                         | Measurement                                           |
| ------------------------------------------ | --------------------------------------------------------------------------- | ----------------------------------------------------- |
| **Intuitive UI for power users**           | Keyboard-driven TUI with visual feedback, parameter preview, batch progress | AC-7 (keyboard-only), AC-3 (preview), AC-5 (progress) |
| **95%+ accuracy**                          | Reuses Epic 1 converter.Convert() API with proven accuracy                  | Epic 1 round-trip tests                               |
| **Cross-platform (Windows, macOS, Linux)** | Bubbletea framework with terminal compatibility, color degradation          | AC-13 cross-platform tests                            |
| **Fast batch operations**                  | Parallel processing using worker pool pattern (same as CLI)                 | AC-5 performance (100 files <2s)                      |
| **Real-time feedback**                     | Live progress bars, parameter preview, color-coded status                   | AC-5 (progress visualization), AC-3 (preview)         |

## Risks, Assumptions, Open Questions

### Risks

**R-1: Charm v2 Experimental API Instability**
- **Risk:** Bubbletea/Lipgloss/Bubbles v2 experimental branches may have breaking changes before stable v2 release
- **Impact:** High - Code may break on dependency updates, requiring migration work
- **Likelihood:** Medium - v2 experimental is pre-release and subject to change
- **Mitigation:**
  - Pin exact v2-exp commit hashes in go.mod for reproducible builds
  - Monitor Charm repos for v2 API changes and breaking announcements
  - Budget time for migration when v2 reaches stable release
  - Use vendor directory (`go mod vendor`) to lock dependencies
  - Document all v2-specific API usage for easier migration tracking
  - Test thoroughly before updating any Charm dependencies

**R-2: Terminal Compatibility Variability**
- **Risk:** TUI may render incorrectly on obscure/legacy terminals (old xterm, tmux edge cases)
- **Impact:** Medium - Some users may see visual glitches or broken layouts
- **Likelihood:** Low - Bubbletea handles most terminal quirks automatically
- **Mitigation:**
  - Test on common terminals: Windows Terminal, PowerShell, iTerm2, GNOME Terminal, Alacritty
  - Document minimum terminal requirements: ANSI color support, 80x24 minimum size
  - Provide fallback text-only mode for terminals without color support
  - Include troubleshooting guide for known terminal issues

**R-3: Bubbletea Framework Learning Curve**
- **Risk:** Elm architecture pattern may be unfamiliar to developers (Model-Update-View paradigm)
- **Impact:** Medium - Slower development if team unfamiliar with pattern
- **Likelihood:** Medium - Elm architecture is different from imperative UI patterns
- **Mitigation:**
  - Study Bubbletea documentation and examples before implementation
  - Start with simple components (file browser) before complex ones (batch progress)
  - Use official Bubbles components where possible (less custom code)
  - Pair programming for initial components to share knowledge

**R-4: Complex State Management in Batch Operations**
- **Risk:** Managing parallel conversions with real-time UI updates could lead to race conditions
- **Impact:** High - Data corruption or TUI crashes if state management incorrect
- **Likelihood:** Low - Bubbletea message-passing architecture prevents races
- **Mitigation:**
  - Use Bubbletea's message-passing exclusively (no shared mutable state)
  - All async operations send messages back to Update() via channels
  - Leverage Elm architecture guarantee: Update() is single-threaded (no locks needed)
  - Comprehensive integration tests for batch scenarios

**R-5: Terminal Resize During Critical Operations**
- **Risk:** Resizing terminal during batch conversion could cause UI corruption or state loss
- **Impact:** Low - Annoying visual glitches but no data corruption
- **Likelihood:** Medium - Users may resize window unintentionally
- **Mitigation:**
  - Bubbletea handles window resize events gracefully (built-in feature)
  - Test resize during all TUI states (browser, preview, progress, confirm)
  - Ensure layout recalculation doesn't discard state (just re-renders)
  - Progress bar and file status persist through resizes

**R-6: User Expectation Mismatch (Mouse Support)**
- **Risk:** Users may expect mouse support (click files, drag selections) but TUI is keyboard-only
- **Impact:** Low - Slight friction for users unfamiliar with keyboard-driven interfaces
- **Likelihood:** Medium - Modern users accustomed to mouse/touch
- **Mitigation:**
  - Prominently display "Keyboard-only interface" in help text
  - Comprehensive keyboard shortcuts visible at all times
  - Consider adding mouse support in future version if users request it
  - Target audience (power users) likely comfortable with keyboard

### Assumptions

**A-1: Terminal Size Minimum**
- **Assumption:** Users run TUI in terminal at least 80x24 characters
- **Validation:** Display warning if terminal smaller, but allow degraded operation
- **Impact if False:** Layout may break with very small terminals (<70x20)

**A-2: ANSI Color Support**
- **Assumption:** Most target terminals support basic ANSI colors (16-color minimum)
- **Validation:** Bubbletea detects capabilities and degrades gracefully
- **Impact if False:** TUI falls back to monochrome (text-only) mode

**A-3: Conversion Engine Stability**
- **Assumption:** `internal/converter` API from Epic 1 is stable and accurate (95%+)
- **Validation:** Epic 1 round-trip tests confirm accuracy
- **Impact if False:** TUI would propagate converter bugs, but this is out of Epic 4 scope (Epic 1 responsibility)

**A-4: File System Performance**
- **Assumption:** Local file system has acceptable performance (not network drives with high latency)
- **Validation:** Preview loading may be slow on network drives but still functional
- **Impact if False:** Show "Loading..." spinner for longer, but no functional breakage

**A-5: Single-User Workflow**
- **Assumption:** TUI is for single-user, single-session workflows (no multi-user coordination)
- **Validation:** Design is stateless, no persistent state or collaboration features
- **Impact if False:** Would need to add file locking, session management (out of scope)

**A-6: Bubbletea v2 Framework Maturity**
- **Assumption:** Bubbletea v2 experimental branch is stable enough for development
- **Validation:** v2 experimental is actively maintained with community testing, though pre-release
- **Impact if False:** Breaking changes before v2 stable release would require code updates. Migration from v2-exp to v2 stable expected to be minimal.
- **⚠️ Risk:** v2 experimental versions may have breaking changes. Code will need updates when v2 reaches stable release.

### Open Questions

**Q-1: Should TUI support configuration persistence?**
- **Context:** Remember last used directory, default target format, window size preferences
- **Trade-offs:**
  - Pro: Better UX, fewer repeated selections
  - Con: Adds complexity (config file management), violates stateless design
- **Decision:** **Defer to post-MVP**. Stateless design simpler for initial release. Add config file (`~/.recipe-tui.yaml`) if users request it.

**Q-2: Should TUI support custom color schemes?**
- **Context:** Allow users to customize colors (Solarized, Dracula, Gruvbox themes)
- **Trade-offs:**
  - Pro: Personalization, accessibility (high contrast for vision impairment)
  - Con: Adds complexity, Lipgloss defaults work well for most users
- **Decision:** **Defer**. Use Lipgloss defaults for MVP. Add theming support via config file if demanded.

**Q-3: Should TUI show real-time file watcher?**
- **Context:** Automatically detect new preset files added to directory and refresh list
- **Trade-offs:**
  - Pro: Convenient for workflows where files added externally
  - Con: Adds complexity (fsnotify library), potential performance impact
- **Decision:** **Out of scope for MVP**. User can manually refresh (future: press 'r' key to reload directory).

**Q-4: Should TUI support undo for batch operations?**
- **Context:** Ability to undo batch conversion (delete generated files, restore state)
- **Trade-offs:**
  - Pro: Safety net for accidental conversions
  - Con: Requires state tracking, file deletion logic (risky)
- **Decision:** **Out of scope**. Conversions are intentional (confirmation screen required). User can manually delete unwanted files.

**Q-5: Should parameter preview show visual comparison?**
- **Context:** Side-by-side diff showing parameter changes (before/after conversion simulation)
- **Trade-offs:**
  - Pro: More insight into conversion accuracy
  - Con: Requires rendering two previews simultaneously, complex UI
- **Decision:** **Defer to Epic 5** (Data Extraction & Inspection - diff tool). Current preview shows source parameters, which is sufficient for MVP.

**Q-6: Should TUI support batch rename patterns?**
- **Context:** User specifies output filename pattern (e.g., `preset_{index}.np3`, `{name}_converted.xmp`)
- **Trade-offs:**
  - Pro: More control over output naming
  - Con: Adds complexity to confirmation screen, potential for user error
- **Decision:** **Defer**. Current behavior (preserve name, change extension) is simple and predictable. Add rename patterns if users request.

**Q-7: Should TUI support subdirectory recursion in batch?**
- **Context:** Select all preset files in current directory AND all subdirectories
- **Example:**
  - Current: Only files in `presets/` selected
  - Recursive: Files in `presets/`, `presets/portraits/`, `presets/landscapes/`, etc.
- **Trade-offs:**
  - Pro: Powerful for large preset libraries
  - Con: Potentially dangerous (accidentally convert entire drive), complex file browser UI
- **Decision:** **Defer**. File browser shows subdirectories, user can navigate into each. Add recursive batch mode with safeguards (max depth, explicit confirmation) if requested.

**Q-8: How should TUI handle very long file names?**
- **Context:** File name longer than terminal width (e.g., `my_extremely_long_preset_name_for_landscape_photography_golden_hour.xmp`)
- **Options:**
  - Truncate: `my_extremely_long_preset_na...xmp`
  - Ellipsis: `my_extremely_...golden_hour.xmp` (middle truncation)
  - Wrap: Multi-line display
- **Decision:** **Truncate with ellipsis** (Bubbles list default behavior). Middle truncation preferred (`...`) to show file extension.

**Q-9: Should TUI support conversion history log?**
- **Context:** Persistent log of all conversions (date, files, success/errors) saved to `~/.recipe-tui-history.log`
- **Trade-offs:**
  - Pro: Auditability, debugging, undo reference
  - Con: Persistent state (violates stateless design), privacy concern
- **Decision:** **Out of scope for MVP**. User can enable verbose logging with `--debug` flag if needed (session-only). Consider persistent history in future if demanded.

### Risk Prioritization Summary

| Risk ID                         | Impact | Likelihood | Priority   | Mitigation Status                                                     |
| ------------------------------- | ------ | ---------- | ---------- | --------------------------------------------------------------------- |
| R-1: Terminal Compatibility     | Medium | Low        | **Medium** | Test common terminals, document requirements, fallback mode planned   |
| R-2: Bubbletea Learning Curve   | Medium | Medium     | **High**   | Study docs before dev, start simple, use official components          |
| R-3: Batch State Management     | High   | Low        | **Medium** | Bubbletea message-passing prevents races, comprehensive tests planned |
| R-4: Terminal Resize During Ops | Low    | Medium     | **Low**    | Bubbletea handles gracefully, test all resize scenarios               |
| R-5: Mouse Support Expectation  | Low    | Medium     | **Low**    | Prominent keyboard-only messaging, target audience comfortable        |

### Next Steps for Risk Mitigation

1. **Before Development:**
   - Complete Bubbletea tutorial and build simple prototype (file browser only)
   - Set up terminal testing matrix (Windows/macOS/Linux, various emulators)
   - Review Elm architecture pattern with team (Model-Update-View paradigm)

2. **During Development:**
   - Implement file browser first (simpler), validate pattern before complex components
   - Test terminal resize on every component as developed
   - Use Bubbles components where possible (less custom code = fewer bugs)
   - Comprehensive logging in debug mode for troubleshooting

3. **Before Release:**
   - Manual testing on all target terminals (see AC-13)
   - Stress test batch operations (1,000 files, terminal resize during batch)
   - User testing with target audience (power users) for feedback
   - Document all assumptions and limitations in README

## Test Strategy Summary

### Testing Philosophy

**Comprehensive TUI Validation:** TUI is the visual interface for terminal users - correctness and responsiveness are critical. Test strategy emphasizes real-world terminal compatibility, UI state consistency, and performance under batch operations.

**Shared Converter Trust:** TUI delegates all conversion logic to `internal/converter` (Epic 1). Focus TUI tests on UI state management, keyboard handling, and visual rendering. Conversion accuracy is validated by Epic 1 tests.

**Elm Architecture Benefits:** Bubbletea's Model-Update-View pattern enables deterministic testing. State transitions can be tested in isolation by calling `Update()` directly with messages.

**⚠️ v2 Experimental Testing:** All tests use Bubbletea/Lipgloss/Bubbles v2 experimental branches. API may change before stable v2 release. Pin exact commit hashes in go.mod and update tests if breaking changes occur upstream.

### Test Levels

#### 1. Unit Tests

**Scope:** Individual functions and state transitions in isolation

**Components Under Test:**

**A. State Transition Logic (`cmd/tui/update.go`)**
```go
func TestUpdate_FileSelection(t *testing.T) {
    // Initial state
    model := Model{
        FileBrowser: list.New([]list.Item{
            FileItem{Name: "file1.xmp", Selected: false},
            FileItem{Name: "file2.xmp", Selected: false},
        }, nil, 0, 0),
        SelectedFiles: []FileItem{},
    }

    // Simulate space key press (select file)
    msg := tea.KeyMsg{Type: tea.KeySpace}
    newModel, _ := model.Update(msg)

    // Assertions
    assert.True(t, newModel.FileBrowser.Items()[0].(FileItem).Selected)
    assert.Len(t, newModel.SelectedFiles, 1)
    assert.Equal(t, "file1.xmp", newModel.SelectedFiles[0].Name)
}

func TestUpdate_FormatCycling(t *testing.T) {
    model := Model{TargetFormat: "np3", CurrentView: ViewPreview}

    // Press Tab 3 times
    for i := 0; i < 3; i++ {
        msg := tea.KeyMsg{Type: tea.KeyTab}
        model, _ = model.Update(msg)
    }

    // Should cycle back to np3 (np3 → xmp → lrtemplate → np3)
    assert.Equal(t, "np3", model.TargetFormat)
}
```

**B. File System Operations (`cmd/tui/filebrowser.go`)**
```go
func TestLoadDirectory(t *testing.T) {
    // Setup: Create temp directory with mixed files
    tmpDir := t.TempDir()
    createFile(tmpDir, "preset1.xmp")
    createFile(tmpDir, "preset2.np3")
    createFile(tmpDir, "image.jpg")  // Should be filtered out
    createFile(tmpDir, "document.pdf")  // Should be filtered out

    // Load directory
    files, err := loadDirectory(tmpDir)

    // Assertions
    assert.NoError(t, err)
    assert.Len(t, files, 2, "Should only load preset files")
    assert.Contains(t, fileNames(files), "preset1.xmp")
    assert.Contains(t, fileNames(files), "preset2.np3")
}

func TestDetectFileFormat(t *testing.T) {
    tests := []struct {
        filename string
        want     string
    }{
        {"preset.np3", "np3"},
        {"preset.xmp", "xmp"},
        {"preset.lrtemplate", "lrtemplate"},
        {"unknown.txt", ""},
    }

    for _, tt := range tests {
        t.Run(tt.filename, func(t *testing.T) {
            got, _ := detectFileFormat(tt.filename)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**C. Preview Parsing (`cmd/tui/preview.go`)**
```go
func TestParsePresetFile_XMP(t *testing.T) {
    data, _ := os.ReadFile("testdata/xmp/portrait.xmp")

    content, err := parsePresetFile(data, "xmp")

    assert.NoError(t, err)
    assert.Equal(t, "xmp", content.DetectedFormat)
    assert.NotEmpty(t, content.Parameters)
    assert.Contains(t, content.Parameters, "Exposure")
    assert.Contains(t, content.Parameters, "Contrast")
}

func TestFormatPreviewContent_WithWarnings(t *testing.T) {
    content := &PreviewContent{
        FileName: "test.xmp",
        Parameters: map[string]interface{}{
            "Exposure": 0.5,
            "Grain": 15,  // Unmappable to NP3
        },
        Warnings: []string{"Parameter 'Grain' not supported in NP3 format"},
    }

    output := formatPreviewContent(content)

    assert.Contains(t, output, "Exposure: 0.5")
    assert.Contains(t, output, "⚠ Grain: 15")
    assert.Contains(t, output, "not supported in NP3 format")
}
```

**D. Rendering Logic (`cmd/tui/view.go`)**
```go
func TestView_BrowseScreen(t *testing.T) {
    model := Model{
        CurrentView: ViewBrowse,
        FileBrowser: /* ... populated with files ... */,
        WindowWidth: 80,
        WindowHeight: 24,
    }

    output := model.View()

    assert.Contains(t, output, "↑↓ navigate")  // Help text
    assert.Contains(t, output, "Space select")
    assert.Contains(t, output, "[NP3]")  // Format badge
    assert.Greater(t, len(output), 100, "Should render substantial content")
}
```

**Tooling:**
- Go standard `testing` package
- `testify/assert` for readable assertions (optional, can use stdlib)
- Table-driven tests for multiple scenarios
- Temp directories for file system tests

---

#### 2. Integration Tests

**Scope:** Full TUI workflows with real file I/O and state transitions

**Test Scenarios:**

**A. Single File Conversion Workflow (AC-1, AC-3, AC-4, AC-6)**
```go
func TestTUI_SingleFileConversion(t *testing.T) {
    tmpDir := t.TempDir()
    inputFile := filepath.Join(tmpDir, "portrait.xmp")
    copyFile("testdata/xmp/portrait.xmp", inputFile)

    // Initialize model
    model := New()
    model, _ = model.Init()()  // Run Init command

    // 1. Load directory
    msg := dirChangedMsg{path: tmpDir}
    model, _ = model.Update(msg)
    assert.Equal(t, ViewBrowse, model.CurrentView)
    assert.Greater(t, len(model.FileBrowser.Items()), 0)

    // 2. Select file (Enter)
    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, cmd := model.Update(msg)
    // Execute loadPreview command (async)
    previewMsg := cmd()
    model, _ = model.Update(previewMsg)
    assert.Equal(t, ViewPreview, model.CurrentView)
    assert.NotNil(t, model.PreviewData)

    // 3. Select target format (Tab to np3)
    model.TargetFormat = "xmp"  // Default
    msg = tea.KeyMsg{Type: tea.KeyTab}
    model, _ = model.Update(msg)
    assert.Equal(t, "np3", model.TargetFormat)

    // 4. Confirm conversion (Enter)
    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, _ = model.Update(msg)
    assert.Equal(t, ViewConfirm, model.CurrentView)

    // 5. Confirm (Y key)
    msg = tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'y'}}
    model, cmd = model.Update(msg)
    assert.Equal(t, ViewProgress, model.CurrentView)

    // Execute conversion command
    convMsg := cmd()
    model, _ = model.Update(convMsg)

    // 6. Verify completion
    assert.Equal(t, ViewComplete, model.CurrentView)
    assert.FileExists(t, filepath.Join(tmpDir, "portrait.np3"))
}
```

**B. Batch Conversion Workflow (AC-2, AC-5, AC-6)**
```go
func TestTUI_BatchConversion(t *testing.T) {
    tmpDir := t.TempDir()

    // Setup: Copy 10 sample files
    for i := 1; i <= 10; i++ {
        copyFile("testdata/xmp/sample.xmp", filepath.Join(tmpDir, fmt.Sprintf("file%d.xmp", i)))
    }

    model := New()
    model, _ = model.Init()()

    // Load directory
    msg := dirChangedMsg{path: tmpDir}
    model, _ = model.Update(msg)

    // Select 10 files (Space bar on each)
    for i := 0; i < 10; i++ {
        msg = tea.KeyMsg{Type: tea.KeySpace}
        model, _ = model.Update(msg)
        msg = tea.KeyMsg{Type: tea.KeyDown}
        model, _ = model.Update(msg)
    }
    assert.Len(t, model.SelectedFiles, 10)

    // Trigger batch conversion ('c' key shortcut)
    msg = tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'c'}}
    model, _ = model.Update(msg)
    assert.Equal(t, ViewConfirm, model.CurrentView)

    // Confirm batch (Y)
    msg = tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'y'}}
    model, cmd := model.Update(msg)
    assert.Equal(t, ViewProgress, model.CurrentView)

    // Execute batch conversion (async)
    startTime := time.Now()
    convMsg := cmd()
    model, _ = model.Update(convMsg)
    elapsed := time.Since(startTime)

    // Verify results
    assert.Equal(t, ViewComplete, model.CurrentView)
    assert.Equal(t, 10, model.SuccessCount)
    assert.Equal(t, 0, model.ErrorCount)
    assert.Less(t, elapsed, 2*time.Second, "Batch should complete in <2s")

    // Verify all output files created
    files, _ := filepath.Glob(filepath.Join(tmpDir, "*.np3"))
    assert.Len(t, files, 10)
}
```

**C. Error Handling (AC-9)**
```go
func TestTUI_CorruptedFileError(t *testing.T) {
    tmpDir := t.TempDir()
    corruptedFile := filepath.Join(tmpDir, "corrupted.xmp")
    os.WriteFile(corruptedFile, []byte("INVALID XML"), 0644)

    model := New()
    model, _ = model.Init()()

    // Load directory and select corrupted file
    msg := dirChangedMsg{path: tmpDir}
    model, _ = model.Update(msg)

    msg = tea.KeyMsg{Type: tea.KeyEnter}
    model, cmd := model.Update(msg)

    // Execute preview load (should fail)
    previewMsg := cmd()
    model, _ = model.Update(previewMsg)

    // Verify error view shown
    assert.Equal(t, ViewError, model.CurrentView)
    assert.NotNil(t, model.LastError)
    assert.Contains(t, model.LastError.Error(), "parse")

    // Test recovery (B key to go back)
    msg = tea.KeyMsg{Type: tea.KeyRune, Runes: []rune{'b'}}
    model, _ = model.Update(msg)
    assert.Equal(t, ViewBrowse, model.CurrentView)
    assert.Nil(t, model.LastError)
}
```

**D. Terminal Resize (AC-8)**
```go
func TestTUI_TerminalResize(t *testing.T) {
    model := New()
    model.WindowWidth = 80
    model.WindowHeight = 24

    // Simulate resize to large
    msg := tea.WindowSizeMsg{Width: 200, Height: 60}
    model, _ = model.Update(msg)
    assert.Equal(t, 200, model.WindowWidth)
    assert.Equal(t, 60, model.WindowHeight)

    // Verify View() can render (no panic)
    output := model.View()
    assert.NotEmpty(t, output)

    // Simulate resize to small
    msg = tea.WindowSizeMsg{Width: 70, Height: 20}
    model, _ = model.Update(msg)
    assert.Equal(t, 70, model.WindowWidth)

    // Should still render (degraded layout acceptable)
    output = model.View()
    assert.NotEmpty(t, output)
}
```

**Tooling:**
- `testify/assert` for assertions
- Temp directories (`t.TempDir()`) for isolated file systems
- Bubbletea testing pattern: Call Update() directly with messages
- Time measurement for performance validation

---

#### 3. Manual Testing (Critical for TUI)

**Scope:** Visual validation and user experience testing on real terminals

**Test Matrix:**

| Platform          | Terminal         | Color Support     | Resolution    | Tested By | Status |
| ----------------- | ---------------- | ----------------- | ------------- | --------- | ------ |
| **Windows 10/11** | Windows Terminal | Full (24-bit)     | 80x24, 120x40 | Developer | ✓      |
|                   | PowerShell       | Basic (16-color)  | 80x24         | Developer | ✓      |
|                   | CMD              | Minimal (8-color) | 80x24         | Developer | ✓      |
| **macOS**         | Terminal.app     | Full (256-color)  | 80x24, 200x60 | Developer | ✓      |
|                   | iTerm2           | Full (24-bit)     | Various       | Developer | ✓      |
| **Linux**         | GNOME Terminal   | Full (24-bit)     | 80x24, 120x40 | Developer | ✓      |
|                   | Konsole (KDE)    | Full (24-bit)     | 80x24         | Tester    | ✓      |
|                   | xterm            | Basic (256-color) | 80x24         | Tester    | ✓      |
|                   | Alacritty        | Full (24-bit)     | Various       | Developer | ✓      |

**Visual Validation Checklist:**

1. **File Browser:**
   - [ ] File list renders correctly (no truncation, aligned columns)
   - [ ] Format badges visible ([NP3], [XMP], [LRT])
   - [ ] Cursor highlight visible and moves smoothly
   - [ ] Selection checkboxes (☑) appear correctly
   - [ ] Help text at bottom is readable
   - [ ] Directory icons/indicators clear

2. **Parameter Preview:**
   - [ ] Side-by-side layout (file details | parameters)
   - [ ] Parameter values aligned correctly
   - [ ] Warning icons (⚠) visible in yellow/orange
   - [ ] Scrollable with ↑↓ arrows (viewport works)
   - [ ] Format selector highlights current choice

3. **Batch Progress:**
   - [ ] Progress bar animates smoothly (10 FPS)
   - [ ] Color coding: ✓ green, ✗ red, ⏳ blue
   - [ ] File list scrollable if >10 files
   - [ ] Time estimates update (elapsed, remaining)
   - [ ] Percentage updates accurately

4. **Confirmation Screen:**
   - [ ] Box drawing characters render correctly (borders)
   - [ ] All file details visible
   - [ ] Y/N/E options highlighted clearly
   - [ ] Batch file list truncated correctly (...and X more)

5. **Error View:**
   - [ ] Red ✗ icon visible
   - [ ] Error message readable (word wrap if long)
   - [ ] Suggestions listed clearly
   - [ ] R/B options visible

**User Experience Validation:**

- **Workflow Timing:** Time full conversion workflow (browse → preview → confirm → convert)
  - Target: <30 seconds for experienced user, <60 seconds for new user
- **Keyboard Efficiency:** Count keystrokes for common tasks
  - Single file: ~5-7 keystrokes (navigate, enter, tab, enter, y)
  - Batch 10 files: ~25-30 keystrokes (10x space, c, enter, y)
- **Error Recovery:** Verify user can recover from all error states without restart
- **Help Discoverability:** New user should discover '?' for help within 30 seconds

**Terminal-Specific Issues to Watch:**

- **Windows CMD:** Limited color support (verify fallback to basic colors)
- **tmux/screen:** Nested terminal sessions (verify resize events propagate)
- **SSH sessions:** Remote terminal (verify no rendering lag)
- **Dark vs Light themes:** Test on both terminal background colors

---

#### 4. Performance Benchmarks (AC-14)

**Scope:** Quantitative validation of performance targets

**Benchmark Tests:**

```go
func BenchmarkTUI_DirectoryScan_1000Files(b *testing.B) {
    tmpDir := setupDirWith1000Files(b)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        files, _ := loadDirectory(tmpDir)
        if len(files) != 1000 {
            b.Fatalf("Expected 1000 files, got %d", len(files))
        }
    }
}

func BenchmarkTUI_PreviewLoad(b *testing.B) {
    data, _ := os.ReadFile("testdata/xmp/portrait.xmp")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        content, err := parsePresetFile(data, "xmp")
        if err != nil {
            b.Fatal(err)
        }
        _ = formatPreviewContent(content)
    }
}

func BenchmarkTUI_StateUpdate_KeyPress(b *testing.B) {
    model := New()
    model.FileBrowser = /* ... populated with 100 files ... */

    msg := tea.KeyMsg{Type: tea.KeyDown}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        model, _ = model.Update(msg)
    }
}
```

**Performance Targets:**

| Benchmark                    | Target            | Command                                     |
| ---------------------------- | ----------------- | ------------------------------------------- |
| Directory scan (1,000 files) | <50ms/op          | `go test -bench=BenchmarkTUI_DirectoryScan` |
| Preview load + format        | <100ms/op         | `go test -bench=BenchmarkTUI_PreviewLoad`   |
| Keystroke → Update           | <16ms/op (60 FPS) | `go test -bench=BenchmarkTUI_StateUpdate`   |
| View rendering               | <16ms/op (60 FPS) | `go test -bench=BenchmarkTUI_View`          |

**Profiling for Optimization:**

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace execution
go test -trace=trace.out -bench=.
go tool trace trace.out
```

---

### Test Coverage Goals

| Component                | Unit Test Coverage           | Integration Test Coverage   | Total Goal |
| ------------------------ | ---------------------------- | --------------------------- | ---------- |
| `cmd/tui/update.go`      | 90%+ (all state transitions) | 95%+ (full workflows)       | **95%+**   |
| `cmd/tui/view.go`        | 80%+ (rendering functions)   | 90%+ (visual validation)    | **85%+**   |
| `cmd/tui/filebrowser.go` | 95%+ (file I/O logic)        | 100% (directory navigation) | **95%+**   |
| `cmd/tui/preview.go`     | 90%+ (parsing, formatting)   | 95%+ (preview workflows)    | **90%+**   |
| `cmd/tui/progress.go`    | 85%+ (progress calculation)  | 100% (batch scenarios)      | **90%+**   |
| `cmd/tui/confirm.go`     | 90%+ (screen rendering)      | 95%+ (confirmation flows)   | **90%+**   |
| **Overall TUI**          | **90%+**                     | **95%+**                    | **90%+**   |

**Note:** Shared converter (`internal/converter`) is validated by Epic 1 tests (not counted in Epic 4 coverage).

---

### Test Automation Strategy

**Makefile Targets:**

```makefile
.PHONY: test test-unit test-integration test-bench test-coverage test-manual

# Run all automated tests
test:
	go test -v ./cmd/tui/...

# Run only unit tests (fast feedback)
test-unit:
	go test -v -short ./cmd/tui/...

# Run integration tests (slower)
test-integration:
	go test -v -run Integration ./cmd/tui/...

# Run performance benchmarks
test-bench:
	go test -bench=. -benchmem ./cmd/tui/...

# Generate coverage report
test-coverage:
	go test -coverprofile=coverage.out ./cmd/tui/...
	go tool cover -html=coverage.out -o coverage.html

# Manual testing checklist
test-manual:
	@echo "Manual Testing Checklist:"
	@echo "1. Test on Windows Terminal (full colors)"
	@echo "2. Test on PowerShell (basic colors)"
	@echo "3. Test on macOS Terminal.app"
	@echo "4. Test on Linux GNOME Terminal"
	@echo "5. Verify resize behavior (80x24 → 200x60 → 70x20)"
	@echo "6. Test batch conversion with 100 files"
	@echo "7. Verify error handling (corrupted file)"
	@echo "8. Check keyboard-only navigation (no mouse)"
	@echo "See docs/test-checklist.md for full details"
```

---

### CI/CD Integration

**GitHub Actions Workflow:**

```yaml
name: TUI Tests

on: [push, pull_request]

jobs:
  test-unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run Unit Tests
        run: make test-unit

      - name: Run Integration Tests
        run: make test-integration

      - name: Generate Coverage
        run: make test-coverage

      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  test-bench:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run Benchmarks
        run: make test-bench

      - name: Check Performance Targets
        run: |
          # Parse benchmark output and fail if any target exceeded
          go test -bench=. -benchmem ./cmd/tui/ > bench.txt
          # (Add script to validate targets: <50ms directory scan, <100ms preview, <16ms keystroke)
```

---

### Success Criteria for Test Strategy

✅ **90%+ code coverage** for TUI codebase (`cmd/tui/`)
✅ **100% of acceptance criteria** have corresponding automated tests
✅ **Cross-platform manual tests** completed on Windows/macOS/Linux terminals
✅ **Performance targets met** (<16ms keystroke, <50ms directory scan, <100ms preview)
✅ **Zero regressions** in existing Epic 1 conversion accuracy
✅ **User experience validation** by target audience (power users)
✅ **Manual test checklist** completed before release (AC-13 terminal matrix)
