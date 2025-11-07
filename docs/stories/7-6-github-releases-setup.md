# Story 7.6: GitHub Releases Setup

**Epic:** Epic 7 - Documentation & Deployment (FR-7)
**Story ID:** 7.6
**Status:** ready-for-dev
**Created:** 2025-11-06
**Complexity:** Medium (2-3 days)

---

## Story

As a **Recipe user**,
I want **multi-platform CLI binaries automatically built and published to GitHub Releases whenever a version tag is pushed**,
so that **I can download and run Recipe CLI on my platform (Linux/macOS/Windows, amd64/arm64) without needing to install Go or build from source**.

---

## Business Value

GitHub Releases automation is Recipe's **zero-cost CLI distribution infrastructure**, providing instant access to production-ready binaries for all supported platforms with every version release.

**Strategic Value:**
- **Zero Friction Distribution:** Users download pre-built binaries (no Go installation, no compilation)
- **Multi-Platform Support:** 6 binaries covering 99%+ of desktop/server platforms (Linux/macOS/Windows × amd64/arm64)
- **Automatic Versioning:** Semantic version tags trigger releases (clear version history)
- **Zero Cost:** GitHub Releases free tier supports unlimited releases, 2GB per release
- **Professional Delivery:** GitHub Release pages provide changelog, download counts, asset management

**Developer Impact:**
- Eliminates manual binary builds (automation replaces ~30 min manual process per release)
- Provides version history and rollback capability (GitHub Release archive)
- Enables faster release cadence (tag → binaries in <10 minutes)
- Reduces distribution errors (consistent build environment, reproducible builds)

**User Impact:**
- **Simple Installation:** Download binary → Run (no dependencies)
- **Platform Choice:** Users choose native binary for their OS/architecture
- **Version Control:** Users select specific version to download (stable releases, beta versions)
- **Trusted Source:** GitHub Releases = official distribution channel (reduces malware risk)

**Risk Mitigation:**
- Build failures visible immediately (GitHub Actions red X, email notification)
- Checksums provided for download verification (security against tampering)
- Release process documented (CHANGELOG.md, semantic versioning strategy)

---

## Acceptance Criteria

### AC-1: GitHub Actions Workflow Triggers on Version Tag Push

**Given** a semantic version tag is created and pushed to the repository
**When** the tag push completes
**Then**:
- ✅ **Workflow File Exists:**
  - File path: `.github/workflows/release.yml`
  - Committed to repository (not local-only)
- ✅ **Trigger Configuration:**
  ```yaml
  on:
    push:
      tags: ['v*']
  ```
- ✅ **Tag Pattern Matching:**
  - Matches tags starting with `v` (e.g., v1.0.0, v0.1.0, v2.5.3)
  - Does NOT trigger on other tags (e.g., test-tag, build-123)
  - Does NOT trigger on branch pushes (main, feature branches)
- ✅ **GitHub Actions Log Visible:**
  - Navigate to repository → Actions tab
  - See workflow run listed with tag name
  - Click workflow run to view detailed logs

**Validation:**
- Create test tag: `git tag v0.0.1 && git push origin v0.0.1`
- Verify workflow run appears in GitHub Actions within 30 seconds
- Verify workflow does NOT run on branch push
- Verify workflow does NOT run on non-version tag

---

### AC-2: Workflow Builds 6 CLI Binaries

**Given** the release workflow is triggered by a version tag
**When** the workflow executes the build matrix
**Then**:
- ✅ **Build Matrix Configuration:**
  ```yaml
  strategy:
    matrix:
      os: [linux, darwin, windows]
      arch: [amd64, arm64]
  ```
- ✅ **6 Binaries Built:**
  1. `recipe-linux-amd64` (Linux 64-bit Intel/AMD)
  2. `recipe-linux-arm64` (Linux 64-bit ARM - Raspberry Pi, AWS Graviton)
  3. `recipe-darwin-amd64` (macOS Intel)
  4. `recipe-darwin-arm64` (macOS Apple Silicon - M1/M2/M3)
  5. `recipe-windows-amd64.exe` (Windows 64-bit Intel/AMD)
  6. `recipe-windows-arm64.exe` (Windows 64-bit ARM - Surface Pro X)
- ✅ **Build Command:**
  ```yaml
  - name: Build
    run: |
      BINARY_NAME=recipe-${{ matrix.os }}-${{ matrix.arch }}
      if [ "${{ matrix.os }}" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
      fi
      GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -ldflags="-s -w" -o $BINARY_NAME cmd/cli/main.go
  ```
- ✅ **Build Flags:**
  - `-ldflags="-s -w"` strips debug symbols (reduces binary size ~30%)
  - `GOOS` and `GOARCH` set for cross-compilation
  - Output binary named by platform: `recipe-{os}-{arch}`
- ✅ **Parallel Execution:**
  - 6 build jobs run in parallel (faster than sequential)
  - Total build time <10 minutes (AC-6 requirement)

**Validation:**
- Check workflow logs for "Build" step in each matrix job
- Verify 6 build jobs executed (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64, windows-arm64)
- Verify binaries created (workflow artifacts or release assets)
- Locally test cross-compilation: `GOOS=linux GOARCH=amd64 go build -o recipe-linux-amd64 cmd/cli/main.go`

---

### AC-3: Binaries are Uploaded to GitHub Release

**Given** all binaries are built successfully
**When** the workflow executes the release upload step
**Then**:
- ✅ **GitHub Release Action:**
  ```yaml
  - name: Upload to GitHub Release
    uses: softprops/action-gh-release@v1
    with:
      files: |
        recipe-linux-amd64
        recipe-linux-arm64
        recipe-darwin-amd64
        recipe-darwin-arm64
        recipe-windows-amd64.exe
        recipe-windows-arm64.exe
  ```
- ✅ **Release Created:**
  - GitHub Release created automatically with tag name (e.g., v0.0.1)
  - Release visible at: `https://github.com/{user}/recipe/releases/tag/{tag}`
- ✅ **Assets Attached:**
  - All 6 binaries attached to release as downloadable assets
  - Asset names match binary names (recipe-linux-amd64, etc.)
  - Download URLs functional (users can download binaries)
- ✅ **Release Page:**
  - Release title: Tag name (e.g., v0.0.1)
  - Release description: Changelog excerpt (if provided) or auto-generated
  - Assets section lists all 6 binaries with download counts

**Validation:**
- Navigate to: Repository → Releases
- Verify release exists with correct tag name
- Click release → Verify 6 binaries listed in "Assets" section
- Click binary download link → Verify file downloads successfully
- Check binary filename matches expected pattern

---

### AC-4: Release Includes CHANGELOG Excerpt

**Given** a GitHub Release is created
**When** viewing the release page
**Then**:
- ✅ **Release Notes Content:**
  - Release description includes changes for this version
  - Extracted from CHANGELOG.md (if present)
  - OR: Auto-generated from git commit messages
  - OR: Placeholder text prompting manual edit
- ✅ **CHANGELOG.md Format:**
  - Follows Keep a Changelog standard: https://keepachangelog.com/
  - Sections: Added, Changed, Fixed, Removed
  - Version entry format: `## [1.0.0] - YYYY-MM-DD`
- ✅ **Automatic Extraction (Optional):**
  - Workflow extracts version section from CHANGELOG.md
  - Release body populated with excerpt
  - Fallback: Manual edit after release creation

**Example Release Notes:**
```markdown
## Changes in v0.1.0

### Added
- NP3 binary parser with full parameter support
- XMP XML parser for Lightroom presets
- lrtemplate Lua parser with develop settings extraction
- CLI convert command with batch processing
- Web interface with drag-and-drop upload

### Fixed
- Parameter mapping accuracy for tone curve conversion
- WASM binary size reduced by 30% with ldflags optimization

Full changelog: https://github.com/{user}/recipe/blob/main/CHANGELOG.md
```

**Validation:**
- View release page → Verify description not empty
- If CHANGELOG.md exists → Verify excerpt matches version entry
- If CHANGELOG.md missing → Verify placeholder or auto-generated notes
- Verify "Full changelog" link points to CHANGELOG.md

---

### AC-5: Semantic Versioning Followed

**Given** a new version release is planned
**When** creating the version tag
**Then**:
- ✅ **Tag Format:**
  - Pattern: `v{MAJOR}.{MINOR}.{PATCH}` (e.g., v1.0.0, v1.1.0, v1.1.1)
  - Prefix: `v` required (matches workflow trigger pattern `v*`)
- ✅ **Versioning Rules:**
  - **MAJOR:** Breaking changes (format incompatibility, API changes)
  - **MINOR:** New features (backward compatible, new format support)
  - **PATCH:** Bug fixes (no new features, backward compatible)
- ✅ **CHANGELOG.md Documents Strategy:**
  - Versioning section explains semantic versioning
  - References: https://semver.org/spec/v2.0.0.html
  - Examples of MAJOR/MINOR/PATCH changes
- ✅ **Initial Version:**
  - First release: v0.1.0 (beta/experimental) or v1.0.0 (stable)
  - Recommendation: v0.1.0 for MVP (allows breaking changes)

**Semantic Versioning Examples:**
- `v0.1.0` → `v0.2.0` - Added TUI interface (new feature, backward compatible)
- `v1.0.0` → `v1.0.1` - Fixed WASM conversion bug (bug fix)
- `v1.5.0` → `v2.0.0` - Changed NP3 format structure (breaking change)

**Validation:**
- Review CHANGELOG.md → Verify versioning strategy documented
- Check git tags → Verify all tags follow vX.Y.Z pattern
- Test workflow → Verify only `v*` tags trigger release

---

### AC-6: Build Completes in <10 Minutes

**Given** a version tag is pushed
**When** the release workflow runs from start to finish
**Then**:
- ✅ **Total Duration <10 Minutes:**
  - Measured from: GitHub Actions workflow start
  - Measured to: All binaries uploaded to release (workflow success)
  - Target: <10 minutes total
  - Acceptable: <15 minutes (set workflow timeout to 15 minutes as safety)
- ✅ **Timing Breakdown:**
  - Checkout repository: <10 seconds
  - Setup Go: <1 minute
  - Build matrix (6 jobs in parallel): 3-5 minutes each
  - Upload binaries: <1 minute
  - Total: ~5-8 minutes typically
- ✅ **Timeout Configuration:**
  ```yaml
  jobs:
    build:
      runs-on: ubuntu-latest
      timeout-minutes: 15  # Fail if exceeds 15 minutes
  ```
- ✅ **Parallel Builds:**
  - Matrix jobs run concurrently (not sequential)
  - Total time = slowest job + upload time (not sum of all jobs)

**Validation:**
- Monitor GitHub Actions run duration (displayed in Actions tab)
- Verify workflow completes in <10 minutes for typical release
- Verify workflow times out at 15 minutes if stuck (safety net)
- Document timing in workflow logs for future reference

---

### AC-7: README.md Updated with Installation Instructions

**Given** CLI binaries are available on GitHub Releases
**When** a user reads README.md
**Then**:
- ✅ **Installation Section Exists:**
  - Section title: "Installation" or "Getting Started"
  - Instructions for downloading and running CLI
- ✅ **Platform-Specific Instructions:**
  - **Linux/macOS:**
    ```bash
    # Download latest release (replace {version} and {os}-{arch})
    curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-{os}-{arch}

    # Make executable
    chmod +x recipe-{os}-{arch}

    # Move to PATH (optional)
    sudo mv recipe-{os}-{arch} /usr/local/bin/recipe

    # Verify installation
    recipe --version
    ```
  - **Windows:**
    ```powershell
    # Download from GitHub Releases page
    # https://github.com/{user}/recipe/releases/latest

    # Add to PATH (optional)
    # Move recipe-windows-amd64.exe to C:\Program Files\recipe\
    # Add C:\Program Files\recipe\ to system PATH

    # Verify installation
    recipe.exe --version
    ```
- ✅ **Latest Release Link:**
  - Direct link to: `https://github.com/{user}/recipe/releases/latest`
  - Users can browse all releases: `https://github.com/{user}/recipe/releases`
- ✅ **Build from Source (Alternative):**
  ```bash
  # Requirements: Go 1.24+
  git clone https://github.com/{user}/recipe.git
  cd recipe
  go build cmd/cli/main.go -o recipe
  ```

**Validation:**
- Review README.md → Verify installation section present
- Follow instructions manually → Verify binary downloads and runs
- Check links → Verify GitHub Releases URLs functional
- Test on each platform (Linux, macOS, Windows)

---

### AC-8: CHANGELOG.md Maintained

**Given** Recipe project uses semantic versioning
**When** preparing a new release
**Then**:
- ✅ **CHANGELOG.md File Exists:**
  - Location: Repository root (same directory as README.md)
  - Format: Markdown (.md)
- ✅ **Follows Keep a Changelog Format:**
  - Header: "# Changelog"
  - Intro: "All notable changes to Recipe will be documented in this file."
  - Format reference: "The format is based on [Keep a Changelog](https://keepachangelog.com/)"
  - Versioning reference: "This project adheres to [Semantic Versioning](https://semver.org/)"
- ✅ **Version Entries:**
  - Latest version at top (reverse chronological order)
  - Format: `## [Version] - YYYY-MM-DD`
  - Sections: Added, Changed, Fixed, Removed, Deprecated, Security
  - Unreleased section for in-progress changes
- ✅ **Example Content:**
  ```markdown
  # Changelog
  All notable changes to Recipe will be documented in this file.
  The format is based on [Keep a Changelog](https://keepachangelog.com/),
  and this project adheres to [Semantic Versioning](https://semver.org/).

  ## [Unreleased]
  ### Added
  - (Changes in main branch not yet released)

  ## [0.1.0] - 2025-11-06
  ### Added
  - Initial release with NP3, XMP, lrtemplate support
  - CLI interface with convert and batch commands
  - Web interface with drag-and-drop upload
  - WASM conversion engine

  ## [0.0.1] - 2025-11-06
  ### Added
  - Initial pre-release for testing
  ```
- ✅ **Update Process:**
  - Before creating version tag: Update CHANGELOG.md with version entry
  - Commit CHANGELOG.md: `git commit -m "chore: prepare v0.1.0 release"`
  - Create tag: `git tag v0.1.0`
  - Push tag: `git push origin v0.1.0`

**Validation:**
- Verify CHANGELOG.md exists in repository root
- Verify format matches Keep a Changelog standard
- Verify version entries use semantic versioning
- Verify Unreleased section exists for future changes

---

## Tasks / Subtasks

### Task 1: Create CHANGELOG.md (AC-8)

- [ ] **Create CHANGELOG.md File:**
  ```bash
  touch CHANGELOG.md
  ```

- [ ] **Write Initial Changelog Content:**
  ```markdown
  # Changelog

  All notable changes to Recipe will be documented in this file.

  The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
  and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

  ## [Unreleased]
  ### Added
  ### Changed
  ### Fixed
  ### Removed

  ## [0.1.0] - 2025-11-06
  ### Added
  - Universal Recipe data model for format-agnostic parameter representation
  - NP3 binary parser and generator (Nik Collection presets)
  - XMP XML parser and generator (Adobe Lightroom presets)
  - lrtemplate Lua parser and generator (Lightroom templates)
  - Parameter mapping rules for bidirectional conversion between formats
  - Metadata field implementation (description, author, keywords)
  - Web interface with drag-and-drop file upload
  - File upload handling with 10MB size limit
  - Format auto-detection for NP3, XMP, lrtemplate
  - Parameter preview display with expandable categories
  - Target format selection with compatibility warnings
  - WASM conversion execution (client-side, zero-latency)
  - File download trigger for converted presets
  - Error handling UI with user-friendly messages
  - Privacy messaging (zero tracking, client-side processing)
  - Responsive design (mobile, tablet, desktop)
  - CLI interface with Cobra framework
  - Convert command for single file conversion
  - Batch processing for multiple files with progress tracking
  - Format auto-detection in CLI
  - Verbose logging with structured slog
  - JSON output mode for programmatic use
  - Cloudflare Pages deployment automation
  - GitHub Releases setup for CLI binary distribution

  ### Performance
  - WASM conversion: <100ms average (target met)
  - Batch processing: 37ms for 100 files (53x faster than target)
  - Format detection: 1.60ms average (1000x+ faster than target)

  ### Testing
  - 1,501 sample files tested across all formats
  - 95%+ conversion accuracy achieved
  - Round-trip testing validates bidirectional conversion

  ## [0.0.1] - 2025-11-06
  ### Added
  - Initial pre-release for testing infrastructure
  ```

- [ ] **Commit CHANGELOG.md:**
  ```bash
  git add CHANGELOG.md
  git commit -m "docs: Add CHANGELOG.md with Keep a Changelog format"
  git push origin main
  ```

**Validation:**
- CHANGELOG.md exists in repository root
- Content follows Keep a Changelog format
- Unreleased section present
- Initial version (0.1.0) documented with comprehensive changes
- Links to Keep a Changelog and Semantic Versioning present

---

### Task 2: Create GitHub Actions Workflow File (AC-1, AC-2, AC-3)

- [ ] **Create Workflow Directory (if not exists):**
  ```bash
  mkdir -p .github/workflows
  ```

- [ ] **Create Workflow File:**
  ```bash
  touch .github/workflows/release.yml
  ```

- [ ] **Write Workflow Configuration:**
  ```yaml
  name: Build and Release CLI Binaries

  on:
    push:
      tags:
        - 'v*'

  jobs:
    build:
      runs-on: ubuntu-latest
      timeout-minutes: 15

      strategy:
        matrix:
          os: [linux, darwin, windows]
          arch: [amd64, arm64]

      steps:
        - name: Checkout repository
          uses: actions/checkout@v4

        - name: Setup Go
          uses: actions/setup-go@v5
          with:
            go-version: '1.24'

        - name: Build binary
          run: |
            BINARY_NAME=recipe-${{ matrix.os }}-${{ matrix.arch }}
            if [ "${{ matrix.os }}" = "windows" ]; then
              BINARY_NAME="${BINARY_NAME}.exe"
            fi
            GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -ldflags="-s -w" -o $BINARY_NAME cmd/cli/main.go

        - name: Generate checksums
          run: |
            BINARY_NAME=recipe-${{ matrix.os }}-${{ matrix.arch }}
            if [ "${{ matrix.os }}" = "windows" ]; then
              BINARY_NAME="${BINARY_NAME}.exe"
            fi
            sha256sum $BINARY_NAME > $BINARY_NAME.sha256

        - name: Upload binaries to release
          uses: softprops/action-gh-release@v1
          with:
            files: |
              recipe-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.os == 'windows' && '.exe' || '' }}
              recipe-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.os == 'windows' && '.exe' || '' }}.sha256
          env:
            GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  ```

- [ ] **Workflow Configuration Details:**
  - **Trigger:** Push to tags matching `v*` pattern
  - **Runner:** Ubuntu latest (GitHub-hosted, free tier)
  - **Timeout:** 15 minutes (safety net, typical build <10 minutes)
  - **Matrix:** 3 OS × 2 architectures = 6 parallel jobs
  - **Steps:**
    1. Checkout repository (actions/checkout@v4)
    2. Setup Go 1.24 (actions/setup-go@v5)
    3. Build binary (GOOS/GOARCH cross-compilation)
    4. Generate SHA256 checksums (security verification)
    5. Upload to GitHub Release (softprops/action-gh-release@v1)

**Validation:**
- Workflow file exists at `.github/workflows/release.yml`
- YAML syntax valid (use online YAML validator or `yamllint`)
- Workflow committed to repository (not local-only)

---

### Task 3: Test Workflow with Test Tag (AC-1, AC-2, AC-3, AC-6)

- [ ] **Create Test Tag:**
  ```bash
  git tag v0.0.1
  git push origin v0.0.1
  ```

- [ ] **Monitor Workflow Execution:**
  - Navigate to: Repository → Actions → "Build and Release CLI Binaries"
  - Click workflow run triggered by v0.0.1 tag
  - Expand matrix jobs to view 6 parallel builds

- [ ] **Check Build Matrix Jobs:**
  - Verify 6 jobs running:
    - linux-amd64
    - linux-arm64
    - darwin-amd64
    - darwin-arm64
    - windows-amd64
    - windows-arm64
  - Each job should complete successfully (green checkmark)

- [ ] **Check Workflow Logs:**
  - **Setup Go:** Verify Go 1.24 installed
  - **Build binary:** Verify binary created (recipe-{os}-{arch})
  - **Generate checksums:** Verify .sha256 file created
  - **Upload to release:** Verify upload successful

- [ ] **Verify GitHub Release Created:**
  - Navigate to: Repository → Releases
  - Verify release "v0.0.1" exists
  - Click release → Verify 6 binaries + 6 checksum files attached (12 assets total)

- [ ] **Test Binary Downloads:**
  - Download one binary (e.g., recipe-linux-amd64)
  - Download corresponding checksum file
  - Verify checksum:
    ```bash
    sha256sum -c recipe-linux-amd64.sha256
    # Output: recipe-linux-amd64: OK
    ```

- [ ] **Measure Workflow Duration:**
  - Check GitHub Actions run duration (top of workflow page)
  - Verify duration <10 minutes
  - Document actual time

**Validation:**
- Test tag triggers workflow successfully
- 6 binaries built in parallel
- All binaries uploaded to GitHub Release
- Checksums generated for each binary
- Workflow completes in <10 minutes
- Binaries downloadable and verified

---

### Task 4: Update README.md with Installation Instructions (AC-7)

- [ ] **Add Installation Section to README.md:**

  Replace or enhance existing installation section with:

  ```markdown
  ## Installation

  Recipe is distributed as pre-built binaries for Linux, macOS, and Windows. Choose the binary for your platform and architecture.

  ### Download Pre-Built Binaries

  **Latest Release:** [Download from GitHub Releases](https://github.com/{user}/recipe/releases/latest)

  #### Linux / macOS

  ```bash
  # Download latest release (choose your platform)
  # Linux amd64 (Intel/AMD 64-bit)
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-linux-amd64

  # Linux arm64 (ARM 64-bit - Raspberry Pi, AWS Graviton)
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-linux-arm64

  # macOS amd64 (Intel Mac)
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-darwin-amd64

  # macOS arm64 (Apple Silicon - M1/M2/M3)
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-darwin-arm64

  # Make executable
  chmod +x recipe-*

  # Move to PATH (optional)
  sudo mv recipe-* /usr/local/bin/recipe

  # Verify installation
  recipe --version
  ```

  #### Windows

  1. Visit [GitHub Releases](https://github.com/{user}/recipe/releases/latest)
  2. Download `recipe-windows-amd64.exe` (Intel/AMD 64-bit) or `recipe-windows-arm64.exe` (ARM 64-bit - Surface Pro X)
  3. (Optional) Add to PATH:
     - Move `recipe-windows-amd64.exe` to `C:\Program Files\recipe\`
     - Add `C:\Program Files\recipe\` to system PATH
  4. Verify installation:
     ```powershell
     recipe.exe --version
     ```

  ### Verify Download Integrity (Optional)

  Each binary includes a SHA256 checksum file for verification:

  ```bash
  # Download checksum file
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-linux-amd64.sha256

  # Verify checksum
  sha256sum -c recipe-linux-amd64.sha256
  # Output: recipe-linux-amd64: OK
  ```

  ### Build from Source

  **Requirements:** Go 1.24+

  ```bash
  git clone https://github.com/{user}/recipe.git
  cd recipe
  go build -o recipe cmd/cli/main.go
  ```

  ### Usage

  ```bash
  # Convert single file
  recipe convert input.np3 output.xmp

  # Batch convert directory
  recipe convert --batch input_dir/ output_dir/

  # View help
  recipe --help
  ```

  ---

  ## Web Interface

  Recipe is also available as a web application with zero installation:

  **Live Web App:** https://recipe.pages.dev

  - **100% Client-Side:** Your files never leave your device
  - **Drag-and-Drop:** Upload presets, select target format, download converted file
  - **Cross-Platform:** Works in any modern browser (Chrome, Firefox, Safari, Edge)
  ```

- [ ] **Commit README Update:**
  ```bash
  git add README.md
  git commit -m "docs: Add CLI installation instructions to README"
  git push origin main
  ```

**Validation:**
- README.md includes installation section
- Platform-specific instructions provided (Linux, macOS, Windows)
- Latest release link functional
- Checksum verification documented
- Build from source instructions included
- Web interface URL documented

---

### Task 5: Document Versioning Strategy in CHANGELOG.md (AC-5)

- [ ] **Add Versioning Strategy Section to CHANGELOG.md:**

  Add after the header, before version entries:

  ```markdown
  # Changelog

  All notable changes to Recipe will be documented in this file.

  The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
  and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

  ## Versioning Strategy

  Recipe uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html): `vMAJOR.MINOR.PATCH`

  - **MAJOR:** Breaking changes (format incompatibility, API changes, CLI flag changes)
    - Example: v1.0.0 → v2.0.0 (NP3 format structure changed, old files incompatible)
  - **MINOR:** New features (backward compatible, new format support, new CLI commands)
    - Example: v1.0.0 → v1.1.0 (Added DNG format support, existing functionality unchanged)
  - **PATCH:** Bug fixes (no new features, backward compatible)
    - Example: v1.1.0 → v1.1.1 (Fixed WASM conversion bug for specific XMP files)

  ### Pre-Release Versions

  - **v0.x.y:** Beta/experimental releases (breaking changes allowed between minor versions)
  - **v1.0.0:** First stable release (API stability commitment begins)

  ### Release Process

  1. Update CHANGELOG.md with version entry
  2. Commit: `git commit -m "chore: prepare vX.Y.Z release"`
  3. Create tag: `git tag vX.Y.Z`
  4. Push tag: `git push origin vX.Y.Z`
  5. GitHub Actions builds and publishes binaries automatically

  ---

  ## [Unreleased]
  ### Added
  ### Changed
  ### Fixed
  ### Removed

  ## [0.1.0] - 2025-11-06
  ...
  ```

- [ ] **Commit CHANGELOG Update:**
  ```bash
  git add CHANGELOG.md
  git commit -m "docs: Add versioning strategy to CHANGELOG"
  git push origin main
  ```

**Validation:**
- CHANGELOG.md includes versioning strategy section
- Semantic versioning explained (MAJOR.MINOR.PATCH)
- Examples provided for each version type
- Pre-release versioning documented
- Release process documented

---

### Task 6: Add Release Notes Template (AC-4)

- [ ] **Create Release Notes Template:**

  GitHub Actions can auto-generate release notes from commits. To enhance with CHANGELOG excerpt, add workflow step:

  Update `.github/workflows/release.yml` to extract CHANGELOG:

  ```yaml
  jobs:
    build:
      runs-on: ubuntu-latest
      timeout-minutes: 15

      strategy:
        matrix:
          os: [linux, darwin, windows]
          arch: [amd64, arm64]

      steps:
        - name: Checkout repository
          uses: actions/checkout@v4

        - name: Setup Go
          uses: actions/setup-go@v5
          with:
            go-version: '1.24'

        - name: Build binary
          run: |
            BINARY_NAME=recipe-${{ matrix.os }}-${{ matrix.arch }}
            if [ "${{ matrix.os }}" = "windows" ]; then
              BINARY_NAME="${BINARY_NAME}.exe"
            fi
            GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -ldflags="-s -w" -o $BINARY_NAME cmd/cli/main.go

        - name: Generate checksums
          run: |
            BINARY_NAME=recipe-${{ matrix.os }}-${{ matrix.arch }}
            if [ "${{ matrix.os }}" = "windows" ]; then
              BINARY_NAME="${BINARY_NAME}.exe"
            fi
            sha256sum $BINARY_NAME > $BINARY_NAME.sha256

        - name: Upload binaries to release
          uses: softprops/action-gh-release@v1
          with:
            files: |
              recipe-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.os == 'windows' && '.exe' || '' }}
              recipe-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.os == 'windows' && '.exe' || '' }}.sha256
            body: |
              ## Recipe ${{ github.ref_name }}

              Pre-built binaries for Linux, macOS, and Windows.

              ### Installation

              Download the binary for your platform and architecture:
              - `recipe-linux-amd64` - Linux (Intel/AMD 64-bit)
              - `recipe-linux-arm64` - Linux (ARM 64-bit)
              - `recipe-darwin-amd64` - macOS (Intel)
              - `recipe-darwin-arm64` - macOS (Apple Silicon)
              - `recipe-windows-amd64.exe` - Windows (Intel/AMD 64-bit)
              - `recipe-windows-arm64.exe` - Windows (ARM 64-bit)

              ### Verify Download

              ```bash
              sha256sum -c recipe-{platform}-{arch}.sha256
              ```

              ### Changelog

              See [CHANGELOG.md](https://github.com/${{ github.repository }}/blob/main/CHANGELOG.md) for full details.
          env:
            GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  ```

- [ ] **Commit Workflow Update:**
  ```bash
  git add .github/workflows/release.yml
  git commit -m "feat(release): Add release notes template with installation instructions"
  git push origin main
  ```

**Note:** For MVP, using auto-generated release notes is acceptable. Future enhancement: Extract CHANGELOG excerpt programmatically.

**Validation:**
- Workflow updated with release body template
- Release notes include installation instructions
- Checksum verification instructions included
- CHANGELOG.md link provided

---

### Task 7: Test Trigger Behavior (AC-1)

- [ ] **Test Non-Version Tag (Should NOT Trigger):**
  ```bash
  git tag test-tag
  git push origin test-tag
  ```
  - Navigate to: Repository → Actions
  - Verify workflow does NOT run (no new workflow run for test-tag)

- [ ] **Test Branch Push (Should NOT Trigger):**
  ```bash
  echo "test" >> README.md
  git add README.md
  git commit -m "test: trigger test"
  git push origin main
  ```
  - Navigate to: Repository → Actions
  - Verify release workflow does NOT run (only deploy-pages.yml runs)

- [ ] **Test Version Tag (Should Trigger):**
  ```bash
  git tag v0.0.2
  git push origin v0.0.2
  ```
  - Navigate to: Repository → Actions
  - Verify release workflow runs (triggered by v0.0.2 tag)

- [ ] **Cleanup Test Tags:**
  ```bash
  # Delete local tags
  git tag -d test-tag v0.0.1 v0.0.2

  # Delete remote tags
  git push origin --delete test-tag v0.0.1 v0.0.2

  # Delete GitHub Releases (via GitHub UI)
  # Navigate to: Repository → Releases
  # Click "..." on v0.0.1 and v0.0.2 → Delete
  ```

**Validation:**
- Non-version tag does NOT trigger workflow
- Branch push does NOT trigger workflow
- Version tag (v*) triggers workflow successfully
- Test tags cleaned up

---

### Task 8: Binary Validation Testing (AC-2)

**Goal:** Ensure binaries are executable and functional on target platforms.

- [ ] **Download All Binaries:**
  ```bash
  # From GitHub Release page or via curl
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-linux-amd64
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-linux-arm64
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-darwin-amd64
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-darwin-arm64
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-windows-amd64.exe
  curl -LO https://github.com/{user}/recipe/releases/latest/download/recipe-windows-arm64.exe
  ```

- [ ] **Test Linux amd64 Binary:**
  ```bash
  chmod +x recipe-linux-amd64
  ./recipe-linux-amd64 --version
  # Expected: recipe version v0.0.1

  ./recipe-linux-amd64 --help
  # Expected: CLI help text with commands

  # Test conversion (if sample files available)
  ./recipe-linux-amd64 convert testdata/np3/sample.np3 output.xmp
  ```

- [ ] **Test macOS arm64 Binary (Apple Silicon):**
  ```bash
  chmod +x recipe-darwin-arm64
  ./recipe-darwin-arm64 --version
  ./recipe-darwin-arm64 --help
  ```

- [ ] **Test Windows amd64 Binary:**
  ```powershell
  .\recipe-windows-amd64.exe --version
  .\recipe-windows-amd64.exe --help
  ```

- [ ] **Verify Binary Sizes:**
  ```bash
  ls -lh recipe-*
  # Expected: Each binary <50MB (typically 10-20MB after stripping debug symbols)
  ```

**Validation:**
- All binaries executable on target platforms
- `--version` outputs correct version (matches tag)
- `--help` displays CLI usage
- Conversion functionality works (if tested)
- Binary sizes reasonable (<50MB each)

---

### Task 9: Prepare First Production Release (AC-5, AC-8)

- [ ] **Review CHANGELOG.md:**
  - Verify all Epic 1-7 changes documented
  - Ensure version entry complete (all Added/Changed/Fixed sections)
  - Proofread for spelling/grammar errors

- [ ] **Finalize Version Number:**
  - Decision: v0.1.0 (beta/experimental) or v1.0.0 (stable)?
  - Recommendation: v0.1.0 (allows breaking changes in v0.2.0 if needed)

- [ ] **Update CHANGELOG.md for v0.1.0:**
  ```markdown
  ## [0.1.0] - 2025-11-06
  ### Added
  - Universal Recipe data model for format-agnostic parameter representation
  - NP3 binary parser and generator (Nik Collection presets)
  - XMP XML parser and generator (Adobe Lightroom presets)
  - lrtemplate Lua parser and generator (Lightroom templates)
  - Parameter mapping rules for bidirectional conversion between formats
  - Metadata field implementation (description, author, keywords)
  - Web interface with drag-and-drop file upload
  - CLI interface with convert and batch commands
  - Cloudflare Pages deployment automation
  - GitHub Releases setup for CLI binary distribution

  ### Performance
  - WASM conversion: <100ms average
  - Batch processing: 37ms for 100 files
  - 95%+ conversion accuracy across 1,501 test files
  ```

- [ ] **Commit Final CHANGELOG:**
  ```bash
  git add CHANGELOG.md
  git commit -m "chore: prepare v0.1.0 release"
  git push origin main
  ```

- [ ] **Create Production Tag:**
  ```bash
  git tag v0.1.0
  git push origin v0.1.0
  ```

- [ ] **Monitor Release Workflow:**
  - Navigate to: Repository → Actions → "Build and Release CLI Binaries"
  - Verify workflow runs successfully
  - Verify 6 binaries + 6 checksums uploaded

- [ ] **Verify GitHub Release:**
  - Navigate to: Repository → Releases → v0.1.0
  - Verify release notes include installation instructions
  - Verify all 12 assets present (6 binaries + 6 checksums)
  - Verify download URLs functional

**Validation:**
- CHANGELOG.md finalized for v0.1.0
- Production tag created and pushed
- GitHub Release created automatically
- All binaries available for download
- Release notes clear and comprehensive

---

### Task 10: Update sprint-status.yaml

- [ ] **Mark Story 7-6 as "drafted":**
  - Load `docs/sprint-status.yaml` completely
  - Find `7-6-github-releases-setup: backlog`
  - Change to: `7-6-github-releases-setup: drafted  # Story created: docs/stories/7-6-github-releases-setup.md (2025-11-06)`
  - Preserve all comments and structure

- [ ] **Commit Sprint Status Update:**
  ```bash
  git add docs/sprint-status.yaml
  git commit -m "chore: Mark story 7-6 (GitHub Releases setup) as drafted"
  git push origin main
  ```

**Validation:**
- sprint-status.yaml updated
- Story status changed from "backlog" to "drafted"
- No other lines modified
- Comments preserved

---

## Dev Notes

### Learnings from Previous Story

**From Story 7-5-cloudflare-pages-deployment (Status: drafted)**

Story 7-5 established GitHub Actions CI/CD for web deployment. Story 7-6 extends this pattern for CLI binary distribution.

**Shared Patterns:**
- **GitHub Actions Workflow:** Both use `.github/workflows/{name}.yml`
- **Trigger Mechanism:** 7-5 triggers on branch push, 7-6 triggers on tag push
- **Go Toolchain:** Both use `actions/setup-go@v5` with Go 1.24
- **Build Optimization:** Both use `-ldflags="-s -w"` to strip debug symbols
- **Timeout Safety:** Both set workflow timeout (7-5: 10 min, 7-6: 15 min)

**Key Differences:**
- **7-5 Output:** Single WASM binary (web/recipe.wasm)
- **7-6 Output:** 6 CLI binaries (multi-platform cross-compilation)
- **7-5 Deploy:** Cloudflare Pages (static site hosting)
- **7-6 Deploy:** GitHub Releases (binary artifact hosting)
- **7-5 Timing:** <5 minutes target
- **7-6 Timing:** <10 minutes target (more binaries to build)

**Workflow Similarities:**
```
Developer action → GitHub Actions triggered
    ↓
Checkout repository
    ↓
Setup Go 1.24
    ↓
Build Go code (WASM or CLI)
    ↓
Deploy to hosting platform (Cloudflare or GitHub Releases)
    ↓
Verify deployment success (commit status, notification)
```

**Reuse from Story 7-5:**
- README.md structure (add installation section similar to deployment section)
- Workflow patterns (trigger, setup, build, deploy, verify)
- Documentation style (code examples, validation steps)

[Source: stories/7-5-cloudflare-pages-deployment.md#Dev-Notes]

---

### Architecture Alignment

**Follows Tech Spec Epic 7:**
- GitHub Releases setup satisfies NFR-7.2 (all 8 ACs)
- Implements automated CLI binary distribution
- Completes deployment architecture for both web and CLI

**Epic 7 Deployment Strategy:**
```
Recipe Deployment Architecture:

CLI Interface (Epic 3)
    ↓
Multi-Platform Build (GitHub Actions) ← YOU ARE HERE (Story 7-6)
    ↓
GitHub Releases Deployment
    ↓
Binaries Available: https://github.com/{user}/recipe/releases
```

**From PRD (Section: NFR-7.2 CLI Distribution):**
> NFR-7.2: CLI binaries available for download from GitHub Releases
> - Platforms: Linux, macOS, Windows
> - Architectures: amd64, arm64
> - Automated build on version tag push

Story 7-6 implements this requirement with:
- GitHub Actions workflow (`.github/workflows/release.yml`)
- Build matrix (3 OS × 2 arch = 6 binaries)
- Cross-compilation (`GOOS` and `GOARCH` environment variables)
- Semantic versioning (vMAJOR.MINOR.PATCH tags)
- Checksum generation (SHA256 for download verification)

**From Architecture (Section: Deployment Architecture):**
> GitHub Actions CI/CD:
> - Implements build matrix for CLI distribution (os: [linux, darwin, windows], arch: [amd64, arm64])
> - Produces release artifacts as specified in Architecture: recipe-{os}-{arch} binaries
> - Integrates with GitHub Releases for zero-cost artifact hosting

Story 7-6 implements:
- Build matrix with 6 platform/architecture combinations
- Binary naming convention: `recipe-{os}-{arch}` (matches Architecture spec)
- GitHub Releases integration via `softprops/action-gh-release@v1`
- Automated checksum generation for security

**Zero-Cost Infrastructure:**
- GitHub Actions: Free tier (2,000 minutes/month for public repos, unlimited for ubuntu runners)
- GitHub Releases: Free tier (unlimited releases, 2GB file size limit per release)
- Total cost: $0/month for MVP

**Build Flow:**
```
Developer creates version tag (v0.1.0)
    ↓
git tag v0.1.0
git push origin v0.1.0
    ↓
GitHub Actions: release.yml triggered
    ↓
Build matrix: 6 jobs (3 OS × 2 arch)
    ↓
Parallel builds:
  - recipe-linux-amd64
  - recipe-linux-arm64
  - recipe-darwin-amd64
  - recipe-darwin-arm64
  - recipe-windows-amd64.exe
  - recipe-windows-arm64.exe
    ↓
Generate SHA256 checksums (6 files)
    ↓
Upload binaries + checksums to GitHub Release (12 assets total)
    ↓
Release published: https://github.com/{user}/recipe/releases/tag/v0.1.0
    ↓
Users can download binaries (5-10 minutes total from tag push to release)
```

**Platform Coverage:**
- **Linux amd64:** Most common server/desktop platform (Intel/AMD processors)
- **Linux arm64:** Raspberry Pi, AWS Graviton instances, ARM servers
- **macOS amd64:** Intel Mac (pre-2020 and some current models)
- **macOS arm64:** Apple Silicon (M1/M2/M3 chips, 2020+)
- **Windows amd64:** Standard Windows desktop/laptop (Intel/AMD)
- **Windows arm64:** Surface Pro X, Windows on ARM devices

**Total Platform Coverage:** 99%+ of desktop/server environments

---

### Dependencies

**Internal Dependencies:**
- Story 7-5 (Cloudflare Pages Deployment) - Establishes GitHub Actions pattern (COMPLETED - drafted)
- Epic 3 (CLI Interface) - Provides `cmd/cli/main.go` for binary builds (COMPLETED - 3-1 through 3-5 done, 3-6 ready-for-dev)
- Epic 1 (Conversion Engine) - Provides conversion logic for CLI (COMPLETED)

**External Dependencies:**
- GitHub repository (already exists)
- Go 1.24 toolchain (GitHub Actions provides via actions/setup-go@v5)
- `softprops/action-gh-release@v1` GitHub Action (public, stable)

**Blockers:**
- None - All prerequisites exist (cmd/cli/main.go, repository, GitHub Actions)

**Version Tag Dependency:**
- No release until user creates first version tag (v0.0.1 or v0.1.0)
- Workflow tested with test tags (v0.0.1), production tag (v0.1.0) created in Task 9

---

### Testing Strategy

**Manual Testing (Primary Method):**
- **Workflow Trigger:** Create version tag → Verify workflow runs
- **Build Matrix:** Check logs → Verify 6 jobs execute in parallel
- **Binary Creation:** Check logs → Verify binaries built successfully
- **Upload Success:** Check logs → Verify binaries uploaded to release
- **Release Creation:** Navigate to Releases → Verify release exists
- **Download Test:** Download binaries → Verify executables work
- **Checksum Verification:** Verify SHA256 checksums match
- **Timing:** Measure workflow duration → Verify <10 minutes

**GitHub Actions Testing:**
- **Trigger Test:** Push non-version tag → Verify workflow does NOT run
- **Tag Pattern Test:** Push v* tag → Verify workflow runs
- **Build Failure Test:** Break cmd/cli/main.go → Verify workflow fails gracefully
- **Timeout Test:** Introduce infinite loop → Verify workflow times out at 15 minutes

**Platform Testing:**
- **Linux amd64:** Download binary → Run `./recipe-linux-amd64 --version`
- **macOS arm64:** Download binary → Run `./recipe-darwin-arm64 --version`
- **Windows amd64:** Download binary → Run `recipe-windows-amd64.exe --version`
- **Cross-Platform:** Test on each platform (or Docker containers for Linux)

**Checksum Verification:**
- Download binary + checksum file
- Run `sha256sum -c recipe-linux-amd64.sha256`
- Verify output: "recipe-linux-amd64: OK"

**Acceptance:**
- All 8 ACs verified (trigger, build matrix, upload, changelog, versioning, timing, README, CHANGELOG)
- Workflow triggers on version tag push
- 6 binaries built and uploaded
- GitHub Release created with all assets
- Binaries executable on target platforms
- Checksums verify successfully
- Build completes in <10 minutes
- README.md and CHANGELOG.md updated

---

### Technical Debt / Future Enhancements

**Deferred to Post-MVP:**
- **GPG Signature Verification:** Sign binaries with GPG key for enhanced security
- **Package Manager Distribution:** Submit to Homebrew (macOS), Scoop/Chocolatey (Windows), apt/yum (Linux)
- **Automatic CHANGELOG Extraction:** Programmatically extract version section from CHANGELOG.md for release notes
- **Build Caching:** Cache Go modules between workflow runs (faster builds)
- **Binary Compression:** Use UPX or similar to reduce binary size further
- **Release Draft Mode:** Create release as draft, manually publish after review
- **Notarization:** Notarize macOS binaries (required for macOS Gatekeeper on future OS versions)

**GPG Signature Verification (Future Enhancement):**
Users can verify binary authenticity with GPG signatures.

**Configuration:**
1. Generate GPG key for Recipe project
2. Add signing step to workflow:
   ```yaml
   - name: Sign binary
     run: gpg --detach-sign --armor recipe-linux-amd64
   ```
3. Upload `.asc` signature files to release
4. Document verification in README:
   ```bash
   gpg --verify recipe-linux-amd64.asc recipe-linux-amd64
   ```

**Benefits:**
- Proves binaries published by trusted source
- Detects tampering with downloaded files
- Industry standard for open-source distribution

**Tradeoff:**
- Requires managing GPG keys securely
- Users must import public key
- Additional complexity for casual users

**Package Manager Distribution (Future Enhancement):**
Distribute via platform-specific package managers for easier installation.

**Homebrew (macOS/Linux):**
1. Create Homebrew formula (Ruby script)
2. Submit to homebrew-core or create tap (custom repository)
3. Users install: `brew install recipe`

**Scoop/Chocolatey (Windows):**
1. Create manifest file (JSON for Scoop, XML for Chocolatey)
2. Submit to package repository
3. Users install: `scoop install recipe` or `choco install recipe`

**Benefits:**
- Automatic updates (package manager handles versioning)
- Simplified installation (one command)
- Integration with OS package ecosystem

**Tradeoff:**
- Requires maintaining separate package manifests
- Submission/approval process for official repositories
- Version lag (package updated after release)

**Recommendation:** Defer to post-MVP, focus on GitHub Releases for initial launch.

**Binary Compression (Future Enhancement):**
Use UPX (Ultimate Packer for eXecutables) to reduce binary size.

**Implementation:**
```yaml
- name: Compress binary
  run: upx --best --lzma recipe-linux-amd64
```

**Benefits:**
- 50-70% size reduction (10MB → 3-5MB)
- Faster downloads for users
- Reduced GitHub storage usage

**Tradeoff:**
- Slower startup time (decompress on execution)
- Some antivirus software flags UPX-compressed binaries as suspicious
- Not supported on all platforms (macOS arm64 requires special handling)

**Recommendation:** Test in post-MVP, evaluate startup time impact.

---

### References

- [Source: docs/tech-spec-epic-7.md#NFR-7.2] - GitHub Releases setup requirements (8 ACs)
- [Source: docs/PRD.md#NFR-7.2] - CLI distribution requirements (GitHub Releases, multi-platform)
- [Source: docs/architecture.md#Deployment-Architecture] - GitHub Actions CI/CD design
- [Source: Keep a Changelog] - https://keepachangelog.com/en/1.0.0/
- [Source: Semantic Versioning] - https://semver.org/spec/v2.0.0.html
- [Source: GitHub Actions Documentation] - https://docs.github.com/en/actions
- [Source: softprops/action-gh-release@v1] - https://github.com/softprops/action-gh-release

**GitHub Actions Features:**
- Free tier: 2,000 minutes/month for public repos, unlimited on ubuntu runners
- Build matrix: Parallel jobs for multi-platform builds
- Secrets: Auto-provided `GITHUB_TOKEN` for release uploads (no manual configuration)
- Artifact storage: Release assets stored indefinitely (no expiration)

**Cross-Compilation in Go:**
- `GOOS`: Target operating system (linux, darwin, windows)
- `GOARCH`: Target architecture (amd64, arm64, 386, arm, etc.)
- Example: `GOOS=linux GOARCH=arm64 go build` produces Linux ARM binary
- No cross-compiler needed (Go toolchain handles all platforms)

**Semantic Versioning Best Practices:**
- Start with v0.1.0 for initial release (signals beta/experimental)
- Increment MAJOR version for breaking changes (v1.0.0 → v2.0.0)
- Increment MINOR version for new features (v1.0.0 → v1.1.0)
- Increment PATCH version for bug fixes (v1.1.0 → v1.1.1)
- Pre-release versions: v1.0.0-alpha, v1.0.0-beta, v1.0.0-rc1

---

### Known Issues / Blockers

**None** - This story has no technical blockers. All required infrastructure exists:
- GitHub repository exists
- CLI interface code exists (`cmd/cli/main.go` from Epic 3)
- GitHub Actions available (free tier for public repos)
- Conversion engine stable (Epic 1 completed, tested with 1,501 sample files)

**Epic 3 CLI Completion:**
- Stories 3-1 through 3-5 completed (CLI functional)
- Story 3-6 (JSON output mode) ready-for-dev but not required for binary builds
- Binary builds tested locally (Go cross-compilation verified)

**GitHub Token Permissions:**
- `GITHUB_TOKEN` auto-provided by GitHub Actions (no manual configuration)
- Token has sufficient permissions for release creation and asset uploads
- No additional secrets required (unlike Story 7-5 which needed Cloudflare secrets)

**Workflow Timeout:**
- Default timeout: 360 minutes (6 hours)
- Configured timeout: 15 minutes (safety net)
- Rationale: Prevent runaway builds, typical build time <10 minutes

**Binary Size:**
- Current CLI binary size: ~10-20MB (after `-ldflags="-s -w"`)
- Well within GitHub Releases 2GB per-file limit
- Acceptable download size for users

---

### Cross-Story Coordination

**Dependencies:**
- Story 7-5 (Cloudflare Pages Deployment) - Establishes GitHub Actions pattern
- Epic 3 (CLI Interface) - Provides CLI source code for binary builds
- Epic 1 (Conversion Engine) - Provides conversion logic

**Enables:**
- Public distribution of Recipe CLI (users can download without building)
- Version management (semantic versioning, release history)
- Professional project image (official releases, checksums, documentation)

**Completes Epic 7:**
Story 7-6 is the final implementation story in Epic 7. After completion:
- Documentation complete (Stories 7-1 through 7-4)
- Web deployment automated (Story 7-5)
- CLI distribution automated (Story 7-6)
- Recipe fully deployable and distributable

**Architectural Consistency:**
GitHub Releases completes Recipe's deployment architecture:
- **Web Interface:** Cloudflare Pages (Story 7-5, zero-cost static hosting)
- **CLI Binaries:** GitHub Releases (Story 7-6, zero-cost artifact hosting)
- **Zero Cost:** Both use free tiers (Cloudflare Pages + GitHub Actions + GitHub Releases)
- **Automated:** Both triggered by git events (push to main, tag push)
- **Global Distribution:** Cloudflare CDN (web), GitHub CDN (binaries)

---

### Project Structure Notes

**New Files Created:**
```
.github/workflows/
├── deploy-pages.yml   # Cloudflare Pages deployment (Story 7-5, existing)
├── release.yml        # GitHub Releases workflow (NEW)

CHANGELOG.md           # Version history (NEW)

docs/stories/
├── 7-6-github-releases-setup.md   # This story document (NEW)
```

**Modified Files:**
```
README.md              # Add installation section (MODIFIED)
docs/sprint-status.yaml   # Mark 7-6 as "drafted" (MODIFIED)
```

**No Structural Changes:** This story adds CI/CD automation and documentation. No changes to source code structure.

**Workflow Location:**
- GitHub Actions workflows: `.github/workflows/` (standard location)
- Workflow naming: `release.yml` (descriptive, action-oriented)
- Future workflows: `test.yml` (optional CI testing)

**CHANGELOG Location:**
- Repository root (same directory as README.md, go.mod)
- Standard location for Keep a Changelog format
- Easily accessible to users and automated tools

---

## Dev Agent Record

**Story Context**: See `docs/stories/7-6-github-releases-setup.context.xml` for complete implementation context including documentation artifacts, code integration points, interfaces, constraints, and testing standards. Generated 2025-11-06.

### Context Reference

- `docs/stories/7-6-github-releases-setup.context.xml` - Story context with documentation artifacts (6), code references (6), interfaces (6), constraints (10), and comprehensive test ideas for all 8 acceptance criteria. Generated 2025-11-06.

### Agent Model Used

<!-- To be filled by dev agent -->

### Debug Log References

<!-- Dev agent will document:
- CHANGELOG.md creation and initial content
- GitHub Actions workflow file creation (.github/workflows/release.yml)
- Test tag creation and workflow trigger verification
- Build matrix execution monitoring (6 parallel jobs)
- Binary download and verification testing
- Checksum verification process
- README.md installation section update
- Versioning strategy documentation in CHANGELOG.md
- Production release preparation (v0.1.0)
- Workflow timing measurements
-->

### Completion Notes List

<!-- Dev agent will document:
- CHANGELOG.md created with Keep a Changelog format
- Versioning strategy documented (semantic versioning, MAJOR.MINOR.PATCH)
- GitHub Actions workflow created (.github/workflows/release.yml)
- Build matrix configured (3 OS × 2 arch = 6 binaries)
- Test tags created and workflow triggered successfully (v0.0.1, v0.0.2)
- 6 binaries built in parallel (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64, windows-arm64)
- SHA256 checksums generated for each binary
- Binaries uploaded to GitHub Release (12 assets total: 6 binaries + 6 checksums)
- Workflow completed in [X] minutes (target: <10 minutes)
- Binaries tested on multiple platforms (Linux, macOS, Windows)
- Checksum verification successful (sha256sum -c passed)
- README.md updated with installation instructions (platform-specific)
- Production release v0.1.0 created with complete CHANGELOG
- Release notes include installation instructions and checksum verification
- sprint-status.yaml updated (7-6: backlog → drafted)
-->

### File List

<!-- Dev agent will document files created/modified/deleted:
**NEW:**
- `CHANGELOG.md` - Version history with Keep a Changelog format
- `.github/workflows/release.yml` - GitHub Releases workflow for CLI binaries
- `docs/stories/7-6-github-releases-setup.md` - Story document

**MODIFIED:**
- `README.md` - Added installation section with platform-specific instructions
- `docs/sprint-status.yaml` - Updated 7-6 from "backlog" to "drafted"

**DELETED:**
- (none)
-->

---

## Change Log

- **2025-11-06:** Story created from Epic 7 Tech Spec (Sixth story in Epic 7, final implementation story, implements automated GitHub Releases for CLI binary distribution with multi-platform support)
