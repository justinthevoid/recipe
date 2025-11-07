package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/converter"
	"github.com/justin/recipe/internal/formats/lrtemplate"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/inspect"
	"github.com/justin/recipe/internal/models"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [file]",
	Short: "Extract and display preset parameters as JSON or binary hex dump",
	Long: `Inspect parses a preset file and outputs all parameters as JSON or binary hex dump.

The inspect command provides a way to analyze preset parameters programmatically,
validate conversions, and learn how presets work internally.

Supports NP3, XMP, and lrtemplate formats with automatic format detection.

JSON Output (default):
  - Metadata (source file, format, timestamp, version)
  - All parameters from the UniversalRecipe model (50+ fields)

Binary Mode (NP3 only):
  - Annotated hex dump with byte offsets
  - Field labels and human-readable values
  - Useful for reverse engineering and parser validation

Examples:
  recipe inspect portrait.np3
  recipe inspect portrait.np3 --output portrait.json
  recipe inspect portrait.xmp | jq '.parameters.contrast'
  recipe inspect portrait.np3 --binary
  recipe inspect portrait.np3 --binary --output hex_dump.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Optional flags
	inspectCmd.Flags().StringP("output", "o", "", "Write output to file instead of stdout")
	inspectCmd.Flags().Bool("pretty", true, "Pretty-print JSON (default: true)")
	inspectCmd.Flags().Bool("binary", false, "Show hex dump with field annotations (NP3 only)")
}

func runInspect(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Parse flags
	outputPath, _ := cmd.Flags().GetString("output")
	binaryMode, _ := cmd.Flags().GetBool("binary")

	// Read input file (AC-5)
	logger.Debug("reading input", "file", inputPath)
	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s: %w", inputPath, err)
	}

	// Auto-detect format (AC-5)
	logger.Debug("detecting format", "file", inputPath)
	format, err := detectFormat(inputPath)
	if err != nil {
		return fmt.Errorf("unable to detect format for '%s'\nSupported formats: .np3, .xmp, .lrtemplate", inputPath)
	}
	logger.Debug("detected format", "format", format, "file", inputPath)

	var outputBytes []byte

	// Binary mode (AC-3, AC-4)
	if binaryMode {
		logger.Debug("generating binary hex dump", "format", format)

		// Binary dump (validates NP3-only internally)
		hexDump, err := inspect.BinaryDump(inputBytes, format)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		outputBytes = []byte(hexDump)
	} else {
		// JSON mode (default)
		// Parse based on format (AC-3)
		logger.Debug("parsing file", "format", format, "file", inputPath)
		recipe, err := parseFile(inputBytes, format)
		if err != nil {
			// Wrap with ConversionError for consistent error handling (AC-6)
			return &converter.ConversionError{
				Operation: "parse",
				Format:    format,
				Cause:     err,
			}
		}

		// Generate JSON output with metadata (AC-1, AC-2)
		logger.Debug("generating JSON output")
		jsonBytes, err := inspect.ToJSONWithMetadata(recipe, inputPath, format, version)
		if err != nil {
			return fmt.Errorf("failed to generate JSON: %w", err)
		}

		outputBytes = jsonBytes
	}

	// Output to file or stdout (AC-4)
	if outputPath != "" {
		// Write to file
		logger.Debug("writing output", "file", outputPath)

		// Create parent directories if needed (AC-4)
		if err := ensureOutputDir(outputPath); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Write file with 0644 permissions (AC-4)
		if err := os.WriteFile(outputPath, outputBytes, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		// Success message to stderr (AC-4)
		if binaryMode {
			fmt.Fprintf(os.Stderr, "✓ Binary dump saved to %s\n", outputPath)
		} else {
			fmt.Fprintf(os.Stderr, "✓ Saved to %s\n", outputPath)
		}

		logger.Info("inspection completed", "file", outputPath, "mode", binaryModeStr(binaryMode))
	} else {
		// Write to stdout
		fmt.Print(string(outputBytes))

		logger.Debug("inspection completed", "output", "stdout", "mode", binaryModeStr(binaryMode))
	}

	return nil
}

// binaryModeStr returns a string representation of binary mode for logging.
func binaryModeStr(binary bool) string {
	if binary {
		return "binary"
	}
	return "json"
}

// parseFile parses input bytes based on detected format.
// Returns UniversalRecipe or error (AC-3).
func parseFile(input []byte, format string) (*models.UniversalRecipe, error) {
	switch format {
	case converter.FormatNP3:
		return np3.Parse(input)
	case converter.FormatXMP:
		return xmp.Parse(input)
	case converter.FormatLRTemplate:
		return lrtemplate.Parse(input)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
