stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: []
workflowType: 'research'
lastStep: 6
research_type: 'technical'
research_topic: 'Bulk apply Nikon NP3 Picture Control recipes to NEF using NX Studio sidecars (NKSC) and optional NX-look JPG/TIFF export'
research_goals: 'Identify the most reliable approach to: (1) apply NP3 picture controls non-destructively (NEF copy + NX-compatible sidecar), and (2) optionally render JPG/TIFF matching NX Studio, preferably headless on Windows where NX Studio is installed.'
user_name: 'Justin'
date: '2026-01-06'
web_research_enabled: true
source_verification: true
---

# Research Report: technical

**Date:** 2026-01-06
**Author:** Justin
**Research Type:** technical

---

## Research Overview

[Research overview and methodology will be appended here]

## Technical Research Scope Confirmation

**Research Topic:** Bulk apply Nikon NP3 Picture Control recipes to NEF using NX Studio sidecars (NKSC) and optional NX-look JPG/TIFF export
**Research Goals:** Identify the most reliable approach to: (1) apply NP3 picture controls non-destructively (NEF copy + NX-compatible sidecar), and (2) optionally render JPG/TIFF matching NX Studio, preferably headless on Windows where NX Studio is installed.

**Technical Research Scope:**

- Architecture Analysis - design patterns, system architecture options
- Implementation Approaches - practical automation/prototyping approaches
- Technology Stack - tools/libraries/platform constraints (WSL dev, Windows runtime)
- Integration Patterns - how to glue NP3→sidecar→NX-render together
- Performance Considerations - batching, throughput, and failure handling patterns

**Research Methodology:**

- Current web data where available (with URL citations) to verify and supplement
- Repo-backed evidence (existing reverse-engineering docs, sample sidecars, and scripts) clearly labeled as such
- Confidence levels for uncertain areas, and explicit “next experiment” recommendations to eliminate uncertainty

**Scope Confirmed:** 2026-01-06

---

<!-- Content will be appended sequentially through research workflow steps -->

## Step 2 — Technology Stack Analysis

### Web search analysis (sources used + why they matter)

This step relies on a mix of (a) vendor/help documentation for the “ground truth” capabilities and constraints of NX Studio, and (b) metadata tooling documentation to validate what we can safely read/write.

- Nikon NX Studio download page + license terms: confirms the supported OS requirement (Windows 11 64-bit) and includes explicit restrictions against reverse engineering / decompiling, which affects the risk profile of any runtime-hooking approach.
	- https://downloadcenter.nikonimglib.com/en/download/sw/275.html
- Nikon NX Studio help: confirms that NX Studio exports edited images to JPEG/HEIF/TIFF and supports managing Custom Picture Controls (import/export).
	- https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html
- ExifTool Nikon tag documentation: confirms ExifTool recognizes Nikon NKSC sidecar XMP namespaces/tags including `ast:XMLPackets` and `nine:NineEdits` (and the NineEdits FilterParameters* tags), which strongly supports a file-based, sidecar-centric integration strategy.
	- https://exiftool.org/TagNames/Nikon.html
- ExifTool overview: establishes ExifTool as a general-purpose, cross-platform CLI/library for reading/writing XMP and other metadata (useful for validation and inspection tooling).
	- https://exiftool.org/
- Frida documentation: confirms Frida’s purpose as a dynamic instrumentation toolkit with multiple language bindings; relevant only for “last resort” introspection when file formats or behavior can’t be derived from samples.
	- https://frida.re/docs/home/
- Microsoft UI Automation documentation: establishes the standard Windows UI automation APIs, relevant only if NX Studio must be driven headlessly via UI automation (last-resort fallback).
	- https://learn.microsoft.com/en-us/windows/win32/winauto/entry-uiauto-win32
	- https://learn.microsoft.com/en-us/dotnet/framework/ui-automation/ui-automation-overview

### Programming languages

- Go (primary): The repository is Go-based (`go.mod`) and already targets a fast, portable CLI. Go is a good fit for:
	- Structured binary parsing/serialization (NP3/NCP-like payloads)
	- Deterministic batch processing and concurrency controls
	- Cross-compilation for Windows runners
- JavaScript/TypeScript (web UI + WASM glue): The web app uses Vite + Svelte (see `web/package.json`) and the README states conversions run in-browser via WebAssembly.
- Python (research/forensics tooling): The repo contains Python scripts under `scripts/` and `.reverse-engineering/scripts/` used for inspecting sidecars and decoding payloads.

### Frameworks and major libraries

- CLI and batching
	- Cobra for CLI parsing (already in `go.mod`).
- Metadata / sidecar inspection
	- ExifTool can read/write XMP metadata and explicitly supports Nikon NKSC as “XMP-based” and identifies Nikon NX Studio’s NKSC tag families (including `ast` and `nine`) (https://exiftool.org/TagNames/Nikon.html, https://exiftool.org/).
- Nikon sidecar (NKSC) schema anchors
	- ExifTool documents Nikon “ast Tags” (notably `XMLPackets`) and Nikon “nine Tags” (notably `NineEdits`), with “Nikon NineEdits Tags” such as `FilterParametersExportExportData` and related fields, which align with our current repo findings that Picture Control application is carried in a `NineEdits` edit graph.
		- https://exiftool.org/TagNames/Nikon.html
- NX Studio (vendor application)
	- NX Studio can process RAW and export the results to common formats including TIFF/JPEG/HEIF (https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html).
	- NX Studio is Windows 11 64-bit per Nikon’s system requirements on the download page (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).

### Database and storage solutions

- Primary storage: filesystem only.
	- Inputs: `.nef`/`.nrw` and `.np3`.
	- Outputs:
		- “Non-destructive” path: NEF copy + generated/modified `.nksc` sidecar.
		- Optional rendered output: JPEG/TIFF/HEIF via NX Studio (export).
- Optional job state persistence (if needed later): a small SQLite DB or simple JSONL manifest could track batches, retries, and per-file status. No DB is required for an MVP batch runner.

### Cloud infrastructure

- Not required. The stated project direction is local-first processing.
- If remote execution is ever needed, the most realistic model is a Windows runner (VM or bare-metal) with NX Studio installed; but this would be an operational choice, not a product requirement.

### Development tools and DevOps

- Go toolchain: build/test/format.
- Node toolchain (Vite/Svelte) for the web UI.
- ExifTool: validate generated NKSC/XMP and inspect diffs (https://exiftool.org/).
- Reverse-engineering / introspection (only if needed)
	- Frida can instrument native apps at runtime (https://frida.re/docs/home/).
	- Caution: Nikon’s license terms include restrictions against reverse engineering, decompiling, disassembly, etc. (https://downloadcenter.nikonimglib.com/en/download/sw/275.html). This pushes us toward preferring file-format-level interoperability and black-box behavioral testing over invasive instrumentation.

### Integration and architecture patterns

The goal implies two “pipelines” that share parsing/mapping logic:

1) Apply NP3 → generate/patch NKSC (non-destructive)
	 - Treat NKSC as the integration boundary. ExifTool’s identification of `ast:XMLPackets` and `nine:NineEdits` supports this as a stable artifact that NX Studio understands (https://exiftool.org/TagNames/Nikon.html).
	 - Internal architecture:
		 - Parse NP3 (and existing NKSC, if present)
		 - Produce a new `nine:NineEdits` graph node (or update existing) for Picture Control data
		 - Preserve unrelated edits/metadata to avoid clobbering

2) Optional render/export via NX Studio
	 - NX Studio explicitly exports enhanced/resized pictures to JPEG/HEIF/TIFF (https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html).
	 - Integration options (in preference order):
		 - If NX Studio exposes any supported CLI/API/import/export hooks: use that.
		 - Otherwise: UI Automation for scripted exports on Windows (https://learn.microsoft.com/en-us/windows/win32/winauto/entry-uiauto-win32).
		 - Runtime instrumentation (Frida) is technically possible (https://frida.re/docs/home/) but has legal/compliance risk due to Nikon license restrictions (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).

### Security and compliance

- License / reverse-engineering restrictions
	- Nikon’s download/license terms for NX Studio explicitly restrict decompiling, reverse engineering, disassembly, and similar activities (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).
	- Implication: prefer interoperability work based on documented/observed file formats and standard metadata tooling; treat invasive app introspection as last-resort and potentially out-of-scope for shared/public automation.
- Data privacy
	- The project is local-first; no cloud upload required.
- Automation safety
	- If UI automation is used, it must be hardened against focus/desktop/session issues; Microsoft notes UI Automation does not enable communication across processes started by different users (relevant to service accounts / “Run as” patterns) (https://learn.microsoft.com/en-us/dotnet/framework/ui-automation/ui-automation-overview).

### Additional technical considerations

- Platform split is real:
	- Dev environment can be Linux/WSL.
	- “Exact NX output render” is likely Windows-only due to NX Studio requirements (Windows 11 64-bit) (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).
- Validation strategy
	- Use ExifTool to verify that generated sidecars are structurally sound and include expected Nikon XMP families (https://exiftool.org/TagNames/Nikon.html).
	- Use NX Studio itself as the final arbiter: does it show the edit, does it render as expected.

---

**Step 2 complete.**

---

## Step 6 — Technical Synthesis and Completion

### Web verification recap (key external anchors)

This synthesis is grounded in a small set of “anchor facts” verified against current public sources:

- NX Studio supports processing RAW images and exporting enhanced/resized pictures to JPEG, HEIF, or TIFF.
	- https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html
- NX Studio (1) is distributed for Windows 11 (pre-installed 64-bit editions only) and (2) includes license restrictions that explicitly prohibit decompiling / reverse engineering / disassembly.
	- https://downloadcenter.nikonimglib.com/en/download/sw/275.html
- ExifTool explicitly documents Nikon “ast Tags” as used by Nikon NX Studio in Nikon NKSC sidecar files and trailers, and Nikon “nine Tags” including `NineEdits`, described as XML-based tags used to store editing information.
	- https://exiftool.org/TagNames/Nikon.html

### Executive Summary

This research converges on a practical, compatibility-first strategy: treat NX Studio’s `.nksc` sidecar as the primary integration artifact and build a deterministic “NP3 → NKSC” pipeline around it. ExifTool’s Nikon tag documentation explicitly treats NKSC as NX Studio’s sidecar format and highlights the `ast` and `nine` namespaces, including `nine:NineEdits` as the container for edit instructions (https://exiftool.org/TagNames/Nikon.html). That supports an artifact-first approach: generate/patch the sidecar NX Studio already understands, rather than trying to embed edits inside NEF or depend on undocumented internal APIs.

For “pixel-perfect Nikon output”, Nikon’s own help confirms NX Studio can export enhanced images to JPEG/HEIF/TIFF (https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html). This makes “NX-matching JPEG/TIFF/HEIF rendering” feasible but realistically Windows-hosted and operationally separate from the core sidecar generator.

Finally, Nikon’s download page provides two constraints that influence the engineering approach: NX Studio’s published Windows system requirement, and explicit license language prohibiting reverse engineering / decompilation / disassembly (https://downloadcenter.nikonimglib.com/en/download/sw/275.html). Combined, these constraints push the solution toward file-format interoperability + black-box acceptance testing, and away from invasive runtime instrumentation.

### Synthesis: recommended operating model

The design naturally splits into two pipelines that share the same NP3 parsing and mapping logic:

1) **Non-destructive apply (core product):**
	- Inputs: `.nef` + `.np3`
	- Outputs: NEF copy (optional) + generated/modified `.nksc` sidecar
	- Goal: NX Studio recognizes the edit and displays the intended Picture Control “recipe” without altering the RAW.

2) **Optional NX Studio render/export (capability add-on):**
	- Inputs: `.nef` + `.nksc` (produced by pipeline #1)
	- Outputs: JPEG/HEIF/TIFF
	- Justification: NX Studio explicitly supports exporting enhanced pictures in these formats (https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html).

### Synthesis: recommended system architecture

- **Core library (Go):**
	- NP3 parsing → stable internal “recipe model”
	- NKSC generation/patching from internal model (NineEdits-centered)
	- Deterministic serialization so batch operations are repeatable

- **CLI (Go):**
	- Batch traversal + concurrency controls
	- Atomic writes (temp + rename) to prevent partial `.nksc` artifacts
	- Idempotency (skip when same recipe already applied)
	- Structured per-file results output suitable for future UI/daemon integration

- **Validation harness (tooling):**
	- ExifTool structural checks to validate the NKSC is well-formed and contains expected Nikon namespaces (`ast`, `nine`, `NineEdits`) (https://exiftool.org/TagNames/Nikon.html).
	- NX Studio acceptance checks on a small fixture set (manual or semi-automated) as the compatibility “truth test”.

- **Optional Windows NX export runner:**
	- Kept separate so it can run only where NX Studio is installed and usable (Windows 11 requirement per Nikon) (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).
	- Treated as a bounded capability with explicit operational expectations (interactive session vs service constraints, flakiness management).

### Risks and mitigations

- **NKSC schema drift (NX Studio updates):** mitigate with a golden fixture corpus and regular NX Studio acceptance checks.
- **Exact look parity:** mitigate by offering optional NX export for users who require Nikon-rendered output.
- **Automation/legal constraints:** avoid invasive “inside NX” techniques where possible; prefer sidecar-first interoperability and black-box behavior verification (Nikon terms explicitly prohibit reverse engineering-like activities) (https://downloadcenter.nikonimglib.com/en/download/sw/275.html).

### Implementation roadmap (phased)

1) Sidecar-first MVP: batch apply NP3 → generate `.nksc` reliably.
2) Hardening: deterministic output, atomic writes, idempotency, resumable manifests.
3) Validation: ExifTool-based structural checks + minimal NX Studio acceptance workflow.
4) Optional Windows export runner: export to JPEG/HEIF/TIFF where NX rendering is required (https://nikonimglib.com/nxstdo/onlinehelp/en/what_nx_studio_can_do_for_you_1.html).

---

**Step 6 complete.**

## Step 3 — Integration Patterns Analysis

### Web search analysis (sources used + why they matter)

This step focuses on integration boundaries and protocols for a batch pipeline that spans Linux/WSL dev + Windows (NX Studio) runtime. Sources used:

- ExifTool Nikon tag documentation: corroborates NKSC as XMP-based and enumerates the relevant tag families (`XMP-ast` / `XMP-nine`) including `ast:XMLPackets` and `nine:NineEdits`, plus the `FilterParameters...ExportData/CustomData` fields that match our observed Picture Control payload placement. This supports a file/artifact-based integration boundary (generate `.nksc` sidecars) as the primary interoperability strategy.
	- https://exiftool.org/TagNames/Nikon.html
- RFC 6455 (WebSocket): establishes the standards-track protocol semantics for a bidirectional, persistent client/server channel (handshake, framing, security considerations). This is relevant if we introduce a local/remote “job runner” service that streams progress/events to a UI.
	- https://www.rfc-editor.org/rfc/rfc6455
- RFC 9457 (Problem Details for HTTP APIs): defines `application/problem+json` for machine-readable HTTP errors. This is useful if we introduce an HTTP API surface (local daemon or remote runner) and want consistent, debuggable error payloads.
	- https://www.rfc-editor.org/rfc/rfc9457
- Microsoft Azure REST API Guidelines (non-normative, but widely adopted patterns): provides concrete patterns for long-running operations (LROs), polling, `operation-location`, idempotency, and `retry-after`. These are relevant if “export via NX Studio” becomes an asynchronous job.
	- https://github.com/microsoft/api-guidelines/blob/vNext/azure/Guidelines.md

### Integration boundary choices (what integrates with what)

For this project, the cleanest integration boundary is *files as contracts*:

- Inputs: `NEF` + `NP3`.
- Primary output: `NEF copy` + `NKSC sidecar` (NX Studio-compatible).
- Optional secondary outputs: rendered `JPEG/TIFF` generated by NX Studio.

This artifact-first boundary is robust because it avoids undocumented “live” APIs and uses the same contract NX Studio already understands (the NKSC sidecar structure corroborated by ExifTool’s Nikon XMP tag taxonomy) (https://exiftool.org/TagNames/Nikon.html).

### API design patterns

#### 1) File-based “API” (recommended default)

For batch apply, the “API” is a directory layout + naming convention:

- Deterministic output paths (e.g., `out/<basename>.NEF` + `out/<basename>.NEF.nksc`).
- A manifest (JSON/JSONL) describing per-file status, timing, and produced artifacts.
- Atomic writes (write temp file then rename) to keep outputs safe under interruption.

This yields simple, portable integration with other tools (shell, Make, CI runners) and keeps “correctness” testable by inspecting the generated sidecars with ExifTool.

#### 2) Local/remote job runner HTTP API (optional, when a UI needs orchestration)

If we add a service that runs NX Studio export jobs (Windows-only runner), the API is best modeled as a long-running job system:

- `POST /jobs` (create) → returns a job resource + a polling URL.
- `GET /jobs/{id}` (status) → returns state machine: `NotStarted | Running | Succeeded | Failed | Canceled`.

For LRO-style behavior, the Azure API guidelines describe returning an `operation-location` and using `retry-after` for client polling cadence (https://github.com/microsoft/api-guidelines/blob/vNext/azure/Guidelines.md).

For error payloads, use RFC 9457 problem details (`application/problem+json`) so clients can reliably parse error types, titles, and per-occurrence detail without inventing a custom error envelope (https://www.rfc-editor.org/rfc/rfc9457).

#### 3) Streaming progress/events (optional)

If the UI needs realtime progress (per-file conversion events, ETA, logs), add a WebSocket endpoint (e.g., `/jobs/{id}/events`). WebSocket’s bidirectional semantics and handshake/framing are defined in RFC 6455 (https://www.rfc-editor.org/rfc/rfc6455).

### Communication protocols

- Filesystem + manifests (default): simplest and most robust for CLI-first workflows.
- HTTP polling (optional): appropriate for asynchronous, Windows-hosted NX export.
- WebSockets (optional): appropriate for realtime UI progress and log streaming; RFC 6455 defines the wire protocol and security model (https://www.rfc-editor.org/rfc/rfc6455).

### Data formats and standards

- NKSC sidecars (domain format): treat the `.nksc` XML as the interoperability contract. ExifTool’s Nikon tag documentation confirms the `XMP-ast`/`XMP-nine` families (`ast:XMLPackets`, `nine:NineEdits`) and the filter parameter fields that align to Picture Control payload embedding (https://exiftool.org/TagNames/Nikon.html).
- JSON for job APIs/manifests (optional): practical for internal orchestration and portable across languages.
- Problem Details for API errors (optional): RFC 9457 standardizes the structure of machine-readable HTTP API error objects (https://www.rfc-editor.org/rfc/rfc9457).

### System interoperability approaches

#### Point-to-point (CLI → filesystem) (recommended)

- `recipe-cli` reads NP3 and writes `.nksc` next to NEFs.
- Downstream tools (NX Studio, other validators) consume the files.

This minimizes coupling and avoids long-lived services.

#### Service boundary (Linux/WSL coordinator → Windows runner) (optional)

If we must drive NX Studio for exports, treat Windows as an execution node:

- The coordinator submits a “render job” referencing file paths (or shared storage) and a target export config.
- The runner reports status via polling endpoints and optionally WebSockets.

This is the point where long-running operation patterns become useful (polling + `operation-location` + `retry-after`) (https://github.com/microsoft/api-guidelines/blob/vNext/azure/Guidelines.md).

### Microservices integration patterns (only if we go distributed)

If a Windows runner is introduced, keep the service surface minimal:

- Single “runner” service responsible only for “apply sidecar + export”.
- Strict idempotency on job creation (same request should not create multiple jobs).
- Clear separation of responsibilities: coordinator handles batching, runner handles NX execution.

The Azure guidelines emphasize idempotency and “exactly once” semantics as a prerequisite for fault-tolerant clients that retry requests (https://github.com/microsoft/api-guidelines/blob/vNext/azure/Guidelines.md).

### Event-driven integration (optional)

Event-driven patterns are useful if we need to fan out many exports across runners, or if a UI must react to granular progress:

- Job lifecycle events: `job.created`, `job.started`, `file.exported`, `job.completed`, `job.failed`.
- Transport: WebSockets for direct UI streaming (RFC 6455) or a message broker if scaling beyond one machine.

If we keep to a single-machine workflow, events can remain “local” (written to JSONL logs) to avoid introducing unnecessary infrastructure.

### Integration security patterns (only if there is a network API)

If we introduce a network-facing service (even on localhost), implement a basic security posture:

- Default bind to localhost only.
- If remote access is required, use TLS and a standard auth mechanism.

For standards-backed auth patterns, OAuth 2.0 (RFC 6749) and JWT (RFC 7519) are commonly used building blocks, but they’re likely unnecessary for a purely local CLI → file pipeline.
	- https://www.rfc-editor.org/rfc/rfc6749
	- https://www.rfc-editor.org/rfc/rfc7519

---

**Step 3 complete.**

---

## Step 4 — Architectural Patterns

### Web search analysis (sources used + why they matter)

This step focuses on architectural patterns that keep the system (a) robust as a batch pipeline, and (b) realistic about the hard constraint that NX Studio rendering is a Windows GUI application.

- Microsoft Azure Architecture Center (Microservices architecture style): provides a practical taxonomy of microservice components, benefits, and tradeoffs (independent deployability vs overall system complexity) that helps decide if/when a Windows “runner” should be treated as a separate service boundary.
	- https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/microservices
- Microsoft Azure Architecture Center (Event-driven architecture style): summarizes producer/consumer decoupling, publish-subscribe vs streaming, and topology choices (broker vs mediator) that map directly onto “batch job progress events” and orchestration design.
	- https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/event-driven
- Microsoft documentation on Session 0 isolation: confirms that Windows services run in session 0 and that session 0 does not support processes that interact with the user. This is critical for NX Studio automation architecture because NX Studio is GUI-driven and UI automation is typically tied to an interactive session.
	- https://learn.microsoft.com/en-us/windows/win32/services/service-changes-for-windows-vista
- Linux `rename(2)` semantics: documents atomic replacement behavior when renaming over an existing path. This supports crash-safe sidecar generation (write temp → atomic rename), reducing corruption and “half-written” NKSC outputs.
	- https://man7.org/linux/man-pages/man2/rename.2.html
- The Twelve-Factor App methodology: useful as a set of operational design principles (separate build/run stages, externalize config, treat logs as event streams, run admin tasks as one-off processes) if we introduce a runner daemon or background agent.
	- https://12factor.net/

### Architectural patterns and principles that fit this problem

#### 1) Artifact-first pipeline (preferred baseline)

Use files as the primary contract and treat every batch run as a deterministic pipeline that produces inspectable artifacts:

- Input artifacts: `NEF` + `NP3` (+ optional preexisting `NEF.nksc`).
- Output artifacts:
	- Non-destructive: copied `NEF` + generated/updated `.nksc`.
	- Optional: rendered `JPEG/TIFF` exported by NX Studio (Windows-only).
	- Batch metadata: per-file manifest + logs.

This architecture keeps the core “apply recipe” path independent of any GUI automation.

Crash safety note: write outputs to a temp path and then rename into place. On Linux, `rename()` specifies that if the destination exists, it is “atomically replaced” (no window where another process sees the destination missing), which is exactly what we want for `.nksc` writes (https://man7.org/linux/man-pages/man2/rename.2.html).

#### 2) Ports-and-adapters (hexagonal) for core logic + OS-specific runners

Treat “NP3 → internal recipe model → NKSC edit graph” as the domain core, surrounded by adapters:

- Domain core (pure logic): parse/validate NP3 semantics, map to a stable internal representation, synthesize a NineEdits update.
- Adapters:
	- Input adapters: filesystem, ZIP/bundles, existing NKSC reader.
	- Output adapters: NKSC writer, manifest writer.
	- Optional execution adapter: “NX export runner” interface that can be implemented via Windows-only automation.

This makes it possible to unit-test most behavior on Linux/WSL while keeping the NX-export problem isolated behind an interface.

#### 3) Coordinator/worker batching (single-machine) + job model (multi-machine optional)

Even without a server, model work as jobs:

- A “batch job” expands into per-file “work items” (apply sidecar, verify, optional export).
- A bounded worker pool processes items concurrently, with per-file retry policies.

If you later introduce a Windows runner, that runner becomes an independently deployed component (in the same spirit as the “independent deployability” discussed in microservices guidance), but you should only pay the distributed-systems complexity cost when it is truly needed (https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/microservices).

#### 4) Windows runner must assume an interactive session (avoid “Windows service drives NX Studio”)

If NX Studio export automation is required, architectural feasibility hinges on Windows session boundaries.

Microsoft’s Session 0 isolation note states that session 0 is reserved for services and “does not support processes that interact with the user”, and services cannot directly display UI such as dialog boxes (https://learn.microsoft.com/en-us/windows/win32/services/service-changes-for-windows-vista). Practically, this argues for a runner design that executes inside an interactive user session (agent/tray app / scheduled task “only when user is logged on”), rather than a pure Windows Service that tries to drive NX Studio UI in session 0.

This is an architectural constraint, not an implementation detail: it should shape how “headless” we promise the export path can be.

#### 5) Event-driven progress reporting as an internal contract (optional)

Even if everything is local, structure progress as a stream of events:

- Event producer: the batch engine emits events like `file.started`, `sidecar.written`, `export.started`, `export.completed`, `file.failed`.
- Consumers: CLI renderer, web UI, JSONL log writer, metrics.

This aligns with event-driven architecture benefits around decoupling producers/consumers (https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/event-driven) while keeping the transport simple (often just logs or an in-process bus). If the project grows into a distributed runner setup, these same events can be published over a network channel later without redesigning the core.

The event-driven style guidance also notes two broad topology choices: a broker-like broadcast (highly decoupled) vs a mediator that manages state and error handling (more control but more coupling) (https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/event-driven). For this project, a mediator-like coordinator tends to fit better because export orchestration is stateful and error-prone.

### Deployment and operations patterns

If (and only if) a Windows runner component is introduced, adopt “boring ops” practices so it’s maintainable:

- Externalize config (paths, concurrency, export presets) rather than hard-coding per-environment settings (https://12factor.net/).
- Treat logs as event streams (JSONL is a good fit) so the CLI/UI can tail progress and parse structured errors (https://12factor.net/).
- Separate build/release/run stages and keep admin tasks (cleanup, cache prune, validation runs) as one-off processes, which maps well to a CLI-centric toolchain (https://12factor.net/).

### Concrete architecture recommendation (MVP)

1) Keep the primary deliverable as a pure CLI batch pipeline:

- Reads NP3 + NEF paths
- Generates/upgrades NKSC sidecars
- Emits deterministic manifests and logs

2) Make NX export a separately deployable “runner” *only when needed*:

- Runner executes in an interactive Windows session (per Session 0 isolation constraints) (https://learn.microsoft.com/en-us/windows/win32/services/service-changes-for-windows-vista)
- Runner accepts jobs via filesystem drop folder (simplest) or a minimal local HTTP API

3) Use event-driven progress internally so the UX can evolve without refactoring core processing (https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/event-driven).

---

**Step 4 complete.**

---

## Implementation Approaches and Technology Adoption

### Technology Adoption Strategies

For this project, the safest adoption strategy is explicitly incremental:

- Treat “NP3 → NKSC sidecar generation” as the first shippable capability, because it’s file/artifact-based and doesn’t require automating NX Studio.
- Add “NX export automation” only after sidecar generation is correct and stable, because export introduces Windows GUI/session constraints and significantly higher flake risk.
- Keep transitional architecture intentionally: for a while, you’ll likely support multiple paths (e.g., generate sidecars only; optionally also do export on machines where NX Studio is available).

This maps closely to the Strangler Fig approach: add new capability in small slices, create seams/interfaces, and gradually shift responsibility into the new pipeline instead of trying to “replace everything at once” (https://martinfowler.com/bliki/StranglerFigApplication.html).

_Source: https://martinfowler.com/bliki/StranglerFigApplication.html_

### Development Workflows and Tooling

Recommended workflow posture:

- Automate build/test on every pull request and on mainline merges.
- Use a build matrix (at least Linux + Windows) once any Windows runner components exist, to catch OS-specific issues early.
- Keep developer workflows fast: `go test ./...` as the default local check, with heavier integration tests gated behind flags or separate CI jobs.

GitHub Actions provides the core primitives for this: YAML-defined workflows triggered by repo events, executed on GitHub-hosted runners (or self-hosted runners if you need NX Studio present for export tests) (https://docs.github.com/en/actions/get-started/understand-github-actions).

_Source: https://docs.github.com/en/actions/get-started/understand-github-actions_

### Testing and Quality Assurance

Testing strategy should match risk:

- Unit tests for:
	- NP3 parsing (including boundary cases)
	- NKSC XML construction and “don’t clobber unrelated fields”
	- Deterministic serialization (stable ordering, stable whitespace rules if needed)
- Golden-file tests for known “good” NKSC outputs, validated via ExifTool + NX Studio spot checks.
- Fuzz tests for parsers/decoders (especially anything handling binary blobs / opaque payloads) to prevent panics and discover edge cases early; Go’s `testing` package supports fuzzing via `FuzzXxx(*testing.F)` and `go test` fuzz flags (https://pkg.go.dev/testing).
- Benchmarks for bulk operations if throughput becomes a bottleneck, using `BenchmarkXxx(*testing.B)` and `go test -bench` (https://pkg.go.dev/testing).

Go’s built-in testing conventions (`_test.go`, `TestXxx`, `go test`) are designed for incremental correctness hardening (https://go.dev/doc/tutorial/add-a-test).

_Sources: https://go.dev/doc/tutorial/add-a-test and https://pkg.go.dev/testing_

### Deployment and Operations Practices

For the CLI-only path, “deployment” is mostly packaging/distribution. If a Windows runner service/agent is introduced, operations should be treated explicitly:

- Observability:
	- Emit structured logs and a per-batch manifest so failures are diagnosable.
	- If you have a daemon, define minimal service health signals and track failures as first-class metrics.
	- Use the “four golden signals” framing (latency, traffic, errors, saturation) to decide what to measure (even for batch systems: job duration/queue depth/error rate/resource saturation) (https://sre.google/sre-book/monitoring-distributed-systems/).
- Reliability:
	- Design for restartability and idempotency at the per-file level.
	- Prefer “resume” semantics using manifests rather than “start over and hope.”
- Operational excellence:
	- Automate routine tasks (build/test/release) and reduce manual steps.

The Azure Well-Architected pillars provide a useful checklist lens (reliability, security, cost optimization, operational excellence, performance efficiency) even if you are not deploying to Azure (https://learn.microsoft.com/en-us/azure/architecture/framework/).

_Sources: https://sre.google/sre-book/monitoring-distributed-systems/ and https://learn.microsoft.com/en-us/azure/architecture/framework/_

### Team Organization and Skills

If you keep an artifact-first CLI pipeline, a small team can ship it end-to-end. If you add a Windows runner component, treat it as a separate bounded capability with clear ownership:

- Keep responsibilities clean (core sidecar logic vs Windows export orchestration).
- Avoid over-splitting into microservices unless you genuinely need distributed scaling; microservices guidance emphasizes that while services can be independently deployable, overall system complexity rises quickly (https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/microservices).

_Source: https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/microservices_

### Cost Optimization and Resource Management

The main cost levers here aren’t cloud spend; they’re engineering/ops time, Windows machine availability, and test data management:

- Default to the artifact-first flow because it avoids always-on infrastructure and minimizes automation flake.
- If NX export is needed, treat it as an opt-in capability that runs only where NX Studio is installed.
- Use batching controls (worker pool size, IO scheduling) to avoid saturating disks/CPUs.

The “Cost Optimization” pillar framing is still a helpful lens: optimize for the constraint that actually costs you (debug time, manual babysitting, reruns) (https://learn.microsoft.com/en-us/azure/architecture/framework/).

_Source: https://learn.microsoft.com/en-us/azure/architecture/framework/_

### Risk Assessment and Mitigation

Key risks and mitigations:

- Format risk (NX Studio accepts/rejects NKSC): build fixture-based tests and validate generated sidecars against NX Studio regularly.
- Automation risk (NX export via GUI): keep export optional; isolate it behind a runner interface; expect interactive-session requirements.
- Batch reliability risk: idempotent per-file operations + resumable manifests + structured logging/metrics.

This aligns with the general monitoring guidance to keep alerting/diagnostics simple and actionable (https://sre.google/sre-book/monitoring-distributed-systems/).

_Source: https://sre.google/sre-book/monitoring-distributed-systems/_

## Technical Research Recommendations

### Implementation Roadmap

1) Ship sidecar generation as the primary deliverable (bulk apply NP3 → NKSC).
2) Add a validation harness (ExifTool structural checks + a small set of NX Studio acceptance checks).
3) Harden with fuzzing and golden-file tests.
4) Only then prototype NX export automation behind a Windows runner boundary (filesystem job queue first; HTTP/WebSocket later only if needed).

### Technology Stack Recommendations

- Keep core logic in Go (fits existing repo).
- Use Go `testing` + `go test` as the main correctness harness; add fuzzing for parsers and benchmarks where needed.
- Keep orchestration simple (filesystem contracts) unless multi-machine coordination becomes a real requirement.

### Skill Development Requirements

- Go (testing, fuzzing, benchmarks, concurrency patterns).
- XMP/NKSC structure handling and “don’t clobber unrelated metadata” discipline.
- Windows automation basics if NX export becomes mandatory.

### Success Metrics and KPIs

- Sidecar acceptance rate in NX Studio: % of outputs where NX Studio recognizes the applied picture control/edit.
- Determinism: repeated runs produce byte-identical NKSC for the same inputs (when applicable).
- Batch reliability: error rate per 1,000 files; mean/95p job duration.
- Operational load: time-to-diagnose failures (log quality), rerun frequency, and “babysitting required” incidents.

---

**Step 5 complete.**
