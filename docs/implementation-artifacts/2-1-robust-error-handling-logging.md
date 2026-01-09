# Story 2.1: Robust Error Handling & Logging

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a DevOps User (Jordan),
I want detailed logs differentiating between file failures and system errors,
So that I can diagnose why specific images failed without digging through source code.

## Acceptance Criteria

1.  **Given** A batch run with some corrupt NEF files.
2.  **When** I run the tool with `--log-level debug`.
3.  **Then** The output includes structured logs using `slog`.
4.  **And** Errors are wrapped with file context and phase (e.g., "IMG_123.NEF: parse_error: invalid header").
5.  **And** The batch does NOT abort on individual file failures (Continue-on-error).
6.  **And** A summary is printed at the end (e.g., "Processed: 100, Success: 98, Failed: 2").

## Tasks / Subtasks

- [x] Define Structured Error Types and Package
    - [x] Create `internal/apperr` (or similar) to hold domain error types (Phase, File, Cause).
    - [x] Ensure errors implement `Unwrap()` compatibility.
- [x] Implement Logging Infrastructure
    - [x] Update `cmd/nx` root/apply to handle `--log-level` flag (debug, info, warn, error).
    - [x] Configure `slog` handler (JSON or Text) based on environment/flag (Text for CLI default).
- [x] Enhance Orchestartor Robustness (`internal/batch`)
    - [x] Refactor `ProcessBatch` loop to use `continue` on file-level errors.
    - [x] Wrap errors with file context (filename, path) before logging/collecting.
    - [x] Implement failure counters and collection.
- [x] Add Contextual Logging
    - [x] Instrument `internal/formats/np3` parsing with debug logs (e.g. "Parsing header", "Reading variable").
    - [x] Instrument `internal/formats/nksc` generation.
    - [x] Instrument file copying/writing operations.
- [x] Implementing Summary Reporting
    - [x] Ensure `Orchestrator` returns strict stats struct.
    - [x] Print formatted summary to Stdout at end of run.

## Dev Notes

- **Architecture Pattern**: Use `log/slog` for structured logging.
- **Error Handling**: Use `errors.As` and `errors.Is` for checking error types. Define explicit error variables/types for known failure modes (e.g., `ErrInvalidNP3`, `ErrCorruptNEF`).
- **Resiliency**: The batch loop must be wrapped in a way that a panic in a file processor is caught (optional, but good for robustness) or at least standard errors don't break the loop.
- **CLI Output**: "Info" level should be user-friendly (progress). "Debug" level should be verbose for devs/support.

### Project Structure Notes

- `internal/apperr` or `internal/pkg/errors` is recommended to avoid circular dependencies if errors are shared.
- Reuse `internal/batch` but ensure it doesn't log directly to stdout/stderr if possible, but takes a Logger instance or uses global `slog` properly configured.

### References

- [Epics: Story 2.1](docs/project-planning-artifacts/epics.md#story-21-robust-error-handling--logging)
- [Architecture: Error Strategy](docs/project-planning-artifacts/architecture.md#error-handling-strategy)

## Dev Agent Record

### Agent Model Used

Gemini 3 Pro (Preview)

### Debug Log References

- Created `internal/apperr` package for structured error handling.
- Integrated `slog` into `cmd/nx` (Global setup in `PersistentPreRun`).
- Enhanced `internal/batch/orchestrator.go` to use `apperr`, wrap errors, and continue on failure.
- Added debug logging to `internal/formats/np3` and `nksc`.
- Verified error wrapping and robustness with `orchestrator_test.go`.

### Completion Notes List

- Implemented standard logging infrastructure using `log/slog`.
- `cmd/nx` now accepts `--log-level`.
- Helper package `apperr` created to standardize errors across `process`.
- Batch processing is now robust against single file failures.

### File List

- internal/apperr/doc.go
- internal/apperr/error.go
- internal/apperr/error_test.go
- cmd/nx/main.go
- cmd/nx/main_test.go
- cmd/nx/apply.go
- internal/batch/orchestrator.go
- internal/batch/orchestrator_test.go
- internal/formats/np3/parse.go
- internal/formats/nksc/recipe.go
