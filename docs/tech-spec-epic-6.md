# Epic Technical Specification: Validation & Testing

Date: 2025-11-06
Author: Justin
Epic ID: epic-6
Status: Draft

---

## Overview

Epic 6 establishes a comprehensive validation and testing framework for Recipe, ensuring the accuracy, performance, and cross-platform reliability of the conversion engine. This epic implements automated testing against the complete sample file corpus (1,501 files), performance benchmarking infrastructure, visual regression testing, and cross-browser compatibility verification for the Web interface.

The testing framework validates the core promise of Recipe: 95%+ conversion accuracy with <100ms performance in-browser. It provides continuous validation throughout development and establishes regression detection to maintain quality as the codebase evolves. All testing infrastructure integrates with CI/CD pipelines (GitHub Actions) to ensure every change is validated before deployment.

## Objectives and Scope

### Primary Objectives
1. **Automated Test Suite**: Validate conversion accuracy against all 1,501 sample files with 100% parse success rate
2. **Visual Regression**: Compare converted presets visually with color delta E <5 for critical colors
3. **Performance Benchmarking**: Measure and document conversion time, memory usage, and WASM binary size with automated regression detection
4. **Browser Compatibility**: Verify Web interface functionality across Chrome, Firefox, Safari (latest 2 versions) covering 90%+ browser market

### In Scope
- Go test suite with table-driven tests for all 1,501 sample files (22 NP3, 913 XMP, 566 lrtemplate, 22 DNG)
- Round-trip conversion validation (A → B → A produces functionally identical output)
- Performance benchmarks using `go test -bench` with targets: <100ms WASM, <20ms CLI, <2s batch (100 files)
- Visual regression framework using reference images with Lightroom/NX Studio comparison
- Browser compatibility test matrix (Chrome, Firefox, Safari) testing FileReader API, WASM execution, file download
- CI/CD integration with GitHub Actions for automated test runs on every PR and push
- Test coverage reporting with 90%+ goal for core conversion logic
- Regression detection to prevent performance degradation

### Out of Scope
- Manual QA processes (automated testing preferred)
- Load testing (single-user tool, not multi-user service)
- Security penetration testing (static site with no server-side processing)
- Accessibility compliance testing (basic keyboard nav sufficient, not WCAG 2.1 AA required)
- Mobile device testing (desktop/tablet focus, mobile secondary)
- End-to-end testing frameworks (Selenium, Playwright) - simple functional tests sufficient

### Success Criteria
- 100% parse success on all 1,501 valid sample files
- ≥95% round-trip conversion accuracy (tolerance ±1 for parameter rounding)
- Test suite completes in <10 seconds
- Performance benchmarks confirm <100ms WASM, <20ms CLI targets
- Visual regression shows color delta E <5 for critical colors
- Web interface functional in Chrome, Firefox, Safari (latest 2 versions)
- CI/CD pipeline blocks merges on test failures
- Test coverage ≥90% for internal/converter and internal/formats packages

## System Architecture Alignment

Epic 6 implements the **validation and quality assurance layer** that ensures all architectural patterns and decisions are correctly implemented and maintained:

### Architecture Alignment Points

1. **Testing Strategy** (aligns with Architecture Pattern 7: Testing Strategy)
   - Table-driven tests using 1,501 real sample files from `testdata/`
   - Round-trip validation (A → B → A produces identical output)
   - Comprehensive coverage of all format parsers (np3, xmp, lrtemplate)
   - Validates 95%+ accuracy goal across all conversions

2. **Performance Validation** (aligns with Architecture Section: Performance Considerations)
   - Go benchmarks validate <100ms WASM conversion target
   - CLI benchmarks validate <20ms single file, <2s batch (100 files) targets
   - Memory profiling ensures <50MB WASM heap, <100MB CLI batch processing
   - WASM binary size verification (<2MB stripped, <800KB gzipped)

3. **Error Handling Validation** (aligns with Pattern 5: Error Handling)
   - Tests verify all conversion failures wrapped in `ConversionError` type
   - Validates format-specific error context (Operation, Format, Cause)
   - Ensures graceful failure with user-friendly messages

4. **Browser Compatibility** (aligns with NFR-5: Browser Compatibility)
   - Validates FileReader API and Blob download across Chrome, Firefox, Safari
   - Confirms WebAssembly MVP support detection
   - Tests drag-and-drop events in all supported browsers
   - Verifies 90%+ browser market coverage

5. **CI/CD Integration** (aligns with Architecture Section: Deployment Architecture)
   - GitHub Actions workflow runs full test suite on every PR
   - Automated deployment blocked on test failures
   - Coverage reports generated and tracked over time
   - Performance regression detection prevents degradation

6. **Quality Gates**
   - 100% parse success on all valid sample files (gates format parser quality)
   - 90%+ test coverage for `internal/converter` and `internal/formats` (gates code quality)
   - <10 second test suite runtime (gates developer experience)
   - Zero test failures required for deployment (gates production stability)

## Detailed Design

### Services and Modules

Epic 6 adds testing infrastructure across multiple layers:

| Module                   | Responsibility                             | Inputs                                     | Outputs                                  | Owner                                |
| ------------------------ | ------------------------------------------ | ------------------------------------------ | ---------------------------------------- | ------------------------------------ |
| **Test Suite (Go)**      | Table-driven tests for all format parsers  | Sample files from testdata/                | Pass/fail results, coverage report       | internal/formats/{format}_test.go    |
| **Round-Trip Tests**     | Validate bidirectional conversion fidelity | Original files, converted files            | Accuracy metrics, diff reports           | internal/converter/roundtrip_test.go |
| **Benchmark Suite**      | Measure conversion performance             | Sample files, iteration count              | Timing (ns/op), memory (B/op, allocs/op) | internal/converter/benchmark_test.go |
| **Visual Regression**    | Compare preset outputs visually            | Reference images, presets                  | Color delta E metrics, visual diffs      | scripts/visual-regression/           |
| **Browser Compat Tests** | Validate Web interface cross-browser       | Test scenarios (upload, convert, download) | Browser support matrix                   | web/tests/compatibility/             |
| **CI/CD Pipeline**       | Automated test execution on every change   | Git commits, PRs                           | Test results, deployment gates           | .github/workflows/test.yml           |
| **Coverage Reporter**    | Track test coverage over time              | Test execution results                     | Coverage percentage, trend reports       | .github/workflows/coverage.yml       |

**Module Dependencies:**
- Test Suite depends on: testdata/ sample files, internal/converter, internal/formats
- Round-Trip Tests depend on: Test Suite, all format parsers
- Benchmark Suite depends on: internal/converter
- Visual Regression depends on: External tools (Lightroom, NX Studio), reference images
- Browser Compat depends on: web/ interface, WASM binary
- CI/CD Pipeline depends on: All test modules

### Data Models and Contracts

**Test Result Schema:**
```go
// TestResult captures outcome of a single test case
type TestResult struct {
    TestName     string        `json:"test_name"`
    FilePath     string        `json:"file_path"`
    Format       string        `json:"format"`
    Status       TestStatus    `json:"status"`  // Pass, Fail, Skip
    Duration     time.Duration `json:"duration"`
    ErrorMessage string        `json:"error_message,omitempty"`
}

type TestStatus string

const (
    TestPass TestStatus = "pass"
    TestFail TestStatus = "fail"
    TestSkip TestStatus = "skip"
)

// BenchmarkResult captures performance metrics
type BenchmarkResult struct {
    Name          string  `json:"name"`
    Iterations    int     `json:"iterations"`
    NsPerOp       int64   `json:"ns_per_op"`
    BytesPerOp    int64   `json:"bytes_per_op"`
    AllocsPerOp   int64   `json:"allocs_per_op"`
    MeetsTarget   bool    `json:"meets_target"`
    TargetNs      int64   `json:"target_ns"`
}

// VisualRegressionResult for color accuracy validation
type VisualRegressionResult struct {
    PresetName    string  `json:"preset_name"`
    SourceFormat  string  `json:"source_format"`
    TargetFormat  string  `json:"target_format"`
    DeltaE        float64 `json:"delta_e"`      // Color difference metric
    PassThreshold float64 `json:"pass_threshold"` // 5.0 for critical colors
    Passed        bool    `json:"passed"`
}

// BrowserCompatResult for cross-browser testing
type BrowserCompatResult struct {
    Browser      string `json:"browser"`
    Version      string `json:"version"`
    FileUpload   bool   `json:"file_upload"`
    WasmLoad     bool   `json:"wasm_load"`
    Conversion   bool   `json:"conversion"`
    FileDownload bool   `json:"file_download"`
    UIRender     bool   `json:"ui_render"`
    Overall      bool   `json:"overall"`  // All tests passed
}
```

**Test Coverage Schema:**
```go
type CoverageReport struct {
    Timestamp    time.Time           `json:"timestamp"`
    TotalLines   int                 `json:"total_lines"`
    CoveredLines int                 `json:"covered_lines"`
    Percentage   float64             `json:"percentage"`
    ByPackage    map[string]Coverage `json:"by_package"`
    MeetsGoal    bool                `json:"meets_goal"`  // ≥90%
}

type Coverage struct {
    Package      string  `json:"package"`
    Lines        int     `json:"lines"`
    Covered      int     `json:"covered"`
    Percentage   float64 `json:"percentage"`
}
```

### APIs and Interfaces

**Go Testing API (Standard Library):**

```go
// Table-driven test pattern for format parsers
func TestParseNP3(t *testing.T) {
    files, err := filepath.Glob("../../../testdata/np3/*.np3")
    if err != nil {
        t.Fatal(err)
    }

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

            // Validate parameter ranges
            validateRecipe(t, recipe)
        })
    }
}

// Round-trip test pattern
func TestRoundTrip_NP3_XMP_NP3(t *testing.T) {
    files, _ := filepath.Glob("../../../testdata/np3/*.np3")

    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            // Step 1: Parse original NP3
            origData, _ := os.ReadFile(file)
            orig, err := np3.Parse(origData)
            if err != nil {
                t.Fatalf("NP3 parse failed: %v", err)
            }

            // Step 2: Convert to XMP
            xmpData, err := xmp.Generate(orig)
            if err != nil {
                t.Fatalf("XMP generate failed: %v", err)
            }

            // Step 3: Parse XMP back
            xmpRecipe, err := xmp.Parse(xmpData)
            if err != nil {
                t.Fatalf("XMP parse failed: %v", err)
            }

            // Step 4: Convert back to NP3
            np3Data, err := np3.Generate(xmpRecipe)
            if err != nil {
                t.Fatalf("NP3 generate failed: %v", err)
            }

            // Step 5: Compare parameters (tolerance ±1)
            compareRecipes(t, orig, xmpRecipe, 1)
        })
    }
}

// Benchmark API
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

**CI/CD Integration API (GitHub Actions):**

```yaml
# .github/workflows/test.yml
name: Test Suite

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run Tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Check Coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 90" | bc -l) )); then
            echo "Coverage $coverage% is below 90% threshold"
            exit 1
          fi

      - name: Run Benchmarks
        run: go test -bench=. -benchmem ./internal/converter/ > benchmarks.txt

      - name: Validate Performance Targets
        run: |
          # Check WASM conversion < 100ms (100,000,000 ns)
          # Parse benchmark output and verify
          ./scripts/validate-benchmarks.sh benchmarks.txt
```

**Browser Compatibility Test Interface:**

```javascript
// web/tests/compatibility/test-suite.js
const testSuite = {
    async testFileUpload(browser) {
        // Test drag-and-drop and file picker
        const file = new File(["test"], "test.np3", { type: "application/octet-stream" });
        const result = await uploadFile(file);
        return result.success;
    },

    async testWasmLoad(browser) {
        // Test WASM initialization
        const loaded = await loadWasm();
        return loaded;
    },

    async testConversion(browser) {
        // Test actual conversion
        const result = await convertFile(testData, "np3", "xmp");
        return result !== null;
    },

    async testFileDownload(browser) {
        // Test download trigger
        const downloaded = await downloadResult(testData, "test.xmp");
        return downloaded;
    },

    async runAllTests(browser) {
        const results = {
            browser: browser.name,
            version: browser.version,
            fileUpload: await this.testFileUpload(browser),
            wasmLoad: await this.testWasmLoad(browser),
            conversion: await this.testConversion(browser),
            fileDownload: await this.testFileDownload(browser),
            uiRender: true  // Visual check
        };
        results.overall = Object.values(results).slice(2).every(v => v === true);
        return results;
    }
};
```

### Workflows and Sequencing

**Test Execution Workflow:**

```
Developer Push/PR
    ↓
GitHub Actions Trigger
    ↓
┌─────────────────────────────────┐
│ 1. Setup Environment            │
│    - Checkout code              │
│    - Install Go 1.24            │
│    - Cache dependencies         │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 2. Run Automated Test Suite     │
│    - Parse all 1,501 files      │
│    - Round-trip conversions     │
│    - Parameter validation       │
│    Duration: ~5-8 seconds       │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 3. Coverage Analysis            │
│    - Generate coverage.out      │
│    - Check ≥90% threshold       │
│    - Block merge if below       │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 4. Performance Benchmarks       │
│    - Run Go benchmarks          │
│    - Validate <100ms WASM       │
│    - Detect regressions         │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 5. Build WASM                   │
│    - Compile to WebAssembly     │
│    - Verify size <2MB stripped  │
│    - Test basic functionality   │
└─────────────────────────────────┘
    ↓
    ├─── All Pass → ✅ Merge Allowed
    └─── Any Fail → ❌ Merge Blocked
```

**Visual Regression Workflow (Manual):**

```
1. Prepare Reference Images
   - Select representative photos (portrait, landscape, etc.)
   - Apply preset in Lightroom → Export reference image
   
2. Convert Preset
   - Use Recipe to convert (e.g., XMP → NP3)
   
3. Apply Converted Preset
   - Load preset in Nikon NX Studio
   - Apply to same reference photo
   - Export test image
   
4. Compare Images
   - Use image diff tool (ImageMagick, SSIM)
   - Calculate color delta E
   - Verify delta E < 5 for critical colors
   
5. Document Results
   - Record delta E metrics
   - Note visual differences
   - Update regression test database
```

**Browser Compatibility Test Workflow:**

```
For each browser (Chrome, Firefox, Safari):
    ↓
┌─────────────────────────────────┐
│ 1. Manual Setup                 │
│    - Open browser               │
│    - Navigate to local server   │
│    - Open DevTools              │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 2. Test File Upload             │
│    - Drag file to drop zone     │
│    - Verify format detection    │
│    - Try file picker fallback   │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 3. Test WASM Execution          │
│    - Monitor WASM load time     │
│    - Verify conversion succeeds │
│    - Check console for errors   │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 4. Test File Download           │
│    - Verify download triggers   │
│    - Check file integrity       │
│    - Validate filename correct  │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│ 5. Test UI Rendering            │
│    - Verify responsive layout   │
│    - Check parameter preview    │
│    - Test error states          │
└─────────────────────────────────┘
    ↓
    Record results in compatibility matrix
```

**Regression Detection Sequence:**

```
On every CI run:
    ↓
1. Run current benchmarks → Save results
    ↓
2. Load baseline benchmarks from main branch
    ↓
3. Compare performance metrics:
   - If current > baseline * 1.1 → Flag regression (10% slower)
   - If current < baseline * 0.9 → Flag improvement
    ↓
4. Update baseline if on main branch
    ↓
5. Report regression in PR comments
```

## Non-Functional Requirements

### Performance

**Test Suite Execution Speed:**
- **Target**: Complete test suite in <10 seconds
- **Measurement**: `time go test ./...`
- **Breakdown**:
  - Format parser tests: <3 seconds (1,501 files)
  - Round-trip tests: <4 seconds (subset of critical conversions)
  - Benchmark execution: <3 seconds
- **Optimization**: Parallel test execution with `go test -parallel=8`

**Benchmark Performance Targets:**
- **WASM Conversion**: <100ms (100,000,000 ns/op)
  - Measured via `BenchmarkConvert_WASM_NP3_to_XMP`
  - Target: 50,000-100,000 ns/op (0.05-0.1ms)
- **CLI Conversion**: <20ms (20,000,000 ns/op)
  - Measured via `BenchmarkConvert_CLI_NP3_to_XMP`
  - Target: 5,000-20,000 ns/op (0.005-0.02ms)
- **Batch Processing**: <2s for 100 files
  - Measured via `BenchmarkBatchConvert_100_Files`
  - Target: 500,000,000-2,000,000,000 ns/op (0.5-2s)

**Memory Constraints:**
- **WASM Heap**: <50MB during conversion
- **CLI Memory**: <100MB for 1,000 file batch
- **Test Memory**: <500MB for full test suite
- **Measurement**: `go test -benchmem` reports B/op and allocs/op

**CI/CD Pipeline Performance:**
- **Total Pipeline**: <5 minutes from push to completion
- **Test Job**: <2 minutes
- **Benchmark Job**: <1 minute
- **Coverage Job**: <1 minute
- **Build WASM Job**: <1 minute

### Security

**Test Data Privacy:**
- All 1,501 sample files are publicly shared presets or generated test data
- No proprietary or commercial presets in test suite
- Test data committed to repository (transparent, reproducible)

**CI/CD Security:**
- GitHub Actions runs in isolated containers
- No secrets required for test execution
- Read-only access to repository during PR tests
- CLOUDFLARE_API_TOKEN only used in deployment (not testing)

**Dependency Security:**
- Zero external test dependencies (stdlib only)
- No npm packages for browser tests (manual verification)
- Go module checksums verified (go.sum)
- Dependabot alerts enabled for Go dependencies

**WASM Testing Security:**
- WASM binary tested in sandbox environment
- No network access during test execution
- File access limited to test fixtures
- Memory safety validated via fuzzing (future enhancement)

### Reliability/Availability

**Test Stability:**
- **Flake Rate**: <1% (tests must be deterministic)
- **False Positives**: Zero tolerance (failing tests indicate real issues)
- **Test Isolation**: Each test runs independently, no shared state
- **Retry Logic**: CI retries flaky tests once before failing

**Test Coverage Stability:**
- Coverage percentage must not decrease on any PR
- Baseline coverage tracked in main branch
- New code requires corresponding tests (no uncovered code merged)

**CI/CD Reliability:**
- **Uptime Target**: 99%+ (GitHub Actions SLA)
- **Fallback**: Manual test execution if CI unavailable
- **Timeout Protection**: Tests timeout after 15 minutes (prevent hung jobs)
- **Resource Limits**: 2 CPU cores, 7GB RAM per job (GitHub Actions limits)

**Regression Prevention:**
- All bug fixes require regression test
- Performance regressions blocked (>10% slowdown fails build)
- Accuracy regressions blocked (round-trip tolerance violations fail build)

**Test Data Integrity:**
- Sample files in `testdata/` are immutable (version controlled)
- File checksums validated before test execution
- Corrupted test files detected and reported

### Observability

**Test Execution Visibility:**
- **Real-Time Logs**: CI job logs show test progress in real-time
- **Failure Details**: Failed tests display file path, expected vs actual, stack trace
- **Subtest Reporting**: Individual file tests reported separately (1,501 subtests visible)
- **Duration Tracking**: Each test reports execution time

**Coverage Reporting:**
- **Package-Level Coverage**: Coverage reported per package (converter, formats/np3, formats/xmp, etc.)
- **Line-Level Coverage**: `go tool cover -html=coverage.out` generates HTML report
- **Trend Tracking**: Coverage percentage tracked over time in CI artifacts
- **Visualization**: GitHub Actions summary shows coverage badge

**Performance Monitoring:**
- **Benchmark Results**: Saved as CI artifacts for historical comparison
- **Trend Analysis**: Performance tracked over time (detect gradual degradation)
- **Comparison Reports**: PR benchmarks compared to baseline (main branch)
- **Metrics Exported**: ns/op, B/op, allocs/op saved to JSON for analysis

**Test Metrics Dashboard (Future Enhancement):**
- Test pass rate over time
- Flake rate tracking
- Average test duration
- Coverage trends
- Performance regression alerts

**Browser Compatibility Matrix:**
- Manual test results documented in `docs/browser-compatibility.md`
- Updated after each major release
- Includes browser version, OS, test results

**Alerting:**
- GitHub Actions sends failure notifications to PR author
- Critical failures (coverage drop, performance regression) mentioned in PR comments
- Deployment blocked status visible in PR status checks

## Dependencies and Integrations

**Go Standard Library (v1.24+):**
- `testing` - Test framework, benchmarking, coverage
- `testing/iotest` - I/O testing utilities
- `os`, `path/filepath` - File operations for test data
- `bytes`, `encoding/json` - Result serialization
- No external test dependencies required

**CI/CD Platform:**
- **GitHub Actions** (free tier)
  - `actions/checkout@v4` - Repository checkout
  - `actions/setup-go@v5` - Go environment setup
  - `actions/upload-artifact@v4` - Test result artifacts
  - Workflow file: `.github/workflows/test.yml`

**Test Data:**
- **testdata/ Directory** (version controlled)
  - 22 NP3 files (Nikon Picture Control samples)
  - 913 XMP files (Lightroom CC preset samples)
  - 566 lrtemplate files (Lightroom Classic preset samples)
  - Total size: ~13 MB (acceptable for repository)

**External Tools (Optional, Manual Testing Only):**
- **Adobe Lightroom CC/Classic** - Visual regression testing (apply presets to reference images)
- **Nikon NX Studio** - Validate generated NP3 files open correctly
- **ImageMagick** - Image comparison for visual regression
- **Browser DevTools** - Browser compatibility testing (Chrome, Firefox, Safari)

**Performance Baseline:**
- Benchmark results from main branch stored in CI artifacts
- Format: JSON file with ns/op, B/op, allocs/op metrics
- Updated on every merge to main
- Used for regression detection in PRs

**Integration Points:**
- **internal/converter** - All conversion logic being tested
- **internal/formats/** - All format parsers/generators being tested
- **web/recipe.wasm** - WASM binary performance validation
- **cmd/cli/** - CLI performance benchmarks

**No Runtime Dependencies:**
- Tests use same zero-dependency philosophy as main code
- Standard library only (no external test frameworks)
- Browser tests use vanilla JavaScript (no Selenium, Playwright)

## Acceptance Criteria (Authoritative)

### AC-1: Automated Test Suite Coverage
**Given** the complete test suite with 1,501 sample files  
**When** tests are executed via `go test ./...`  
**Then**:
- ✅ 100% parse success on all 1,501 valid sample files (22 NP3 + 913 XMP + 566 lrtemplate)
- ✅ Zero test failures on valid input files
- ✅ Test suite completes in <10 seconds
- ✅ All format parsers have corresponding test files (np3_test.go, xmp_test.go, lrtemplate_test.go)
- ✅ Tests discoverable via subtests (1,501 individual subtests visible in output)

### AC-2: Round-Trip Conversion Validation
**Given** sample files for each format  
**When** round-trip conversions are executed (A → B → A)  
**Then**:
- ✅ NP3 → XMP → NP3 produces functionally identical output (tolerance ±1 for rounding)
- ✅ XMP → lrtemplate → XMP preserves all critical parameters (Exposure, Contrast, Saturation, HSL)
- ✅ lrtemplate → NP3 → lrtemplate maintains parameter accuracy ≥95%
- ✅ Round-trip tests cover all format pairs (6 conversion paths)
- ✅ Failed round-trips report specific parameter mismatches with expected vs actual values

### AC-3: Test Coverage Metrics
**Given** coverage analysis via `go test -cover ./...`  
**When** coverage report is generated  
**Then**:
- ✅ Overall test coverage ≥90%
- ✅ `internal/converter` package coverage ≥90%
- ✅ `internal/formats/np3` package coverage ≥90%
- ✅ `internal/formats/xmp` package coverage ≥90%
- ✅ `internal/formats/lrtemplate` package coverage ≥90%
- ✅ Coverage report exported to `coverage.out` and `coverage.html`
- ✅ Uncovered lines documented with justification (if any)

### AC-4: Performance Benchmarks
**Given** benchmark suite via `go test -bench=. -benchmem`  
**When** benchmarks are executed  
**Then**:
- ✅ Single file WASM conversion: <100ms (100,000,000 ns/op)
- ✅ Single file CLI conversion: <20ms (20,000,000 ns/op)
- ✅ Batch 100 files CLI: <2s (2,000,000,000 ns/op)
- ✅ Memory allocation: <4096 B/op per conversion
- ✅ Memory allocations: <12 allocs/op per conversion
- ✅ Benchmark results saved to `benchmarks.txt` artifact
- ✅ Performance regression detection: PR fails if >10% slower than baseline

### AC-5: Visual Regression Testing (Manual)
**Given** reference images with applied presets  
**When** converted presets are applied to same images  
**Then**:
- ✅ Visual similarity ≥95% (subjective assessment)
- ✅ Color delta E <5 for critical colors (reds, blues, greens, skin tones)
- ✅ At least 5 representative presets tested (portrait, landscape, vintage, b&w, HDR)
- ✅ Results documented in `docs/visual-regression-results.md`
- ✅ Known differences documented with explanation (e.g., grain not supported in NP3)

### AC-6: Browser Compatibility Testing
**Given** Web interface deployed to local server  
**When** tested in Chrome, Firefox, Safari (latest 2 versions)  
**Then**:
- ✅ File upload works (drag-and-drop + file picker) in all browsers
- ✅ WASM loads and initializes successfully in all browsers
- ✅ Conversion executes without errors in all browsers
- ✅ File download triggers correctly in all browsers
- ✅ UI renders correctly (responsive layout, parameter preview) in all browsers
- ✅ Unsupported browser detection works (shows clear message for IE11)
- ✅ Results documented in `docs/browser-compatibility.md`

### AC-7: CI/CD Integration
**Given** GitHub Actions workflow configured  
**When** PR is created or code is pushed to main  
**Then**:
- ✅ Test suite runs automatically on every PR
- ✅ Coverage report generated and checked against 90% threshold
- ✅ Benchmarks run and compared to baseline
- ✅ WASM binary compiled and size validated (<2MB stripped)
- ✅ PR merge blocked if any test fails
- ✅ PR merge blocked if coverage drops below 90%
- ✅ PR merge blocked if performance regresses >10%
- ✅ Test results visible in PR status checks
- ✅ CI job completes in <5 minutes

### AC-8: Test Documentation
**Given** testing infrastructure implementation  
**When** developers need to run or add tests  
**Then**:
- ✅ README includes "Running Tests" section with examples
- ✅ Test file organization documented (pattern: {format}_test.go)
- ✅ How to add new test files to testdata/ documented
- ✅ How to run specific test suites documented (go test ./internal/formats/np3/)
- ✅ How to generate coverage reports documented
- ✅ How to run benchmarks documented
- ✅ CI/CD workflow explained in CONTRIBUTING.md

## Traceability Mapping

| Acceptance Criteria                        | Spec Section(s)                                                                | Component(s)/API(s)                                                                                                      | Test Idea                                                                                                       |
| ------------------------------------------ | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------- |
| **AC-1: Automated Test Suite Coverage**    | Services and Modules: Test Suite (Go), APIs: Table-driven test pattern         | `internal/formats/np3/np3_test.go`, `internal/formats/xmp/xmp_test.go`, `internal/formats/lrtemplate/lrtemplate_test.go` | Run `go test ./...` and verify 1,501 subtests pass; validate execution time <10s via `time` command             |
| **AC-2: Round-Trip Conversion Validation** | Services and Modules: Round-Trip Tests, APIs: Round-trip test pattern          | `internal/converter/roundtrip_test.go`, all format parsers/generators                                                    | Implement `TestRoundTrip_NP3_XMP_NP3`, `TestRoundTrip_XMP_LRT_XMP`, verify parameter equality with tolerance ±1 |
| **AC-3: Test Coverage Metrics**            | Services and Modules: Coverage Reporter, Workflows: Test Execution             | `.github/workflows/coverage.yml`, `go tool cover`                                                                        | Run `go test -coverprofile=coverage.out ./...`; parse coverage output, verify ≥90% for each package             |
| **AC-4: Performance Benchmarks**           | Services and Modules: Benchmark Suite, NFR: Performance                        | `internal/converter/benchmark_test.go`, benchmark validation script                                                      | Implement `BenchmarkConvert_*` functions, run `go test -bench=.`, validate ns/op meets targets                  |
| **AC-5: Visual Regression Testing**        | Services and Modules: Visual Regression, Workflows: Visual Regression Workflow | `scripts/visual-regression/`, external tools (Lightroom, NX Studio)                                                      | Manual process: export reference images, apply converted presets, compare with ImageMagick/SSIM                 |
| **AC-6: Browser Compatibility Testing**    | Services and Modules: Browser Compat Tests, Workflows: Browser Compat Workflow | `web/tests/compatibility/test-suite.js`, manual browser testing                                                          | Open Web interface in each browser, execute test suite, document results in compatibility matrix                |
| **AC-7: CI/CD Integration**                | Services and Modules: CI/CD Pipeline, Workflows: Test Execution                | `.github/workflows/test.yml`, GitHub Actions jobs                                                                        | Create PR, verify tests run automatically, check status checks block merge on failure                           |
| **AC-8: Test Documentation**               | N/A (documentation)                                                            | `README.md`, `CONTRIBUTING.md`                                                                                           | Review documentation for completeness, verify all test commands work as documented                              |

**PRD Requirement Mapping:**

| PRD Section                           | Tech Spec Component                                    | Stories                           |
| ------------------------------------- | ------------------------------------------------------ | --------------------------------- |
| FR-6.1: Automated Test Suite          | AC-1, AC-2, AC-3                                       | 6-1-automated-test-suite          |
| FR-6.2: Visual Regression Testing     | AC-5                                                   | 6-2-visual-regression-testing     |
| FR-6.3: Performance Benchmarking      | AC-4                                                   | 6-3-performance-benchmarking      |
| FR-6.4: Browser Compatibility Testing | AC-6                                                   | 6-4-browser-compatibility-testing |
| NFR-3: Reliability                    | AC-1, AC-2 (100% parse success, round-trip validation) | All stories                       |
| NFR-1: Performance                    | AC-4 (<100ms WASM, <20ms CLI, <2s batch)               | 6-3-performance-benchmarking      |
| NFR-5: Browser Compatibility          | AC-6 (90%+ browser market)                             | 6-4-browser-compatibility-testing |
| NFR-6.2: Test Coverage                | AC-3 (90%+ coverage)                                   | 6-1-automated-test-suite          |

## Risks, Assumptions, Open Questions

### Risks

**Risk 1: Test Suite Runtime Exceeds 10 Seconds**
- **Impact:** Medium - Slows developer iteration, frustration with slow CI
- **Probability:** Low - Current Python tests run in ~5s, Go should be faster
- **Mitigation:** Parallel test execution (`go test -parallel=8`), subset of round-trip tests (not all 1,501 files), cache testdata/ in CI

**Risk 2: Visual Regression Testing is Subjective and Manual**
- **Impact:** Medium - Hard to automate "looks correct" validation, requires manual review
- **Probability:** High - Color science comparison is inherently subjective
- **Mitigation:** Document visual regression process clearly, use color delta E metrics (quantitative), test with 5-10 representative presets (not all 1,501), accept manual validation for MVP

**Risk 3: Browser Compatibility Matrix is Incomplete**
- **Impact:** Low - Some users may encounter unsupported browsers
- **Probability:** Medium - Testing all browser versions is impractical
- **Mitigation:** Focus on latest 2 versions of Chrome, Firefox, Safari (90%+ market), clear "unsupported browser" message, progressive enhancement strategy

**Risk 4: Performance Regression Detection is Noisy**
- **Impact:** Low - False positives block merges unnecessarily
- **Probability:** Medium - Benchmark variance can cause flaky failures
- **Mitigation:** Set threshold at 10% (not 5%), run benchmarks multiple times and average, allow manual override with justification

**Risk 5: CI/CD Pipeline Costs Exceed Free Tier**
- **Impact:** Low - GitHub Actions free tier limits (2,000 minutes/month)
- **Probability:** Low - Tests run in ~2 minutes, ~30 PRs/month = 60 minutes
- **Mitigation:** Optimize test runtime, cache dependencies, run full suite only on main branch (subset on PRs if needed)

### Assumptions

**Assumption 1: Sample files in testdata/ are representative**
- All 1,501 sample files accurately represent real-world preset usage
- Files cover edge cases and boundary conditions
- **Validation:** Community testing with user-submitted presets in later phases

**Assumption 2: 90% coverage threshold is sufficient**
- Code coverage ≥90% indicates high test quality
- Remaining 10% is error handling and edge cases that are hard to test
- **Validation:** Track defect density over time, adjust threshold if needed

**Assumption 3: Go benchmarks accurately reflect WASM performance**
- Go benchmark results predict WASM performance within 10-20%
- WASM overhead is consistent and predictable
- **Validation:** Compare Go benchmark to actual WASM execution in browser

**Assumption 4: Manual browser testing is acceptable for MVP**
- Automated browser testing (Selenium, Playwright) is overkill for simple UI
- Manual testing with 3 browsers is sufficient for MVP
- **Validation:** If bugs arise from browser differences, invest in automation

**Assumption 5: GitHub Actions free tier is sufficient**
- Current usage fits within 2,000 minutes/month limit
- Free tier provides adequate resources (2 CPU cores, 7GB RAM)
- **Validation:** Monitor usage in GitHub billing, upgrade if needed

### Open Questions

**Question 1: Should we implement automated visual regression testing?**
- **Context:** Tools like ImageMagick, SSIM can compare images programmatically
- **Tradeoff:** Complexity vs automation value
- **Decision Needed:** After MVP, evaluate if manual visual regression is too slow

**Question 2: What is the acceptable flake rate for tests?**
- **Context:** Occasional test failures due to timing, resource contention
- **Tradeoff:** Zero tolerance vs pragmatic acceptance of <1% flake rate
- **Decision Needed:** Monitor flake rate in first month, set policy

**Question 3: Should performance benchmarks be required for all PRs?**
- **Context:** Benchmarks add ~1 minute to CI runtime
- **Tradeoff:** Every PR validated vs faster CI for non-performance changes
- **Decision Needed:** Run benchmarks only on main branch or for PRs touching converter code?

**Question 4: How do we handle new browser versions breaking compatibility?**
- **Context:** Browsers auto-update, WebAssembly APIs may change
- **Tradeoff:** Continuous monitoring vs reactive fixes when users report issues
- **Decision Needed:** Set up automated browser version monitoring or rely on user reports?

**Question 5: Should we invest in E2E testing framework (Playwright)?**
- **Context:** Current plan is manual browser testing
- **Tradeoff:** Investment in automation vs simple manual verification
- **Decision Needed:** Evaluate after MVP based on regression frequency

## Test Strategy Summary

### Test Pyramid Structure

```
                    ┌─────────────────┐
                    │  Manual Tests   │ (Browser compat, visual regression)
                    │   ~5% effort    │
                    └─────────────────┘
                           ▲
                ┌──────────────────────┐
                │  Integration Tests   │ (Round-trip conversions)
                │     ~25% effort      │
                └──────────────────────┘
                           ▲
            ┌──────────────────────────────┐
            │      Unit Tests              │ (Format parsers, 1,501 files)
            │       ~70% effort            │
            └──────────────────────────────┘
```

### Testing Levels

**Level 1: Unit Tests (Highest Volume)**
- **Target:** All format parsers (np3, xmp, lrtemplate)
- **Method:** Table-driven tests with 1,501 sample files
- **Coverage:** Parse success, parameter extraction, range validation
- **Execution:** `go test ./internal/formats/...`
- **Runtime:** ~3 seconds
- **Success Criteria:** 100% pass rate on valid files

**Level 2: Integration Tests (Medium Volume)**
- **Target:** Bidirectional conversion paths (6 combinations)
- **Method:** Round-trip testing (A → B → A)
- **Coverage:** Parameter fidelity, format compatibility, edge cases
- **Execution:** `go test ./internal/converter/...`
- **Runtime:** ~4 seconds
- **Success Criteria:** ≥95% accuracy (tolerance ±1)

**Level 3: Performance Tests (Continuous)**
- **Target:** Conversion speed, memory usage, binary size
- **Method:** Go benchmarks with standardized test data
- **Coverage:** WASM (<100ms), CLI (<20ms), batch (<2s)
- **Execution:** `go test -bench=. -benchmem`
- **Runtime:** ~3 seconds
- **Success Criteria:** Meet or exceed performance targets

**Level 4: Visual Regression (Manual, Periodic)**
- **Target:** Color accuracy, visual similarity
- **Method:** Apply presets to reference images, compare outputs
- **Coverage:** 5-10 representative presets (portrait, landscape, etc.)
- **Execution:** Manual process with Lightroom/NX Studio
- **Runtime:** ~30 minutes per test cycle
- **Success Criteria:** Delta E <5, visual similarity ≥95%

**Level 5: Browser Compatibility (Manual, Release Validation)**
- **Target:** Web interface functionality across browsers
- **Method:** Manual testing in Chrome, Firefox, Safari
- **Coverage:** Upload, WASM load, conversion, download, UI
- **Execution:** Manual checklist per browser
- **Runtime:** ~15 minutes per browser
- **Success Criteria:** All features functional in 90%+ browser market

### Test Frameworks and Tools

**Primary: Go Testing Package**
- Standard library `testing` package (zero dependencies)
- Table-driven tests for data-driven validation
- Subtests for granular failure reporting
- Coverage analysis via `go test -cover`
- Benchmark framework via `go test -bench`

**CI/CD: GitHub Actions**
- Automated test execution on every PR/push
- Coverage threshold enforcement (≥90%)
- Performance regression detection (>10% slowdown fails)
- Artifact storage for benchmark history

**Visual Regression: Manual + ImageMagick (Optional)**
- Adobe Lightroom CC/Classic for preset application
- Nikon NX Studio for NP3 validation
- ImageMagick for image comparison (future automation)

**Browser Testing: Manual Verification**
- Chrome DevTools for WASM debugging
- Firefox DevTools for compatibility testing
- Safari Web Inspector for iOS/macOS validation

### Test Data Strategy

**Sample File Corpus:**
- 1,501 total sample files in `testdata/` directory
- 22 NP3 files (Nikon Picture Control)
- 913 XMP files (Lightroom CC presets)
- 566 lrtemplate files (Lightroom Classic presets)
- All files version controlled (reproducible tests)

**Test Data Characteristics:**
- Real-world presets from public sources
- Edge cases (min/max parameter values)
- Boundary conditions (empty fields, missing tags)
- Corrupted files (invalid magic bytes, truncated data)

**Test Data Organization:**
```
testdata/
├── np3/
│   ├── official/         # Nikon-provided samples
│   ├── community/        # User-created presets
│   └── edge-cases/       # Boundary/error cases
├── xmp/
│   ├── lightroom-cc/     # Adobe preset samples
│   ├── community/        # User-created presets
│   └── edge-cases/       # Boundary/error cases
└── lrtemplate/
    ├── lightroom-classic/
    ├── community/
    └── edge-cases/
```

### Coverage Goals

**Code Coverage:**
- Overall: ≥90%
- `internal/converter`: ≥95% (critical path)
- `internal/formats/*`: ≥90% (parser/generator logic)
- `cmd/*`: ≥70% (CLI/TUI interfaces, lower priority)

**Functional Coverage:**
- 100% of supported parameters tested
- 100% of error paths exercised
- 100% of format combinations validated
- 100% of browser compatibility requirements verified

### Regression Prevention Strategy

**Automated Regression Detection:**
1. Every bug fix requires regression test added to suite
2. Performance benchmarks tracked over time, >10% slowdown fails build
3. Coverage percentage cannot decrease (new code must include tests)
4. Round-trip accuracy cannot degrade (tolerance tightens over time)

**Manual Regression Validation:**
1. Visual regression testing repeated before each release
2. Browser compatibility re-validated with major browser updates
3. User-reported issues added to test corpus

### Test Execution Schedule

**On Every PR/Push:**
- Unit tests (1,501 file parsing)
- Integration tests (round-trip conversions)
- Coverage analysis (≥90% threshold)
- Performance benchmarks (regression detection)

**Before Release:**
- Visual regression testing (manual, 5-10 presets)
- Browser compatibility testing (Chrome, Firefox, Safari)
- End-to-end smoke testing (complete user workflow)

**Periodic (Monthly):**
- Full benchmark suite with historical comparison
- Test data corpus review (add new community presets)
- Browser version update testing (new releases)

### Success Metrics

**Quality Metrics:**
- Test pass rate: 100% (zero failures allowed in main branch)
- Coverage: ≥90% overall, ≥95% for core converter
- Round-trip accuracy: ≥95% (tolerance ±1 for rounding)
- Visual similarity: ≥95% (subjective assessment)

**Performance Metrics:**
- Test suite runtime: <10 seconds
- WASM conversion: <100ms
- CLI conversion: <20ms
- Batch processing: <2s for 100 files

**Reliability Metrics:**
- Flake rate: <1% (tests are deterministic)
- CI uptime: ≥99% (GitHub Actions SLA)
- Regression escape rate: <1% (bugs caught before production)

---

**Epic 6 establishes Recipe's quality foundation, ensuring 95%+ conversion accuracy, <100ms performance, and 90%+ browser compatibility through comprehensive automated testing, performance benchmarking, and validation infrastructure.**
