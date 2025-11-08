//go:build ignore
// +build ignore

// test_no_color_grading.go - Generate XMP without color grading zones to test
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/converter"
	"github.com/justin/recipe/internal/formats/np3"
)

func main() {
	// Read original NP3
	originalNP3, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading NP3: %v\n", err)
		os.Exit(1)
	}

	// Parse to UniversalRecipe
	recipe, err := np3.Parse(originalNP3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing NP3: %v\n", err)
		os.Exit(1)
	}

	// Remove color grading zones (keep blending/balance)
	fmt.Println("Original Color Grading:")
	fmt.Printf("  Highlights: Hue=%d Chroma=%d Brightness=%d\n",
		recipe.ColorGrading.Highlights.Hue,
		recipe.ColorGrading.Highlights.Chroma,
		recipe.ColorGrading.Highlights.Brightness)
	fmt.Printf("  Midtone: Hue=%d Chroma=%d Brightness=%d\n",
		recipe.ColorGrading.Midtone.Hue,
		recipe.ColorGrading.Midtone.Chroma,
		recipe.ColorGrading.Midtone.Brightness)
	fmt.Printf("  Shadows: Hue=%d Chroma=%d Brightness=%d\n",
		recipe.ColorGrading.Shadows.Hue,
		recipe.ColorGrading.Shadows.Chroma,
		recipe.ColorGrading.Shadows.Brightness)

	// Zero out color grading zones
	recipe.ColorGrading.Highlights.Hue = 0
	recipe.ColorGrading.Highlights.Chroma = 0
	recipe.ColorGrading.Highlights.Brightness = 0
	recipe.ColorGrading.Midtone.Hue = 0
	recipe.ColorGrading.Midtone.Chroma = 0
	recipe.ColorGrading.Midtone.Brightness = 0
	recipe.ColorGrading.Shadows.Hue = 0
	recipe.ColorGrading.Shadows.Chroma = 0
	recipe.ColorGrading.Shadows.Brightness = 0

	fmt.Println("\nGenerating XMP without color grading zones...")

	// Convert to XMP via the converter
	// (We need to convert back to NP3 first, then to XMP)
	np3Data, err := converter.Convert(nil, "universal", "np3")
	if err != nil {
		// Direct generation not supported, use the recipe's Generate
		fmt.Println("Using direct np3.Generate...")
		np3Data, err = np3.Generate(recipe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating NP3: %v\n", err)
			os.Exit(1)
		}
	}

	xmpData, err := converter.Convert(np3Data, "np3", "xmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to XMP: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	outputPath := "testdata/FIlmstill's Nostalgic Negative - NO_COLOR_GRADING.xmp"
	if err := os.WriteFile(outputPath, xmpData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing XMP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ XMP without color grading zones written to: %s\n", outputPath)
	fmt.Println("This version only has Color Blender (8-color HSL) but no zone-based grading")
}
