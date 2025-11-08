// test_roundtrip_colors.go - Verify color data preservation in full round-trip
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

	// Parse original to get baseline colors
	originalRecipe, err := np3.Parse(originalNP3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing original NP3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ORIGINAL NP3 Color Values:")
	fmt.Println("==========================")
	fmt.Printf("Red:     Hue=%3d Sat=%3d Lum=%3d\n", originalRecipe.Red.Hue, originalRecipe.Red.Saturation, originalRecipe.Red.Luminance)
	fmt.Printf("Orange:  Hue=%3d Sat=%3d Lum=%3d\n", originalRecipe.Orange.Hue, originalRecipe.Orange.Saturation, originalRecipe.Orange.Luminance)
	fmt.Printf("Yellow:  Hue=%3d Sat=%3d Lum=%3d\n", originalRecipe.Yellow.Hue, originalRecipe.Yellow.Saturation, originalRecipe.Yellow.Luminance)
	fmt.Printf("Blending: %d\n", originalRecipe.ColorGrading.Blending)
	fmt.Printf("Balance:  %d\n", originalRecipe.ColorGrading.Balance)

	// Convert NP3 → XMP
	xmpData, err := converter.Convert(originalNP3, "np3", "xmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to XMP: %v\n", err)
		os.Exit(1)
	}

	// Convert XMP → NP3
	finalNP3, err := converter.Convert(xmpData, "xmp", "np3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting back to NP3: %v\n", err)
		os.Exit(1)
	}

	// Parse final NP3
	finalRecipe, err := np3.Parse(finalNP3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing final NP3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nFINAL NP3 Color Values (after round-trip):")
	fmt.Println("==========================================")
	fmt.Printf("Red:     Hue=%3d Sat=%3d Lum=%3d\n", finalRecipe.Red.Hue, finalRecipe.Red.Saturation, finalRecipe.Red.Luminance)
	fmt.Printf("Orange:  Hue=%3d Sat=%3d Lum=%3d\n", finalRecipe.Orange.Hue, finalRecipe.Orange.Saturation, finalRecipe.Orange.Luminance)
	fmt.Printf("Yellow:  Hue=%3d Sat=%3d Lum=%3d\n", finalRecipe.Yellow.Hue, finalRecipe.Yellow.Saturation, finalRecipe.Yellow.Luminance)
	fmt.Printf("Blending: %d\n", finalRecipe.ColorGrading.Blending)
	fmt.Printf("Balance:  %d\n", finalRecipe.ColorGrading.Balance)

	// Compare
	fmt.Println("\nCOMPARISON:")
	fmt.Println("===========")
	redMatch := originalRecipe.Red.Hue == finalRecipe.Red.Hue &&
		originalRecipe.Red.Saturation == finalRecipe.Red.Saturation &&
		originalRecipe.Red.Luminance == finalRecipe.Red.Luminance
	orangeMatch := originalRecipe.Orange.Hue == finalRecipe.Orange.Hue &&
		originalRecipe.Orange.Saturation == finalRecipe.Orange.Saturation &&
		originalRecipe.Orange.Luminance == finalRecipe.Orange.Luminance
	yellowMatch := originalRecipe.Yellow.Hue == finalRecipe.Yellow.Hue &&
		originalRecipe.Yellow.Saturation == finalRecipe.Yellow.Saturation &&
		originalRecipe.Yellow.Luminance == finalRecipe.Yellow.Luminance
	blendingMatch := originalRecipe.ColorGrading.Blending == finalRecipe.ColorGrading.Blending
	balanceMatch := originalRecipe.ColorGrading.Balance == finalRecipe.ColorGrading.Balance

	fmt.Printf("Red colors match:     %v\n", redMatch)
	fmt.Printf("Orange colors match:  %v\n", orangeMatch)
	fmt.Printf("Yellow colors match:  %v\n", yellowMatch)
	fmt.Printf("Blending matches:     %v\n", blendingMatch)
	fmt.Printf("Balance matches:      %v\n", balanceMatch)

	if redMatch && orangeMatch && yellowMatch && blendingMatch && balanceMatch {
		fmt.Println("\n✓ SUCCESS: Color data preserved through round-trip!")
	} else {
		fmt.Println("\n✗ FAILURE: Color data lost in round-trip")
		os.Exit(1)
	}
}
