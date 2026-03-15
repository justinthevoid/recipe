// Package inspect provides JSON inspection and serialization for photo preset parameters.
package inspect

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/justin/recipe/internal/models"
)

// ToJSON serializes a UniversalRecipe to indented JSON.
// Returns pretty-printed JSON with 2-space indentation (AC-1).
func ToJSON(recipe *models.UniversalRecipe) ([]byte, error) {
	if recipe == nil {
		return nil, fmt.Errorf("recipe cannot be nil")
	}

	jsonBytes, err := json.MarshalIndent(recipe, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonBytes, nil
}

// ToJSONWithMetadata serializes a UniversalRecipe with metadata wrapper to indented JSON.
// This is the primary function used by the inspect CLI command (AC-2).
//
// Parameters:
//   - recipe: The parsed UniversalRecipe to serialize
//   - sourceFile: The original file path (will be converted to basename)
//   - format: The detected format ("np3" or "xmp")
//   - version: The Recipe tool version string
//
// Returns pretty-printed JSON with 2-space indentation.
func ToJSONWithMetadata(recipe *models.UniversalRecipe, sourceFile, format, version string) ([]byte, error) {
	if recipe == nil {
		return nil, fmt.Errorf("recipe cannot be nil")
	}

	// Create output structure with metadata
	output := InspectOutput{
		Metadata: Metadata{
			SourceFile:    filepath.Base(sourceFile), // Use basename only (not full path)
			SourceFormat:  format,
			ParsedAt:      time.Now().UTC().Format(time.RFC3339), // ISO 8601 in UTC
			RecipeVersion: version,
		},
		Parameters: recipe,
	}

	// Serialize to pretty-printed JSON
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonBytes, nil
}
