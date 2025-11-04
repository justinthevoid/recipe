# Development Backlog

This file tracks action items and follow-ups identified during code reviews and other development activities.

---

## Story 1-2-np3-binary-parser - Code Review Follow-ups
**Date**: 2025-11-04
**Review Outcome**: BLOCKED
**Reviewer**: Senior Developer (Code Review Workflow)

### High Priority Action Items

- [ ] **[H1] Fix chunk parser offset/logic** - Chunk parser finds 0 chunks in all test files. Must investigate binary structure at offset 0x2C and correct parsing logic or offset.
  - File: `internal/formats/np3/parse.go:270-334`
  - Impact: Parameter extraction completely broken without working chunk parser
  - Verification: Run `test_chunks.go` - should find >0 chunks per file

- [ ] **[H2] Implement actual parameter extraction from chunks** - Current `estimateParametersFromChunks()` uses hardcoded fallbacks, never actual chunk data.
  - File: `internal/formats/np3/parse.go:168-267`
  - Impact: All files return identical parameter values
  - Verification: Run `test_params.go` - should show diverse parameters across files
  - Depends on: H1 (needs working chunk parser first)

- [ ] **[H3] Add Brightness/Hue mapping** - Lines 362, 365 have TODO comments for brightness/hue but no implementation
  - File: `internal/formats/np3/parse.go:362,365`
  - Impact: Missing 2 of 5 core parameters (violates AC1)
  - Evidence: Commented code `// Brightness: params.Brightness, // TODO: map from [-100,100] to heuristic range`

- [ ] **[H4] Increase test coverage to ≥95%** - Current coverage is 88.5%, target is 95%
  - File: `internal/formats/np3/np3_test.go`
  - Missing coverage: Error paths, edge cases, chunk parsing validation
  - Verification: `go test -coverprofile=coverage.out && go tool cover -func=coverage.out`

### Medium Priority Action Items

- [ ] **[M1] Fix Parse() error handling** - AC3 requires graceful degradation but current implementation returns errors that halt conversion
  - File: `internal/formats/np3/parse.go:52-130`
  - Need: Fallback values when chunks missing, continue with defaults
  - Test case: Verify Parse() succeeds even with malformed chunk data

- [ ] **[M2] Add validation warnings** - AC4 requires logging invalid values detected. No validation warnings currently logged.
  - File: `internal/formats/np3/parse.go:132-166`
  - Need: Detect out-of-range chunk values, log warnings via logger
  - Test case: Verify warning logged for value >100 or <-100

- [ ] **[M3] Document binary format findings** - Code comments claim chunk structure understanding but parser doesn't work
  - File: `internal/formats/np3/parse.go` (add package-level doc comment)
  - Need: Document actual binary structure discovered through hex analysis
  - Include: Offset locations, chunk format, parameter encoding details

### Low Priority Action Items

- [ ] **[L1] Add benchmark tests** - Architecture patterns recommend benchmarking for binary parsers
  - File: Create `internal/formats/np3/parse_bench_test.go`
  - Benchmarks: Parse(), parseChunks(), extractParameters()
  - Target: <10ms for typical NP3 file parse

---

## Story 1-5-xmp-xml-generator - Code Review Follow-ups
**Date**: 2025-11-04
**Review Outcome**: CHANGES_REQUIRED
**Reviewer**: claude-sonnet-4.5 (Senior Developer Code Review Agent)
**Model**: claude-sonnet-4-5-20250929

### High Priority Action Items

- [ ] **[H1] Implement ToneCurve generation** - AC FR-5.1 and Task 4.2 marked complete but NOT implemented
  - File: `internal/formats/xmp/generate.go:287-288`
  - Impact: BLOCKING - Core AC requirement not satisfied, story cannot move to DONE
  - Evidence: Comment states "// Tone Curve - omitted for now (complex array format)"
  - Implementation Required:
    - Parse ToneCurve array from UniversalRecipe.ToneCurve (array of Point{X, Y})
    - Research XMP ToneCurve format specification
    - Format as XMP-compliant string
    - Add to buildXMPDocument() at generate.go:287
    - Add test cases for ToneCurve generation
    - Verify round-trip with ToneCurve data
  - **Alternative**: Formally descope ToneCurve from FR-5 and update AC (requires SM approval)

### Medium Priority Action Items

None.

### Low Priority Action Items

- [ ] **[L1] Update AC FR-3.5 Temperature range documentation** - AC states incorrect range
  - File: `docs/stories/1-5-xmp-xml-generator.md` line 37
  - Current AC: "(-100 to +100)"
  - Should be: "(2000 to 50000 Kelvin)"
  - Impact: Documentation accuracy only, implementation is correct
  - Note: XMP Temperature uses Kelvin units (2000-50000K), not slider values

- [ ] **[L2] Clarify AC FR-3.6 Tint range specification** - Implementation uses wider range
  - File: `docs/stories/1-5-xmp-xml-generator.md` line 38
  - AC states: "(-100 to +100)"
  - Implementation: "(-150 to +150)" at generate.go:150-152
  - Action: Verify if [-150, +150] is correct per XMP specification
  - Resolution: Either update AC to match implementation OR constrain code to [-100, +100]

---

**Notes**:
- Story 1-5: HIGH priority item H1 is BLOCKING - must be completed before story can be re-reviewed
- Story 1-5: LOW items are documentation clarifications only
- Story 1-2: All HIGH priority items must be completed before story can be re-reviewed. MEDIUM items should be addressed in same iteration. LOW items can be deferred to epic retrospective if time-constrained.
