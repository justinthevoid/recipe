# Epic Technical Specification: Documentation & Deployment

Date: November 6, 2025
Author: Justin
Epic ID: 7
Status: Draft

---

## Overview

Epic 7 establishes Recipe's public presence through comprehensive documentation and automated deployment infrastructure. This epic delivers the landing page, user documentation, legal disclosures, and CI/CD pipelines required to deploy Recipe to Cloudflare Pages (web) and GitHub Releases (CLI).

The epic transforms Recipe from a functional tool into a publicly accessible, legally compliant, and professionally documented product. Key deliverables include a clear landing page explaining Recipe's purpose and privacy promise, comprehensive format compatibility documentation, FAQ addressing legal and technical questions, automated Cloudflare Pages deployment from the main branch, and GitHub Actions workflow for multi-platform CLI binary releases.

## Objectives and Scope

### In Scope

**Documentation (FR-7.1 - FR-7.4):**
- Landing page with project description, 3-step usage guide, and privacy promise
- Format compatibility matrix showing parameter mapping across NP3/XMP/lrtemplate
- FAQ section answering common legal, privacy, and technical questions
- Legal disclaimer covering reverse engineering disclosure and no-warranty statement

**Deployment Infrastructure (NFR-7.1 - NFR-7.3):**
- Cloudflare Pages deployment with automatic builds on main branch push
- GitHub Actions workflow for multi-platform CLI binary releases (Linux/macOS/Windows, amd64/arm64)
- Semantic versioning strategy (vMAJOR.MINOR.PATCH)
- CHANGELOG.md maintenance process

**Content Requirements:**
- Non-technical user-friendly language for landing page and FAQ
- Clear visual format compatibility table (scannable, includes approximations)
- Legally reviewed disclaimer (if possible) or conservative "use at own risk" language
- README.md updates reflecting deployment URLs and installation instructions

### Out of Scope

- **User analytics/tracking:** Privacy-first approach, no tracking (covered in FR-2.9)
- **Internationalization:** English-only documentation for MVP
- **Package manager submissions:** Homebrew/Scoop/Chocolatey distribution deferred to post-MVP
- **Custom domain setup:** Using default Cloudflare Pages domain (recipe.pages.dev)
- **Advanced CI/CD features:** No A/B testing, canary deployments, or multi-stage pipelines
- **User support forum/ticketing:** GitHub Issues only
- **Marketing materials:** No press releases, blog posts, or promotional content

## System Architecture Alignment

Epic 7 completes the deployment architecture defined in the Architecture document (Section: "Deployment Architecture").

**Cloudflare Pages Integration:**
- Aligns with Architecture decision to use Cloudflare Pages for static hosting (zero cost, global CDN, auto HTTPS)
- Deploys `web/` directory containing index.html, main.js, style.css, recipe.wasm (built in Epic 2)
- Leverages automatic gzip compression (WASM reduced 70% per Architecture spec)
- Provides global CDN with sub-100ms latency worldwide

**GitHub Actions CI/CD:**
- Implements build matrix for CLI distribution (os: [linux, darwin, windows], arch: [amd64, arm64])
- Produces release artifacts as specified in Architecture: recipe-{os}-{arch} binaries
- Automates WASM build step: `GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go`
- Integrates with GitHub Releases for zero-cost artifact hosting

**Technology Stack Consistency:**
- Static HTML/CSS/JS (no build framework) - aligns with Epic 2's vanilla JavaScript decision
- Markdown-based documentation (README.md, CHANGELOG.md) - aligns with Architecture documentation strategy
- Go 1.24+ toolchain for WASM and CLI builds - consistent with Epic 1-6

**Constraints Met:**
- Zero-cost deployment infrastructure (Cloudflare Pages free tier, GitHub Actions free tier)
- No backend services (static site only)
- Privacy-preserving (no analytics, no user data collection)
- Single-binary CLI distribution (no dependencies)

## Detailed Design

### Services and Modules

Epic 7 has no runtime services - it consists of documentation artifacts and CI/CD configuration files.

| Module                        | Responsibility                                | Inputs                                                   | Outputs                                        | Owner                 |
| ----------------------------- | --------------------------------------------- | -------------------------------------------------------- | ---------------------------------------------- | --------------------- |
| **Landing Page**              | User-facing documentation explaining Recipe   | PRD requirements, Epic 1-2 features                      | `web/index.html` (enhanced)                    | Story 7-1             |
| **Format Matrix**             | Document parameter mapping compatibility      | Epic 1 mapping rules, PRD compatibility notes            | Markdown table in landing page or separate doc | Story 7-2             |
| **FAQ Documentation**         | Answer legal, privacy, technical questions    | PRD NFRs, legal research, user feedback                  | Markdown section in landing page or FAQ.md     | Story 7-3             |
| **Legal Disclaimer**          | Reverse engineering disclosure, no-warranty   | Legal review (if available), conservative template       | Markdown section in landing page               | Story 7-4             |
| **Cloudflare Pages Workflow** | Automate web deployment on main branch push   | `.github/workflows/deploy-pages.yml`, Cloudflare secrets | Deployed site at recipe.pages.dev              | Story 7-5             |
| **GitHub Releases Workflow**  | Build multi-platform CLI binaries on tag push | `.github/workflows/release.yml`, Go build matrix         | Release artifacts (6 binaries)                 | Story 7-6             |
| **README.md**                 | Repository overview, installation guide       | Epic 1-6 features, deployment URLs                       | Updated README.md                              | Stories 7-1, 7-5, 7-6 |
| **CHANGELOG.md**              | Version history tracking                      | Git commits, semantic versioning                         | CHANGELOG.md                                   | Story 7-6             |

### Data Models and Contracts

**No runtime data models** - Epic 7 deals with documentation and configuration files.

**Configuration Schemas:**

**1. Cloudflare Pages Workflow (`.github/workflows/deploy-pages.yml`):**
```yaml
name: Deploy to Cloudflare Pages
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build WASM
        run: GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/recipe.wasm cmd/wasm/main.go
      - name: Deploy to Cloudflare Pages
        uses: cloudflare/pages-action@v1
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          projectName: recipe
          directory: web
```

**2. GitHub Releases Workflow (`.github/workflows/release.yml`):**
```yaml
name: Build Release Binaries
on:
  push:
    tags: ['v*']
jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build
        run: GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -o recipe-${{ matrix.os }}-${{ matrix.arch }}${{ matrix.os == 'windows' && '.exe' || '' }} cmd/cli/main.go
      - name: Upload to GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: recipe-*
```

**3. Semantic Version Tag Format:**
- Pattern: `v{MAJOR}.{MINOR}.{PATCH}` (e.g., v1.0.0, v1.1.0, v1.1.1)
- MAJOR: Breaking changes (format incompatibility)
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

**4. CHANGELOG.md Format (Keep a Changelog standard):**
```markdown
# Changelog
All notable changes to Recipe will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Fixed

## [1.0.0] - YYYY-MM-DD
### Added
- Initial release
```

### APIs and Interfaces

**No programmatic APIs** - Epic 7 provides user-facing documentation and CI/CD automation.

**GitHub Actions Secrets (Story 7-5):**
- `CLOUDFLARE_API_TOKEN`: Cloudflare API token with Pages write permission
- `CLOUDFLARE_ACCOUNT_ID`: Cloudflare account ID
- `GITHUB_TOKEN`: Auto-provided by GitHub Actions for release uploads

**Cloudflare Pages API (implicit via cloudflare/pages-action@v1):**
- POST deployment to Cloudflare Pages
- Inputs: `directory` (web/), `projectName` (recipe), `apiToken`, `accountId`
- Outputs: Deployment URL (https://recipe.pages.dev), preview URLs for PRs

**GitHub Releases API (implicit via softprops/action-gh-release@v1):**
- POST release creation on tag push
- Inputs: `tag_name` (v*), `files` (recipe-* binaries)
- Outputs: Release URL, download URLs for each binary

**User-Facing URLs:**
- Production web app: `https://recipe.pages.dev` (or custom domain if configured later)
- Latest CLI release: `https://github.com/{user}/recipe/releases/latest`
- Specific binary download: `https://github.com/{user}/recipe/releases/latest/download/recipe-{os}-{arch}`

### Workflows and Sequencing

**Workflow 1: Cloudflare Pages Deployment (Story 7-5)**

```
Developer pushes to main branch
    ↓
GitHub Actions: deploy-pages.yml triggered
    ↓
Checkout repository
    ↓
Setup Go 1.24
    ↓
Build WASM: go build -o web/recipe.wasm cmd/wasm/main.go
    ↓
Cloudflare Pages Action: Deploy web/ directory
    ↓
Cloudflare Pages: Build and deploy to CDN
    ↓
Live at https://recipe.pages.dev (2-5 minutes total)
    ↓
Notify GitHub commit status: ✓ Deployed
```

**Workflow 2: GitHub Release Creation (Story 7-6)**

```
Developer creates semantic version tag (v1.0.0)
    ↓
git tag v1.0.0
git push origin v1.0.0
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
Upload all binaries to GitHub Release
    ↓
Release published: https://github.com/{user}/recipe/releases/tag/v1.0.0
    ↓
Users can download binaries (5-10 minutes total)
```

**Workflow 3: Documentation Update (Stories 7-1 to 7-4)**

```
SM creates documentation content (landing page, FAQ, disclaimer)
    ↓
Update web/index.html with new sections
    ↓
Update README.md with deployment URLs and installation instructions
    ↓
Commit and push to main branch
    ↓
Triggers Cloudflare Pages deployment (Workflow 1)
    ↓
Documentation live on recipe.pages.dev
```

**Workflow 4: Version Release Process (Story 7-6)**

```
Team decides to release new version
    ↓
Update CHANGELOG.md with version details
    ↓
Commit changes: git commit -m "chore: prepare v1.1.0 release"
    ↓
Create tag: git tag v1.1.0
    ↓
Push tag: git push origin v1.1.0
    ↓
Triggers GitHub Release workflow (Workflow 2)
    ↓
Binaries built and released
    ↓
Update README.md if installation instructions changed
```

## Non-Functional Requirements

### Performance

**NFR-7.1: Deployment Speed (per PRD NFR-7.1)**
- **Requirement:** Cloudflare Pages deployment completes in <5 minutes from push to live
- **Target Breakdown:**
  - GitHub Actions trigger: <30 seconds
  - WASM build: <2 minutes (Go compilation)
  - Cloudflare Pages deploy: <2 minutes (CDN propagation)
  - Total: <5 minutes
- **Validation:** Monitor GitHub Actions run times, add workflow timeout if exceeded
- **Measurement:** GitHub Actions duration from push to deployment success status

**NFR-7.2: Binary Build Speed**
- **Requirement:** GitHub Release workflow completes all 6 binaries in <10 minutes
- **Target Breakdown:**
  - Parallel builds (6 jobs): 3-5 minutes each
  - Upload to release: <1 minute
  - Total: <10 minutes
- **Validation:** GitHub Actions build matrix timing
- **Optimization:** Use Go build cache between runs

**NFR-7.3: Documentation Load Time**
- **Requirement:** Landing page loads in <2 seconds (first visit), <500ms (cached)
- **Target:**
  - HTML size: <50KB
  - Total page size (HTML + CSS + images): <200KB
  - Time to First Contentful Paint: <1.5 seconds
  - Cloudflare CDN latency: <100ms globally
- **Validation:** Lighthouse performance audit, WebPageTest
- **Constraint:** No JavaScript frameworks (vanilla HTML/CSS only for docs sections)

### Security

**NFR-7.4: Secrets Management**
- **Requirement:** GitHub Actions secrets securely stored, never logged
- **Implementation:**
  - `CLOUDFLARE_API_TOKEN`: GitHub repository secret (encrypted at rest)
  - `CLOUDFLARE_ACCOUNT_ID`: GitHub repository secret (not sensitive but consistent storage)
  - `GITHUB_TOKEN`: Auto-provided by GitHub Actions (automatic rotation)
- **Validation:** Audit workflow logs for secret exposure, use `::add-mask::` if needed
- **Access Control:** Repository admin only (Justin) can modify secrets

**NFR-7.5: Download Integrity**
- **Requirement:** CLI binaries verifiable via checksums (future: GPG signatures)
- **Implementation (MVP):**
  - Generate SHA256 checksums for each binary
  - Include checksums in release notes
  - Future: Add GPG signature verification
- **Validation:** Download binary, verify SHA256 matches release notes
- **User Command:** `sha256sum recipe-linux-amd64` (compare against GitHub Release)

**NFR-7.6: HTTPS Enforcement**
- **Requirement:** All traffic to recipe.pages.dev served over HTTPS
- **Implementation:** Cloudflare Pages automatic HTTPS (no configuration needed)
- **Validation:** HTTP requests auto-redirect to HTTPS, verify via curl
- **Certificate:** Cloudflare-managed TLS certificate (auto-renewal)

### Reliability/Availability

**NFR-7.7: Deployment Rollback**
- **Requirement:** Ability to rollback failed deployments within 5 minutes
- **Implementation:**
  - Cloudflare Pages: Rollback via dashboard to previous deployment
  - GitHub Releases: Delete tag, re-tag previous version
  - Manual process (no automation required for MVP)
- **Validation:** Test rollback scenario in development
- **Recovery Time Objective (RTO):** <5 minutes for web, <10 minutes for CLI

**NFR-7.8: Static Site Availability**
- **Requirement:** Web interface available 99.9%+ (inherent to Cloudflare Pages)
- **Target:** Leverage Cloudflare's global CDN (250+ data centers)
- **Monitoring:** Cloudflare Pages dashboard (uptime metrics)
- **No SLA Commitment:** Best-effort availability (free tier)
- **Degradation:** If Cloudflare down, users can still use CLI (offline-capable)

**NFR-7.9: Binary Availability**
- **Requirement:** GitHub Releases available 99.9%+ (inherent to GitHub infrastructure)
- **Fallback:** Users can build from source (`go build cmd/cli/main.go`)
- **Retention:** All releases preserved indefinitely (GitHub free tier policy)
- **Mirror:** No CDN mirror required (GitHub CDN sufficient)

### Observability

**NFR-7.10: Deployment Monitoring**
- **Requirement:** Visibility into deployment success/failure
- **Implementation:**
  - GitHub Actions: Built-in logs and status badges
  - Cloudflare Pages: Dashboard deployment history
  - Email notifications on workflow failure (GitHub default)
- **Metrics Tracked:**
  - Deployment duration (GitHub Actions timing)
  - Build success rate (workflow runs / successful runs)
  - WASM binary size (log in build output)
- **Alerting:** GitHub email notifications on failure (no custom alerting needed)

**NFR-7.11: Download Metrics**
- **Requirement:** Track CLI download counts (GitHub Insights)
- **Implementation:** GitHub Releases download statistics (automatic)
- **Privacy:** Aggregate counts only, no user tracking
- **Visibility:** Public via GitHub Release page

**NFR-7.12: Documentation Accessibility**
- **Requirement:** No tracking on recipe.pages.dev (privacy-first)
- **Implementation:** Zero analytics/tracking scripts (per FR-2.9)
- **Validation:** Audit web/ directory for tracking scripts, confirm none present
- **Exception:** None - strict no-tracking policy

## Dependencies and Integrations

### External Dependencies

**Go Toolchain:**
- **Version:** Go 1.24+ (current: 1.25.1 per go.mod)
- **Purpose:** WASM compilation, CLI binary builds
- **Integration Point:** GitHub Actions workflows (setup-go@v5)
- **Constraint:** Must support GOARCH=wasm for WebAssembly builds

**GitHub Actions:**
- **Actions Used:**
  - `actions/checkout@v4` - Repository checkout
  - `actions/setup-go@v5` - Go toolchain installation
  - `cloudflare/pages-action@v1` - Cloudflare Pages deployment
  - `softprops/action-gh-release@v1` - GitHub Release creation
- **Integration:** Triggered by push events (main branch, version tags)
- **Credentials:** GitHub-managed secrets (CLOUDFLARE_API_TOKEN, CLOUDFLARE_ACCOUNT_ID)

**Cloudflare Pages:**
- **Service:** Static site hosting and CDN
- **Plan:** Free tier (unlimited bandwidth, 500 builds/month)
- **Integration:** Via `cloudflare/pages-action@v1` GitHub Action
- **Configuration:** Project name: `recipe`, directory: `web/`
- **Output:** HTTPS endpoint at recipe.pages.dev

**GitHub Releases:**
- **Service:** Binary artifact hosting
- **Plan:** Free tier (unlimited releases, 2GB file size limit per release)
- **Integration:** Via `softprops/action-gh-release@v1` GitHub Action
- **Artifacts:** 6 binaries (Linux/macOS/Windows × amd64/arm64)

### Internal Dependencies

**Epic 1 (Core Conversion Engine):**
- **Dependency:** WASM binary built from Epic 1 conversion engine code
- **Files Used:** `cmd/wasm/main.go`, `internal/converter/*`, `internal/formats/*`, `internal/models/*`
- **Integration Point:** Story 7-5 builds WASM during deployment workflow
- **Constraint:** Epic 1 must be complete and stable

**Epic 2 (Web Interface):**
- **Dependency:** Web UI files (index.html, main.js, style.css)
- **Files Deployed:** `web/index.html`, `web/main.js`, `web/style.css`, `web/recipe.wasm`
- **Integration Point:** Story 7-1 enhances index.html with documentation sections
- **Constraint:** Epic 2 complete, web/ directory production-ready

**Epic 3 (CLI Interface):**
- **Dependency:** CLI source code for binary builds
- **Files Used:** `cmd/cli/main.go`, `cmd/cli/root.go`, `cmd/cli/convert.go`
- **Integration Point:** Story 7-6 builds CLI binaries from cmd/cli/main.go
- **Constraint:** Epic 3 stories 3-1 to 3-6 complete (CLI functional)

### Configuration Files

**Required Configuration:**
1. `.github/workflows/deploy-pages.yml` - Cloudflare Pages deployment (Story 7-5)
2. `.github/workflows/release.yml` - GitHub Release automation (Story 7-6)
3. `CHANGELOG.md` - Version history (Story 7-6)
4. Updated `README.md` - Installation instructions (Stories 7-1, 7-5, 7-6)

**GitHub Repository Secrets (Manual Setup):**
- `CLOUDFLARE_API_TOKEN` - Create at Cloudflare dashboard → API Tokens
- `CLOUDFLARE_ACCOUNT_ID` - Found in Cloudflare dashboard → Account ID

**No External APIs:**
- No runtime API calls from deployed application
- All integrations happen at build/deploy time via GitHub Actions

## Acceptance Criteria (Authoritative)

### FR-7.1: Landing Page (Story 7-1)

**AC-1:** Landing page contains clear project description ("What is Recipe?")
- Description explains Recipe as a photo preset converter
- Mentions support for NP3, XMP, and lrtemplate formats
- Written in non-technical language

**AC-2:** Landing page includes 3-step usage guide ("How to use")
- Step 1: Upload preset file
- Step 2: Select target format
- Step 3: Download converted file

**AC-3:** Privacy promise is prominently displayed
- States "100% client-side processing"
- Clarifies "Your files never leave your device"
- References WASM architecture for privacy

**AC-4:** Format compatibility matrix is visible and accessible
- Links to or embeds compatibility table
- Easy to navigate from landing page

**AC-5:** FAQ section is linked or embedded
- Accessible from landing page navigation
- Covers legal, privacy, and technical questions

**AC-6:** Page is readable by non-technical users
- No jargon without explanation
- Clear headings and structure
- Mobile-responsive layout (per Epic 2)

**AC-7:** Links to technical documentation are provided
- GitHub repository link
- README.md link
- Optional: Architecture/PRD links for developers

### FR-7.2: Format Compatibility Matrix (Story 7-2)

**AC-1:** Matrix shows parameter support across all 3 formats (NP3, XMP, lrtemplate)
- Table format with formats as columns, parameters as rows
- Clear visual indicators (✓, ✗, ~)

**AC-2:** Matrix is easy to scan
- Organized by parameter category (Basic, Tone Curve, Color, etc.)
- Sortable or filterable (optional for MVP)

**AC-3:** Approximations are clearly noted
- Symbol/indicator for approximated conversions (e.g., "~" or "⚠")
- Footnote explaining approximation meaning

**AC-4:** Unmappable features are documented
- Lists parameters that don't convert between specific formats
- Explains why (format limitation, no equivalent)

### FR-7.3: FAQ Documentation (Story 7-3)

**AC-1:** FAQ answers "Is this legal? (reverse engineering)"
- Explains reverse engineering context
- Cites fair use or research exemptions (if applicable)
- Recommends private use until legal assessment complete

**AC-2:** FAQ answers "Is my data private?"
- Confirms yes, WASM client-side processing
- No server uploads, no tracking
- References FR-2.9 privacy implementation

**AC-3:** FAQ answers "Why doesn't [feature] convert?"
- Explains format limitations
- References compatibility matrix
- Provides examples of unmappable parameters

**AC-4:** FAQ answers "How accurate is conversion?"
- States 95%+ accuracy target (from PRD)
- Explains approximations for unmappable parameters
- Encourages visual validation

**AC-5:** Answers are clear and concise
- 2-4 sentences per question
- Links to technical details for deep dives

**AC-6:** FAQ updated based on user feedback (post-launch)
- Process defined for adding new questions
- GitHub Issues monitored for common questions

### FR-7.4: Legal Disclaimer (Story 7-4)

**AC-1:** Disclaimer includes reverse engineering disclosure
- States "This tool uses reverse-engineered file formats"
- Clarifies no affiliation with original vendors

**AC-2:** No-warranty statement is present
- Standard "AS IS" disclaimer
- No guarantee of conversion accuracy or compatibility

**AC-3:** Recommends private use
- Suggests use for personal projects
- Notes legal assessment incomplete for commercial use

**AC-4:** Disclaimer is visible on landing page
- Placed in footer or prominent section
- Not hidden behind clicks (directly visible)

**AC-5:** Legally reviewed (if possible) or conservative template used
- Attempt to get legal review (optional)
- If not available, use conservative open-source template

### NFR-7.1: Cloudflare Pages Deployment (Story 7-5)

**AC-1:** GitHub Actions workflow triggers on push to main branch
- Workflow file: `.github/workflows/deploy-pages.yml`
- Triggers on push events to `main` only

**AC-2:** Workflow builds WASM binary
- Runs `GOOS=js GOARCH=wasm go build -o web/recipe.wasm cmd/wasm/main.go`
- Binary placed in `web/` directory

**AC-3:** Workflow deploys `web/` directory to Cloudflare Pages
- Uses `cloudflare/pages-action@v1`
- Deploys to project `recipe`

**AC-4:** Deployment completes in <5 minutes
- Measured from push to live site
- GitHub Actions workflow timeout set to 10 minutes

**AC-5:** Site is accessible at https://recipe.pages.dev
- HTTPS enforced automatically
- Cloudflare CDN enabled

**AC-6:** GitHub repository secrets configured
- `CLOUDFLARE_API_TOKEN` set
- `CLOUDFLARE_ACCOUNT_ID` set

**AC-7:** Deployment status visible in GitHub commit status
- Green checkmark on successful deployment
- Red X on failure with logs

### NFR-7.2: GitHub Releases Setup (Story 7-6)

**AC-1:** GitHub Actions workflow triggers on version tag push
- Workflow file: `.github/workflows/release.yml`
- Triggers on tags matching `v*` pattern

**AC-2:** Workflow builds 6 CLI binaries
- recipe-linux-amd64
- recipe-linux-arm64
- recipe-darwin-amd64
- recipe-darwin-arm64
- recipe-windows-amd64.exe
- recipe-windows-arm64.exe (optional if Epic 3 supports)

**AC-3:** Binaries are uploaded to GitHub Release
- Attached to release matching tag (e.g., v1.0.0)
- Download URLs functional

**AC-4:** Release includes CHANGELOG excerpt
- Automated from CHANGELOG.md or manual
- Shows changes for this version

**AC-5:** Semantic versioning followed
- Tags follow vMAJOR.MINOR.PATCH format
- CHANGELOG.md documents versioning strategy

**AC-6:** Build completes in <10 minutes
- Parallel matrix build
- Workflow timeout set appropriately

**AC-7:** README.md updated with installation instructions
- Includes download URLs for binaries
- Platform-specific installation steps

**AC-8:** CHANGELOG.md maintained
- Follows Keep a Changelog format
- Updated before each release

## Traceability Mapping

| AC ID                                    | Spec Section                                 | Component/File                     | Test Idea                                                                          |
| ---------------------------------------- | -------------------------------------------- | ---------------------------------- | ---------------------------------------------------------------------------------- |
| **FR-7.1: Landing Page**                 |
| AC-1                                     | Overview, Services/Landing Page              | web/index.html                     | Manual review: Verify project description present and clear                        |
| AC-2                                     | Workflows/Documentation Update               | web/index.html                     | Manual review: Count steps in usage guide (must be 3)                              |
| AC-3                                     | Overview, NFR Security                       | web/index.html                     | Manual review: Search for "client-side" and "never leave" text                     |
| AC-4                                     | Services/Format Matrix                       | web/index.html                     | Manual test: Click compatibility matrix link, verify accessible                    |
| AC-5                                     | Services/FAQ Documentation                   | web/index.html                     | Manual test: Click FAQ link, verify reachable                                      |
| AC-6                                     | NFR Performance/Documentation Load           | web/index.html                     | User testing: Non-technical user comprehension test                                |
| AC-7                                     | APIs/User-Facing URLs                        | web/index.html                     | Manual review: Verify GitHub link present and functional                           |
| **FR-7.2: Format Compatibility Matrix**  |
| AC-1                                     | Services/Format Matrix, Epic 1 mapping rules | web/index.html or docs             | Manual review: Verify all 3 formats (NP3/XMP/lrtemplate) as columns                |
| AC-2                                     | Data Models/No runtime models                | Table in HTML                      | Manual test: Scan table for readability, verify category grouping                  |
| AC-3                                     | Detailed Design/Format Matrix                | Table content                      | Manual review: Search for approximation symbols (~, ⚠), verify footnote            |
| AC-4                                     | Services/Format Matrix                       | Table content                      | Manual review: Verify unmappable features listed with explanations                 |
| **FR-7.3: FAQ Documentation**            |
| AC-1                                     | Services/FAQ, Legal Disclaimer               | web/index.html or FAQ.md           | Manual review: Search for "legal" question, verify answer present                  |
| AC-2                                     | Overview/Privacy promise, NFR Observability  | FAQ content                        | Manual review: Search for "privacy" question, verify WASM mentioned                |
| AC-3                                     | Services/Format Matrix                       | FAQ content                        | Manual review: Search for "doesn't convert" question, verify answer                |
| AC-4                                     | Overview, Epic 1 dependencies                | FAQ content                        | Manual review: Search for "accuracy" question, verify 95%+ mentioned               |
| AC-5                                     | Objectives/Content Requirements              | FAQ content                        | Manual review: Count sentences per answer (2-4), verify conciseness                |
| AC-6                                     | Workflows/Documentation Update               | Process documentation              | Post-launch: Monitor GitHub Issues for FAQ updates                                 |
| **FR-7.4: Legal Disclaimer**             |
| AC-1                                     | Services/Legal Disclaimer                    | web/index.html footer              | Manual review: Search for "reverse-engineered" text                                |
| AC-2                                     | Services/Legal Disclaimer                    | Disclaimer content                 | Manual review: Search for "AS IS" or "no warranty" text                            |
| AC-3                                     | Objectives/Content Requirements              | Disclaimer content                 | Manual review: Search for "private use" recommendation                             |
| AC-4                                     | Workflows/Documentation Update               | web/index.html                     | Manual test: Load page, verify disclaimer visible without scrolling/clicking       |
| AC-5                                     | Objectives/Content Requirements              | Disclaimer content                 | If legal review obtained: Verify review date/signature; else verify template used  |
| **NFR-7.1: Cloudflare Pages Deployment** |
| AC-1                                     | Data Models/Cloudflare Workflow              | .github/workflows/deploy-pages.yml | CI test: Push to main, verify workflow triggered (GitHub Actions log)              |
| AC-2                                     | Workflows/Cloudflare Deployment              | Workflow YAML                      | CI test: Verify WASM build step runs, check web/recipe.wasm exists                 |
| AC-3                                     | APIs/Cloudflare Pages API                    | Workflow YAML                      | CI test: Verify deploy step executes, check Cloudflare action logs                 |
| AC-4                                     | NFR Performance/Deployment Speed             | GitHub Actions timing              | CI monitoring: Measure workflow duration, fail if >5 minutes                       |
| AC-5                                     | APIs/User-Facing URLs                        | Deployed site                      | Manual test: Visit https://recipe.pages.dev, verify site loads                     |
| AC-6                                     | APIs/GitHub Actions Secrets                  | GitHub repository settings         | Manual review: Verify both secrets present in Settings → Secrets                   |
| AC-7                                     | Workflows/Cloudflare Deployment              | GitHub commit status               | CI test: Verify commit shows green checkmark after deployment                      |
| **NFR-7.2: GitHub Releases Setup**       |
| AC-1                                     | Data Models/GitHub Releases Workflow         | .github/workflows/release.yml      | CI test: Push tag v0.0.1, verify workflow triggered                                |
| AC-2                                     | Workflows/GitHub Release Creation            | Workflow YAML build matrix         | CI test: Verify 6 binaries built (check workflow artifacts)                        |
| AC-3                                     | APIs/GitHub Releases API                     | Workflow YAML upload step          | CI test: Verify binaries attached to release, download one                         |
| AC-4                                     | Data Models/CHANGELOG Format                 | Release description                | Manual review: Verify CHANGELOG excerpt in release notes                           |
| AC-5                                     | Data Models/Semantic Version Format          | Tag naming, CHANGELOG.md           | Manual review: Verify tag matches vMAJOR.MINOR.PATCH, CHANGELOG documents strategy |
| AC-6                                     | NFR Performance/Binary Build Speed           | GitHub Actions timing              | CI monitoring: Measure workflow duration, fail if >10 minutes                      |
| AC-7                                     | Services/README.md                           | README.md                          | Manual review: Verify installation section includes binary download URLs           |
| AC-8                                     | Services/CHANGELOG.md, Data Models           | CHANGELOG.md                       | Manual review: Verify Keep a Changelog format, version entries present             |

## Risks, Assumptions, Open Questions

### Risks

**RISK-1: Cloudflare Pages Free Tier Limitations**
- **Description:** Free tier caps at 500 builds/month; frequent commits to main could exceed limit
- **Likelihood:** Medium (if rapid development continues post-Epic 7)
- **Impact:** Medium (deployment blocked until next month or upgrade required)
- **Mitigation:** 
  - Monitor build usage in Cloudflare dashboard
  - Use feature branches, merge to main only for releases
  - Document upgrade path to paid tier if needed ($20/month for unlimited builds)
- **Contingency:** Deploy manually via Cloudflare CLI if limit reached

**RISK-2: GitHub Actions Secrets Exposure**
- **Description:** Accidental logging of CLOUDFLARE_API_TOKEN in workflow output
- **Likelihood:** Low (GitHub automatically masks known secrets)
- **Impact:** Critical (API token compromise, unauthorized deployments)
- **Mitigation:**
  - Use `::add-mask::` for any custom secret handling
  - Audit workflow logs before making public
  - Rotate tokens quarterly as best practice
- **Contingency:** Immediately revoke and regenerate token if exposed

**RISK-3: Legal Challenges to Reverse Engineering**
- **Description:** Vendor (Adobe, Nik Software) could issue DMCA or cease-and-desist
- **Likelihood:** Low (reverse engineering for interoperability generally protected)
- **Impact:** High (project takedown, legal costs)
- **Mitigation:**
  - Prominent legal disclaimer (Story 7-4)
  - Recommend private use only until legal assessment
  - No commercial distribution or monetization
  - Fair use and research exemptions (DMCA 1201)
- **Contingency:** Consult IP attorney if challenged, potentially take private

**RISK-4: Incomplete CLI Implementation for Story 7-6**
- **Description:** Epic 3 not fully complete when Story 7-6 attempts to build binaries
- **Likelihood:** Low (Epic 3 already in progress, Story 3-1 done)
- **Impact:** Medium (release workflow fails to build functional binaries)
- **Mitigation:**
  - Verify Epic 3 Stories 3-1 to 3-6 complete before starting Story 7-6
  - Test manual binary build before implementing workflow
  - Use feature flags if some CLI features incomplete
- **Contingency:** Defer Story 7-6 until Epic 3 fully complete

**RISK-5: WASM Binary Size Exceeds Expectations**
- **Description:** WASM binary >3MB compressed (NFR-1.5 target), slow page loads
- **Likelihood:** Low (current size ~1.03MB per Epic 2)
- **Impact:** Medium (violates performance NFR, poor UX)
- **Mitigation:**
  - Build with `-ldflags="-s -w"` to strip debug symbols
  - Monitor binary size in deployment workflow logs
  - Test with Lighthouse after each deployment
- **Contingency:** Investigate further optimization (TinyGo, code splitting)

### Assumptions

**ASSUMPTION-1: Cloudflare Pages Project Already Created**
- **Assumption:** Cloudflare Pages project named "recipe" exists before Story 7-5
- **Validation:** Manually create project in Cloudflare dashboard before implementing workflow
- **Impact if False:** Workflow fails with "project not found" error
- **Action:** Document project creation in Story 7-5 prerequisites

**ASSUMPTION-2: Epic 1 and Epic 2 Fully Stable**
- **Assumption:** No breaking changes to WASM interface or web/ directory after Epic 7 starts
- **Validation:** Epic 1 and Epic 2 retrospectives marked complete, no major bugs filed
- **Impact if False:** Deployment workflow publishes broken code to production
- **Action:** Code freeze Epic 1-2 during Epic 7 implementation, regression tests before deploy

**ASSUMPTION-3: Go 1.24+ Supports All Target Platforms**
- **Assumption:** Go toolchain can compile for all 6 target platforms (Linux/macOS/Windows × amd64/arm64)
- **Validation:** Manual test builds on local machine before implementing workflow
- **Impact if False:** Some binaries fail to build in release workflow
- **Action:** Test matrix build locally, remove unsupported platforms from workflow

**ASSUMPTION-4: No Custom Domain Required for MVP**
- **Assumption:** recipe.pages.dev default domain acceptable for launch
- **Validation:** User (Justin) confirms no branding/domain requirements
- **Impact if False:** Additional DNS configuration and Cloudflare setup needed
- **Action:** Document custom domain setup process for post-MVP

**ASSUMPTION-5: English-Only Documentation Sufficient**
- **Assumption:** English documentation adequate for target audience (photographers)
- **Validation:** PRD specifies English-only for MVP (out of scope: i18n)
- **Impact if False:** Non-English users unable to use tool
- **Action:** Accept limitation for MVP, defer internationalization to future epic

### Open Questions

**QUESTION-1: Legal Review Availability?**
- **Question:** Can we obtain legal review of disclaimer (FR-7.4 AC-5) before launch?
- **Options:**
  - A) Consult IP attorney (cost: $500-2000, time: 2-4 weeks)
  - B) Use conservative open-source disclaimer template (MIT/Apache style)
  - C) Defer public launch until legal review obtained
- **Decision Needed By:** Before Story 7-4 implementation
- **Recommendation:** Option B for MVP, pursue Option A post-launch if budget allows

**QUESTION-2: README.md Content Ownership?**
- **Question:** Should Stories 7-1, 7-5, 7-6 each update README.md, or consolidate in one story?
- **Options:**
  - A) Each story updates its own section (Story 7-1: docs, 7-5: web URL, 7-6: CLI install)
  - B) Story 7-1 owns all README updates, other stories provide content
  - C) Separate Story 7-7 for comprehensive README rewrite
- **Decision Needed By:** Before Story 7-1 starts
- **Recommendation:** Option A (distributed ownership), resolve merge conflicts as needed

**QUESTION-3: CHANGELOG.md Initial Version?**
- **Question:** What version should initial release be tagged? v0.1.0 (beta) or v1.0.0 (stable)?
- **Options:**
  - A) v0.1.0 - signals beta/experimental status, allows breaking changes
  - B) v1.0.0 - signals production-ready, commits to API stability
  - C) v0.0.1 - alpha status, no guarantees
- **Decision Needed By:** Before Story 7-6 implementation
- **Recommendation:** Option A (v0.1.0) - allows iteration based on user feedback

**QUESTION-4: Binary Checksums Generation Method?**
- **Question:** How to generate SHA256 checksums in release workflow (NFR-7.5)?
- **Options:**
  - A) Add workflow step: `sha256sum recipe-* > checksums.txt`, upload checksums.txt
  - B) Manual generation post-release, update release notes
  - C) Defer to post-MVP (not critical for initial release)
- **Decision Needed By:** During Story 7-6 implementation
- **Recommendation:** Option A (automated) - minimal workflow overhead, better security

**QUESTION-5: Preview Deployments for PRs?**
- **Question:** Should Cloudflare Pages create preview URLs for pull requests?
- **Options:**
  - A) Enable preview deployments (Cloudflare default) - costs more builds/month
  - B) Disable preview deployments - only deploy main branch
  - C) Enable selectively for specific PR labels
- **Decision Needed By:** During Story 7-5 workflow configuration
- **Recommendation:** Option B for MVP (conserve build quota), enable later if needed

## Test Strategy Summary

### Testing Approach

Epic 7 testing focuses on **infrastructure validation** (deployment pipelines) and **documentation quality** (user comprehension). No unit tests required - testing is primarily manual verification and CI/CD monitoring.

### Test Levels

**1. Documentation Review (Stories 7-1 to 7-4)**
- **Method:** Manual peer review + non-technical user testing
- **Scope:** Landing page, compatibility matrix, FAQ, legal disclaimer
- **Success Criteria:**
  - Non-technical user can understand "What is Recipe?" in <2 minutes
  - User can find answer to "Is my data private?" in <30 seconds
  - All links functional, no broken references
  - Spelling/grammar errors: 0 (use Grammarly or similar)
- **Testers:** SM (content review), external non-technical user (comprehension test)
- **Tools:** Browser (manual), Lighthouse (accessibility), Grammarly (grammar)

**2. CI/CD Pipeline Testing (Stories 7-5, 7-6)**
- **Method:** Integration testing via GitHub Actions
- **Scope:** Cloudflare deployment workflow, GitHub Release workflow
- **Test Scenarios:**
  - **Scenario 1:** Push to main → Verify deployment triggered, site live in <5 min
  - **Scenario 2:** Push tag v0.0.1 → Verify release created, 6 binaries uploaded
  - **Scenario 3:** Workflow failure → Verify email notification sent
  - **Scenario 4:** Rollback deployment → Verify Cloudflare Pages rollback works
- **Success Criteria:**
  - 100% workflow success rate on valid pushes
  - Deployment time <5 minutes (Cloudflare), <10 minutes (GitHub Releases)
  - Zero secret exposure in logs
- **Testers:** SM (workflow author), Dev (verification)
- **Tools:** GitHub Actions logs, Cloudflare Pages dashboard

**3. Binary Validation (Story 7-6)**
- **Method:** Manual download and execution testing
- **Scope:** All 6 CLI binaries (Linux/macOS/Windows × amd64/arm64)
- **Test Matrix:**
  | OS      | Arch  | Test Command                         | Expected Output |
  | ------- | ----- | ------------------------------------ | --------------- |
  | Linux   | amd64 | `./recipe-linux-amd64 --version`     | v0.0.1          |
  | Linux   | arm64 | `./recipe-linux-arm64 --version`     | v0.0.1          |
  | macOS   | amd64 | `./recipe-darwin-amd64 --version`    | v0.0.1          |
  | macOS   | arm64 | `./recipe-darwin-arm64 --version`    | v0.0.1          |
  | Windows | amd64 | `recipe-windows-amd64.exe --version` | v0.0.1          |
  | Windows | arm64 | `recipe-windows-arm64.exe --version` | v0.0.1          |
- **Success Criteria:**
  - Binary executes without errors on target platform
  - Version output matches tagged release
  - File size reasonable (<50MB uncompressed)
- **Testers:** Dev (multi-platform testing) or CI/CD (automated)
- **Tools:** Docker (Linux), macOS VM, Windows VM

**4. Performance Validation (NFR-7.1 to NFR-7.3)**
- **Method:** Automated + manual performance testing
- **Metrics:**
  - Deployment speed: GitHub Actions duration
  - Binary build speed: GitHub Actions matrix job duration
  - Documentation load time: Lighthouse audit
- **Targets:**
  - Cloudflare deployment: <5 minutes
  - GitHub Release build: <10 minutes
  - Landing page TTFB: <1.5 seconds, FCP: <1.5 seconds
- **Tools:** GitHub Actions timing, Lighthouse CLI, WebPageTest

**5. Security Validation (NFR-7.4 to NFR-7.6)**
- **Method:** Manual audit + automated checks
- **Test Cases:**
  - **Secrets Audit:** Review workflow logs, search for "CLOUDFLARE", verify masked
  - **Checksum Verification:** Download binary, run `sha256sum`, compare to release notes
  - **HTTPS Enforcement:** `curl -I http://recipe.pages.dev` → verify 301 redirect to HTTPS
  - **TLS Certificate:** `openssl s_client -connect recipe.pages.dev:443` → verify Cloudflare cert
- **Success Criteria:**
  - Zero secrets exposed in logs
  - 100% checksum match for all binaries
  - HTTPS redirect works, valid TLS certificate
- **Testers:** SM (security review)
- **Tools:** curl, openssl, GitHub Actions log viewer

### Coverage Goals

**Functional Coverage:**
- **Documentation (Stories 7-1 to 7-4):** 100% manual review (all ACs verified)
- **Deployment (Stories 7-5, 7-6):** 100% CI/CD integration tests (all workflows tested)

**Non-Functional Coverage:**
- **Performance:** 100% NFRs validated (deployment speed, load time)
- **Security:** 100% security controls verified (secrets, HTTPS, checksums)
- **Reliability:** Rollback tested once (manual verification)

### Test Exclusions

**Out of Scope for Epic 7:**
- **Unit tests:** No code logic to test (documentation + config files only)
- **Load testing:** Static site, CDN handles scale automatically
- **Compatibility testing:** No browser-specific code in docs (vanilla HTML/CSS)
- **Accessibility testing:** Covered in Epic 2 (Story 2-10), inherited by documentation

### Edge Cases and Error Scenarios

**Edge Case 1: Workflow Failure Recovery**
- **Scenario:** WASM build fails in Cloudflare deployment workflow
- **Expected:** Workflow fails, GitHub commit shows red X, email sent
- **Test:** Introduce syntax error in Go code, push to test branch, verify failure handling

**Edge Case 2: Binary Build Platform Failure**
- **Scenario:** One platform (e.g., Windows arm64) fails to build
- **Expected:** Partial release created, error logged, other binaries available
- **Test:** Exclude one platform from matrix, verify remaining binaries upload

**Edge Case 3: Deployment During Cloudflare Outage**
- **Scenario:** Cloudflare Pages API unavailable during deployment
- **Expected:** Workflow retries (up to 3 attempts), fails gracefully if still down
- **Test:** Manual simulation (disconnect network during deploy), verify retry logic

**Edge Case 4: Large WASM Binary**
- **Scenario:** WASM binary exceeds Cloudflare Pages size limit (25MB)
- **Expected:** Deployment fails with clear error message
- **Test:** Not tested (current size 1.03MB, unlikely to reach 25MB)

### Acceptance Testing

**Definition of Done for Epic 7:**
1. ✅ All 6 stories complete (7-1 through 7-6)
2. ✅ All 38 ACs verified PASS (see Traceability Mapping)
3. ✅ Cloudflare deployment successful (site live at recipe.pages.dev)
4. ✅ GitHub Release created (6 binaries downloadable)
5. ✅ Performance NFRs met (deployment <5 min, build <10 min, load <2s)
6. ✅ Security NFRs met (secrets secure, HTTPS enforced, checksums valid)
7. ✅ Documentation reviewed by non-technical user (comprehension confirmed)
8. ✅ Zero P0/P1 bugs outstanding

**Sign-off Required:**
- SM: Documentation quality and completeness
- Dev: CI/CD workflows functional
- User (Justin): Legal disclaimer acceptable, deployment URLs correct

### Regression Testing

**Pre-Deployment Checks (before each Cloudflare deployment):**
1. Verify Epic 2 Web UI still functional (manual smoke test)
2. Verify WASM conversion works (upload test file, download result)
3. Verify responsive design intact (mobile + desktop view)
4. Lighthouse audit score >90 (performance, accessibility, best practices)

**Pre-Release Checks (before creating GitHub Release):**
1. Verify Epic 3 CLI functional (manual smoke test on local build)
2. Verify `recipe convert` command works (test with sample file)
3. Verify `--help` and `--version` flags work
4. Run `go build` manually on each platform before automated workflow

### Test Artifacts

**Artifacts Produced:**
- GitHub Actions workflow logs (deploy-pages.yml, release.yml)
- Lighthouse audit reports (HTML export)
- Manual test checklist (markdown document per story)
- Binary checksum verification results (SHA256 hashes)
- User comprehension test notes (non-technical user feedback)

**Artifact Retention:**
- Workflow logs: GitHub default (90 days)
- Lighthouse reports: Committed to repo (docs/test-reports/)
- Manual checklists: Committed to repo (docs/stories/7-X-*.checklist.md)
- Verification results: Included in release notes
