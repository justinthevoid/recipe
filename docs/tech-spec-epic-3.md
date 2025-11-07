# Epic Technical Specification: CLI Interface

Date: 2025-11-06
Author: Justin
Epic ID: 3
Status: Draft

---

## Overview

Epic 3 delivers a professional command-line interface (CLI) for Recipe, enabling automated preset conversion workflows for power users, batch processing scenarios, and integration into photography pipelines. Built with the industry-standard Cobra framework, this CLI provides Unix-convention commands for single and batch file conversions with intelligent format detection, progress tracking, and machine-readable output.

The CLI reuses the proven conversion engine from Epic 1 (`internal/converter`), ensuring 95%+ accuracy and consistency across all Recipe interfaces (Web, TUI, CLI). It targets photographers who need scripting capabilities, automation workflows, and high-performance batch operations that would be impractical through the web interface.

## Objectives and Scope

**In Scope:**
- Cobra-based CLI application with convert command (`recipe convert INPUT --to FORMAT`)
- Single file conversion with auto-format detection and explicit target specification
- Batch processing with parallel CPU utilization (`recipe convert --batch *.xmp --to np3`)
- Progress indicators for batch operations (file count, success/error tracking)
- Verbose logging mode for debugging (`--verbose` flag outputs to stderr)
- JSON output mode for scripting integration (`--json` flag outputs structured data)
- Cross-platform single-binary distribution (Windows, macOS, Linux for amd64 and arm64)
- Exit codes for scripting (0=success, 1=error)
- Help text and version information (`--help`, `--version`)

**Out of Scope:**
- Interactive prompts (defer to TUI in Epic 4)
- Visual parameter preview (defer to TUI in Epic 4)
- Real-time file watching / daemon mode
- Network-based conversion (all processing remains local)
- Configuration file management (all options via CLI flags)
- Plugin system or custom format adapters
- GUI wrapper around CLI (separate project if needed)

## System Architecture Alignment

The CLI leverages the existing hub-and-spoke architecture from Epic 1, calling `converter.Convert()` as the single source of truth for all conversion logic. This ensures format consistency and validation rules are identical across Web, TUI, and CLI interfaces.

**Architecture References:**
- **Shared Library:** `internal/converter/converter.go` - Stateless Convert() function used by all interfaces
- **Format Parsers:** `internal/formats/{np3,xmp,lrtemplate}/parse.go` - Already implemented and tested in Epic 1
- **Format Generators:** `internal/formats/{np3,xmp,lrtemplate}/generate.go` - Already implemented and tested in Epic 1
- **CLI Framework:** Cobra (github.com/spf13/cobra) - Standard Go CLI pattern used by kubectl, hugo, gh
- **Project Structure:** `cmd/cli/main.go` as entry point, following standard Go project layout

**Key Constraints from Architecture:**
- Zero OS dependencies in converter package (enables WASM compatibility maintained from Epic 2)
- All file I/O happens in CLI layer (`cmd/cli/`), not in converter
- Errors wrapped as `ConversionError` type with operation/format context
- Logging via `slog` with structured fields (Go 1.21+ stdlib)
- Table-driven tests with real sample files from `testdata/` (1,501 files)

## Detailed Design

### Services and Modules

| Module                                    | Responsibility                            | Inputs                            | Outputs                                    | Owner           |
| ----------------------------------------- | ----------------------------------------- | --------------------------------- | ------------------------------------------ | --------------- |
| **cmd/cli/main.go**                       | CLI entry point, Cobra app initialization | Command-line args                 | Exit code (0 or 1)                         | Epic 3          |
| **cmd/cli/convert.go**                    | Convert command implementation            | File path(s), format flags        | Converted files on disk                    | Story 3-2       |
| **cmd/cli/root.go**                       | Root command, global flags, version info  | CLI flags                         | Cobra root command                         | Story 3-1       |
| **cmd/cli/batch.go**                      | Batch processing orchestration            | File glob patterns                | Multiple converted files                   | Story 3-3       |
| **cmd/cli/format.go**                     | Format detection utilities                | File path or bytes                | Format string ("np3", "xmp", "lrtemplate") | Story 3-4       |
| **cmd/cli/output.go**                     | Output formatting (normal, verbose, JSON) | Conversion result                 | Formatted output to stdout/stderr          | Story 3-5, 3-6  |
| **internal/converter**                    | **Shared conversion engine (Epic 1)**     | Byte array, source/target formats | Converted byte array or error              | Epic 1 (reused) |
| **internal/formats/{np3,xmp,lrtemplate}** | **Format parsers/generators (Epic 1)**    | Format-specific bytes             | UniversalRecipe or format bytes            | Epic 1 (reused) |

**Key Design Principles:**
- **Thin CLI Layer:** All business logic in `internal/converter`, CLI only handles I/O and formatting
- **Cobra Command Hierarchy:** Root command with subcommands (convert, batch future: inspect, diff)
- **Stateless Operations:** Each convert call is independent, no global state
- **Parallel Batch Processing:** Use goroutines + worker pool pattern for batch operations
- **Structured Logging:** slog with JSON output in verbose mode, minimal output in normal mode

### Data Models and Contracts

**CLI does not introduce new data models.** All conversion data models are defined in Epic 1:

| Model                | Definition                     | Usage in CLI                                               |
| -------------------- | ------------------------------ | ---------------------------------------------------------- |
| **UniversalRecipe**  | `internal/model/recipe.go`     | Intermediate format (not exposed to CLI user)              |
| **ConversionError**  | `internal/converter/errors.go` | Wrapped errors returned to CLI, formatted for user display |
| **ConversionResult** | New struct for CLI (below)     | Wraps output bytes + metadata for display                  |

**New CLI-Specific Model:**

```go
// cmd/cli/types.go
package main

type ConversionResult struct {
    InputFile    string        `json:"input"`
    OutputFile   string        `json:"output"`
    SourceFormat string        `json:"source_format"`
    TargetFormat string        `json:"target_format"`
    Success      bool          `json:"success"`
    Error        string        `json:"error,omitempty"`
    Duration     time.Duration `json:"duration_ms"`
    FileSize     int64         `json:"file_size_bytes"`
    Warnings     []string      `json:"warnings,omitempty"`
}

type BatchResult struct {
    TotalFiles     int                 `json:"total_files"`
    SuccessCount   int                 `json:"success_count"`
    ErrorCount     int                 `json:"error_count"`
    TotalDuration  time.Duration       `json:"total_duration_ms"`
    Results        []ConversionResult  `json:"results"`
}
```

**Flag Structures:**

```go
// Command flags (managed by Cobra)
type ConvertFlags struct {
    From       string  // --from, -f (source format, optional if auto-detect)
    To         string  // --to, -t (target format, required)
    Output     string  // --output, -o (output file path, optional)
    Verbose    bool    // --verbose, -v
    JSON       bool    // --json
    Overwrite  bool    // --overwrite (default: false, fail if output exists)
}

type BatchFlags struct {
    To         string  // --to, -t (target format, required)
    OutputDir  string  // --output-dir (optional, default: same directory)
    Parallel   int     // --parallel (number of workers, default: NumCPU)
    Verbose    bool    // --verbose, -v
    JSON       bool    // --json
    ContinueOnError bool  // --continue-on-error (default: true)
}
```

### APIs and Interfaces

**CLI-to-Converter API (Existing from Epic 1):**

```go
// internal/converter/converter.go (already implemented)
func Convert(input []byte, sourceFormat, targetFormat string) ([]byte, error)
```

**New CLI Internal APIs:**

```go
// cmd/cli/convert.go
func convertSingleFile(inputPath, outputPath, sourceFormat, targetFormat string, flags ConvertFlags) (*ConversionResult, error)
// Orchestrates: Read file → Detect format → Call converter.Convert() → Write output → Return result

// cmd/cli/batch.go
func convertBatch(inputFiles []string, targetFormat string, flags BatchFlags) (*BatchResult, error)
// Orchestrates: Worker pool → Parallel conversions → Progress tracking → Aggregate results

// cmd/cli/format.go
func detectFormat(filePath string) (string, error)
// Returns: "np3", "xmp", "lrtemplate" based on file extension
// Fallback: Read magic bytes if extension ambiguous

func detectFormatFromBytes(data []byte) (string, error)
// Returns: Format based on content inspection (magic bytes, XML root element)

// cmd/cli/output.go
func formatResult(result *ConversionResult, flags ConvertFlags) string
// Normal mode: "✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)"
// Verbose mode: Add parameter details, warnings
// JSON mode: Marshal result to JSON

func formatBatchResult(result *BatchResult, flags BatchFlags) string
// Normal mode: "✓ Converted 100 files: 98 success, 2 errors (1.5s total)"
// Verbose mode: List each file result
// JSON mode: Marshal full batch result to JSON
```

**Cobra Command Signatures:**

```go
// cmd/cli/root.go
var rootCmd = &cobra.Command{
    Use:   "recipe",
    Short: "Convert photo presets between formats",
    Long:  "Recipe is a universal photo preset converter...",
}

// cmd/cli/convert.go
var convertCmd = &cobra.Command{
    Use:   "convert [input] [output]",
    Short: "Convert a preset file between formats",
    Args:  cobra.MinimumNArgs(1),
    RunE:  runConvert,
}

// cmd/cli/batch.go (future story if batch separated from convert)
var batchCmd = &cobra.Command{
    Use:   "batch [pattern]",
    Short: "Convert multiple files in batch",
    Args:  cobra.ExactArgs(1),
    RunE:  runBatch,
}
```

**Exit Codes:**

| Code | Meaning     | When                                                        |
| ---- | ----------- | ----------------------------------------------------------- |
| 0    | Success     | All conversions succeeded                                   |
| 1    | Error       | Any conversion failed (unless --continue-on-error in batch) |
| 2    | Usage Error | Invalid flags, missing required args                        |

### Workflows and Sequencing

**Single File Conversion Flow:**

```
User Command: recipe convert portrait.xmp --to np3
    ↓
1. Parse CLI flags (Cobra)
    ↓
2. Validate args (input file exists, target format valid)
    ↓
3. Detect source format (from extension or --from flag)
    ↓
4. Read input file (os.ReadFile)
    ↓
5. Call converter.Convert(bytes, "xmp", "np3")
    ↓
6. Write output file (os.WriteFile)
    ↓
7. Display result: "✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)"
    ↓
8. Exit code 0
```

**Batch Conversion Flow:**

```
User Command: recipe convert --batch *.xmp --to np3
    ↓
1. Parse CLI flags (Cobra)
    ↓
2. Glob pattern to file list (filepath.Glob)
    ↓
3. Create worker pool (N = runtime.NumCPU() or --parallel flag)
    ↓
4. Start progress tracking goroutine
    ↓
5. Distribute files to workers via channel
    ↓
For each worker (parallel):
    6a. Read input file
    6b. Detect format
    6c. Call converter.Convert()
    6d. Write output file
    6e. Send result to results channel
    6f. Update progress counter (atomic)
    ↓
7. Aggregate results (success/error counts)
    ↓
8. Display summary: "✓ Converted 100 files: 98 success, 2 errors (1.5s total)"
    ↓
9. Exit code 0 (if --continue-on-error) or 1 (if any errors and default)
```

**Verbose Mode Sequence:**

```
Normal: ✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)

Verbose (-v flag):
[INFO] Reading input: portrait.xmp
[INFO] Detected format: xmp
[INFO] Parsing XMP file...
[INFO] Extracted parameters: Exposure=+0.5, Contrast=+15, Saturation=-10
[INFO] Converting xmp → np3...
[WARN] Parameter 'Grain' not supported in NP3 format (omitted)
[INFO] Generating NP3 binary...
[INFO] Writing output: portrait.np3
✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)
```

**JSON Mode Sequence:**

```
Normal output: Human-readable text

JSON output (--json flag):
{
  "input": "portrait.xmp",
  "output": "portrait.np3",
  "source_format": "xmp",
  "target_format": "np3",
  "success": true,
  "duration_ms": 15,
  "file_size_bytes": 1234,
  "warnings": ["Parameter 'Grain' not supported in NP3 format (omitted)"]
}
```

**Error Handling Sequence:**

```
If file not found:
    → Error: "Input file not found: portrait.xmp"
    → Exit code 1

If format detection fails:
    → Error: "Unknown file format: portrait.unknown (expected .np3, .xmp, .lrtemplate)"
    → Exit code 1

If conversion fails (Epic 1 error):
    → Error: "Conversion failed: parse xmp error: invalid XML structure"
    → Exit code 1

If output file exists (no --overwrite):
    → Error: "Output file already exists: portrait.np3 (use --overwrite to replace)"
    → Exit code 1
```

## Non-Functional Requirements

### Performance

**CLI Execution Speed:**
- **Single File Conversion:** Target <20ms per file (10-100x faster than Python v1.0 baseline of 50-200ms)
- **Batch Processing:** Convert 100 files in <2 seconds using parallel processing (vs 5-20s Python)
- **Cold Start:** CLI initialization <100ms (single binary, no interpreter startup)
- **Memory Efficiency:** CLI process uses <100MB for 1,000 file batch operations
- **Measurement:** Go benchmarks (`go test -bench`) and time command validation

**Command Responsiveness:**
- Help text display: <10ms
- Flag parsing: <5ms
- Error messages: Instant display (<1ms)
- Progress indicators: Update every 100ms during batch operations

**Format Detection:**
- File extension detection: <1ms
- Content-based detection (fallback): <5ms
- No user-perceptible delay in workflow

**Output Generation:**
- JSON output formatting: <10ms additional overhead
- Verbose logging: <15% performance impact
- Normal output: <5ms formatting time

**Optimization Strategies:**
- Worker pool pattern for batch processing (utilize all CPU cores)
- Parallel file I/O where possible
- Minimal allocations in hot path (converter.Convert() core)
- Pre-compiled binary (zero JIT warm-up)

### Security

**Local Processing Only:**
- Zero network requests during conversion
- No telemetry, analytics, or crash reporting
- No auto-update mechanisms (user-controlled updates)
- All processing happens on local machine

**Input Validation:**
- File size limits: Max 10MB per file (prevent memory exhaustion)
- Format validation: Magic bytes verification before parsing
- Path traversal prevention: Validate all file paths
- Buffer overflow protection: Go's memory safety model

**File System Access:**
- Read-only access to input files
- Write-only access to specified output paths
- No directory traversal beyond user-specified paths
- Fail-safe: Refuse overwrite without --overwrite flag

**Dependency Security:**
- Cobra framework: Well-maintained, security-audited
- Zero custom crypto: No encryption/decryption capabilities
- Go stdlib only for core conversion: Minimal attack surface
- Vendored dependencies: Reproducible builds with go.mod/go.sum

**Error Information Disclosure:**
- Error messages don't reveal system paths (sanitize)
- No stack traces in production builds
- Verbose mode: Additional details only when explicitly requested
- JSON output: No internal state leakage

### Reliability/Availability

**Conversion Accuracy:**
- 95%+ parameter fidelity target (validated via round-trip testing)
- Graceful handling of unmappable parameters (user warnings)
- Consistent results across platforms (Windows, macOS, Linux)
- Deterministic output (same input always produces same output)

**Error Handling:**
- Clear error messages for all failure modes
- Actionable suggestions (e.g., "Use --overwrite to replace existing file")
- No silent failures (all errors reported)
- Exit codes: 0=success, 1=error, 2=usage error

**Cross-Platform Consistency:**
- Identical behavior on Windows, macOS, Linux
- Same conversion results regardless of platform
- Path handling: Respect OS-specific conventions (/ vs \)
- Line endings: Preserve source format conventions

**Batch Processing Resilience:**
- Continue-on-error mode: Process all files even if some fail (default)
- Fail-fast mode: Stop on first error (optional flag)
- Summary report: Count of successes/failures
- Individual file errors logged with context

**File Integrity:**
- Atomic writes: Use temp files, rename on success
- No partial writes on error
- Original files never modified
- Output validation before write

**Recovery:**
- Stateless operation: No persistent state to corrupt
- Retry-safe: Re-running command produces same result
- Crash-safe: No cleanup needed if process killed

### Observability

**Logging Framework:**
- Structured logging with Go stdlib `slog` (Go 1.21+)
- Log levels: Debug, Info, Warn, Error
- Machine-readable JSON output in verbose mode
- Human-readable output in normal mode

**Normal Mode Output:**
```
✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)
```

**Verbose Mode Output (-v flag):**
```
[INFO] Reading input: portrait.xmp
[INFO] Detected format: xmp
[INFO] Parsing XMP file...
[INFO] Extracted parameters: Exposure=+0.5, Contrast=+15, Saturation=-10
[INFO] Converting xmp → np3...
[WARN] Parameter 'Grain' not supported in NP3 format (omitted)
[INFO] Generating NP3 binary...
[INFO] Writing output: portrait.np3
✓ Converted portrait.xmp → portrait.np3 (1.2 KB, 15ms)
```

**JSON Mode Output (--json flag):**
```json
{
  "input": "portrait.xmp",
  "output": "portrait.np3",
  "source_format": "xmp",
  "target_format": "np3",
  "success": true,
  "duration_ms": 15,
  "file_size_bytes": 1234,
  "warnings": ["Parameter 'Grain' not supported in NP3 format (omitted)"]
}
```

**Batch Progress Tracking:**
- Real-time progress: "Processing 45/100 files..."
- Success/error counters updated live
- Final summary: "✓ Converted 100 files: 98 success, 2 errors (1.5s total)"

**Performance Metrics:**
- Duration tracking per file (milliseconds)
- Total batch duration
- File size reporting
- Available in JSON output for scripting analysis

**Debugging Support:**
- Verbose mode shows parameter extraction details
- Warnings for approximated/omitted parameters
- Error messages include format-specific context
- Stack traces available in debug builds (not production)

**Scripting Integration:**
- JSON output parsable by jq, Python json module
- Exit codes for shell scripting
- Structured errors in JSON mode
- Machine-readable warnings array

## Dependencies and Integrations

### External Dependencies

**Cobra CLI Framework:**
- **Package:** `github.com/spf13/cobra`
- **Version:** v1.8+ (check go.mod for exact version)
- **Purpose:** Command-line interface structure, flag parsing, help generation
- **Rationale:** Industry standard (used by kubectl, hugo, gh), active maintenance
- **Installation:** `go get -u github.com/spf13/cobra@latest`

**Go Standard Library (No Version Constraints):**
- `os` - File I/O operations
- `path/filepath` - Cross-platform path handling, glob pattern matching
- `errors` - Error wrapping and unwrapping
- `log/slog` - Structured logging (Go 1.21+)
- `encoding/json` - JSON marshaling for --json output
- `time` - Duration tracking, performance metrics
- `runtime` - CPU count for parallel processing (runtime.NumCPU())
- `sync` - Goroutine synchronization (worker pools)

### Internal Dependencies (Shared with All Interfaces)

**Core Conversion Engine:**
- **Package:** `github.com/justin/recipe/internal/converter`
- **API:** `Convert([]byte, string, string) ([]byte, error)`
- **Shared By:** CLI (Epic 3), Web/WASM (Epic 2), TUI (Epic 4)
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
- **Read Operations:** `os.ReadFile()` for input files
- **Write Operations:** `os.WriteFile()` for output files
- **Directory Operations:** `os.MkdirAll()` for output directories (if needed)
- **Glob Patterns:** `filepath.Glob()` for batch file selection

**Terminal/Shell:**
- **Standard Output:** `os.Stdout` for normal output
- **Standard Error:** `os.Stderr` for verbose logging and errors
- **Exit Codes:** `os.Exit(0)` for success, `os.Exit(1)` for errors

**Conversion Engine (Critical Integration):**
- **Interface Contract:**
  ```go
  func Convert(input []byte, sourceFormat, targetFormat string) ([]byte, error)
  ```
- **Input:** File bytes read via `os.ReadFile()`
- **Output:** Converted bytes written via `os.WriteFile()`
- **Error Handling:** Wrap conversion errors with file context

**Parallel Processing:**
- **Worker Pool Pattern:** Goroutines + channels for batch operations
- **Concurrency:** `runtime.NumCPU()` workers (default) or user-specified via --parallel flag
- **Synchronization:** `sync.WaitGroup` for worker coordination

### Build Dependencies

**Go Toolchain:**
- **Minimum Version:** Go 1.24.0+
- **Reason:** Enhanced WASM support (shared with Epic 2), slog stdlib
- **Verification:** `go version` must output go1.24.0 or higher

**Build Tools:**
- **Make:** Optional but recommended for build automation
- **Cobra CLI Generator:** `cobra-cli` for scaffolding commands
  - Installation: `go install github.com/spf13/cobra-cli@latest`

### Development Dependencies

**Testing:**
- `testing` (stdlib) - Go test framework
- `testdata/` directory - 1,501 sample files for validation

**Benchmarking:**
- `testing` (stdlib) - Benchmark functions
- `time` (stdlib) - Performance measurement

**Linting/Formatting:**
- `gofmt` (stdlib) - Code formatting
- `go vet` (stdlib) - Static analysis

### Distribution Dependencies

**Cross-Platform Compilation:**
- No runtime dependencies required
- Single binary distribution (static linking)
- Supported platforms: Windows (amd64), macOS (amd64, arm64), Linux (amd64, arm64)

**Package Managers (Optional):**
- Homebrew (macOS/Linux) - For `brew install recipe` distribution
- Scoop (Windows) - For `scoop install recipe` distribution
- Neither required for manual installation (download binary from GitHub Releases)

### Zero Dependencies Policy

**Explicitly NOT Using:**
- No external logging frameworks (using stdlib slog)
- No external CLI frameworks beyond Cobra
- No ORM or database libraries (CLI doesn't persist data)
- No network/HTTP libraries (all processing is local)
- No configuration file parsers (flags only)

**Rationale:**
- Minimal attack surface
- Faster builds
- Simpler maintenance
- Consistent with project philosophy (privacy-first, zero-telemetry)

## Acceptance Criteria (Authoritative)

### AC-1: Basic CLI Command Structure

**Given** a user has Recipe CLI installed
**When** they run `recipe convert portrait.xmp --to np3`
**Then** the CLI successfully converts the file and outputs `portrait.np3` in the same directory
**And** displays success message "✓ Converted portrait.xmp → portrait.np3 (X KB, Xms)"
**And** exits with code 0

**Validation:**
- Test with all format combinations (np3↔xmp, np3↔lrtemplate, xmp↔lrtemplate)
- Verify output file exists and is valid
- Confirm success message displays correctly
- Check exit code equals 0

---

### AC-2: Format Auto-Detection

**Given** a user has a preset file with standard extension (.np3, .xmp, .lrtemplate)
**When** they run `recipe convert INPUT --to FORMAT` without specifying --from
**Then** the CLI correctly detects the source format from file extension
**And** performs conversion without requiring --from flag
**And** displays detected format in verbose output

**Validation:**
- Test all three format extensions
- Confirm no --from flag required
- Verify error message if extension is ambiguous/unknown
- Check verbose mode shows "Detected format: {format}"

---

### AC-3: Batch Processing with Glob Patterns

**Given** a directory contains 100 XMP files
**When** user runs `recipe convert --batch *.xmp --to np3`
**Then** all 100 files are converted in parallel
**And** conversion completes in <2 seconds total
**And** output displays summary: "✓ Converted 100 files: X success, Y errors (Zs total)"
**And** exit code is 0 if all succeed, 1 if any fail (unless --continue-on-error)

**Validation:**
- Benchmark with 100 real sample files from testdata/
- Measure total time with `time` command
- Verify parallel processing (check CPU usage during conversion)
- Test error handling with mixed valid/invalid files

---

### AC-4: Verbose Logging Mode

**Given** a user wants detailed conversion information
**When** they run `recipe convert INPUT --to FORMAT --verbose`
**Then** CLI outputs detailed logs to stderr including:
- File reading confirmation
- Format detection
- Parameter extraction details
- Warnings for unmappable parameters
- Generation and write confirmation
**And** final success message to stdout

**Validation:**
- Compare normal vs verbose output
- Confirm logs go to stderr (success message to stdout)
- Verify structured log format (slog)
- Check warnings appear for known unmappable parameters

---

### AC-5: JSON Output for Scripting

**Given** a user wants machine-readable output
**When** they run `recipe convert INPUT --to FORMAT --json`
**Then** CLI outputs valid JSON to stdout containing:
```json
{
  "input": "portrait.xmp",
  "output": "portrait.np3",
  "source_format": "xmp",
  "target_format": "np3",
  "success": true,
  "duration_ms": 15,
  "file_size_bytes": 1234,
  "warnings": ["..."]
}
```
**And** no human-readable text is mixed with JSON output
**And** JSON is parsable by `jq` and Python's `json.load()`

**Validation:**
- Pipe output to `jq` and verify parsing succeeds
- Use Python to parse: `import json; json.loads(output)`
- Verify all required fields present
- Check no extra text outside JSON object

---

### AC-6: Error Handling - Invalid File

**Given** a user provides an invalid or corrupted preset file
**When** they run `recipe convert INVALID_FILE --to np3`
**Then** CLI displays clear error message:
- "Error: Invalid file format: expected .xmp, .lrtemplate, or .np3"
  OR
- "Error: Conversion failed: parse xmp error: {specific parse error}"
**And** exits with code 1
**And** no output file is created

**Validation:**
- Test with non-preset files (e.g., .txt, .jpg)
- Test with corrupted files (truncated, wrong magic bytes)
- Verify error messages are user-friendly (no stack traces)
- Confirm exit code equals 1

---

### AC-7: Overwrite Protection

**Given** an output file already exists
**When** user runs `recipe convert INPUT --to FORMAT` (without --overwrite)
**Then** CLI displays error: "Error: Output file already exists: {filename} (use --overwrite to replace)"
**And** exits with code 1
**And** existing file is NOT modified

**When** user runs same command WITH `--overwrite` flag
**Then** CLI overwrites existing file
**And** displays success message

**Validation:**
- Create dummy output file before test
- Confirm error without --overwrite flag
- Verify file unchanged after error
- Test --overwrite successfully replaces file

---

### AC-8: Cross-Platform Compatibility

**Given** Recipe CLI is installed on Windows, macOS, and Linux
**When** user converts same input file on all three platforms
**Then** output files are byte-identical across platforms
**And** CLI commands use OS-appropriate path separators (\ on Windows, / on Unix)
**And** line endings are preserved according to source format conventions

**Validation:**
- Build binaries for all platforms (windows-amd64, darwin-amd64, darwin-arm64, linux-amd64)
- Run same conversion on each platform
- SHA256 hash output files and compare
- Test with files containing different line endings

---

### AC-9: Help and Version Information

**Given** a new user wants to learn CLI usage
**When** they run `recipe --help`
**Then** CLI displays comprehensive help text including:
- Available commands (convert, batch, etc.)
- Flag descriptions and examples
- Supported formats
- Exit codes

**When** they run `recipe --version`
**Then** CLI displays version number: "Recipe CLI v{version}"

**Validation:**
- Verify all flags documented in --help output
- Confirm examples are accurate
- Check version matches release tag

---

### AC-10: Performance Targets

**Given** a single preset file (<50KB)
**When** user runs `recipe convert INPUT --to FORMAT`
**Then** conversion completes in <20ms (measured via benchmark)

**Given** 100 preset files
**When** user runs `recipe convert --batch *.xmp --to np3`
**Then** batch conversion completes in <2 seconds total

**Validation:**
- Run `go test -bench` for single file benchmark
- Use `time` command for batch benchmark
- Test on reference hardware (document specs)
- Confirm performance is 10-100x faster than Python v1.0 baseline

---

### AC-11: Exit Code Consistency

**Given** various CLI scenarios
**Then** exit codes follow this contract:
- 0: All conversions successful
- 1: Any conversion error (file not found, parse error, write error)
- 2: Usage error (invalid flags, missing required arguments)

**Validation:**
- Test successful conversion: `echo $?` should print 0
- Test file not found: `echo $?` should print 1
- Test invalid flags: `echo $?` should print 2
- Script integration test using exit codes for control flow

## Traceability Mapping

| Acceptance Criteria                    | Tech Spec Section                                                                                                            | Component/Module                                                                         | Test Approach                                                                                                                                                    |
| -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **AC-1: Basic CLI Command Structure**  | Services and Modules → cmd/cli/convert.go<br>APIs and Interfaces → convertSingleFile()                                       | `cmd/cli/convert.go`<br>`cmd/cli/root.go`<br>`internal/converter`                        | Integration test: Run CLI with sample files, verify output file created, check success message format, validate exit code 0                                      |
| **AC-2: Format Auto-Detection**        | Services and Modules → cmd/cli/format.go<br>APIs and Interfaces → detectFormat()                                             | `cmd/cli/format.go`<br>`detectFormat()`<br>`detectFormatFromBytes()`                     | Unit test: Pass files with different extensions, verify correct format returned. Edge case: Unknown extension returns error                                      |
| **AC-3: Batch Processing**             | Services and Modules → cmd/cli/batch.go<br>Workflows → Batch Conversion Flow<br>APIs and Interfaces → convertBatch()         | `cmd/cli/batch.go`<br>`convertBatch()`<br>Worker pool pattern                            | Performance test: Benchmark 100 files from testdata/, measure total time with `time` command, verify <2s target. Test parallel execution by monitoring CPU usage |
| **AC-4: Verbose Logging Mode**         | Services and Modules → cmd/cli/output.go<br>NFR → Observability → Verbose Mode Output<br>Data Models → ConvertFlags.Verbose  | `cmd/cli/output.go`<br>`formatResult()`<br>slog structured logging                       | Integration test: Compare normal vs verbose output, verify logs to stderr, check parameter details appear, validate slog JSON format in verbose mode             |
| **AC-5: JSON Output**                  | Services and Modules → cmd/cli/output.go<br>Data Models → ConversionResult struct<br>APIs and Interfaces → formatResult()    | `cmd/cli/output.go`<br>`ConversionResult` JSON marshal<br>`BatchResult` JSON marshal     | Unit test: Parse JSON output with `encoding/json`, verify all required fields present. Integration test: Pipe to `jq` and Python json.loads()                    |
| **AC-6: Error Handling**               | Services and Modules → cmd/cli/convert.go<br>Data Models → ConversionError<br>Workflows → Error Handling Sequence            | `cmd/cli/convert.go`<br>`internal/converter/errors.go`<br>`ConversionError` type         | Negative test: Provide corrupted files, invalid formats, verify user-friendly error messages, confirm exit code 1, check no output file created                  |
| **AC-7: Overwrite Protection**         | Services and Modules → cmd/cli/convert.go<br>Data Models → ConvertFlags.Overwrite<br>Workflows → Single File Conversion Flow | `cmd/cli/convert.go`<br>File existence check<br>`--overwrite` flag logic                 | Integration test: Create existing file, run without --overwrite (expect error), verify file unchanged, run with --overwrite (expect success)                     |
| **AC-8: Cross-Platform Compatibility** | System Architecture Alignment → Cross-platform single-binary<br>NFR → Reliability → Cross-Platform Consistency               | Go cross-compilation<br>`filepath` stdlib (OS-agnostic paths)<br>Build matrix (OS, arch) | Cross-platform test: Build for Windows/macOS/Linux, run same input, SHA256 hash outputs, verify byte-identical results                                           |
| **AC-9: Help and Version**             | Services and Modules → cmd/cli/root.go<br>Cobra framework                                                                    | `cmd/cli/root.go`<br>Cobra `--help` auto-generation<br>`--version` flag                  | Manual test: Run `recipe --help` and `recipe --version`, verify output completeness and accuracy                                                                 |
| **AC-10: Performance Targets**         | NFR → Performance → CLI Execution Speed<br>Workflows → Batch Conversion Flow (parallel)                                      | `internal/converter/converter.go`<br>`cmd/cli/batch.go` worker pool<br>Go benchmarks     | Benchmark test: `go test -bench=. ./internal/converter/` for single file (<20ms). Time command: `time recipe convert --batch *.xmp --to np3` for 100 files (<2s) |
| **AC-11: Exit Code Consistency**       | Workflows → Exit Codes table<br>Services and Modules → Error handling in main.go                                             | `cmd/cli/main.go`<br>`os.Exit()` calls<br>Cobra error handling                           | Shell script test: Test success (check $? == 0), file not found ($? == 1), invalid flags ($? == 2)                                                               |

### Traceability to PRD Requirements

| PRD Requirement                     | Epic 3 Implementation                                                               | Verification                           |
| ----------------------------------- | ----------------------------------------------------------------------------------- | -------------------------------------- |
| **FR-3.1: Command Structure**       | cmd/cli/convert.go with Cobra framework, syntax: `recipe convert INPUT --to FORMAT` | AC-1, AC-9 (help text)                 |
| **FR-3.2: Batch Processing**        | cmd/cli/batch.go with worker pool, `--batch *.xmp` glob pattern support             | AC-3 (100 files <2s)                   |
| **FR-3.3: Format Auto-Detection**   | cmd/cli/format.go detectFormat() from extension or magic bytes                      | AC-2 (no --from required)              |
| **FR-3.4: Verbose Mode**            | slog structured logging to stderr with --verbose flag                               | AC-4 (detailed logs)                   |
| **FR-3.5: JSON Output**             | ConversionResult/BatchResult JSON marshaling with --json flag                       | AC-5 (jq compatible)                   |
| **NFR-1.3: CLI Performance**        | Single file <20ms, batch 100 files <2s, cold start <100ms                           | AC-10 (benchmarks)                     |
| **NFR-2.1: Zero Data Exfiltration** | All processing local, no network access in CLI code                                 | Code review (no net/http imports)      |
| **NFR-3.1: Conversion Accuracy**    | Uses shared converter.Convert() API (95%+ accuracy from Epic 1)                     | Round-trip tests in internal/converter |
| **Cross-Platform Distribution**     | Go cross-compilation for Windows/macOS/Linux amd64/arm64                            | AC-8 (byte-identical outputs)          |

### Test Coverage Matrix

| Component            | Unit Tests                                                          | Integration Tests                            | Benchmarks             | Manual Tests            |
| -------------------- | ------------------------------------------------------------------- | -------------------------------------------- | ---------------------- | ----------------------- |
| `cmd/cli/convert.go` | ✓ Flag parsing<br>✓ Path handling                                   | ✓ End-to-end conversion<br>✓ Error scenarios | -                      | ✓ Real file workflows   |
| `cmd/cli/batch.go`   | ✓ Worker pool logic<br>✓ Result aggregation                         | ✓ 100 file batch test                        | ✓ Parallel performance | -                       |
| `cmd/cli/format.go`  | ✓ detectFormat() with extensions<br>✓ detectFormatFromBytes() magic | ✓ Unknown format errors                      | -                      | -                       |
| `cmd/cli/output.go`  | ✓ formatResult() variations<br>✓ JSON marshaling                    | ✓ Verbose/JSON modes                         | -                      | ✓ Human-readable output |
| `cmd/cli/root.go`    | ✓ Cobra setup<br>✓ Flag definitions                                 | ✓ --help output                              | -                      | ✓ Version display       |
| **Shared Converter** | (Epic 1 tests)                                                      | ✓ CLI calls converter.Convert()              | ✓ <20ms single file    | -                       |

### PRD Success Criteria Mapping

| PRD Success Criteria                            | Epic 3 Contribution                                                                      | Measurement                |
| ----------------------------------------------- | ---------------------------------------------------------------------------------------- | -------------------------- |
| **10-100x performance improvement over Python** | CLI performance targets: Single file <20ms (vs 50-200ms), batch 100 files <2s (vs 5-20s) | AC-10 benchmarks           |
| **Cross-platform (Windows, macOS, Linux)**      | Single binary distribution for all platforms using Go cross-compilation                  | AC-8 cross-platform tests  |
| **95%+ accuracy**                               | Reuses Epic 1 converter.Convert() API with proven accuracy                               | Epic 1 round-trip tests    |
| **Single binary distribution**                  | Go produces standalone executable, zero dependencies at runtime                          | GitHub Release artifacts   |
| **Easy installation**                           | Download binary from GitHub Releases, chmod +x, move to PATH                             | Distribution documentation |

## Risks, Assumptions, Open Questions

### Risks

**R-1: Cross-Platform Path Handling Complexity**
- **Risk:** Windows uses backslashes (\), Unix uses forward slashes (/) - potential for path bugs
- **Impact:** Medium - Could cause file not found errors on Windows
- **Likelihood:** Low - Go's `filepath` package handles this automatically
- **Mitigation:**
  - Use `filepath.Join()` instead of string concatenation
  - Use `filepath.Glob()` which respects OS conventions
  - Test on all three platforms (Windows, macOS, Linux)
  - Add cross-platform integration tests in CI

**R-2: Large Batch Performance Variability**
- **Risk:** 100 file target may not scale linearly to 1,000 or 10,000 files due to file I/O bottlenecks
- **Impact:** Medium - User expectations may not match performance at scale
- **Likelihood:** Medium - Disk I/O can become bottleneck with thousands of files
- **Mitigation:**
  - Document performance targets clearly (100 files <2s, larger batches proportional)
  - Implement adaptive worker pool sizing based on I/O vs CPU bottleneck detection
  - Add --parallel flag for user to tune worker count
  - Consider streaming I/O for very large batches (future optimization)

**R-3: Cobra Framework Breaking Changes**
- **Risk:** Cobra v2.0 may introduce breaking API changes
- **Impact:** Low - Would require code updates but not fundamental redesign
- **Likelihood:** Low - Cobra is stable, major version changes are infrequent
- **Mitigation:**
  - Pin exact Cobra version in go.mod
  - Monitor Cobra release notes
  - Vendor dependencies (go mod vendor) for reproducible builds

**R-4: User Expectation Mismatch on Error Recovery**
- **Risk:** Users may expect partial conversion recovery (resume from failure point)
- **Impact:** Low - Current design is stateless, no resume capability
- **Likelihood:** Medium - Power users with large batches may request this
- **Mitigation:**
  - Document that CLI is stateless (can re-run safely)
  - Continue-on-error mode processes all files regardless of individual failures
  - Consider checkpoint/resume feature in future epic if demanded

**R-5: JSON Output Schema Stability**
- **Risk:** JSON schema changes could break downstream scripts
- **Impact:** High - Breaking change for scripting users
- **Likelihood:** Low - Schema designed comprehensively upfront
- **Mitigation:**
  - Version JSON schema in output (add `"schema_version": "1.0"` field)
  - Additive-only changes (never remove fields, only add new ones)
  - Document JSON schema in README with stability guarantee
  - Semantic versioning: MAJOR bump if schema breaks compatibility

### Assumptions

**A-1: File System Access**
- **Assumption:** CLI has read/write permissions to specified file paths
- **Validation:** Document permission errors clearly, suggest chmod/chown solutions
- **Impact if False:** CLI cannot function without file access

**A-2: File Size Limits**
- **Assumption:** Preset files are <10MB (typical range: 1KB - 100KB)
- **Validation:** 10MB hard limit enforced in parser validation
- **Impact if False:** Memory exhaustion on very large files (unlikely for presets)

**A-3: Go 1.24+ Availability**
- **Assumption:** Users can install/upgrade to Go 1.24+ for development
- **Validation:** Document minimum Go version clearly in README
- **Impact if False:** Developers on older Go versions cannot build from source (can still use pre-built binaries)

**A-4: Single-Threaded I/O Acceptable**
- **Assumption:** File I/O is fast enough that parallelization at conversion level (not I/O level) is sufficient
- **Validation:** Benchmark confirms <2s for 100 files
- **Impact if False:** May need async I/O or memory-mapped files (optimization for future)

**A-5: Conversion Engine Stability**
- **Assumption:** `internal/converter` API from Epic 1 is stable and accurate (95%+)
- **Validation:** Epic 1 round-trip tests confirm accuracy
- **Impact if False:** CLI would propagate converter bugs, but this is out of Epic 3 scope (Epic 1 responsibility)

**A-6: User Familiarity with CLI**
- **Assumption:** Target users (power users, automation workflows) are comfortable with command-line tools
- **Validation:** PRD defines CLI as "for automation, scripting, and batch processing"
- **Impact if False:** Web interface (Epic 2) and TUI (Epic 4) serve less technical users

### Open Questions

**Q-1: Should CLI support stdin/stdout piping?**
- **Context:** Unix convention allows `cat file.xmp | recipe convert - --to np3 > output.np3`
- **Trade-offs:**
  - Pro: More flexible, composable with shell tools
  - Con: Adds complexity, requires streaming implementation
- **Decision:** Defer to post-MVP. File paths sufficient for initial release. Consider if users request.

**Q-2: Should batch mode preserve directory structure?**
- **Context:** If input is `presets/**/*.xmp`, should output mirror directory tree?
- **Example:**
  - Input: `presets/portraits/vintage.xmp`
  - Output Option 1: `vintage.np3` (flat, current design)
  - Output Option 2: `presets/portraits/vintage.np3` (preserve structure)
- **Decision:** **Flat structure for MVP** (simpler). Add `--preserve-structure` flag if users need it. Document in README.

**Q-3: What level of progress detail in batch mode?**
- **Context:** Should CLI show individual file names or just count?
- **Options:**
  - Minimal: "Processing 45/100 files..." (current design)
  - Detailed: "Converting file 45/100: portrait_vintage.xmp..."
  - Silent: No progress, only final summary
- **Decision:** **Minimal for normal mode**, detailed in verbose mode. Allows users to choose noise level.

**Q-4: Should CLI validate output format before writing?**
- **Context:** Could re-parse generated file to ensure it's valid before writing to disk
- **Trade-offs:**
  - Pro: Extra safety, catch generator bugs
  - Con: Doubles conversion time (parse → generate → parse again)
- **Decision:** **No re-validation for MVP**. Epic 1 tests already validate generators. Trust shared converter. Add `--validate` flag if paranoid users request.

**Q-5: How to handle format version differences?**
- **Context:** XMP has different versions (crs:Version="7.0" vs "12.0"), same for lrtemplate
- **Current Behavior:** Parser handles all known versions, generator uses latest
- **Question:** Should CLI expose `--xmp-version` flag to control output version?
- **Decision:** **Defer**. Generate latest version by default. Add flag only if compatibility issues reported by users.

**Q-6: Should --json mode support streaming (newline-delimited JSON)?**
- **Context:** For large batches, single JSON object could be huge. NDJSON alternative:
  ```
  {"input": "file1.xmp", ...}
  {"input": "file2.xmp", ...}
  ```
- **Trade-offs:**
  - Pro: Streamable, works with `jq -c` and similar tools
  - Con: Not valid JSON (requires line-by-line parsing)
- **Decision:** **Single JSON object for MVP** (ConversionResult or BatchResult). Add `--json-stream` flag if users need NDJSON.

**Q-7: Should CLI auto-detect batch mode without --batch flag?**
- **Context:** Could detect glob patterns (* or ?) and switch to batch mode automatically
- **Example:** `recipe convert *.xmp --to np3` (without --batch) auto-detects batch
- **Trade-offs:**
  - Pro: Simpler UX, fewer flags to remember
  - Con: Ambiguity if `*.xmp` matches single file
- **Decision:** **Keep --batch explicit for MVP**. Prevents accidental batch mode. Clearer intent. Reconsider if users find it tedious.

### Risk Prioritization Summary

| Risk ID                     | Impact | Likelihood | Priority   | Mitigation Status                                   |
| --------------------------- | ------ | ---------- | ---------- | --------------------------------------------------- |
| R-1: Path Handling          | Medium | Low        | **Medium** | Use filepath stdlib, cross-platform tests planned   |
| R-2: Batch Performance      | Medium | Medium     | **High**   | Document targets, add --parallel flag, benchmark CI |
| R-3: Cobra Breaking Changes | Low    | Low        | **Low**    | Pin version in go.mod, monitor releases             |
| R-4: Resume Capability      | Low    | Medium     | **Medium** | Document stateless design, continue-on-error mode   |
| R-5: JSON Schema            | High   | Low        | **Medium** | Version schema, additive-only changes, semver       |

### Next Steps for Risk Mitigation

1. **Before Development:**
   - Define exact JSON schema with version field
   - Set up cross-platform CI matrix (Windows/macOS/Linux)

2. **During Development:**
   - Implement comprehensive path handling tests
   - Benchmark batch mode at 100, 500, 1000 file scales
   - Add --parallel flag for user tuning

3. **Before Release:**
   - Document all assumptions and limitations in README
   - Create troubleshooting guide for common errors (permissions, path issues)
   - Review open questions with potential beta users

## Test Strategy Summary

### Testing Philosophy

**Comprehensive Validation:** CLI is the automation interface - reliability is critical. Test strategy emphasizes real-world scenarios, cross-platform consistency, and performance validation.

**Shared Converter Trust:** CLI delegates all conversion logic to `internal/converter` (Epic 1). Focus CLI tests on I/O, flag parsing, and output formatting. Conversion accuracy is validated by Epic 1 tests.

### Test Levels

#### 1. Unit Tests

**Scope:** Individual functions in isolation

**Components Under Test:**
- `cmd/cli/format.go` - Format detection logic
  - `TestDetectFormat()` - Extension-based detection (np3, xmp, lrtemplate)
  - `TestDetectFormatFromBytes()` - Magic byte detection (fallback)
  - Edge cases: Unknown extensions, empty files, ambiguous formats

- `cmd/cli/output.go` - Output formatting
  - `TestFormatResult()` - Normal, verbose, JSON modes
  - `TestFormatBatchResult()` - Batch summary formatting
  - JSON marshaling validation (all required fields present)

- `cmd/cli/convert.go` - Flag parsing and validation
  - `TestValidateFlags()` - Required flags, incompatible combinations
  - `TestOutputPathGeneration()` - Extension replacement logic
  - Path handling edge cases (special characters, Unicode)

**Tooling:**
- Go standard `testing` package
- Table-driven tests for multiple scenarios
- `encoding/json` for JSON validation

**Example:**
```go
func TestDetectFormat(t *testing.T) {
    tests := []struct {
        filename string
        want     string
        wantErr  bool
    }{
        {"portrait.np3", "np3", false},
        {"preset.xmp", "xmp", false},
        {"classic.lrtemplate", "lrtemplate", false},
        {"unknown.txt", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.filename, func(t *testing.T) {
            got, err := detectFormat(tt.filename)
            if (err != nil) != tt.wantErr {
                t.Errorf("detectFormat() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("detectFormat() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

#### 2. Integration Tests

**Scope:** End-to-end CLI execution with real sample files

**Test Scenarios:**

**A. Basic Conversion (AC-1)**
```go
func TestCLI_BasicConversion(t *testing.T) {
    // Setup: Copy sample file to temp directory
    tmpDir := t.TempDir()
    inputPath := filepath.Join(tmpDir, "portrait.xmp")
    copyFile("testdata/xmp/portrait.xmp", inputPath)

    // Execute CLI
    cmd := exec.Command("recipe", "convert", inputPath, "--to", "np3")
    output, err := cmd.CombinedOutput()

    // Assertions
    assert.NoError(t, err)
    assert.Contains(t, string(output), "✓ Converted")
    assert.FileExists(t, filepath.Join(tmpDir, "portrait.np3"))
    assert.Equal(t, 0, cmd.ProcessState.ExitCode())
}
```

**B. Batch Processing (AC-3)**
```go
func TestCLI_BatchConversion(t *testing.T) {
    tmpDir := t.TempDir()

    // Setup: Copy 100 sample XMP files
    for i := 0; i < 100; i++ {
        copyFile("testdata/xmp/sample.xmp", filepath.Join(tmpDir, fmt.Sprintf("file%d.xmp", i)))
    }

    // Execute batch conversion
    start := time.Now()
    cmd := exec.Command("recipe", "convert", "--batch", filepath.Join(tmpDir, "*.xmp"), "--to", "np3")
    output, err := cmd.CombinedOutput()
    elapsed := time.Since(start)

    // Assertions
    assert.NoError(t, err)
    assert.Contains(t, string(output), "✓ Converted 100 files")
    assert.Less(t, elapsed, 2*time.Second, "Batch conversion should complete in <2s")

    // Verify all output files created
    files, _ := filepath.Glob(filepath.Join(tmpDir, "*.np3"))
    assert.Len(t, files, 100)
}
```

**C. Error Scenarios (AC-6, AC-7)**
```go
func TestCLI_ErrorHandling(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        wantErr  string
        exitCode int
    }{
        {
            name:     "file not found",
            args:     []string{"convert", "nonexistent.xmp", "--to", "np3"},
            wantErr:  "no such file or directory",
            exitCode: 1,
        },
        {
            name:     "invalid format",
            args:     []string{"convert", "test.txt", "--to", "np3"},
            wantErr:  "invalid file format",
            exitCode: 1,
        },
        {
            name:     "overwrite protection",
            args:     []string{"convert", "existing.xmp", "--to", "np3"},
            wantErr:  "already exists",
            exitCode: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command("recipe", tt.args...)
            output, _ := cmd.CombinedOutput()

            assert.Contains(t, string(output), tt.wantErr)
            assert.Equal(t, tt.exitCode, cmd.ProcessState.ExitCode())
        })
    }
}
```

**D. Output Modes (AC-4, AC-5)**
- Verbose mode: Check stderr contains detailed logs
- JSON mode: Parse output with `encoding/json`, validate schema
- Normal mode: Verify stdout contains success message only

**Tooling:**
- `os/exec` to run CLI as subprocess
- `testify/assert` for readable assertions (optional, can use stdlib)
- Temp directories for isolated test environments

---

#### 3. Cross-Platform Tests (AC-8)

**Scope:** Verify identical behavior on Windows, macOS, Linux

**Test Matrix:**
- **Platforms:** windows-amd64, darwin-amd64, darwin-arm64, linux-amd64, linux-arm64
- **Test Cases:** Basic conversion, batch processing, error handling
- **Validation:** SHA256 hash of output files must match across platforms

**CI Configuration (GitHub Actions):**
```yaml
name: Cross-Platform Tests

on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.24']

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run Tests
        run: go test -v ./cmd/cli/...

      - name: Build CLI
        run: go build -o recipe ./cmd/cli/

      - name: Test Conversion Output Consistency
        run: |
          ./recipe convert testdata/xmp/portrait.xmp --to np3
          sha256sum portrait.np3 > checksum-${{ matrix.os }}.txt

      - name: Upload Checksums
        uses: actions/upload-artifact@v3
        with:
          name: checksums
          path: checksum-*.txt

  verify-checksums:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v3
      - name: Compare Checksums
        run: |
          # Verify all platforms produced identical output
          diff checksum-ubuntu-latest.txt checksum-macos-latest.txt
          diff checksum-ubuntu-latest.txt checksum-windows-latest.txt
```

---

#### 4. Performance Benchmarks (AC-10)

**Scope:** Validate performance targets (<20ms single file, <2s batch 100 files)

**Go Benchmarks:**
```go
// cmd/cli/convert_test.go
func BenchmarkCLI_SingleFileConversion(b *testing.B) {
    inputPath := "testdata/xmp/portrait.xmp"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cmd := exec.Command("recipe", "convert", inputPath, "--to", "np3", "-o", "/tmp/bench_output.np3")
        if err := cmd.Run(); err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkCLI_BatchConversion100Files(b *testing.B) {
    // Setup: Create temp directory with 100 XMP files
    tmpDir := b.TempDir()
    for i := 0; i < 100; i++ {
        copyFile("testdata/xmp/sample.xmp", filepath.Join(tmpDir, fmt.Sprintf("file%d.xmp", i)))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cmd := exec.Command("recipe", "convert", "--batch", filepath.Join(tmpDir, "*.xmp"), "--to", "np3")
        if err := cmd.Run(); err != nil {
            b.Fatal(err)
        }
    }
}
```

**Run Benchmarks:**
```bash
go test -bench=. -benchmem ./cmd/cli/
```

**Expected Output:**
```
BenchmarkCLI_SingleFileConversion-8         100     15000000 ns/op     (15ms)
BenchmarkCLI_BatchConversion100Files-8       1     1500000000 ns/op    (1.5s)
```

**Validation:**
- Single file: <20ms/op ✓
- Batch 100: <2s/op ✓

**Alternative (Shell Timing):**
```bash
# Single file
time recipe convert testdata/xmp/portrait.xmp --to np3

# Batch 100 files
time recipe convert --batch testdata/xmp/*.xmp --to np3
```

---

#### 5. Manual Tests

**Scope:** User experience validation not easily automated

**Test Cases:**

**A. Help Text Completeness (AC-9)**
- Run `recipe --help`
- Verify all flags documented
- Check examples are accurate
- Confirm formatting is readable

**B. Version Display (AC-9)**
- Run `recipe --version`
- Verify version matches git tag

**C. Real-World Workflow Validation**
- User story: "Convert my entire preset library"
  - Create directory with mixed formats (100+ files)
  - Run batch conversion
  - Verify all files converted correctly
  - Check warnings for unmappable parameters

**D. Error Message Clarity**
- Intentionally cause errors (missing file, invalid format)
- Verify error messages are user-friendly
- Confirm actionable suggestions provided

**E. Cross-Platform Path Handling**
- Windows: Test with backslashes (`recipe convert C:\presets\file.xmp --to np3`)
- Unix: Test with forward slashes and spaces (`recipe convert "/path/with spaces/file.xmp" --to np3`)

---

### Test Coverage Goals

| Component            | Unit Test Coverage      | Integration Test Coverage | Total Goal |
| -------------------- | ----------------------- | ------------------------- | ---------- |
| `cmd/cli/convert.go` | 80%+                    | 95%+                      | **90%+**   |
| `cmd/cli/batch.go`   | 85%+                    | 100% (critical path)      | **90%+**   |
| `cmd/cli/format.go`  | 95%+ (pure logic)       | 80%                       | **90%+**   |
| `cmd/cli/output.go`  | 90%+                    | 85%                       | **90%+**   |
| `cmd/cli/root.go`    | 70% (Cobra boilerplate) | 90%                       | **80%+**   |
| **Overall CLI**      | **85%+**                | **90%+**                  | **90%+**   |

**Note:** Shared converter (`internal/converter`) is validated by Epic 1 tests (not counted in Epic 3 coverage).

---

### Test Execution Strategy

**Development (Local):**
```bash
# Run all tests
go test ./cmd/cli/...

# Run with coverage
go test -cover ./cmd/cli/...

# Run specific test
go test -run TestCLI_BasicConversion ./cmd/cli/

# Run benchmarks
go test -bench=. ./cmd/cli/
```

**CI/CD (Automated):**
1. **On Every Commit:**
   - Unit tests (all platforms)
   - Integration tests (Linux only for speed)
   - Code coverage report (upload to Codecov)

2. **On Pull Request:**
   - Full integration tests (all platforms)
   - Cross-platform output consistency check
   - Performance benchmarks (compare vs baseline)

3. **On Release Tag:**
   - Full test suite (all platforms)
   - Performance regression check
   - Manual smoke tests
   - Build release binaries

**Test Data:**
- Use sample files from `testdata/` directory (1,501 files)
- Focus on representative subset for fast tests (e.g., 10 files per format)
- Full suite for CI (all 1,501 files)

---

### Test Automation

**Makefile Targets:**
```makefile
.PHONY: test test-unit test-integration test-bench test-coverage

# Run all tests
test:
	go test -v ./cmd/cli/...

# Run only unit tests (fast)
test-unit:
	go test -v -short ./cmd/cli/...

# Run integration tests (slower)
test-integration:
	go test -v -run Integration ./cmd/cli/...

# Run benchmarks
test-bench:
	go test -bench=. -benchmem ./cmd/cli/

# Generate coverage report
test-coverage:
	go test -coverprofile=coverage.out ./cmd/cli/...
	go tool cover -html=coverage.out -o coverage.html
```

**Usage:**
```bash
make test            # All tests
make test-unit       # Quick feedback loop
make test-bench      # Performance validation
make test-coverage   # Coverage report
```

---

### Edge Cases and Corner Cases

**Path Handling:**
- Files with spaces in name: `"my preset.xmp"`
- Unicode characters: `"プリセット.xmp"`
- Windows special characters: `"preset:name.xmp"` (invalid)
- Relative vs absolute paths
- Symlinks (follow or error?)

**File Size:**
- Empty file (0 bytes) - expect error
- Minimal valid file (smallest parsable)
- Maximum valid file (10MB limit)

**Concurrent Batch:**
- 1 file batch (edge of parallelization)
- 10,000 file batch (stress test)
- All files invalid (error handling at scale)

**Output Scenarios:**
- Output directory doesn't exist - create or error?
- Output to read-only location - permission error
- Overwrite existing file (with/without --overwrite flag)

**Format Detection:**
- File with wrong extension (`.xmp` but actually `.np3` content)
- No extension: `preset` (fallback to content detection)
- Multiple extensions: `preset.backup.xmp`

---

### Success Criteria for Test Strategy

✅ **90%+ code coverage** for CLI codebase (`cmd/cli/`)
✅ **100% of acceptance criteria** have corresponding automated tests
✅ **Cross-platform consistency** validated in CI (SHA256 hashes match)
✅ **Performance targets met** (<20ms single, <2s batch 100)
✅ **Zero regressions** in existing Epic 1 conversion accuracy
✅ **Manual test checklist** completed before release (help text, UX, error messages)
