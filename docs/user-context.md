# User-Provided Context

**Generated:** 2025-11-03

## Architectural Direction

### UI Framework Options

**CLI Implementation (Preferred):**
- **Framework:** Bubbletea TUI (Terminal User Interface)
- **Language:** Go
- **Rationale:** Interactive terminal-based interface for file conversion and review
- **Reference:** https://github.com/charmbracelet/bubbletea

**GUI Implementation (Alternative):**
- **Framework:** Wails or similar
- **Language:** Go + Web frontend (HTML/CSS/JS)
- **Rationale:** Desktop application with graphical interface
- **Reference:** https://wails.io/

**Decision Status:** To be determined during planning/architecture phase

## Focus Areas

### Scan Scope
- **Primary:** `legacy/` folder (complete Python implementation)
- **Secondary:** `examples/` folder (sample files of each format)
- **Purpose:** Understand existing functionality before Go migration

### Migration Priorities
1. Understand existing Python converter logic
2. Document file format specifications (NP3, DNG, XMP, lrtemplate)
3. Plan Go architecture (CLI vs GUI vs both)
4. Design conversion engine architecture
5. Implement parser/decoder for each format
