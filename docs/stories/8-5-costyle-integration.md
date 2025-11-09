# Story 8.5: Capture One CLI/TUI/Web Integration

Status: ready-for-dev

## Story

As a **photographer**,
I want **Capture One .costyle format fully integrated into Recipe's CLI, TUI, and web interfaces**,
so that **I can convert Capture One presets using my preferred interface (command line, terminal UI, or browser) with the same seamless experience as other supported formats (XMP, lrtemplate, NP3)**.

## Acceptance Criteria

**AC-1: CLI Integration**
- ✅ CLI accepts .costyle files: `recipe convert input.costyle --to xmp`
- ✅ CLI accepts .costylepack bundles: `recipe convert bundle.costylepack --to xmp-bundle`
- ✅ Format auto-detection works for .costyle and .costylepack extensions
- ✅ Help text includes Capture One examples: `recipe convert --help`
- ✅ Verbose output logs Capture One-specific parsing/generation steps
- ✅ JSON output mode includes Capture One format metadata

**AC-2: TUI Integration**
- ✅ TUI format menu includes "Capture One (.costyle)" option
- ✅ Format badge uses purple color (Capture One brand color)
- ✅ File browser detects and displays .costyle files correctly
- ✅ Batch mode supports .costylepack bundles (shows file count)
- ✅ Preview screen displays Capture One parameters (exposure, contrast, saturation, etc.)
- ✅ Progress indicators work for .costylepack batch conversions

**AC-3: Web Integration**
- ✅ Web UI accepts .costyle uploads via drag-and-drop
- ✅ Web UI accepts .costylepack uploads (batch conversion)
- ✅ Format detection automatically identifies .costyle files
- ✅ Parameter preview displays Capture One adjustments
- ✅ Target format selector includes Capture One option
- ✅ Download triggers correctly for generated .costyle files

**AC-4: Format Detection**
- ✅ Detect .costyle by extension (.costyle)
- ✅ Detect .costylepack by extension (.costylepack)
- ✅ Validate XML structure (check for xmpmeta, RDF, Description elements)
- ✅ Return format type "costyle" or "costylepack" from detector
- ✅ Format detection works in CLI, TUI, and Web contexts
- ✅ Handle ambiguous files gracefully (e.g., .xml extension but costyle structure)

**AC-5: Converter Integration**
- ✅ `Converter.Convert()` function handles "costyle" format type
- ✅ Route to `costyle.Parse()` for input files
- ✅ Route to `costyle.Generate()` for output files
- ✅ Handle .costylepack bundles (unpack → convert → pack if needed)
- ✅ Support all conversion paths:
  - costyle → xmp, np3, lrtemplate, dcp
  - xmp/np3/lrtemplate/dcp → costyle
  - costylepack → multi-file outputs
- ✅ Preserve bundle structure for pack → pack conversions

**AC-6: Documentation Updates**
- ✅ README.md includes Capture One format in supported formats list
- ✅ CLI help text (`recipe convert --help`) includes .costyle examples
- ✅ Web FAQ includes Capture One format questions
- ✅ Parameter mapping docs include .costyle mappings (already done in 8-1)
- ✅ Format compatibility matrix includes Capture One row/column
- ✅ Deployment docs note Capture One support (landing page, changelog)

**AC-7: Error Handling**
- ✅ Invalid .costyle XML → Clear error message ("malformed .costyle file")
- ✅ Corrupt .costylepack ZIP → Clear error message ("corrupt .costylepack bundle")
- ✅ Unsupported Capture One parameters → Warning logged, conversion continues
- ✅ Empty .costyle → Convert to neutral preset (all parameters zero)
- ✅ CLI/TUI/Web all display Capture One-specific errors consistently

**AC-8: Integration Testing**
- ✅ CLI integration test: Convert .costyle to all formats
- ✅ TUI integration test: Navigate to .costyle file, convert successfully
- ✅ Web integration test: Upload .costyle via drag-drop, convert, download
- ✅ Batch test: Convert .costylepack bundle to multiple outputs
- ✅ Round-trip test: CLI convert → Web upload → CLI convert (verify consistency)

## Tasks / Subtasks

### Task 1: Update Converter Package (AC-5)
- [ ] Open `internal/converter/converter.go`
- [ ] Add "costyle" and "costylepack" to format type constants:
  ```go
  const (
      FormatNP3 = "np3"
      FormatXMP = "xmp"
      FormatLRTemplate = "lrtemplate"
      FormatCostyle = "costyle"          // NEW
      FormatCostylepack = "costylepack"  // NEW
      // ...
  )
  ```
- [ ] Update `Convert()` function to handle Capture One formats:
  ```go
  func (c *Converter) Convert(inputData []byte, sourceFormat, targetFormat string) ([]byte, error) {
      // ... existing code ...

      // Parse input
      var recipe *universal.Recipe
      switch sourceFormat {
      case FormatCostyle:
          recipe, err = costyle.Parse(inputData)
      case FormatCostylepack:
          recipes, err := costyle.Unpack(inputData)
          // Handle multiple recipes (batch conversion)
      // ... other formats ...
      }

      // Generate output
      var outputData []byte
      switch targetFormat {
      case FormatCostyle:
          outputData, err = costyle.Generate(recipe)
      case FormatCostylepack:
          outputData, err = costyle.Pack(recipes, filenames)
      // ... other formats ...
      }

      return outputData, nil
  }
  ```
- [ ] Handle .costylepack batch conversions:
  - Unpack bundle → slice of recipes
  - Convert each recipe to target format
  - If target is bundle format, pack outputs
  - If target is individual format, return slice of outputs
- [ ] Import costyle package: `import "github.com/jpoechill/recipe/internal/formats/costyle"`

### Task 2: Update Format Detection (AC-4)
- [ ] Open `internal/formats/detection.go` (or equivalent detection module)
- [ ] Add .costyle detection logic:
  ```go
  func DetectFormat(data []byte, filename string) (string, error) {
      // Check extension first
      if strings.HasSuffix(filename, ".costyle") {
          return FormatCostyle, nil
      }
      if strings.HasSuffix(filename, ".costylepack") {
          return FormatCostylepack, nil
      }

      // Check XML magic bytes for .costyle
      if bytes.HasPrefix(data, []byte("<?xml")) {
          // Parse XML to check for Capture One structure
          if strings.Contains(string(data[:1000]), "<x:xmpmeta") &&
             strings.Contains(string(data[:1000]), "adobe:ns:meta") {
              return FormatCostyle, nil
          }
      }

      // Check ZIP magic bytes for .costylepack
      if bytes.HasPrefix(data, []byte{0x50, 0x4B, 0x03, 0x04}) {
          // Check if ZIP contains .costyle files
          // ... inspect ZIP contents ...
          return FormatCostylepack, nil
      }

      // ... other formats ...
  }
  ```
- [ ] Add Capture One to format validator (verify structure)
- [ ] Test detection with .costyle, .costylepack, and ambiguous files

### Task 3: CLI Integration (AC-1)
- [ ] Open `cmd/cli/convert.go` (Cobra command file)
- [ ] Verify format detection already works (should auto-detect .costyle)
- [ ] Update help text to include Capture One examples:
  ```go
  Long: `Convert photo editing presets between formats.

  Supported formats:
    - Nikon NP3 (.np3)
    - Adobe XMP (.xmp)
    - Lightroom Classic (.lrtemplate)
    - Capture One (.costyle, .costylepack)  // NEW
    - DNG Camera Profile (.dcp)

  Examples:
    # Convert Capture One style to XMP
    recipe convert preset.costyle --to xmp

    # Convert bundle to individual XMP files
    recipe convert bundle.costylepack --to xmp

    # Round-trip test
    recipe convert original.costyle temp.xmp && recipe convert temp.xmp output.costyle
  `,
  ```
- [ ] Test CLI commands:
  - `recipe convert sample.costyle --to xmp`
  - `recipe convert bundle.costylepack --to xmp`
  - `recipe convert sample.xmp --to costyle`
- [ ] Verify verbose output logs Capture One steps
- [ ] Verify JSON output includes Capture One metadata

### Task 4: TUI Integration (AC-2)
- [ ] Open `internal/tui/formats.go` (format menu definitions)
- [ ] Add Capture One format to menu:
  ```go
  var SupportedFormats = []Format{
      {Name: "Nikon NP3", Extension: ".np3", Color: YellowBadge},
      {Name: "Adobe XMP", Extension: ".xmp", Color: BlueBadge},
      {Name: "Lightroom Classic", Extension: ".lrtemplate", Color: CyanBadge},
      {Name: "Capture One", Extension: ".costyle", Color: PurpleBadge},  // NEW
      {Name: "DNG Camera Profile", Extension: ".dcp", Color: GreenBadge},
  }
  ```
- [ ] Update file browser to display .costyle files:
  - Add .costyle to file filter list
  - Show Capture One icon/badge next to .costyle files
- [ ] Update preview screen to display Capture One parameters:
  - Add parameter display logic (similar to XMP preview)
  - Show: Exposure, Contrast, Saturation, Temperature, Tint, Clarity
- [ ] Test batch mode with .costylepack bundles:
  - Display bundle file count: "Bundle: 5 presets"
  - Show progress: "Converting 3 of 5..."
- [ ] Verify TUI navigation and conversion flow

### Task 5: Web Integration (AC-3)
- [ ] Open `web/js/format-detection.js`
- [ ] Add .costyle detection:
  ```javascript
  function detectFormat(file) {
      const extension = file.name.split('.').pop().toLowerCase();

      if (extension === 'costyle') {
          return 'costyle';
      }
      if (extension === 'costylepack') {
          return 'costylepack';
      }

      // Check file contents for Capture One XML structure
      // ... read first 1KB, check for xmpmeta/RDF elements ...

      return 'unknown';
  }
  ```
- [ ] Update `web/js/file-upload.js`:
  - Accept .costyle and .costylepack file uploads
  - Show Capture One format badge on upload
- [ ] Update `web/js/parameter-preview.js`:
  - Parse Capture One parameters from .costyle XML
  - Display exposure, contrast, saturation, temperature, tint, clarity
  - Format values for display (e.g., "Exposure: +0.5", "Contrast: +15")
- [ ] Update `web/js/conversion.js`:
  - Handle .costyle → other formats conversion
  - Handle .costylepack bundle uploads (show file count)
  - Trigger download for generated .costyle files
- [ ] Update target format selector:
  - Add "Capture One (.costyle)" option to dropdown
- [ ] Test web UI flow:
  - Upload .costyle via drag-drop
  - Preview parameters
  - Select target format
  - Convert and download

### Task 6: Documentation Updates (AC-6)
- [ ] Update `README.md`:
  - Add Capture One to supported formats list:
    ```markdown
    ## Supported Formats

    - Nikon Picture Control (.np3)
    - Adobe XMP (.xmp)
    - Lightroom Classic (.lrtemplate)
    - **Capture One (.costyle, .costylepack)** ← NEW
    - DNG Camera Profile (.dcp)
    ```
  - Add Capture One conversion examples
- [ ] Update `web/index.html` (landing page):
  - Add Capture One logo/badge to format grid
  - Update feature list: "Convert Capture One styles to XMP, lrtemplate, NP3, DCP"
- [ ] Update `web/faq.html`:
  - Add FAQ entry: "Does Recipe support Capture One .costyle files?"
  - Answer: "Yes! Recipe supports both individual .costyle files and .costylepack bundles..."
- [ ] Update `docs/format-compatibility-matrix.md`:
  - Add Capture One row/column
  - Mark supported conversion paths (costyle ↔ xmp, costyle ↔ np3, etc.)
- [ ] Update `CHANGELOG.md`:
  - Add entry: "Added Capture One .costyle format support"

### Task 7: Error Handling (AC-7)
- [ ] Implement Capture One-specific error messages in converter:
  ```go
  if err := costyle.Parse(data); err != nil {
      return nil, fmt.Errorf("failed to parse .costyle file: %w", err)
  }
  ```
- [ ] Add error message mapping in CLI:
  - "malformed .costyle file" → User-friendly message with fix suggestions
  - "corrupt .costylepack bundle" → Suggest re-downloading file
- [ ] Add error display in TUI:
  - Show error dialog with Capture One icon
  - Display specific error (e.g., "Invalid XML structure in preset.costyle")
- [ ] Add error display in Web UI:
  - Show error banner with red styling
  - Display user-friendly error message
  - Provide help link to FAQ or documentation
- [ ] Test error paths:
  - Upload corrupt .costyle file (malformed XML)
  - Upload corrupt .costylepack file (truncated ZIP)
  - Convert .costyle with unsupported parameters (verify warning logged)

### Task 8: Integration Testing (AC-8)
- [ ] Write CLI integration test:
  ```go
  func TestCLI_ConvertCostyle(t *testing.T) {
      // Convert .costyle to XMP
      output, err := runCLI("convert", "testdata/sample.costyle", "--to", "xmp")
      require.NoError(t, err)
      assert.Contains(t, output, "Conversion successful")

      // Verify output XMP file exists and is valid
      xmpData, err := os.ReadFile("sample.xmp")
      require.NoError(t, err)
      recipe, err := xmp.Parse(xmpData)
      require.NoError(t, err)
      assert.NotNil(t, recipe)
  }
  ```
- [ ] Write TUI integration test (if automated TUI testing available):
  - Navigate to .costyle file in file browser
  - Select file, choose target format
  - Verify conversion completes successfully
- [ ] Write Web integration test (using Playwright or Selenium):
  ```javascript
  test('Upload and convert .costyle file', async ({ page }) => {
      await page.goto('http://localhost:8080');

      // Upload .costyle file
      await page.setInputFiles('input[type=file]', 'testdata/sample.costyle');

      // Verify format detected
      await expect(page.locator('.format-badge')).toHaveText('Capture One');

      // Select target format
      await page.selectOption('#target-format', 'xmp');

      // Click convert
      await page.click('#convert-button');

      // Wait for download
      const download = await page.waitForEvent('download');
      assert.equal(download.suggestedFilename(), 'sample.xmp');
  });
  ```
- [ ] Write batch conversion test:
  - Convert .costylepack to multiple XMP files
  - Verify all files converted successfully
  - Verify file count matches bundle count
- [ ] Write round-trip integration test:
  - CLI: Convert .costyle → .xmp
  - Web: Upload .xmp, convert → .costyle
  - CLI: Convert .costyle → verify matches original
- [ ] Run all integration tests in CI

### Task 9: Final Validation (AC-1 to AC-8)
- [ ] Manual testing checklist:
  - [ ] CLI: Convert .costyle to all formats (xmp, np3, lrtemplate, dcp)
  - [ ] CLI: Convert .costylepack bundle to XMP files
  - [ ] TUI: Browse to .costyle file, preview parameters, convert
  - [ ] TUI: Batch convert .costylepack bundle
  - [ ] Web: Drag-drop .costyle file, preview, convert, download
  - [ ] Web: Upload .costylepack bundle, convert batch
  - [ ] Format detection: Upload .xml file with costyle structure (verify detected)
  - [ ] Error handling: Upload corrupt .costyle (verify error message)
- [ ] Verify all interfaces show consistent behavior:
  - Same conversion results (CLI vs. TUI vs. Web)
  - Same error messages
  - Same parameter preview format
- [ ] Performance check:
  - CLI conversion <1 second for single .costyle
  - Web upload and preview <2 seconds
  - Batch .costylepack conversion <5 seconds for 50 files

## Dev Notes

### Learnings from Previous Story

**From Story 8-4-costyle-round-trip-testing (Status: drafted)**

- **Round-Trip Testing**: Automated tests verify 95%+ accuracy
- **Test Samples**: Real-world .costyle files in `testdata/costyle/real-world/`
- **Accuracy Metrics**: Parameter-by-parameter comparison with tolerance
- **Documentation**: known-conversion-limitations.md documents lossy conversions
- **Visual Validation**: Manual spot-check in Capture One Pro software

**Reuse from Story 8-4:**
- Test samples (use for integration testing)
- Round-trip test pattern (verify CLI/TUI/Web conversions preserve accuracy)
- Error handling patterns (consistent error messages across interfaces)

[Source: docs/stories/8-4-costyle-round-trip-testing.md#Dev-Notes]

### Architecture Alignment

**Tech Spec Epic 8 Alignment:**

Story 8-5 implements **AC-5 (CLI/TUI/Web Integration)** from tech-spec-epic-8.md.

**Integration Flow:**
```
CLI Input → Format Detection → Converter.Convert() → costyle.Parse/Generate → Output
TUI Input → Format Detection → Converter.Convert() → costyle.Parse/Generate → Output
Web Input → Format Detection → Converter.Convert() → costyle.Parse/Generate → Output
```

**Converter Extension:**
```go
// Before (Epics 1-7)
func Convert(data []byte, source, target string) ([]byte, error) {
    // Supports: np3, xmp, lrtemplate
}

// After (Epic 8)
func Convert(data []byte, source, target string) ([]byte, error) {
    // Supports: np3, xmp, lrtemplate, costyle, costylepack
}
```

**Format Detection Priority:**
1. Extension check (.costyle, .costylepack)
2. Magic bytes (XML header, ZIP header)
3. Content structure (xmpmeta/RDF elements)

[Source: docs/tech-spec-epic-8.md#System-Architecture-Alignment]

### Interface Patterns (from Epics 1-7)

**CLI Pattern (from Epic 3):**
- Auto-detect format from extension
- `--to` flag specifies target format
- Verbose logging with `--verbose` flag
- JSON output with `--json` flag
- Help text includes format examples

**TUI Pattern (from Epic 4):**
- Format badge with color coding (purple for Capture One)
- File browser with extension filtering
- Parameter preview screen
- Batch mode progress indicators
- Error dialogs with format-specific icons

**Web Pattern (from Epic 2):**
- Drag-and-drop file upload
- Format detection and badge display
- Parameter preview panel
- Target format selector dropdown
- Download button triggers file save

**Reuse Existing Components:**
- `internal/converter/converter.go` - Extend Convert() function
- `internal/formats/detection.go` - Add .costyle detection
- `cmd/cli/convert.go` - Add help text examples
- `internal/tui/formats.go` - Add Capture One format entry
- `web/js/format-detection.js` - Add .costyle detection logic

[Source: docs/tech-spec-epic-2.md, docs/tech-spec-epic-3.md, docs/tech-spec-epic-4.md]

### Project Structure Notes

**Modified Files (Story 8-5):**
```
internal/converter/
└── converter.go           # Add costyle format handling (MODIFIED)

internal/formats/
└── detection.go           # Add .costyle detection (MODIFIED)

cmd/cli/
└── convert.go             # Update help text, add examples (MODIFIED)

internal/tui/
└── formats.go             # Add Capture One format entry (MODIFIED)

web/js/
├── format-detection.js    # Add .costyle detection (MODIFIED)
├── file-upload.js         # Accept .costyle uploads (MODIFIED)
├── parameter-preview.js   # Display Capture One params (MODIFIED)
└── conversion.js          # Handle .costyle conversions (MODIFIED)

docs/
├── README.md              # Add Capture One to supported formats (MODIFIED)
├── format-compatibility-matrix.md  # Add Capture One row/column (MODIFIED)
└── CHANGELOG.md           # Add Epic 8 entry (MODIFIED)

web/
├── index.html             # Add Capture One badge (MODIFIED)
└── faq.html               # Add Capture One FAQ (MODIFIED)
```

**No New Files Created**: This story integrates existing components.

[Source: docs/tech-spec-epic-8.md#Components]

### Testing Strategy

**Unit Tests:**
- Converter tests: `TestConvert_Costyle()` for all conversion paths
- Detection tests: `TestDetectFormat_Costyle()` for .costyle and .costylepack

**Integration Tests (Required for AC-8):**
- CLI test: `TestCLI_ConvertCostyle()` - CLI conversion flow
- TUI test: `TestTUI_ConvertCostyle()` - TUI navigation and conversion (if automated TUI testing available)
- Web test: `TestWeb_ConvertCostyle()` - Web upload, convert, download (Playwright/Selenium)
- Batch test: `TestBatch_Costylepack()` - Bundle conversion
- Round-trip test: CLI → Web → CLI consistency

**Manual Validation:**
- Test all interfaces (CLI, TUI, Web) with same .costyle file
- Verify consistent results across interfaces
- Test error handling (corrupt files, unsupported parameters)
- Performance check (conversion times <1s for CLI, <2s for Web)

[Source: docs/tech-spec-epic-8.md#Test-Strategy-Summary]

### Known Risks

**RISK-10: Interface inconsistencies across CLI/TUI/Web**
- **Impact**: Confusing user experience, different behavior per interface
- **Mitigation**: Centralize conversion logic in Converter.Convert() (single source of truth)
- **Testing**: Round-trip integration tests verify consistency

**RISK-11: Format detection conflicts with XMP**
- **Impact**: .costyle XML structure may be misidentified as XMP
- **Mitigation**: Check for Capture One-specific elements (xmpmeta namespace, RDF structure)
- **Priority**: Extension check first, content check as fallback

**RISK-12: Web UI performance with large .costylepack bundles**
- **Impact**: Browser may hang on 50+ file bundles
- **Mitigation**: Show progress indicator, use Web Workers for conversion (if available)
- **Target**: <5 seconds for 50-file bundle (same as CLI/TUI)

[Source: docs/tech-spec-epic-8.md#Risks-Assumptions-Open-Questions]

### References

- [Source: docs/tech-spec-epic-8.md#Acceptance-Criteria] - AC-5: CLI/TUI/Web Integration
- [Source: internal/converter/converter.go] - Converter.Convert() function (extend for .costyle)
- [Source: internal/formats/detection.go] - Format detection logic (add .costyle)
- [Source: cmd/cli/convert.go] - CLI convert command (add help text)
- [Source: internal/tui/formats.go] - TUI format menu (add Capture One)
- [Source: web/js/format-detection.js] - Web format detection (add .costyle)
- [Source: docs/stories/3-2-convert-command.md] - CLI conversion patterns (reference)
- [Source: docs/stories/4-1-bubbletea-file-browser.md] - TUI file browser patterns (reference)
- [Source: docs/stories/2-2-file-upload-handling.md] - Web upload patterns (reference)

## Dev Agent Record

### Context Reference

- `docs/stories/8-5-costyle-integration.context.xml` - Generated 2025-11-09

### Agent Model Used

claude-sonnet-4-5-20250929

### Debug Log References

### Completion Notes List

### File List
