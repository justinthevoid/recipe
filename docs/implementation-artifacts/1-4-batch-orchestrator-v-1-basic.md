# Story 1.4: Batch Orchestrator v1 (Basic)

Status: review

## Story

As a User,
I want the tool to find all NEF and NRW files in a folder structure and process them,
so that I don't have to specify files one by one and can process entire shoots at once.

## Acceptance Criteria

1.  **Given** An input directory containing `.NEF` or `.NRW` files (case-insensitive) nested in subfolders.
2.  **When** I run the batch orchestrator with input and output paths.
3.  **Then** It recursively discovers all raw files.
4.  **And** It replicates the relative folder structure in the output directory.
5.  **And** For each file, it performs an "Exact Copy" (content + modification time + permissions) to the output location.
6.  **And** It generates a corresponding `.nksc` sidecar file next to the copied raw file (using the implementation from Story 1.3).
7.  **And** If an individual file fails (e.g., read permission), it logs the error and continues processing the rest of the batch (Partial Success).
8.  **And** It returns a summary report (Total Found, Processed, Failed).

## Tasks / Subtasks

- [x] Create `internal/batch` package
    - [x] Define `Orchestrator` struct and configuration options
    - [x] Implement `FindFiles` using `filepath.WalkDir` with case-insensitive matching (`.nef`, `.nrw`)
- [x] Implement File Operation Helpers
    - [x] `CopyFile(src, dst)`: atomic-like copy ensuring content, `ModTime`, and `Mode` are preserved
- [x] Implement `ProcessBatch` Logic
    - [x] Loop through discovered files
    - [x] Replicate subfolder structure in output (`os.MkdirAll`)
    - [x] Execute Copy step
    - [x] Execute Sidecar Generation step using `nksc.NewNKSCRecipe` facade
    - [x] Capture errors in a `BatchResult` struct without stopping loop
- [x] Add Unit Tests
    - [x] `TestFindFiles_Recursive`: Setup temp dir with nested files, verify all are found
    - [x] `TestCopyFile_Metadata`: Verify modtime is preserved after copy
    - [x] `TestProcessBatch_PartialFailure`: Simulate one unreadable file, verify others succeed
    - [x] Verify integration with `nksc` package using mock or actual NP3 data

## Technical Implementation Notes

### Use of `internal/batch`
Create a new package `internal/batch` to encapsulate this logic. This keeps `cmd/` clean and allows for easier testing.

### File Discovery
Use `filepath.WalkDir`. It is more efficient than `Walk` as it avoids `os.Stat` calls for every file.

### "Exact Copy" requirements
To satisfy "Exact Copy" (preserving metadata bits):
```go
// After io.Copy...
info, _ := os.Stat(src)
os.Chtimes(dst, info.ModTime(), info.ModTime())
os.Chmod(dst, info.Mode())
```

### Sidecar Integration
Since Story 1.3 is done, we will import `recipe/internal/formats/nksc`.
Logic:
```go
// Load NP3 once at start (np3.Metadata)
// For each file:
//   // Use the NKSCRecipe facade for sidecar logic
//   recipe := nksc.NewNKSCRecipe(np3Metadata, targetNEFPath)
//   // Generate XML bytes
//   xmlBytes, err := recipe.MarshalXML()
//   // Write to sidecar path
```
*Correction*: Story 1.2 "NP3 Parsing & NKSC Model" should have provided the way to get `NKSCRecipe` from NP3. If not, this story implicitly requires wiring that up. Assuming `nksc` package has a constructor or mapper.

### Error Handling
Do not verify inputs inside the critical loop. Use a `files processed` counter and an `errors []error` slice. Return a `BatchSummary` struct.

## Dev Agent Record

### Agent Model Used
Gemini 3 Pro (Preview)

### Plan
1. Create `internal/batch` directory.
2. Implement file walker.
3. Implement robust copy.
4. Wire up NKSC generation.
### Completion Notes
- Implemented `internal/batch` package with `Orchestrator` struct.
- Implemented recursive `FindFiles` for .NEF/.NRW using `filepath.WalkDir`.
- Implemented `CopyFile` helper ensuring preservation of permissions and modification times.
- Implemented `ProcessBatch` logic to coordinate file finding, copying, and sidecar generation.
- Added robust error handling in `ProcessBatch` to allow partial success (counting failed files without aborting).
- Created comprehensive unit tests: `TestOrchestrator_FindFiles`, `TestCopyFile_Metadata`, `TestProcessBatch_PartialFailure`.
- Verified integration with `nksc` package for sidecar generation.
- **Fixed validation issues (2026-01-08)**: Added recursion protection to `FindFiles`, strictly enforced NP3 recipe presence, and committed missing dependencies.

## File List
- internal/batch/doc.go
- internal/batch/orchestrator.go
- internal/batch/orchestrator_test.go
- internal/batch/file_ops.go
- internal/batch/file_ops_test.go
- internal/batch/process_test.go

## Change Log
- 2026-01-08: Implemented Batch Orchestrator v1 (Basic) including recursive file finding, metadata-preserving copy, and sidecar generation integration. Added full unit test coverage.
