//go:build ignore
// +build ignore

// test_profile_xmp.go - Generate profile-based XMP (like Adobe's Dream profile)
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

	fmt.Println("Generating Profile-based XMP with embedded 3D LUT...")
	fmt.Println("This creates a profile/look (like Adobe's Dream) that:")
	fmt.Println("  • Uses Nikon's Camera Flexible Color profile as the base")
	fmt.Println("  • Applies NP3 transformations via embedded 3D LUT")
	fmt.Println("  • Adds +1000K temperature offset to match NX Studio warmth")
	fmt.Println("  • Loads as a profile in Lightroom (not just a preset)")
	fmt.Println()

	// Generate profile-based XMP with Camera Flexible Color + temperature compensation
	// Temperature offset of +1000K compensates for baseline difference between
	// Lightroom's Camera Flexible Color and NX Studio's native rendering
	xmpData, err := xmp.GenerateProfileWithLUT(recipe, "", 1000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating profile XMP: %v\n", err)
		os.Exit(1)
	}

	// Write XMP to file
	outputPath := "testdata/FIlmstill's Nostalgic Negative - PROFILE.xmp"
	if err := os.WriteFile(outputPath, xmpData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing XMP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Profile-based XMP generated successfully!\n")
	fmt.Printf("✓ Output: %s\n", outputPath)
	fmt.Printf("✓ Size: %.1f KB\n", float64(len(xmpData))/1024.0)
	fmt.Println()
	fmt.Println("How this differs from the previous version:")
	fmt.Println("  • Specifies Camera Flexible Color as base profile")
	fmt.Println("  • Uses Nikon's color science instead of Adobe's")
	fmt.Println("  • Adds +1000K temperature to compensate for baseline warmth difference")
	fmt.Println("  • Should achieve 90-95% accuracy vs NX Studio")
	fmt.Println()
	fmt.Println("To test in Lightroom:")
	fmt.Println("  1. Copy this XMP to Lightroom's CameraRaw/Settings folder")
	fmt.Println("  2. Restart Lightroom")
	fmt.Println("  3. Apply as a profile (not a preset) to a Nikon Z f RAW file")
	fmt.Println("  4. Compare with NX Studio rendering")
}
