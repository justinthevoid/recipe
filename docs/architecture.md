# Architecture

**Last Updated:** 2026-02-04
**Documentation Generated:** Exhaustive Scan (Full Rescan)

## Executive Summary

Recipe is a universal photo preset converter built with a privacy-first, client-side architecture. The system provides multiple interfaces (CLI, NX Studio Integration, Web) that all share a common conversion engine using a hub-and-spoke pattern with a universal intermediate representation (UniversalRecipe).

**Key Architectural Characteristics:**
- **Performance**: <100ms WASM conversions, sub-millisecond CLI conversions (0.003-0.079ms per file)
- **Privacy**: Zero server uploads, all processing client-side (Web) or local (CLI)
- **Accuracy**: 98%+ conversion fidelity for core adjustments via round-trip testing
- **Simplicity**: Single conversion API (`converter.Convert()`) shared across all interfaces
- **Minimal Dependencies**: Core library uses Go standard library + minimal external packages

**Technology Foundation:**
- Go 1.25.1 (latest WASM support with go:wasmexport)
- Svelte 5.43.8 + Vite 7.2.4 (Web interface)
- Cloudflare Pages (static hosting for Web interface)
- Cobra CLI framework (command-line interface)
- WebAssembly for browser-based conversion

**Build Commands:**
```bash
# Build CLI for current platform
make cli

# Build CLI for all platforms
make cli-all

# Build WASM module (production)
make wasm

# Build NX Studio integration
make build-nx

# Run tests
make test

# Run tests with coverage
make coverage
```

---

## Decision Summary

| Category | Decision | Version/Choice | Rationale |
|----------|----------|----------------|-----------|
| **Critical Decisions** | | | |
| Frontend Language | Vanilla JavaScript (ES6+) | No frameworks | Zero build step, simple deployment, <100ms WASM goal achievable |
| Language Version | Go 1.24.0+ | Released Feb 2025 | go:wasmexport directive, reduced WASM memory, enhanced type support |
| Project Structure | cmd/{cli,tui,wasm}/ + internal/ + web/ | Standard Go layout | Clear separation, shared internal/ packages, WASM-friendly |
| Shared Library Design | `converter.Convert([]byte, string, string) ([]byte, error)` | Single API | Stateless, no OS deps, works in WASM, easy to test |
| **Important Decisions** | | | |
| Error Handling | Wrapped errors with custom ConversionError type | stdlib errors | Type-safe error checking, format-specific context, no dependencies |
| WebAssembly Bridge | go:wasmexport directive | Go 1.24+ | Direct memory access, <100ms goal, zero reflection overhead |
| File I/O (Web) | FileReader API → ArrayBuffer pattern | Standard Web API | Works with go:wasmexport, zero dependencies, 90%+ browser support |
| Testing Strategy | Table-driven tests with 1,501 real sample files | Standard Go | Comprehensive validation, 95%+ accuracy goal, idiomatic Go |
| Logging | slog with structured fields | Go stdlib | Type-safe, log levels, zero dependencies, Go 1.21+ standard |
| Build System | Makefile with targets for all interfaces | make | Simple, universal, 90%+ dev env support, no learning curve |
| **Nice-to-Have Decisions** | | | |
| CSS Framework | Vanilla CSS with CSS custom properties | No framework | Simple, fast, no build step, <10 KB, consistent with JS choice |
| State Management | DOM-centric with vanilla JS | No framework | Direct manipulation, fast, no state library needed for simple UI |
| Validation Strategy | Inline validation in parsers | No separate layer | Fail-fast, clear error messages, simpler architecture |
| Performance Monitoring | Benchmarks only (`go test -bench`) | No runtime metrics | Development-focused, sufficient for <100ms validation |
| CI/CD Platform | GitHub Actions with Cloudflare Pages integration | Free tier | Zero cost, auto-deploy Web, build all interfaces, standard choice |
| Linting/Formatting | gofmt + go vet (stdlib) | No external linters | Zero dependencies, good enough, fast, standard Go practice |
| Pre-commit Hooks | None initially | Can add later | Simpler workflow, faster iteration, add if team grows |
| Documentation | README + examples/ + godoc comments | No separate docs | Godoc for API reference, README for users, examples for learning |

---

## Project Structure

```
recipe/
├── cmd/                           # Entry points for all interfaces
│   ├── cli/                       # Main CLI application (Cobra)
│   │   ├── main.go                # Entry point
│   │   ├── root.go                # Root command definition
│   │   ├── convert.go             # Convert command
│   │   ├── batch.go               # Batch conversion command
│   │   ├── inspect.go             # Preset inspection command
│   │   ├── diff.go                # Preset diff command
│   │   └── format.go              # Format detection command
│   ├── nx/                        # NX Studio integration CLI
│   │   ├── main.go                # Entry point for recipe-nx
│   │   ├── apply.go               # Apply recipe to NEF files
│   │   ├── batch.go               # Batch processing
│   │   └── verify.go              # Verify applied recipes
│   ├── wasm/
│   │   └── main.go                # WASM export entry point
│   └── debug_curve/               # Development utilities
│       └── main.go
│
├── internal/                      # Private application code (65 Go files)
│   ├── apperr/                    # Application error types
│   │   └── error.go
│   ├── batch/                     # Batch processing orchestrator
│   │   ├── orchestrator.go        # Parallel file processing
│   │   ├── manifest.go            # Idempotent processing manifest
│   │   └── file_ops.go            # File operations
│   ├── converter/                 # Hub-and-spoke conversion engine
│   │   ├── converter.go           # Convert([]byte, string, string) ([]byte, error)
│   │   └── error.go               # ConversionError type
│   ├── formats/                   # Format parsers/generators
│   │   ├── np3/                   # Nikon Picture Control (.np3)
│   │   │   ├── parse.go           # Binary parsing with exact offsets
│   │   │   ├── generate.go        # Binary generation
│   │   │   ├── offsets.go         # Byte offset constants
│   │   │   ├── curvegen.go        # Tone curve generation
│   │   │   └── metadata.go        # NP3 metadata extraction
│   │   ├── xmp/                   # Adobe XMP (.xmp)
│   │   │   ├── parse.go           # XML parsing with struct tags
│   │   │   └── generate.go        # XML generation
│   │   └── (other format packages archived)
│   ├── inspect/                   # Preset inspection utilities
│   ├── lut/                       # LUT processing
│   ├── models/                    # Core data structures
│   │   └── recipe.go              # UniversalRecipe (157 lines, 50+ fields)
│   ├── testutil/                  # Test utilities
│   ├── utils/                     # Shared utilities
│   └── verify/                    # Verification utilities
│
├── web/                           # Svelte web interface
│   ├── index.html                 # Main HTML entry
│   ├── package.json               # Svelte 5 + Vite 7
│   ├── vite.config.js             # Vite configuration
│   ├── svelte.config.js           # Svelte configuration
│   ├── public/
│   │   ├── recipe.wasm            # Compiled WASM binary
│   │   └── wasm_exec.js           # Go WASM runtime
│   └── src/
│       ├── App.svelte             # Main application component
│       ├── main.js                # Entry point
│       └── lib/                   # Svelte components and utilities
│           ├── wasm.js            # WASM initialization
│           ├── converter.js       # Conversion logic
│           └── components/        # UI components
│
├── testdata/                      # Test fixtures (302 items)
│   ├── np3/                       # Nikon Picture Control samples
│   ├── xmp/                       # Adobe XMP samples
│   └── nx-fixtures/               # NX Studio integration fixtures
│
├── docs/                          # Documentation (25 files)
├── bin/                           # Build output directory
├── Makefile                       # Build automation (139 lines)
├── go.mod                         # Go 1.25.1
├── go.sum
├── README.md                      # Main documentation (35KB)
└── CHANGELOG.md                   # Version history
```

**Key Design Decisions:**
- **cmd/** separates interfaces (cli, nx, wasm) - all use `internal/converter`
- **internal/** prevents external imports, enforces API boundaries
- **internal/converter** is the single source of truth for conversion logic
- **internal/formats** follows identical pattern for each format (parse.go, generate.go)
- **internal/batch** provides parallel processing with idempotency via manifests
- **web/** uses Svelte 5 + Vite 7 for modern component-based UI

---

## Technology Stack Details

### Core Technology: Go 1.25.1

**Selection Rationale:**
- **go:wasmexport**: Direct memory access for WASM exports (zero reflection overhead)
- **Reduced WASM Binary Size**: Optimized WASM builds with `-ldflags="-s -w"`
- **Enhanced Type Support**: Better TypeScript interop from WASM exports
- **slog**: Structured logging built into stdlib (Go 1.21+)

**Go Dependencies (go.mod):**
```go
require (
    github.com/google/tiff v0.0.0-20161109161721-4b31f3041d9a  // TIFF handling
    github.com/lucasb-eyer/go-colorful v1.3.0                  // Color space conversions
    github.com/spf13/cobra v1.10.1                             // CLI framework
    golang.org/x/image v0.34.0                                 // Image processing
    golang.org/x/term v0.36.0                                  // Terminal detection
    gopkg.in/yaml.v3 v3.0.1                                    // YAML parsing
)
```

**Validation:**
```bash
# Verify Go version
go version  # Should output: go version go1.24.0 or higher

# Test WASM build
GOOS=js GOARCH=wasm go build -o test.wasm ./cmd/wasm
```

### Frontend: Vanilla JavaScript (ES6+)

**Selection Rationale:**
- **Zero Build Step**: No webpack, no npm, no node_modules - just HTML/JS/CSS
- **90%+ Browser Support**: FileReader API, WebAssembly supported since 2017
- **Performance**: Direct DOM manipulation faster than framework diffing
- **Simplicity**: Fewer moving parts, easier debugging, faster iteration

**Browser Compatibility:**
- Chrome 57+ (March 2017)
- Firefox 52+ (March 2017)
- Safari 11+ (September 2017)
- Edge 16+ (September 2017)

### Deployment: Cloudflare Pages

**Selection Rationale:**
- **Zero Cost**: Free tier supports unlimited static sites
- **Auto-Deploy**: Push to `main` branch → automatic build & deploy
- **Global CDN**: Sub-100ms latency worldwide
- **Simple Setup**: Connect GitHub repo, no configuration needed

**Deployment Flow:**
1. Push code to GitHub repository
2. Cloudflare Pages detects push to `main` branch
3. Automatic build (if build command configured, otherwise direct deploy)
4. Deploy to https://recipe.pages.dev (or custom domain)

### CLI Framework: Cobra

**Selection Rationale:**
- **Official Generator**: `cobra-cli init` scaffolds complete structure
- **Standard Pattern**: Used by kubectl, hugo, gh - developers familiar
- **Command Hierarchy**: Easy to add subcommands (convert, batch, validate)
- **Flag Management**: Automatic help generation, completion scripts

**Initialization:**
```bash
cobra-cli init
cobra-cli add convert
cobra-cli add batch
```

### TUI Framework: Bubbletea

**Selection Rationale:**
- **Elm Architecture**: Predictable state updates (Model, Update, View)
- **Composable**: Combine components from Bubbles library
- **Cross-Platform**: Works on Windows, macOS, Linux terminals
- **Lipgloss Styling**: CSS-like styling for terminal UIs

---

## Implementation Patterns

### Pattern 1: Package Naming Convention

**Rule**: All packages under `internal/` use lowercase, singular nouns

**Examples:**
```go
// CORRECT
internal/converter/
internal/model/
internal/formats/np3/

// INCORRECT
internal/converters/        # Plural
internal/Models/            # CamelCase
internal/formats/NP3/       # Uppercase acronym
```

**Rationale**: Go convention, easier to import, clearer package purpose

---

### Pattern 2: Function and Method Naming

**Rule**:
- Exported functions/methods: CamelCase starting with capital
- Unexported functions/methods: camelCase starting with lowercase
- Test functions: `Test` + function name

**Examples:**
```go
// CORRECT - Exported
func Parse(data []byte) (*UniversalRecipe, error)
func Generate(recipe *UniversalRecipe) ([]byte, error)
func Convert(input []byte, from, to string) ([]byte, error)

// CORRECT - Unexported
func normalizeRange(val, min, max int) int
func validateFormat(format string) error

// CORRECT - Tests
func TestParse(t *testing.T)
func TestGenerate(t *testing.T)
func TestParseInvalidMagic(t *testing.T)

// INCORRECT
func parse()                 # Unexported but should be public API
func PARSE()                 # All caps
func test_parse()            # Snake case
```

---

### Pattern 3: Type Naming

**Rule**:
- Structs: CamelCase, descriptive nouns
- Interfaces: CamelCase, often end with -er (Parser, Generator)
- Errors: End with Error suffix

**Examples:**
```go
// CORRECT - Structs
type UniversalRecipe struct { ... }
type ConversionError struct { ... }
type ColorAdjustment struct { ... }

// CORRECT - Interfaces
type Parser interface { ... }
type Generator interface { ... }

// CORRECT - Errors
type ConversionError struct {
    Operation string
    Format    string
    Cause     error
}

// INCORRECT
type universal_recipe struct { ... }  # Snake case
type IParser interface { ... }        # Hungarian notation
type ErrorConversion struct { ... }   # Wrong order
```

---

### Pattern 4: File Structure for Format Packages

**Rule**: Every format package has identical structure

```
internal/formats/{format}/
├── parse.go          # Parse([]byte) (*UniversalRecipe, error)
├── generate.go       # Generate(*UniversalRecipe) ([]byte, error)
└── {format}_test.go  # Table-driven tests
```

**Examples:**

**parse.go:**
```go
package np3

import (
    "fmt"
    "recipe/internal/model"
)

// Parse converts NP3 binary data to UniversalRecipe
func Parse(data []byte) (*model.UniversalRecipe, error) {
    if len(data) < 1024 {
        return nil, fmt.Errorf("invalid NP3 file: too small")
    }

    // Parse implementation
    return &model.UniversalRecipe{ /* ... */ }, nil
}
```

**generate.go:**
```go
package np3

import (
    "fmt"
    "recipe/internal/model"
)

// Generate converts UniversalRecipe to NP3 binary format
func Generate(recipe *model.UniversalRecipe) ([]byte, error) {
    if recipe == nil {
        return nil, fmt.Errorf("recipe is nil")
    }

    // Generate implementation
    return []byte{ /* ... */ }, nil
}
```

**np3_test.go:**
```go
package np3

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParse(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/np3/*.np3")

    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            data, err := os.ReadFile(file)
            if err != nil {
                t.Fatal(err)
            }

            recipe, err := Parse(data)
            if err != nil {
                t.Errorf("Parse() error = %v", err)
            }

            if recipe.Contrast < -100 || recipe.Contrast > 100 {
                t.Errorf("Contrast out of range: %d", recipe.Contrast)
            }
        })
    }
}
```

**Rationale**: Identical structure across formats makes it easy to add new formats. Copy-paste pattern, update implementation.

---

### Pattern 5: Error Handling

**Rule**: Always return wrapped `ConversionError` for conversion failures

**Implementation:**
```go
// internal/converter/errors.go
package converter

import "fmt"

type ConversionError struct {
    Operation string  // "parse", "generate", "validate"
    Format    string  // "np3", "xmp"
    Cause     error   // Underlying error
}

func (e *ConversionError) Error() string {
    return fmt.Sprintf("%s %s failed: %v", e.Operation, e.Format, e.Cause)
}

func (e *ConversionError) Unwrap() error {
    return e.Cause
}
```

**Usage:**
```go
// internal/converter/converter.go
func Convert(input []byte, from, to string) ([]byte, error) {
    recipe, err := parse(input, from)
    if err != nil {
        return nil, &ConversionError{
            Operation: "parse",
            Format:    from,
            Cause:     err,
        }
    }

    output, err := generate(recipe, to)
    if err != nil {
        return nil, &ConversionError{
            Operation: "generate",
            Format:    to,
            Cause:     err,
        }
    }

    return output, nil
}
```

**Rationale**: Type-safe error checking, format-specific context, no external dependencies

---

### Pattern 6: Validation Strategy

**Rule**: Validate inline in parsers, fail fast with descriptive errors

**Examples:**
```go
func Parse(data []byte) (*model.UniversalRecipe, error) {
    // Validate magic bytes FIRST
    if len(data) < 2 || string(data[0:2]) != "NP" {
        return nil, fmt.Errorf("invalid NP3 magic: expected 'NP', got '%s'", data[0:2])
    }

    // Validate file size
    if len(data) < 1024 {
        return nil, fmt.Errorf("invalid NP3 file: expected >= 1024 bytes, got %d", len(data))
    }

    // Validate parameter ranges during extraction
    contrast := normalizeRange(int(data[66]), -100, 100)
    if contrast < -100 || contrast > 100 {
        return nil, fmt.Errorf("contrast out of range: %d (expected -100 to +100)", contrast)
    }

    return &model.UniversalRecipe{Contrast: contrast}, nil
}
```

**Rationale**: Fail-fast prevents invalid data from propagating, clear error messages aid debugging

---

### Pattern 7: Testing Strategy

**Rule**: Table-driven tests using 1,501 real sample files from `testdata/`

**Examples:**
```go
func TestParse(t *testing.T) {
    // Discover all test files
    files, err := filepath.Glob("../../../testdata/np3/*.np3")
    if err != nil {
        t.Fatal(err)
    }

    if len(files) == 0 {
        t.Fatal("no test files found in testdata/np3/")
    }

    // Run subtest for each file
    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            data, err := os.ReadFile(file)
            if err != nil {
                t.Fatalf("failed to read %s: %v", file, err)
            }

            recipe, err := Parse(data)
            if err != nil {
                t.Errorf("Parse() error = %v", err)
                return
            }

            // Validate all critical fields
            if recipe.Contrast < -100 || recipe.Contrast > 100 {
                t.Errorf("Contrast out of range: %d", recipe.Contrast)
            }

            if recipe.Saturation < -100 || recipe.Saturation > 100 {
                t.Errorf("Saturation out of range: %d", recipe.Saturation)
            }
        })
    }
}

func TestRoundTrip_XMP_LRTemplate(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/xmp/*.xmp")

    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            // Step 1: Parse XMP
            origData, _ := os.ReadFile(file)
            orig, err := xmp.Parse(origData)
            if err != nil {
                t.Fatalf("XMP parse failed: %v", err)
            }

            // Step 2: Generate NP3
            np3Data, err := np3.Generate(orig)
            if err != nil {
                t.Fatalf("NP3 generate failed: %v", err)
            }

            // Step 3: Parse NP3 back
            recovered, err := np3.Parse(np3Data)
            if err != nil {
                t.Fatalf("NP3 parse failed: %v", err)
            }

            // Step 4: Compare critical fields
            tolerance := 1 // Allow ±1 due to rounding
            if abs(orig.Contrast - recovered.Contrast) > tolerance {
                t.Errorf("Contrast mismatch: orig=%d, recovered=%d", orig.Contrast, recovered.Contrast)
            }
        })
    }
}
```

**Run Tests:**
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/formats/np3/

# Run specific test
go test -run TestParse ./internal/formats/np3/

# Benchmark (for <100ms validation)
go test -bench=. ./internal/converter/
```

**Rationale**:
- 1,501 real sample files = comprehensive validation
- Round-trip tests ensure 95%+ accuracy goal
- Table-driven pattern is idiomatic Go
- Subtests provide granular failure reporting

---

### Pattern 8: Logging Strategy

**Rule**: Use `slog` with consistent structured fields

**Examples:**
```go
import "log/slog"

// Initialize logger
var logger = slog.Default()

// Log conversion start
logger.Info("starting conversion",
    slog.String("file", inputPath),
    slog.String("from", sourceFormat),
    slog.String("to", targetFormat),
)

// Log conversion success
logger.Info("conversion complete",
    slog.String("file", inputPath),
    slog.Duration("elapsed", time.Since(start)),
)

// Log conversion error
logger.Error("conversion failed",
    slog.String("file", inputPath),
    slog.String("error", err.Error()),
)

// Debug logging (only in verbose mode)
logger.Debug("parsed recipe",
    slog.String("name", recipe.Name),
    slog.Int("contrast", recipe.Contrast),
    slog.Int("saturation", recipe.Saturation),
)
```

**Log Levels:**
- **Debug**: Internal state, parsed values (verbose mode only)
- **Info**: Conversion start/complete, file counts
- **Warn**: Missing optional fields, approximated values
- **Error**: Conversion failures, invalid input

**Rationale**: Structured logging with slog (Go 1.21+) provides type safety, log levels, zero dependencies

---

### Pattern 9: WASM Export Pattern

**Rule**: Use `go:wasmexport` for direct memory access, return (ptr, len, error) tuple

**Implementation:**
```go
// cmd/wasm/main.go
package main

import "unsafe"

//go:wasmexport convertPreset
func convertPreset(inputPtr, inputLen uint32, srcFormat, dstFormat string) (uint32, uint32, string) {
    // Step 1: Reconstruct input slice from pointer
    input := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(inputPtr))), inputLen)

    // Step 2: Call internal converter
    output, err := converter.Convert(input, srcFormat, dstFormat)
    if err != nil {
        return 0, 0, err.Error()
    }

    // Step 3: Return pointer and length to output
    outputPtr := uint32(uintptr(unsafe.Pointer(&output[0])))
    outputLen := uint32(len(output))

    return outputPtr, outputLen, ""
}

func main() {
    // WASM exports don't need main, but Go requires it
}
```

**JavaScript Bridge:**
```javascript
// web/main.js
async function convertFile(inputBytes, sourceFormat, targetFormat) {
    // Allocate memory in WASM
    const inputPtr = await wasmInstance.exports.malloc(inputBytes.length);
    const inputBuf = new Uint8Array(wasmMemory.buffer, inputPtr, inputBytes.length);
    inputBuf.set(inputBytes);

    // Call WASM export
    const result = await wasmInstance.exports.convertPreset(
        inputPtr,
        inputBytes.length,
        sourceFormat,
        targetFormat
    );

    // Extract results
    const [outputPtr, outputLen, errorMsg] = result;

    if (errorMsg) {
        throw new Error(errorMsg);
    }

    // Read output from WASM memory
    const outputBytes = new Uint8Array(wasmMemory.buffer, outputPtr, outputLen);

    // Free memory
    await wasmInstance.exports.free(inputPtr);

    return outputBytes;
}
```

**Rationale**:
- `go:wasmexport` provides direct memory access (zero reflection overhead)
- Faster than JSON serialization (<100ms goal)
- Standard pattern for Go 1.24+ WASM exports

---

### Pattern 10: CLI Command Pattern

**Rule**: Every CLI command follows Cobra structure with RunE for error handling

**Implementation:**
```go
// cmd/cli/convert.go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "recipe/internal/converter"
)

var convertCmd = &cobra.Command{
    Use:   "convert [input] [output]",
    Short: "Convert preset between formats",
    Long: `Convert photo presets between NP3 and XMP formats.

Examples:
  recipe convert portrait.np3 portrait.xmp
  recipe convert --from np3 --to xmp portrait.np3`,
    Args: cobra.MinimumNArgs(1),
    RunE: runConvert,
}

func init() {
    convertCmd.Flags().StringP("from", "f", "", "Source format (auto-detected if omitted)")
    convertCmd.Flags().StringP("to", "t", "", "Target format (required if single argument)")
    convertCmd.Flags().StringP("output", "o", "", "Output file path")
}

func runConvert(cmd *cobra.Command, args []string) error {
    // Parse flags
    fromFormat, _ := cmd.Flags().GetString("from")
    toFormat, _ := cmd.Flags().GetString("to")
    outputPath, _ := cmd.Flags().GetString("output")

    inputPath := args[0]

    // Auto-detect formats if not specified
    if fromFormat == "" {
        fromFormat = detectFormat(inputPath)
    }

    if toFormat == "" && len(args) > 1 {
        toFormat = detectFormat(args[1])
        outputPath = args[1]
    }

    if toFormat == "" {
        return fmt.Errorf("target format required (use --to or provide output file)")
    }

    // Read input
    input, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input: %w", err)
    }

    // Convert
    output, err := converter.Convert(input, fromFormat, toFormat)
    if err != nil {
        return fmt.Errorf("convert: %w", err)
    }

    // Write output
    if outputPath == "" {
        outputPath = replaceExtension(inputPath, toFormat)
    }

    if err := os.WriteFile(outputPath, output, 0644); err != nil {
        return fmt.Errorf("write output: %w", err)
    }

    fmt.Printf("✓ Converted %s → %s\n", inputPath, outputPath)
    return nil
}
```

**Rationale**: RunE returns errors for graceful handling by Cobra, consistent command structure

---

## Consistency Rules for AI Agents

**Critical**: All AI agents implementing this system MUST follow these rules to ensure consistency.

### Rule 1: Single API Rule
**All interfaces MUST use `converter.Convert()` - no direct format parser calls**

```go
// CORRECT - CLI
output, err := converter.Convert(input, "np3", "xmp")

// CORRECT - TUI
output, err := converter.Convert(input, sourceFormat, targetFormat)

// CORRECT - WASM
output, err := converter.Convert(input, srcFormat, dstFormat)

// INCORRECT - Direct parser call bypasses shared logic
recipe, _ := np3.Parse(input)
output, _ := xmp.Generate(recipe)
```

### Rule 2: Error Type Rule
**All conversion errors MUST be wrapped in `ConversionError`**

```go
// CORRECT
return &ConversionError{
    Operation: "parse",
    Format:    "np3",
    Cause:     err,
}

// INCORRECT - Bare error loses context
return err

// INCORRECT - Generic error loses type safety
return fmt.Errorf("parse failed: %w", err)
```

### Rule 3: Testing Rule
**All format parsers MUST have table-driven tests using real sample files**

```go
// CORRECT
func TestParse(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/np3/*.np3")
    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            // Test implementation
        })
    }
}

// INCORRECT - Hardcoded test data
func TestParse(t *testing.T) {
    data := []byte{0x4E, 0x50, ...}  // Magic bytes hardcoded
    // ...
}
```

### Rule 4: Logging Rule
**All conversion operations MUST log with structured slog fields**

```go
// CORRECT
logger.Info("conversion complete",
    slog.String("file", path),
    slog.Duration("elapsed", elapsed),
)

// INCORRECT - Unstructured string
logger.Info(fmt.Sprintf("converted %s in %v", path, elapsed))
```

### Rule 5: Validation Rule
**All parsers MUST validate input inline and fail fast**

```go
// CORRECT - Validate immediately
func Parse(data []byte) (*UniversalRecipe, error) {
    if len(data) < 1024 {
        return nil, fmt.Errorf("invalid file: too small")
    }
    // Continue parsing
}

// INCORRECT - Parse first, validate later
func Parse(data []byte) (*UniversalRecipe, error) {
    recipe := &UniversalRecipe{}
    // Parse entire file
    return validate(recipe)  // Validation deferred
}
```

### Rule 6: Package Structure Rule
**All format packages MUST follow identical structure: parse.go, generate.go, {format}_test.go**

```
// CORRECT
internal/formats/np3/
├── parse.go
├── generate.go
└── np3_test.go

// INCORRECT - Mixed structure
internal/formats/np3/
├── np3.go           # Combined parse/generate
└── test.go          # Wrong test name
```

### Rule 7: WASM Export Rule
**All WASM exports MUST use `go:wasmexport` and return (ptr, len, error) tuple**

```go
// CORRECT
//go:wasmexport convertPreset
func convertPreset(inputPtr, inputLen uint32, src, dst string) (uint32, uint32, string)

// INCORRECT - Wrong signature
func convertPreset(input []byte, src, dst string) ([]byte, error)
```

### Rule 8: Naming Convention Rule
**All package names MUST be lowercase singular, all types MUST be CamelCase**

```go
// CORRECT
package converter
type UniversalRecipe struct { }
func Parse(data []byte) (*UniversalRecipe, error)

// INCORRECT
package Converter            // CamelCase package
type universal_recipe struct // Snake case type
func parse()                 # Unexported when should be public
```

---

## Data Architecture

### UniversalRecipe Structure

**Purpose**: Central data model for all conversions. Supports all parameters from all formats.

**Complete Structure:**
```go
// internal/model/recipe.go
package model

type UniversalRecipe struct {
    // Metadata
    Name         string `json:"name"`
    SourceFormat string `json:"source_format"`  // "np3", "xmp"

    // Basic Adjustments
    Exposure     float64 `json:"exposure"`       // -5.0 to +5.0
    Contrast     int     `json:"contrast"`       // -100 to +100
    Highlights   int     `json:"highlights"`     // -100 to +100
    Shadows      int     `json:"shadows"`        // -100 to +100
    Whites       int     `json:"whites"`         // -100 to +100
    Blacks       int     `json:"blacks"`         // -100 to +100

    // Color Adjustments
    Saturation   int     `json:"saturation"`     // -100 to +100
    Vibrance     int     `json:"vibrance"`       // -100 to +100

    // Clarity & Sharpness
    Clarity      int     `json:"clarity"`        // -100 to +100
    Sharpness    int     `json:"sharpness"`      // 0 to 150

    // HSL Color Adjustments (8 colors)
    Red          ColorAdjustment `json:"red"`
    Orange       ColorAdjustment `json:"orange"`
    Yellow       ColorAdjustment `json:"yellow"`
    Green        ColorAdjustment `json:"green"`
    Aqua         ColorAdjustment `json:"aqua"`
    Blue         ColorAdjustment `json:"blue"`
    Purple       ColorAdjustment `json:"purple"`
    Magenta      ColorAdjustment `json:"magenta"`

    // Temperature & Tint
    Temperature  int     `json:"temperature"`    // -100 to +100
    Tint         int     `json:"tint"`           // -100 to +100

    // Tone Curve
    ToneCurve    []Point `json:"tone_curve,omitempty"`

    // Split Toning
    SplitShadowHue        int `json:"split_shadow_hue"`        // 0 to 360
    SplitShadowSaturation int `json:"split_shadow_saturation"` // 0 to 100
    SplitHighlightHue     int `json:"split_highlight_hue"`     // 0 to 360
    SplitHighlightSaturation int `json:"split_highlight_saturation"` // 0 to 100

    // Metadata for unknown fields
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type ColorAdjustment struct {
    Hue        int `json:"hue"`         // -100 to +100
    Saturation int `json:"saturation"`  // -100 to +100
    Luminance  int `json:"luminance"`   // -100 to +100
}

type Point struct {
    X int `json:"x"`  // 0-255
    Y int `json:"y"`  // 0-255
}
```

**Parameter Mapping Strategy:**

| NP3 Parameter | XMP Parameter | Range | Offset |
|---------------|---------------|-------|--------|
| Contrast (-100 to +100) | crs:Contrast2012 | -100 to +100 | 0x110 |
| Highlights (-100 to +100) | crs:Highlights2012 | -100 to +100 | 0x11A |
| Shadows (-100 to +100) | crs:Shadows2012 | -100 to +100 | 0x124 |
| White Level (-100 to +100) | crs:Whites2012 | -100 to +100 | 0x12E |
| Black Level (-100 to +100) | crs:Blacks2012 | -100 to +100 | 0x138 |
| Saturation (±3) | crs:Saturation | -100 to +100 | 0x142 |
| Hue (±9°) | crs:HueAdjustment* | -180 to +180 | - |
| Sharpness (0-9) | crs:Sharpness | 0 to 150 | - |
| Brightness (±1) | crs:Exposure2012 | -5.0 to +5.0 | - |

**Note on Tone Parameters:** The Contrast, Highlights, Shadows, White Level, and Black Level parameters are the foundation of our XMP → NP3 conversion strategy. NP3 supports either these basic tone parameters OR a custom tone curve, but not both simultaneously. We prioritize direct parameter mapping for accuracy and simplicity.

**Metadata Dictionary Usage:**

For fields that don't map 1:1 between formats, use the metadata dictionary:

```go
// Preserve unknown NP3 bytes
recipe.Metadata["np3_unknown_bytes"] = "3A7F..."

// Preserve XMP tone curves when converting to NP3 (for round-trip fidelity)
// NP3 limitation: Cannot use tone curve AND basic parameters simultaneously
// Strategy: Use basic parameters (Highlights/Shadows/etc.), preserve curves in metadata
recipe.Metadata["xmp_point_curves"] = pointCurveJSON
recipe.Metadata["xmp_parametric_curve"] = parametricCurveJSON

// Preserve format-specific fields
recipe.Metadata["format_version"] = "7.0"
```

---

## API Contracts

### Core Conversion API

**Function Signature:**
```go
func Convert(input []byte, sourceFormat, targetFormat string) ([]byte, error)
```

**Parameters:**
- `input`: Raw bytes of source file
- `sourceFormat`: One of "np3", "xmp"
- `targetFormat`: One of "np3", "xmp"

**Returns:**
- `[]byte`: Converted file bytes (ready to write to disk)
- `error`: ConversionError if conversion fails, nil on success

**Error Types:**
```go
type ConversionError struct {
    Operation string  // "parse", "generate", "validate"
    Format    string  // "np3", "xmp"
    Cause     error   // Underlying error
}
```

**Example Usage:**

```go
// CLI usage
input, _ := os.ReadFile("portrait.np3")
output, err := converter.Convert(input, "np3", "xmp")
if err != nil {
    var convErr *converter.ConversionError
    if errors.As(err, &convErr) {
        fmt.Printf("Conversion failed: %s %s: %v\n", convErr.Operation, convErr.Format, convErr.Cause)
    }
}
os.WriteFile("portrait.xmp", output, 0644)

// WASM usage
output, err := converter.Convert(inputBytes, srcFormat, dstFormat)
if err != nil {
    return 0, 0, err.Error()  // Return error to JavaScript
}
```

**Contract Guarantees:**

1. **Idempotency**: Converting A → B → A produces functionally equivalent output to A
2. **Validation**: All inputs validated before conversion, errors returned immediately
3. **Thread Safety**: Stateless API, safe to call from multiple goroutines
4. **No Side Effects**: No file I/O, no global state modification, pure function
5. **Round-Trip Fidelity**: 95%+ accuracy for all critical parameters when converting A → B → A

---

## Security Architecture

### Privacy-First Design

**Principle**: Zero data leaves the user's device

**Web Interface:**
- All conversions happen client-side via WebAssembly
- No network requests during conversion
- No analytics, no telemetry, no tracking
- Files never uploaded to any server

**Implementation:**
```javascript
// web/main.js - NO fetch() calls, NO XMLHttpRequest
async function convertFile(file) {
    // Step 1: Read file locally (FileReader API)
    const bytes = await file.arrayBuffer();

    // Step 2: Convert in WASM (runs in browser)
    const result = await wasmConvert(bytes, file.type, targetFormat);

    // Step 3: Download result (never touches network)
    downloadFile(result, file.name);
}
```

**CLI/TUI:**
- All processing local, no network access
- No telemetry, no crash reporting, no auto-updates
- Files read/written only to user-specified paths

### Input Validation

**Defense Against Malicious Files:**

```go
func Parse(data []byte) (*UniversalRecipe, error) {
    // Step 1: Validate file size (prevent memory exhaustion)
    if len(data) > 10*1024*1024 {  // 10 MB max
        return nil, fmt.Errorf("file too large: %d bytes (max 10 MB)", len(data))
    }

    // Step 2: Validate magic bytes (prevent wrong format)
    if len(data) < 2 || string(data[0:2]) != "NP" {
        return nil, fmt.Errorf("invalid magic bytes")
    }

    // Step 3: Validate structure before parsing
    if len(data) < 1024 {
        return nil, fmt.Errorf("file too small, expected >= 1024 bytes")
    }

    // Step 4: Validate parameter ranges during parsing
    contrast := int(data[66])
    if contrast < 0 || contrast > 255 {
        return nil, fmt.Errorf("invalid contrast byte: %d", contrast)
    }

    return &UniversalRecipe{}, nil
}
```

### WASM Security

**Memory Safety:**
- Go's memory safety prevents buffer overflows
- No unsafe pointer arithmetic exposed to JavaScript
- Bounds checking on all array accesses

**Sandboxing:**
- WASM runs in browser sandbox
- No access to filesystem, network, or OS
- Limited to allocated memory only

---

## Performance Considerations

### Performance Goals

| Operation | Target | Measurement |
|-----------|--------|-------------|
| Single file conversion (Web) | <100ms | JavaScript performance.now() |
| Single file conversion (CLI) | <20ms | Go benchmark |
| Batch 100 files (CLI) | <2s | time command |
| WASM binary size | <2 MB | wasm-opt --strip |
| Web page load | <500ms | Lighthouse |

### Optimization Strategies

**1. WASM Optimization:**

```bash
# Build with optimizations
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go

# Strip debug symbols (reduces 50%)
wasm-opt --strip web/recipe.wasm -o web/recipe.wasm

# Compress (gzip reduces 70%)
gzip -9 web/recipe.wasm
```

**Expected Results:**
- Raw WASM: ~3-5 MB
- Stripped: ~1.5-2.5 MB
- Gzipped: ~500 KB - 800 KB

**2. Parallel Batch Processing (CLI/TUI):**

```go
func ConvertBatch(files []string, targetFormat string) error {
    numWorkers := runtime.NumCPU()  // Use all CPU cores
    jobs := make(chan string, len(files))
    results := make(chan error, len(files))

    // Start worker pool
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for file := range jobs {
                err := convertSingle(file, targetFormat)
                results <- err
            }
        }()
    }

    // Send jobs
    for _, file := range files {
        jobs <- file
    }
    close(jobs)

    wg.Wait()
    close(results)

    // Collect errors
    for err := range results {
        if err != nil {
            return err
        }
    }

    return nil
}
```

**Expected Speedup:**
- 8-core CPU: 100 files in ~500ms (vs 2s sequential)
- Scales linearly with CPU count

**3. Lazy Loading (TUI):**

```go
type Model struct {
    files        []FileInfo   // Metadata only (name, size, format)
    previewCache map[int]*UniversalRecipe  // Cache parsed previews
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.KeyMsg:
        if msg.String() == "down" {
            m.selected++

            // Lazy load preview for selected file
            if _, cached := m.previewCache[m.selected]; !cached {
                return m, loadPreview(m.files[m.selected])
            }
        }
}
```

**4. Benchmarking:**

```go
func BenchmarkConvert_NP3_to_XMP(b *testing.B) {
    input, _ := os.ReadFile("testdata/np3/portrait.np3")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := converter.Convert(input, "np3", "xmp")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**Run Benchmarks:**
```bash
go test -bench=. -benchmem ./internal/converter/
```

**Expected Output:**
```
BenchmarkConvert_NP3_to_XMP-8    50000    25000 ns/op    4096 B/op    12 allocs/op
```

This indicates:
- ~25µs (0.025ms) per conversion
- Well under 100ms goal
- Minimal memory allocation

---

## Deployment Architecture

### Web Interface (Cloudflare Pages)

**Deployment Flow:**

```
GitHub Push (main branch)
    ↓
Cloudflare Pages Trigger
    ↓
Build (optional) - None for static HTML/JS/WASM
    ↓
Deploy to CDN
    ↓
Live at https://recipe.pages.dev
```

**Configuration:**

```yaml
# .github/workflows/deploy-pages.yml (optional, for build step)
name: Deploy to Cloudflare Pages

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build WASM
        run: |
          GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go

      - name: Deploy to Cloudflare Pages
        uses: cloudflare/pages-action@v1
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          projectName: recipe
          directory: web
```

**Directory Structure for Deployment:**

```
web/
├── index.html       # Main page
├── main.js          # WASM loader and UI
├── style.css        # Styling
└── recipe.wasm      # Compiled WASM binary
```

**Cloudflare Pages automatically serves**:
- index.html at https://recipe.pages.dev/
- Gzip compression enabled (WASM reduced 70%)
- Global CDN (sub-100ms latency worldwide)
- HTTPS by default

---

### CLI Distribution (GitHub Releases)

**Build Matrix:**

```yaml
# .github/workflows/release.yml
name: Build Release Binaries

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build
        run: |
          GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -o recipe-${{ matrix.os }}-${{ matrix.arch }} cmd/cli/main.go

      - name: Upload to GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: recipe-*
```

**Release Artifacts:**
- `recipe-linux-amd64` (Linux x86_64)
- `recipe-linux-arm64` (Linux ARM64)
- `recipe-darwin-amd64` (macOS Intel)
- `recipe-darwin-arm64` (macOS Apple Silicon)
- `recipe-windows-amd64.exe` (Windows x64)

**Installation:**

```bash
# macOS (Apple Silicon)
curl -L https://github.com/user/recipe/releases/latest/download/recipe-darwin-arm64 -o recipe
chmod +x recipe
sudo mv recipe /usr/local/bin/

# Linux
curl -L https://github.com/user/recipe/releases/latest/download/recipe-linux-amd64 -o recipe
chmod +x recipe
sudo mv recipe /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/user/recipe/releases/latest/download/recipe-windows-amd64.exe -OutFile recipe.exe
Move-Item recipe.exe C:\Windows\System32\
```

---

## Development Environment

### Prerequisites

**Required:**
- Go 1.24+ (`go version` to verify)
- Git
- Make (optional but recommended)

**Recommended:**
- VS Code with Go extension
- gofmt configured on save

**Installation:**

```bash
# macOS
brew install go
brew install make

# Linux
sudo apt install golang make

# Windows
choco install golang make
```

### Setup

**Clone Repository:**
```bash
git clone https://github.com/user/recipe.git
cd recipe
```

**Verify Go Version:**
```bash
go version
# Should output: go version go1.24.0 or higher
```

**Install Cobra CLI (optional):**
```bash
go install github.com/spf13/cobra-cli@latest
```

### Development Workflow

**Build CLI:**
```bash
make cli
# or
go build -o recipe cmd/cli/main.go
./recipe --help
```

**Build TUI:**
```bash
make tui
# or
go build -o recipe-tui cmd/tui/main.go
./recipe-tui
```

**Build WASM:**
```bash
make wasm
# or
GOOS=js GOARCH=wasm go build -o web/recipe.wasm cmd/wasm/main.go
```

**Run Tests:**
```bash
make test
# or
go test ./...

# With coverage
go test -cover ./...

# With race detector
go test -race ./...
```

**Run Benchmarks:**
```bash
make bench
# or
go test -bench=. -benchmem ./internal/converter/
```

**Format Code:**
```bash
make fmt
# or
gofmt -s -w .
go vet ./...
```

**Serve Web Interface Locally:**
```bash
# Simple HTTP server for testing WASM
cd web
python3 -m http.server 8080
# Open http://localhost:8080
```

### Makefile

**Complete Makefile:**

```makefile
.PHONY: all cli tui wasm test bench fmt clean

# Default target
all: cli tui wasm

# Build CLI
cli:
	go build -o recipe cmd/cli/main.go

# Build TUI
tui:
	go build -o recipe-tui cmd/tui/main.go

# Build WASM
wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go

# Run tests
test:
	go test -v -cover ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem ./internal/converter/

# Format code
fmt:
	gofmt -s -w .
	go vet ./...

# Clean build artifacts
clean:
	rm -f recipe recipe-tui web/recipe.wasm

# Install CLI to system
install: cli
	sudo mv recipe /usr/local/bin/
```

**Usage:**
```bash
make              # Build all interfaces
make test         # Run tests
make bench        # Run benchmarks
make fmt          # Format code
make install      # Install CLI
make clean        # Remove build artifacts
```

---

## Architecture Decision Records (ADRs)

### ADR-001: Go 1.24+ for Enhanced WASM Support

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need Go version with best WebAssembly support for <100ms conversion goal

**Decision**: Use Go 1.24.0+ as minimum required version

**Rationale**:
- `go:wasmexport` directive enables direct memory access (zero reflection overhead)
- 20-30% smaller WASM binary size vs Go 1.23
- Enhanced type support for WASM exports (better TypeScript interop)
- slog built into stdlib (Go 1.21+) for structured logging

**Consequences**:
- ✅ Faster WASM performance (<100ms goal achievable)
- ✅ Smaller binary size (better download time)
- ✅ Better JavaScript/TypeScript integration
- ⚠️ Requires developers use Go 1.24+ (released Feb 2025)

**Validation**: Web search confirmed Go 1.24 released February 11, 2025 with go:wasmexport feature

---

### ADR-002: Vanilla JavaScript (No Frameworks)

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need frontend technology for Web interface with minimal build complexity

**Decision**: Use vanilla JavaScript (ES6+) with no frameworks

**Rationale**:
- Zero build step (no webpack, no npm, no node_modules)
- 90%+ browser support for FileReader API and WebAssembly (since 2017)
- Direct DOM manipulation faster than framework diffing for simple UI
- Simpler deployment (just HTML/JS/CSS files)
- Consistent with project philosophy (zero dependencies)

**Consequences**:
- ✅ No build tooling complexity
- ✅ Faster iteration (edit file, refresh browser)
- ✅ Simpler deployment (static files only)
- ⚠️ Manual DOM updates required
- ⚠️ No type checking (can add TypeScript later if needed)

**Alternatives Considered**:
- React: Too heavy, requires build step, overkill for simple UI
- Svelte: Better than React but still requires build step
- Vue: Better than React but still requires build step

---

### ADR-003: Hub-and-Spoke Architecture

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need conversion architecture that scales with number of formats

**Decision**: Use hub-and-spoke pattern with UniversalRecipe intermediate representation

**Rationale**:
- Reduces converter complexity from O(N²) to O(N)
- Easy to add new formats (1 parser + 1 generator)
- Single source of truth for parameter mappings
- Centralized validation and normalization

**Consequences**:
- ✅ Adding 4th format requires 2 functions (not 6 converters)
- ✅ Parameter mapping logic in one place
- ✅ Easy to test (each format independently)
- ⚠️ Extra conversion step (format → universal → format)
- ⚠️ Must maintain comprehensive UniversalRecipe (50+ fields)

**Mitigation**:
- Metadata dictionary for unknown fields (prevents data loss)
- Round-trip tests validate fidelity (95%+ accuracy goal)

---

### ADR-004: Cloudflare Pages for Hosting

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need simple, fast, free hosting for Web interface

**Decision**: Use Cloudflare Pages for static site hosting

**Rationale**:
- Zero cost (free tier sufficient)
- Auto-deploy on push to main branch
- Global CDN (sub-100ms latency worldwide)
- Gzip compression enabled by default (70% WASM size reduction)
- HTTPS by default
- No configuration needed (just connect GitHub repo)

**Consequences**:
- ✅ Zero hosting cost
- ✅ Fast deployment (push to main → live in <1 min)
- ✅ Global performance
- ⚠️ Locked to Cloudflare (but can migrate to any static host)

**Alternatives Considered**:
- GitHub Pages: Similar but slower CDN, less compression
- Netlify: Similar but free tier more restrictive
- Vercel: Similar but optimized for Next.js (overkill)

---

### ADR-005: Table-Driven Tests with Real Sample Files

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need testing strategy that validates 95%+ accuracy goal

**Decision**: Use table-driven tests with 1,501 real sample files from testdata/

**Rationale**:
- Real-world validation (not synthetic test data)
- Comprehensive coverage (73 NP3, 914 XMP)
- Idiomatic Go testing pattern
- Round-trip tests ensure fidelity
- Catches edge cases synthetic data misses

**Consequences**:
- ✅ High confidence in conversion accuracy
- ✅ Real-world validation
- ✅ Easy to add new test files (drop in testdata/)
- ⚠️ Slower test runs (987 files to process)
- ⚠️ Large testdata/ directory

**Optimization**:
- Run quick tests by default (`go test`)
- Run full suite in CI (`go test ./...`)
- Use subtests for granular failure reporting

---

### ADR-006: Wrapped Errors with ConversionError Type

**Date**: 2025-11-03
**Status**: Accepted
**Context**: Need error handling that provides format-specific context

**Decision**: Use custom `ConversionError` type that wraps stdlib errors

**Rationale**:
- Type-safe error checking (can use errors.As)
- Format-specific context (operation, format, cause)
- No external dependencies (stdlib errors package)
- Implements error interface and Unwrap()

**Consequences**:
- ✅ Clear error messages ("parse np3 failed: invalid magic bytes")
- ✅ Type-safe error handling
- ✅ Zero dependencies
- ⚠️ Slightly more verbose error creation

**Example**:
```go
return &ConversionError{
    Operation: "parse",
    Format:    "np3",
    Cause:     fmt.Errorf("invalid magic bytes"),
}
```

---

**End of Architecture Document**
