---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments:
  - docs/analysis/brainstorming-session-2026-01-06T11-29-51.md
  - docs/project-planning-artifacts/research/technical-np3-nksc-nx-studio-batch-apply-research-2026-01-06.md
  - docs/analysis/product-brief-recipe-2025-12-18.md
  - docs/PRD.md
  - docs/architecture.md
  - docs/user-context.md
  - docs/technology-stack.md
  - README.md
date: '2026-01-06'
author: Justin
---

# Product Brief: recipe

<!-- Content will be appended sequentially through collaborative workflow steps -->

## Executive Summary

Recipe is expanding beyond preset-format conversion into a workflow automation tool for Nikon Z shooters: bulk-applying NP3 Picture Control “recipes” to NEF RAW files without the manual, one-by-one process in Nikon NX Studio.

Today, applying an NP3 look to a shoot means repeatedly opening each NEF in NX Studio and applying the Picture Control manually. For occasional batches of ~20–200 images, this is pure time/effort overhead and blocks a “pre-applied look” review workflow.

The vision is a batch-first pipeline that produces NX Studio–compatible, non-destructive results (NEF + sidecar) so images open in NX Studio with the desired Picture Control already applied. As an optional add-on, Recipe can run a Windows-only export step (on machines with NX Studio installed) to render NX-matching JPEG/TIFF outputs.

---

## Core Vision

### Problem Statement

Applying NP3 Picture Controls to NEF RAW photos is a manual, repetitive NX Studio workflow (open photo → apply NP3 → repeat) that does not scale to even modest shoot sizes.

### Problem Impact

- Wasted time on repetitive, low-value clicks for 20–200 image batches
- Delayed review: the look isn’t “pre-applied” until each file is opened and edited
- Prevents an automation-friendly workflow for Nikon Picture Controls on RAW files

### Why Existing Solutions Fall Short

- The baseline solution is NX Studio manual editing (one image at a time)
- There is no known reliable batch workflow that (a) applies NP3 looks at scale and (b) remains compatible with NX Studio

### Proposed Solution

A batch-capable feature in Recipe that:

- Takes NEF files + an NP3 Picture Control as inputs
- Produces NX Studio–compatible, non-destructive outputs so the look is already applied when browsing/opening in NX Studio (sidecar-based interoperability)
- Optionally runs an export step on Windows (where NX Studio is installed) to render NX-matching JPEG/TIFF outputs for users who want final deliverables

Constraints and positioning:

- Windows-only is acceptable for the NX-dependent parts of the workflow
- The core value is time savings and “pre-applied look” review; export/render is an optional capability

### Key Differentiators

- Batch-first workflow designed around “apply once, process many” for occasional 20–200 file runs
- Compatibility-first approach: outputs designed to be opened and trusted in NX Studio
- Non-destructive by default: preserves RAW data while enabling “pre-applied” viewing
- Optional NX-native rendering path for users who require Nikon’s output look

## Target Users

### Primary Users

#### Persona 1: Alex — “Nikon Film Recipe Enthusiast”

- **Context:** Shoots with a Nikon Z body that supports Picture Controls and collects NP3 “recipes” to get specific looks. Shoots personal projects and trips; edits in NX Studio when needed.
- **Motivations/Goals:** Wants a consistent look applied across a set of NEFs before review, without spending time clicking through each image.
- **Problem Experience:** Today must open NX Studio and apply a Picture Control one-by-one. This is tedious even for 20–50 images and discourages experimenting with looks.
- **Success Vision:** Picks a folder of NEFs + one NP3 recipe, runs a batch, then opens NX Studio and sees the look already applied across the whole set.

#### Persona 2: Maya — “Occasional Event Shooter”

- **Context:** Shoots occasional events/family sessions (20–200 NEFs per batch). Prefers Nikon’s rendering/look and uses NX Studio as the “truth” for Nikon color.
- **Motivations/Goals:** Save time and reduce mistakes while keeping a Nikon-native workflow. Wants repeatability across a shoot.
- **Problem Experience:** Manual apply is repetitive and error-prone (missed images, inconsistent application). Delays review and selection.
- **Success Vision:** Batch apply once, then immediately review/cull with the intended look already present; optionally export a set of JPEGs/TIFFs when needed.

#### Persona 3: Sam — “Hybrid Creator / Small-Batch Producer”

- **Context:** Shoots photos for content or small client deliverables. Likes the Nikon look and wants fast turnaround without learning a complex post workflow.
- **Motivations/Goals:** Speed and consistency; wants “good enough Nikon look” applied reliably.
- **Problem Experience:** NX Studio manual steps don’t fit quick production bursts.
- **Success Vision:** A repeatable batch command that outputs NX-compatible results and (optionally) export-ready images.

### Secondary Users

- **Clients / recipients of exports:** Benefit indirectly if the optional NX-matching JPEG/TIFF export is used (faster delivery, more consistent look).
- **Future “power users” / automation folks:** May integrate this into a scripted workflow (watch folders, per-shoot presets), even if the first release is CLI-first.

### User Journey

- **Discovery:** Finds Recipe via Nikon communities, GitHub, or “NP3 recipe” / “Picture Control” workflow searches.
- **Onboarding (first run):** Installs/opens the tool on Windows (NX-dependent workflow is Windows-only), selects a folder of NEFs and an NP3 Picture Control, and chooses output behavior (sidecar-only vs sidecar + optional export).
- **Core Usage:** Runs batch apply for 20–200 files occasionally; expects a simple per-file report (success/failure) and predictable output locations.
- **Success Moment (“aha!”):** Opens NX Studio and sees the look already applied across the set without manual per-image edits.
- **Long-term Use:** Reuses saved “recipes” / batch presets for repeatable looks per shoot type; uses export only when deliverables are needed.

## Success Metrics

### User Success (Outcomes + Behaviors)

- **Time saved vs manual NX Studio:** Users can apply one NP3 Picture Control to a batch of 20–200 NEFs in minutes instead of doing per-image manual work.
  - Baseline expectation: manual workflow requires opening each file and applying the Picture Control one-by-one.
  - Success threshold (initial): “Batch apply completes with <5 minutes of user attention for 100 files” (setup + start + verify).
- **Pre-applied review achieved:** After running the batch, users can open/browse the images in NX Studio and see the intended Picture Control already applied without repeating manual steps.
- **Repeatability:** Users can rerun the same batch configuration on a new shoot and get consistent behavior and outputs (predictable file locations and naming).

### Quality / Compatibility

- **NX Studio acceptance rate (core):** ≥ 95% of processed NEFs open in NX Studio with the Picture Control visibly applied as expected (initial target; aim to raise to ≥ 99% as fixtures grow).
- **Non-destructive guarantee:** The RAW image data is not altered; changes are represented via NX Studio–compatible sidecars and/or non-destructive metadata paths.
- **Idempotency (nice-to-have early, required later):** Re-running the tool on an already-processed batch should not degrade or duplicate edits; it should detect and skip or safely overwrite.

### Reliability / Operational Robustness

- **Continue-on-error batching:** Batch runs complete even if individual files fail; failures are reported clearly with per-file status.
- **Failure rate:** < 5% file failures per batch on supported camera/NEF variants (initial), trending toward < 1% as format coverage and fixtures expand.
- **Clear diagnostics:** For any failed file, the tool provides actionable output (what failed, which file, suggested next step).

### Performance

- **Sidecar-only mode performance:** 100 files complete in < 60 seconds on a typical modern Windows machine (initial target).
- **Export mode performance (optional, NX Studio dependent):** Export completes reliably but may be slower and environment-dependent; success is defined primarily by correctness + completion without manual babysitting, with time-to-complete tracked for optimization.

### Business Objectives

- **3-month objective (ship value):** Deliver a stable, Windows-first batch workflow that removes the one-by-one NX Studio apply process for the majority of supported NEF inputs.
- **Adoption objective:** At least a small cohort of Nikon users can run it end-to-end successfully (self + early testers), producing repeatable results across multiple shoots.
- **Quality objective:** Build a fixture corpus and acceptance checks sufficient to confidently prevent regressions in NX compatibility.

### Key Performance Indicators

- **Batch completion rate:** % of runs that complete (even with partial failures) without requiring manual intervention beyond starting the job.
- **Per-file success rate:** Successful files / total files per batch (target ≥ 95% initially).
- **NX acceptance rate:** % of outputs confirmed as “opens in NX Studio with look applied” on the fixture set (target ≥ 95% initially).
- **Median runtime (sidecar-only):** Median time to process 100 files (target < 60s initially).
- **User effort proxy:** Number of required user interactions per batch after configuration (target: start job + optional review only).

## MVP Scope

### Core Features

- **Batch apply NP3 → NX Studio viewing (non-destructive):** Apply a selected `.np3` Picture Control to a batch of `.nef` files so NX Studio opens/browses them with the look already applied.
- **CLI + folder workflow:** Primary UX is command-line driven, taking a folder (or set of files) as input and producing outputs into a specified output folder.
- **Never modify originals:** Original input NEFs are not changed in-place. Outputs are written separately (e.g., copied NEFs and NX Studio–compatible sidecar(s)).
- **NX-matching export (optional but supported in MVP):** On Windows machines with NX Studio installed, provide an export step that renders Nikon/NX Studio matching outputs (e.g., JPEG/TIFF) for the processed batch.

### Out of Scope for MVP

- **Editing inside NEF / destructive changes:** No in-place modification of original RAW files.
- **Cross-platform NX export:** NX-dependent rendering/export is Windows-only.
- **Full GUI/Web UX:** CLI-first only; no polished UI beyond basic terminal output.
- **Always-on automation:** No watch-folders, background daemons, or cloud processing in MVP.

### MVP Success Criteria

- **Pre-applied look works:** After processing, NX Studio shows the intended Picture Control look without manual per-image steps.
- **Non-destructive guarantee holds:** Original NEFs remain untouched; outputs are created safely and predictably.
- **Batch robustness:** Processing completes for a folder-sized batch (typical 20–200 files), continues on per-file failures, and reports which files succeeded/failed.
- **Export produces usable deliverables (when enabled):** On supported Windows setups, NX export completes and produces the expected JPEG/TIFF outputs.

### Future Vision

- **Richer reporting and automation ergonomics:** Dry-run mode and a structured per-file report.
- **Better UX surfaces:** A simple folder UI while retaining a CLI core.
- **Runner model if export needs orchestration:** Optional Windows runner patterns (interactive-session constraints) for more reliable export automation at scale.
