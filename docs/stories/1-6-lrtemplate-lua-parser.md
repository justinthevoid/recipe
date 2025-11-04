# Story 1.6: lrtemplate Lua Parser

Status: review

## Story

As a developer implementing the Recipe conversion engine,
I want a robust lrtemplate Lua parser that extracts all parameters from Lightroom Classic preset files into UniversalRecipe data,
so that users can convert Lightroom Classic presets to other formats (NP3, XMP) with 95%+ accuracy.

## Acceptance Criteria

### Functional Requirements

**FR-1: Lua Table Parsing**
- ✅ Parse valid lrtemplate files starting with `s = {`
- ✅ Handle quoted strings and escaped characters correctly (backslashes, quotes, newlines)
- ✅ Extract all 50+ parameters using regex patterns (no external Lua libraries)
- ✅ Return clear error if file doesn't match lrtemplate syntax
- ✅ Support both single-line and multi-line Lua table formats

**FR-2: Core Parameter Extraction (Basic Adjustments)**
- ✅ Extract `Exposure` from `Exposure2012` or `Exposure` field (-5.0 to +5.0, float)
- ✅ Extract `Contrast` from `Contrast2012` field (-100 to +100, integer)
- ✅ Extract `Highlights` from `Highlights2012` field (-100 to +100, integer)
- ✅ Extract `Shadows` from `Shadows2012` field (-100 to +100, integer)
- ✅ Extract `Whites` from `Whites2012` field (-100 to +100, integer)
- ✅ Extract `Blacks` from `Blacks2012` field (-100 to +100, integer)

**FR-3: Color Parameter Extraction**
- ✅ Extract `Saturation` field (-100 to +100, integer)
- ✅ Extract `Vibrance` field (-100 to +100, integer)
- ✅ Extract `Clarity` from `Clarity2012` field (-100 to +100, integer)
- ✅ Extract `Sharpness` field (0 to 150, integer)
- ✅ Extract `Temperature` field (2000 to 50000 Kelvin, integer)
- ✅ Extract `Tint` field (-150 to +150, integer)

**FR-4: HSL Color Adjustments Extraction**
- ✅ Extract HSL adjustments for all 8 colors (Red, Orange, Yellow, Green, Aqua, Blue, Purple, Magenta)
- ✅ For each color, extract:
  - `Hue{Color}` field (-100 to +100, integer)
  - `Saturation{Color}` field (-100 to +100, integer)
  - `Luminance{Color}` field (-100 to +100, integer)
- ✅ Example fields: `HueAdjustmentRed`, `SaturationAdjustmentRed`, `LuminanceAdjustmentRed`
- ✅ Handle both abbreviated (Red) and full (Red) color names in field names

**FR-5: Advanced Features Extraction**
- ✅ Extract `ToneCurvePV2012` array if present (array of {x, y} coordinate pairs)
- ✅ Parse ToneCurve array syntax: `ToneCurvePV2012 = { { 0, 0 }, { 255, 255 } }`
- ✅ Extract Split Toning parameters:
  - `SplitToningShadowHue` (0 to 360, integer)
  - `SplitToningShadowSaturation` (0 to 100, integer)
  - `SplitToningHighlightHue` (0 to 360, integer)
  - `SplitToningHighlightSaturation` (0 to 100, integer)

**FR-6: Data Type Handling**
- ✅ Parse floating-point values correctly (Exposure uses Lua number format)
- ✅ Parse integer values correctly (Contrast, Saturation, etc.)
- ✅ Handle missing fields gracefully (use zero values for UniversalRecipe)
- ✅ Validate all numeric values are within expected ranges during extraction
- ✅ Handle negative numbers in Lua syntax (e.g., `Exposure = -1.5`)

**FR-7: Error Handling**
- ✅ Return clear error if file doesn't start with `s = {`
- ✅ Return clear error for malformed Lua syntax
- ✅ Wrap all errors in ConversionError with context (following pattern from Story 1-4/1-5)
- ✅ Include field name in error messages for debugging
- ✅ Handle edge cases: empty file, truncated file, invalid characters

### Non-Functional Requirements

**NFR-1: Performance**
- ✅ Parse single lrtemplate file in <20ms (target from tech spec)
- ✅ Pre-compile regex patterns for efficiency (avoid repeated compilation)
- ✅ Validation via Go benchmarks

**NFR-2: Test Coverage**
- ✅ Parse all 566 lrtemplate samples from testdata/lrtemplate/ without errors
- ✅ Test coverage ≥90% for lrtemplate package (matching parser coverage)
- ✅ Table-driven tests following Pattern 7 (Architecture doc)
- ✅ Round-trip validation: lrtemplate → parse → generate → compare produces identical output (±1 tolerance)

**NFR-3: Code Quality**
- ✅ Follow Pattern 4: File structure (parse.go in lrtemplate package)
- ✅ Follow Pattern 5: Error handling with ConversionError wrapper (reuse type from xmp/parse.go)
- ✅ Follow Pattern 6: Inline validation, fail fast
- ✅ Use stdlib regex package only (zero external dependencies, no Lua interpreter)
- ✅ Code passes gofmt and go vet without warnings

## Tasks / Subtasks

- [ ] Task 1: Design Lua parsing strategy with regex patterns (AC: FR-1)
  - [ ] 1.1: Research lrtemplate file format structure from sample files
  - [ ] 1.2: Design regex patterns for key-value extraction
  - [ ] 1.3: Plan handling of escaped characters and quoted strings
  - [ ] 1.4: Determine strategy for nested tables (ToneCurve array)
  - [ ] 1.5: Document Lua syntax compatibility requirements

- [ ] Task 2: Implement core Parse() function (AC: FR-1, FR-2, FR-3)
  - [ ] 2.1: Create function signature matching Tech Spec: `Parse([]byte) (*model.UniversalRecipe, error)`
  - [ ] 2.2: Validate file starts with `s = {` prefix
  - [ ] 2.3: Pre-compile regex patterns for all parameter types
  - [ ] 2.4: Extract core parameters (Exposure, Contrast, Highlights, Shadows, Whites, Blacks)
  - [ ] 2.5: Extract color parameters (Saturation, Vibrance, Clarity, Sharpness, Temperature, Tint)
  - [ ] 2.6: Implement inline validation for extracted values

- [ ] Task 3: Implement HSL color extraction (AC: FR-4)
  - [ ] 3.1: Create regex patterns for HSL field names (HueAdjustmentRed, etc.)
  - [ ] 3.2: Extract Red, Orange, Yellow, Green color adjustments
  - [ ] 3.3: Extract Aqua, Blue, Purple, Magenta color adjustments
  - [ ] 3.4: Validate HSL ranges for all colors (-100 to +100)
  - [ ] 3.5: Map extracted values to UniversalRecipe ColorAdjustment structs

- [ ] Task 4: Implement advanced features extraction (AC: FR-5)
  - [ ] 4.1: Create regex pattern for ToneCurvePV2012 array syntax
  - [ ] 4.2: Parse nested coordinate pairs `{ {0,0}, {255,255} }`
  - [ ] 4.3: Extract Split Toning parameters (Shadow/Highlight Hue and Saturation)
  - [ ] 4.4: Handle optional/missing fields gracefully
  - [ ] 4.5: Validate array formats and coordinate ranges

- [ ] Task 5: Implement error handling with ConversionError (AC: FR-7, NFR-3)
  - [ ] 5.1: Reuse ConversionError type from xmp/parse.go (same package location)
  - [ ] 5.2: Wrap all errors with Operation="parse", Format="lrtemplate"
  - [ ] 5.3: Include field names in error context for debugging
  - [ ] 5.4: Add validation for file format and syntax errors
  - [ ] 5.5: Handle edge cases (empty file, truncated, invalid UTF-8)

- [ ] Task 6: Implement string and character handling (AC: FR-1, FR-6)
  - [ ] 6.1: Handle escaped characters in quoted strings (\\, \", \n, \r, \t)
  - [ ] 6.2: Support both single-line and multi-line Lua tables
  - [ ] 6.3: Handle negative numbers correctly in Lua syntax
  - [ ] 6.4: Validate quoted string formats
  - [ ] 6.5: Test with edge cases (empty strings, special characters)

- [ ] Task 7: Write comprehensive tests (AC: NFR-2)
  - [ ] 7.1: Create table-driven tests using testdata/lrtemplate/ samples (566 files)
  - [ ] 7.2: Implement round-trip tests (lrtemplate → parse → generate → compare)
  - [ ] 7.3: Add edge case tests (invalid syntax, missing fields, malformed arrays)
  - [ ] 7.4: Create performance benchmarks targeting <20ms
  - [ ] 7.5: Validate test coverage ≥90%

- [ ] Task 8: Documentation and code quality (AC: NFR-3)
  - [ ] 8.1: Add GoDoc comments to Parse() function
  - [ ] 8.2: Document regex patterns in code comments with examples
  - [ ] 8.3: Run gofmt and go vet
  - [ ] 8.4: Verify code follows all architectural patterns

## Dev Notes

### Technical Approach

**Phase 1: Analysis (Learn from Previous Stories)**
1. **Study sample lrtemplate files** - Understand actual file structure from 566 samples
2. **Design regex patterns** - Extract key-value pairs without full Lua parser
3. **Identify edge cases** - Nested tables, escaped strings, special characters
4. **Map to UniversalRecipe** - Same parameter set as XMP (50+ fields)
5. **Plan validation strategy** - Inline checks, fail fast on invalid syntax

**Phase 2: Core Parser Implementation**
1. **Define Parse() function signature**:
   ```go
   func Parse(data []byte) (*model.UniversalRecipe, error)
   ```

2. **Validation steps** (fail fast):
   - Check file starts with `s = {`
   - Validate basic Lua table syntax
   - Return ConversionError on invalid input

3. **Regex-based extraction**:
   - Pre-compile all regex patterns (performance optimization)
   - Extract parameters using named groups
   - Convert string values to appropriate types (int, float)
   - Build UniversalRecipe struct

4. **Pattern examples**:
   ```go
   // Simple key-value: Exposure = -1.50
   exposureRegex := regexp.MustCompile(`Exposure2012\s*=\s*(-?\d+\.?\d*)`)

   // Integer parameter: Contrast = 20
   contrastRegex := regexp.MustCompile(`Contrast2012\s*=\s*(-?\d+)`)

   // Tone curve array: ToneCurvePV2012 = { { 0, 0 }, { 255, 255 } }
   toneCurveRegex := regexp.MustCompile(`ToneCurvePV2012\s*=\s*\{(.+?)\}`)
   ```

**Phase 3: Testing & Validation**
1. **Round-trip tests** - Core validation strategy:
   - Parse existing lrtemplate file → UniversalRecipe
   - Generate lrtemplate from recipe → new lrtemplate bytes
   - Parse generated lrtemplate → recovered UniversalRecipe
   - Compare original vs. recovered (tolerance ±1 for rounding)

2. **Edge case testing**:
   - Invalid syntax → error
   - Missing fields → zero values
   - Malformed arrays → error
   - Escaped characters → correct parsing

3. **Performance benchmarking** to validate <20ms target

### Key Technical Decisions

**Decision 1: Use regex instead of full Lua parser**
- Rationale: lrtemplate files are simple Lua tables, not complex Lua code
- Benefit: Zero external dependencies, faster than full parser, sufficient for format
- Risk: May miss edge cases in complex Lua syntax (mitigated by 566-file test suite)
- Alternative: Use gopher-lua library - adds dependency, slower, overkill for format

**Decision 2: Pre-compile all regex patterns**
- Rationale: Regex compilation is expensive, patterns are static
- Implementation: Define patterns as package-level variables, compile once
- Benefit: Significant performance gain (regex compilation can be 50%+ of parse time)
- Pattern:
  ```go
  var (
      exposureRegex  = regexp.MustCompile(`Exposure2012\s*=\s*(-?\d+\.?\d*)`)
      contrastRegex  = regexp.MustCompile(`Contrast2012\s*=\s*(-?\d+)`)
      // ... all other patterns
  )
  ```

**Decision 3: Follow XMP parameter naming (2012 process version)**
- Rationale: Lightroom Classic uses same parameter names as XMP
- Fields: `Exposure2012`, `Contrast2012`, `ToneCurvePV2012` (PV = Process Version)
- Benefit: Consistency with XMP parser, easier round-trip validation

**Decision 4: Handle escaped characters in regex**
- Rationale: Quoted strings may contain `\"`, `\\`, `\n` etc.
- Strategy: Extract quoted strings first, then unescape
- Pattern: `fieldName\s*=\s*"([^"]*)"`
- Post-process: Replace escape sequences

### Learnings from Previous Story (1-5-xmp-xml-generator)

**From Story 1-5-xmp-xml-generator (Status: done, Code Review: APPROVED)**

- **ConversionError Pattern**: Defined at parse.go:40-61 - MUST reuse this type
  ```go
  type ConversionError struct {
      Operation string  // "parse" for parser
      Format    string  // "lrtemplate"
      Field     string  // Optional: specific field that caused error
      Cause     error   // Underlying error
  }
  ```
  - All errors MUST be wrapped with Operation="parse", Format="lrtemplate"
  - Include Field parameter for debugging context

- **File Structure Pattern**: Follow same layout as NP3 and XMP packages
  ```
  internal/formats/lrtemplate/
  ├── parse.go           # ⬅ THIS STORY
  ├── generate.go        # Story 1-7 (next)
  └── lrtemplate_test.go # Extend with parser tests
  ```

- **Parameter Mapping Consistency**: Use same parameter names as XMP
  - Temperature is Kelvin units (2000-50000K), not slider values
  - Tint range is [-150, +150], wider than some documentation suggests
  - HSL fields use format: `HueAdjustmentRed`, `SaturationAdjustmentOrange`, etc.

- **Test Coverage Achievement**: XMP parser reached 90.6%, generator 92.3%
  - lrtemplate parser should match or exceed 90% coverage
  - Use table-driven tests with all 566 sample files
  - Round-trip testing is PRIMARY validation strategy

- **Performance Target**: XMP parser achieves 0.045ms (much faster than 30ms target)
  - lrtemplate parser <20ms target is very achievable
  - Pre-compiled regex patterns are key to performance
  - Benchmark against XMP parser performance

- **Nullable Types Gotcha**: Temperature is `*int` (nullable) in some formats
  - Check XMP parser implementation for Temperature handling
  - UniversalRecipe may use regular `int` or `*int` - verify model definition

- **New Services/Patterns Created in Previous Stories**:
  - `internal/formats/xmp/parse.go` (650+ lines) - XMP parser reference
  - `internal/formats/xmp/generate.go` (410+ lines) - XMP generator reference
  - ConversionError type now available for reuse
  - Round-trip test pattern established in xmp_test.go

- **Architectural Deviations**:
  - None reported from XMP story - stayed aligned with architecture patterns
  - Follow same pattern: Pattern 4 (file structure), Pattern 5 (errors), Pattern 6 (validation), Pattern 7 (testing)

- **Technical Debt from XMP Stories**:
  - None reported - all ACs satisfied, code review approved
  - lrtemplate parser should maintain this quality standard

- **Files to Reference**:
  - `internal/formats/xmp/parse.go` - Example of parser implementation
  - `internal/formats/xmp/xmp_test.go` - Example of table-driven tests
  - `testdata/xmp/sample.xmp` - Example test file structure
  - Use as templates for lrtemplate parser implementation

### Project Structure Notes

**File Locations** (from Architecture doc Pattern 4):
```
internal/formats/lrtemplate/
├── parse.go              # ⬅ THIS STORY
├── generate.go           # Story 1-7 (next)
└── lrtemplate_test.go    # Parser tests
```

**Integration Points**:
- Parser will be called by `converter.Convert()` (Story 1-9)
- Works alongside NP3 and XMP parsers in hub-and-spoke pattern
- Outputs to all interfaces: CLI (Epic 3), TUI (Epic 4), Web (Epic 2)

**Testing Alignment**:
- Use testdata/lrtemplate/ directory (566 sample files)
- Round-trip tests validate parser + generator work together
- Table-driven tests follow Go best practices

**Reuse Opportunities**:
- ConversionError type from xmp package
- UniversalRecipe struct from model package
- Testing patterns from xmp_test.go
- DO NOT recreate these - reuse existing code

### References

**Technical References**:
- **Tech Spec**: `docs/tech-spec-epic-1.md` - Section "Module: internal/formats/lrtemplate"
- **Architecture**: `docs/architecture.md` - Pattern 4, 5, 6, 7
- **PRD**: `docs/PRD.md` - FR-1.3: lrtemplate Format Support
- **Previous Story**: `docs/stories/1-5-xmp-xml-generator.md` - Generator implementation (DONE, APPROVED)
- **XMP Parser**: `internal/formats/xmp/parse.go` - Reference implementation for parser structure
- **Model**: `internal/model/recipe.go` - UniversalRecipe struct definition

**Sample Files**:
- `testdata/lrtemplate/*.lrtemplate` - 566 real Lightroom Classic preset files
- Study several files to understand actual format variations

**External Resources**:
- Lightroom Classic Lua table format (documented in community forums)
- Go regexp package: https://pkg.go.dev/regexp
- Lua string escape sequences: https://www.lua.org/pil/2.4.html

### Success Metrics

- **566/566 sample files parse successfully** (100% parse rate)
- **≥90% test coverage** for lrtemplate package (including parser code)
- **<20ms parse time** validated by benchmarks
- **Zero parse errors** for valid lrtemplate inputs
- **Round-trip accuracy ≥95%** (lrtemplate → parse → generate → parse matches original)

## Dev Agent Record

### Context Reference

- docs/stories/1-6-lrtemplate-lua-parser.context.xml

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
