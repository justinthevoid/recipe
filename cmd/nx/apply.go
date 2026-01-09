package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/justin/recipe/internal/batch"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/utils"
	"github.com/spf13/cobra"
)

// dependency injection for testing
var openFolderFunc = utils.OpenFolder

func newApplyCmd() *cobra.Command {
	var (
		np3Path      string
		inputDir     string
		outputDir    string
		manifestName string
		overwrite    bool
		force        bool
		workers      int
		dryRun       bool
		strict       bool
		exportGuide  bool
	)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply recipe to specific files",
		Long:  `Apply an NP3 recipe to a directory of NEF files by generating NKSC sidecars.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use default logger configured in root command
			logger := slog.Default()

			// 1. Load NP3 Recipe
			logger.Info("loading recipe", "path", np3Path)
			np3Data, err := os.ReadFile(np3Path)
			if err != nil {
				return fmt.Errorf("failed to read NP3 file: %w", err)
			}

			recipe, err := np3.ParseMetadata(np3Data)
			if err != nil {
				return fmt.Errorf("failed to parse NP3 file: %w", err)
			}

			// 2. Validate Input Directory
			info, err := os.Stat(inputDir)
			if err != nil {
				return fmt.Errorf("failed to access input directory: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("input path is not a directory: %s", inputDir)
			}

			// 3. Create Output Directory
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// 4. Initialize Orchestrator
			logger.Info("starting batch processing",
				"input", inputDir,
				"output", outputDir,
				"overwrite", overwrite,
				"workers", workers,
			)

			cfg := batch.Config{
				InputPath:    inputDir,
				OutputPath:   outputDir,
				Recipe:       recipe,
				Overwrite:    overwrite,
				Force:        force,
				ManifestName: manifestName,
				Workers:      workers,
				DryRun:       dryRun,
				Strict:       strict,
			}
			orch := batch.NewOrchestrator(cfg)

			// 5. Run Processing
			startTime := time.Now()
			result, processErr := orch.ProcessBatch(cmd.Context())
			duration := time.Since(startTime)

			// 6. Print Summary if we have results (even partial ones from strict mode failure)
			if result != nil {
				logger.Info("batch summary",
					"total", result.TotalFiles,
					"success", result.Processed,
					"failed", result.Failed,
					"skipped", result.Skipped,
					"duration", duration.String(),
				)

				if len(result.Errors) > 0 {
					logger.Error("some files failed", "count", len(result.Errors))
					for i, e := range result.Errors {
						if i >= 5 {
							logger.Error("...and more errors")
							break
						}
						logger.Error("file error", "err", e)
					}
				}
			}

			if processErr != nil {
				return fmt.Errorf("batch processing failed: %w", processErr)
			}

			if exportGuide {
				if result != nil && result.Processed == 0 && dryRun {
					logger.Warn("skipping export guide in dry-run mode settings", "dry_run", dryRun)
					fmt.Fprintln(cmd.OutOrStdout(), "\n(Export guide skipped: Dry-run active, no sidecars generated)")
					return nil
				}
				// 1. Open Output Folder
				if err := openFolderFunc(outputDir); err != nil {
					logger.Warn("failed to open output folder", "err", err)
					fmt.Fprintf(cmd.OutOrStdout(), "Warning: Could not open folder: %v\n", err)
				}

				// 2. Print Instructions
				fmt.Fprintln(cmd.OutOrStdout(), "\n--- NX Studio Export Guide ---")
				fmt.Fprintln(cmd.OutOrStdout(), "1. Open NX Studio")
				fmt.Fprintf(cmd.OutOrStdout(), "2. Navigate to: %s\n", outputDir)
				fmt.Fprintln(cmd.OutOrStdout(), "3. Select All images")
				fmt.Fprintln(cmd.OutOrStdout(), "4. Export as JPEG (Standard)")
				fmt.Fprintln(cmd.OutOrStdout(), "------------------------------")

				// 3. Interactive Prompt
				fmt.Fprint(cmd.OutOrStdout(), "Press Enter when export is finished...")
				scanner := bufio.NewScanner(cmd.InOrStdin())
				scanner.Scan()

				// 4. Verification
				fmt.Fprintln(cmd.OutOrStdout(), "\nVerifying exports...")

				entries, err := os.ReadDir(outputDir)
				if err != nil {
					logger.Warn("failed to read output directory for verification", "err", err)
				} else {
					count := 0
					for _, entry := range entries {
						if entry.IsDir() {
							continue
						}
						ext := strings.ToLower(filepath.Ext(entry.Name()))
						if ext == ".jpg" || ext == ".jpeg" || ext == ".tif" || ext == ".tiff" {
							count++
						}
					}

					if result != nil && count == result.Processed {
						fmt.Fprintf(cmd.OutOrStdout(), "✅ Verification Successful: Found %d/%d exported images\n", count, result.Processed)
					} else {
						expected := 0
						if result != nil {
							expected = result.Processed
						}
						fmt.Fprintf(cmd.OutOrStdout(), "⚠️ Verification Warning: Expected %d images, found %d\n", expected, count)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&np3Path, "np3", "n", "", "Path to NP3 recipe file (required)")
	cmd.Flags().StringVarP(&inputDir, "input", "i", "", "Input directory containing NEF files (required)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for NKSC sidecars (required)")
	cmd.Flags().StringVar(&manifestName, "manifest", "manifest.json", "Name of the manifest file")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing NKSC files")
	cmd.Flags().BoolVar(&force, "force", false, "Force re-processing of all files (ignores idempotency checks)")
	cmd.Flags().IntVar(&workers, "workers", 0, "Number of parallel workers (0 = auto-detect)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without modifying files")
	cmd.Flags().BoolVar(&strict, "strict", false, "Abort immediately on first error")
	cmd.Flags().BoolVar(&exportGuide, "export-guide", false, "Show interactive export guide after processing")

	_ = cmd.MarkFlagRequired("np3")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
