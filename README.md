# Recipe

Convert photo presets between Nikon NP3 and Adobe Lightroom XMP formats.

All processing happens locally on your device. Files are never uploaded to any server.

**Web app:** [recipe.pages.dev](https://recipe.pages.dev) -- no install needed, runs entirely in your browser via WebAssembly.

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

A VSCode extension for NP3 preset editing is in progress. See [extension/](extension/) for details.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## Legal

The NP3 format was analyzed through clean-room methods for interoperability purposes, protected under DMCA Section 1201(f).

## License

[MIT](LICENSE)
