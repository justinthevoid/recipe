# Epic Technical Specification: Data Extraction & Inspection

Date: 2025-11-06
Author: Justin
Epic ID: epic-5
Status: Draft

---

## Overview

Epic 5 implements data extraction and inspection capabilities for Recipe, enabling advanced users to examine preset parameters, understand file structure internals, and validate conversion accuracy. While Epic 1 focuses on conversion functionality, Epic 5 provides transparency tools that reveal the "what" and "why" behind preset files. This epic delivers three CLI-based inspection tools:
- **Parameter Inspector** (`recipe inspect FILE`) - Extracts and displays all parameters as JSON for programmatic analysis
- **Binary Structure Visualizer** (`recipe inspect FILE --binary`) - Shows NP3 hex dumps with field annotations for reverse engineering
- **Diff Tool** (`recipe diff FILE1 FILE2`) - Compares parameters between files to validate conversion accuracy

These tools serve power users learning color science, developers extending Recipe's format support, and photographers validating preset conversions. All tools leverage the existing `internal/converter` engine and `UniversalRecipe` model from Epic 1, adding read-only inspection capabilities without modifying the core conversion logic.

**Core Value Proposition**: Transparency and education through data extraction, enabling users to understand preset internals, validate conversions, and contribute to the open-source color science community.

## Objectives and Scope

### Primary Objectives
1. **Parameter Inspection**: Provide JSON export of all preset parameters for programmatic analysis and learning
2. **Binary Transparency**: Enable hex-level inspection of NP3 files with field annotations for reverse engineering
3. **Conversion Validation**: Offer diff capabilities to compare parameters across formats and validate conversion accuracy
4. **Educational Value**: Help users understand color science by revealing internal parameter representations

### In Scope
- `recipe inspect FILE` command with JSON output of all UniversalRecipe fields
- `recipe inspect FILE --binary` command with annotated hex dump for NP3 files
- `recipe diff FILE1 FILE2` command with parameter comparison across any format pair
- Metadata inclusion (format, version, source info) in inspect output
- Cross-format diff support (e.g., compare NP3 to XMP)
- Human-readable diff formatting with color-coded changes
- Integration with existing converter.Convert() and format parsers
- Comprehensive error handling for invalid files

### Out of Scope
- Visual diff rendering (text-based only, no GUI)
- Batch inspection operations (single file at a time)
- File modification capabilities (read-only inspection)
- Export to formats other than JSON (inspect) and text (diff)
- Advanced diffing algorithms (unified diff, word-level diff)
- Web interface for inspection tools (CLI only for Epic 5)
- Binary visualization for XMP/lrtemplate (text formats don't need hex dumps)

### Success Criteria
- Inspect command outputs valid, parseable JSON for all supported formats
- Binary mode shows byte offsets, hex values, and field labels for all known NP3 fields
- Diff tool accurately identifies all parameter changes between file pairs
- Cross-format diff correctly maps parameters (e.g., NP3 Contrast vs XMP Contrast2012)
- Performance: Inspect completes in <50ms, diff completes in <100ms
- Error messages clearly identify file format issues and parsing failures
- JSON output includes 100% of UniversalRecipe fields (no data omission)

## System Architecture Alignment

Epic 5 extends the CLI interface established in Epic 3 with three new inspection commands that leverage the existing hub-and-spoke conversion architecture:

```
┌─────────────────────────────────────────────────────────┐
│ CLI Commands (Epic 3 + Epic 5)                          │
│  • recipe convert FILE1 FILE2    (Epic 3)              │
│  • recipe inspect FILE           (Epic 5 - NEW)        │
│  • recipe inspect FILE --binary  (Epic 5 - NEW)        │
│  • recipe diff FILE1 FILE2       (Epic 5 - NEW)        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │ internal/converter/   │
         │   Parse → Universal   │◄──── Reuses Epic 1 parsers
         └───────────┬───────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │ internal/inspect/     │◄──── NEW: Epic 5 module
         │  • ToJSON()           │
         │  • BinaryDump()       │
         │  • Diff()             │
         └───────────────────────┘
```

### Architecture Alignment Points

1. **CLI Pattern Consistency** (aligns with Epic 3: CLI Interface)
   - All commands use Cobra framework with consistent flag patterns
   - `inspect` and `diff` commands follow same structure as `convert` command
   - Automatic format detection reuses Epic 3's detection logic
   - Error handling uses same ConversionError type from Epic 1

2. **Reuses Existing Parsers** (aligns with Epic 1: Core Conversion Engine)
   - No new parsing logic required - uses `np3.Parse()`, `xmp.Parse()`, `lrtemplate.Parse()`
   - Operates on `UniversalRecipe` struct from `internal/model`
   - All inspection happens AFTER parsing (read-only operations)

3. **New Module: internal/inspect/** (extends architecture)
   - `inspect.ToJSON(recipe *model.UniversalRecipe) ([]byte, error)` - Serializes to JSON
   - `inspect.BinaryDump(data []byte, format string) (string, error)` - Annotated hex dump
   - `inspect.Diff(recipe1, recipe2 *model.UniversalRecipe) (string, error)` - Parameter comparison

4. **CLI Command Structure** (aligns with Pattern 10: CLI Command Pattern)
   ```
   cmd/cli/
   ├── main.go          (Epic 3 - root command)
   ├── convert.go       (Epic 3 - conversion)
   ├── inspect.go       (Epic 5 - NEW)
   └── diff.go          (Epic 5 - NEW)
   ```

5. **No Architecture Changes Required**
   - Epic 5 is purely additive (no modifications to Epic 1 or Epic 3 code)
   - Inspection tools are read-only (no risk of breaking conversion logic)
   - Performance constraints align with existing <100ms conversion goals

## Detailed Design

### Services and Modules

Epic 5 introduces one new internal module and two new CLI command files:

| Module/Command         | Responsibility             | Key Functions                                                                                                          | Dependencies                                                       |
| ---------------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| **internal/inspect/**  | Core inspection logic      | `ToJSON()` - Serialize UniversalRecipe to JSON<br>`BinaryDump()` - Annotate hex dump<br>`Diff()` - Compare two recipes | `internal/model`, `encoding/json`                                  |
| **cmd/cli/inspect.go** | CLI command for inspection | `inspectCmd` - Cobra command<br>`runInspect()` - Command handler<br>`formatJSON()` - Pretty print                      | `internal/inspect`, `internal/formats/*`, `github.com/spf13/cobra` |
| **cmd/cli/diff.go**    | CLI command for diffing    | `diffCmd` - Cobra command<br>`runDiff()` - Command handler<br>`colorizeOutput()` - Terminal coloring                   | `internal/inspect`, `internal/formats/*`, `github.com/spf13/cobra` |

**Module Details:**

**internal/inspect/inspect.go**
- `ToJSON(recipe *model.UniversalRecipe) ([]byte, error)` - Marshals UniversalRecipe to indented JSON
- `ToJSONWithMetadata(recipe *model.UniversalRecipe, sourceFile, format string) ([]byte, error)` - Includes file metadata wrapper
- Input: Parsed UniversalRecipe from any format
- Output: Pretty-printed JSON with all fields
- Error handling: Returns error if JSON marshaling fails

**internal/inspect/binary.go**
- `BinaryDump(data []byte, format string) (string, error)` - Creates annotated hex dump
- Only supports NP3 format (returns error for XMP/lrtemplate)
- Output format: `[offset] hex_bytes  field_name (value)`
- Example: `[0x0000] 4E 50  Magic ("NP")`
- Uses hardcoded NP3 field map from reverse engineering documentation

**internal/inspect/diff.go**
- `Diff(recipe1, recipe2 *model.UniversalRecipe) (string, error)` - Compares all fields
- Returns formatted string with changed parameters only
- Format: `FieldName: old_value → new_value`
- Supports tolerance for float comparisons (±0.001)
- Highlights significant differences (>5% change for percentages)

### Data Models and Contracts

Epic 5 does not introduce new data models - it operates entirely on existing structures from Epic 1:

**Primary Data Model: UniversalRecipe** (from `internal/model/recipe.go`)
- All inspection operates on this existing struct
- No modifications required to UniversalRecipe definition
- JSON serialization uses struct tags already defined in Epic 1

**New Type: InspectOutput** (added to `internal/inspect/types.go`)
```go
type InspectOutput struct {
    Metadata struct {
        SourceFile   string `json:"source_file"`
        SourceFormat string `json:"source_format"`
        ParsedAt     string `json:"parsed_at"`     // ISO 8601 timestamp
        RecipeVersion string `json:"recipe_version"` // Recipe tool version
    } `json:"metadata"`
    Parameters *model.UniversalRecipe `json:"parameters"`
}
```

**New Type: DiffResult** (added to `internal/inspect/types.go`)
```go
type DiffResult struct {
    Field      string      `json:"field"`
    OldValue   interface{} `json:"old_value"`
    NewValue   interface{} `json:"new_value"`
    ChangeType string      `json:"change_type"` // "modified", "added", "removed"
    Significant bool       `json:"significant"` // true if >5% change
}
```

**Binary Field Map** (constant in `internal/inspect/binary.go`)
```go
var np3FieldMap = map[int]string{
    0x0000: "Magic Bytes (NP)",
    0x0002: "Version",
    0x0042: "Contrast (-3 to +3)",
    0x0043: "Brightness (-1 to +1)",
    0x0044: "Saturation (-3 to +3)",
    0x0045: "Hue (-9° to +9°)",
    0x0046: "Sharpness (0-9)",
    // ... additional fields from reverse engineering
}
```

### APIs and Interfaces

**CLI Command Interfaces:**

**1. Inspect Command**
```bash
# Basic usage - JSON output to stdout
recipe inspect portrait.np3

# Binary mode - hex dump with annotations (NP3 only)
recipe inspect portrait.np3 --binary

# Save JSON to file
recipe inspect portrait.xmp --output portrait.json

# Pretty print (default is already pretty, this forces color)
recipe inspect portrait.lrtemplate --pretty
```

**Cobra Command Signature:**
```go
var inspectCmd = &cobra.Command{
    Use:   "inspect [file]",
    Short: "Extract and display preset parameters as JSON",
    Long:  `Inspect parses a preset file and outputs all parameters as JSON.
Supports NP3, XMP, and lrtemplate formats. Use --binary flag for hex dump (NP3 only).`,
    Args:  cobra.ExactArgs(1),
    RunE:  runInspect,
}

// Flags
inspectCmd.Flags().Bool("binary", false, "Show hex dump with field annotations (NP3 only)")
inspectCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
inspectCmd.Flags().Bool("pretty", true, "Pretty-print JSON (default: true)")
```

**2. Diff Command**
```bash
# Compare two files (same or different formats)
recipe diff original.np3 converted.xmp

# Unified output (shows all fields, not just changes)
recipe diff --unified portrait1.xmp portrait2.xmp

# JSON output for programmatic parsing
recipe diff --format=json original.lrtemplate converted.lrtemplate

# Tolerance for float comparisons (default: 0.001)
recipe diff --tolerance=0.01 file1.xmp file2.xmp
```

**Cobra Command Signature:**
```go
var diffCmd = &cobra.Command{
    Use:   "diff [file1] [file2]",
    Short: "Compare parameters between two preset files",
    Long:  `Diff compares all parameters between two preset files.
Files can be different formats (e.g., NP3 vs XMP). Shows only changed parameters by default.`,
    Args:  cobra.ExactArgs(2),
    RunE:  runDiff,
}

// Flags
diffCmd.Flags().Bool("unified", false, "Show all fields, not just changes")
diffCmd.Flags().String("format", "text", "Output format: text or json")
diffCmd.Flags().Float64("tolerance", 0.001, "Tolerance for float comparisons")
diffCmd.Flags().Bool("no-color", false, "Disable colored output")
```

**Internal Package APIs:**

**internal/inspect/inspect.go**
```go
// ToJSON serializes a UniversalRecipe to pretty-printed JSON
func ToJSON(recipe *model.UniversalRecipe) ([]byte, error)

// ToJSONWithMetadata adds file metadata wrapper
func ToJSONWithMetadata(recipe *model.UniversalRecipe, sourceFile, format string) ([]byte, error)
```

**internal/inspect/binary.go**
```go
// BinaryDump creates annotated hex dump (NP3 only)
func BinaryDump(data []byte, format string) (string, error)
```

**internal/inspect/diff.go**
```go
// Diff compares two UniversalRecipe structs
func Diff(recipe1, recipe2 *model.UniversalRecipe) ([]DiffResult, error)

// FormatDiff renders diff results as human-readable text
func FormatDiff(results []DiffResult, colorize bool) string

// FormatDiffJSON renders diff results as JSON
func FormatDiffJSON(results []DiffResult) ([]byte, error)
```

### Workflows and Sequencing

**Workflow 1: Inspect Command (JSON Mode)**
```
User: recipe inspect portrait.np3
  ↓
1. CLI parses command and flags
  ↓
2. Read file: os.ReadFile("portrait.np3")
  ↓
3. Detect format: detectFormat("portrait.np3") → "np3"
  ↓
4. Parse file: np3.Parse(data) → *UniversalRecipe
  ↓
5. Convert to JSON: inspect.ToJSONWithMetadata(recipe, "portrait.np3", "np3")
  ↓
6. Output to stdout (or file if --output specified)
  ↓
Result: Pretty-printed JSON with all parameters
```

**Workflow 2: Inspect Command (Binary Mode)**
```
User: recipe inspect portrait.np3 --binary
  ↓
1. CLI parses command and flags
  ↓
2. Read file: os.ReadFile("portrait.np3")
  ↓
3. Detect format: detectFormat("portrait.np3") → "np3"
  ↓
4. Validate NP3 format (error if XMP/lrtemplate)
  ↓
5. Generate hex dump: inspect.BinaryDump(data, "np3")
  ↓
6. Output annotated hex dump to stdout
  ↓
Result: Hex dump with field names and values
Example:
[0x0000] 4E 50                Magic ("NP")
[0x0002] 03 00                Version (3)
[0x0042] 80                   Contrast (0, normalized from 128)
[0x0043] 80                   Brightness (0, normalized from 128)
```

**Workflow 3: Diff Command**
```
User: recipe diff original.np3 converted.xmp
  ↓
1. CLI parses command and flags
  ↓
2. Read both files
  ↓
3. Detect formats: "np3", "xmp"
  ↓
4. Parse file1: np3.Parse(data1) → recipe1
  ↓
5. Parse file2: xmp.Parse(data2) → recipe2
  ↓
6. Compare: inspect.Diff(recipe1, recipe2) → []DiffResult
  ↓
7. Format output:
   - If --format=json: inspect.FormatDiffJSON(results)
   - Else: inspect.FormatDiff(results, !noColor)
  ↓
8. Output to stdout
  ↓
Result: Comparison showing changed parameters
Example (text mode):
Contrast: 0 → 15
Saturation: 0 → -10
Highlights: -50 → -45

No changes: Exposure, Shadows, Whites, Blacks, ...
```

**Error Handling Sequences:**

**Invalid File Format:**
```
User: recipe inspect invalid.txt
  ↓
detectFormat() → error: "unknown format"
  ↓
Output: Error: unable to detect format for 'invalid.txt'
Supported formats: .np3, .xmp, .lrtemplate
```

**Binary Mode on Non-NP3:**
```
User: recipe inspect portrait.xmp --binary
  ↓
detectFormat() → "xmp"
  ↓
BinaryDump(data, "xmp") → error: "binary mode only supports NP3"
  ↓
Output: Error: --binary flag only works with NP3 files
XMP and lrtemplate are text-based formats
```

## Non-Functional Requirements

### Performance

**Target Metrics:**

| Operation                      | Target | Rationale                                              |
| ------------------------------ | ------ | ------------------------------------------------------ |
| `recipe inspect FILE`          | <50ms  | Faster than conversion (<100ms), only parsing required |
| `recipe inspect FILE --binary` | <10ms  | Simple byte-to-string formatting, no parsing           |
| `recipe diff FILE1 FILE2`      | <100ms | Two parses + comparison, same budget as one conversion |
| JSON serialization             | <5ms   | Standard library json.Marshal() is very fast           |
| Memory usage (inspect)         | <10 MB | Single UniversalRecipe in memory                       |
| Memory usage (diff)            | <20 MB | Two UniversalRecipe structs in memory                  |

**Performance Optimizations:**

1. **Reuse Existing Parsers** - No additional parsing overhead, leverages Epic 1 optimizations
2. **Streaming Output** - Write JSON/diff directly to stdout/file (no intermediate buffering)
3. **Lazy Binary Annotation** - Only annotate known NP3 fields, skip unknown bytes
4. **Smart Diff** - Only compare fields that exist in both recipes (skip nil checks where possible)
5. **No Pretty-Printing Overhead** - Use json.MarshalIndent() which is optimized in Go stdlib

**Benchmarking Strategy:**
```go
func BenchmarkInspectJSON(b *testing.B) {
    data, _ := os.ReadFile("testdata/np3/portrait.np3")
    recipe, _ := np3.Parse(data)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := inspect.ToJSON(recipe)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**Expected Results:**
- Inspect: ~20-30ms (15ms parse + 5ms JSON)
- Binary: ~2-5ms (just hex formatting)
- Diff: ~50-70ms (30ms two parses + 5ms compare + 5ms format)

### Security

**Input Validation:**
- All file reading uses same validation as Epic 1 parsers (magic bytes, size limits, format checks)
- No new attack vectors introduced (read-only operations only)
- Binary mode does not execute any code, just displays hex bytes
- JSON output is properly escaped (Go's json.Marshal handles this automatically)

**Specific Security Considerations:**

1. **File Size Limits** (inherited from Epic 1)
   - Maximum file size: 10 MB (prevents memory exhaustion)
   - Enforced in parser before passing to inspect functions

2. **Path Traversal Prevention**
   - CLI accepts only file paths provided by user (no automatic directory traversal)
   - No wildcard expansion in inspect/diff commands
   - Output file paths validated (no overwriting system files)

3. **JSON Injection Prevention**
   - All JSON output uses json.Marshal() (auto-escapes special characters)
   - No string concatenation or manual JSON construction
   - Prevents XSS if JSON displayed in web context

4. **Binary Dump Safety**
   - Binary mode only displays data, does not interpret or execute
   - No buffer overflows (Go's bounds checking)
   - Hex output is pure text (no escape sequences that could affect terminal)

5. **Diff Output Safety**
   - Diff output is sanitized text only
   - No execution of shell commands
   - ANSI color codes only used if terminal supports them (--no-color flag available)

**Privacy:**
- Same privacy-first design as Epic 1 (all processing local, no network access)
- No telemetry, no logging of file contents
- Inspect/diff output goes to stdout or user-specified file only

### Reliability/Availability

**Error Handling:**

All commands follow Epic 1's error handling pattern with ConversionError types:

```go
// Example: inspect command error handling
func runInspect(cmd *cobra.Command, args []string) error {
    inputPath := args[0]

    // Step 1: Read file
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    // Step 2: Detect format
    format, err := detectFormat(inputPath)
    if err != nil {
        return fmt.Errorf("unable to detect format: %w", err)
    }

    // Step 3: Parse (may return ConversionError)
    recipe, err := parseFile(data, format)
    if err != nil {
        var convErr *converter.ConversionError
        if errors.As(err, &convErr) {
            return fmt.Errorf("parse failed (%s): %w", convErr.Format, convErr.Cause)
        }
        return err
    }

    // Step 4: Convert to JSON
    output, err := inspect.ToJSONWithMetadata(recipe, inputPath, format)
    if err != nil {
        return fmt.Errorf("JSON serialization failed: %w", err)
    }

    // Step 5: Output
    fmt.Println(string(output))
    return nil
}
```

**Graceful Degradation:**

- If binary mode fails on corrupt NP3, fall back to raw hex dump without annotations
- If diff finds incomparable fields (e.g., nil vs value), mark as "added" or "removed" instead of crashing
- If JSON output fails, provide error message suggesting file might be corrupt

**Availability:**

- CLI commands are fully offline (no network dependencies)
- No external services required
- Works on any system with Go runtime (Windows, macOS, Linux)
- Exit codes follow Unix conventions (0 = success, 1 = error)

### Observability

**Logging Strategy:**

Epic 5 uses same slog-based logging as Epic 3 CLI:

```go
import "log/slog"

// Inspect command logging
logger.Info("inspecting file",
    slog.String("file", inputPath),
    slog.String("format", format),
    slog.Bool("binary_mode", binaryMode),
)

logger.Info("inspection complete",
    slog.String("file", inputPath),
    slog.Duration("elapsed", time.Since(start)),
    slog.Int("json_size_bytes", len(output)),
)

// Diff command logging
logger.Info("comparing files",
    slog.String("file1", file1),
    slog.String("file2", file2),
    slog.String("format1", format1),
    slog.String("format2", format2),
)

logger.Info("diff complete",
    slog.Int("changes_found", len(diffResults)),
    slog.Int("significant_changes", countSignificant(diffResults)),
)
```

**Verbose Mode:**

When `--verbose` flag is set (from Epic 3):

```bash
recipe inspect portrait.np3 --verbose
```

Additional debug logging:
```go
logger.Debug("parsed recipe",
    slog.String("name", recipe.Name),
    slog.Int("contrast", recipe.Contrast),
    slog.Int("saturation", recipe.Saturation),
    slog.Int("non_zero_fields", countNonZero(recipe)),
)

logger.Debug("binary field map loaded",
    slog.Int("known_fields", len(np3FieldMap)),
)
```

**Error Logging:**

```go
logger.Error("inspection failed",
    slog.String("file", inputPath),
    slog.String("error", err.Error()),
    slog.String("format", format),
)
```

**Metrics (via --json-output mode in diff):**

Diff command with `--format=json` outputs structured data suitable for monitoring:

```json
{
  "file1": "original.np3",
  "file2": "converted.xmp",
  "changes": 3,
  "significant_changes": 1,
  "fields_compared": 50,
  "differences": [...]
}
```

This enables scripting/automation to track conversion accuracy over time.

## Dependencies and Integrations

**External Dependencies:**

Epic 5 introduces no new external dependencies beyond what Epic 3 already requires:

| Dependency             | Version | Purpose                | Used By                                      |
| ---------------------- | ------- | ---------------------- | -------------------------------------------- |
| github.com/spf13/cobra | Latest  | CLI framework          | inspect.go, diff.go commands                 |
| Go standard library    | 1.24+   | All core functionality | JSON marshaling, file I/O, string formatting |

**Internal Dependencies (from existing epics):**

```
internal/inspect/          (Epic 5 - NEW)
    ↓ depends on
internal/model/           (Epic 1)
    ↓ UniversalRecipe struct
    
internal/inspect/          (Epic 5 - NEW)
    ↓ depends on
internal/formats/np3/     (Epic 1)
internal/formats/xmp/     (Epic 1)
internal/formats/lrtemplate/ (Epic 1)
    ↓ Parse() functions

cmd/cli/inspect.go        (Epic 5 - NEW)
cmd/cli/diff.go           (Epic 5 - NEW)
    ↓ depend on
cmd/cli/main.go           (Epic 3)
    ↓ Cobra root command
```

**Integration Points:**

1. **Epic 1 Integration (Core Conversion Engine)**
   - Uses `np3.Parse()`, `xmp.Parse()`, `lrtemplate.Parse()` without modification
   - Operates on `model.UniversalRecipe` struct
   - No changes required to Epic 1 code

2. **Epic 3 Integration (CLI Interface)**
   - Adds two new commands to existing Cobra CLI structure
   - Reuses format detection logic from Epic 3
   - Inherits error handling patterns and logging setup
   - Shares `--verbose` flag behavior

3. **No Integration with Epic 2 (Web Interface)**
   - Inspection tools are CLI-only in this epic
   - Web-based inspection could be added in future epic if needed

4. **No Integration with Epic 4 (TUI Interface)**
   - TUI could potentially call inspect functions for parameter display
   - Not in scope for Epic 5 (TUI uses direct UniversalRecipe access)

**Build System Integration:**

```makefile
# Makefile additions for Epic 5 (no changes needed, inspect/diff built with CLI)
cli:
	go build -o recipe cmd/cli/main.go  # Includes inspect.go and diff.go

test:
	go test ./internal/inspect/  # New test package
	go test ./...
```

**Version Requirements:**
- Go 1.24+ (same as Epic 1 requirement for WASM support)
- No minimum version change
- Compatible with all platforms (Windows, macOS, Linux)

## Acceptance Criteria (Authoritative)

Extracted from PRD FR-5 and normalized into testable statements:

### AC-1: Parameter Inspection Tool (FR-5.1)

**AC-1.1:** `recipe inspect FILE` command outputs complete parameter set as valid JSON
- All UniversalRecipe fields present in output
- JSON is parseable by standard JSON parsers
- Pretty-printed with indentation for readability

**AC-1.2:** Inspect output includes metadata wrapper with source file info
- `metadata.source_file`: Original filename
- `metadata.source_format`: Detected format ("np3", "xmp", "lrtemplate")
- `metadata.parsed_at`: ISO 8601 timestamp
- `metadata.recipe_version`: Recipe tool version

**AC-1.3:** Inspect command supports all three formats (NP3, XMP, lrtemplate)
- Successfully parses and outputs JSON for all formats
- No format-specific errors when used correctly

**AC-1.4:** `--output FILE` flag writes JSON to specified file instead of stdout
- File created with correct permissions (0644)
- Overwrites existing files with confirmation
- Creates parent directories if needed

### AC-2: Binary Structure Visualization (FR-5.2)

**AC-2.1:** `recipe inspect FILE --binary` shows hex dump with byte offsets
- Format: `[0xOFFSET] HEX_BYTES  FIELD_NAME (VALUE)`
- Example: `[0x0000] 4E 50  Magic ("NP")`
- All bytes displayed in hex format

**AC-2.2:** Binary mode labels all known NP3 fields from reverse engineering
- Minimum fields annotated: Magic, Version, Contrast, Brightness, Saturation, Hue, Sharpness
- Unknown bytes shown without labels (raw hex only)
- Field values shown in human-readable form (e.g., normalized ranges)

**AC-2.3:** Binary mode only works with NP3 files
- Returns clear error if used with XMP or lrtemplate
- Error message explains that XMP/lrtemplate are text formats

### AC-3: Diff Tool (FR-5.3)

**AC-3.1:** `recipe diff FILE1 FILE2` shows added/removed/changed parameters
- Changed parameters: `FieldName: old_value → new_value`
- Added parameters: `FieldName: (none) → new_value`
- Removed parameters: `FieldName: old_value → (none)`
- Unchanged parameters listed in summary line

**AC-3.2:** Diff works across different formats (cross-format comparison)
- Example: Can compare NP3 to XMP, XMP to lrtemplate, etc.
- Parameters correctly mapped via UniversalRecipe intermediate
- No false positives from format differences

**AC-3.3:** Diff highlights significant differences (>5% change for numeric values)
- Significant changes marked with indicator (e.g., `*` or color)
- Minor changes (<5%) shown but not highlighted
- Tolerance configurable via `--tolerance` flag

**AC-3.4:** `--format=json` flag outputs diff as structured JSON
- JSON contains array of DiffResult objects
- Each result has: field, old_value, new_value, change_type, significant
- Parseable for automation/scripting

**AC-3.5:** `--unified` flag shows all fields, not just changes
- Displays complete parameter comparison
- Unchanged fields shown as: `FieldName: value = value`

### AC-4: Error Handling and User Experience

**AC-4.1:** Clear error messages for unsupported file formats
- Message identifies the problematic file
- Lists supported formats (.np3, .xmp, .lrtemplate)

**AC-4.2:** Parse errors return format-specific context
- Uses ConversionError type from Epic 1
- Error message includes: operation, format, underlying cause

**AC-4.3:** Binary mode error on non-NP3 files is user-friendly
- Explains that binary mode only supports NP3
- Suggests using JSON mode instead for text formats

### AC-5: Performance Requirements

**AC-5.1:** `recipe inspect FILE` completes in <50ms
- Measured with Go benchmarks on reference hardware
- Includes file read, parse, and JSON serialization

**AC-5.2:** `recipe inspect FILE --binary` completes in <10ms
- Hex formatting is lightweight operation
- No parsing overhead in binary mode

**AC-5.3:** `recipe diff FILE1 FILE2` completes in <100ms
- Two file parses + comparison + formatting
- Same performance budget as single conversion

## Traceability Mapping

Maps acceptance criteria to technical design and test coverage:

| AC         | PRD Requirement                        | Spec Section                         | Component/API                                      | Test Coverage                        |
| ---------- | -------------------------------------- | ------------------------------------ | -------------------------------------------------- | ------------------------------------ |
| **AC-1.1** | FR-5.1: Complete parameter set as JSON | APIs & Interfaces → Inspect Command  | `inspect.ToJSON()`                                 | `TestToJSON_AllFields()`             |
| **AC-1.2** | FR-5.1: Includes metadata              | Data Models → InspectOutput          | `inspect.ToJSONWithMetadata()`                     | `TestToJSONWithMetadata()`           |
| **AC-1.3** | FR-5.1: Supports all formats           | APIs & Interfaces → Inspect Command  | `np3.Parse()`, `xmp.Parse()`, `lrtemplate.Parse()` | `TestInspect_AllFormats()`           |
| **AC-1.4** | FR-5.1: Output flag                    | APIs & Interfaces → Inspect Command  | `inspectCmd.Flags()`                               | `TestInspect_OutputFile()`           |
| **AC-2.1** | FR-5.2: Hex dump with offsets          | APIs & Interfaces → Inspect Binary   | `inspect.BinaryDump()`                             | `TestBinaryDump_Format()`            |
| **AC-2.2** | FR-5.2: Labels known fields            | Data Models → Binary Field Map       | `np3FieldMap` constant                             | `TestBinaryDump_KnownFields()`       |
| **AC-2.3** | FR-5.2: NP3 only                       | Workflows → Error Handling           | Binary mode validation                             | `TestBinaryDump_NonNP3Error()`       |
| **AC-3.1** | FR-5.3: Shows changes                  | APIs & Interfaces → Diff Command     | `inspect.Diff()`                                   | `TestDiff_DetectsChanges()`          |
| **AC-3.2** | FR-5.3: Cross-format                   | System Architecture → Reuses Parsers | UniversalRecipe intermediate                       | `TestDiff_CrossFormat()`             |
| **AC-3.3** | FR-5.3: Highlights significant         | Data Models → DiffResult.Significant | Threshold comparison                               | `TestDiff_SignificantChanges()`      |
| **AC-3.4** | FR-5.3: JSON output                    | APIs & Interfaces → Diff Command     | `inspect.FormatDiffJSON()`                         | `TestDiff_JSONFormat()`              |
| **AC-3.5** | FR-5.3: Unified mode                   | APIs & Interfaces → Diff Command     | `--unified` flag                                   | `TestDiff_UnifiedMode()`             |
| **AC-4.1** | Error handling                         | NFR: Reliability → Error Handling    | Format detection                                   | `TestInspect_InvalidFormat()`        |
| **AC-4.2** | Error handling                         | NFR: Reliability → Error Handling    | ConversionError wrapping                           | `TestInspect_ParseError()`           |
| **AC-4.3** | Error handling                         | NFR: Reliability → Error Handling    | Binary mode validation                             | `TestBinaryDump_UserFriendlyError()` |
| **AC-5.1** | Performance                            | NFR: Performance                     | `inspect.ToJSON()`                                 | `BenchmarkInspectJSON()`             |
| **AC-5.2** | Performance                            | NFR: Performance                     | `inspect.BinaryDump()`                             | `BenchmarkBinaryDump()`              |
| **AC-5.3** | Performance                            | NFR: Performance                     | `inspect.Diff()`                                   | `BenchmarkDiff()`                    |

**Coverage Summary:**
- **Functional Coverage**: 15 ACs mapped to 15 test cases
- **Component Coverage**: All 3 new modules (`inspect.go`, `binary.go`, `diff.go`) covered
- **Integration Coverage**: Epic 1 parser integration, Epic 3 CLI integration
- **Performance Coverage**: 3 benchmark tests for all primary operations

## Risks, Assumptions, Open Questions

### Risks

**R-1: NP3 Field Map Completeness** (Medium)
- **Description**: Binary field map may not cover all NP3 byte positions
- **Impact**: Unknown fields show as unlabeled hex bytes (still functional, but less educational)
- **Mitigation**: Start with known fields from Epic 1 reverse engineering, expand based on user feedback
- **Fallback**: Raw hex dump still provides value even without full annotation

**R-2: JSON Output Size for Large Metadata** (Low)
- **Description**: Some XMP files have extensive metadata that could bloat JSON output
- **Impact**: Stdout output may be unwieldy, file sizes larger than expected
- **Mitigation**: Metadata dictionary in UniversalRecipe already filters to relevant fields
- **Fallback**: Users can pipe output to file or use `--output` flag

**R-3: Diff Tolerance Tuning** (Low)
- **Description**: Default tolerance (0.001) may not suit all use cases
- **Impact**: False positives (noise from rounding differences) or false negatives (missing real changes)
- **Mitigation**: Make tolerance configurable via `--tolerance` flag, document recommended values
- **Fallback**: Users can experiment with tolerance values for their workflows

### Assumptions

**A-1: Epic 1 Parsers Are Stable**
- Epic 5 assumes Epic 1's parsers are production-ready and thoroughly tested
- Any parser bugs will surface in inspection output
- Validation: Epic 1 has 1,479 sample files tested (high confidence)

**A-2: CLI Interface (Epic 3) Is Implemented**
- Epic 5 assumes Cobra CLI structure from Epic 3 exists
- Dependencies on root command setup, flag patterns, error handling
- Validation: Epic 3 tech spec defines complete CLI architecture

**A-3: UniversalRecipe Struct Is Comprehensive**
- Epic 5 assumes UniversalRecipe captures all convertible parameters
- JSON output completeness depends on this struct's design
- Validation: Epic 1 defines UniversalRecipe with 50+ fields covering all formats

**A-4: Users Have Basic CLI Knowledge**
- Assumes users understand stdout/stderr, file redirection, flag syntax
- Documentation will provide examples for common use cases
- Validation: Target audience includes developers and power users (appropriate assumption)

### Open Questions

**Q-1: Should binary mode support partial annotation?**
- Current design: Annotate known fields only, show raw hex for unknown bytes
- Alternative: Skip unknown bytes entirely (cleaner output but less complete)
- **Resolution Needed**: User research - do users want completeness or clarity?
- **Proposal**: Keep current design, add `--known-only` flag to hide unknown bytes

**Q-2: Should diff support side-by-side output?**
- Current design: Vertical list format (field: old → new)
- Alternative: Two-column side-by-side comparison
- **Resolution Needed**: Terminal width constraints, readability testing
- **Proposal**: Defer to future enhancement, start with simple vertical format

**Q-3: Should inspect support output format selection (YAML, TOML)?**
- Current design: JSON only
- Alternative: Support `--format=yaml` or `--format=toml`
- **Resolution Needed**: Assess demand, additional dependencies required
- **Proposal**: Start with JSON (zero dependencies), add formats if requested

**Q-4: Should diff support ignore lists (e.g., --ignore-metadata)?**
- Current design: Compare all fields in UniversalRecipe
- Alternative: Allow users to exclude certain fields from comparison
- **Resolution Needed**: Common use cases for partial diff
- **Proposal**: Defer to future enhancement, accept feedback on which fields to ignore

### Decision Log

**D-1: Binary Mode NP3-Only** ✅ Decided
- **Decision**: Binary visualization only for NP3 (not XMP/lrtemplate)
- **Rationale**: XMP and lrtemplate are text formats (users can view with text editor)
- **Impact**: Simpler implementation, clearer error messages
- **Date**: 2025-11-06

**D-2: JSON as Default Output** ✅ Decided
- **Decision**: Inspect command outputs JSON by default (no `--format` flag needed)
- **Rationale**: JSON is universal, parseable by all languages, zero dependencies
- **Impact**: Simple CLI interface, wide compatibility
- **Date**: 2025-11-06

**D-3: Diff Shows Changes Only by Default** ✅ Decided
- **Decision**: Default diff output shows only changed fields, `--unified` for all fields
- **Rationale**: Most users care about differences, not unchanged fields
- **Impact**: Cleaner output, faster scanning for changes
- **Date**: 2025-11-06

## Test Strategy Summary

### Test Levels

**Unit Tests** (Primary)
- Test each function in `internal/inspect/` package independently
- Mock UniversalRecipe structs for predictable test cases
- Validate JSON serialization, binary formatting, diff logic

**Integration Tests**
- Test CLI commands end-to-end (file input → stdout output)
- Verify integration with Epic 1 parsers
- Validate error handling across module boundaries

**Performance Tests** (Benchmarks)
- Measure inspect, binary dump, and diff operations
- Validate <50ms, <10ms, <100ms targets respectively
- Run on CI to catch performance regressions

### Test Organization

```
internal/inspect/
├── inspect_test.go        # ToJSON, ToJSONWithMetadata tests
├── binary_test.go         # BinaryDump tests
├── diff_test.go           # Diff, FormatDiff tests
└── testdata/
    ├── sample.np3         # Known-good NP3 for testing
    ├── sample.xmp         # Known-good XMP for testing
    └── expected_*.json    # Golden files for JSON comparison

cmd/cli/
├── inspect_test.go        # End-to-end inspect command tests
└── diff_test.go           # End-to-end diff command tests
```

### Key Test Cases

**Inspect Command Tests:**
```go
TestToJSON_AllFields              // Verify all UniversalRecipe fields present
TestToJSONWithMetadata           // Verify metadata wrapper structure
TestInspect_AllFormats           // NP3, XMP, lrtemplate all work
TestInspect_OutputFile           // --output flag creates file
TestInspect_InvalidFormat        // Error on unknown format
TestInspect_ParseError           // ConversionError handling
```

**Binary Dump Tests:**
```go
TestBinaryDump_Format            // Verify [offset] hex_bytes format
TestBinaryDump_KnownFields       // All NP3 fields annotated
TestBinaryDump_NonNP3Error       // Error on XMP/lrtemplate
TestBinaryDump_CorruptFile       // Graceful degradation
```

**Diff Tests:**
```go
TestDiff_DetectsChanges          // Identifies modified/added/removed
TestDiff_CrossFormat             // NP3 vs XMP comparison
TestDiff_SignificantChanges      // >5% threshold marking
TestDiff_JSONFormat              // --format=json output
TestDiff_UnifiedMode             // --unified shows all fields
TestDiff_Tolerance               // --tolerance flag works
```

**Performance Benchmarks:**
```go
BenchmarkInspectJSON             // Target: <50ms total
BenchmarkBinaryDump              // Target: <10ms
BenchmarkDiff                    // Target: <100ms
```

### Test Data Strategy

**Reuse Epic 1 Sample Files:**
- Use existing testdata/np3/*.np3 files (22 files)
- Use existing testdata/xmp/*.xmp files (913 files)
- Use existing testdata/lrtemplate/*.lrtemplate files (544 files)
- Ensures consistency with conversion testing

**Golden File Testing:**
- Store expected JSON output in testdata/expected_*.json
- Compare actual output to golden files byte-for-byte
- Regenerate golden files when UniversalRecipe struct changes

**Synthetic Test Cases:**
- Create minimal UniversalRecipe structs for unit tests
- Test edge cases (all zeros, all max values, nil fields)
- Test error paths (corrupt data, invalid JSON)

### CI/CD Integration

```yaml
# .github/workflows/test.yml
test-inspect:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.24'
    - run: go test ./internal/inspect/ -v
    - run: go test ./cmd/cli/ -run TestInspect -v
    - run: go test ./cmd/cli/ -run TestDiff -v
    - run: go test -bench=. ./internal/inspect/
```

### Coverage Targets

- **Unit Test Coverage**: ≥90% for `internal/inspect/` package
- **Integration Test Coverage**: 100% of CLI commands (inspect, diff)
- **Error Path Coverage**: 100% of error handling branches
- **Performance Coverage**: All 3 primary operations benchmarked

### Manual Testing Checklist

Pre-release validation by human tester:

- [ ] `recipe inspect` works on all three formats
- [ ] JSON output is valid (test with `jq` or `python -m json.tool`)
- [ ] `--binary` flag shows readable hex dump for NP3
- [ ] `--binary` flag errors appropriately on XMP/lrtemplate
- [ ] `recipe diff` correctly identifies changes
- [ ] Cross-format diff works (e.g., NP3 vs XMP)
- [ ] `--format=json` outputs parseable JSON
- [ ] Performance meets targets (<50ms inspect, <10ms binary, <100ms diff)
- [ ] Error messages are clear and actionable
- [ ] Help text (`--help`) is accurate

### Acceptance Test Execution

Each AC from "Acceptance Criteria (Authoritative)" section maps to automated test:

| AC     | Test Function                    | Pass Criteria                                 |
| ------ | -------------------------------- | --------------------------------------------- |
| AC-1.1 | TestToJSON_AllFields             | All 50+ UniversalRecipe fields in JSON        |
| AC-1.2 | TestToJSONWithMetadata           | Metadata section present with 4 fields        |
| AC-1.3 | TestInspect_AllFormats           | No errors on NP3, XMP, lrtemplate             |
| AC-1.4 | TestInspect_OutputFile           | File created, content matches stdout          |
| AC-2.1 | TestBinaryDump_Format            | Output matches `[0xOFFSET] HEX` pattern       |
| AC-2.2 | TestBinaryDump_KnownFields       | Min 7 fields annotated (Magic, Version, etc.) |
| AC-2.3 | TestBinaryDump_NonNP3Error       | Error message mentions "NP3 only"             |
| AC-3.1 | TestDiff_DetectsChanges          | Modified/added/removed all identified         |
| AC-3.2 | TestDiff_CrossFormat             | NP3→XMP diff finds correct matches            |
| AC-3.3 | TestDiff_SignificantChanges      | >5% changes marked, <5% not marked            |
| AC-3.4 | TestDiff_JSONFormat              | JSON output is valid, parseable               |
| AC-3.5 | TestDiff_UnifiedMode             | Unchanged fields shown in output              |
| AC-4.1 | TestInspect_InvalidFormat        | Error lists .np3/.xmp/.lrtemplate             |
| AC-4.2 | TestInspect_ParseError           | ConversionError type returned                 |
| AC-4.3 | TestBinaryDump_UserFriendlyError | Message explains text format issue            |
| AC-5.1 | BenchmarkInspectJSON             | <50ms on reference hardware                   |
| AC-5.2 | BenchmarkBinaryDump              | <10ms on reference hardware                   |
| AC-5.3 | BenchmarkDiff                    | <100ms on reference hardware                  |
