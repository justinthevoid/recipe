// test_gray_transform.go - Test how neutral gray is transformed
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
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
	_, err = np3.Parse(np3Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing NP3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Testing Neutral Gray Transformation ===")
	fmt.Println()

	// Test with neutral gray
	r, g, b := 0.5, 0.5, 0.5
	fmt.Printf("Input RGB: [%.4f, %.4f, %.4f]\n", r, g, b)

	// Convert to HSL
	color := colorful.Color{R: r, G: g, B: b}
	h, s, l := color.Hsl()
	fmt.Printf("Input HSL: [%.2f°, %.4f, %.4f]\n", h, s, l)
	fmt.Println()

	// The problem: When saturation is 0 (neutral gray), the hue is undefined!
	// Different libraries may return different hue values for neutral colors
	if s < 0.01 {
		fmt.Println("⚠ WARNING: Saturation is near zero - hue is undefined!")
		fmt.Printf("  Library returned hue=%.2f°, but this is meaningless for gray\n", h)
		fmt.Println("  Applying color-specific adjustments to gray will produce incorrect results")
		fmt.Println()
	}

	// Test what happens when we convert back
	outColor := colorful.Hsl(h, s, l)
	fmt.Printf("Roundtrip RGB: [%.4f, %.4f, %.4f]\n", outColor.R, outColor.G, outColor.B)
	fmt.Println()

	// Now let's trace through the color transformation
	fmt.Println("=== Applying Color Blender to Gray ===")
	fmt.Println("The issue: Color Blender applies adjustments based on hue,")
	fmt.Println("even though gray has no hue (saturation = 0)")
	fmt.Println()

	// The green adjustment has huge negative luminance
	fmt.Printf("Green adjustment: L=-76 (darkens green by 76%%)\n")
	fmt.Printf("If the undefined hue happens to be near green (120°), gray gets darkened!\n")
	fmt.Println()

	// Test with a defined green color
	r, g, b = 0.0, 0.5, 0.0
	color = colorful.Color{R: r, G: g, B: b}
	h, s, l = color.Hsl()
	fmt.Printf("Pure Green RGB: [%.4f, %.4f, %.4f]\n", r, g, b)
	fmt.Printf("Pure Green HSL: [%.2f°, %.4f, %.4f]\n", h, s, l)
	fmt.Println("This would correctly receive the green luminance adjustment")
	fmt.Println()

	fmt.Println("=== The Solution ===")
	fmt.Println("We need to skip color-specific adjustments when saturation is very low")
	fmt.Println("(e.g., s < 0.01) because the hue is undefined for neutral colors.")
	fmt.Println()
	fmt.Println("This affects:")
	fmt.Println("  • Color Blender (8-color HSL adjustments)")
	fmt.Println("  • Color Grading zone hue shifts")
	fmt.Println()
	fmt.Println("For neutral/near-neutral colors, we should only apply:")
	fmt.Println("  • Tone curve (contrast, highlights, shadows)")
	fmt.Println("  • Luminance/brightness adjustments")
	fmt.Println("  • Global saturation (which will keep them near-neutral)")
}
