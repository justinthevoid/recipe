# Epic Technical Specification: Core Conversion Engine

Date: 2025-11-04
Author: Justin
Epic ID: epic-1
Status: Draft

---

## Overview

Epic 1 implements the core conversion engine for Recipe, a universal photo preset converter. This epic establishes the foundation for all format conversions using a hub-and-spoke architecture with a universal intermediate representation (UniversalRecipe). The conversion engine supports three format families:
- Nikon Picture Control (.np3) - Binary format with 5 core parameters
- Adobe Lightroom CC (.xmp) - XML/RDF format with 50+ parameters
- Lightroom Classic (.lrtemplate) - Lua table format with same parameters as XMP

The engine provides a single API (`converter.Convert()`) that all interfaces (CLI, TUI, Web) will use, ensuring consistent behavior and maintainability. All processing is stateless and uses only Go standard library (zero dependencies), enabling WASM compilation for client-side web execution.

**Core Value Proposition**: Bidirectional conversion between any format pair with ≥95% accuracy, <100ms conversion time, and graceful handling of unmappable parameters.

## Objectives and Scope

### Primary Objectives
1. **Format Support**: Parse and generate valid NP3, XMP, and lrtemplate files that open correctly in their respective applications
2. **Bidirectional Conversion**: Support all conversion paths (NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate) with parameter mapping rules
3. **Accuracy**: Achieve ≥95% visual similarity for common adjustments through round-trip testing with 1,501 sample files
4. **Performance**: Complete conversions in <100ms with <5ms overhead from intermediate representation
5. **Error Transparency**: Report unmappable features to users with no silent data loss

### In Scope
- Universal intermediate representation (UniversalRecipe struct) supporting superset of all format capabilities
- NP3 binary parser/generator with magic number validation
- XMP XML/RDF parser/generator with namespace preservation
- lrtemplate Lua table parser/generator with syntax compatibility
- Parameter mapping rules with direct mapping and approximation strategies
- Comprehensive error handling with custom ConversionError type
- Table-driven tests using 1,501 real sample files (22 NP3, 913 XMP, 566 lrtemplate)
- Round-trip conversion validation (format A → B → A produces identical output)

### Out of Scope
- Batch processing (handled by CLI/TUI interfaces in Epic 3/4)
- UI for parameter preview (handled by Web interface in Epic 2)
- Additional format support beyond NP3, XMP, lrtemplate (future epics)
- Server-side conversion (architecture is client-side/local only)
- Real-time preview rendering (visualization tools in Epic 5)

### Success Criteria
- All 1,501 sample files parse without errors
- Round-trip conversions produce binary-identical output (tolerance ±1 for rounding)
- Generated files open correctly in Nikon NX Studio, Lightroom CC, and Lightroom Classic
- Conversion API completes in <100ms (validated via Go benchmarks)
- Zero dependencies beyond Go standard library
- Test coverage ≥90% for all format packages

## System Architecture Alignment

Epic 1 implements the **hub-and-spoke pattern** as defined in the system architecture:

```
       ┌─────────────┐
       │ NP3 Parser  │
       └──────┬──────┘
              │
              ▼
       ┌──────────────┐        ┌─────────────┐
       │              │◄───────┤ XMP Parser  │
       │ Universal    │        └─────────────┘
       │   Recipe     │
       │   (Hub)      │        ┌─────────────────┐
       │              │◄───────┤ lrtemplate      │
       └──────────────┘        │ Parser          │
              │                └─────────────────┘
              ▼
       ┌──────────────┐
       │ NP3 Generator│
       └──────────────┘
       ┌──────────────┐
       │ XMP Generator│
       └──────────────┘
       ┌──────────────┐
       │ lrtemplate   │
       │ Generator    │
       └──────────────┘
```

### Architecture Alignment Points

1. **Project Structure** (aligns with Section 3: Project Structure)
   - `internal/converter/` - Contains Convert() API and ConversionError type
   - `internal/formats/{np3,xmp,lrtemplate}/` - Each format has parse.go, generate.go, {format}_test.go
   - `internal/model/` - Contains UniversalRecipe struct definition

2. **Single API Rule** (aligns with Pattern 10: Consistency Rules)
   - All interfaces MUST call `converter.Convert(input []byte, from, to string) ([]byte, error)`
   - No direct parser calls allowed (enforces hub-and-spoke pattern)

3. **Error Handling** (aligns with Pattern 5: Error Handling)
   - All conversion failures wrapped in `ConversionError{Operation, Format, Cause}`
   - Type-safe error checking, format-specific context

4. **Testing Strategy** (aligns with Pattern 7: Testing Strategy)
   - Table-driven tests with real sample files from `testdata/{np3,xmp,lrtemplate}/`
   - Round-trip tests validate accuracy goal (≥95%)

5. **WASM Compatibility** (aligns with Pattern 9: WASM Export Pattern)
   - Stateless design (no OS dependencies, no file I/O in core library)
   - Only uses Go stdlib (works in WASM environment)
   - Will be exported via `go:wasmexport` in cmd/wasm/main.go (Epic 2)

6. **Zero Dependencies** (aligns with Decision Summary: Critical Decisions)
   - Only `encoding/xml` for XMP parsing
   - Only `encoding/binary` for NP3 binary parsing
   - No external libraries for Lua parsing (custom parser using regex)

## Detailed Design

### Services and Modules

#### Module: internal/converter
**Responsibility**: Orchestrate conversion between formats using hub-and-spoke pattern

**Public API**:
```go
package converter

// Convert transforms input bytes from source format to target format
// Supported formats: "np3", "xmp", "lrtemplate"
// Returns converted bytes or ConversionError on failure
func Convert(input []byte, from, to string) ([]byte, error)
```

**Internal Flow**:
```go
1. Validate format strings (from, to ∈ {"np3", "xmp", "lrtemplate"})
2. Call appropriate parser:
   - from="np3"        → np3.Parse(input)
   - from="xmp"        → xmp.Parse(input)
   - from="lrtemplate" → lrtemplate.Parse(input)
3. Receive UniversalRecipe struct (hub)
4. Call appropriate generator:
   - to="np3"        → np3.Generate(recipe)
   - to="xmp"        → xmp.Generate(recipe)
   - to="lrtemplate" → lrtemplate.Generate(recipe)
5. Return generated bytes
6. Wrap all errors in ConversionError for context
```

**Error Handling**:
```go
// ConversionError provides context for conversion failures
type ConversionError struct {
    Operation string  // "parse", "generate", "validate"
    Format    string  // "np3", "xmp", "lrtemplate"
    Cause     error   // Underlying error
}

func (e *ConversionError) Error() string {
    return fmt.Sprintf("%s %s: %v", e.Operation, e.Format, e.Cause)
}

func (e *ConversionError) Unwrap() error {
    return e.Cause
}
```

#### Module: internal/formats/np3
**Responsibility**: Parse and generate Nikon Picture Control (.np3) binary files

**File Structure**:
```
internal/formats/np3/
├── parse.go       # Parse([]byte) (*model.UniversalRecipe, error)
├── generate.go    # Generate(*model.UniversalRecipe) ([]byte, error)
└── np3_test.go    # Table-driven tests with 22 sample files
```

**NP3 Format Specification**:
- File size: 1024 bytes (fixed)
- Magic bytes: First 4 bytes identify NP3 format
- Parameters (5 core):
  - Sharpening: 1 byte (0-9)
  - Contrast: 1 byte (-3 to +3, stored as 0-6)
  - Brightness: 1 byte (-1 to +1, stored as 0-2)
  - Saturation: 1 byte (-3 to +3, stored as 0-6)
  - Hue: 1 byte (-9° to +9°, stored as 0-18)

**Parse Implementation**:
```go
func Parse(data []byte) (*model.UniversalRecipe, error) {
    // Step 1: Validate file size
    if len(data) != 1024 {
        return nil, fmt.Errorf("invalid NP3 file: expected 1024 bytes, got %d", len(data))
    }

    // Step 2: Validate magic bytes
    if !bytes.HasPrefix(data, np3MagicBytes) {
        return nil, fmt.Errorf("invalid NP3 file: wrong magic bytes")
    }

    // Step 3: Extract parameters from known byte offsets
    recipe := &model.UniversalRecipe{
        SourceFormat: "np3",
    }

    // Sharpening at byte offset 0x100 (0-9 → 0-150 range for UniversalRecipe)
    sharpening := int(data[0x100])
    if sharpening > 9 {
        return nil, fmt.Errorf("invalid sharpening value: %d (expected 0-9)", sharpening)
    }
    recipe.Sharpness = sharpening * 15  // Map 0-9 → 0-150

    // Contrast at byte offset 0x104 (-3 to +3 → -100 to +100)
    contrast := int(data[0x104])
    if contrast > 6 {
        return nil, fmt.Errorf("invalid contrast value: %d (expected 0-6)", contrast)
    }
    recipe.Contrast = (contrast - 3) * 33  // Map 0-6 → -100 to +100

    // ... similar for Brightness, Saturation, Hue

    return recipe, nil
}
```

**Generate Implementation**:
```go
func Generate(recipe *model.UniversalRecipe) ([]byte, error) {
    if recipe == nil {
        return nil, fmt.Errorf("recipe is nil")
    }

    // Step 1: Initialize 1024-byte buffer with zeros
    data := make([]byte, 1024)

    // Step 2: Write magic bytes
    copy(data[0:4], np3MagicBytes)

    // Step 3: Map UniversalRecipe parameters to NP3 range
    // Sharpness 0-150 → 0-9
    sharpening := recipe.Sharpness / 15
    if sharpening > 9 {
        sharpening = 9
    }
    data[0x100] = byte(sharpening)

    // Contrast -100 to +100 → 0-6
    contrast := (recipe.Contrast / 33) + 3
    if contrast < 0 {
        contrast = 0
    }
    if contrast > 6 {
        contrast = 6
    }
    data[0x104] = byte(contrast)

    // ... similar for other parameters

    return data, nil
}
```

**Domain Constraints**:
- NP3 format is proprietary, reverse-engineered through binary analysis
- Fixed 1024-byte structure, cannot be extended
- Limited to 5 parameters (cannot map 50+ XMP parameters)
- Approximation required when generating from XMP/lrtemplate

#### Module: internal/formats/xmp
**Responsibility**: Parse and generate Adobe Lightroom CC (.xmp) XML/RDF files

**File Structure**:
```
internal/formats/xmp/
├── parse.go       # Parse([]byte) (*model.UniversalRecipe, error)
├── generate.go    # Generate(*model.UniversalRecipe) ([]byte, error)
└── xmp_test.go    # Table-driven tests with 913 sample files
```

**XMP Format Specification**:
- XML/RDF structure with Adobe XMP namespaces
- Parameters: 50+ including Exposure, Contrast, Highlights, Shadows, Whites, Blacks, Saturation, Vibrance, HSL (8 colors × 3 properties), Tone Curves, Split Toning
- Namespace declarations: `xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"`

**Parse Implementation**:
```go
func Parse(data []byte) (*model.UniversalRecipe, error) {
    // Step 1: Parse XML structure
    var xmpDoc XMPDocument
    if err := xml.Unmarshal(data, &xmpDoc); err != nil {
        return nil, fmt.Errorf("invalid XMP XML: %w", err)
    }

    // Step 2: Validate namespace declarations
    if !strings.Contains(string(data), "http://ns.adobe.com/camera-raw-settings") {
        return nil, fmt.Errorf("missing camera-raw-settings namespace")
    }

    // Step 3: Extract parameters from XML attributes/elements
    recipe := &model.UniversalRecipe{
        SourceFormat: "xmp",
    }

    // Parse crs:Exposure attribute
    if val, ok := xmpDoc.GetAttribute("crs:Exposure"); ok {
        exposure, _ := strconv.ParseFloat(val, 64)
        recipe.Exposure = exposure
    }

    // Parse crs:Contrast attribute
    if val, ok := xmpDoc.GetAttribute("crs:Contrast"); ok {
        contrast, _ := strconv.Atoi(val)
        recipe.Contrast = contrast
    }

    // ... extract all 50+ parameters

    return recipe, nil
}
```

**Generate Implementation**:
```go
func Generate(recipe *model.UniversalRecipe) ([]byte, error) {
    if recipe == nil {
        return nil, fmt.Errorf("recipe is nil")
    }

    // Step 1: Build XMP structure with namespaces
    xmpDoc := &XMPDocument{
        Namespaces: []string{
            "xmlns:x=\"adobe:ns:meta/\"",
            "xmlns:crs=\"http://ns.adobe.com/camera-raw-settings/1.0/\"",
            "xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\"",
        },
    }

    // Step 2: Set parameters from UniversalRecipe
    xmpDoc.SetAttribute("crs:Exposure", fmt.Sprintf("%.2f", recipe.Exposure))
    xmpDoc.SetAttribute("crs:Contrast", strconv.Itoa(recipe.Contrast))
    // ... set all 50+ parameters

    // Step 3: Marshal to XML with indentation
    data, err := xml.MarshalIndent(xmpDoc, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("marshal XMP: %w", err)
    }

    return data, nil
}
```

**Domain Constraints**:
- Must follow Adobe XMP specification
- Namespace declarations required for valid XMP
- Support both sidecar and embedded XMP
- Case-sensitive attribute names

#### Module: internal/formats/lrtemplate
**Responsibility**: Parse and generate Lightroom Classic (.lrtemplate) Lua table files

**File Structure**:
```
internal/formats/lrtemplate/
├── parse.go           # Parse([]byte) (*model.UniversalRecipe, error)
├── generate.go        # Generate(*model.UniversalRecipe) ([]byte, error)
└── lrtemplate_test.go # Table-driven tests with 566 sample files
```

**lrtemplate Format Specification**:
- Lua table syntax: `s = { key = value, ... }`
- Parameters: Same as XMP (50+), different syntax
- Quoted strings and escaped characters
- Nested tables for HSL color adjustments

**Parse Implementation**:
```go
func Parse(data []byte) (*model.UniversalRecipe, error) {
    // Step 1: Validate Lua table syntax
    if !bytes.HasPrefix(data, []byte("s = {")) {
        return nil, fmt.Errorf("invalid lrtemplate: expected 's = {'")
    }

    // Step 2: Parse Lua table using regex patterns
    recipe := &model.UniversalRecipe{
        SourceFormat: "lrtemplate",
    }

    // Extract Exposure field
    if match := exposureRegex.FindSubmatch(data); len(match) > 1 {
        exposure, _ := strconv.ParseFloat(string(match[1]), 64)
        recipe.Exposure = exposure
    }

    // Extract Contrast field
    if match := contrastRegex.FindSubmatch(data); len(match) > 1 {
        contrast, _ := strconv.Atoi(string(match[1]))
        recipe.Contrast = contrast
    }

    // ... extract all 50+ parameters using regex

    return recipe, nil
}
```

**Generate Implementation**:
```go
func Generate(recipe *model.UniversalRecipe) ([]byte, error) {
    if recipe == nil {
        return nil, fmt.Errorf("recipe is nil")
    }

    var buf bytes.Buffer
    buf.WriteString("s = {\n")

    // Write parameters in sorted order
    buf.WriteString(fmt.Sprintf("\tExposure = %.2f,\n", recipe.Exposure))
    buf.WriteString(fmt.Sprintf("\tContrast = %d,\n", recipe.Contrast))
    // ... write all 50+ parameters

    buf.WriteString("}\n")
    return buf.Bytes(), nil
}
```

**Domain Constraints**:
- Lua syntax compatibility required
- Quoted strings must handle escaped characters
- Field order doesn't matter, but sorted order preferred for consistency
- Must handle nested tables for complex structures

#### Module: internal/model
**Responsibility**: Define UniversalRecipe data structure (central hub)

**File Structure**:
```
internal/model/
├── recipe.go       # UniversalRecipe and ColorAdjustment structs
└── recipe_test.go  # Unit tests for data model
```

**Complete Data Model**:
```go
package model

type UniversalRecipe struct {
    // Metadata
    Name         string `json:"name"`
    SourceFormat string `json:"source_format"`  // "np3", "xmp", "lrtemplate"

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
    X float64 `json:"x"`  // 0.0 to 1.0
    Y float64 `json:"y"`  // 0.0 to 1.0
}
```

**Design Rationale**:
- Supports superset of all format capabilities (50+ parameters)
- JSON tags for serialization (useful for debugging, future API)
- Zero values are meaningful (no adjustment)
- Extensible via Metadata map for unknown fields
- Type safety (int for discrete values, float64 for continuous)

### Data Models and Contracts

#### Parameter Mapping Rules

**Direct Mapping** (parameters are equivalent across formats):
```
XMP Exposure      → UniversalRecipe.Exposure     → lrtemplate Exposure
XMP Contrast      → UniversalRecipe.Contrast     → lrtemplate Contrast
XMP Highlights    → UniversalRecipe.Highlights   → lrtemplate Highlights
... (40+ parameters with 1:1 mapping)
```

**Approximation Mapping** (parameters don't have 1:1 equivalents):
```
NP3 Sharpening (0-9) → UniversalRecipe.Sharpness (0-150)
  Generate: Sharpness / 15 (with bounds checking)
  Parse:    Sharpening * 15

NP3 Saturation (-3 to +3) → UniversalRecipe.Saturation (-100 to +100)
  Generate: Saturation / 33 (with bounds checking)
  Parse:    Saturation * 33

XMP Vibrance → NP3 Saturation (approximation, warn user)
  Vibrance and Saturation are similar but not identical
  Use Saturation as approximation, reduce by 20% to avoid oversaturation
```

**Unmappable Parameters** (no equivalent in target format):
```
XMP Grain        → NP3 (no equivalent, warn user, store in Metadata)
XMP Tone Curve   → NP3 (no equivalent, warn user, discard)
XMP Split Toning → NP3 (no equivalent, warn user, discard)

When generating NP3 from XMP/lrtemplate:
- Return ConversionError with unmappable fields listed
- OR: Add warnings to Metadata field in UniversalRecipe
```

**Validation Rules**:
```go
// Validate all parameters are within allowed ranges
func (r *UniversalRecipe) Validate() error {
    if r.Exposure < -5.0 || r.Exposure > 5.0 {
        return fmt.Errorf("exposure out of range: %.2f (expected -5.0 to +5.0)", r.Exposure)
    }
    if r.Contrast < -100 || r.Contrast > 100 {
        return fmt.Errorf("contrast out of range: %d (expected -100 to +100)", r.Contrast)
    }
    // ... validate all fields
    return nil
}
```

### APIs and Interfaces

#### Public API: converter.Convert()

**Signature**:
```go
func Convert(input []byte, from, to string) ([]byte, error)
```

**Parameters**:
- `input []byte`: Raw file bytes (NP3 binary, XMP XML, or lrtemplate Lua)
- `from string`: Source format identifier ("np3", "xmp", "lrtemplate")
- `to string`: Target format identifier ("np3", "xmp", "lrtemplate")

**Returns**:
- `[]byte`: Converted file bytes in target format
- `error`: ConversionError if conversion fails, nil on success

**Usage Examples**:
```go
// CLI usage
input, _ := os.ReadFile("portrait.np3")
output, err := converter.Convert(input, "np3", "xmp")
if err != nil {
    log.Fatal(err)
}
os.WriteFile("portrait.xmp", output, 0644)

// TUI usage
output, err := converter.Convert(selectedFile.Bytes, selectedFile.Format, targetFormat)

// WASM usage (called from JavaScript via go:wasmexport)
output, err := converter.Convert(wasmInputBytes, "xmp", "lrtemplate")
```

**Error Handling**:
```go
output, err := converter.Convert(input, "np3", "xmp")
if err != nil {
    var convErr *ConversionError
    if errors.As(err, &convErr) {
        // Type-safe error checking
        log.Printf("Failed to %s %s: %v", convErr.Operation, convErr.Format, convErr.Cause)
    } else {
        // Unexpected error type
        log.Fatal(err)
    }
}
```

**Thread Safety**: Convert() is stateless and thread-safe (can be called concurrently)

**Performance Contract**: Complete in <100ms (validated via benchmarks)

#### Internal APIs

**Format Parser Interface** (not exported, internal consistency):
```go
// Each format package implements:
func Parse(data []byte) (*model.UniversalRecipe, error)

// Contract:
// - Validate input immediately (fail fast)
// - Return ConversionError on failure
// - Populate SourceFormat field in UniversalRecipe
// - Return nil recipe on error
```

**Format Generator Interface** (not exported, internal consistency):
```go
// Each format package implements:
func Generate(recipe *model.UniversalRecipe) ([]byte, error)

// Contract:
// - Validate recipe is not nil
// - Return ConversionError on failure
// - Generate valid output that opens in target application
// - Handle missing fields gracefully (use zero values)
```

### Workflows and Sequencing

#### Conversion Flow

**Normal Flow** (happy path):
```
User Code
    │
    ▼
converter.Convert(input, "np3", "xmp")
    │
    ├─► Validate format strings ("np3", "xmp")
    │
    ├─► np3.Parse(input)
    │     │
    │     ├─► Validate file size (1024 bytes)
    │     ├─► Validate magic bytes
    │     ├─► Extract parameters from byte offsets
    │     └─► Return UniversalRecipe
    │
    ├─► Receive UniversalRecipe (hub)
    │
    ├─► xmp.Generate(recipe)
    │     │
    │     ├─► Build XMP structure with namespaces
    │     ├─► Set parameters from UniversalRecipe
    │     ├─► Marshal to XML
    │     └─► Return XMP bytes
    │
    └─► Return XMP bytes to user
```

**Error Flow** (parse failure):
```
User Code
    │
    ▼
converter.Convert(input, "np3", "xmp")
    │
    ├─► Validate format strings
    │
    ├─► np3.Parse(input)
    │     │
    │     ├─► Validate file size
    │     │     └─► FAIL: len(input) != 1024
    │     │
    │     └─► Return nil, fmt.Errorf("invalid NP3 file: expected 1024 bytes")
    │
    ├─► Wrap error in ConversionError
    │     └─► Operation: "parse", Format: "np3", Cause: original error
    │
    └─► Return nil, ConversionError
          │
          ▼
User Code: Check error type, log message
```

**Round-Trip Flow** (accuracy validation):
```
Original XMP File
    │
    ▼
xmp.Parse(xmpData)
    │
    ▼
UniversalRecipe (hub)
    │
    ├─────────────────┐
    ▼                 ▼
lrtemplate.Generate   xmp.Generate
    │                 │
    ▼                 ▼
lrtemplate bytes  XMP bytes (round-trip)
    │                 │
    ▼                 │
lrtemplate.Parse      │
    │                 │
    ▼                 │
UniversalRecipe       │
    │                 │
    ├─────────────────┘
    │
    ▼
Compare original vs recovered
    │
    ├─► Exposure: -0.5 vs -0.5 ✓
    ├─► Contrast: +20 vs +20 ✓
    ├─► Saturation: +10 vs +10 ✓
    └─► Accuracy: 100% for this file
```

## Non-Functional Requirements

### Performance

**Target**: <100ms per conversion (measured via Go benchmarks)

**Performance Budget**:
- NP3 parse: <10ms (fixed 1024-byte binary read)
- XMP parse: <30ms (XML parsing with encoding/xml)
- lrtemplate parse: <20ms (regex-based Lua parsing)
- UniversalRecipe overhead: <5ms (struct copy, no allocation)
- NP3 generate: <5ms (fixed 1024-byte binary write)
- XMP generate: <20ms (XML marshaling with encoding/xml)
- lrtemplate generate: <10ms (string concatenation)
- **Total budget**: 100ms worst case (XMP → lrtemplate)

**Validation Strategy**:
```bash
# Run benchmarks to validate <100ms target
go test -bench=BenchmarkConvert ./internal/converter/

# Expected output:
# BenchmarkConvert_NP3_to_XMP-8        20000    85000 ns/op  (85ms) ✓
# BenchmarkConvert_XMP_to_lrtemplate-8 15000    92000 ns/op  (92ms) ✓
# BenchmarkConvert_NP3_to_lrtemplate-8 25000    67000 ns/op  (67ms) ✓
```

**Optimization Strategy**:
- Use `encoding/xml` (stdlib) instead of reflection-based parsers
- Pre-compile regex patterns for lrtemplate parsing (avoid repeated compilation)
- Avoid unnecessary memory allocations (reuse buffers where possible)
- No goroutines (overhead not justified for <100ms target)

**WASM Consideration**: WASM builds may be 10-20% slower than native, budget allows for this overhead

### Security

**Input Validation**:
- **File size limits**: Reject files >10MB (prevents memory exhaustion)
- **Format validation**: Validate magic bytes, XML structure, Lua syntax before processing
- **Parameter bounds checking**: Validate all extracted values are within expected ranges
- **No code execution**: Lua parsing uses regex only (does not execute Lua code)

**Memory Safety**:
- Go's memory safety prevents buffer overflows
- No unsafe pointer operations (except in WASM export layer, Epic 2)
- Bounded allocations (no unbounded slices/maps)

**Error Information Disclosure**:
- ConversionError messages safe to display to users
- No internal file paths, memory addresses, or sensitive data in errors
- Stack traces only in debug mode (not production)

**Denial of Service Protection**:
- File size validation prevents excessive memory usage
- No recursion in parsers (prevents stack overflow)
- Timeout not required (<100ms guarantee)

**Privacy**:
- No network calls (all processing local/client-side)
- No telemetry, no logging of file contents
- Stateless design (no data retention)

### Reliability/Availability

**Error Handling**:
- Fail fast: Validate input immediately, return clear errors
- No silent failures: All errors wrapped in ConversionError
- Recoverable: All errors return to caller (no panics)
- Deterministic: Same input always produces same output or same error

**Fault Tolerance**:
- Invalid input handled gracefully (does not crash application)
- Partial files rejected (no partial conversion)
- Corrupt files rejected with clear error messages

**Availability**:
- No external dependencies (stdlib only)
- No network calls (cannot fail due to network issues)
- No database (cannot fail due to DB unavailability)
- Stateless (no state corruption possible)

**Testing for Reliability**:
```bash
# Test with 1,501 real sample files
go test ./internal/formats/np3/     # 22 NP3 files
go test ./internal/formats/xmp/     # 913 XMP files
go test ./internal/formats/lrtemplate/  # 566 lrtemplate files

# Test round-trip conversions
go test -run TestRoundTrip ./internal/converter/

# Expected: 0 failures, all files parse correctly
```

### Observability

**Logging** (CLI/TUI only, not Web):
```go
import "log/slog"

// Log conversion start
slog.Info("starting conversion",
    slog.String("file", inputPath),
    slog.String("from", sourceFormat),
    slog.String("to", targetFormat),
)

// Log conversion success
slog.Info("conversion complete",
    slog.String("file", inputPath),
    slog.Duration("elapsed", time.Since(start)),
)

// Log conversion error
slog.Error("conversion failed",
    slog.String("file", inputPath),
    slog.String("error", err.Error()),
)

// Debug logging (verbose mode only)
slog.Debug("parsed recipe",
    slog.String("name", recipe.Name),
    slog.Int("contrast", recipe.Contrast),
    slog.Int("saturation", recipe.Saturation),
)
```

**Metrics** (development only, no runtime metrics):
- Benchmarks: `go test -bench=. ./...` (validates <100ms target)
- Test coverage: `go test -cover ./...` (target ≥90%)
- No runtime metrics collection (privacy-first design)

**Debugging**:
- UniversalRecipe JSON serialization for debugging (json tags on struct)
- Verbose logging mode in CLI/TUI (slog.Debug level)
- Clear error messages with context (ConversionError type)

## Dependencies and Integrations

**External Dependencies**: None (Go stdlib only)

**Internal Dependencies**:
- `encoding/xml` - XMP XML parsing and generation
- `encoding/binary` - NP3 binary parsing (little-endian, big-endian)
- `bytes` - Buffer operations for lrtemplate generation
- `regexp` - Lua parsing for lrtemplate (pre-compiled patterns)
- `fmt` - Error formatting
- `strconv` - String to int/float conversion

**Downstream Consumers** (other epics):
- Epic 2 (Web Interface) - Will call converter.Convert() via WASM export
- Epic 3 (CLI Interface) - Will call converter.Convert() directly
- Epic 4 (TUI Interface) - Will call converter.Convert() directly
- Epic 5 (Data Extraction) - Will use UniversalRecipe struct for inspection

**Integration Points**:
```
┌──────────────┐
│   Epic 2     │
│ (Web/WASM)   │──┐
└──────────────┘  │
                  │
┌──────────────┐  │    ┌──────────────────────┐
│   Epic 3     │  │    │   Epic 1             │
│   (CLI)      │──┼───►│ converter.Convert()  │
└──────────────┘  │    └──────────────────────┘
                  │
┌──────────────┐  │
│   Epic 4     │  │
│   (TUI)      │──┘
└──────────────┘
```

**No Breaking Changes**: converter.Convert() signature is stable, will not change

## Acceptance Criteria (Authoritative)

### FR-1.1: NP3 Format Support
- ✅ Parse all 22 sample NP3 files from testdata/np3/ without errors
- ✅ Extract all 5 supported parameters (Sharpening, Contrast, Brightness, Saturation, Hue) accurately
- ✅ Generate NP3 files that open in Nikon NX Studio without errors or warnings
- ✅ Round-trip conversion (NP3 → XMP → NP3) produces binary-identical output (tolerance ±1 byte for rounding)
- ✅ Validate magic bytes on parse (reject invalid files immediately)
- ✅ Validate file size is exactly 1024 bytes (reject files that are too small or too large)

### FR-1.2: XMP Format Support
- ✅ Parse all 913 sample XMP files from testdata/xmp/ without errors
- ✅ Handle both sidecar XMP and embedded XMP (detect format automatically)
- ✅ Generate XMP files that load in Lightroom CC without errors or warnings
- ✅ Preserve XMP namespace declarations (xmlns:crs, xmlns:rdf, xmlns:x)
- ✅ Extract all 50+ parameters correctly (Exposure, Contrast, Highlights, Shadows, Whites, Blacks, Saturation, Vibrance, HSL colors, etc.)
- ✅ Round-trip conversion (XMP → lrtemplate → XMP) produces semantically identical output (tolerance ±1 for integer rounding)

### FR-1.3: lrtemplate Format Support
- ✅ Parse all 566 sample lrtemplate files from testdata/lrtemplate/ without errors
- ✅ Handle quoted strings and escaped characters correctly (backslashes, quotes)
- ✅ Generate lrtemplate files that load in Lightroom Classic without errors or warnings
- ✅ Maintain parameter parity with XMP (same 50+ parameters)
- ✅ Validate Lua table syntax (reject files that don't start with "s = {")
- ✅ Round-trip conversion (lrtemplate → XMP → lrtemplate) produces identical output

### FR-1.4: Universal Intermediate Representation
- ✅ Single UniversalRecipe struct captures superset of all format capabilities
- ✅ Supports 50+ parameters from XMP/lrtemplate
- ✅ Supports 5 core parameters from NP3
- ✅ Graceful handling of format-specific features (store in Metadata map if unmappable)
- ✅ Extensible for future formats (add new fields without breaking existing code)
- ✅ Intermediate representation adds <5ms overhead per conversion (validated via benchmarks)

### FR-1.5: Bidirectional Conversion
- ✅ All conversion paths work bidirectionally (NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate)
- ✅ Parameter mapping documented in code comments (direct, approximation, unmappable)
- ✅ Accuracy ≥95% for core parameters (validated via round-trip tests)
- ✅ Unmappable parameters reported to user via ConversionError or Metadata field
- ✅ Handle missing parameters gracefully (use zero values, no crashes)

### FR-1.6: Parameter Mapping & Approximation
- ✅ Mapping rules documented in code (comments in generate.go files)
- ✅ Visual similarity ≥95% for common adjustments (validated via round-trip tests)
- ✅ User receives warnings for unmappable features (XMP Grain → NP3 not supported)
- ✅ No silent data loss (all unmappable fields stored in Metadata or reported as error)
- ✅ Approximation examples implemented:
  - XMP Vibrance → NP3 Saturation (reduce by 20% to avoid oversaturation)
  - XMP Grain → NP3 (not mappable, warn user)

### Testing Criteria
- ✅ Test coverage ≥90% for all format packages (go test -cover)
- ✅ All 1,501 sample files parse without errors
- ✅ All round-trip tests pass (tolerance ±1 for rounding)
- ✅ Benchmarks validate <100ms conversion time (go test -bench)
- ✅ Table-driven tests for each format package

### Performance Criteria
- ✅ Conversion completes in <100ms (validated via benchmarks)
- ✅ Intermediate representation overhead <5ms (measured separately)
- ✅ Zero memory leaks (validated via go test -memprofile)

## Traceability Mapping

### PRD Requirements to Architecture Components

| PRD Requirement | Architecture Component | Implementation Details |
|-----------------|------------------------|------------------------|
| FR-1.1: NP3 Format Support | internal/formats/np3/ | parse.go, generate.go, np3_test.go with 22 sample files |
| FR-1.2: XMP Format Support | internal/formats/xmp/ | parse.go, generate.go, xmp_test.go with 913 sample files |
| FR-1.3: lrtemplate Format Support | internal/formats/lrtemplate/ | parse.go, generate.go, lrtemplate_test.go with 566 sample files |
| FR-1.4: Universal Intermediate Representation | internal/model/recipe.go | UniversalRecipe struct with 50+ fields, ColorAdjustment, Point types |
| FR-1.5: Bidirectional Conversion | internal/converter/converter.go | Convert() function orchestrates parse → hub → generate |
| FR-1.6: Parameter Mapping & Approximation | internal/formats/{format}/generate.go | Mapping logic in each generator, warnings via ConversionError |

### Architecture Decisions to Implementation

| Architecture Decision | Rationale | Implementation |
|-----------------------|-----------|----------------|
| Hub-and-Spoke Pattern | Simplifies adding new formats (N parsers + N generators vs N² converters) | converter.Convert() orchestrates parse → UniversalRecipe → generate |
| Single API (converter.Convert()) | Ensures consistent behavior across CLI/TUI/Web | All interfaces MUST call Convert(), no direct parser calls |
| Zero Dependencies | WASM compatibility, no external security risks | Only encoding/xml, encoding/binary, bytes, regexp from stdlib |
| ConversionError Type | Type-safe error handling with context | Wraps all errors with Operation, Format, Cause fields |
| Table-Driven Tests | Comprehensive validation with real files | 1,501 sample files in testdata/, round-trip tests |
| Stateless Design | Thread-safe, WASM-compatible, no state corruption | Convert() takes input bytes, returns output bytes, no side effects |

### Test Strategy to Acceptance Criteria

| Test Type | Purpose | Coverage |
|-----------|---------|----------|
| Unit Tests (parse.go) | Validate individual parsers | 22 NP3 + 913 XMP + 566 lrtemplate = 1,501 files |
| Unit Tests (generate.go) | Validate individual generators | Same 1,501 files, validate output opens in target app |
| Round-Trip Tests | Validate accuracy (≥95%) | All conversion paths (NP3↔XMP, XMP↔lrtemplate, NP3↔lrtemplate) |
| Benchmarks | Validate performance (<100ms) | BenchmarkConvert for each conversion path |
| Coverage Tests | Validate test coverage (≥90%) | go test -cover ./... |

## Risks, Assumptions, Open Questions

### Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| NP3 format is reverse-engineered, Nikon may change format | HIGH | Round-trip testing validates correctness, monitor for format changes in new Nikon software versions |
| XMP namespace changes in future Lightroom versions | MEDIUM | Follow Adobe XMP specification, include version detection in parser |
| lrtemplate Lua syntax edge cases not covered by 566 samples | MEDIUM | Expand test coverage as new edge cases discovered, use fuzzing to find corner cases |
| Performance <100ms may not be achievable in WASM | MEDIUM | Budget allows 10-20% WASM overhead, optimize critical paths if needed |
| Parameter approximation (XMP Vibrance → NP3 Saturation) may not achieve 95% visual similarity | LOW | User testing with photographers to validate approximations, adjust formulas if needed |

### Assumptions

1. **Sample Files Representative**: 1,501 sample files cover majority of real-world presets
   - Validation: Add more samples as users report issues
2. **Go 1.24 WASM Support Stable**: go:wasmexport is production-ready
   - Validation: Test WASM build in Epic 2, have fallback plan (older WASM export method)
3. **No Server-Side Processing**: All conversions are client-side (Web) or local (CLI/TUI)
   - Validation: Architecture decision, no change expected
4. **Parameter Ranges Consistent**: XMP/lrtemplate use same ranges (-100 to +100, etc.)
   - Validation: Verified via Adobe documentation, sample files
5. **Binary Format Stability**: NP3 format won't change frequently
   - Validation: Monitor Nikon software releases, maintain test suite

### Open Questions

1. **How to handle future format versions?**
   - Decision: Include version detection in parsers, maintain backward compatibility
   - Owner: Dev (Story 1-2, 1-4, 1-6)

2. **Should we support partial conversions (ignore unmappable fields silently)?**
   - Decision: No silent data loss, always report unmappable fields via ConversionError
   - Owner: SM (clarify in Story 1-9)

3. **How to optimize Lua parsing without external library?**
   - Decision: Use pre-compiled regex patterns, benchmark against alternatives
   - Owner: Dev (Story 1-6)

4. **Should UniversalRecipe include ALL XMP fields or just common ones?**
   - Decision: Start with 50+ common fields, add more as needed (extensible via Metadata map)
   - Owner: Dev (Story 1-1)

5. **How to validate generated files open correctly in target applications?**
   - Decision: Manual testing with Nikon NX Studio, Lightroom CC, Lightroom Classic
   - Owner: TEA (test plan in Story 1-2, 1-4, 1-6)

## Test Strategy Summary

### Unit Testing
- **Framework**: Go stdlib testing package (`testing`)
- **Pattern**: Table-driven tests with t.Run() for each sample file
- **Coverage Target**: ≥90% for all format packages
- **Sample Files**: 1,501 real files from testdata/ directory

**Example Test Structure**:
```go
func TestParse(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/np3/*.np3")
    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            data, _ := os.ReadFile(file)
            recipe, err := Parse(data)
            if err != nil {
                t.Errorf("Parse() error = %v", err)
            }
            // Validate recipe fields are within expected ranges
        })
    }
}
```

### Round-Trip Testing
- **Purpose**: Validate ≥95% accuracy for bidirectional conversions
- **Method**: Format A → UniversalRecipe → Format B → UniversalRecipe → Compare
- **Tolerance**: ±1 for integer rounding (e.g., Contrast: 20 vs 21 is acceptable)
- **Paths Tested**: NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate

**Example Round-Trip Test**:
```go
func TestRoundTrip_XMP_lrtemplate(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/xmp/*.xmp")
    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            // Parse original XMP
            origData, _ := os.ReadFile(file)
            orig, _ := xmp.Parse(origData)

            // Generate lrtemplate
            lrtData, _ := lrtemplate.Generate(orig)

            // Parse lrtemplate back
            recovered, _ := lrtemplate.Parse(lrtData)

            // Compare critical fields
            tolerance := 1
            if abs(orig.Contrast - recovered.Contrast) > tolerance {
                t.Errorf("Contrast mismatch: %d vs %d", orig.Contrast, recovered.Contrast)
            }
        })
    }
}
```

### Performance Testing
- **Framework**: Go benchmarks (`go test -bench`)
- **Target**: <100ms per conversion
- **Validation**: Run benchmarks in CI/CD pipeline

**Example Benchmark**:
```go
func BenchmarkConvert_NP3_to_XMP(b *testing.B) {
    data, _ := os.ReadFile("../../../testdata/np3/sample.np3")
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = converter.Convert(data, "np3", "xmp")
    }
}
```

### Integration Testing
- **Scope**: Validate generated files open in target applications
- **Method**: Manual testing with Nikon NX Studio, Lightroom CC, Lightroom Classic
- **Frequency**: Before each release, after format parser changes
- **Test Cases**: 5-10 representative files per format

### Regression Testing
- **Trigger**: Any change to format parsers/generators
- **Method**: Run full test suite (1,501 files + round-trip tests)
- **Expected**: 0 failures, no new errors introduced

### Test Execution
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./internal/converter/

# Run specific format tests
go test ./internal/formats/np3/
go test ./internal/formats/xmp/
go test ./internal/formats/lrtemplate/

# Run round-trip tests only
go test -run TestRoundTrip ./internal/converter/
```

### Test Data Management
- **Location**: testdata/{np3,xmp,lrtemplate}/
- **Organization**: Sample files grouped by format
- **Version Control**: All sample files committed to git (binary files)
- **Maintenance**: Add new samples as edge cases discovered

---

## Story Breakdown (Reference)

Epic 1 is decomposed into the following user stories (tracked in sprint-status.yaml):

1. **Story 1-1**: Universal Recipe Data Model (DONE)
2. **Story 1-2**: NP3 Binary Parser (IN REVIEW)
3. **Story 1-3**: NP3 Binary Generator (IN REVIEW)
4. **Story 1-4**: XMP XML Parser (BACKLOG)
5. **Story 1-5**: XMP XML Generator (BACKLOG)
6. **Story 1-6**: lrtemplate Lua Parser (BACKLOG)
7. **Story 1-7**: lrtemplate Lua Generator (BACKLOG)
8. **Story 1-8**: Parameter Mapping Rules (BACKLOG)
9. **Story 1-9**: Bidirectional Conversion API (BACKLOG)

Each story will reference this tech spec for detailed technical requirements and acceptance criteria.
