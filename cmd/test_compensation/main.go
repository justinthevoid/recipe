package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/formats/xmp"
)

func main() {
	// Read XMP file
	xmpPath := "testdata/visual-regression/images/JMH_0079.xmp"
	fmt.Printf("Reading XMP: %s\n", xmpPath)

	xmpData, err := os.ReadFile(xmpPath)
	if err != nil {
		fmt.Printf("Error reading XMP: %v\n", err)
		os.Exit(1)
	}

	// Parse XMP
	fmt.Println("Parsing XMP...")
	recipe, err := xmp.Parse(xmpData)
	if err != nil {
		fmt.Printf("Error parsing XMP: %v\n", err)
		os.Exit(1)
	}

	// Set preset name to match original NP3 structure
	recipe.Name = "Agfachrome RSX 2 v15"
	fmt.Printf("Setting preset name: %s\n", recipe.Name)

	// Enable baseline compensation
	fmt.Println("Enabling baseline compensation for Flexible Color profile...")
	if recipe.Metadata == nil {
		recipe.Metadata = make(map[string]interface{})
	}
	recipe.Metadata["baseline_compensation"] = "flexible_color"

	// Generate NP3 with compensation
	fmt.Println("Generating NP3 with baseline compensation...")
	np3Data, err := np3.Generate(recipe)
	if err != nil {
		fmt.Printf("Error generating NP3: %v\n", err)
		os.Exit(1)
	}

	// Save to output file
	outputPath := "testdata/visual-regression/test_output/Agfachrome RSX 200-compensated.np3"
	fmt.Printf("Saving compensated NP3: %s\n", outputPath)

	err = os.WriteFile(outputPath, np3Data, 0644)
	if err != nil {
		fmt.Printf("Error writing NP3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Success! Compensated NP3 file generated.")
	fmt.Printf("  Size: %d bytes\n", len(np3Data))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Open JMH_0079.NEF in NX Studio")
	fmt.Println("  2. Load the compensated preset: " + outputPath)
	fmt.Println("  3. Ensure 'Flexible Color' profile is active")
	fmt.Println("  4. Export as JPEG/TIFF")
	fmt.Println("  5. Run: python3 scripts/advanced_visual_compare.py")
}
