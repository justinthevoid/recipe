//go:build ignore
// +build ignore

// test_validation_debug.go - Debug which validation triggers fallback
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Extract all parameters that are validated
	decodeSigned8 := func(b byte) int {
		return int(b) - 0x80
	}
	decodeScaled4 := func(b byte) float64 {
		return float64(int(b)-0x80) / 4.0
	}

	fmt.Println("Validation Debug for Nostalgic Negative")
	fmt.Println("========================================")

	// Basic Adjustments
	sharpening := decodeScaled4(data[82])
	clarity := decodeScaled4(data[92])
	fmt.Printf("Sharpening: %.2f (range: -3.0 to 9.0) - %s\n", sharpening,
		validRange(sharpening, -3.0, 9.0))
	fmt.Printf("Clarity:    %.2f (range: -5.0 to 5.0) - %s\n", clarity,
		validRange(clarity, -5.0, 5.0))

	// Advanced Adjustments
	midRangeSharpening := decodeScaled4(data[242])
	contrast := decodeSigned8(data[272])
	highlights := decodeSigned8(data[282])
	shadows := decodeSigned8(data[292])
	whiteLevel := decodeSigned8(data[302])
	blackLevel := decodeSigned8(data[312])
	saturation := decodeSigned8(data[322])

	fmt.Printf("MidRangeSharpening: %.2f (range: -5.0 to 5.0) - %s\n", midRangeSharpening,
		validRange(midRangeSharpening, -5.0, 5.0))
	fmt.Printf("Contrast:    %3d (range: -100 to 100) - %s\n", contrast,
		validRangeInt(contrast, -100, 100))
	fmt.Printf("Highlights:  %3d (range: -100 to 100) - %s\n", highlights,
		validRangeInt(highlights, -100, 100))
	fmt.Printf("Shadows:     %3d (range: -100 to 100) - %s\n", shadows,
		validRangeInt(shadows, -100, 100))
	fmt.Printf("WhiteLevel:  %3d (range: -100 to 100) - %s\n", whiteLevel,
		validRangeInt(whiteLevel, -100, 100))
	fmt.Printf("BlackLevel:  %3d (range: -100 to 100) - %s\n", blackLevel,
		validRangeInt(blackLevel, -100, 100))
	fmt.Printf("Saturation:  %3d (range: -100 to 100) - %s\n", saturation,
		validRangeInt(saturation, -100, 100))

	// Color Blender - check a few
	redHue := decodeSigned8(data[332])
	redChroma := decodeSigned8(data[333])
	redBrightness := decodeSigned8(data[334])
	fmt.Printf("\nRed Hue:        %3d (range: -100 to 100) - %s\n", redHue,
		validRangeInt(redHue, -100, 100))
	fmt.Printf("Red Chroma:     %3d (range: -100 to 100) - %s\n", redChroma,
		validRangeInt(redChroma, -100, 100))
	fmt.Printf("Red Brightness: %3d (range: -100 to 100) - %s\n", redBrightness,
		validRangeInt(redBrightness, -100, 100))

	// Color Grading
	decodeHue12 := func(b1, b2 byte) int {
		return (int(b1&0x0F) << 8) | int(b2)
	}

	highlightsHue := decodeHue12(data[368], data[369])
	highlightsChroma := decodeSigned8(data[370])
	highlightsBrightness := decodeSigned8(data[371])

	midtoneHue := decodeHue12(data[372], data[373])
	midtoneChroma := decodeSigned8(data[374])
	midtoneBrightness := decodeSigned8(data[375])

	shadowsHue := decodeHue12(data[376], data[377])
	shadowsChroma := decodeSigned8(data[378])
	shadowsBrightness := decodeSigned8(data[379])

	blending := int(data[384])
	balance := decodeSigned8(data[386])

	fmt.Printf("\nColor Grading Highlights: Hue=%3d Chroma=%3d Brightness=%3d\n",
		highlightsHue, highlightsChroma, highlightsBrightness)
	fmt.Printf("  Hue valid (0-360)?     %s\n", validRangeInt(highlightsHue, 0, 360))
	fmt.Printf("  Chroma valid (-100-100)? %s\n", validRangeInt(highlightsChroma, -100, 100))
	fmt.Printf("  Brightness valid?      %s\n", validRangeInt(highlightsBrightness, -100, 100))

	fmt.Printf("\nColor Grading Midtone: Hue=%3d Chroma=%3d Brightness=%3d\n",
		midtoneHue, midtoneChroma, midtoneBrightness)
	fmt.Printf("  Hue valid (0-360)?     %s\n", validRangeInt(midtoneHue, 0, 360))
	fmt.Printf("  Chroma valid (-100-100)? %s\n", validRangeInt(midtoneChroma, -100, 100))
	fmt.Printf("  Brightness valid?      %s\n", validRangeInt(midtoneBrightness, -100, 100))

	fmt.Printf("\nColor Grading Shadows: Hue=%3d Chroma=%3d Brightness=%3d\n",
		shadowsHue, shadowsChroma, shadowsBrightness)
	fmt.Printf("  Hue valid (0-360)?     %s\n", validRangeInt(shadowsHue, 0, 360))
	fmt.Printf("  Chroma valid (-100-100)? %s\n", validRangeInt(shadowsChroma, -100, 100))
	fmt.Printf("  Brightness valid?      %s\n", validRangeInt(shadowsBrightness, -100, 100))

	fmt.Printf("\nBlending: %3d (range: 0 to 100) - %s\n", blending,
		validRangeInt(blending, 0, 100))
	fmt.Printf("Balance:  %3d (range: -100 to 100) - %s\n", balance,
		validRangeInt(balance, -100, 100))

	fmt.Println("\n========================================")
	fmt.Println("Check for INVALID values above - those trigger needsFallback=true")
}

func validRange(v float64, min, max float64) string {
	if v < min || v > max {
		return fmt.Sprintf("INVALID (out of range)")
	}
	return "OK"
}

func validRangeInt(v, min, max int) string {
	if v < min || v > max {
		return fmt.Sprintf("INVALID (out of range)")
	}
	return "OK"
}
