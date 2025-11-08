//go:build ignore
// +build ignore

// test_convert.go - Test full conversion with fix
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/converter"
)

func main() {
	// Read NP3 file
	np3Data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading NP3: %v\n", err)
		os.Exit(1)
	}

	// Convert NP3 to XMP
	xmpData, err := converter.Convert(np3Data, "np3", "xmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %v\n", err)
		os.Exit(1)
	}

	// Write XMP to file
	outputPath := "testdata/FIlmstill's Nostalgic Negative - FIXED.xmp"
	if err := os.WriteFile(outputPath, xmpData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing XMP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Conversion successful!\n")
	fmt.Printf("✓ Output written to: %s\n", outputPath)
	fmt.Printf("✓ File size: %d bytes\n", len(xmpData))
	fmt.Println("\nXMP Preview (first 1500 chars):")
	fmt.Println("================================")
	if len(xmpData) > 1500 {
		fmt.Println(string(xmpData[:1500]) + "...")
	} else {
		fmt.Println(string(xmpData))
	}
}
