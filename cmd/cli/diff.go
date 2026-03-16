package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/inspect"
	"github.com/justin/recipe/internal/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var diffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "Compare parameters between two preset files",
	Long: `Diff compares all parameters between two preset files across any format.

Files can be different formats (e.g., NP3 vs XMP). Shows only changed parameters
by default. Use --unified to see all fields including unchanged ones.

Examples:
  recipe diff original.np3 converted.xmp
  recipe diff file1.xmp file2.xmp --unified
  recipe diff original.np3 converted.xmp --format=json
  recipe diff file1.xmp file2.xmp --tolerance=0.01 --no-color`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	// Add flags
	diffCmd.Flags().Bool("unified", false, "Show all fields, not just changes")
	diffCmd.Flags().String("format", "text", "Output format: text or json")
	diffCmd.Flags().Float64("tolerance", 0.001, "Float comparison tolerance (default: 0.001)")
	diffCmd.Flags().Bool("no-color", false, "Disable colored output")

	// Register command with root
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	file1Path := args[0]
	file2Path := args[1]

	// Parse flags
	unified, _ := cmd.Flags().GetBool("unified")
	format, _ := cmd.Flags().GetString("format")
	tolerance, _ := cmd.Flags().GetFloat64("tolerance")
	noColor, _ := cmd.Flags().GetBool("no-color")

	// Validate format flag
	if format != "text" && format != "json" {
		return fmt.Errorf("invalid format: %q (must be 'text' or 'json')", format)
	}

	// Validate files exist
	if _, err := os.Stat(file1Path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file1Path)
	}
	if _, err := os.Stat(file2Path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file2Path)
	}

	// Read both files
	logger.Debug("reading files", "file1", file1Path, "file2", file2Path)
	data1, err := os.ReadFile(file1Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1Path, err)
	}

	data2, err := os.ReadFile(file2Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2Path, err)
	}

	// Detect formats
	format1, err := detectFormat(file1Path)
	if err != nil {
		return fmt.Errorf("unable to detect format for %s: %w\nSupported formats: .np3, .xmp", file1Path, err)
	}
	logger.Debug("detected format", "file", file1Path, "format", format1)

	format2, err := detectFormat(file2Path)
	if err != nil {
		return fmt.Errorf("unable to detect format for %s: %w\nSupported formats: .np3, .xmp", file2Path, err)
	}
	logger.Debug("detected format", "file", file2Path, "format", format2)

	// Parse files to UniversalRecipe
	recipe1, err := parseDiffFile(data1, format1, file1Path)
	if err != nil {
		return err
	}

	recipe2, err := parseDiffFile(data2, format2, file2Path)
	if err != nil {
		return err
	}

	// Run diff comparison
	logger.Debug("comparing recipes", "tolerance", tolerance)
	results, err := inspect.Diff(recipe1, recipe2, tolerance)
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}

	// Check if any differences found (for exit code)
	hasDifferences := len(results) > 0

	// Format output
	if format == "json" {
		return outputDiffJSON(results, file1Path, file2Path, unified)
	}

	// Text format
	colorize := !noColor && shouldUseColor()
	output := inspect.FormatDiff(results, unified, colorize)
	fmt.Println(output)

	// Exit codes:
	// 0 = no differences
	// 1 = differences found
	// 2 = error (handled by cobra error return)
	if hasDifferences {
		os.Exit(1)
	}

	return nil
}

// parseDiffFile parses a file based on its format and returns a UniversalRecipe.
func parseDiffFile(data []byte, format, filePath string) (*models.UniversalRecipe, error) {
	logger.Debug("parsing file", "format", format, "size", len(data))

	var recipe *models.UniversalRecipe
	var err error

	switch format {
	case FormatNP3:
		recipe, err = np3.Parse(data)
	case FormatXMP:
		recipe, err = xmp.Parse(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w\nFile may be corrupted. Try re-exporting from source application.", filePath, err)
	}

	return recipe, nil
}

// outputDiffJSON formats and outputs diff results as JSON.
func outputDiffJSON(results []inspect.DiffResult, file1, file2 string, unified bool) error {
	// Count changes and significant changes
	changes := 0
	significantChanges := 0
	var unchanged []string

	for _, r := range results {
		if r.ChangeType != "unchanged" {
			changes++
			if r.Significant {
				significantChanges++
			}
		} else if unified {
			unchanged = append(unchanged, r.Field)
		}
	}

	// Create output structure
	output := inspect.DiffOutput{
		File1:              file1,
		File2:              file2,
		Changes:            changes,
		SignificantChanges: significantChanges,
		FieldsCompared:     len(results),
		Differences:        results,
	}

	if unified {
		output.Unchanged = unchanged
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// shouldUseColor determines if colored output should be used.
// Returns false if:
// - stdout is not a terminal (piped or redirected)
// - NO_COLOR environment variable is set
func shouldUseColor() bool {
	// Check NO_COLOR environment variable (Unix convention)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	return term.IsTerminal(int(os.Stdout.Fd()))
}
