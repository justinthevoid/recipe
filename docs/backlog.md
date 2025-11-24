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

---

## Story 10-1-landing-page-redesign - Code Review Follow-ups
**Date**: 2025-11-10
**Review Outcome**: APPROVED (Production-Ready)
**Reviewer**: Senior Developer Code Review Workflow
**Model**: claude-sonnet-4-5-20250929

### Low Priority Action Items (Non-Blocking Improvements)

- [ ] **[L1] Refactor CSS into modular files** - Tech spec recommends modular CSS structure
  - Current: All styles in single `web/static/style.css` (2463 lines)
  - Recommended Structure:
    - `web/css/main.css` - Global styles, CSS variables, reset
    - `web/css/components.css` - Reusable badges, buttons, cards
    - `web/css/layout.css` - Responsive grid, breakpoints
    - `web/css/legacy.css` - Story 2 styles to be deprecated
  - Impact: Improves maintainability for future Epic 10 stories
  - Effort: ~2 hours
  - Epic: 10
  - References: `docs/tech-spec-epic-10.md:37-46`

- [ ] **[L2] Consolidate duplicate badge styles** - Two badge implementations exist
  - Current:
    - `.badge` (Story 10-1, style.css:186-223) - Modern, WCAG compliant
    - `.format-badge` (Story 2-3, style.css:724-761) - Legacy
  - Action: Merge into single badge system, update all references
  - Impact: Reduces maintenance overhead, ensures consistency
  - Effort: ~1 hour
  - Epic: 10
  - References: Tech spec AC-2 (consistent styling across interface)

- [ ] **[L3] Run WebPageTest performance validation** - AC-5 performance target validation
  - Test: Load time <2 seconds on 3G connection (1.6 Mbps, 300ms RTT)
  - URL: https://recipe.justins.studio (or local build)
  - Tools: [WebPageTest](https://www.webpagetest.org/)
  - Metrics to verify:
    - First Contentful Paint (FCP): <1.5s
    - Largest Contentful Paint (LCP): <2.0s
    - Total Blocking Time (TBT): <300ms
  - Impact: Informational - implementation follows all best practices, likely meets target
  - Effort: ~30 minutes
  - Epic: 10
  - References: `docs/tech-spec-epic-10.md:422-444` (AC-5 Performance requirements)

**Notes**:
- Story 10-1: All follow-ups are LOW priority, non-blocking improvements
- Story 10-1: Implementation is production-ready and approved for deployment
- Story 10-1: L1 and L2 improve code maintainability for future Epic 10 work
- Story 10-1: L3 is informational validation only (implementation already follows best practices)

---

**Notes**:
- Story 1-5: HIGH priority item H1 is BLOCKING - must be completed before story can be re-reviewed
- Story 1-5: LOW items are documentation clarifications only
- Story 1-2: All HIGH priority items must be completed before story can be re-reviewed. MEDIUM items should be addressed in same iteration. LOW items can be deferred to epic retrospective if time-constrained.

---

## Story 11-1-css-filter-mapping - Code Review Follow-ups
**Date**: 2025-11-15
**Review Outcome**: BLOCKED
**Reviewer**: Senior Developer Code Review Workflow
**Model**: claude-sonnet-4-5-20250929
**Epic**: 11 (Image Preview System)

### High Priority Action Items

- [ ] **[H1] Integrate preview with upload flow** - Task 7 marked complete but NO integration code exists
  - File: `web/static/main.js` or `web/static/upload.js`
  - Impact: BLOCKING - AC-3 requires preview applied when file uploaded/dropped
  - Evidence: No calls to `applyPreviewFilter()` from upload handling code
  - Implementation Required:
    - Import `applyPreviewFilter()` from `preview.js`
    - Call it in `handlePresetUpload()` or equivalent function after preset parsing
    - Pass parsed UniversalRecipe to the function
    - Verify preview image filter updates correctly
  - Related AC: AC-3 (Preview applied when file uploaded/dropped)
  - Related Task: Task 7 (Integration with Preset Upload Flow)
  - Verification: Upload NP3/XMP/lrtemplate file → preview image filter updates

- [ ] **[H2] Add "Preview" button to file cards in upload interface** - UI element to trigger preview missing
  - File: `web/index.html`, `web/static/upload.js`
  - Impact: BLOCKING - AC-3 requires preview visible when file uploaded
  - Evidence: No button in file card HTML, no click handler to show preview modal
  - Implementation Required:
    - Add "Preview" button to each uploaded file card in HTML
    - Add click handler that calls `showPreviewModal()` and `applyPreviewFilter(recipe)`
    - Style button consistently with existing UI (badge/button styles)
    - Test with multiple uploaded files (each card gets preview button)
  - Related AC: AC-3
  - Related Task: Task 7
  - Verification: Upload file → "Preview" button appears on file card → clicking shows modal with filter applied

- [ ] **[H3] Execute performance testing on iPhone 8 and Android mid-range device** - AC-6 performance validation required
  - File: Test report to be added to story Dev Notes or separate test log
  - Impact: BLOCKING - AC-6 requires <100ms filter application, 60fps slider responsiveness
  - Evidence: Task 5 correctly marked pending [ ], no test results documented
  - Testing Required:
    - iPhone 8 (iOS 12+): Measure filter application time, slider responsiveness
    - Android mid-range (e.g., Samsung Galaxy A52): Same measurements
    - GPU acceleration verification (use browser DevTools Performance tab)
    - Document results in story file or separate test report
  - Device Requirements:
    - iPhone 8 running iOS 12 or later
    - Android mid-range device (2019-2021 era, e.g., Samsung Galaxy A52, Google Pixel 4a)
  - Related AC: AC-6 (Performance Optimization)
  - Related Task: Task 5 (Performance Testing)
  - Verification: Filter application <100ms, slider updates at 60fps, GPU acceleration active

### Medium Priority Action Items

- [ ] **[M1] Document manual browser testing results** - AC-5 validation requires evidence
  - File: Story Dev Notes section or separate test report
  - Impact: Task 4 marked [x] complete but no manual testing evidence provided
  - Evidence: No screenshots, test log, or documented results in story file
  - Testing Required:
    - Chrome (latest 2 versions): Desktop + Android
    - Firefox (latest 2 versions): Desktop
    - Safari (latest 2 versions): Desktop + iOS
    - Edge (latest 2 versions): Desktop
    - Verify CSS filter support detection works correctly
    - Document fallback behavior for unsupported browsers (if any)
  - Related AC: AC-5 (Browser Compatibility Testing)
  - Related Task: Task 4 (Browser Compatibility Detection)
  - Verification: Add section to story with browser versions tested, screenshots or test log

**Notes**:
- Story 11-1: All HIGH priority items (H1, H2, H3) are BLOCKING - must be completed before story can move to DONE
- Story 11-1: Implementation quality is excellent (92% test coverage, 51/51 tests pass, type validation, GPU acceleration)
- Story 11-1: Core CSS filter mapping function is complete and thoroughly tested
- Story 11-1: Integration gap (H1, H2) is critical - preview feature exists but not wired into upload flow
- Story 11-1: Performance testing (H3) must be done on actual devices, not just desktop browsers
- Story 11-1: MEDIUM item M1 should be completed with H3 in same test session
