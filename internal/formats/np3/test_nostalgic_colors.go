// test_nostalgic_colors.go - Debug color extraction from Nostalgic Negative
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Color Blender offsets (332-355)
	fmt.Println("Color Blender Data (offsets 332-355):")
	fmt.Println("==========================================")

	// Helper function to decode Signed8 (byte - 0x80)
	decodeSigned8 := func(b byte) int {
		return int(b) - 0x80
	}

	// Red (332-334)
	if len(data) > 334 {
		redHue := decodeSigned8(data[332])
		redChroma := decodeSigned8(data[333])
		redBrightness := decodeSigned8(data[334])
		fmt.Printf("Red:    Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			redHue, redChroma, redBrightness, data[332], data[333], data[334])
	}

	// Orange (335-337)
	if len(data) > 337 {
		orangeHue := decodeSigned8(data[335])
		orangeChroma := decodeSigned8(data[336])
		orangeBrightness := decodeSigned8(data[337])
		fmt.Printf("Orange: Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			orangeHue, orangeChroma, orangeBrightness, data[335], data[336], data[337])
	}

	// Yellow (338-340)
	if len(data) > 340 {
		yellowHue := decodeSigned8(data[338])
		yellowChroma := decodeSigned8(data[339])
		yellowBrightness := decodeSigned8(data[340])
		fmt.Printf("Yellow: Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			yellowHue, yellowChroma, yellowBrightness, data[338], data[339], data[340])
	}

	// Green (341-343)
	if len(data) > 343 {
		greenHue := decodeSigned8(data[341])
		greenChroma := decodeSigned8(data[342])
		greenBrightness := decodeSigned8(data[343])
		fmt.Printf("Green:  Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			greenHue, greenChroma, greenBrightness, data[341], data[342], data[343])
	}

	// Cyan (344-346)
	if len(data) > 346 {
		cyanHue := decodeSigned8(data[344])
		cyanChroma := decodeSigned8(data[345])
		cyanBrightness := decodeSigned8(data[346])
		fmt.Printf("Cyan:   Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			cyanHue, cyanChroma, cyanBrightness, data[344], data[345], data[346])
	}

	// Blue (347-349)
	if len(data) > 349 {
		blueHue := decodeSigned8(data[347])
		blueChroma := decodeSigned8(data[348])
		blueBrightness := decodeSigned8(data[349])
		fmt.Printf("Blue:   Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			blueHue, blueChroma, blueBrightness, data[347], data[348], data[349])
	}

	// Purple (350-352)
	if len(data) > 352 {
		purpleHue := decodeSigned8(data[350])
		purpleChroma := decodeSigned8(data[351])
		purpleBrightness := decodeSigned8(data[352])
		fmt.Printf("Purple: Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			purpleHue, purpleChroma, purpleBrightness, data[350], data[351], data[352])
	}

	// Magenta (353-355)
	if len(data) > 355 {
		magentaHue := decodeSigned8(data[353])
		magentaChroma := decodeSigned8(data[354])
		magentaBrightness := decodeSigned8(data[355])
		fmt.Printf("Magenta: Hue=%3d Chroma=%3d Brightness=%3d (raw: %02x %02x %02x)\n",
			magentaHue, magentaChroma, magentaBrightness, data[353], data[354], data[355])
	}

	fmt.Println("\n==========================================")
	fmt.Println("Expected: All values should be NON-ZERO if there are color adjustments")
	fmt.Println("If all are 0, then the file has neutral colors")
}
