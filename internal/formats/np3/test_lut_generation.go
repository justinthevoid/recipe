// test_lut_generation.go - Test 3D LUT generation from NP3
package main

import (
	"fmt"
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

	fmt.Println("Generating 3D LUT...")
	fmt.Printf("  Grid size: %dx%dx%d = %d color mappings\n", lut.LUTSize, lut.LUTSize, lut.LUTSize, lut.LUTPoints)

	// Generate 3D LUT
	lutData, err := lut.Generate3DLUT(recipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating LUT: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ LUT generated: %d bytes\n", len(lutData))

	// Compress and encode
	fmt.Println("\nCompressing and encoding LUT...")
	tableID, encoded, err := lut.CompressAndEncodeLUT(lutData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing LUT: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Table ID: %s\n", tableID)
	fmt.Printf("✓ Encoded size: %d characters\n", len(encoded))
	fmt.Printf("✓ Compression ratio: %.1f%%\n", 100.0*float64(len(encoded))/float64(len(lutData)))

	fmt.Println("\nEncoded data preview (first 100 chars):")
	if len(encoded) > 100 {
		fmt.Println(encoded[:100] + "...")
	} else {
		fmt.Println(encoded)
	}

	fmt.Println("\n✓ 3D LUT generation successful!")
	fmt.Println("Next: Integrate into XMP generator")
}
