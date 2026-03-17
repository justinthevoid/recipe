# Recipe NP3 Editor — VSCode Extension

A custom visual editor for Nikon NP3 picture control files (`.np3`) directly inside VSCode. Open any NP3 file to get an interactive parameter editor with live photo preview — no external tools needed.

> **Status:** Fully functional for editing and saving NP3 presets. The live photo preview applies filters in real-time but is not yet color-accurate compared to Nikon NX Studio's rendering. Not published to the VSCode Marketplace.

## Features

### Parameter Editing
- Visual sliders for all supported NP3 parameters
- Exposure, Contrast, Highlights, Shadows, Whites, Blacks
- HSL Color Mixer (Hue, Saturation, Luminance per channel)
- Color Grading (Shadows, Midtones, Highlights wheels)
- Sharpening and Mid-Range Sharpening
- Clarity and Saturation
- Tone curve visualization

### Live Photo Preview
- Real-time filter preview powered by WebAssembly
- Toggle between Editor (two-column parameters) and Preview (photo + controls) layouts
- LUT-based filter generation for instant visual feedback

### Editing Workflow
- Full undo/redo with 100-entry edit history and debounced coalescing
- Copy/paste parameters between open NP3 files
- Dirty-state tracking with `[Modified]` indicator
- Save, Save As, and Reset All commands
- Automatic `.bak` backup file on open
- Multi-tab editing (up to 20 concurrent files, configurable)
- Keyboard shortcuts: `Ctrl+S` save, `Ctrl+Z`/`Ctrl+Y` undo/redo, `Ctrl+C`/`Ctrl+V` copy/paste parameters

### Reliability
- Crash recovery with automatic binary restart (up to 3 retries)
- Corrupted file detection with hex dump inspection UI
- Webview state persists across tab switches

## Architecture

```
┌─────────────┐     JSONL IPC     ┌──────────┐
│  VSCode     │◄────────────────►│  np3tool  │
│  Extension  │                   │  (Go)    │
│  (TS)       │                   └──────────┘
└──────┬──────┘
       │ postMessage
┌──────▼──────┐     WASM      ┌──────────┐
│  Webview    │◄─────────────►│  recipe   │
│  (Svelte 5) │               │  (Go)    │
└─────────────┘               └──────────┘
```

- **Extension host** (TypeScript): Manages editor panels, file I/O, and binary lifecycle
- **np3tool** (Go): JSONL-based IPC server for NP3 parsing, patching, and saving
- **Webview** (Svelte 5 + Tailwind): Interactive parameter UI with real-time preview
- **WASM module** (Go): Photo filter generation for live preview

## Building

Requires Go 1.25.1+, Node.js 18+, and [Bun](https://bun.sh).

```bash
# From the repository root:

# 1. Build the Go backend binary
go build -o extension/bin/np3tool cmd/np3tool/main.go

# 2. Build the WASM module for preview
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o webview/public/recipe.wasm cmd/wasm/main.go

# 3. Install dependencies and build the webview
cd webview && bun install && bun run build && cd ..

# 4. Build the extension
cd extension && bun install && bun run build && cd ..
```

### Development

```bash
# Terminal 1: Watch webview changes
cd webview && bun run dev

# Terminal 2: Watch extension changes
cd extension && bun run dev

# Then press F5 in VSCode to launch the Extension Development Host
```

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `recipe.maxConcurrentPanels` | 5 | Maximum number of NP3 files open simultaneously (1-20). Each file runs a Go backend process. |

## Commands

| Command | Shortcut | Description |
|---------|----------|-------------|
| Recipe: Save As... | `Ctrl+Shift+S` | Save the current preset to a new file |
| Recipe: Reset All Parameters | — | Reset all parameters to their original values |

## Known Limitations

- **Photo preview is not color-accurate**: Filter rendering approximates NP3 adjustments but does not match Nikon NX Studio's proprietary rendering pipeline
- **NP3 format constraints apply**: Temperature, Tint, Vibrance, custom tone curves, and vignette are not available in the standard 480-byte NP3 format (see [docs/known-conversion-limitations.md](../docs/known-conversion-limitations.md))
- **Requires workspace trust**: The extension executes a local Go binary for file parsing

## License

[MIT](../LICENSE)
