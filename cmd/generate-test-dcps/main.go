package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/justin/recipe/internal/formats/dcp"
	"github.com/justin/recipe/internal/models"
)

func main() {
	// Output directory
	outputDir := "testdata/dcp/generated"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==========================================")
	fmt.Println("Generating Test DCPs for Manual Validation")
	fmt.Println("==========================================")
	fmt.Printf("Output directory: %s\n\n", outputDir)

	// Define 3 test presets with unique profile names
	testPresets := []struct {
		filename    string
		profileName string
		recipe      *models.UniversalRecipe
	}{
		{
			filename:    "neutral.dcp",
			profileName: "Recipe Test - Neutral",
			recipe: &models.UniversalRecipe{
				Exposure:   0.0,
				Contrast:   0,
				Highlights: 0,
				Shadows:    0,
				Metadata: map[string]interface{}{
					"profile_name":   "Recipe Test - Neutral",
					"camera_model": "Nikon Z f", // Match user's camera
				},
			},
		},
		{
			filename:    "portrait.dcp",
			profileName: "Recipe Test - Portrait",
			recipe: &models.UniversalRecipe{
				Exposure:   0.5,  // +0.5 stops
				Contrast:   30,   // +30
				Highlights: -20,  // -20
				Shadows:    0,    // 0
				Metadata: map[string]interface{}{
					"profile_name":   "Recipe Test - Portrait",
					"camera_model": "Nikon Z f", // Match user's camera
				},
			},
		},
		{
			filename:    "landscape.dcp",
			profileName: "Recipe Test - Landscape",
			recipe: &models.UniversalRecipe{
				Exposure:   0.3, // +0.3 stops
				Contrast:   0,   // 0
				Highlights: 0,   // 0
				Shadows:    20,  // +20
				Metadata: map[string]interface{}{
					"profile_name":   "Recipe Test - Landscape",
					"camera_model": "Nikon Z f", // Match user's camera
				},
			},
		},
	}

	// Generate each DCP
	success := 0
	for i, preset := range testPresets {
		fmt.Printf("[%d/3] Generating %s\n", i+1, preset.filename)
		fmt.Printf("  Profile Name: %s\n", preset.profileName)
		fmt.Printf("  Parameters: Exposure=%.1f, Contrast=%d, Highlights=%d, Shadows=%d\n",
			preset.recipe.Exposure, preset.recipe.Contrast, preset.recipe.Highlights, preset.recipe.Shadows)

		// Generate DCP
		dcpData, err := dcp.Generate(preset.recipe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Error generating %s: %v\n", preset.filename, err)
			continue
		}

		// Write to file
		outputPath := filepath.Join(outputDir, preset.filename)
		if err := os.WriteFile(outputPath, dcpData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Error writing %s: %v\n", preset.filename, err)
			continue
		}

		fmt.Printf("  ✓ Generated successfully (%d bytes)\n\n", len(dcpData))
		success++
	}

	// Summary
	fmt.Println("==========================================")
	fmt.Println("Verification")
	fmt.Println("==========================================")
	if success == 3 {
		fmt.Println("✓ All 3 test DCPs generated successfully!")
		fmt.Println()
		fmt.Println("Generated files:")
		for _, preset := range testPresets {
			outputPath := filepath.Join(outputDir, preset.filename)
			if info, err := os.Stat(outputPath); err == nil {
				fmt.Printf("  %s (%d KB) - Profile: \"%s\"\n",
					preset.filename, info.Size()/1024, preset.profileName)
			}
		}
		fmt.Println()
		fmt.Println("Next steps for manual validation (Story 9-4):")
		fmt.Println("1. Install Adobe Camera Raw and Lightroom Classic 13.0+")
		fmt.Println("2. Copy DCPs to Adobe Camera Profiles folder:")
		fmt.Println("   Windows: C:\\Users\\<username>\\AppData\\Roaming\\Adobe\\CameraRaw\\CameraProfiles\\")
		fmt.Println("   macOS: ~/Library/Application Support/Adobe/CameraRaw/CameraProfiles/")
		fmt.Println("3. Restart Lightroom Classic")
		fmt.Println("4. Open RAW files and apply DCPs from Profile Browser")
		fmt.Println("   - Look for \"Recipe Test - Neutral\"")
		fmt.Println("   - Look for \"Recipe Test - Portrait\"")
		fmt.Println("   - Look for \"Recipe Test - Landscape\"")
		fmt.Println("5. Take screenshots and document results in validation-report.md")
		fmt.Println()
		fmt.Println("See testdata/dcp/MANUAL_TESTING_GUIDE.md for complete testing procedure")
	} else {
		fmt.Printf("✗ Only %d/3 DCPs were generated successfully\n", success)
		os.Exit(1)
	}
}
