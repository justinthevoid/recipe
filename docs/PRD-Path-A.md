# recipe - Product Requirements Document (Path A Enhancements)

**Author:** Justin
**Date:** 2025-11-08
**Version:** 1.0
**Enhancement Track:** Path A - Format Expansion & UX Enhancement

---

## Executive Summary

Path A represents Recipe's evolution from a solid conversion utility to a polished, professional-grade tool that photographers love to use. Building on Recipe's core achievement - the only universal converter supporting Nikon's NP3 format - Path A expands format coverage to include Capture One (serving professional photographers) and DCP camera profiles, while dramatically improving the user experience.

**The Enhancement Vision:** Recipe already converts presets with 98%+ accuracy across three formats. Path A takes this foundation and adds the capabilities that transform a functional tool into a delightful experience: visual preset previews, expanded format support for professionals, and a web interface that feels modern and polished.

**Core Value Proposition:**
- **For Nikon → Capture One users**: Finally, a way to use your Picture Control recipes in professional editing software
- **For all users**: See what a preset does before converting (no more blind conversions)
- **For Recipe itself**: A professional, modern interface that matches the quality of the conversion engine

### What Makes This Special

**The Magic Moments:**

1. **Visual Discovery** - A photographer drags a preset file into Recipe. Instead of blindly converting, they see a before/after preview on sample images - the preset's warm tones, lifted shadows, and film-like grain instantly visible. They adjust a slider to compare, thinking "Yes, this is exactly what I want." They click convert with confidence.

2. **Professional Polish** - A Capture One user visits Recipe for the first time. The interface feels modern and intentional - clean typography, smooth interactions, visual format badges. Batch uploads show progress elegantly. The experience whispers "this tool was crafted with care" rather than "this was built by a developer for developers."

**What Makes Path A Special:**
- **Confidence through transparency**: Preview removes the guesswork from conversion
- **Craft over features**: Every interaction polished, every detail considered
- **Expanding reach**: Capture One support brings Recipe to professional photographers
- **Maintaining principles**: Still privacy-first, still local-only, still zero dependencies

---

## Project Classification

**Technical Type:** Web Application Enhancement (CLI/TUI/Web multi-interface)
**Domain:** Photography / Digital Asset Management
**Complexity:** Medium (expanding existing architecture with new formats and UI)

**Project Context:**
Path A is an enhancement initiative for an existing, proven product. Recipe already has:
- Established architecture (hub-and-spoke conversion pattern)
- Production deployment (Cloudflare Pages for web, GitHub Releases for binaries)
- Test infrastructure (1,531 real sample files, 89.5% coverage)
- User base (GitHub stars, web traffic)

This PRD focuses on **strategic expansion** rather than greenfield development. The foundation is solid; Path A builds upward with four coordinated enhancements that share the goal of making Recipe more comprehensive and delightful to use.

{{#if domain_context_summary}}

### Domain Context

{{domain_context_summary}}
{{/if}}

---

## Success Criteria

Path A success is measured by **craft and utility**, not vanity metrics. Since this is a personal project driven by continuous improvement, success means delivering enhancements that genuinely improve the tool and maintain code quality.

**Primary Success Criteria:**

1. **Preview Feature Utility** ✓
   - Users can make informed conversion decisions based on visual preview
   - Preview accuracy is sufficient for understanding preset characteristics (not necessarily pixel-perfect)
   - Feature is used (not ignored as decoration) - evidenced by user feedback or analytics
   - Implementation is performant (<100ms preview rendering)

2. **Code Quality Maintained** ✓
   - Hub-and-spoke architecture preserved and extended cleanly
   - Test coverage remains ≥85% across internal packages
   - All conversions maintain <100ms performance target
   - New format packages follow established patterns (`internal/formats/{format}/parse.go`, `generate.go`)
   - No technical debt introduced; refactoring opportunities identified and documented

3. **Professional Web Interface** ✓
   - Interface inspires confidence on first visit (modern, intentional design)
   - Mobile-responsive (works elegantly on phones/tablets)
   - Smooth interactions (no janky animations, proper loading states)
   - Accessibility basics covered (keyboard navigation, screen reader friendly)
   - Feel matches quality of conversion engine ("crafted with care")

**Secondary Success Indicators:**

- At least one Capture One user successfully converts presets
- GitHub stars or positive user feedback increases
- Zero regression bugs in existing NP3/XMP/lrtemplate conversion
- Documentation comprehensively updated (user guides + technical docs)

{{#if business_metrics}}

### Business Metrics

{{business_metrics}}
{{/if}}

---

## Product Scope

### MVP - Minimum Viable Product (Path A Complete)

Path A consists of four integrated enhancements, all essential for the enhancement vision:

**Epic 1: Capture One Format Support**
- Parse .costyle files (XML-based Capture One preset format)
- Generate .costyle files from UniversalRecipe intermediate representation
- Support .costylepack bundles (zip archives containing multiple styles)
- Round-trip conversion testing with real Capture One sample files
- Integration across all interfaces (CLI, TUI, Web)
- Update converter.Convert() to support new format

**Epic 2: DCP Camera Profile Support**
- Parse DCP (DNG Camera Profile) files per Adobe DNG Specification 1.6
- Generate DCP files from UniversalRecipe color adjustments
- Support XML-based camera profile structure
- Validate compatibility with Adobe Camera Raw and Lightroom
- Document DCP parameter mapping and limitations

**Epic 3: Enhanced Web UI/UX**
- Redesigned landing page with modern visual design
- Visual format badges (color-coded pills showing supported formats)
- Batch file upload with drag-and-drop support
- Progress indicators for multi-file conversions
- Mobile-responsive design (works on phones, tablets, desktop)
- Improved conversion flow with clear format selection
- Before/after comparison slider UI component

**Epic 4: Image Preview System (Phase 1)**
- CSS filter-based preview using browser native capabilities
- Map UniversalRecipe parameters to CSS filter values:
  - Exposure → brightness()
  - Contrast → contrast()
  - Saturation → saturate()
  - Hue adjustments → hue-rotate()
- Instant preview on bundled reference images (portrait, landscape, product)
- Approximate accuracy (good enough for decision-making)
- Clear communication that preview is approximate

**Acceptance Criteria:**
- All four epics functionally complete
- Quality gates pass (coverage, performance, testing)
- Documentation updated
- Deployed to production (web) and released (CLI/TUI binaries)

### Growth Features (Post-MVP / Path A)

Features deferred from Path A but considered for future phases:

**Enhanced Preview (Epic 4 Phase 2):**
- WebAssembly-based accurate preview using Photon library
- Pixel-perfect preview matching actual conversion output
- Support for complex adjustments (tone curves, HSL, parametric adjustments)
- Bundle size optimization (~500KB target)

**Additional Format Support:**
- DxO PhotoLab .dop format
- ON1 Photo RAW .onpreset format
- Affinity Photo .afphoto presets
- Additional format requests from community

**Advanced Batch Operations:**
- Batch rename presets during conversion
- Batch metadata editing
- Conversion profiles (save conversion settings for reuse)

### Vision (Future / Paths B & C)

Long-term capabilities explored in roundtable but deferred to future enhancement tracks:

**Path B - Market Expansion:**
- 3D LUT export (.cube format for DaVinci Resolve, video color grading)
- API/REST endpoints for third-party integrations
- Browser extensions for one-click conversion from preset marketplaces
- Mobile native apps (iOS/Android)
- Lightroom plugin for in-app conversion

**Path C - Platform Features:**
- Preset library management and organization
- Search and tag system for preset collections
- Cloud sync/backup (using user's own storage: S3, Dropbox)
- Collaborative features or preset sharing
- Preset marketplace integration
- Community-contributed presets

---

{{#if domain_considerations}}

## Domain-Specific Requirements

{{domain_considerations}}

This section shapes all functional and non-functional requirements below.
{{/if}}

---

{{#if innovation_patterns}}

## Innovation & Novel Patterns

{{innovation_patterns}}

### Validation Approach

{{validation_approach}}
{{/if}}

---

{{#if project_type_requirements}}

## {{project_type}} Specific Requirements

{{project_type_requirements}}

{{#if endpoint_specification}}

### API Specification

{{endpoint_specification}}
{{/if}}

{{#if authentication_model}}

### Authentication & Authorization

{{authentication_model}}
{{/if}}

{{#if platform_requirements}}

### Platform Support

{{platform_requirements}}
{{/if}}

{{#if device_features}}

### Device Capabilities

{{device_features}}
{{/if}}

{{#if tenant_model}}

### Multi-Tenancy Architecture

{{tenant_model}}
{{/if}}

{{#if permission_matrix}}

### Permissions & Roles

{{permission_matrix}}
{{/if}}
{{/if}}

---

## User Experience Principles

Path A's enhanced web interface is guided by principles that reflect Recipe's core values while elevating the user experience to professional standards.

### Design Philosophy

**Craft Over Flash**
The interface should feel intentionally designed, not over-designed. Every visual element serves a purpose. Smooth, subtle animations guide attention without distracting. Typography is clean and readable. Color is used deliberately to communicate format types, status, and hierarchy - not as decoration.

**Confidence Through Clarity**
Photographers should never feel uncertain about what Recipe is doing or what will happen next. Upload states are obvious. Conversion progress is visible. Format badges instantly communicate compatibility. Preview shows what a preset does before commitment. Error messages are helpful, not cryptic.

**Privacy-First Transparency**
The interface actively communicates Recipe's privacy-first architecture. Users see that processing happens locally - no spinners waiting for server responses, no "uploading..." states. The experience reinforces trust through instant, local processing.

**Progressive Disclosure**
Start simple, reveal complexity only when needed. The landing page shows core value immediately. Batch processing features appear when multiple files are dragged. Advanced options are accessible but not overwhelming. Mobile users see streamlined flows; desktop users get full capabilities.

### Visual Personality

**Modern but Timeless**
- Clean sans-serif typography (system fonts for performance)
- Generous whitespace that guides the eye
- Subtle shadows and depth cues (not flat, not skeuomorphic)
- Responsive grid that adapts elegantly across devices
- Color palette that feels photographic: neutral grays, accent colors derived from format types

**Format Badge System**
Visual format badges are a key differentiator - color-coded pills that make format compatibility instantly recognizable:
- **NP3**: Nikon yellow (#FFC107) - warm, inviting (Nikon's brand color)
- **XMP**: Adobe blue (#0073E6) - trustworthy, professional
- **lrtemplate**: Classic magenta (#D81B60) - creative energy
- **Capture One**: Elegant purple (#9C27B0) - professional sophistication
- **DCP**: Technical green (#4CAF50) - camera precision

### Key Interactions

**1. File Upload & Conversion Flow**

**Initial Landing:**
- Hero section immediately shows value: "Convert Photo Presets. Instantly. Privately."
- Visual format badges prominently displayed
- Single large drop zone: "Drop preset files here or click to browse"
- Supported formats listed clearly with badges

**Drag-and-Drop Experience:**
- Drop zone highlights on drag-over with subtle scale animation
- Multiple files create individual cards (batch mode)
- Each card shows: filename, detected format badge, file size, status
- Clear "Convert to..." dropdown for each file (or batch convert all)

**Conversion Process:**
- Instant format detection (happens client-side immediately)
- Progress indicator for batch conversions: "Converting 3 of 10..."
- Per-file status: queued → processing → complete → download ready
- Smooth transitions between states (no jarring updates)

**Completion:**
- Success state with download buttons per file
- Option to "Convert More" or "Start Fresh"
- Batch download as ZIP for multiple files

**2. Preview Feature**

**Before Converting:**
- After file upload, "Preview" button appears alongside "Convert"
- Click preview → modal opens with before/after slider
- Reference images: Portrait, Landscape, Product (tabs to switch)
- Draggable slider reveals original vs. preset-applied image
- Preset parameters displayed: "Exposure +0.7 • Contrast +15 • Warmth +10"
- Clear label: "Approximate preview using CSS filters"

**Interaction Details:**
- Smooth slider drag with visual feedback
- Keyboard accessible (arrow keys move slider)
- Mobile: Tap-and-hold to compare, or tap sides to snap
- "Looks good? Convert now" button in modal
- Close modal returns to conversion screen

**3. Mobile-Responsive Design**

**Phone (< 768px):**
- Stacked layout, single column
- Simplified upload: tap to browse (drag-drop optional)
- Preview slider uses full width
- Format badges stack vertically if needed
- Batch files shown as list (not grid)

**Tablet (768px - 1024px):**
- Two-column grid for batch files
- Preview slider uses 80% width modal
- Format badges in horizontal row

**Desktop (> 1024px):**
- Three-column grid for batch operations
- Side-by-side comparison in preview (optional layout)
- Full feature set visible without scrolling

**4. Accessibility Fundamentals**

- Keyboard navigation: Tab through all interactive elements
- Focus indicators: Clear outline on focused elements (not browser default blue)
- Screen reader labels: All icons have aria-labels
- Color contrast: WCAG AA compliant (4.5:1 for text)
- Skip to main content link for keyboard users
- Error messages read by screen readers
- Form labels properly associated

---

## Functional Requirements

All functional requirements are organized by epic and numbered for traceability. Each requirement includes acceptance criteria that will be decomposed into testable user stories.

### Epic 1: Capture One Format Support

**FR-1.1: Parse Capture One .costyle Files**

The system shall parse XML-based Capture One .costyle preset files and extract all supported adjustments into UniversalRecipe intermediate representation.

*Acceptance Criteria:*
- Parse XML structure per Capture One style specification
- Extract exposure, contrast, saturation, temperature, tint, clarity adjustments
- Extract color balance (shadows, midtones, highlights)
- Extract tone curve points if present
- Handle missing or optional parameters gracefully
- Validate XML structure and report parsing errors
- Support .costyle format versions currently in use (2023-2025)

**FR-1.2: Generate Capture One .costyle Files**

The system shall generate valid Capture One .costyle XML files from UniversalRecipe representation.

*Acceptance Criteria:*
- Generate valid XML structure matching Capture One specification
- Map UniversalRecipe parameters to .costyle equivalents
- Include required XML elements (version, metadata)
- Generate human-readable XML (formatted, not minified)
- Validate generated XML against schema
- Generated files load successfully in Capture One software
- Handle edge cases (missing parameters, out-of-range values)

**FR-1.3: Support .costylepack Bundles**

The system shall parse and generate .costylepack files (ZIP archives containing multiple .costyle presets).

*Acceptance Criteria:*
- Unzip .costylepack archives and extract individual .costyle files
- Parse each .costyle file within bundle
- Generate .costylepack by bundling multiple .costyle files into ZIP
- Maintain bundle metadata (name, description if present)
- Handle large bundles (50+ styles) efficiently
- Validate ZIP structure and report extraction errors

**FR-1.4: Round-Trip Conversion Testing**

The system shall support round-trip conversion (Capture One → UniversalRecipe → Capture One) with minimal data loss.

*Acceptance Criteria:*
- Round-trip conversion preserves 95%+ of parameter values
- Key adjustments (exposure, contrast, saturation) preserved exactly
- Document known limitations of lossy conversions
- Test suite includes real-world .costyle samples from Etsy/marketplaces
- Automated tests verify round-trip accuracy for all test files
- Visual validation in Capture One software confirms output quality

**FR-1.5: CLI/TUI/Web Integration**

The system shall support Capture One format conversion across all interfaces (CLI, TUI, Web).

*Acceptance Criteria:*
- CLI: `recipe convert input.costyle --to xmp` works correctly
- TUI: Format menu includes Capture One option
- Web: Upload .costyle files via drag-drop, convert to other formats
- Format detection automatically identifies .costyle files
- Converter.Convert() function extended to handle Capture One format
- Help text and documentation updated for new format

### Epic 2: DCP Camera Profile Support

**FR-2.1: Parse DCP Files**

The system shall parse DCP (DNG Camera Profile) files per Adobe DNG Specification 1.6.

*Acceptance Criteria:*
- Read TIFF-based DCP file structure
- Extract XML camera profile data from TIFF tags
- Parse color matrices (forward, color, calibration)
- Parse tone curve adjustments
- Parse hue/saturation/value tables if present
- Handle DCP v1.x format variations
- Validate DCP structure and report parsing errors
- Support both embedded (in DNG) and standalone DCP files

**FR-2.2: Generate DCP Files**

The system shall generate valid DCP files from UniversalRecipe color adjustments.

*Acceptance Criteria:*
- Generate TIFF-based DCP file structure per Adobe spec
- Embed XML camera profile data in TIFF tags
- Map UniversalRecipe color parameters to DCP equivalents
- Generate required matrices (identity matrices if not calibrating)
- Create tone curves from exposure/contrast/highlights/shadows
- Validate generated DCP against Adobe spec
- Generated DCPs load in Adobe Camera Raw and Lightroom
- Document mapping limitations and best practices

**FR-2.3: DCP Parameter Mapping**

The system shall define clear mapping between UniversalRecipe and DCP parameters.

*Acceptance Criteria:*
- Document which UniversalRecipe parameters map to DCP
- Identify unsupported DCP features (e.g., dual illuminant profiles)
- Define conversion formulas for tone curves
- Handle color space transformations correctly
- Create reference documentation for DCP mapping
- Include examples of common conversions
- Test mapping with real DCP samples from Adobe

**FR-2.4: Compatibility Validation**

The system shall validate DCP compatibility with Adobe software.

*Acceptance Criteria:*
- Generated DCPs load without errors in Adobe Camera Raw
- Generated DCPs load without errors in Lightroom Classic
- Preset adjustments render visually similar to original
- Performance: DCP generation completes in <200ms
- Test with multiple camera models (Nikon, Canon samples)
- Document known compatibility issues or edge cases

### Epic 3: Enhanced Web UI/UX

**FR-3.1: Redesigned Landing Page**

The system shall present a modern, visually appealing landing page that communicates Recipe's value immediately.

*Acceptance Criteria:*
- Hero section with clear value proposition: "Convert Photo Presets. Instantly. Privately."
- Visual format badges displayed prominently with brand colors
- Single-page layout (no navigation to other pages for core conversion)
- Clean typography using system fonts for performance
- Responsive design: works on mobile (320px+), tablet (768px+), desktop (1024px+)
- Fast load time: <2 seconds on 3G connection
- No external dependencies (no CDN fonts, no analytics trackers)

**FR-3.2: Visual Format Badges**

The system shall display color-coded format badges for instant format recognition.

*Acceptance Criteria:*
- Badge system implemented with defined colors:
  - NP3: #FFC107 (Nikon yellow)
  - XMP: #0073E6 (Adobe blue)
  - lrtemplate: #D81B60 (Magenta)
  - Capture One: #9C27B0 (Purple)
  - DCP: #4CAF50 (Green)
- Badges shown on landing page, in upload cards, in conversion dropdowns
- Accessible: Color not sole indicator (includes format name text)
- Responsive: Badges scale/stack appropriately on mobile
- Consistent styling across all interface elements

**FR-3.3: Batch File Upload with Drag-and-Drop**

The system shall support uploading multiple preset files simultaneously via drag-and-drop or file picker.

*Acceptance Criteria:*
- Large drop zone on landing page invites drag-and-drop
- Visual feedback on drag-over (highlight, scale animation)
- Support multiple file selection via file picker
- Accept all supported formats (.np3, .xmp, .lrtemplate, .costyle, .dcp)
- Reject unsupported files with clear error message
- No file size limit (client-side processing)
- Display uploaded files as individual cards in grid layout
- Each card shows: filename, detected format badge, file size, conversion status

**FR-3.4: Progress Indicators**

The system shall show clear progress indicators during multi-file conversions.

*Acceptance Criteria:*
- Batch conversion shows overall progress: "Converting 3 of 10..."
- Per-file status indicators: queued, processing, complete, error
- Smooth transitions between states (no jarky updates)
- Visual feedback during processing (spinner, progress bar)
- Completion state with download buttons
- Error state shows specific error message per file
- Users can cancel in-progress batch conversions

**FR-3.5: Mobile-Responsive Design**

The system shall adapt layout elegantly across device sizes.

*Acceptance Criteria:*
- Mobile (<768px): Single column, stacked layout, tap-to-browse upload
- Tablet (768-1024px): Two-column grid for batch files
- Desktop (>1024px): Three-column grid, full features visible
- Touch-friendly targets (44px minimum) on mobile
- No horizontal scrolling on any device size
- Readable text without zooming (16px minimum body text)
- Test on real devices: iPhone, Android, iPad, Desktop browsers

**FR-3.6: Before/After Comparison Slider**

The system shall provide an interactive slider to compare original vs. preset-applied images.

*Acceptance Criteria:*
- Slider implemented with draggable handle
- Smooth drag interaction with visual feedback
- Keyboard accessible (arrow keys move slider)
- Mobile: Tap-and-hold to compare, or tap sides to snap 50/50
- Slider position persists during modal session
- Clear visual indicators (before/after labels)
- Works across all browsers (Chrome, Firefox, Safari, Edge)

**FR-3.7: Improved Conversion Flow**

The system shall streamline the conversion process with clear format selection and immediate results.

*Acceptance Criteria:*
- Format detection happens instantly on upload (client-side)
- "Convert to..." dropdown shows only valid target formats
- Batch convert: Single action converts all files to same format
- Individual convert: Per-file format selection available
- Conversion happens instantly (<100ms per file)
- Download buttons appear immediately after conversion
- Option to "Convert More" without page reload
- "Start Fresh" clears all files and resets interface

### Epic 4: Image Preview System (Phase 1)

**FR-4.1: CSS Filter-Based Preview**

The system shall render approximate preset previews using browser-native CSS filters.

*Acceptance Criteria:*
- Map UniversalRecipe parameters to CSS filter functions:
  - Exposure → brightness()
  - Contrast → contrast()
  - Saturation → saturate()
  - Hue adjustments → hue-rotate()
  - Temperature/tint → sepia() + hue-rotate() approximation
- Preview renders in <100ms (instant, no processing delay)
- Preview updates in real-time as preset is selected
- Clear label: "Approximate preview using CSS filters"
- Works across all modern browsers (Chrome, Firefox, Safari, Edge)

**FR-4.2: Reference Image Bundle**

The system shall include bundled reference images for preview demonstration.

*Acceptance Criteria:*
- Three reference images: Portrait, Landscape, Product/Still-life
- Images optimized for web (<200KB each, total <600KB)
- Images representative of common photography genres
- Images embedded in web bundle (no external requests)
- Licensing: Public domain or created specifically for Recipe
- Images work well with typical preset adjustments (neutral starting point)

**FR-4.3: Preview Modal Interface**

The system shall display preview in a modal dialog with intuitive controls.

*Acceptance Criteria:*
- Modal opens when "Preview" button clicked after upload
- Modal shows before/after slider with reference image
- Tabs to switch between reference images (Portrait, Landscape, Product)
- Preset parameters displayed: "Exposure +0.7 • Contrast +15 • Warmth +10"
- "Convert now" button in modal proceeds to conversion
- Close/cancel button returns to upload screen
- Modal keyboard accessible (Esc to close, Tab navigation)
- Mobile: Full-screen modal, touch-friendly controls

**FR-4.4: Accuracy Communication**

The system shall clearly communicate that preview is approximate, not pixel-perfect.

*Acceptance Criteria:*
- Label: "Approximate preview using CSS filters"
- Tooltip/help text explains CSS filter limitations
- No misleading claims about preview accuracy
- Documentation explains preview vs. actual conversion differences
- Preview limitations listed (e.g., tone curves not supported in Phase 1)
- User expectations managed through transparent communication

**FR-4.5: Performance Optimization**

The system shall render preview instantly without noticeable delay.

*Acceptance Criteria:*
- Preview rendering completes in <100ms
- No blocking JavaScript during preview render
- CSS filters applied via hardware acceleration
- Reference images cached after first load
- Smooth slider interaction (60fps minimum)
- No performance degradation with multiple preview sessions
- Works on mid-range mobile devices (tested on 3-year-old phones)

---

## Non-Functional Requirements

### Performance

**NFR-P1: Conversion Performance**
- All format conversions complete in <100ms per file (WASM target)
- Batch conversions process files concurrently where possible
- No UI blocking during conversion (use Web Workers if needed)
- Memory efficient: Handle 100+ file batch without memory issues
- Benchmark regression tests prevent performance degradation

**NFR-P2: Web Interface Load Time**
- Initial page load completes in <2 seconds on 3G connection
- WASM binary cached aggressively (long cache headers)
- Critical CSS inlined for first paint
- No render-blocking resources
- Lighthouse performance score ≥90

**NFR-P3: Preview Rendering Performance**
- CSS filter preview renders in <100ms
- Slider drag maintains 60fps minimum
- No jank during tab switching between reference images
- Preview modal opens in <200ms
- Works smoothly on mid-range mobile devices (3-year-old phones)

### Security

**NFR-S1: Privacy-First Architecture**
- All processing happens client-side (browser or local CLI/TUI)
- Zero data transmission to external servers
- No analytics, tracking, or telemetry
- No external resource loading (fonts, scripts, images hosted locally)
- Documentation explicitly states privacy guarantees

**NFR-S2: File Handling Security**
- Validate file types before processing (magic number checks, not just extension)
- Sanitize XML input to prevent XXE (XML External Entity) attacks
- Limit file size parsing to prevent DoS (e.g., 10MB max per file)
- Graceful error handling for malformed files (no crashes, no sensitive error details)
- ZIP bomb protection for .costylepack bundles (limit uncompressed size ratio)

**NFR-S3: Web Application Security**
- Content Security Policy (CSP) headers prevent XSS
- No eval() or Function() constructors in JavaScript
- WASM binary integrity verification (subresource integrity if served)
- HTTPS-only deployment (enforced by Cloudflare Pages)
- No user-generated content or uploads to server

### Accessibility

**NFR-A1: WCAG 2.1 Level AA Compliance**
- Color contrast ratio ≥4.5:1 for normal text, ≥3:1 for large text
- All interactive elements keyboard accessible (Tab, Enter, Space, Esc)
- Focus indicators clearly visible (not browser default outline)
- Skip to main content link for keyboard users
- No flashing content (seizure prevention)

**NFR-A2: Screen Reader Support**
- All images have descriptive alt text
- All icons have aria-labels
- Form labels properly associated with inputs
- Error messages announced to screen readers (aria-live regions)
- Dynamic content updates communicated via ARIA
- Tested with NVDA (Windows) and VoiceOver (macOS/iOS)

**NFR-A3: Motor Impairment Accommodation**
- Touch targets ≥44x44px on mobile (Apple/Google guidelines)
- No hover-only interactions (mobile friendly)
- Generous click areas, no precision required
- Slider accessible via keyboard (arrow keys)
- No time-based interactions (auto-advancing carousels, etc.)

### Code Quality & Maintainability

**NFR-M1: Test Coverage**
- Unit test coverage ≥85% across internal packages
- Round-trip conversion tests for all format combinations
- Integration tests for CLI, TUI, Web interfaces
- Automated regression tests prevent breaking changes
- CI/CD pipeline blocks merges if tests fail

**NFR-M2: Architectural Consistency**
- Hub-and-spoke pattern maintained for all formats
- New format packages follow established structure:
  - `internal/formats/{format}/parse.go`
  - `internal/formats/{format}/generate.go`
  - `internal/formats/{format}/types.go`
  - `internal/formats/{format}/testdata/`
- UniversalRecipe remains single source of truth
- No circular dependencies between packages
- Clean separation: formats, converter, interfaces (CLI/TUI/Web)

**NFR-M3: Documentation Standards**
- All public functions have GoDoc comments
- Format specifications documented (NP3, XMP, Capture One, DCP)
- Parameter mapping tables maintained (docs/parameter-mapping.md)
- User-facing documentation updated (README, web landing page)
- Technical documentation for contributors (CLAUDE.md, architecture.md)
- CHANGELOG.md updated for all notable changes

**NFR-M4: Code Readability**
- Go code follows standard formatting (gofmt, golangci-lint)
- JavaScript follows consistent style (Prettier or similar)
- Clear variable and function names (no cryptic abbreviations)
- Comments explain "why" not "what" (code should be self-documenting)
- Complex algorithms include ASCII diagrams or examples
- No "clever" code; prefer clarity over brevity

### Browser & Platform Compatibility

**NFR-C1: Browser Support**
- **Tier 1 (Full support)**: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- **Tier 2 (Graceful degradation)**: Older versions show fallback message
- WebAssembly required: Check for WASM support, show error if unavailable
- CSS Grid and Flexbox used (no IE11 support needed)
- ES6+ JavaScript (no transpilation to ES5)

**NFR-C2: Operating System Support**
- **CLI/TUI Binaries**: Windows 10+, macOS 11+, Linux (Ubuntu 20.04+, Fedora 35+)
- **Web Interface**: Platform-agnostic (runs in browser)
- Binary releases for: windows/amd64, darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
- No platform-specific dependencies in core conversion logic

**NFR-C3: Mobile Device Support**
- **iOS**: Safari 14+ on iPhone 8 and newer
- **Android**: Chrome 90+ on Android 9+ devices
- Touch interactions work correctly (no desktop-only assumptions)
- Responsive design tested on real devices (not just browser DevTools)
- Performance acceptable on mid-range devices (not flagship-only)

### Deployment & Operations

**NFR-D1: Deployment Process**
- Web interface deploys to Cloudflare Pages automatically (GitHub integration)
- CLI/TUI binaries released via GitHub Releases with automated builds
- Versioning follows semantic versioning (SemVer 2.0.0)
- No downtime deployments for web interface
- Rollback capability if deployment issues detected

**NFR-D2: Monitoring & Error Handling**
- No crash logs collected (privacy-first, no telemetry)
- User-facing errors are clear and actionable ("File format not recognized" not "Error 0x04F2")
- Internal errors logged to browser console for debugging
- Documentation includes troubleshooting guide
- GitHub Issues for user-reported bugs and feature requests

**NFR-D3: Backwards Compatibility**
- Conversion engine maintains backward compatibility (existing conversions continue working)
- UniversalRecipe format versioned; parsers handle legacy versions
- Breaking changes documented in CHANGELOG, migration guide provided
- Deprecated features marked clearly, removed only in major versions

---

## Implementation Planning

### Epic Breakdown Required

Requirements must be decomposed into epics and bite-sized stories (200k context limit).

Path A consists of **4 major epics** with the following implementation sequence recommended:

**Recommended Implementation Order:**

1. **Epic 1: Capture One Format Support** (Weeks 1-3)
   - Foundation: Parse/generate .costyle XML files
   - Extends proven pattern from XMP (also XML-based)
   - Deliverable: Full Capture One conversion working in CLI/TUI/Web
   - Milestone: Professional photographers can convert NP3 → Capture One

2. **Epic 2: DCP Camera Profile Support** (Weeks 4-6)
   - Build on XML parsing experience from Epic 1
   - More complex: TIFF embedding, color matrices
   - Deliverable: DCP generation and parsing functional
   - Milestone: Camera profiles can be created from presets

3. **Epic 3: Enhanced Web UI/UX** (Weeks 7-9)
   - Redesigned landing page, format badges, batch processing
   - Can proceed in parallel with Epic 4
   - Deliverable: Professional, polished web interface
   - Milestone: First-time users impressed by quality and clarity

4. **Epic 4: Image Preview System** (Weeks 10-11)
   - Phase 1: CSS filter-based preview
   - Requires Web UI foundation from Epic 3
   - Deliverable: Approximate preset preview working
   - Milestone: Users can see preset effects before converting

**Total Estimated Timeline:** 11 weeks (flexible, personal project)

**Dependencies:**
- Epic 4 depends on Epic 3 (needs UI components like modal, slider)
- Epic 2 can start after Epic 1 (leverages XML parsing patterns)
- Epics 3 and 4 can overlap partially (UI components first, then preview integration)

**Quality Gates:**
- All epics must pass test coverage gate (≥85%)
- All epics must pass performance benchmarks (<100ms conversions)
- Web UI must pass browser compatibility testing (Chrome, Firefox, Safari, Edge)
- Documentation must be updated for each epic before marking complete

### Story Creation Process

After PRD approval, epics will be broken down into bite-sized user stories following BMad Method workflow:

1. **Epic Technical Specification** - Create detailed tech specs per epic
2. **Story Generation** - Break each epic into stories targeting 200k context limit
3. **Story Context Assembly** - Generate Story Context XML for each story
4. **Implementation** - Developer agent executes stories with DoD validation

**Next Step:** Run BMad workflow to create Architecture document, then Epic breakdown, then Stories.

---

## References

### Core Documentation

- **Product Brief**: `docs/product-brief-recipe-2025-11-08.md`
- **Original PRD**: `docs/PRD.md` (Recipe core functionality, dated 2025-11-03)
- **Architecture**: `docs/architecture.md` (Hub-and-spoke pattern, UniversalRecipe design)
- **Development Guide**: `CLAUDE.md` (Technical onboarding for contributors)

### Format Specifications

- **NP3 Format**: `docs/np3-format-specification.md` (Nikon Picture Control binary format)
- **Parameter Mapping**: `docs/parameter-mapping.md` (Cross-format parameter equivalence tables)
- **Adobe DNG Specification 1.6**: External reference for DCP format
- **Capture One .costyle**: Community-documented XML format (research required)

### Research & Planning Documents

- **Agent Roundtable Discussion**: Path A exploration (2025-11-08)
  - Market research: Capture One 8% market share, Lightroom 46.8%
  - CSS filter feasibility research
  - Competitive landscape analysis (Picture Instruments $39 converter)
  - Privacy-first architecture validation

### External References

- **W3C CSS Filter Effects Module Level 1**: CSS filter specification for preview feature
- **Adobe Camera Raw**: DCP compatibility testing target
- **Capture One Software**: Trial version for validation testing
- **Etsy Preset Marketplaces**: Source for .costyle sample files

### Technical Standards

- **SemVer 2.0.0**: Versioning standard for Recipe releases
- **WCAG 2.1 Level AA**: Accessibility compliance target
- **CommonMark**: Markdown standard for documentation

---

## Next Steps

### Immediate Next Steps

1. **Architecture Document** - Run: `/bmad:bmm:workflows:architecture`
   - Define format package structure for Capture One and DCP
   - Design CSS filter preview architecture
   - Plan web UI component structure
   - Document integration points between epics

2. **Epic Technical Specifications** - Run: `/bmad:bmm:workflows:epic-tech-context`
   - Create detailed tech spec for each of the 4 epics
   - Define acceptance criteria mappings
   - Identify integration requirements
   - Plan test strategies per epic

3. **UX Design** (if needed) - Run: `/bmad:bmm:workflows:create-ux-design`
   - Create visual mockups for enhanced web UI
   - Design format badge system
   - Prototype preview modal and slider
   - Test responsive layouts

4. **Story Breakdown** - Run: `/bmad:bmm:workflows:create-story`
   - Generate bite-sized stories from epics
   - Create Story Context XML for each story
   - Assign to sprint backlog
   - Begin implementation with developer agent

### Success Criteria Reminder

Path A is complete when:
- ✅ All four epics functionally complete
- ✅ Quality gates pass (coverage ≥85%, performance <100ms, browser compatibility)
- ✅ Documentation comprehensively updated
- ✅ Deployed to production (web) and released (CLI/TUI binaries via GitHub Releases)
- ✅ At least one Capture One user successfully converts presets
- ✅ Preview feature genuinely helps decision-making (not ignored)
- ✅ Web interface feels professional and polished

---

_This PRD captures the essence of Recipe Path A: transforming a solid conversion utility into a polished, professional-grade tool that photographers love to use._

_Path A brings visual discovery through preview, professional polish through enhanced UI, expanded reach through Capture One support, and maintains Recipe's core principles: privacy-first, local-only processing, zero dependencies._

_Created through collaborative discovery between Justin and AI facilitator using BMad Method workflows._
