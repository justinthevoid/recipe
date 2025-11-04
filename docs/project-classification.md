# Project Classification - Recipe

**Generated:** 2025-11-03

## Repository Structure

**Type:** Monolith (single cohesive codebase)

## Project Classification

**Project Type:** CLI Tool/Library (Python → Go Migration)

**Primary Purpose:** Photo preset format converter and extractor for Nikon and Lightroom color grading profiles

## Current State

**Legacy Implementation (Python):**
- Location: `legacy/` folder
- Languages: Python 3.7+
- Dependencies: Standard library only (no external deps)
- Status: Functional converter with 70+ preset collection

**Target Implementation (Go):**
- Status: Greenfield (not yet started)
- Purpose: Rebuild Python tools in Go + add GUI application layer

## Core Functionality

**Supported File Formats:**
- `.np3` - Nikon Picture Control files (camera color profiles)
- `.dng` - Adobe DNG camera raw format
- `.xmp` - Adobe Lightroom presets (modern)
- `.lrtemplate` - Adobe Lightroom templates (legacy)

**Current Capabilities (Python):**
- ✅ NP3 → XMP conversion
- ✅ NP3 → lrtemplate conversion
- ✅ XMP ↔ lrtemplate bidirectional conversion
- ⏳ DNG extraction (planned)
- ⏳ Reverse conversion to NP3 (planned)

**Target Capabilities (Go):**
- Extract and review .np3 color grading
- Extract Lightroom .dng, .xmp, .lrtemplate files
- Bidirectional conversion: .np3 ↔ Lightroom formats
- GUI application for visual review

## Technology Stack

### Legacy (Python)
- **Language:** Python 3.7+
- **Dependencies:** None (stdlib only)
- **Key Scripts:**
  - `nikon_picture_control_decoder.py` - Binary NP3 parser
  - `recipe_converter.py` - Universal format converter

### Target (Go)
- **Language:** Go (to be determined: version, frameworks)
- **Architecture:** CLI + GUI application
- **Status:** Planning phase

## Project Parts

### Part 1: CLI Tool (cli)
- **Root:** `C:\Users\Justin\void\recipe`
- **Type:** cli
- **Implementation:** To be developed in Go
