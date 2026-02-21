//go:build ignore
// +build ignore

// test_full_pipeline.go - Test full NP3→UniversalRecipe→XMP pipeline
package main

import (
	"fmt"
	"os"
	"encoding/xml"

	np3 "github.com/justin/recipe/internal/formats/np3"
	xmpgen "github.com/justin/recipe/internal/formats/xmp"
)

func main() {
	// Step 1: Read NP3 file
	data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading NP3: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Parse NP3 to UniversalRecipe
	recipe, err := np3.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing NP3: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Display UniversalRecipe color values
	fmt.Println("UniversalRecipe Color Values:")
	fmt.Println("==============================")
	fmt.Printf("Red:     Hue=%3d Sat=%3d Lum=%3d\n", recipe.Red.Hue, recipe.Red.Saturation, recipe.Red.Luminance)
	fmt.Printf("Orange:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Orange.Hue, recipe.Orange.Saturation, recipe.Orange.Luminance)
	fmt.Printf("Yellow:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Yellow.Hue, recipe.Yellow.Saturation, recipe.Yellow.Luminance)
	fmt.Printf("Green:   Hue=%3d Sat=%3d Lum=%3d\n", recipe.Green.Hue, recipe.Green.Saturation, recipe.Green.Luminance)
	fmt.Printf("Aqua:    Hue=%3d Sat=%3d Lum=%3d\n", recipe.Aqua.Hue, recipe.Aqua.Saturation, recipe.Aqua.Luminance)
	fmt.Printf("Blue:    Hue=%3d Sat=%3d Lum=%3d\n", recipe.Blue.Hue, recipe.Blue.Saturation, recipe.Blue.Luminance)
	fmt.Printf("Purple:  Hue=%3d Sat=%3d Lum=%3d\n", recipe.Purple.Hue, recipe.Purple.Saturation, recipe.Purple.Luminance)
	fmt.Printf("Magenta: Hue=%3d Sat=%3d Lum=%3d\n", recipe.Magenta.Hue, recipe.Magenta.Saturation, recipe.Magenta.Luminance)
	fmt.Println()

	// Step 4: Generate XMP
	xmpData, err := xmpgen.Generate(recipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XMP: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Parse XMP to display actual values
	type XMPDescription struct {
		HueRed        string `xml:"HueRed,attr"`
		SaturationRed string `xml:"SaturationRed,attr"`
		LuminanceRed  string `xml:"LuminanceRed,attr"`
	}
	type XMP struct {
		XMLName xml.Name        `xml:"xmpmeta"`
		RDF     struct {
			Description XMPDescription `xml:"Description"`
		} `xml:"RDF"`
	}

	var xmp XMP
	if err := xml.Unmarshal(xmpData, &xmp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing XMP: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("XMP Output Values (Red color only):")
	fmt.Println("====================================")
	fmt.Printf("HueRed=%s SaturationRed=%s LuminanceRed=%s\n",
		xmp.RDF.Description.HueRed,
		xmp.RDF.Description.SaturationRed,
		xmp.RDF.Description.LuminanceRed)
	fmt.Println()

	fmt.Println("Expected from NP3 binary:")
	fmt.Println("=========================")
	fmt.Println("Red: Hue= 25 Chroma= 35 Brightness=  4")
	fmt.Println()
	fmt.Println("If UniversalRecipe shows 0,0,0 but NP3 binary has 25,35,4,")
	fmt.Println("then the problem is in NP3 parsing.")
	fmt.Println()
	fmt.Println("If UniversalRecipe shows 25,35,4 but XMP shows 0,0,0,")
	fmt.Println("then the problem is in XMP generation.")
}
