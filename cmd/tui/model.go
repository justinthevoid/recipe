package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

// model holds the TUI application state
type model struct {
	currentDir  string
	files       []FileInfo
	cursor      int
	selected    map[string]bool // Map of file paths to selection state
	termWidth   int
	termHeight  int
	showHelp    bool
	err         error

	// Preview pane fields (Story 4-2)
	showPreview     bool   // True when terminal width >= 120 columns
	previewFile     string // Currently previewed file path
	previewContent  string // Rendered parameter display text
	previewLoading  bool   // Loading indicator state
	scrollOffset    int    // Current scroll position in preview pane
	viewportHeight  int    // Visible lines in preview pane
	previewCache    map[string]string // Cache of formatted preview content by file path
	previewFocused  bool   // True when preview pane has focus (for scrolling)

	// Batch conversion fields (Story 4-3)
	showFormatPrompt bool              // Show format selection menu
	showOutputPrompt bool              // Show output directory prompt
	showConfirmation bool              // Show conversion confirmation screen
	converting       bool              // Batch conversion in progress
	currentFile      int               // Current file index (0-based)
	totalFiles       int               // Total files in batch
	completedFiles   int               // Successfully completed count
	errorCount       int               // Failed conversions count
	warningCount     int               // Conversions with warnings
	startTime        time.Time         // Batch start timestamp
	elapsedTime      time.Duration     // Updated every second
	estimatedRemaining time.Duration   // Calculated from average
	currentFileName  string            // File being processed
	currentStatus    string            // "converting", "success", "warning", "error"
	completedList    []string          // Last 5 completed filenames
	results          []ConversionResult // Full results list
	cancelling       bool              // Cancel requested
	cancelChan       chan bool         // Signal to stop batch
	showSummary      bool              // Display summary screen
	targetFormat     string            // Selected output format ("np3", "xmp", "lrtemplate")
	outputDir        string            // Output directory path
	formatMenuCursor int               // Current selection in format menu (0-2)

	// Validation screen fields (Story 4-4)
	showValidation       bool              // Show validation screen
	validationPassed     bool              // All validation checks passed
	validationWarnings   []Warning         // Detected unmappable parameter warnings
	validationFiles      []ValidationFile  // Enriched file info with warnings
	validationPlan       ConversionPlan    // Batch statistics and estimates
	showSettingsEditor   bool              // Settings editor active
	editorCursor         int               // Current field in editor (0=format, 1=dir, 2=files)
	editedTargetFormat   string            // Temporary format during editing
	editedOutputDir      string            // Temporary directory during editing
	editedFileSelection  map[string]bool   // Temporary file selection during editing
	fileListScrollOffset int               // Scroll position in validation file list
	showDirectoryPrompt  bool              // Show directory creation/overwrite prompt
	directoryIssue       string            // "missing", "overwrite", or ""
	overwriteFiles       []OverwriteInfo   // Files that will be overwritten
}

// FileInfo represents a file or directory entry
type FileInfo struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
	IsDir   bool
	Format  string // "np3", "xmp", "lrtemplate", "dir"
}

// ConversionResult represents the result of a single file conversion
type ConversionResult struct {
	File         string // Filename
	Status       string // "success", "warning", "error", "cancelled"
	Message      string // Error/warning message if applicable
	SourceFormat string // Original format
	TargetFormat string // Converted format
}

// ValidationFile represents a file with enriched validation information (Story 4-4)
type ValidationFile struct {
	Name         string   // Filename
	Path         string   // Full path
	Size         int64    // File size in bytes
	SourceFormat string   // "np3", "xmp", "lrtemplate"
	TargetFormat string   // Target format
	HasWarnings  bool     // True if unmappable parameters detected
	Warnings     []string // Specific warning messages
}

// Warning represents an unmappable parameter warning (Story 4-4)
type Warning struct {
	File           string   // Filename
	ParameterCount int      // Number of unmappable parameters
	Parameters     []string // Specific parameter names
	Severity       string   // "minor" (<3 params) or "significant" (≥3 params)
	Description    string   // Human-readable explanation
}

// ConversionPlan represents batch conversion statistics and estimates (Story 4-4)
type ConversionPlan struct {
	FileCount           int           // Total files
	TotalInputSize      int64         // Sum of input file sizes
	EstimatedOutputSize int64         // Estimated output size
	EstimatedTime       time.Duration // Estimated conversion time
	AvailableDiskSpace  int64         // Available disk space in bytes
	CrossFormatCount    int           // Cross-format conversions
	SameFormatCount     int           // Same-format updates
}

// OverwriteInfo represents a file that will be overwritten (Story 4-4)
type OverwriteInfo struct {
	File         string // Filename
	ExistingSize int64  // Current file size
	NewSize      int64  // Estimated new file size
}

// initialModel creates the initial model state
func initialModel() model {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	return model{
		currentDir:       cwd,
		selected:         make(map[string]bool),
		termWidth:        80,  // Default, will be updated on first render
		termHeight:       24,  // Default, will be updated on first render
		showHelp:         false,
		err:              nil,
		showPreview:      false,  // Will be set based on terminal width
		previewCache:     make(map[string]string),
		previewFocused:   false,
		// Batch conversion fields initialized
		showFormatPrompt: false,
		showOutputPrompt: false,
		showConfirmation: false,
		converting:       false,
		cancelChan:       make(chan bool, 1),
		formatMenuCursor: 0,
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return loadFilesCmd(m.currentDir)
}

// WindowSizeMsgInterface is an interface for window size messages
type WindowSizeMsgInterface interface {
	GetWidth() int
	GetHeight() int
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		// Update preview visibility based on terminal width
		m.showPreview = m.termWidth >= 120
		// Calculate viewport height for preview pane (reserve space for header, footer)
		if m.showPreview {
			m.viewportHeight = m.termHeight - 8 // Reserve for headers, borders, status
		}
		return m, nil

	case WindowSizeMsgInterface:
		// For testing
		m.termWidth = msg.GetWidth()
		m.termHeight = msg.GetHeight()
		// Update preview visibility based on terminal width
		m.showPreview = m.termWidth >= 120
		// Calculate viewport height for preview pane
		if m.showPreview {
			m.viewportHeight = m.termHeight - 8
		}
		return m, nil

	case filesLoadedMsg:
		m.files = msg.files
		m.err = msg.err
		// Reset cursor if it's out of bounds
		if m.cursor >= len(m.files) {
			m.cursor = 0
		}
		return m, nil

	case previewLoadedMsg:
		// Handle async preview loading results
		if msg.err != nil {
			// Show error in preview pane
			m.previewContent = fmt.Sprintf("  Error loading preview:\n  %v", msg.err)
		} else {
			// Cache the formatted content
			m.previewCache[msg.filePath] = msg.content
			// Update display if this is still the current file
			if msg.filePath == m.previewFile {
				m.previewContent = msg.content
			}
		}
		m.previewLoading = false
		return m, nil

	case tickMsg:
		// Update time during conversion
		if m.converting {
			m.elapsedTime = time.Since(m.startTime)
			m.estimatedRemaining = estimateRemainingTime(m.startTime, m.currentFile, m.totalFiles)
			return m, tickCmd() // Re-schedule tick
		}
		return m, nil

	case conversionCompleteMsg:
		// Handle single file completion (for real-time updates)
		m.results = append(m.results, msg.result)
		m.currentFile++

		// Update counts
		switch msg.result.Status {
		case "success":
			m.completedFiles++
			// Add to completed list (keep last 5)
			m.completedList = append(m.completedList, msg.result.File)
			if len(m.completedList) > 5 {
				m.completedList = m.completedList[1:]
			}
		case "error":
			m.errorCount++
		case "warning":
			m.warningCount++
		}

		// Check if batch is complete
		if m.currentFile >= m.totalFiles {
			m.converting = false
			m.showSummary = true
		}

		return m, nil

	case batchCompleteMsg:
		// Batch conversion complete
		m.converting = false
		m.showSummary = true
		return m, nil

	case tea.QuitMsg:
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI
func (m model) View() tea.View {
	var content string

	if m.showSummary {
		// Show batch conversion summary
		content = renderSummaryScreen(m.results, m.elapsedTime)
	} else if m.converting {
		// Show batch conversion progress
		content = renderConversionScreen(m)
	} else if m.showSettingsEditor {
		// Show settings editor (Story 4-4)
		content = renderSettingsEditor(m)
	} else if m.showValidation {
		// Show validation screen (Story 4-4)
		content = renderValidationScreen(m)
	} else if m.showConfirmation {
		// Show confirmation screen
		content = confirmationScreen(m)
	} else if m.showFormatPrompt {
		// Show format selection menu
		content = formatMenuOptions()
	} else if m.showHelp {
		content = renderHelp(m)
	} else if m.err != nil {
		content = fmt.Sprintf("Error: %v\n\nPress 'q' to quit", m.err)
	} else if m.showPreview {
		// Show split-pane layout with preview (Story 4-2)
		content = renderSplitView(m)
	} else {
		// Show full-width file list (Story 4-1)
		content = renderFileList(m)
	}

	return tea.NewView(content)
}
