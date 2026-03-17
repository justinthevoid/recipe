# Recipe

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Cloudflare Pages](https://img.shields.io/badge/Deployed-Cloudflare%20Pages-F38020?logo=cloudflare&logoColor=white)](https://recipe.pages.dev)
[![WebAssembly](https://img.shields.io/badge/WebAssembly-654FF0?logo=webassembly&logoColor=white)](https://webassembly.org)
[![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?logo=svelte&logoColor=white)](https://svelte.dev)

Convert photo presets between Nikon NP3 and Adobe Lightroom XMP formats.

All processing happens locally on your device. Files are never uploaded to any server.

**Web app:** [recipe.shuttercoach.app](https://recipe.shuttercoach.app) -- no install needed, runs entirely in your browser via WebAssembly.

## Features

- **Privacy-first** -- all conversions run client-side (WebAssembly in browser, or local CLI)
- **Fast** -- sub-millisecond conversions (0.003-0.079ms per file)
- **Accurate** -- 98%+ conversion fidelity via exact offset mapping for 48 NP3 parameters
- **Web + CLI** -- browser-based converter and command-line tool for scripting/batch workflows

## Quick Start

### Web

Visit [recipe.shuttercoach.app](https://recipe.shuttercoach.app), drag and drop your preset files, and download the converted output.

### CLI

Download a pre-built binary from [GitHub Releases](https://github.com/justinthevoid/recipe/releases/latest), or build from source:

```bash
make cli
./recipe convert portrait.xmp --to np3
./recipe convert preset.np3 --to xmp
./recipe batch *.xmp --to np3
```

## Building from Source

Requires Go 1.25.1+ and Node.js 18+.

```bash
# CLI
make cli

# WASM module for web interface
make wasm

# Web dev server (hot reload)
cd web && npm install && npm run dev
```

## How It Works

Recipe uses a hub-and-spoke conversion pattern with a `UniversalRecipe` intermediate representation:

```
NP3 --Parse--> UniversalRecipe --Generate--> XMP
XMP --Parse--> UniversalRecipe --Generate--> NP3
```

All conversions flow through a single entry point: `converter.Convert(input, from, to)`. Format parsers and generators are isolated in `internal/formats/`.

## Supported Formats

| Format | Extension | Used In            | Direction |
| ------ | --------- | ------------------ | --------- |
| NP3    | .np3      | Nikon Z cameras    | NP3 <-> XMP |
| XMP    | .xmp      | Adobe Lightroom CC | XMP <-> NP3 |

### Known Limitations

NP3 is a compact binary format with fewer parameters than XMP. Some XMP features cannot be represented in NP3:

- Vibrance, Temperature/Tint, Grain, Vignette, Parametric Tone Curves

See [docs/known-conversion-limitations.md](docs/known-conversion-limitations.md) for the full list.

## VSCode Extension

A custom editor for NP3 preset files directly in VSCode. Open any `.np3` file to get a visual parameter editor with live photo preview.

**Features:**
- Visual sliders for all NP3 parameters (exposure, contrast, HSL color, color grading, sharpening, etc.)
- Live photo preview with real-time filter updates via WebAssembly
- Tone curve visualization
- Undo/redo with full edit history
- Copy/paste parameters between open presets
- Multi-tab editing with dirty-state tracking
- Automatic backup file creation on open
- Crash recovery with auto-restart

**Architecture:** The extension runs a local Go binary (`np3tool`) for file parsing and a Svelte 5 webview for the UI, communicating over JSONL IPC.

**Status:** Fully functional for editing and saving NP3 presets. The photo preview applies filters in real-time but is not yet color-accurate compared to Nikon NX Studio's rendering. Not yet published to the VSCode Marketplace.

See [extension/](extension/) for build instructions.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## Legal

The NP3 format was analyzed through clean-room methods for interoperability purposes, protected under DMCA Section 1201(f).

## Support

If you find Recipe useful, consider supporting development:

[![GitHub Sponsors](https://img.shields.io/badge/Sponsor-GitHub-ea4aaa?logo=github)](https://github.com/sponsors/justinthevoid)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?logo=buymeacoffee&logoColor=black)](https://buymeacoffee.com/justinthevoid)

## License

[MIT](LICENSE)
