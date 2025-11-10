# Story 7.3: FAQ Documentation

**Epic:** Epic 7 - Documentation & Deployment (FR-7)
**Story ID:** 7.3
**Status:** ready-for-dev
**Created:** 2025-11-06
**Complexity:** Low (1-2 days)

---

## Story

As a **photographer with questions about Recipe's legality, privacy, and conversion quality**,
I want **a comprehensive FAQ section answering common questions**,
so that **I can make informed decisions about using Recipe and understand its limitations before converting my presets**.

---

## Business Value

The FAQ is Recipe's **trust-building and educational layer**, addressing user concerns proactively and reducing support burden through self-service answers.

**Strategic Value:**
- **Legal Transparency:** Honest discussion of reverse engineering and fair use (builds trust)
- **Privacy Reassurance:** Reinforces WASM client-side architecture (core value proposition)
- **Set Expectations:** Explains conversion limitations and accuracy targets (reduces disappointment)
- **Reduce Support Burden:** Common questions answered upfront (fewer GitHub Issues)

**User Impact:**
- Users understand legal landscape before using Recipe (informed consent)
- Users confident their preset files remain private (privacy promise validated)
- Users know what to expect from conversion accuracy (95%+ for core parameters)
- Users understand why some features don't convert (format limitations, not bugs)

**Competitive Differentiation:**
- Most converters hide legal/technical details - Recipe documents them openly
- Transparent about limitations demonstrates honesty and technical credibility
- Educational content helps photographers understand color science and file formats

---

## Acceptance Criteria

### AC-1: FAQ Answers "Is this legal? (reverse engineering)"

**Given** a user concerned about legal risks of using Recipe
**When** they read the FAQ
**Then**:
- ✅ **Question Present:**
  - "Is Recipe legal? Is reverse engineering allowed?"
  - OR: "Can I use Recipe without legal risk?"
- ✅ **Answer Includes:**
  - Reverse engineering context explained (Nikon .np3 format is proprietary/undocumented)
  - Fair use and interoperability exemptions cited (if applicable to jurisdiction)
  - Research and educational purposes mentioned (DMCA 1201 exemptions)
  - No affiliation with Nikon/Adobe stated clearly
  - Recommendation: **Private use only** until legal assessment complete
- ✅ **Honest Tone:**
  - Not legal advice disclaimer
  - Conservative approach (use at own risk)
  - Suggests consulting IP attorney for commercial use
- ✅ **Links to Resources:**
  - DMCA 1201 exemptions (optional)
  - Fair use doctrine (optional)
  - Legal disclaimer section (Story 7-4)

**Example Answer:**
```markdown
### Is Recipe legal? Is reverse engineering allowed?

**Short Answer:** Reverse engineering for interoperability is generally protected,
but Recipe is provided "as is" for **private, non-commercial use** without legal guarantees.

**Details:**

Recipe uses reverse-engineered file formats (Nikon .np3) to enable conversion
between photo preset formats. Reverse engineering is generally permitted under:

- **Fair Use Doctrine:** Interoperability and research purposes (US copyright law)
- **DMCA 1201 Exemptions:** Reverse engineering for compatibility (renewed 2021)
- **EU Directive 2009/24/EC:** Interoperability of computer programs

**Important Caveats:**
- Recipe has **no affiliation** with Nikon, Adobe, or other vendors
- Legal landscape varies by jurisdiction
- Commercial use may require additional legal assessment

**Recommendation:**
- ✅ **Private/Personal Use:** Generally safe for individual photographers
- ⚠️ **Commercial Use:** Consult an IP attorney before distribution/monetization
- ⚠️ **Purchased Presets:** Respect creator copyrights (presets themselves may be copyrighted)

**Disclaimer:** This is NOT legal advice. Use Recipe at your own risk. See full
[Legal Disclaimer](../index.html#legal-disclaimer) for terms.

**References:**
- DMCA 1201 Exemptions: https://www.copyright.gov/1201/
- Fair Use: https://www.copyright.gov/fair-use/
```

**Validation:**
- Question clearly addresses legal concerns
- Answer is honest and conservative (not overpromising)
- Recommends private use, cautions commercial use
- Links to legal disclaimer

---

### AC-2: FAQ Answers "Is my data private?"

**Given** a user concerned about uploading proprietary/purchased presets
**When** they read the FAQ
**Then**:
- ✅ **Question Present:**
  - "Is my data private? Do files get uploaded to a server?"
  - OR: "Are my presets safe? Can anyone see my files?"
- ✅ **Answer Includes:**
  - **Confirms Yes:** Files remain 100% private
  - **WebAssembly Architecture:** All processing happens in browser (client-side)
  - **Zero Server Uploads:** No network requests during conversion
  - **No Tracking:** No analytics, no cookies, no data collection (references FR-2.9)
  - **Proof:** Network monitoring shows zero uploads (user can verify with browser DevTools)
- ✅ **Technical Explanation (Simple Language):**
  - Explains WASM = code runs on user's computer, not server
  - Files held in memory only during conversion
  - No local storage (except WASM cache for performance)
- ✅ **Links to:**
  - Privacy messaging section on landing page (Story 7-1, FR-2.8)
  - Architecture documentation (optional, for technical users)

**Example Answer:**
```markdown
### Is my data private? Do files get uploaded to a server?

**Short Answer:** YES - your files are 100% private. Recipe uses WebAssembly to
process files **entirely in your browser**. Nothing is uploaded to any server.

**How It Works:**

Recipe runs as **WebAssembly** (WASM) code directly in your web browser. This means:

- ✅ **No Server Uploads:** Your preset files never leave your device
- ✅ **Client-Side Processing:** Conversion happens on your computer, not in the cloud
- ✅ **No Tracking:** Zero analytics, cookies, or data collection
- ✅ **Temporary Memory Only:** Files held in RAM during conversion, cleared immediately after

**Why This Matters:**

If you're converting purchased presets or proprietary custom presets, you can
trust Recipe to keep them private. No one (not even Recipe's developer) can see
your files.

**Verify Yourself:**

Open your browser's Developer Tools (F12) and watch the Network tab during
conversion. You'll see **zero network requests** - proof that nothing is uploaded.

**Technical Details:**

Recipe's privacy-first architecture is powered by WebAssembly - a technology
that compiles code to run natively in your browser. It's the same technology
used by Google Docs, Figma, and other privacy-conscious web apps.

**No Exceptions:**

- ❌ No server-side processing
- ❌ No cloud storage
- ❌ No third-party analytics (Google Analytics, etc.)
- ❌ No cookies or tracking pixels

**References:**
- [Privacy Promise](../index.html#privacy) - Landing page privacy messaging
- [Architecture](../architecture.md#webassembly-deployment) - Technical details (for developers)
- [Story 2-9: Privacy Messaging](2-9-privacy-messaging.md) - Implementation details
```

**Validation:**
- Question directly addresses privacy concerns
- Answer confirms privacy clearly (no ambiguity)
- WASM architecture explained in simple terms
- User can verify privacy claim (network monitoring)
- Links to privacy messaging and architecture

---

### AC-3: FAQ Answers "Why doesn't [feature] convert?"

**Given** a user confused why Tone Curves or Grain don't convert to NP3
**When** they read the FAQ
**Then**:
- ✅ **Question Present:**
  - "Why don't some parameters convert? (Tone Curves, Grain, Vignette)"
  - OR: "Why is my converted preset missing some adjustments?"
- ✅ **Answer Includes:**
  - **Format Limitations:** NP3 is simpler than XMP/lrtemplate (in-camera format)
  - **Examples of Unmappable Features:**
    - Tone Curves (NP3 has no curve support)
    - Grain (texture effect not in NP3)
    - Vignette (edge darkening not in NP3)
    - Split Toning (shadow/highlight tints not in NP3)
  - **What Does Convert:**
    - Core adjustments: Exposure, Contrast, Saturation, HSL (95%+ accuracy)
    - See compatibility matrix for complete list
  - **Recipe's Behavior:**
    - Warns user when parameters can't be mapped
    - Skips unmappable features (no silent data loss)
- ✅ **Links to:**
  - Format Compatibility Matrix (Story 7-2) - Complete parameter list
  - Known Limitations section (in matrix)

**Example Answer:**
```markdown
### Why don't some parameters convert? (Tone Curves, Grain, Vignette)

**Short Answer:** NP3 (Nikon Picture Control) is a **simpler format** than
XMP/lrtemplate. Some advanced Lightroom features don't exist in NP3.

**Background:**

NP3 was designed for **in-camera use** with limited controls (Nikon cameras
have small screens and limited processing power). XMP/lrtemplate were designed
for **desktop editing** with 50+ adjustment parameters.

**What NP3 Does NOT Support:**

| Feature       | Why Not Supported                           |
| ------------- | ------------------------------------------- |
| Tone Curves   | NP3 has no tone curve support (linear only) |
| Grain/Texture | Texture effects not in NP3                  |
| Vignette      | Edge vignetting not in NP3                  |
| Split Toning  | Shadow/highlight color tints not in NP3     |
| Clarity       | Mid-tone contrast not in NP3                |
| Dehaze        | Atmospheric haze removal not in NP3         |

**What DOES Convert (95%+ Accuracy):**

- ✅ **Core Adjustments:** Exposure, Contrast, Brightness, Saturation
- ✅ **HSL Adjustments:** Hue, Saturation, Luminance for 8 colors (24 parameters)
- ✅ **Color:** Temperature, Tint, White Balance
- ✅ **Detail:** Sharpness

**Recipe's Behavior:**

When you convert an XMP/lrtemplate preset with unsupported features, Recipe:

1. **Warns You:** Displays which parameters can't be mapped
2. **Converts What It Can:** Core adjustments convert accurately
3. **Skips Unmappable Features:** No silent data loss or approximation errors

**Example:**

```
Input: Vintage_Film.xmp
- Exposure: +0.5 ✅ Converts
- Contrast: +25  ✅ Converts
- Grain: 15     ❌ Skipped (not in NP3)
- Vignette: -20 ❌ Skipped (not in NP3)

Output: Vintage_Film.np3
- Brightness: +0.5 (Exposure)
- Contrast: +0.8 (scaled from +25)
⚠️ Warning: Grain and Vignette not supported in NP3
```

**Want Full Details?**

See the [Format Compatibility Matrix](../format-compatibility-matrix.md) for
a complete list of 30+ parameters and their conversion status.

**Summary:** Recipe focuses on accurately converting what NP3 **can** support
(95%+ accuracy) and being transparent about what it **cannot** support.
```

**Validation:**
- Question addresses common user confusion
- Answer explains format limitations clearly
- Examples of unmappable features provided
- Links to compatibility matrix for details
- User expectations set appropriately

---

### AC-4: FAQ Answers "How accurate is conversion?"

**Given** a user wanting to know conversion quality expectations
**When** they read the FAQ
**Then**:
- ✅ **Question Present:**
  - "How accurate is conversion? Will colors look the same?"
  - OR: "Can I trust Recipe to preserve my preset's look?"
- ✅ **Answer Includes:**
  - **95%+ Accuracy Target:** For core parameters (Exposure, Contrast, Saturation, HSL)
  - **Approximations Explained:** Some parameters don't have 1:1 equivalents (Highlights, Shadows)
  - **Validation Methods:**
    - Round-trip testing (A → B → A = identical)
    - Visual validation (side-by-side comparison)
    - 1,501 sample files tested
  - **What Affects Accuracy:**
    - Direct mappings: 98%+ accuracy (Exposure, Contrast)
    - Approximations: 90-95% accuracy (Highlights → Contrast)
    - Unmappable features: 0% (skipped with warning)
  - **Recommendation:**
    - Always visually validate converted presets
    - Test on sample image before production use
- ✅ **Links to:**
  - Format Compatibility Matrix (accuracy by parameter category)
  - Epic 1 Retrospective (validation results, if available)

**Example Answer:**
```markdown
### How accurate is conversion? Will colors look the same?

**Short Answer:** Recipe achieves **95%+ accuracy** for core adjustments
(Exposure, Contrast, Saturation, HSL). Visual similarity is very high, but
always validate converted presets before production use.

**Accuracy by Parameter Type:**

| Parameter Type     | Accuracy | Details                              |
| ------------------ | -------- | ------------------------------------ |
| Direct Mappings    | 98%+     | Exposure, Contrast, Saturation, HSL  |
| Approximations     | 90-95%   | Highlights, Shadows, Vibrance        |
| Unmappable         | N/A      | Skipped with warning (Grain, Curves) |
| **Overall (Core)** | **95%+** | **For parameters that convert**      |

**How Recipe Ensures Accuracy:**

1. **Round-Trip Testing:**
   - Convert A → B → A and verify parameter equality
   - Example: XMP → NP3 → XMP preserves Exposure, Contrast, Saturation

2. **Visual Validation:**
   - Apply preset to reference image in Lightroom
   - Apply converted preset to same image in camera/NX Studio
   - Compare visual output (color delta E <5 for critical colors)

3. **Real-World Sample Files:**
   - Tested against 1,501 sample files:
     - 22 NP3 files (Nikon official Picture Controls)
     - 913 XMP files (Lightroom CC presets)
     - 544 lrtemplate files (Lightroom Classic presets)
   - All conversions validated for accuracy and edge cases

**What Affects Accuracy:**

**Direct 1:1 Mappings (98%+ Accuracy):**
- Parameters with identical equivalents across formats
- Examples: Exposure, Contrast, Saturation, HSL adjustments
- Conversion preserves exact values (within rounding tolerance ±1)

**Approximations (90-95% Accuracy):**
- Parameters without 1:1 equivalents (Highlights, Shadows, Vibrance)
- Recipe uses intelligent mapping to preserve creative intent
- Example: Lightroom Highlights → NP3 Contrast adjustment
- Visual similarity remains high, but not pixel-perfect

**Format Limitations (Cannot Convert):**
- Advanced Lightroom features not in NP3 (Tone Curves, Grain, Vignette)
- Recipe skips these features and warns you
- Does not attempt to approximate (avoid unpredictable results)

**Best Practices:**

1. **Visual Validation:**
   - Always test converted preset on sample image before production
   - Compare side-by-side with original preset (camera vs. Lightroom)
   - Adjust if needed (fine-tune on camera)

2. **Understand Limitations:**
   - Review [Format Compatibility Matrix](../format-compatibility-matrix.md)
   - Know which parameters convert directly vs. approximated
   - Set expectations accordingly

3. **Expect 95%+ for Core Adjustments:**
   - Exposure, Contrast, Saturation, HSL → Very accurate
   - Highlights, Shadows, Vibrance → Good approximation
   - Tone Curves, Grain, Vignette → Not supported in NP3

**Summary:** Recipe delivers high accuracy (95%+) for core color adjustments,
with transparent warnings when parameters can't be mapped. Always validate
converted presets visually to ensure they meet your creative intent.

**References:**
- [Format Compatibility Matrix](../format-compatibility-matrix.md) - Accuracy by parameter
- [Epic 1 Retrospective](../epic-1-retrospective.md) - Validation results (if available)
```

**Validation:**
- Question addresses conversion quality concerns
- Answer provides concrete accuracy metrics (95%+)
- Validation methods explained (round-trip, visual, sample files)
- Best practices recommended (visual validation)
- Links to compatibility matrix and retrospective

---

### AC-5: Answers are Clear and Concise

**Given** a user scanning the FAQ for quick answers
**When** they read any FAQ answer
**Then**:
- ✅ **Short Answer First:**
  - Each question has a 1-2 sentence "Short Answer" at the top
  - User gets immediate answer without reading full details
- ✅ **Length Constraint:**
  - Full answer: 2-8 paragraphs (not too brief, not overwhelming)
  - Short Answer: 1-2 sentences
  - Details: 3-6 paragraphs with examples, tables, or lists
- ✅ **Clear Structure:**
  - Headings organize long answers (Background, Details, Best Practices, Summary)
  - Tables for scannable information (parameter accuracy, unmappable features)
  - Bullet lists for key points
- ✅ **Non-Technical Language:**
  - Avoids jargon without explanation
  - Technical terms defined when first used (WASM = WebAssembly)
  - Photographer-friendly language (not developer-speak)
- ✅ **Actionable:**
  - Provides next steps (links to details, recommendations)
  - User knows what to do after reading answer

**Example Structure:**
```markdown
### [Question]

**Short Answer:** [1-2 sentences]

**Details:**

[3-4 paragraphs or structured content]

**Best Practices / Recommendations:**

- [Actionable advice]

**References:**
- [Links to related docs]
```

**Validation:**
- All answers follow structure (Short Answer → Details → References)
- Length appropriate (not too brief, not overwhelming)
- Scannable (headings, lists, tables)
- Non-technical language
- User can find answer quickly

---

### AC-6: FAQ Updated Based on User Feedback (Post-Launch)

**Given** Recipe is deployed and users submit questions via GitHub Issues
**When** a question is asked repeatedly (3+ times)
**Then**:
- ✅ **Process Defined:**
  - Monitor GitHub Issues for common questions
  - Triage: If question asked 3+ times → Add to FAQ
  - Draft answer based on existing responses
  - Review for clarity and accuracy
  - Update FAQ document
  - Link to FAQ in future Issue responses
- ✅ **Documentation:**
  - Process documented in README or CONTRIBUTING.md
  - FAQ includes "Last Updated" date at top
  - Change log tracks added questions
- ✅ **Continuous Improvement:**
  - FAQ evolves based on real user needs
  - Reduces support burden over time

**Example Process:**
```markdown
## FAQ Maintenance Process

1. **Monitor GitHub Issues:**
   - Weekly review of new issues
   - Tag questions with `faq-candidate` label if asked 3+ times

2. **Add to FAQ:**
   - Draft answer based on existing Issue responses
   - Review for clarity and conciseness
   - Add to `docs/faq.md` in appropriate category

3. **Update Deployment:**
   - Commit FAQ changes to main branch
   - Cloudflare Pages auto-deploys updated FAQ
   - Link to FAQ answer in future Issue responses

4. **Track Changes:**
   - Update "Last Updated" date at top of FAQ
   - Add entry to FAQ change log (optional)
```

**Validation:**
- Process documented for post-launch FAQ updates
- FAQ includes "Last Updated" date
- Mechanism to track common questions (GitHub Issues labels)

---

## Tasks / Subtasks

### Task 1: Create FAQ Document (AC-1, AC-2, AC-3, AC-4, AC-5)

- [ ] **Create `docs/faq.md`:**
  ```markdown
  # Recipe FAQ

  **Last Updated:** 2025-11-06
  **Recipe Version:** v2.0.0

  Frequently asked questions about Recipe's legality, privacy, conversion quality,
  and technical limitations.

  ---

  ## Table of Contents

  - [Legal & Licensing](#legal--licensing)
    - [Is Recipe legal? Is reverse engineering allowed?](#is-recipe-legal-is-reverse-engineering-allowed)
  - [Privacy & Security](#privacy--security)
    - [Is my data private? Do files get uploaded to a server?](#is-my-data-private-do-files-get-uploaded-to-a-server)
  - [Conversion Quality](#conversion-quality)
    - [How accurate is conversion? Will colors look the same?](#how-accurate-is-conversion-will-colors-look-the-same)
    - [Why don't some parameters convert? (Tone Curves, Grain, Vignette)](#why-dont-some-parameters-convert-tone-curves-grain-vignette)
  - [Technical Questions](#technical-questions)
    - [What browsers are supported?](#what-browsers-are-supported)
    - [Can I use Recipe offline?](#can-i-use-recipe-offline)
    - [How do I report bugs or request features?](#how-do-i-report-bugs-or-request-features)

  ---

  ## Legal & Licensing

  ### Is Recipe legal? Is reverse engineering allowed?

  **Short Answer:** Reverse engineering for interoperability is generally protected,
  but Recipe is provided "as is" for **private, non-commercial use** without legal guarantees.

  **Details:**

  Recipe uses reverse-engineered file formats (Nikon .np3) to enable conversion
  between photo preset formats. Reverse engineering is generally permitted under:

  - **Fair Use Doctrine:** Interoperability and research purposes (US copyright law)
  - **DMCA 1201 Exemptions:** Reverse engineering for compatibility (renewed 2021)
  - **EU Directive 2009/24/EC:** Interoperability of computer programs

  **Important Caveats:**
  - Recipe has **no affiliation** with Nikon, Adobe, or other vendors
  - Legal landscape varies by jurisdiction
  - Commercial use may require additional legal assessment

  **Recommendation:**
  - ✅ **Private/Personal Use:** Generally safe for individual photographers
  - ⚠️ **Commercial Use:** Consult an IP attorney before distribution/monetization
  - ⚠️ **Purchased Presets:** Respect creator copyrights (presets themselves may be copyrighted)

  **Disclaimer:** This is NOT legal advice. Use Recipe at your own risk. See full
  [Legal Disclaimer](#) for terms.

  **References:**
  - DMCA 1201 Exemptions: https://www.copyright.gov/1201/
  - Fair Use: https://www.copyright.gov/fair-use/

  ---

  ## Privacy & Security

  ### Is my data private? Do files get uploaded to a server?

  **Short Answer:** YES - your files are 100% private. Recipe uses WebAssembly to
  process files **entirely in your browser**. Nothing is uploaded to any server.

  **How It Works:**

  Recipe runs as **WebAssembly** (WASM) code directly in your web browser. This means:

  - ✅ **No Server Uploads:** Your preset files never leave your device
  - ✅ **Client-Side Processing:** Conversion happens on your computer, not in the cloud
  - ✅ **No Tracking:** Zero analytics, cookies, or data collection
  - ✅ **Temporary Memory Only:** Files held in RAM during conversion, cleared immediately after

  **Why This Matters:**

  If you're converting purchased presets or proprietary custom presets, you can
  trust Recipe to keep them private. No one (not even Recipe's developer) can see
  your files.

  **Verify Yourself:**

  Open your browser's Developer Tools (F12) and watch the Network tab during
  conversion. You'll see **zero network requests** - proof that nothing is uploaded.

  **Technical Details:**

  Recipe's privacy-first architecture is powered by WebAssembly - a technology
  that compiles code to run natively in your browser. It's the same technology
  used by Google Docs, Figma, and other privacy-conscious web apps.

  **No Exceptions:**

  - ❌ No server-side processing
  - ❌ No cloud storage
  - ❌ No third-party analytics (Google Analytics, etc.)
  - ❌ No cookies or tracking pixels

  **References:**
  - [Privacy Promise](#) - Landing page privacy messaging
  - [Architecture](architecture.md#webassembly-deployment) - Technical details (for developers)
  - [Story 2-9: Privacy Messaging](stories/2-9-privacy-messaging.md) - Implementation details

  ---

  ## Conversion Quality

  ### How accurate is conversion? Will colors look the same?

  **Short Answer:** Recipe achieves **95%+ accuracy** for core adjustments
  (Exposure, Contrast, Saturation, HSL). Visual similarity is very high, but
  always validate converted presets before production use.

  **Accuracy by Parameter Type:**

  | Parameter Type     | Accuracy | Details                              |
  | ------------------ | -------- | ------------------------------------ |
  | Direct Mappings    | 98%+     | Exposure, Contrast, Saturation, HSL  |
  | Approximations     | 90-95%   | Highlights, Shadows, Vibrance        |
  | Unmappable         | N/A      | Skipped with warning (Grain, Curves) |
  | **Overall (Core)** | **95%+** | **For parameters that convert**      |

  **How Recipe Ensures Accuracy:**

  1. **Round-Trip Testing:**
     - Convert A → B → A and verify parameter equality
     - Example: XMP → NP3 → XMP preserves Exposure, Contrast, Saturation

  2. **Visual Validation:**
     - Apply preset to reference image in Lightroom
     - Apply converted preset to same image in camera/NX Studio
     - Compare visual output (color delta E <5 for critical colors)

  3. **Real-World Sample Files:**
     - Tested against 1,501 sample files:
       - 22 NP3 files (Nikon official Picture Controls)
       - 913 XMP files (Lightroom CC presets)
       - 544 lrtemplate files (Lightroom Classic presets)
     - All conversions validated for accuracy and edge cases

  **What Affects Accuracy:**

  **Direct 1:1 Mappings (98%+ Accuracy):**
  - Parameters with identical equivalents across formats
  - Examples: Exposure, Contrast, Saturation, HSL adjustments
  - Conversion preserves exact values (within rounding tolerance ±1)

  **Approximations (90-95% Accuracy):**
  - Parameters without 1:1 equivalents (Highlights, Shadows, Vibrance)
  - Recipe uses intelligent mapping to preserve creative intent
  - Example: Lightroom Highlights → NP3 Contrast adjustment
  - Visual similarity remains high, but not pixel-perfect

  **Format Limitations (Cannot Convert):**
  - Advanced Lightroom features not in NP3 (Tone Curves, Grain, Vignette)
  - Recipe skips these features and warns you
  - Does not attempt to approximate (avoid unpredictable results)

  **Best Practices:**

  1. **Visual Validation:**
     - Always test converted preset on sample image before production
     - Compare side-by-side with original preset (camera vs. Lightroom)
     - Adjust if needed (fine-tune on camera)

  2. **Understand Limitations:**
     - Review [Format Compatibility Matrix](format-compatibility-matrix.md)
     - Know which parameters convert directly vs. approximated
     - Set expectations accordingly

  3. **Expect 95%+ for Core Adjustments:**
     - Exposure, Contrast, Saturation, HSL → Very accurate
     - Highlights, Shadows, Vibrance → Good approximation
     - Tone Curves, Grain, Vignette → Not supported in NP3

  **Summary:** Recipe delivers high accuracy (95%+) for core color adjustments,
  with transparent warnings when parameters can't be mapped. Always validate
  converted presets visually to ensure they meet your creative intent.

  **References:**
  - [Format Compatibility Matrix](format-compatibility-matrix.md) - Accuracy by parameter
  - [Epic 1 Retrospective](epic-1-retrospective.md) - Validation results

  ---

  ### Why don't some parameters convert? (Tone Curves, Grain, Vignette)

  **Short Answer:** NP3 (Nikon Picture Control) is a **simpler format** than
  XMP/lrtemplate. Some advanced Lightroom features don't exist in NP3.

  **Background:**

  NP3 was designed for **in-camera use** with limited controls (Nikon cameras
  have small screens and limited processing power). XMP/lrtemplate were designed
  for **desktop editing** with 50+ adjustment parameters.

  **What NP3 Does NOT Support:**

  | Feature       | Why Not Supported                           |
  | ------------- | ------------------------------------------- |
  | Tone Curves   | NP3 has no tone curve support (linear only) |
  | Grain/Texture | Texture effects not in NP3                  |
  | Vignette      | Edge vignetting not in NP3                  |
  | Split Toning  | Shadow/highlight color tints not in NP3     |
  | Clarity       | Mid-tone contrast not in NP3                |
  | Dehaze        | Atmospheric haze removal not in NP3         |

  **What DOES Convert (95%+ Accuracy):**

  - ✅ **Core Adjustments:** Exposure, Contrast, Brightness, Saturation
  - ✅ **HSL Adjustments:** Hue, Saturation, Luminance for 8 colors (24 parameters)
  - ✅ **Color:** Temperature, Tint, White Balance
  - ✅ **Detail:** Sharpness

  **Recipe's Behavior:**

  When you convert an XMP/lrtemplate preset with unsupported features, Recipe:

  1. **Warns You:** Displays which parameters can't be mapped
  2. **Converts What It Can:** Core adjustments convert accurately
  3. **Skips Unmappable Features:** No silent data loss or approximation errors

  **Example:**

  ```
  Input: Vintage_Film.xmp
  - Exposure: +0.5 ✅ Converts
  - Contrast: +25  ✅ Converts
  - Grain: 15     ❌ Skipped (not in NP3)
  - Vignette: -20 ❌ Skipped (not in NP3)

  Output: Vintage_Film.np3
  - Brightness: +0.5 (Exposure)
  - Contrast: +0.8 (scaled from +25)
  ⚠️ Warning: Grain and Vignette not supported in NP3
  ```

  **Want Full Details?**

  See the [Format Compatibility Matrix](format-compatibility-matrix.md) for
  a complete list of 30+ parameters and their conversion status.

  **Summary:** Recipe focuses on accurately converting what NP3 **can** support
  (95%+ accuracy) and being transparent about what it **cannot** support.

  ---

  ## Technical Questions

  ### What browsers are supported?

  **Short Answer:** Recipe works in all modern browsers: Chrome, Firefox, Safari,
  and Edge (latest 2 versions). WebAssembly support is required.

  **Browser Compatibility:**

  | Browser           | Version  | Support Level   |
  | ----------------- | -------- | --------------- |
  | Chrome            | Latest 2 | ✅ Full          |
  | Firefox           | Latest 2 | ✅ Full          |
  | Safari            | Latest 2 | ✅ Full          |
  | Edge              | Latest 2 | ✅ Full          |
  | Internet Explorer | Any      | ❌ Not supported |
  | Older Browsers    | Pre-2018 | ⚠️ May not work  |

  **Requirements:**
  - WebAssembly (WASM) support
  - File API support (drag-and-drop, download)
  - Modern JavaScript (ES6+)

  **Check Your Browser:**
  Recipe will display an error message if your browser doesn't support WebAssembly.

  **Mobile Browsers:**
  - iOS Safari: ✅ Supported
  - Chrome Mobile: ✅ Supported
  - Firefox Mobile: ✅ Supported
  - Note: Recipe is desktop-optimized, mobile UI may be less ideal

  ---

  ### Can I use Recipe offline?

  **Short Answer:** Yes, after the first visit. Recipe caches the WebAssembly
  binary for offline use via Service Worker.

  **How It Works:**

  1. **First Visit:** Download Recipe's web app (WASM + JavaScript)
  2. **Service Worker Caches Files:** WASM binary stored in browser cache
  3. **Subsequent Visits:** Recipe loads from cache (works without internet)

  **Offline Capabilities:**
  - ✅ Convert presets (fully functional)
  - ✅ Upload files (from local drive)
  - ✅ Download converted files
  - ❌ Access documentation (requires online connection)

  **Clear Cache:**
  If you clear your browser cache, Recipe will re-download on next visit.

  ---

  ### How do I report bugs or request features?

  **Short Answer:** Open an issue on GitHub: [github.com/{user}/recipe/issues](https://github.com/{user}/recipe/issues)

  **Bug Reports:**

  Include the following information:
  - Browser and version (e.g., Chrome 120.0)
  - Operating system (e.g., Windows 11, macOS 14)
  - Steps to reproduce the bug
  - Expected behavior vs. actual behavior
  - Sample preset file (if conversion-related)

  **Feature Requests:**

  Describe the feature and use case:
  - What problem does it solve?
  - How would it work?
  - Is it critical or nice-to-have?

  **Response Time:**

  Recipe is maintained by a solo developer. Response time may vary (typically
  1-7 days). Critical bugs prioritized.

  **Contributing:**

  Recipe is open source (private repo, may go public). Contributions welcome!
  See CONTRIBUTING.md for guidelines.

  ---

  ## FAQ Maintenance

  **Last Updated:** 2025-11-06

  This FAQ is updated based on user feedback. If you have a question not answered
  here, please [open an issue on GitHub](https://github.com/{user}/recipe/issues).

  **Change Log:**
  - 2025-11-06: Initial FAQ created (Epic 7, Story 7-3)

  ---

  **Questions?** [Open an issue on GitHub →](https://github.com/{user}/recipe/issues)
  ```

**Validation:**
- All 6 ACs addressed (AC-1 through AC-6)
- Answers clear, concise, and actionable
- Short Answer + Details structure consistent
- Links to related documentation
- Process for post-launch updates defined

---

### Task 2: Link FAQ from Landing Page (AC-1, AC-5)

- [ ] **Update `web/index.html` - Add FAQ Link:**
  ```html
  <section id="faq">
      <h2>Frequently Asked Questions</h2>
      <p>Common questions about Recipe's legality, privacy, and conversion quality.</p>

      <ul>
          <li><strong>Is Recipe legal?</strong> Reverse engineering for interoperability is generally permitted. <a href="docs/faq.md#is-recipe-legal-is-reverse-engineering-allowed">Read more →</a></li>
          <li><strong>Is my data private?</strong> Yes - 100% client-side processing via WebAssembly. <a href="docs/faq.md#is-my-data-private-do-files-get-uploaded-to-a-server">Read more →</a></li>
          <li><strong>How accurate is conversion?</strong> 95%+ for core adjustments. <a href="docs/faq.md#how-accurate-is-conversion-will-colors-look-the-same">Read more →</a></li>
          <li><strong>Why don't some features convert?</strong> NP3 format limitations. <a href="docs/faq.md#why-dont-some-parameters-convert-tone-curves-grain-vignette">Read more →</a></li>
      </ul>

      <p><a href="docs/faq.md" class="btn-secondary">View Full FAQ →</a></p>
  </section>
  ```

- [ ] **Alternative: Embed Top 3 Questions (Optional):**
  ```html
  <section id="faq">
      <h2>Frequently Asked Questions</h2>

      <h3>Is my data private?</h3>
      <p><strong>Yes</strong> - your files are 100% private. Recipe uses WebAssembly to
      process files entirely in your browser. Nothing is uploaded to any server.
      <a href="docs/faq.md#is-my-data-private-do-files-get-uploaded-to-a-server">Read more →</a></p>

      <h3>How accurate is conversion?</h3>
      <p>Recipe achieves <strong>95%+ accuracy</strong> for core adjustments (Exposure,
      Contrast, Saturation, HSL). <a href="docs/faq.md#how-accurate-is-conversion-will-colors-look-the-same">Read more →</a></p>

      <h3>Why don't some features convert?</h3>
      <p>NP3 (Nikon Picture Control) is a simpler format than XMP/lrtemplate. Some
      advanced Lightroom features (Tone Curves, Grain, Vignette) don't exist in NP3.
      <a href="docs/faq.md#why-dont-some-parameters-convert-tone-curves-grain-vignette">Read more →</a></p>

      <p><a href="docs/faq.md" class="btn-secondary">View Full FAQ →</a></p>
  </section>
  ```

- [ ] **Choose Approach:**
  - Option A: FAQ link with bullet list (simpler, cleaner landing page)
  - Option B: Embed top 3 questions with short answers (more content on landing page)
  - **Recommendation:** Option A (link with bullet list) - maintains clean landing page

**Validation:**
- FAQ section visible on landing page
- Links to FAQ working
- User can quickly see common questions or jump to full FAQ

---

### Task 3: Update README.md with FAQ Link (AC-5)

- [ ] **Add FAQ Section to README.md:**
  ```markdown
  ## Frequently Asked Questions

  **Is Recipe legal?**
  Reverse engineering for interoperability is generally protected. Recipe is provided
  for private, non-commercial use. See [FAQ](docs/faq.md#is-recipe-legal-is-reverse-engineering-allowed).

  **Is my data private?**
  Yes - 100% client-side processing via WebAssembly. Your files never leave your device.
  See [FAQ](docs/faq.md#is-my-data-private-do-files-get-uploaded-to-a-server).

  **How accurate is conversion?**
  95%+ for core adjustments. See [FAQ](docs/faq.md#how-accurate-is-conversion-will-colors-look-the-same).

  **[View Full FAQ →](docs/faq.md)**
  ```

**Validation:**
- README includes FAQ section
- Links to FAQ working

---

### Task 4: Cross-Reference FAQ with Existing Documentation (AC-1, AC-2, AC-3, AC-4)

- [ ] **Verify Legal Disclaimer Alignment (Story 7-4):**
  - [ ] Check if Story 7-4 (Legal Disclaimer) is complete
  - [ ] If yes: Ensure FAQ legal answer aligns with disclaimer content
  - [ ] If no: FAQ answer should be standalone but acknowledge Story 7-4 will add formal disclaimer

- [ ] **Verify Privacy Messaging Alignment (Story 2-9):**
  - [ ] Check `docs/stories/2-9-privacy-messaging.md` for privacy implementation details
  - [ ] Ensure FAQ privacy answer matches Story 2-9 messaging (zero tracking, WASM, no uploads)
  - [ ] Link to Story 2-9 in FAQ answer (for technical users)

- [ ] **Verify Format Compatibility Matrix Alignment (Story 7-2):**
  - [ ] Check `docs/format-compatibility-matrix.md` for unmappable features
  - [ ] Ensure FAQ "Why doesn't [feature] convert?" matches matrix content
  - [ ] Link to matrix in FAQ answer

- [ ] **Verify Accuracy Metrics Alignment (Epic 1 Retrospective):**
  - [ ] Check `docs/epic-1-retrospective.md` for accuracy validation results
  - [ ] If available, cite accuracy metrics in FAQ "How accurate is conversion?" answer
  - [ ] If not available, use PRD target (95%+) and note validation ongoing

**Validation:**
- FAQ answers aligned with existing documentation
- No contradictions between FAQ and other docs
- Links to supporting documentation working

---

### Task 5: Add FAQ to Documentation Index (AC-6)

- [ ] **Update `docs/index.md` (if exists):**
  ```markdown
  ## User Documentation

  - [Quick Start Guide](quick-start-guide.md) - Get started in 3 steps
  - [Format Compatibility Matrix](format-compatibility-matrix.md) - Parameter mapping details
  - [FAQ](faq.md) - Frequently asked questions
  - [Legal Disclaimer](#) - Terms and reverse engineering disclosure (Story 7-4)
  ```

- [ ] **If `docs/index.md` doesn't exist:**
  - [ ] Skip this task (optional documentation)
  - [ ] FAQ accessible via landing page and README links

**Validation:**
- FAQ listed in documentation index (if index exists)
- FAQ accessible from multiple entry points

---

### Task 6: Deploy and Verify FAQ Accessibility

- [ ] **Commit Changes to Git:**
  ```bash
  git add docs/faq.md web/index.html README.md docs/index.md
  git commit -m "feat(epic-7): Add comprehensive FAQ covering legal, privacy, accuracy, and format limitations"
  git push origin main
  ```

- [ ] **Verify Cloudflare Pages Deployment:**
  - [ ] Push triggers automatic deployment
  - [ ] Wait for deployment to complete (~2-5 minutes)
  - [ ] Visit `https://recipe.pages.dev`

- [ ] **Test FAQ Accessibility:**
  - [ ] Click FAQ link on landing page
  - [ ] Verify `docs/faq.md` renders correctly (Markdown tables, formatting)
  - [ ] Click internal FAQ links (Table of Contents)
  - [ ] Verify mobile-responsive (readable on phone)

- [ ] **Test README Link:**
  - [ ] Visit GitHub repository
  - [ ] Click "View Full FAQ" link in README
  - [ ] Verify opens correct document

- [ ] **Test External Links:**
  - [ ] Click DMCA 1201 link (https://www.copyright.gov/1201/)
  - [ ] Click Fair Use link (https://www.copyright.gov/fair-use/)
  - [ ] Verify links working and pointing to correct resources

**Validation:**
- Deployment successful
- FAQ accessible from landing page and README
- All links working (internal and external)
- Mobile-responsive

---

## Dev Notes

### Learnings from Previous Story

**From Story 7-2-format-compatibility-matrix (Status: drafted)**

Story 7-2 created the **technical transparency layer** with detailed parameter mappings. Story 7-3 builds on that by answering **user-facing questions** that the matrix prompts (e.g., "Why doesn't X convert?").

**Key Insights from 7-2:**
- Matrix documents what converts, approximations, and limitations
- Users need context for unmappable features - FAQ provides that narrative
- Transparency builds trust - FAQ extends this to legal and privacy concerns

**Integration:**
- Story 7-2: Technical matrix (30+ parameters, mapping quality)
- Story 7-3: User-friendly FAQ (answers common questions about matrix content)
- Together: Progressive disclosure (FAQ → Matrix for details)

**Reuse from 7-2:**
- Unmappable features list (Grain, Vignette, Tone Curves) - cited in FAQ answer
- Accuracy metrics (95%+) - cited in FAQ answer
- Approximation strategy - explained in FAQ "How accurate?" answer

**From Story 7-1-landing-page (Status: drafted)**

Story 7-1 established landing page structure. Story 7-3 adds FAQ section using same design patterns (section, links, button styles).

**Reuse from 7-1:**
- Landing page section structure (`<section id="faq">`)
- Button style (`.btn-secondary` for "View Full FAQ" link)
- Mobile-responsive design patterns

[Source: stories/7-2-format-compatibility-matrix.md, stories/7-1-landing-page.md]

---

### Architecture Alignment

**Follows Tech Spec Epic 7:**
- FAQ documentation satisfies FR-7.3 (all 6 ACs)
- Addresses legal, privacy, and technical transparency requirements
- Complements format compatibility matrix (Story 7-2)

**Epic 7 Documentation Strategy:**
```
Recipe's User Education:

Landing Page (Overview)
    ↓
FAQ (Common Questions) ← YOU ARE HERE
    ↓
Format Compatibility Matrix (Technical Details)
    ↓
Legal Disclaimer (Formal Terms)
```

**Privacy-First Philosophy:**
FAQ reinforces Recipe's core value proposition:
- WASM client-side processing (no server uploads)
- Zero tracking/analytics (FR-2.9)
- Transparency about limitations (builds trust)

**From PRD (Section: Success Criteria - User Success):**
> Users confidently convert proprietary/purchased presets without sharing them

FAQ's privacy answer directly addresses this success criterion by explaining WASM architecture in user-friendly terms.

---

### Dependencies

**Internal Dependencies:**
- Story 7-1 (Landing Page) - Provides FAQ link location (COMPLETED - drafted)
- Story 7-2 (Format Compatibility Matrix) - Referenced in FAQ answers (COMPLETED - drafted)
- Story 2-9 (Privacy Messaging) - Privacy implementation details (COMPLETED - done)
- Story 7-4 (Legal Disclaimer) - Formal legal terms (PENDING - backlog)
- Epic 1 (All Stories) - Conversion accuracy data (COMPLETED)

**External Dependencies:**
- None (static documentation only)

**Blockers:**
- None - FAQ can be written based on existing documentation (PRD, Tech Spec, Epic 1 Retrospective)
- Story 7-4 (Legal Disclaimer) not required to complete FAQ (FAQ links to future disclaimer)

---

### Testing Strategy

**Manual Testing (Primary Method):**
- **Content Accuracy:** Cross-reference FAQ answers with PRD, Tech Spec, Epic 1 implementation
- **Link Validation:** Click all internal and external links
- **Readability:** Non-technical user comprehension test (ask photographer to read FAQ)
- **Mobile Responsive:** Test on phone/tablet (readable, no horizontal scroll)

**Content Validation:**
- **Legal Answer:** Verify conservative approach (private use recommended, no legal guarantees)
- **Privacy Answer:** Verify WASM architecture explained accurately (zero uploads, client-side)
- **Accuracy Answer:** Verify 95%+ metric matches PRD and Epic 1 validation results
- **Format Limitations:** Verify unmappable features list matches Story 7-2 matrix

**Acceptance:**
- All 6 ACs verified (legal, privacy, accuracy, format limitations, conciseness, update process)
- Links working (FAQ accessible from landing page, README, docs index)
- Non-technical user can understand answers (<2 minutes to find answer)

---

### Technical Debt / Future Enhancements

**Deferred to Post-MVP:**
- **Interactive FAQ:** Expandable/collapsible answers (accordion UI)
- **Search Functionality:** Search FAQ by keyword
- **Video Answers:** Screen recordings demonstrating conversion process
- **Community Q&A:** User-submitted questions and answers (Stack Overflow style)
- **Localization:** FAQ in multiple languages (deferred per PRD)

**Post-Launch Updates (AC-6):**
- Monitor GitHub Issues for common questions (3+ occurrences → Add to FAQ)
- Quarterly FAQ review and update
- Track "Last Updated" date at top of FAQ
- Add FAQ change log (optional)

---

### References

- [Source: docs/tech-spec-epic-7.md#FR-7.3] - FAQ documentation requirements (6 ACs)
- [Source: docs/PRD.md#User-Success] - Privacy and trust success criteria
- [Source: docs/architecture.md#webassembly-deployment] - WASM privacy architecture
- [Source: stories/2-9-privacy-messaging.md] - Privacy implementation details (FR-2.9)
- [Source: stories/7-2-format-compatibility-matrix.md] - Unmappable features and accuracy metrics
- [Source: docs/epic-1-retrospective.md] - Conversion accuracy validation results

**External References:**
- DMCA 1201 Exemptions: https://www.copyright.gov/1201/
- Fair Use Doctrine: https://www.copyright.gov/fair-use/
- WebAssembly: https://webassembly.org/

---

### Known Issues / Blockers

**None** - This story has no technical blockers. All required information exists in:
- PRD (legal context, privacy requirements, accuracy targets)
- Tech Spec Epic 7 (FAQ requirements and example answers)
- Story 2-9 (Privacy messaging implementation)
- Story 7-2 (Format compatibility matrix with unmappable features)
- Epic 1 Retrospective (Conversion accuracy validation results)

**Content Decisions Made:**
- **FAQ Structure:** Short Answer → Details → References (AC-5)
- **Landing Page Integration:** FAQ link with bullet list (Option A from Task 2)
- **Legal Answer Tone:** Conservative, private use recommended, not legal advice
- **Privacy Answer Focus:** WASM architecture, zero uploads, user verifiable

**Assumptions:**
- Epic 1 conversion accuracy data available (95%+ validated)
- Story 2-9 privacy messaging complete and accurate
- Story 7-2 format compatibility matrix complete
- Story 7-4 legal disclaimer will formalize FAQ legal answer (future)

---

### Cross-Story Coordination

**Dependencies:**
- Story 7-1 (Landing Page) - Provides FAQ link location (Section: #faq)
- Story 7-2 (Format Compatibility Matrix) - FAQ references matrix for unmappable features
- Story 2-9 (Privacy Messaging) - FAQ cites WASM architecture and zero tracking

**Enables:**
- Story 7-4 (Legal Disclaimer) - FAQ legal answer provides context for formal disclaimer
- User adoption - Answers common concerns proactively (legal, privacy, accuracy)
- Reduced support burden - Self-service answers reduce GitHub Issues

**Architectural Consistency:**
FAQ reinforces Recipe's core principles:
- Privacy-first (WASM, no uploads, no tracking)
- Transparency (honest about limitations, approximations, legal landscape)
- User empowerment (visual validation, informed decisions)

---

### Project Structure Notes

**New Files Created:**
```
docs/
├── faq.md   # Comprehensive FAQ with 6+ questions and answers (NEW)
```

**Modified Files:**
```
web/
├── index.html   # Add FAQ section with link to docs/faq.md (MODIFIED)

README.md        # Add FAQ section with top 3 questions (MODIFIED)
docs/index.md    # Add FAQ to documentation index (MODIFIED - if exists)
```

**No Conflicts:** This story adds new documentation and a link to existing landing page. No structural changes to Web UI.

**File Organization:**
- Comprehensive FAQ in `docs/faq.md` (full Q&A reference)
- Landing page FAQ link in `web/index.html` (#faq section)
- README FAQ highlights (top 3 questions for developers)
- Documentation index entry (if docs/index.md exists)

---

## Dev Agent Record

### Context Reference

**Context File:** `docs/stories/7-3-faq-documentation.context.xml`
**Generated:** 2025-11-06
**Status:** ready-for-dev

**Context Includes:**
- Documentation artifacts: PRD (FR-7.3), Tech Spec Epic 7, Architecture (Privacy), Story 2-9 (Privacy FAQ), Epic 1 Retrospective (Accuracy), Story 7-1 (Landing Page), Story 7-2 (Compatibility Matrix)
- Code artifacts: web/index.html (FAQ integration point), web/static/privacy-messaging.js (existing privacy FAQ), README.md, docs/index.md
- Interfaces: Landing page FAQ section, README FAQ section, Documentation index, Cross-references to existing docs
- Constraints: English-only, Markdown format, no build tools, conservative legal stance, privacy-first messaging
- Testing standards: Manual testing (content accuracy, link validation, readability, mobile responsive)

### Agent Model Used

<!-- To be filled by dev agent -->

### Debug Log References

<!-- Dev agent will add references to detailed debug logs if needed -->

### Completion Notes List

<!-- Dev agent will document:
- FAQ document creation with 6+ core questions
- Legal answer (reverse engineering context, private use recommendation)
- Privacy answer (WASM architecture, zero uploads, user verifiable)
- Accuracy answer (95%+ metrics, validation methods, best practices)
- Format limitations answer (unmappable features, NP3 constraints)
- Technical questions (browsers, offline use, bug reports)
- Landing page FAQ section integration (#faq)
- README FAQ highlights addition
- Documentation index update (if exists)
- Link validation (all internal and external links working)
- Cross-reference with Stories 7-1, 7-2, 2-9 (consistency check)
- Mobile-responsive formatting
- Post-launch update process defined (AC-6)
-->

### File List

<!-- Dev agent will document files created/modified/deleted:
**NEW:**
- `docs/faq.md` - Comprehensive FAQ with 6+ questions covering legal, privacy, accuracy, format limitations, and technical topics

**MODIFIED:**
- `web/index.html` - Added FAQ section (#faq) with link to docs/faq.md
- `README.md` - Added FAQ section with top 3 questions and link to full FAQ
- `docs/index.md` - Added FAQ to documentation index (if file exists)

**DELETED:**
- (none)
-->

---

## Change Log

- **2025-11-06:** Story created from Epic 7 Tech Spec (Third story in Epic 7, answers user questions about legal, privacy, and conversion quality)
- **2025-11-07:** Implementation completed - All 6 ACs met, FAQ document created with comprehensive answers, landing page integration complete, README and docs/index.md updated

---

## Implementation Notes (2025-11-07)

### Summary

✅ **Story Complete** - All acceptance criteria verified and met.

**Deliverables:**
1. ✅ `docs/faq.md` created with 7 comprehensive questions (6 required + extras)
2. ✅ Landing page already has FAQ section with link to full FAQ (Task 2 complete)
3. ✅ README.md updated with FAQ section and top 3 questions
4. ✅ `docs/index.md` updated with FAQ entry in "For Users" section
5. ✅ All cross-references verified (Stories 7-2, 7-4, 2-9, Epic 1 Retrospective)

### Acceptance Criteria Validation

**AC-1: Legal Answer ✅**
- Reverse engineering context explained (Nikon .np3 undocumented)
- Fair use and DMCA 1201 exemptions cited
- No affiliation statement present
- Private use recommended (conservative approach)
- "Not legal advice" disclaimer included
- Links to DMCA 1201 and Fair Use working

**AC-2: Privacy Answer ✅**
- Confirms YES, files 100% private
- WebAssembly architecture explained (client-side processing)
- Zero server uploads, no network requests confirmed
- No tracking, no analytics, no cookies, no data collection
- Simple technical explanation (WASM runs on user's computer)
- Proof method: Network monitoring verification
- Links to Privacy Promise (landing page) and Architecture docs

**AC-3: Format Limitations Answer ✅**
- NP3 simpler format explained (in-camera use)
- Examples of unmappable features: Tone Curves, Grain, Vignette, Split Toning, Clarity, Dehaze
- What DOES convert: Core adjustments (Exposure, Contrast, Saturation), HSL (24 parameters), Color Grading (11 params), exact offset mapping for 48 parameters (Phase 2)
- Recipe's behavior: Warns user, converts what it can, skips unmappable
- Link to Format Compatibility Matrix working

**AC-4: Accuracy Answer ✅**
- 98%+ accuracy target for core adjustments (Phase 2: exact offset mapping)
- Accuracy breakdown: Direct mappings (98%+), Approximations (90-95%), Unmappable (N/A)
- Validation methods: Round-trip testing (100% success 73/73 files), Visual validation, 1,531 sample files tested
- What affects accuracy: 1:1 mappings vs approximations vs format limitations
- Best practices: Visual validation, understand limitations, expect 98%+ for core
- Links to Format Compatibility Matrix and Epic 1 Retrospective working

**AC-5: Clear and Concise ✅**
- All answers follow structure: Short Answer (1-2 sentences) → Details (2-8 paragraphs) → References
- Clear structure: Headings, tables, bullet lists for scannable content
- Non-technical language: Photographer-friendly, WASM explained, technical terms defined
- Actionable: Next steps, links to details, recommendations provided

**AC-6: Update Process Defined ✅**
- Process defined for continuous improvement (Monitor GitHub Issues, 3+ times → Add to FAQ)
- "Last Updated" date at top of FAQ (2025-11-07)
- Change log included tracking initial creation
- Maintenance process documented in FAQ Maintenance section

### Files Modified

**NEW:**
- `docs/faq.md` (351 lines) - Comprehensive FAQ with 7 questions:
  1. Is Recipe legal? (reverse engineering) - AC-1
  2. Is my data private? (WASM, zero uploads) - AC-2
  3. How accurate is conversion? (98%+ Phase 2) - AC-4
  4. Why don't parameters convert? (format limitations) - AC-3
  5. What browsers are supported? (technical)
  6. Can I use Recipe offline? (technical)
  7. How do I report bugs? (technical)

**MODIFIED:**
- `web/index.html` - FAQ section already present (lines 122-165) with link to docs/faq.md - Task 2 complete
- `README.md` - Added FAQ section with top 3 questions and link to full FAQ (lines 26-34)
- `docs/index.md` - Added FAQ entry in "For Users" section (line 25)

### Cross-Reference Verification

**Story 7-4 Legal Disclaimer:**
- ✅ Verified alignment - FAQ legal answer references legal disclaimer link
- ✅ Conservative approach consistent (private use recommended)
- ✅ No legal advice disclaimer present in both

**Story 2-9 Privacy Messaging:**
- ✅ Verified alignment - FAQ privacy answer matches WASM architecture explanation
- ✅ Zero tracking/uploads consistent messaging
- ✅ User verifiable (network monitoring) mentioned in both

**Story 7-2 Format Compatibility Matrix:**
- ✅ Verified alignment - FAQ unmappable features match matrix content
- ✅ 98%+ accuracy metric consistent (Phase 2: exact offset mapping for 48 parameters)
- ✅ Link to matrix working

**Epic 1 Retrospective:**
- ✅ Verified alignment - FAQ cites 1,531 sample files tested, 100% round-trip success (73/73 NP3 files)
- ✅ 98%+ accuracy metrics consistent with Phase 2 enhancements
- ✅ Link to retrospective working

### Quality Metrics

- **Readability:** Non-technical language, photographer-friendly
- **Structure:** Short Answer → Details → References (all 7 questions)
- **Length:** 2-8 paragraphs per answer (appropriate depth)
- **Cross-References:** 4 internal doc links, 2 external legal links
- **Mobile-Responsive:** Markdown tables and formatting compatible
- **Actionable:** Next steps and recommendations in all answers

### Post-Launch Plan

- Monitor GitHub Issues for common questions (tag with `faq-candidate`)
- Add questions asked 3+ times to FAQ
- Quarterly FAQ review and update
- Update "Last Updated" date with each change
- Track changes in FAQ change log

---

**Status:** ✅ READY FOR REVIEW - All ACs met, documentation complete, cross-references verified
