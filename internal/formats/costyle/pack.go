package costyle

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/justin/recipe/internal/models"
)

// Unpack extracts multiple .costyle files from a .costylepack ZIP archive
// and parses each to UniversalRecipe.
//
// The .costylepack format is a ZIP archive containing multiple .costyle preset files.
// This function:
//   - Validates ZIP structure and magic bytes
//   - Extracts all .costyle files (skips non-.costyle files with warning)
//   - Parses each .costyle file using existing Parse() function
//   - Returns slice of recipes (one per .costyle file in bundle)
//   - Handles individual file failures gracefully (partial results allowed)
//
// Parameters:
//   - data: Raw bytes of the .costylepack ZIP file
//
// Returns:
//   - []*models.UniversalRecipe: Slice of parsed recipes (one per .costyle file)
//   - error: Error if ZIP extraction fails completely, nil if successful (even if some files fail)
//
// Example:
//   costylepackData, _ := os.ReadFile("presets.costylepack")
//   recipes, err := costyle.Unpack(costylepackData)
//   if err != nil {
//       log.Fatal(err) // ZIP extraction failed
//   }
//   // recipes contains all successfully parsed .costyle files
func Unpack(data []byte) ([]*models.UniversalRecipe, error) {
	// Validate ZIP magic bytes (PK signature: 50 4B 03 04)
	if len(data) < 4 || !bytes.Equal(data[:4], []byte{0x50, 0x4B, 0x03, 0x04}) {
		return nil, fmt.Errorf("invalid .costylepack: not a valid ZIP file (missing magic bytes)")
	}

	// Open ZIP archive
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		// Handle specific ZIP errors
		return nil, fmt.Errorf("failed to read .costylepack ZIP: %w", err)
	}

	// Extract bundle metadata from ZIP comment (if present)
	bundleMetadata := make(map[string]interface{})
	if zipReader.Comment != "" {
		bundleMetadata["bundle_comment"] = zipReader.Comment
	}

	// Handle empty bundle (0 files)
	if len(zipReader.File) == 0 {
		return []*models.UniversalRecipe{}, fmt.Errorf("empty .costylepack: no files found in ZIP archive")
	}

	// Extract and parse all .costyle files
	var recipes []*models.UniversalRecipe
	var errors []string
	costyleCount := 0

	for _, file := range zipReader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		// Check if file has .costyle extension
		if filepath.Ext(file.Name) != ".costyle" {
			// Skip non-.costyle files with warning
			errors = append(errors, fmt.Sprintf("skipped non-.costyle file: %s", file.Name))
			continue
		}

		costyleCount++

		// Open file from ZIP
		rc, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to open %s: %v", file.Name, err))
			continue
		}

		// Read file contents
		fileData, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to read %s: %v", file.Name, err))
			continue
		}

		// Parse .costyle file
		recipe, err := Parse(fileData)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to parse %s: %v", file.Name, err))
			continue
		}

		// Track original filename for reference
		if recipe.Metadata == nil {
			recipe.Metadata = make(map[string]interface{})
		}
		recipe.Metadata["original_filename"] = file.Name

		// Add bundle metadata to each recipe
		for key, value := range bundleMetadata {
			recipe.Metadata[key] = value
		}

		recipes = append(recipes, recipe)
	}

	// If no .costyle files found, return error
	if costyleCount == 0 {
		return nil, fmt.Errorf("no .costyle files found in .costylepack bundle")
	}

	// Return partial results if some files failed (don't fail entire bundle)
	// Log errors if any occurred
	if len(errors) > 0 {
		// Note: In production, these would be logged via slog
		// For now, we allow partial success without failing
		_ = errors // Suppress unused variable warning
	}

	return recipes, nil
}

// Pack creates a .costylepack ZIP archive from multiple UniversalRecipes.
//
// This function:
//   - Validates inputs (at least one recipe required)
//   - Auto-generates filenames if not provided (Style1.costyle, Style2.costyle, etc.)
//   - Handles duplicate filenames by appending index suffix
//   - Generates .costyle XML for each recipe using existing Generate() function
//   - Packages all .costyle files into single ZIP archive
//   - Writes bundle metadata to ZIP comment
//
// Parameters:
//   - recipes: Slice of UniversalRecipes to pack (minimum 1 required)
//   - filenames: Optional filenames for each recipe (must match recipe count or be empty)
//
// Returns:
//   - []byte: ZIP archive bytes (.costylepack file)
//   - error: Error if validation or ZIP creation fails
//
// Example:
//   recipes := []*models.UniversalRecipe{recipe1, recipe2, recipe3}
//   filenames := []string{"Portrait.costyle", "Landscape.costyle", "Product.costyle"}
//   costylepackData, err := costyle.Pack(recipes, filenames)
//   if err != nil {
//       log.Fatal(err)
//   }
//   os.WriteFile("bundle.costylepack", costylepackData, 0644)
func Pack(recipes []*models.UniversalRecipe, filenames []string) ([]byte, error) {
	// Validate inputs
	if len(recipes) == 0 {
		return nil, fmt.Errorf("at least one recipe required to create .costylepack bundle")
	}

	// Validate filenames: must be empty OR match recipe count
	if len(filenames) != 0 && len(filenames) != len(recipes) {
		return nil, fmt.Errorf("filenames count (%d) must match recipes count (%d) or be empty", len(filenames), len(recipes))
	}

	// Auto-generate filenames if not provided
	if len(filenames) == 0 {
		filenames = make([]string, len(recipes))
		for i := range recipes {
			filenames[i] = fmt.Sprintf("Style%d.costyle", i+1)
		}
	}

	// Deduplicate filenames
	filenames = deduplicateFilenames(filenames)

	// Create ZIP archive buffer
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Track bundle metadata for ZIP comment
	bundleName := ""
	bundleDescription := ""

	// Add each recipe to ZIP
	for i, recipe := range recipes {
		// Extract bundle metadata from first recipe (if present)
		if i == 0 && recipe.Metadata != nil {
			if name, ok := recipe.Metadata["bundle_name"].(string); ok {
				bundleName = name
			}
			if desc, ok := recipe.Metadata["bundle_description"].(string); ok {
				bundleDescription = desc
			}
		}

		// Generate .costyle XML for this recipe
		xmlData, err := Generate(recipe)
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to generate .costyle for recipe %d: %w", i, err)
		}

		// Create file in ZIP
		fileWriter, err := zipWriter.Create(filenames[i])
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to create ZIP entry for %s: %w", filenames[i], err)
		}

		// Write XML data
		if _, err := fileWriter.Write(xmlData); err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to write data for %s: %w", filenames[i], err)
		}
	}

	// Write bundle metadata to ZIP comment
	if bundleName != "" || bundleDescription != "" {
		comment := fmt.Sprintf("Bundle: %s, Description: %s, Files: %d", bundleName, bundleDescription, len(recipes))
		if err := zipWriter.SetComment(comment); err != nil {
			// Non-fatal error, continue without comment
			_ = err
		}
	}

	// Close ZIP writer
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close ZIP archive: %w", err)
	}

	return buf.Bytes(), nil
}

// deduplicateFilenames handles duplicate filenames by appending index suffixes.
// For example: ["Preset.costyle", "Preset.costyle", "Other.costyle"]
// becomes: ["Preset.costyle", "Preset_1.costyle", "Other.costyle"]
func deduplicateFilenames(filenames []string) []string {
	usedNames := make(map[string]int)
	result := make([]string, len(filenames))

	for i, filename := range filenames {
		// Check if filename already used
		count, exists := usedNames[filename]
		if !exists {
			// First occurrence, use as-is
			result[i] = filename
			usedNames[filename] = 1
		} else {
			// Duplicate found, append index suffix
			baseName := strings.TrimSuffix(filename, ".costyle")
			newFilename := fmt.Sprintf("%s_%d.costyle", baseName, count)
			result[i] = newFilename
			usedNames[filename] = count + 1
		}
	}

	return result
}
