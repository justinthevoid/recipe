package batch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/justin/recipe/internal/apperr"
	"github.com/justin/recipe/internal/formats/nksc"
	"github.com/justin/recipe/internal/formats/np3"
)

// Config holds configuration for the batch Orchestrator.
type Config struct {
	InputPath    string
	OutputPath   string
	Recipe       *np3.Metadata
	Overwrite    bool
	Force        bool
	ManifestName string
	Workers      int
	DryRun       bool
	Strict       bool
}

// WorkerCount returns the effective number of workers to use.
func (c Config) WorkerCount() int {
	if c.Workers <= 0 {
		return runtime.NumCPU()
	}
	return c.Workers
}

// Orchestrator manages the batch processing of files.
type Orchestrator struct {
	Config          Config
	previousResults map[string]FileResult
}

// NewOrchestrator creates a new Orchestrator with the given config.
func NewOrchestrator(cfg Config) *Orchestrator {
	return &Orchestrator{
		Config:          cfg,
		previousResults: make(map[string]FileResult),
	}
}

// shouldSkip checks if a file should be skipped based on idempotency rules.
// It compares against the loaded manifest from a previous run.
// Checks are performed in the following order:
// 1. Force flag (if true, never skip)
// 2. Existence in previous manifest
// 3. Status in previous run (must be success or skipped)
// 4. Recipe hash (must match)
// 5. File size (must match)
// 6. Modification time (must match)
func (o *Orchestrator) shouldSkip(inputPath string, currentNP3Hash string, size int64, modTime time.Time) bool {
	if o.Config.Force {
		return false
	}

	res, ok := o.previousResults[inputPath]
	if !ok {
		return false // New file
	}

	if res.Status != "success" && res.Status != "skipped" {
		return false // Retry failed
	}

	if res.NP3Hash != currentNP3Hash {
		return false // Recipe changed
	}

	if res.Size != size {
		return false // File changed (size)
	}

	if !res.ModTime.Equal(modTime) {
		return false // File changed (time)
	}

	return true
}

// FindFiles recursively finds all .NEF and .NRW files in the input path.
// Matching is case-insensitive.
// It explicitly excludes files within the OutputPath to avoid processing loops.
func (o *Orchestrator) FindFiles() ([]string, error) {
	var files []string

	absOut, err := filepath.Abs(o.Config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute output path: %w", err)
	}

	err = filepath.WalkDir(o.Config.InputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Avoid processing the output directory if it's nested inside input
		absPath, err := filepath.Abs(path)
		if err == nil && strings.HasPrefix(absPath, absOut) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".nef" || ext == ".nrw" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return files, nil
}

// BatchResult summarizes the outcome of a batch operation.
type BatchResult struct {
	TotalFiles int
	Processed  int
	Skipped    int
	Failed     int
	Errors     []error
}

// ProcessBatch executes the batch processing logic.
func (o *Orchestrator) ProcessBatch(ctx context.Context) (*BatchResult, error) {
	// Support strict mode cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	startTime := time.Now()

	if o.Config.Recipe == nil {
		return nil, fmt.Errorf("configuration error: NP3 recipe is required for sidecar generation")
	}

	manifestName := o.Config.ManifestName
	if manifestName == "" {
		manifestName = "manifest.json"
	}
	manifestPath := filepath.Join(o.Config.OutputPath, manifestName)

	// Load existing manifest for idempotency
	if existing, err := ReadManifest(manifestPath); err == nil {
		for _, f := range existing.Files {
			o.previousResults[f.InputPath] = f
		}
	}

	files, err := o.FindFiles()
	if err != nil {
		return nil, err
	}

	// Calculate Hash
	hasher := sha256.New()
	hasher.Write(o.Config.Recipe.RawBytes)
	np3Hash := hex.EncodeToString(hasher.Sum(nil))

	logger := slog.Default()
	workers := o.Config.WorkerCount()
	logger.Info("starting batch processing", "workers", workers, "files", len(files))

	// Worker Pool
	// Use bounded channels to avoid O(N) memory usage for large file sets
	// Buffer size is workers*2 to ensure workers stay busy without buffering all jobs.
	jobs := make(chan string, workers*2)
	results := make(chan FileResult, workers*2)
	var wg sync.WaitGroup

	// Start Workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				results <- func() FileResult {
					// Recover from panic in processFile to protect the batch run
					defer func() {
						if r := recover(); r != nil {
							logger.Error("worker panic", "panic", r, "file", path)
							// We can't easily return a FileResult here unless we make processFile always return safely
							// or we construct one here.
						}
					}()
					// Do not ignore ctx, use the one passed to ProcessBatch
					return o.processFile(ctx, path, np3Hash)
				}()
			}
		}()
	}

	// Queue Jobs
	go func() {
		defer close(jobs)
		// Check for context cancellation during job queuing
		for _, f := range files {
			select {
			case <-ctx.Done():
				return // Stop queuing if context cancelled
			case jobs <- f:
			}
		}
	}()

	// Wait and Close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect Results
	result := &BatchResult{
		TotalFiles: len(files),
	}
	var fileResults []FileResult

	for res := range results {
		fileResults = append(fileResults, res)
		switch res.Status {
		case StatusSuccess:
			result.Processed++
		case StatusSkipped:
			result.Skipped++
		case StatusError:
			result.Failed++
			// Reconstruct error context for CLI reporting
			var err error
			if res.ErrorCode != "" {
				err = apperr.New(res.ErrorCode, filepath.Base(res.InputPath), fmt.Errorf("%s", res.ErrorMessage))
			} else {
				err = fmt.Errorf("%s: %s", filepath.Base(res.InputPath), res.ErrorMessage)
			}
			result.Errors = append(result.Errors, err)
			logger.Error("file processing failed", "file", filepath.Base(res.InputPath), "err", res.ErrorMessage)

			if o.Config.Strict {
				cancel()
				return result, fmt.Errorf("strict mode abort: %w", err)
			}
		}
	}

	endTime := time.Now()

	// Sort fileResults by InputPath for deterministic manifest
	sort.Slice(fileResults, func(i, j int) bool {
		return fileResults[i].InputPath < fileResults[j].InputPath
	})

	manifest := &BatchManifest{
		Version: ManifestVersion,
		Summary: BatchSummary{
			TotalProcessed: result.Processed + result.Skipped + result.Failed,
			SuccessCount:   result.Processed,
			SkippedCount:   result.Skipped,
			FailureCount:   result.Failed,
			StartTime:      startTime,
			EndTime:        endTime,
			Duration:       endTime.Sub(startTime).String(),
			NP3Hash:        np3Hash,
		},
		Files: fileResults,
	}

	if !o.Config.DryRun {
		if err := WriteManifest(manifestPath, manifest); err != nil {
			return result, fmt.Errorf("failed to write manifest: %w", err)
		}
	}

	return result, nil
}

// processFile handles the logic for a single file. Safe for concurrent use.
func (o *Orchestrator) processFile(ctx context.Context, srcPath string, np3Hash string) FileResult {
	fileName := filepath.Base(srcPath)
	fileLogger := slog.Default().With("file", fileName)

	// Calculate relative path to maintain structure
	relPath, err := filepath.Rel(o.Config.InputPath, srcPath)
	if err != nil {
		return FileResult{
			InputPath:    srcPath,
			Status:       StatusError,
			ErrorCode:    "resolve_path",
			ErrorMessage: err.Error(),
		}
	}

	dstPath := filepath.Join(o.Config.OutputPath, relPath)
	dstDir := filepath.Dir(dstPath)
	sidecarPath := dstPath + ".nksc"

	// Dry Run Check
	if o.Config.DryRun {
		fileLogger.Info("DRY RUN: Would process file", "input", srcPath, "output", dstPath)
		return FileResult{
			InputPath:  srcPath,
			Status:     StatusSkipped,
			OutputPath: sidecarPath,
		}
	}

	// Idempotency Check
	info, statErr := os.Stat(srcPath)
	if statErr == nil && o.shouldSkip(srcPath, np3Hash, info.Size(), info.ModTime()) {
		// Metadata matches, check payload hash and output existence
		currentHash, hashErr := CalculateFileHash(srcPath)
		if hashErr == nil {
			prevRes, ok := o.previousResults[srcPath]
			// Only skip if content hash matches AND sidecar exists
			// Note: prevRes usage is safe because map is read-only during processing
			if ok && prevRes.PayloadHash == currentHash {
				if _, err := os.Stat(sidecarPath); err == nil {
					fileLogger.Debug("skipping unchanged file")
					res := prevRes
					res.Status = StatusSkipped
					return res
				}
			}
		}
	}

	// Create output directory structure
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return FileResult{
			InputPath:    srcPath,
			Status:       StatusError,
			ErrorCode:    "create_dir",
			ErrorMessage: err.Error(),
		}
	}

	// 1. Handle NEF/Image File copying
	nefExists := false
	if _, err := os.Stat(dstPath); err == nil {
		nefExists = true
	}

	if !nefExists {
		fileLogger.Debug("copying source image")
		if err := CopyFile(ctx, srcPath, dstPath); err != nil {
			return FileResult{
				InputPath:    srcPath,
				Status:       StatusError,
				ErrorCode:    "copy_nef",
				ErrorMessage: err.Error(),
			}
		}
	}

	// Double check context before proceeding to potential sidecar generation
	if err := ctx.Err(); err != nil {
		return FileResult{
			InputPath:    srcPath,
			Status:       StatusError,
			ErrorCode:    "context_cancelled",
			ErrorMessage: err.Error(),
		}
	}

	// 2. Handle Sidecar Generation
	sidecarExists := false
	if _, err := os.Stat(sidecarPath); err == nil {
		sidecarExists = true
	}

	// Determine if this is an update to a previously known file.
	_, isKnownFile := o.previousResults[srcPath]

	if sidecarExists && !o.Config.Overwrite && !o.Config.Force && !isKnownFile {
		fileLogger.Debug("skipping existing sidecar")
		return FileResult{
			InputPath:  srcPath,
			Status:     StatusSkipped,
			OutputPath: sidecarPath,
		}
	}

	// Generate Sidecar
	fileLogger.Debug("generating sidecar")
	if err := o.generateSidecar(srcPath, dstPath); err != nil {
		return FileResult{
			InputPath:    srcPath,
			Status:       StatusError,
			ErrorCode:    "generate_sidecar",
			ErrorMessage: err.Error(),
		}
	}

	var size int64
	var modTime time.Time
	if statErr == nil {
		size = info.Size()
		modTime = info.ModTime()
	}

	payloadHash, _ := CalculateFileHash(srcPath)

	return FileResult{
		InputPath:   srcPath,
		Status:      StatusSuccess,
		OutputPath:  sidecarPath,
		Size:        size,
		ModTime:     modTime,
		PayloadHash: payloadHash,
		NP3Hash:     np3Hash,
	}
}

func (o *Orchestrator) generateSidecar(srcPath, targetNEF string) error {
	if o.Config.Recipe == nil {
		return nil
	}

	recipe := nksc.NewNKSCRecipe(o.Config.Recipe, targetNEF)

	// Use the recipe's Write method which handles XMP packet wrapping and atomic writing
	sidecarPath := targetNEF + ".nksc"
	if err := recipe.Write(sidecarPath); err != nil {
		return err
	}
	return nil
}
