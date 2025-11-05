# Story 1.9: bidirectional-conversion-api

Status: drafted

## Story

As a developer using the Recipe conversion engine,
I want a single unified API (`converter.Convert()`) that orchestrates bidirectional conversions between all supported formats,
so that CLI, TUI, and Web interfaces can share consistent conversion logic without duplicating code.

## Acceptance Criteria

### FR-1: Unified Conversion API
- [ ] Implement `converter.Convert(input []byte, from, to string) ([]byte, error)` function
- [ ] Support all format pairs: np3↔xmp, np3↔lrtemplate, xmp↔lrtemplate (6 conversion paths)
- [ ] Validate format strings ("np3", "xmp", "lrtemplate") and return clear error for invalid formats
- [ ] Route to appropriate parser based on `from` parameter
- [ ] Route to appropriate generator based on `to` parameter
- [ ] Return ConversionError with context (Operation, Format, Cause) for all failures

### FR-2: Format Auto-Detection
- [ ] Detect NP3 format by magic bytes (first 4 bytes) and 1024-byte file size
- [ ] Detect XMP format by XML structure and camera-raw-settings namespace
- [ ] Detect lrtemplate format by "s = {" Lua table syntax
- [ ] Implement `converter.DetectFormat(input []byte) (string, error)` helper function
- [ ] Allow Convert() to work with empty `from` parameter (auto-detect)

### FR-3: Error Handling & Transparency
- [ ] Wrap all parser/generator errors in ConversionError with operation context
- [ ] Return clear error messages for invalid input (wrong format, corrupted file)
- [ ] Report unmappable parameters via ConversionError.Warnings field
- [ ] Provide actionable error messages (e.g., "File appears to be XMP but missing required namespace")
- [ ] No silent failures - all conversion issues surfaced to caller

### FR-4: Performance Targets
- [ ] Single conversion completes in <100ms (measured via Go benchmarks)
- [ ] UniversalRecipe hub overhead <5ms per conversion
- [ ] No unnecessary memory allocations (use bytes.Buffer efficiently)
- [ ] Thread-safe (Convert() can be called concurrently without race conditions)
- [ ] Stateless design (no global state, no side effects)

### FR-5: Integration with Existing Parsers/Generators
- [ ] Call np3.Parse() for NP3 input, np3.Generate() for NP3 output
- [ ] Call xmp.Parse() for XMP input, xmp.Generate() for XMP output
- [ ] Call lrtemplate.Parse() for lrtemplate input, lrtemplate.Generate() for lrtemplate output
- [ ] Verify no regressions in existing parser/generator tests
- [ ] All existing format tests continue to pass

### Non-Functional Requirements

**NFR-1: Code Quality**
- [ ] Follows Go naming conventions and idiomatic patterns
- [ ] Clear separation of concerns (validation, routing, error handling)
- [ ] Comprehensive inline documentation (godoc comments)
- [ ] No golint or go vet issues introduced
- [ ] Consistent with existing codebase style

**NFR-2: Testing**
- [ ] Unit tests for Convert() covering all 6 conversion paths
- [ ] Unit tests for DetectFormat() with valid and invalid inputs
- [ ] Unit tests for ConversionError wrapping and unwrapping
- [ ] Integration tests validating end-to-end conversions
- [ ] Performance benchmarks documenting <100ms target achievement
- [ ] Test coverage ≥90% for converter package

**NFR-3: Documentation**
- [ ] Godoc comments on all exported functions
- [ ] Usage examples in package documentation
- [ ] Error handling examples demonstrating ConversionError usage
- [ ] Performance characteristics documented
- [ ] Thread-safety guarantees documented

## Tasks / Subtasks

- [ ] Task 1: Create converter package structure (AC: FR-1)
  - [ ] 1.1: Create internal/converter/ directory
  - [ ] 1.2: Create converter.go with Convert() function signature
  - [ ] 1.3: Create converter_test.go for unit tests
  - [ ] 1.4: Create bench_test.go for performance benchmarks
  - [ ] 1.5: Create error.go with ConversionError type

- [ ] Task 2: Implement ConversionError type (AC: FR-3, NFR-1)
  - [ ] 2.1: Define ConversionError struct with Operation, Format, Cause, Warnings fields
  - [ ] 2.2: Implement Error() method for string representation
  - [ ] 2.3: Implement Unwrap() method for error chain compatibility
  - [ ] 2.4: Add godoc comments explaining usage patterns
  - [ ] 2.5: Create unit tests for error handling

- [ ] Task 3: Implement format validation (AC: FR-1, FR-2)
  - [ ] 3.1: Create validateFormat() helper to check format strings
  - [ ] 3.2: Implement DetectFormat() for auto-detection
  - [ ] 3.3: Add NP3 detection by magic bytes + file size
  - [ ] 3.4: Add XMP detection by XML structure + namespace
  - [ ] 3.5: Add lrtemplate detection by Lua table syntax
  - [ ] 3.6: Create unit tests for all detection scenarios
  - [ ] 3.7: Test with malformed/corrupted files

- [ ] Task 4: Implement Convert() orchestration logic (AC: FR-1, FR-5)
  - [ ] 4.1: Implement format validation (call validateFormat)
  - [ ] 4.2: Implement parser routing based on `from` parameter
  - [ ] 4.3: Call appropriate parser (np3.Parse, xmp.Parse, lrtemplate.Parse)
  - [ ] 4.4: Receive UniversalRecipe from parser
  - [ ] 4.5: Implement generator routing based on `to` parameter
  - [ ] 4.6: Call appropriate generator (np3.Generate, xmp.Generate, lrtemplate.Generate)
  - [ ] 4.7: Return generated bytes to caller
  - [ ] 4.8: Wrap all errors in ConversionError with context

- [ ] Task 5: Create comprehensive unit tests (AC: NFR-2)
  - [ ] 5.1: Test all 6 conversion paths (np3→xmp, np3→lrtemplate, xmp→np3, xmp→lrtemplate, lrtemplate→np3, lrtemplate→xmp)
  - [ ] 5.2: Test with sample files from testdata/ (22 NP3, 913 XMP, 566 lrtemplate)
  - [ ] 5.3: Test format validation with invalid format strings
  - [ ] 5.4: Test DetectFormat() with all format types
  - [ ] 5.5: Test error handling (corrupted files, wrong formats)
  - [ ] 5.6: Test ConversionError wrapping and unwrapping
  - [ ] 5.7: Verify test coverage ≥90%

- [ ] Task 6: Create performance benchmarks (AC: FR-4, NFR-2)
  - [ ] 6.1: Benchmark Convert() for each conversion path
  - [ ] 6.2: Benchmark DetectFormat() overhead
  - [ ] 6.3: Measure UniversalRecipe hub overhead (<5ms target)
  - [ ] 6.4: Test thread-safety with concurrent conversions
  - [ ] 6.5: Document performance characteristics
  - [ ] 6.6: Verify <100ms target achieved for all paths

- [ ] Task 7: Integration validation (AC: FR-5)
  - [ ] 7.1: Run all existing NP3 parser/generator tests
  - [ ] 7.2: Run all existing XMP parser/generator tests
  - [ ] 7.3: Run all existing lrtemplate parser/generator tests
  - [ ] 7.4: Verify no regressions (all tests pass)
  - [ ] 7.5: Test end-to-end conversions with real files
  - [ ] 7.6: Validate round-trip accuracy (A→B→A produces identical output)

- [ ] Task 8: Documentation and examples (AC: NFR-3)
  - [ ] 8.1: Write godoc comments for Convert() function
  - [ ] 8.2: Write godoc comments for DetectFormat() function
  - [ ] 8.3: Write godoc comments for ConversionError type
  - [ ] 8.4: Create package-level documentation with usage examples
  - [ ] 8.5: Document error handling patterns
  - [ ] 8.6: Document performance characteristics
  - [ ] 8.7: Document thread-safety guarantees

## Dev Notes

### Technical Approach

**Implementation Strategy**:
1. **Single Responsibility**: converter package orchestrates, doesn't implement parsing/generation logic
2. **Fail Fast**: Validate inputs immediately, return clear errors
3. **No Surprises**: All errors surfaced to caller, no silent failures
4. **Performance First**: Minimize allocations, reuse buffers, optimize hot paths
5. **Future-Proof**: Extensible for additional formats (Canon, Sony, etc.)

**Package Structure** (internal/converter/):
```
internal/converter/
├── converter.go      # Convert() and DetectFormat() functions
├── error.go          # ConversionError type
├── converter_test.go # Unit tests
└── bench_test.go     # Performance benchmarks
```

**Convert() Implementation Pattern**:
```go
func Convert(input []byte, from, to string) ([]byte, error) {
    // Step 1: Validate format strings
    if err := validateFormat(from); err != nil {
        return nil, &ConversionError{
            Operation: "validate",
            Format:    from,
            Cause:     err,
        }
    }
    if err := validateFormat(to); err != nil {
        return nil, &ConversionError{
            Operation: "validate",
            Format:    to,
            Cause:     err,
        }
    }

    // Step 2: Parse input to UniversalRecipe
    var recipe *model.UniversalRecipe
    var err error

    switch from {
    case "np3":
        recipe, err = np3.Parse(input)
    case "xmp":
        recipe, err = xmp.Parse(input)
    case "lrtemplate":
        recipe, err = lrtemplate.Parse(input)
    }

    if err != nil {
        return nil, &ConversionError{
            Operation: "parse",
            Format:    from,
            Cause:     err,
        }
    }

    // Step 3: Generate output from UniversalRecipe
    var output []byte

    switch to {
    case "np3":
        output, err = np3.Generate(recipe)
    case "xmp":
        output, err = xmp.Generate(recipe)
    case "lrtemplate":
        output, err = lrtemplate.Generate(recipe)
    }

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

**DetectFormat() Implementation Pattern**:
```go
func DetectFormat(input []byte) (string, error) {
    // Check file size and magic bytes for NP3
    if len(input) == 1024 && bytes.HasPrefix(input, np3MagicBytes) {
        return "np3", nil
    }

    // Check for XMP XML structure
    if bytes.Contains(input, []byte("<?xml")) &&
       bytes.Contains(input, []byte("camera-raw-settings")) {
        return "xmp", nil
    }

    // Check for lrtemplate Lua syntax
    if bytes.HasPrefix(input, []byte("s = {")) {
        return "lrtemplate", nil
    }

    return "", fmt.Errorf("unknown format: unable to detect from file content")
}
```

**ConversionError Type**:
```go
type ConversionError struct {
    Operation string   // "validate", "parse", "generate"
    Format    string   // "np3", "xmp", "lrtemplate"
    Cause     error    // Underlying error
    Warnings  []string // Non-fatal issues (unmappable parameters)
}

func (e *ConversionError) Error() string {
    msg := fmt.Sprintf("%s %s: %v", e.Operation, e.Format, e.Cause)
    if len(e.Warnings) > 0 {
        msg += fmt.Sprintf(" (warnings: %v)", e.Warnings)
    }
    return msg
}

func (e *ConversionError) Unwrap() error {
    return e.Cause
}
```

### Architecture Alignment

**Tech Spec Reference** (tech-spec-epic-1.md:119-164):
- Section "Module: internal/converter" specifies the exact API signature
- Convert() orchestrates parse → hub → generate flow
- ConversionError provides context for all failures

**Architecture.md Alignment**:
- Hub-and-spoke pattern: UniversalRecipe is the central hub
- Single API rule: All interfaces MUST call converter.Convert()
- Error handling pattern: Type-safe ConversionError with operation context

**Performance Budget** (tech-spec-epic-1.md:740-770):
- Target: <100ms per conversion (validated via benchmarks)
- Hub overhead: <5ms (measured separately)
- No goroutines needed (overhead not justified for <100ms target)

### Learnings from Previous Story (1-9a)

**From Story 1-9a (Metadata Field Implementation)**:

**New Capabilities Created**:
- ✅ **Metadata Field**: `Metadata map[string]interface{}` added to UniversalRecipe (recipe.go:124)
- ✅ **Purpose**: Stores unmappable format-specific parameters during conversions
- ✅ **Usage**: Preserve data that doesn't have 1:1 mapping between formats

**Implementation Patterns Established**:
- Field added at end of UniversalRecipe struct (after format-specific fields)
- JSON tag: `json:"metadata,omitempty"` (excludes empty maps from JSON)
- XML tag: `xml:"-"` (excluded due to map[string]interface{} incompatibility)
- Zero-value behavior: nil or empty Metadata has zero performance impact

**Testing Approach to Reuse**:
- Comprehensive unit tests with 99.7% coverage
- Story 1-8 documentation examples validated via dedicated tests
- Performance benchmarks documented zero impact
- All existing tests passed (no regressions)

**Key Insights for This Story**:
1. **Metadata Usage**: When converting between formats with unmappable parameters:
   - Store unmappable data in recipe.Metadata
   - Use key naming convention: "format_fieldname" (e.g., "xmp_tone_curve_pv2012")
   - Report unmappable fields via ConversionError.Warnings
   - Example: XMP Grain → NP3 (not supported, store in Metadata or warn user)

2. **Error Handling Pattern**: Follow story 1-9a approach:
   - Validate inputs immediately (fail fast)
   - Comprehensive error context (ConversionError with Operation, Format, Cause)
   - No silent failures (all issues surfaced)

3. **Test Coverage Strategy**:
   - 90%+ coverage target is achievable
   - Table-driven tests with real sample files
   - Validate all documented examples work correctly
   - Performance benchmarks document characteristics

4. **Files Created in 1-9a** (Available for Reference):
   - internal/models/recipe_metadata_test.go (10 tests, all patterns)
   - internal/models/recipe_bench_test.go (5 benchmarks)
   - Can reference these for test structure and benchmark patterns

5. **Architectural Decisions from 1-9a**:
   - Minimal, surgical changes preferred
   - Zero breaking changes to existing code
   - Documentation-driven (implement what's specified)
   - Test-first approach (write tests validating examples)

**Files Modified by 1-9a**:
- internal/models/recipe.go - Added Metadata field (line 124)

**Integration Points**:
- All existing parsers (np3, xmp, lrtemplate) work with Metadata field
- All existing generators can read/write Metadata
- Zero regressions in existing tests

**Warnings for This Story**:
- Ensure Metadata is used when unmappable parameters detected
- Report unmappable fields via ConversionError.Warnings
- Follow key naming convention: "format_fieldname"

[Source: stories/1-9a-metadata-field-implementation.md]

### Testing Strategy

**Unit Tests** (converter_test.go):
1. `TestConvert_AllPaths` - Test all 6 conversion paths with sample files
2. `TestConvert_InvalidFormat` - Test validation with invalid format strings
3. `TestConvert_CorruptedInput` - Test error handling with malformed files
4. `TestConvert_Concurrent` - Test thread-safety with parallel conversions
5. `TestDetectFormat_NP3` - Validate NP3 detection by magic bytes + size
6. `TestDetectFormat_XMP` - Validate XMP detection by XML + namespace
7. `TestDetectFormat_LRTemplate` - Validate lrtemplate detection by Lua syntax
8. `TestDetectFormat_Invalid` - Test with unknown/corrupted files
9. `TestConversionError_Wrapping` - Test error wrapping and unwrapping
10. `TestConversionError_Warnings` - Test unmappable parameter warnings

**Integration Tests**:
1. Run all existing format tests to ensure no regressions
2. Test end-to-end conversions with representative sample files
3. Validate round-trip accuracy (A→B→A produces identical output)
4. Test with all 1,501 sample files (22 NP3, 913 XMP, 566 lrtemplate)

**Performance Benchmarks** (bench_test.go):
1. `BenchmarkConvert_NP3_to_XMP` - Measure NP3→XMP conversion time
2. `BenchmarkConvert_NP3_to_LRTemplate` - Measure NP3→lrtemplate conversion time
3. `BenchmarkConvert_XMP_to_NP3` - Measure XMP→NP3 conversion time
4. `BenchmarkConvert_XMP_to_LRTemplate` - Measure XMP→lrtemplate conversion time
5. `BenchmarkConvert_LRTemplate_to_NP3` - Measure lrtemplate→NP3 conversion time
6. `BenchmarkConvert_LRTemplate_to_XMP` - Measure lrtemplate→XMP conversion time
7. `BenchmarkDetectFormat` - Measure format detection overhead
8. `BenchmarkHubOverhead` - Measure UniversalRecipe hub overhead

**Expected Performance**:
- All conversions <100ms (validated via benchmarks)
- Hub overhead <5ms (measured separately)
- Format detection <1ms (negligible overhead)

### Success Metrics

- ✅ Convert() function completes all 6 conversion paths successfully
- ✅ DetectFormat() correctly identifies all supported formats
- ✅ All conversion errors wrapped in ConversionError with context
- ✅ All existing parser/generator tests pass (no regressions)
- ✅ Performance benchmarks validate <100ms target
- ✅ Test coverage ≥90% for converter package
- ✅ No golint or go vet issues introduced
- ✅ Thread-safe (concurrent conversions work correctly)

### Risks and Mitigations

**Risk 1: Performance Target Not Met**
- Mitigation: Optimize critical paths, minimize allocations, profile with pprof
- Validation: Run benchmarks early, identify bottlenecks before implementation complete

**Risk 2: Breaking Changes to Existing Parsers/Generators**
- Mitigation: Run full test suite after any changes, fix regressions immediately
- Validation: All existing format tests must pass

**Risk 3: Error Handling Complexity**
- Mitigation: Use ConversionError consistently, follow established patterns from 1-9a
- Validation: Comprehensive unit tests for all error scenarios

**Risk 4: Unmappable Parameters Not Reported**
- Mitigation: Use Metadata field + ConversionError.Warnings for transparency
- Validation: Test with known unmappable conversions (XMP Grain → NP3)

### References

**Tech Spec References**:
- tech-spec-epic-1.md:119-164 (Module: internal/converter)
- tech-spec-epic-1.md:569-649 (APIs and Interfaces)
- tech-spec-epic-1.md:647-735 (Workflows and Sequencing)
- tech-spec-epic-1.md:740-770 (Performance Non-Functional Requirements)

**Architecture References**:
- architecture.md (Hub-and-spoke pattern, Single API rule)
- architecture.md (Error handling pattern: ConversionError)

**Story References**:
- stories/1-9a-metadata-field-implementation.md (Metadata usage patterns)
- stories/1-8-parameter-mapping-rules.md (Parameter mapping documentation)

**Code References**:
- internal/models/recipe.go (UniversalRecipe struct with Metadata field)
- internal/formats/np3/parse.go, generate.go (NP3 parser/generator)
- internal/formats/xmp/parse.go, generate.go (XMP parser/generator)
- internal/formats/lrtemplate/parse.go, generate.go (lrtemplate parser/generator)

### Estimated Effort

- **Package Setup**: ~30 minutes (create directories, files)
- **Convert() Implementation**: ~2-3 hours (routing logic, error handling)
- **DetectFormat() Implementation**: ~1 hour (magic byte detection, validation)
- **ConversionError Implementation**: ~1 hour (error type, methods, tests)
- **Unit Tests**: ~3-4 hours (comprehensive test coverage)
- **Integration Validation**: ~1 hour (run existing tests, fix regressions)
- **Performance Benchmarks**: ~2 hours (create benchmarks, validate targets)
- **Documentation**: ~1 hour (godoc comments, usage examples)

**Total Estimate**: 11-14 hours

### Completion Criteria

This story is complete when:
1. ✅ Convert() function implemented and working for all 6 conversion paths
2. ✅ DetectFormat() function implemented and tested
3. ✅ ConversionError type implemented with wrapping/unwrapping
4. ✅ All unit tests pass (≥90% coverage)
5. ✅ All existing parser/generator tests pass (no regressions)
6. ✅ Performance benchmarks validate <100ms target
7. ✅ Godoc documentation complete with usage examples
8. ✅ CLI, TUI, and Web interfaces can use converter.Convert() directly

Once complete, this API will be the foundation for all Recipe interfaces (CLI, TUI, Web), ensuring consistent conversion behavior across all user experiences.

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

## Change Log

- **2025-11-04**: Story 1-9 drafted
  - Defined API signature for converter.Convert()
  - Specified 6 conversion paths (all format pairs)
  - Created comprehensive acceptance criteria
  - Documented integration with previous stories (1-9a Metadata field)
  - Established performance targets (<100ms per conversion)
  - Planned test strategy (unit tests, integration tests, benchmarks)
