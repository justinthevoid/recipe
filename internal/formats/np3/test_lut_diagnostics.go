// test_lut_diagnostics.go - Diagnostic tool to investigate LUT generation issues
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/lut"
)

func main() {
	// Read NP3 file
	np3Data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading NP3: %v\n", err)
		os.Exit(1)
	}

	// Parse to UniversalRecipe
	recipe, err := np3.Parse(np3Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing NP3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== LUT Generation Diagnostics ===")
	fmt.Println()

	// Generate 3D LUT
	lutData, err := lut.Generate3DLUT(recipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating LUT: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ LUT generated successfully\n")
	fmt.Printf("  Size: %d bytes\n", len(lutData))
	fmt.Printf("  Expected size: %d bytes (32^3 × 3 channels × 4 bytes/float32)\n", 32*32*32*3*4)
	fmt.Println()

	// Verify size
	expectedSize := 32 * 32 * 32 * 3 * 4
	if len(lutData) != expectedSize {
		fmt.Printf("❌ SIZE MISMATCH: Expected %d bytes, got %d bytes\n", expectedSize, len(lutData))
		return
	}

	// Sample some LUT values to check for corruption
	fmt.Println("=== LUT Sample Values ===")
	fmt.Println("Testing neutral gray mapping (input [0.5, 0.5, 0.5] should stay close to gray):")
	fmt.Println()

	// Calculate index for mid-gray (16, 16, 16 in 32x32x32 cube)
	midIdx := 16
	offset := (midIdx*32*32 + midIdx*32 + midIdx) * 3 * 4

	r := readFloat32(lutData, offset)
	g := readFloat32(lutData, offset+4)
	b := readFloat32(lutData, offset+8)

	fmt.Printf("  Input:  [0.5, 0.5, 0.5] (mid-gray)\n")
	fmt.Printf("  Output: [%.4f, %.4f, %.4f]\n", r, g, b)
	fmt.Println()

	// Check if output is valid (should be in [0,1] range)
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
		fmt.Println("❌ OUTPUT OUT OF RANGE: Values should be in [0,1]")
	} else {
		fmt.Println("✓ Output values are in valid range [0,1]")
	}

	// Check if output maintains color balance (for neutral input)
	maxDiff := math.Max(math.Abs(float64(r-g)), math.Max(math.Abs(float64(g-b)), math.Abs(float64(r-b))))
	if maxDiff > 0.1 {
		fmt.Printf("⚠ WARNING: Color shift detected (max channel diff: %.4f)\n", maxDiff)
		fmt.Println("  Neutral gray input should produce near-neutral output")
	} else {
		fmt.Printf("✓ Color balance maintained (max channel diff: %.4f)\n", maxDiff)
	}
	fmt.Println()

	// Test extreme values
	fmt.Println("=== Testing Extreme Values ===")

	// Black (0,0,0)
	offset = 0
	r = readFloat32(lutData, offset)
	g = readFloat32(lutData, offset+4)
	b = readFloat32(lutData, offset+8)
	fmt.Printf("Black [0,0,0] → [%.4f, %.4f, %.4f]\n", r, g, b)

	// White (31,31,31)
	offset = (31*32*32 + 31*32 + 31) * 3 * 4
	r = readFloat32(lutData, offset)
	g = readFloat32(lutData, offset+4)
	b = readFloat32(lutData, offset+8)
	fmt.Printf("White [1,1,1] → [%.4f, %.4f, %.4f]\n", r, g, b)
	fmt.Println()

	// Test primary colors
	fmt.Println("=== Testing Primary Colors ===")

	// Pure red (31,0,0)
	offset = (0*32*32 + 0*32 + 31) * 3 * 4
	r = readFloat32(lutData, offset)
	g = readFloat32(lutData, offset+4)
	b = readFloat32(lutData, offset+8)
	fmt.Printf("Pure Red   [1,0,0] → [%.4f, %.4f, %.4f]\n", r, g, b)

	// Pure green (0,31,0)
	offset = (0*32*32 + 31*32 + 0) * 3 * 4
	r = readFloat32(lutData, offset)
	g = readFloat32(lutData, offset+4)
	b = readFloat32(lutData, offset+8)
	fmt.Printf("Pure Green [0,1,0] → [%.4f, %.4f, %.4f]\n", r, g, b)

	// Pure blue (0,0,31)
	offset = (31*32*32 + 0*32 + 0) * 3 * 4
	r = readFloat32(lutData, offset)
	g = readFloat32(lutData, offset+4)
	b = readFloat32(lutData, offset+8)
	fmt.Printf("Pure Blue  [0,0,1] → [%.4f, %.4f, %.4f]\n", r, g, b)
	fmt.Println()

	// Test compression and encoding
	fmt.Println("=== Testing Compression & Encoding ===")
	tableID, encodedLUT, err := lut.CompressAndEncodeLUT(lutData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing LUT: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ LUT compressed and encoded successfully\n")
	fmt.Printf("  Table ID (MD5): %s\n", tableID)
	fmt.Printf("  Encoded size: %d characters\n", len(encodedLUT))
	fmt.Printf("  Compression ratio: %.1f%% (original: %d bytes → %d chars)\n",
		float64(len(encodedLUT))*100/float64(len(lutData)), len(lutData), len(encodedLUT))
	fmt.Println()

	// Check if encoding contains only valid XML characters
	hasInvalidChars := false
	for _, ch := range encodedLUT {
		if ch < 32 && ch != '\n' && ch != '\r' && ch != '\t' {
			hasInvalidChars = true
			break
		}
	}

	if hasInvalidChars {
		fmt.Println("❌ ENCODING ERROR: Contains invalid XML characters")
	} else {
		fmt.Println("✓ Encoding uses valid XML characters")
	}

	fmt.Println()
	fmt.Println("=== Diagnostic Summary ===")
	fmt.Println("If you see any ❌ or ⚠ above, there may be issues with:")
	fmt.Println("  • LUT generation logic (color transformation formulas)")
	fmt.Println("  • Data structure/byte ordering")
	fmt.Println("  • Compression or encoding process")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review the sample values above")
	fmt.Println("  2. Compare with Adobe's reference profiles if possible")
	fmt.Println("  3. Test the warmth issue specifically with diagnostic color samples")
}

func readFloat32(data []byte, offset int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
}
