# Recipe TUI Guide

Interactive file browser for selecting and managing preset files.

## Overview

The Recipe TUI (Terminal User Interface) provides a keyboard-driven, visual file browser for navigating directories and selecting preset files (.np3, .xmp, .lrtemplate) for conversion.

## Features

### Visual File Browser
- Color-coded format badges for easy identification
- File size display (auto-formatted: B/KB/MB)
- Directory icons and indicators
- Scrolling support for large directories
- Terminal resize handling

### Multi-File Selection
- Select/deselect individual files with Space
- Select all files with 'a'
- Deselect all with 'n'
- Visual checkmarks (✓) on selected files
- Selection count display
- Persistent selection across navigation

### Keyboard Navigation
- Arrow keys or Vim-style (j/k) navigation
- Home/End for quick jumps
- Enter to navigate into directories
- Backspace/Left to go up to parent directory
- Live directory refresh with 'r'

## Quick Start

```bash
# Launch TUI in current directory
./recipe-tui

# Navigate with arrow keys or j/k
# Press Space to select files
# Press Enter to navigate into directories
# Press ? for help overlay
# Press q to quit
```

## Keyboard Reference

### Navigation Keys

| Key | Action |
|-----|--------|
| `↑` or `k` | Move cursor up |
| `↓` or `j` | Move cursor down |
| `←` or `Backspace` | Go to parent directory |
| `→` or `Enter` | Enter directory / Select file |
| `Home` | Jump to first item |
| `End` | Jump to last item |

### Selection Keys

| Key | Action |
|-----|--------|
| `Space` | Toggle selection on current file |
| `a` | Select all files in current directory |
| `n` | Deselect all files |

### Action Keys

| Key | Action |
|-----|--------|
| `r` | Refresh file list |
| `?` | Toggle help overlay |
| `q` or `Ctrl+C` | Quit application |
| `Esc` | Close help overlay |

## Visual Layout

```
┌──────────────────────────────────────────────────────────┐
│ Recipe - Preset File Browser                              │
│ Current: /Users/justin/presets                            │
├──────────────────────────────────────────────────────────┤
│                                                            │
│   📁 .. (parent directory)                                │
│ ✓ 🔵 NP3  portrait.np3                  1.2 KB           │
│   🟠 XMP  landscape.xmp                 856 B            │
│ ✓ 🟢 LRT  vintage.lrtemplate            2.3 KB           │
│   📁 subfolder/                                           │
│                                                            │
│ [2 files selected]                                        │
│                                                            │
│ Showing 1-5 of 42 files                                   │
│ Press '?' for help                                        │
└──────────────────────────────────────────────────────────┘
```

## Format Indicators

Files are displayed with color-coded badges:

- 🔵 **NP3** (Blue) - Nikon Picture Control binary format
- 🟠 **XMP** (Orange) - Adobe Lightroom sidecar XML
- 🟢 **LRT** (Green) - Adobe Lightroom Classic template
- 📁 **DIR** (Gray) - Directory

## File Filtering

The TUI automatically filters and displays only supported preset formats:
- `.np3` - Nikon Picture Control files
- `.xmp` - Adobe Lightroom sidecar files
- `.lrtemplate` - Lightroom Classic preset files
- Directories (for navigation)

Other file types are hidden from view.

## Selection Management

### Single Selection
1. Navigate to a file with arrow keys or j/k
2. Press Space to toggle selection
3. A checkmark (✓) appears next to selected files

### Multi-Selection
1. Navigate to first file and press Space
2. Move to next file and press Space
3. Repeat for all desired files
4. Selection count is shown at bottom

### Select All
Press `a` to select all preset files in the current directory. Directories are not selected.

### Clear Selection
Press `n` to deselect all files.

## Directory Navigation

### Enter Directory
1. Navigate to a directory (marked with 📁)
2. Press Enter or →
3. File list updates to show contents

### Go to Parent
1. Press Backspace or ←
2. Or navigate to ".." (parent directory) and press Enter

### Root Protection
When at filesystem root, attempting to go up has no effect.

## Help Overlay

Press `?` to display the help overlay with all keyboard shortcuts.

```
╔═══════════════════════════════════════════════════════════╗
║           Keyboard Shortcuts                               ║
╠═══════════════════════════════════════════════════════════╣
║                                                            ║
║ Navigation:                                                ║
║   ↑/k, ↓/j           Move cursor up/down                  ║
║   ←/Backspace        Go to parent directory               ║
║   →/Enter            Enter directory                      ║
║   Home/End           Jump to first/last item              ║
║                                                            ║
║ Selection:                                                 ║
║   Space              Toggle file selection                ║
║   a                  Select all files                     ║
║   n                  Deselect all                         ║
║                                                            ║
║ Actions:                                                   ║
║   r                  Refresh file list                    ║
║   ?                  Toggle this help                     ║
║   q/Ctrl+C           Quit                                 ║
║                                                            ║
║ Press ? or Esc to close                                    ║
╚═══════════════════════════════════════════════════════════╝
```

Press `?` or `Esc` to close the help overlay.

## Technical Details

### Implementation
- Built with Bubbletea v2 (Elm Architecture)
- Lipgloss v2 for terminal styling
- Model-Update-View pattern
- Fully testable architecture (88.5% coverage)

### Terminal Requirements
- Minimum width: 60 columns
- Minimum height: 10 rows
- UTF-8 support for icons
- ANSI color support (optional, falls back to monochrome)

### Performance
- Instant directory loading (<10ms)
- Smooth scrolling for 1000+ files
- Efficient redraw on terminal resize
- No blocking operations

## Troubleshooting

### Icons Not Displaying
If emoji icons (📁, ✓) don't display correctly:
- Ensure terminal has UTF-8 encoding enabled
- Use a font with emoji support (e.g., Cascadia Code, Fira Code)
- Functionality works without icons, just visual degradation

### Colors Not Showing
If colors don't display:
- Check terminal supports ANSI colors
- Verify TERM environment variable is set
- Application works in monochrome mode as fallback

### Terminal Resize Issues
If layout breaks after resize:
- Press `r` to refresh display
- Ensure terminal size is at least 60x10
- Try quitting and restarting

### File List Not Updating
If files don't appear after adding them:
- Press `r` to manually refresh
- Check file permissions
- Verify files have correct extensions (.np3, .xmp, .lrtemplate)

## Integration with CLI

While the TUI provides file selection, conversion still requires the CLI:

```bash
# 1. Use TUI to browse and identify files
./recipe-tui

# 2. Use CLI to convert selected files
./recipe convert portrait.xmp --to np3
./recipe batch *.xmp --to np3
```

Future versions may integrate conversion directly into the TUI interface.

## Architecture

The TUI follows the Elm Architecture (Model-Update-View):

### Model
State management including:
- Current directory path
- File list with metadata
- Cursor position
- Selection map
- Terminal dimensions
- Help overlay state

### Update
Message handling for:
- Keyboard events
- File loading
- Window resize
- Navigation
- Selection toggling

### View
Rendering logic for:
- File list display
- Format badges
- Selection indicators
- Help overlay
- Scroll indicators

## Testing

Comprehensive test suite with 88.5% coverage:

```bash
# Run all TUI tests
go test ./cmd/tui -v

# Run with coverage report
go test ./cmd/tui -cover

# Run specific test
go test ./cmd/tui -run TestNavigationKeys
```

## Future Enhancements

Planned features for future releases:
- Inline conversion support (convert without CLI)
- Preview pane showing preset parameters
- Batch conversion progress display
- Favorite directories bookmark system
- Search/filter by filename
- Sort by name/size/date
- Color theme customization

## See Also

- [README.md](../README.md) - Main project documentation
- [Epic 4 Tech Spec](tech-spec-epic-4.md) - TUI technical specification
- [Story 4-1](stories/4-1-bubbletea-file-browser.md) - File browser story
