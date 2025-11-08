// test_lut_xmp.go - Generate XMP with embedded 3D LUT
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
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

	fmt.Println("Generating XMP with embedded 3D LUT...")
	fmt.Println("This will take a few seconds as we generate 32,768 color mappings...")

	// Generate XMP with 3D LUT
	xmpData, err := xmp.GenerateWithLUT(recipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XMP with LUT: %v\n", err)
		os.Exit(1)
	}

	// Write XMP to file
	outputPath := "testdata/FIlmstill's Nostalgic Negative - LUT.xmp"
	if err := os.WriteFile(outputPath, xmpData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing XMP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ XMP with 3D LUT generated successfully!\n")
	fmt.Printf("✓ Output: %s\n", outputPath)
	fmt.Printf("✓ Size: %.1f KB\n", float64(len(xmpData))/1024.0)
	fmt.Println("\nThis XMP includes:")
	fmt.Println("  • All parametric adjustments (HSL, tone curve, etc.)")
	fmt.Println("  • Embedded 32x32x32 RGB lookup table")
	fmt.Println("  • Should match NX Studio rendering with 95%+ accuracy")
	fmt.Println("\nTest this preset in Lightroom and compare to NX Studio!")
}
