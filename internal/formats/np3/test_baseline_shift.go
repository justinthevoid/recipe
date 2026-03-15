//go:build ignore
// +build ignore

// test_baseline_shift.go - Analyze the baseline color shift between NX Studio and Lightroom
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/lut"
	"github.com/lucasb-eyer/go-colorful"
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

	fmt.Println("=== Analyzing Color Temperature Shift ===")
	fmt.Println()
	fmt.Println("Based on your screenshot, Lightroom appears cooler/more cyan")
	fmt.Println("compared to NX Studio's warmer/golden rendering.")
	fmt.Println()

	// Test warm colors (sky tones in your image appear to be in this range)
	fmt.Println("Testing sky/sand tones that should be warm:")
	testColor := func(r, g, b float64, name string) {
		// Convert to HSL to analyze
		color := colorful.Color{R: r, G: g, B: b}
		h, s, l := color.Hsl()
		fmt.Printf("\n%s RGB[%.2f, %.2f, %.2f]:\n", name, r, g, b)
		fmt.Printf("  HSL: Hue=%.1f° Sat=%.2f Lum=%.2f\n", h, s, l)

		// Check if this would be affected by blue/cyan cast
		if h >= 180 && h <= 240 {
			fmt.Printf("  ⚠ This is in the blue/cyan range - would appear TOO BLUE\n")
		} else if h >= 30 && h <= 60 {
			fmt.Printf("  ✓ This is in the yellow/orange range - correct warmth\n")
		}
	}

	// Test typical sky/sand colors from a beach scene
	testColor(0.8, 0.75, 0.65, "Warm sand/sky")
	testColor(0.6, 0.65, 0.70, "Cool sky (problematic)")
	testColor(0.55, 0.50, 0.45, "Dry grass")

	fmt.Println()
	fmt.Println("=== Potential Solutions ===")
	fmt.Println()
	fmt.Println("The issue is that Camera Flexible Color in Lightroom has a different")
	fmt.Println("baseline color matrix than NX Studio's rendering engine.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("1. Add a global temperature shift in the XMP (e.g., +800-1200K)")
	fmt.Println("2. Apply a hue shift to specific color ranges (blues → yellows)")
	fmt.Println("3. Try a different Adobe camera profile as baseline")
	fmt.Println("4. Adjust XMP tone parameters to compensate")
	fmt.Println()
	fmt.Println("Most practical: Add Temperature adjustment to XMP")
	fmt.Println("This compensates for the baseline profile difference.")
}
