// Package converter provides a unified API for bidirectional photo editing recipe conversions.
//
// The converter package orchestrates conversions between all supported formats (NP3, XMP, lrtemplate)
// using the UniversalRecipe hub-and-spoke pattern. It provides format auto-detection, error handling
// with operation context, and thread-safe stateless conversion functions.
//
// # Supported Formats
//
// The package supports three photo editing recipe formats:
//   - NP3: NikonCapture NX binary format (300 bytes, magic bytes "NCP")
//   - XMP: Adobe XMP sidecar format (XML with Camera Raw Settings namespace)
//   - lrtemplate: Adobe Lightroom template format (Lua table syntax)
//
// All 6 conversion paths are supported: NP3↔XMP, NP3↔lrtemplate, XMP↔lrtemplate
//
// # Basic Usage
//
// Convert between formats with explicit format specification:
//
//	np3Data, err := os.ReadFile("recipe.np3")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	xmpData, err := converter.Convert(np3Data, converter.FormatNP3, converter.FormatXMP)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err = os.WriteFile("recipe.xmp", xmpData, 0644)
//
// # Auto-Detection
//
// Use format auto-detection by passing an empty string as the source format:
//
//	data, err := os.ReadFile("unknown_recipe")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Auto-detect source format and convert to XMP
//	xmpData, err := converter.Convert(data, "", converter.FormatXMP)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Or detect format explicitly:
//
//	format, err := converter.DetectFormat(data)
//	if err != nil {
//	    log.Fatalf("Unknown format: %v", err)
//	}
//	fmt.Printf("Detected format: %s\n", format)
//
// # Error Handling
//
// The package returns structured ConversionError with operation context:
//
//	xmpData, err := converter.Convert(np3Data, converter.FormatNP3, converter.FormatXMP)
//	if err != nil {
//	    var convErr *converter.ConversionError
//	    if errors.As(err, &convErr) {
//	        // Structured error with operation context
//	        fmt.Printf("Failed during %s of %s format: %v\n",
//	            convErr.Operation, convErr.Format, convErr.Cause)
//
//	        // Check for non-fatal warnings
//	        if len(convErr.Warnings) > 0 {
//	            fmt.Printf("Warnings: %d unmappable parameters\n", len(convErr.Warnings))
//	        }
//	    }
//	    return err
//	}
//
// Error operations include: "validate", "detect", "parse", "generate"
//
// # Thread Safety
//
// All conversion functions are thread-safe and stateless. Multiple goroutines can
// safely perform concurrent conversions without synchronization:
//
//	var wg sync.WaitGroup
//	for i := 0; i < 100; i++ {
//	    wg.Add(1)
//	    go func(data []byte) {
//	        defer wg.Done()
//	        result, err := converter.Convert(data, converter.FormatNP3, converter.FormatXMP)
//	        if err != nil {
//	            log.Printf("Conversion failed: %v", err)
//	            return
//	        }
//	        // Process result...
//	    }(np3Data)
//	}
//	wg.Wait()
//
// No shared state is modified during conversion. Each call operates independently
// on its input data.
//
// # Performance Characteristics
//
// The package is designed for high-performance batch conversions:
//
//   - Target: <100ms per conversion (typical: <15ms)
//   - Hub overhead: <5ms (typical: <1μs)
//   - Format detection: <1ms (typical: <100ns)
//   - Memory: Allocations proportional to file size
//   - Concurrency: Scales linearly with CPU cores
//
// Benchmarks on AMD Ryzen 9 7900X (12-core):
//
//   - Format detection: 3-54ns depending on format
//   - Auto-detect conversion: ~11μs (including parse + generate)
//   - Thread-safety: 100 concurrent conversions complete in <5ms
//
// For optimal performance:
//   - Reuse input buffers when possible
//   - Use explicit format specification to skip auto-detection
//   - Leverage goroutines for batch conversions
//   - Profile with go test -bench to validate performance targets
package converter

import "fmt"

// ConversionError wraps conversion failures with operation context for debugging and user feedback.
// It provides structured error information including the failed operation, target format, root cause,
// and any non-fatal warnings encountered during conversion.
type ConversionError struct {
	// Operation describes what step failed: "validate", "detect", "parse", or "generate"
	Operation string

	// Format specifies which format was being processed when the error occurred
	Format string

	// Cause is the underlying error that triggered the conversion failure
	Cause error

	// Warnings contains non-fatal issues encountered (e.g., unmappable parameters)
	// These are informational only and don't prevent conversion completion
	Warnings []string
}

// Error implements the error interface, providing a human-readable error message
// that includes operation context, format information, and the root cause.
func (e *ConversionError) Error() string {
	msg := fmt.Sprintf("%s %s: %v", e.Operation, e.Format, e.Cause)
	if len(e.Warnings) > 0 {
		msg += fmt.Sprintf(" (warnings: %d unmappable parameters)", len(e.Warnings))
	}
	return msg
}

// Unwrap returns the underlying error, enabling Go 1.13+ error chain unwrapping
// with errors.Is() and errors.As() for programmatic error handling.
func (e *ConversionError) Unwrap() error {
	return e.Cause
}
