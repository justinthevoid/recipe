// test_all_parameters.go - Extract ALL parameters from NP3 to find discrepancies
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
)

func main() {
	data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	recipe, err := np3.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("COMPLETE NP3 PARAMETER EXTRACTION")
	fmt.Println("==================================")

	// Basic adjustments
	fmt.Println("\nBASIC ADJUSTMENTS:")
	fmt.Printf("  Name:       %s\n", recipe.Name)
	fmt.Printf("  Sharpness:  %d\n", recipe.Sharpness)
	fmt.Printf("  Clarity:    %d\n", recipe.Clarity)

	// Advanced adjustments
	fmt.Println("\nADVANCED ADJUSTMENTS:")
	fmt.Printf("  MidRangeSharpening: %d\n", recipe.MidRangeSharpening)
	fmt.Printf("  Contrast:   %d\n", recipe.Contrast)
	fmt.Printf("  Highlights: %d\n", recipe.Highlights)
	fmt.Printf("  Shadows:    %d\n", recipe.Shadows)
	fmt.Printf("  Whites:     %d\n", recipe.Whites)
	fmt.Printf("  Blacks:     %d\n", recipe.Blacks)
	fmt.Printf("  Saturation: %d\n", recipe.Saturation)
	fmt.Printf("  Vibrance:   %d\n", recipe.Vibrance)

	// Color Blender (8 colors)
	fmt.Println("\nCOLOR BLENDER (HSL for 8 colors):")
	fmt.Printf("  Red:     Hue=%3d Sat=%3d Lum=%3d\n", recipe.Red.Hue, recipe.Red.Saturation, recipe.Red.Luminance)
	fmt.Printf("  Orange:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Orange.Hue, recipe.Orange.Saturation, recipe.Orange.Luminance)
	fmt.Printf("  Yellow:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Yellow.Hue, recipe.Yellow.Saturation, recipe.Yellow.Luminance)
	fmt.Printf("  Green:   Hue=%3d Sat=%3d Lum=%3d\n", recipe.Green.Hue, recipe.Green.Saturation, recipe.Green.Luminance)
	fmt.Printf("  Aqua:    Hue=%3d Sat=%3d Lum=%3d\n", recipe.Aqua.Hue, recipe.Aqua.Saturation, recipe.Aqua.Luminance)
	fmt.Printf("  Blue:    Hue=%3d Sat=%3d Lum=%3d\n", recipe.Blue.Hue, recipe.Blue.Saturation, recipe.Blue.Luminance)
	fmt.Printf("  Purple:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Purple.Hue, recipe.Purple.Saturation, recipe.Purple.Luminance)
	fmt.Printf("  Magenta: Hue=%3d Sat=%3d Lum=%3d\n", recipe.Magenta.Hue, recipe.Magenta.Saturation, recipe.Magenta.Luminance)

	// Color Grading (3 zones)
	fmt.Println("\nCOLOR GRADING (3 zones - Flexible Color Picture Control):")
	fmt.Printf("  Highlights: Hue=%3d Chroma=%3d Brightness=%3d\n",
		recipe.ColorGrading.Highlights.Hue,
		recipe.ColorGrading.Highlights.Chroma,
		recipe.ColorGrading.Highlights.Brightness)
	fmt.Printf("  Midtone:    Hue=%3d Chroma=%3d Brightness=%3d\n",
		recipe.ColorGrading.Midtone.Hue,
		recipe.ColorGrading.Midtone.Chroma,
		recipe.ColorGrading.Midtone.Brightness)
	fmt.Printf("  Shadows:    Hue=%3d Chroma=%3d Brightness=%3d\n",
		recipe.ColorGrading.Shadows.Hue,
		recipe.ColorGrading.Shadows.Chroma,
		recipe.ColorGrading.Shadows.Brightness)
	fmt.Printf("  Blending:   %d (smoothness of zone transitions)\n", recipe.ColorGrading.Blending)
	fmt.Printf("  Balance:    %d (shifts overall color balance)\n", recipe.ColorGrading.Balance)

	fmt.Println("\n==================================")
	fmt.Println("Looking for discrepancies that might explain")
	fmt.Println("the 'more vibrant' and 'green tint' in NX Studio")
}
