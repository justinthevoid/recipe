package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/justin/recipe/internal/converter"
	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
	"github.com/justin/recipe/internal/models"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert [input]",
	Short: "Convert a preset file between formats",
	Long: `Convert photo presets between NP3 and XMP formats.

The CLI auto-detects the source format from the file extension.
You must specify the target format with --to.

Examples:
  recipe convert portrait.xmp --to np3
  recipe convert portrait.np3 --to xmp --output custom.xmp`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Parse flags
	toFormat, _ := cmd.Flags().GetString("to")
	fromFormat, _ := cmd.Flags().GetString("from")
	outputPath, _ := cmd.Flags().GetString("output")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	jsonMode := isJSONMode(cmd)

	// Validate target format (AC-7)
	if err := validateFormat(toFormat); err != nil {
		return err
	}

	// Start timing for performance logging (AC-6)
	start := time.Now()

	// Initialize result structure (AC-2)
	result := ConversionResult{
		Input:        inputPath,
		TargetFormat: toFormat,
	}

	// Auto-detect source format if not specified (AC-2)
	if fromFormat == "" {
		logger.Debug("detecting format", "file", inputPath)
		var err error
		fromFormat, err = detectFormat(inputPath)
		if err != nil {
			// AC-3: Failed conversions still output valid JSON
			result.Error = fmt.Sprintf("auto-detect failed: %v", err)
			result.DurationMs = time.Since(start).Milliseconds()
			outputConversionResult(result, jsonMode)
			return fmt.Errorf("auto-detect failed: %w", err)
		}
		logger.Debug("detected format", "format", fromFormat, "file", inputPath)
	} else {
		// Validate explicit source format
		if err := validateFormat(fromFormat); err != nil {
			result.Error = fmt.Sprintf("invalid source format: %v", err)
			result.SourceFormat = fromFormat
			result.DurationMs = time.Since(start).Milliseconds()
			outputConversionResult(result, jsonMode)
			return err
		}
	}
	result.SourceFormat = fromFormat

	// Generate output path if not specified (AC-3)
	if outputPath == "" {
		outputPath = generateOutputPath(inputPath, toFormat)
	}
	result.Output = outputPath

	// Check overwrite protection (AC-4)
	if err := checkOutputExists(outputPath, overwrite); err != nil {
		result.Error = err.Error()
		result.DurationMs = time.Since(start).Milliseconds()
		outputConversionResult(result, jsonMode)
		return err
	}

	// Read input file (AC-5)
	logger.Debug("reading input", "file", inputPath)
	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read input file: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		outputConversionResult(result, jsonMode)
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse and log parameter extraction (AC-4)
	// Note: We parse separately here for logging purposes only.
	// The actual conversion still happens via converter.Convert() which does its own parsing.
	logger.Debug("parsing file", "format", fromFormat, "file", inputPath)
	recipe, parseErr := parseForLogging(inputBytes, fromFormat)
	if parseErr == nil && recipe != nil {
		paramCount := countParameters(recipe)
		logger.Debug("extracted parameters",
			"count", paramCount,
			"format", fromFormat)

		if paramCount > 0 {
			paramSummary := formatParameterSummary(recipe, 10)
			logger.Debug("key parameters", "summary", paramSummary)
		}
	}

	// Convert (single API call to converter - AC-8)
	logger.Debug("converting formats", "from", fromFormat, "to", toFormat)
	outputBytes, err := converter.Convert(inputBytes, fromFormat, toFormat)

	if err != nil {
		// AC-3: Handle conversion errors (invalid format, parse errors)
		result.Error = fmt.Sprintf("conversion failed: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		outputConversionResult(result, jsonMode)
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Log generation (AC-3)
	logger.Debug("generating output", "format", toFormat)

	// Ensure output directory exists (AC-3)
	if err := ensureOutputDir(outputPath); err != nil {
		result.Error = fmt.Sprintf("failed to create output directory: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		outputConversionResult(result, jsonMode)
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write output file (AC-1)
	logger.Debug("writing output", "file", outputPath)
	if err := os.WriteFile(outputPath, outputBytes, 0644); err != nil {
		result.Error = fmt.Sprintf("failed to write output file: %v", err)
		result.DurationMs = time.Since(start).Milliseconds()
		outputConversionResult(result, jsonMode)
		return fmt.Errorf("failed to write output file: %w", err)
	}

	// Success! (AC-2)
	result.Success = true
	result.FileSizeBytes = int64(len(outputBytes))
	result.DurationMs = time.Since(start).Milliseconds()

	// Log completion with timing (AC-6)
	logger.Info("conversion completed",
		"file", outputPath,
		"duration_ms", result.DurationMs,
		"from", fromFormat,
		"to", toFormat)

	// Output result (AC-2, AC-7)
	outputConversionResult(result, jsonMode)

	return nil
}

// generateOutputPath creates output path from input path and target format.
// Example: "portrait.xmp" + "np3" -> "portrait.np3"
func generateOutputPath(inputPath, targetFormat string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + "." + targetFormat
}

// checkOutputExists returns error if file exists and overwrite is false.
// Used for overwrite protection (AC-4).
func checkOutputExists(outputPath string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("output file already exists: %s (use --overwrite to replace)", outputPath)
		}
	}
	return nil
}

// ensureOutputDir creates output directory if it doesn't exist.
// Returns error if directory creation fails.
func ensureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	return os.MkdirAll(dir, 0755)
}

// formatBytes converts bytes to human-readable format (AC-9).
// Examples: 1234 → "1.2 KB", 1048576 → "1.0 MB"
func formatBytes(bytes int) string {
	const kb = 1024
	const mb = kb * 1024

	if bytes < kb {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < mb {
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	}
}

// formatDuration converts duration to human-readable format (AC-9).
// Examples: 15ms, 1.23s
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Required flags
	convertCmd.Flags().StringP("to", "t", "", "Target format (required): np3 or xmp")
	convertCmd.MarkFlagRequired("to")

	// Optional flags
	convertCmd.Flags().StringP("from", "f", "", "Source format (auto-detected if omitted)")
	convertCmd.Flags().StringP("output", "o", "", "Output file path (default: replace input extension)")
	convertCmd.Flags().Bool("overwrite", false, "Overwrite existing output file")
}

// countParameters counts the number of non-zero fields in a UniversalRecipe.
// Used for verbose logging to show how many parameters were extracted (AC-4).
func countParameters(recipe *models.UniversalRecipe) int {
	if recipe == nil {
		return 0
	}

	count := 0
	val := reflect.ValueOf(*recipe)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip metadata fields
		if fieldType.Name == "SourceFormat" || fieldType.Name == "Name" ||
			strings.HasPrefix(fieldType.Name, "NP3") || fieldType.Name == "Metadata" {
			continue
		}

		// Check if field has non-zero value
		if !isZeroValue(field) {
			count++
		}
	}

	return count
}

// isZeroValue checks if a reflect.Value is its zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// For structs, check if all fields are zero
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// formatParameterSummary creates a summary string of key parameters (AC-4).
// Shows first 'limit' parameters with their values for verbose logging.
func formatParameterSummary(recipe *models.UniversalRecipe, limit int) string {
	if recipe == nil {
		return ""
	}

	params := []string{}

	// Collect key parameters with values
	if recipe.Exposure != 0 {
		params = append(params, fmt.Sprintf("Exposure=%.1f", recipe.Exposure))
	}
	if recipe.Contrast != 0 {
		params = append(params, fmt.Sprintf("Contrast=%+d", recipe.Contrast))
	}
	if recipe.Highlights != 0 {
		params = append(params, fmt.Sprintf("Highlights=%+d", recipe.Highlights))
	}
	if recipe.Shadows != 0 {
		params = append(params, fmt.Sprintf("Shadows=%+d", recipe.Shadows))
	}
	if recipe.Saturation != 0 {
		params = append(params, fmt.Sprintf("Saturation=%+d", recipe.Saturation))
	}
	if recipe.Vibrance != 0 {
		params = append(params, fmt.Sprintf("Vibrance=%+d", recipe.Vibrance))
	}
	if recipe.Clarity != 0 {
		params = append(params, fmt.Sprintf("Clarity=%+d", recipe.Clarity))
	}
	if recipe.Sharpness != 0 {
		params = append(params, fmt.Sprintf("Sharpness=%d", recipe.Sharpness))
	}
	if recipe.Temperature != nil {
		params = append(params, fmt.Sprintf("Temperature=%dK", *recipe.Temperature))
	}
	if recipe.Tint != 0 {
		params = append(params, fmt.Sprintf("Tint=%+d", recipe.Tint))
	}

	// Limit to specified number
	if len(params) > limit {
		extra := len(params) - limit
		params = params[:limit]
		params = append(params, fmt.Sprintf("... (%d more)", extra))
	}

	if len(params) == 0 {
		return "(no major adjustments)"
	}

	return strings.Join(params, ", ")
}

// parseForLogging parses input bytes into UniversalRecipe for logging purposes only.
// This allows us to log parameter details without modifying the core converter.
// Returns nil on error (parsing errors are not critical for logging).
func parseForLogging(input []byte, format string) (*models.UniversalRecipe, error) {
	switch format {
	case converter.FormatNP3:
		return np3.Parse(input)
	case converter.FormatXMP:
		return xmp.Parse(input)
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}
