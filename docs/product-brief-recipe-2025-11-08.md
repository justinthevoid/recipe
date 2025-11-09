# Product Brief: recipe

**Date:** 2025-11-08
**Author:** Justin
**Context:** Personal Project

---

## Executive Summary

Recipe is a universal photo preset converter that enables photographers to convert between Nikon NP3, Adobe Lightroom XMP, and lrtemplate formats. Path A represents a strategic enhancement focused on expanding format support and significantly improving user experience.

**The Opportunity:** Recipe currently serves Nikon photographers who want to use Picture Control recipes across different editing platforms. Path A extends this value by adding Capture One support (serving the professional 8% market segment), DCP camera profile capabilities, and a dramatically improved web interface with visual preset preview.

**Core Enhancements:**
1. Capture One .costyle format support
2. DCP (DNG Camera Profile) generation and parsing
3. Redesigned web UI with batch processing and mobile responsiveness
4. Image preview system showing preset effects before conversion

**Success Criteria:** Path A prioritizes craft and quality over growth metrics. Success means the preview feature is genuinely useful (not a gimmick), the codebase remains clean and maintainable, and the web interface feels professional and polished.

**Timeline:** Flexible personal project timeline, estimated 10-13 weeks for complete implementation across four major epics.

This brief captures the vision for Recipe's Path A enhancement, maintaining its core values (privacy-first, open-source, high-performance) while expanding capabilities to serve Nikon photographers more comprehensively.

---

## Core Vision

### Problem Statement

Recipe successfully converts between NP3 (Nikon Picture Control), XMP (Adobe Lightroom), and lrtemplate formats with high fidelity. However, the photography editing ecosystem includes additional major platforms that photographers use, particularly Capture One (used by 8% of photographers, primarily professionals).

Additionally, while Recipe's core conversion functionality is solid, the user experience - especially in the web interface - could be significantly enhanced. Photographers would benefit from being able to preview what a preset actually does before converting it, rather than a blind conversion process.

The current limitation to three formats and basic UI represents missed opportunities to serve photographers more comprehensively and create a more delightful, professional-grade tool.

{{#if problem_impact}}

### Problem Impact

{{problem_impact}}
{{/if}}

{{#if existing_solutions_gaps}}

### Why Existing Solutions Fall Short

{{existing_solutions_gaps}}
{{/if}}

### Proposed Solution

Path A focuses on maximizing value for current and potential Recipe users through four strategic enhancements:

1. **Capture One Format Support (.costyle)** - Expand Recipe's format coverage to include Capture One styles, enabling conversions for professional photographers who use this industry-standard software

2. **DCP Camera Profile Support** - Add support for DNG Camera Profile (DCP) files, allowing photographers to work with camera-specific color calibration profiles alongside their creative presets

3. **Enhanced Web Interface** - Redesign the web UI with modern UX patterns including visual format badges, batch processing with progress indicators, and mobile-responsive design

4. **Image Preview System** - Implement a visual preview feature that shows users what a preset will do to an image before conversion, using CSS filters initially and progressing to WebAssembly-based accurate rendering

These enhancements maintain Recipe's core principles: privacy-first (all processing client-side), high-performance (<100ms conversions), and zero dependencies on external services.

### Key Differentiators

- **Only universal converter supporting NP3 format** - Recipe remains the only tool (free or paid) that can convert Nikon's proprietary NP3 format
- **Complete privacy** - All processing happens locally in the browser or CLI; zero server uploads, zero analytics
- **True cross-platform** - CLI, TUI, and Web interfaces all using the same conversion engine
- **Open source and free** - No paywalls, subscriptions, or feature gating
- **Professional-grade accuracy** - 98%+ conversion fidelity with comprehensive test suite (1,531 real sample files)

---

## Target Users

### Primary Users

**Nikon Camera Owners Who Want Cross-Platform Flexibility**

Recipe's primary users are Nikon photographers who appreciate Nikon's Picture Control recipes and want to use them across different photo editing platforms. These users:

- Own Nikon cameras that support NP3 Picture Control files
- Have invested in Nikon's creative picture styles or downloaded community-created NP3 recipes
- Want to maintain consistent editing aesthetics whether they're using Nikon's software, Adobe Lightroom, or Capture One
- May be transitioning between editing platforms but don't want to lose their preset library
- Value the quality and characteristics of Nikon's color science and want to apply it in their preferred editing environment

These photographers range from enthusiasts to professionals, but share a common thread: they've discovered the power of Nikon's Picture Control system and don't want to be locked into a single editing ecosystem.

{{#if secondary_user_segment}}

### Secondary Users

{{secondary_user_segment}}
{{/if}}

{{#if user_journey}}

### User Journey

{{user_journey}}
{{/if}}

---

## Success Metrics

### Project Success Criteria

Path A will be considered successful when:

1. **Preview Feature Utility** - The image preview feature provides genuine value to users, not just visual decoration. Users should be able to make informed decisions about conversions based on the preview, with sufficient accuracy that the preview reflects the actual converted output.

2. **Code Quality Maintained** - The codebase remains clean, maintainable, and consistent with existing patterns despite adding significant new functionality. New format packages follow the established hub-and-spoke architecture, tests maintain 85%+ coverage, and performance stays within existing benchmarks (<100ms conversions).

3. **Professional Web Interface** - The web UI feels polished and professional, with modern UX patterns that inspire confidence. The interface should be visually appealing, intuitive to use, and responsive across devices (desktop, tablet, mobile).

These criteria prioritize craft and user value over growth metrics, reflecting Recipe's nature as a personal project focused on excellence.

---

## MVP Scope

### Core Features

Path A consists of four major feature areas, all essential to the enhancement vision:

**1. Capture One Format Support**
- Parse .costyle files (XML-based Capture One preset format)
- Generate .costyle files from UniversalRecipe
- Support .costylepack bundles (zip archives of multiple styles)
- Round-trip conversion testing with real Capture One files
- Integration with CLI, TUI, and Web interfaces

**2. DCP Camera Profile Support**
- Parse DCP (DNG Camera Profile) files per Adobe DNG spec 1.6
- Generate DCP files from UniversalRecipe color adjustments
- Support embedded profiles in DNG files
- Validate compatibility with Adobe Camera Raw and Lightroom

**3. Enhanced Web UI/UX**
- Modern landing page with visual format badges
- Batch file upload with drag-and-drop support
- Progress indicators for multi-file conversions
- Mobile-responsive design (works on phones/tablets)
- Improved conversion flow with clear format selection
- Before/after comparison slider

**4. Image Preview System**
- **Phase 1 (MVP)**: CSS filter-based preview using browser native capabilities
  - Map UniversalRecipe parameters to CSS filter values (brightness, contrast, saturation, hue-rotate)
  - Instant preview on reference images
  - Approximate accuracy (good enough for decision-making)
- **Phase 2 (Future)**: WebAssembly-based accurate preview
  - Integration with Photon or custom WASM image processing
  - Pixel-perfect preview matching actual conversion output
  - Support for complex adjustments (tone curves, HSL, etc.)

### Out of Scope for MVP

The following capabilities were discussed during planning but are explicitly deferred to future phases:

**Deferred to Path B (Market Expansion):**
- 3D LUT export (.cube format for video color grading)
- API/REST endpoints for third-party integrations
- Browser extensions or plugins for preset marketplaces
- Mobile native apps (iOS/Android)

**Deferred to Path C (Platform Features):**
- Preset library management and organization
- Cloud sync/backup (even with user's own storage)
- Collaborative features or preset sharing
- Preset marketplace integration

**Technical Scope Constraints:**
- Image preview Phase 2 (WASM-based) is aspirational; CSS filters (Phase 1) are the MVP commitment
- DCP support focuses on generation and parsing; full color science validation is best-effort
- Capture One testing relies on trial version; edge cases may emerge post-launch

### MVP Success Criteria

Path A is complete when all four core features are functional and meet quality standards:

**Functional Completeness:**
- ✅ Capture One .costyle files convert successfully to/from NP3, XMP, lrtemplate
- ✅ DCP files can be generated from UniversalRecipe and parsed back
- ✅ Web UI supports batch uploads with progress tracking
- ✅ CSS filter preview displays approximate preset effects on reference images

**Quality Gates:**
- ✅ Round-trip tests pass for all new format combinations
- ✅ Code coverage remains ≥85% across internal packages
- ✅ Performance: All conversions complete in <100ms (WASM target)
- ✅ Web UI is mobile-responsive and works on Chrome, Firefox, Safari, Edge
- ✅ Documentation updated (CLAUDE.md, format guides, user tutorials)

**User Validation:**
- At least one real Capture One user successfully uses the conversion
- Preview feature helps make conversion decisions (not ignored as useless)
- No major architectural refactoring needed to maintain code quality

### Future Vision

**Beyond Path A:**

While Path A focuses on format expansion and UX polish, the longer-term vision includes:

- **Path B**: Market expansion through LUT export, API integrations, and video color grading support
- **Path C**: Platform features like preset library management, search, and organization tools
- **Advanced Preview**: Full WASM-based image processing for pixel-perfect preset visualization
- **Additional Formats**: DxO PhotoLab, ON1, Affinity Photo as user demand warrants
- **Community Features**: Preset sharing, collaborative editing, marketplace integration (if Recipe grows beyond personal project scope)

---

## Market Context

**Photography Editing Software Landscape (2025):**

- **Adobe Lightroom**: 46.8% market share - dominant platform, largest preset ecosystem
- **Capture One**: 8% market share - preferred by professional photographers for advanced color grading and studio workflows
- **Other platforms**: DxO PhotoLab, ON1, Affinity Photo, Darktable (smaller market share)

**Preset Converter Market:**
- Picture Instruments Preset Converter: $39 (Lightroom → Capture One only)
- Various free tools for specific conversion pairs (limited format support)
- **Recipe is the ONLY tool supporting NP3 (Nikon Picture Control) format**

**Market Opportunity:**
- Nikon camera owners represent a significant photography segment
- Professional photographers using Capture One (8%) are underserved for NP3 conversion
- Photographers switching platforms lose preset investments (economic pain point)
- Privacy-conscious users value local processing (zero cloud uploads)

**Competitive Position:**
Recipe occupies a unique niche as the only free, open-source, privacy-first converter supporting Nikon's proprietary NP3 format. Path A strengthens this position by adding professional-grade formats (Capture One, DCP) while maintaining Recipe's core differentiators.

{{#if financial_considerations}}

## Financial Considerations

{{financial_considerations}}
{{/if}}

## Technical Preferences

**Technology Stack (Existing):**
- **Language**: Go 1.25.1+ (leverages `go:wasmexport` for WASM)
- **CLI Framework**: Cobra
- **TUI Framework**: Bubbletea v2
- **Web**: Vanilla JavaScript (ES6+) + WebAssembly
- **Deployment**: Cloudflare Pages (web), GitHub Releases (binaries)

**Path A Technical Decisions:**

**Format Implementation:**
- Follow existing pattern: `internal/formats/{format}/` with `parse.go` and `generate.go`
- Capture One: XML parsing (similar to XMP implementation)
- DCP: XML embedded in TIFF (use existing XML libraries + TIFF reading)

**Preview System:**
- **Phase 1**: Pure CSS filters (zero dependencies, instant)
  - Map UniversalRecipe → CSS filter functions
  - Trade accuracy for speed and simplicity
- **Phase 2**: Consider Photon (Rust/WASM) or custom image processing
  - Only if Phase 1 proves insufficient
  - Evaluate bundle size impact (~500KB)

**Web UI:**
- Continue with vanilla JavaScript (no React/Vue/framework bloat)
- Modern CSS (Grid, Flexbox, CSS Variables)
- Progressive enhancement (works without JavaScript for basic conversion)

**Testing:**
- Acquire real .costyle files from Etsy/marketplaces for test fixtures
- Download reference DCP files from Adobe
- Maintain >85% test coverage
- Validate with actual Capture One software (trial version)

{{#if organizational_context}}

## Organizational Context

{{organizational_context}}
{{/if}}

## Risks and Assumptions

**Technical Risks:**

1. **Format Specification Accuracy**
   - **Risk**: Capture One .costyle and DCP formats may have undocumented quirks or version differences
   - **Mitigation**: Test with multiple real-world files; validate output in actual software

2. **Preview Accuracy**
   - **Risk**: CSS filter approximation may be too inaccurate to be useful
   - **Mitigation**: Start with Phase 1 (CSS); clearly communicate "approximate preview"; can upgrade to WASM if needed

3. **Test File Acquisition**
   - **Risk**: May not have access to sufficient .costyle or DCP sample files
   - **Mitigation**: Purchase budget samples from Etsy; download from Adobe; community contributions

4. **Browser Compatibility**
   - **Risk**: CSS filters or WASM features may not work consistently across browsers
   - **Mitigation**: Test on Chrome, Firefox, Safari, Edge; graceful degradation

**Assumptions:**

- Capture One .costyle format is stable and well-documented (XML-based)
- Adobe DNG spec 1.6 provides sufficient DCP implementation guidance
- CSS filter approximation is "good enough" for MVP preview feature
- Vanilla JavaScript approach remains viable (no framework needed yet)
- Personal project timeline is flexible (no hard deadlines)
- Test coverage can be maintained while adding new features

**Dependencies:**

- Access to Capture One trial version for validation testing
- Availability of .costyle and DCP sample files for test suite
- Cloudflare Pages continues to support free WASM hosting

{{#if timeline_constraints}}

## Timeline

{{timeline_constraints}}
{{/if}}

## Supporting Materials

**Research Conducted:**
- Agent roundtable discussion (2025-11-08) - comprehensive Path A exploration
- Market research on Capture One vs Lightroom adoption (46.8% vs 8%)
- Technical research on CSS filters and WebAssembly image processing
- Competitive analysis of existing preset converters

**Existing Documentation:**
- `docs/architecture.md` - Current Recipe architecture (hub-and-spoke pattern)
- `docs/np3-format-specification.md` - NP3 binary format details
- `docs/parameter-mapping.md` - Cross-format parameter mapping
- `CLAUDE.md` - Development guide for Recipe codebase

**Referenced Specifications:**
- Adobe DNG Specification 1.6 (DCP format documentation)
- Capture One .costyle format (XML-based, community documented)
- CSS Filter Effects Module Level 1 (W3C specification)

---

_This Product Brief captures the vision and requirements for recipe._

_It was created through collaborative discovery and reflects the unique needs of this Personal Project project._

_Next: PRD workflow will transform this brief into detailed planning artifacts including requirements, architecture, epics, and user stories._
