package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/justin/recipe/internal/converter"
	"github.com/spf13/cobra"
)

// BatchFlags holds all flags for batch processing mode
type BatchFlags struct {
	To              string
	From            string
	OutputDir       string
	Parallel        int
	Verbose         bool
	JSON            bool
	ContinueOnError bool
	FailFast        bool
	Overwrite       bool
}

var batchCmd = &cobra.Command{
	Use:   "batch [pattern]",
	Short: "Convert multiple files in batch using glob patterns",
	Long: `Convert multiple preset files in parallel using glob patterns.

The CLI auto-detects the source format from the file extension.
You must specify the target format with --to.

Examples:
  recipe batch *.xmp --to np3
  recipe batch presets/**/*.xmp --to np3 --output-dir converted
  recipe batch *.xmp --to np3 --parallel 8 --overwrite
  recipe batch testdata/xmp/*.xmp --to np3 --json
  recipe batch *.costyle --to xmp
  recipe batch presets/*.xmp --to costylepack --output-dir bundles`,
	Args: cobra.ExactArgs(1),
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Required flags
	batchCmd.Flags().StringP("to", "t", "", "Target format (required): np3, xmp, lrtemplate, costyle, or costylepack")
	batchCmd.MarkFlagRequired("to")

	// Optional flags
	batchCmd.Flags().StringP("from", "f", "", "Source format (auto-detected if omitted)")
	batchCmd.Flags().String("output-dir", "", "Output directory for converted files")
	batchCmd.Flags().IntP("parallel", "p", 0, "Number of parallel workers (default: NumCPU)")
	batchCmd.Flags().Bool("verbose", false, "Enable verbose logging")
	batchCmd.Flags().Bool("json", false, "Output results as JSON")
	batchCmd.Flags().Bool("continue-on-error", true, "Continue processing on errors")
	batchCmd.Flags().Bool("fail-fast", false, "Stop on first error")
	batchCmd.Flags().Bool("overwrite", false, "Overwrite existing output files")
}

func runBatch(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Parse flags
	flags := BatchFlags{
		To:              mustGetString(cmd, "to"),
		From:            mustGetString(cmd, "from"),
		OutputDir:       mustGetString(cmd, "output-dir"),
		Parallel:        mustGetInt(cmd, "parallel"),
		Verbose:         mustGetBool(cmd, "verbose"),
		JSON:            isJSONMode(cmd), // Use unified helper
		ContinueOnError: mustGetBool(cmd, "continue-on-error"),
		FailFast:        mustGetBool(cmd, "fail-fast"),
		Overwrite:       mustGetBool(cmd, "overwrite"),
	}

	// Validate target format (AC-8)
	if err := validateFormat(flags.To); err != nil {
		return err
	}

	// Expand glob pattern (AC-1, AC-8)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid glob pattern: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files match pattern: %s", pattern)
	}

	// Process batch
	result, err := processBatch(files, flags)
	if err != nil {
		return err
	}

	// Display result using unified output function (AC-4, AC-7)
	outputBatchResult(*result, flags.JSON)

	// Exit code 1 if any errors (AC-5)
	if result.ErrorCount > 0 {
		return fmt.Errorf("batch completed with %d errors", result.ErrorCount)
	}

	return nil
}

// processBatch processes all files using a worker pool pattern (AC-2)
func processBatch(files []string, flags BatchFlags) (*BatchResult, error) {
	// Start timing for batch processing
	startTime := time.Now()

	// Log batch start (AC-7)
	logger.Info("starting batch conversion",
		"count", len(files),
		"target", flags.To)

	// Determine worker pool size (AC-2)
	numWorkers := runtime.NumCPU()
	if flags.Parallel > 0 {
		numWorkers = flags.Parallel
	}
	logger.Debug("worker pool configured", "workers", numWorkers)

	// Channels for work distribution
	jobs := make(chan string, len(files))
	results := make(chan ConversionResult, len(files))
	stopChan := make(chan struct{})

	// Progress tracking (AC-4)
	var processed atomic.Int32
	total := len(files)

	// Worker pool (AC-2)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, stopChan, flags, &processed, total, &wg)
	}

	// Distribute work
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)
	close(stopChan)

	// Aggregate results (AC-7)
	result := aggregateResults(results, total)
	result.DurationMs = time.Since(startTime).Milliseconds()

	// Log batch completion (AC-7)
	logger.Info("batch complete",
		"success", result.SuccessCount,
		"error", result.ErrorCount,
		"duration_ms", result.DurationMs)

	return result, nil
}

// worker processes files from the jobs channel (AC-2)
func worker(
	jobs <-chan string,
	results chan<- ConversionResult,
	stopChan <-chan struct{},
	flags BatchFlags,
	processed *atomic.Int32,
	total int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for inputPath := range jobs {
		// Check for stop signal (fail-fast mode)
		select {
		case <-stopChan:
			return
		default:
		}

		// Process single file
		result := convertSingleFileForBatch(inputPath, flags)
		results <- result

		// Update progress (AC-4)
		p := processed.Add(1)

		// Log per-file progress in verbose mode (AC-7)
		logger.Debug("processing file",
			"index", p,
			"total", total,
			"file", inputPath)

		if !flags.JSON {
			fmt.Fprintf(os.Stderr, "\rProcessing %d/%d files...", p, total)
		}

		// Fail-fast mode: signal stop on error (AC-5)
		if flags.FailFast && !result.Success {
			return
		}
	}
}

// convertSingleFileForBatch converts a single file and returns the result
func convertSingleFileForBatch(inputPath string, flags BatchFlags) ConversionResult {
	start := time.Now()

	result := ConversionResult{
		Input:        inputPath,
		TargetFormat: flags.To,
	}

	// Auto-detect source format if not specified (AC-10 - reuse from Story 3-2)
	sourceFormat := flags.From
	if sourceFormat == "" {
		var err error
		sourceFormat, err = detectFormat(inputPath)
		if err != nil {
			result.Error = fmt.Sprintf("auto-detect failed: %v", err)
			result.DurationMs = time.Since(start).Milliseconds()
			return result
		}
	} else {
		// Validate explicit source format
		if err := validateFormat(sourceFormat); err != nil {
			result.Error = fmt.Sprintf("invalid source format: %v", err)
			result.DurationMs = time.Since(start).Milliseconds()
			return result
		}
	}
	result.SourceFormat = sourceFormat

	// Generate output path (AC-6 - reuse from Story 3-2)
	var outputPath string
	if flags.OutputDir != "" {
		// Custom output directory
		filename := filepath.Base(inputPath)
		outputPath = filepath.Join(flags.OutputDir, generateOutputPath(filename, flags.To))
	} else {
		// Same directory as input
		outputPath = generateOutputPath(inputPath, flags.To)
	}
	result.Output = outputPath

	// Check overwrite protection (AC-9 - reuse from Story 3-2)
	if err := checkOutputExists(outputPath, flags.Overwrite); err != nil {
		result.Error = "file already exists"
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Read input file
	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read input file: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Convert (AC-10 - single API call to converter)
	outputBytes, err := converter.Convert(inputBytes, sourceFormat, flags.To)
	if err != nil {
		result.Error = fmt.Sprintf("conversion failed: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Ensure output directory exists (AC-6)
	if err := ensureOutputDir(outputPath); err != nil {
		result.Error = fmt.Sprintf("failed to create output directory: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Write output file
	if err := os.WriteFile(outputPath, outputBytes, 0644); err != nil {
		result.Error = fmt.Sprintf("failed to write output file: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		return result
	}

	// Success
	result.Success = true
	result.FileSizeBytes = int64(len(outputBytes))
	result.DurationMs = time.Since(start).Milliseconds()
	return result
}

// aggregateResults collects and aggregates all conversion results (AC-7)
func aggregateResults(results <-chan ConversionResult, expectedTotal int) *BatchResult {
	batch := &BatchResult{
		Batch:   true, // AC-4: Batch field identifies batch operation
		Results: make([]ConversionResult, 0, expectedTotal),
	}

	for result := range results {
		batch.Total++
		batch.Results = append(batch.Results, result)

		if result.Success {
			batch.SuccessCount++
		} else {
			batch.ErrorCount++
		}
	}

	return batch
}

// Helper functions to safely extract flag values
func mustGetString(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	return val
}

func mustGetInt(cmd *cobra.Command, name string) int {
	val, _ := cmd.Flags().GetInt(name)
	return val
}

func mustGetBool(cmd *cobra.Command, name string) bool {
	val, _ := cmd.Flags().GetBool(name)
	return val
}
