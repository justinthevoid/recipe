//go:build ignore
// +build ignore

// test_recipe_values.go - Show actual NP3 recipe values
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
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

	fmt.Println("=== NP3 Recipe Values ===")
	fmt.Printf("Name: %s\n\n", recipe.Name)

	fmt.Println("Global Adjustments:")
	fmt.Printf("  Exposure:    %.2f\n", recipe.Exposure)
	fmt.Printf("  Contrast:    %d\n", recipe.Contrast)
	fmt.Printf("  Highlights:  %d\n", recipe.Highlights)
	fmt.Printf("  Shadows:     %d\n", recipe.Shadows)
	fmt.Printf("  Whites:      %d\n", recipe.Whites)
	fmt.Printf("  Blacks:      %d\n", recipe.Blacks)
	fmt.Printf("  Saturation:  %d\n", recipe.Saturation)
	fmt.Println()

	fmt.Println("Color Blender (8-color HSL):")
	fmt.Printf("  Red:     H=%3d  S=%3d  L=%3d\n", recipe.Red.Hue, recipe.Red.Saturation, recipe.Red.Luminance)
	fmt.Printf("  Orange:  H=%3d  S=%3d  L=%3d\n", recipe.Orange.Hue, recipe.Orange.Saturation, recipe.Orange.Luminance)
	fmt.Printf("  Yellow:  H=%3d  S=%3d  L=%3d\n", recipe.Yellow.Hue, recipe.Yellow.Saturation, recipe.Yellow.Luminance)
	fmt.Printf("  Green:   H=%3d  S=%3d  L=%3d\n", recipe.Green.Hue, recipe.Green.Saturation, recipe.Green.Luminance)
	fmt.Printf("  Aqua:    H=%3d  S=%3d  L=%3d\n", recipe.Aqua.Hue, recipe.Aqua.Saturation, recipe.Aqua.Luminance)
	fmt.Printf("  Blue:    H=%3d  S=%3d  L=%3d\n", recipe.Blue.Hue, recipe.Blue.Saturation, recipe.Blue.Luminance)
	fmt.Printf("  Purple:  H=%3d  S=%3d  L=%3d\n", recipe.Purple.Hue, recipe.Purple.Saturation, recipe.Purple.Luminance)
	fmt.Printf("  Magenta: H=%3d  S=%3d  L=%3d\n", recipe.Magenta.Hue, recipe.Magenta.Saturation, recipe.Magenta.Luminance)
	fmt.Println()

	if recipe.ColorGrading != nil {
		fmt.Println("Color Grading (3-zone):")
		fmt.Printf("  Balance:    %d\n", recipe.ColorGrading.Balance)
		fmt.Printf("  Shadows:    H=%3d  C=%3d  B=%3d\n",
			recipe.ColorGrading.Shadows.Hue,
			recipe.ColorGrading.Shadows.Chroma,
			recipe.ColorGrading.Shadows.Brightness)
		fmt.Printf("  Midtone:    H=%3d  C=%3d  B=%3d\n",
			recipe.ColorGrading.Midtone.Hue,
			recipe.ColorGrading.Midtone.Chroma,
			recipe.ColorGrading.Midtone.Brightness)
		fmt.Printf("  Highlights: H=%3d  C=%3d  B=%3d\n",
			recipe.ColorGrading.Highlights.Hue,
			recipe.ColorGrading.Highlights.Chroma,
			recipe.ColorGrading.Highlights.Brightness)
	} else {
		fmt.Println("Color Grading: Not present")
	}
}
