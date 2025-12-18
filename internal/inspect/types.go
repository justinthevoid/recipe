// Package inspect provides JSON inspection and serialization for photo preset parameters.
package inspect

import "github.com/justin/recipe/internal/models"

// InspectOutput represents the JSON output structure for the inspect command.
// It wraps the UniversalRecipe with metadata about the inspection operation.
type InspectOutput struct {
	Metadata   Metadata               `json:"metadata"`
	Parameters *models.UniversalRecipe `json:"parameters"`
}

// Metadata contains information about the inspection operation.
type Metadata struct {
	SourceFile    string `json:"source_file"`     // Original filename
	SourceFormat  string `json:"source_format"`   // Detected format: "np3", "xmp", "lrtemplate"
	ParsedAt      string `json:"parsed_at"`       // ISO 8601 timestamp (UTC)
	RecipeVersion string `json:"recipe_version"`  // Recipe tool version
}
