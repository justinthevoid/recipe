package inspect

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/justin/recipe/internal/models"
)

// TestToJSON verifies JSON serialization of UniversalRecipe
func TestToJSON(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Name:       "Test Preset",
		Exposure:   0.5,
		Contrast:   15,
		Highlights: -20,
		Shadows:    10,
		Saturation: -5,
		Vibrance:   10,
		Clarity:    5,
		Sharpness:  25,
	}

	jsonBytes, err := ToJSON(recipe)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify pretty-printing (should have newlines and indentation)
	jsonStr := string(jsonBytes)
	if !strings.Contains(jsonStr, "\n") {
		t.Error("JSON is not pretty-printed (no newlines)")
	}
	if !strings.Contains(jsonStr, "  ") {
		t.Error("JSON is not indented (no 2-space indent)")
	}

	// Verify key fields present
	if _, ok := parsed["name"]; !ok {
		t.Error("Missing 'name' field in JSON output")
	}
	if _, ok := parsed["exposure"]; !ok {
		t.Error("Missing 'exposure' field in JSON output")
	}
	if _, ok := parsed["contrast"]; !ok {
		t.Error("Missing 'contrast' field in JSON output")
	}
}

// TestToJSON_NilRecipe verifies error handling for nil input
func TestToJSON_NilRecipe(t *testing.T) {
	_, err := ToJSON(nil)
	if err == nil {
		t.Error("Expected error for nil recipe, got nil")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("Error message should mention nil, got: %v", err)
	}
}

// TestToJSONWithMetadata verifies metadata wrapper
func TestToJSONWithMetadata(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Name:     "Test Preset",
		Contrast: 15,
	}

	sourceFile := "/path/to/portrait.np3"
	format := "np3"
	version := "2.0.0"

	jsonBytes, err := ToJSONWithMetadata(recipe, sourceFile, format, version)
	if err != nil {
		t.Fatalf("ToJSONWithMetadata failed: %v", err)
	}

	// Parse output
	var output InspectOutput
	if err := json.Unmarshal(jsonBytes, &output); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	// Verify metadata fields
	if output.Metadata.SourceFile != "portrait.np3" {
		t.Errorf("Expected source_file='portrait.np3', got '%s'", output.Metadata.SourceFile)
	}
	if output.Metadata.SourceFormat != "np3" {
		t.Errorf("Expected source_format='np3', got '%s'", output.Metadata.SourceFormat)
	}
	if output.Metadata.RecipeVersion != "2.0.0" {
		t.Errorf("Expected recipe_version='2.0.0', got '%s'", output.Metadata.RecipeVersion)
	}

	// Verify timestamp is valid ISO 8601
	if output.Metadata.ParsedAt == "" {
		t.Error("parsed_at timestamp is empty")
	}
	parsedTime, err := time.Parse(time.RFC3339, output.Metadata.ParsedAt)
	if err != nil {
		t.Errorf("parsed_at is not valid ISO 8601 format: %v", err)
	}

	// Verify timestamp is recent (within last minute)
	if time.Since(parsedTime) > time.Minute {
		t.Errorf("Timestamp is not recent: %s", output.Metadata.ParsedAt)
	}

	// Verify parameters section present
	if output.Parameters == nil {
		t.Error("Parameters section is nil")
	}
	if output.Parameters.Contrast != 15 {
		t.Errorf("Expected contrast=15, got %d", output.Parameters.Contrast)
	}
}

// TestToJSONWithMetadata_NilRecipe verifies error handling
func TestToJSONWithMetadata_NilRecipe(t *testing.T) {
	_, err := ToJSONWithMetadata(nil, "test.np3", "np3", "1.0")
	if err == nil {
		t.Error("Expected error for nil recipe, got nil")
	}
}

// TestJSONStructure verifies JSON follows conventions
func TestJSONStructure(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Name:       "Test",
		Exposure:   0.5,
		Contrast:   10,
		Saturation: -5,
	}

	jsonBytes, err := ToJSONWithMetadata(recipe, "test.np3", "np3", "1.0")
	if err != nil {
		t.Fatalf("ToJSONWithMetadata failed: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify camelCase keys
	if !strings.Contains(jsonStr, `"source_file"`) {
		t.Error("Expected snake_case 'source_file' in metadata")
	}
	if !strings.Contains(jsonStr, `"source_format"`) {
		t.Error("Expected snake_case 'source_format' in metadata")
	}
	if !strings.Contains(jsonStr, `"parsed_at"`) {
		t.Error("Expected snake_case 'parsed_at' in metadata")
	}
	if !strings.Contains(jsonStr, `"recipe_version"`) {
		t.Error("Expected snake_case 'recipe_version' in metadata")
	}

	// Verify top-level structure has metadata and parameters
	if !strings.Contains(jsonStr, `"metadata"`) {
		t.Error("Missing top-level 'metadata' key")
	}
	if !strings.Contains(jsonStr, `"parameters"`) {
		t.Error("Missing top-level 'parameters' key")
	}

	// Verify 2-space indentation
	if !strings.Contains(jsonStr, "  \"metadata\":") || !strings.Contains(jsonStr, "  \"parameters\":") {
		t.Error("JSON is not indented with 2 spaces at root level")
	}
}

// TestAllUniversalRecipeFields verifies all fields are serialized
func TestAllUniversalRecipeFields(t *testing.T) {
	recipe := &models.UniversalRecipe{
		Name:          "Complete Test",
		SourceFormat:  "np3",
		Exposure:      1.0,
		Contrast:      20,
		Highlights:    -30,
		Shadows:       15,
		Whites:        10,
		Blacks:        -10,
		Texture:       5,
		Clarity:       10,
		Dehaze:        15,
		Vibrance:      20,
		Saturation:    -10,
		Sharpness:     40,
		SharpnessRadius: 1.5,
		SharpnessDetail: 50,
		SharpnessMasking: 30,
		Tint:          -5,
		Red:           models.ColorAdjustment{Hue: 5, Saturation: 10, Luminance: 3},
		Orange:        models.ColorAdjustment{Hue: 10, Saturation: 15, Luminance: 5},
		Yellow:        models.ColorAdjustment{Hue: 0, Saturation: 5, Luminance: 2},
		Green:         models.ColorAdjustment{Hue: -5, Saturation: 8, Luminance: 1},
		Aqua:          models.ColorAdjustment{Hue: 3, Saturation: 12, Luminance: 4},
		Blue:          models.ColorAdjustment{Hue: -3, Saturation: 6, Luminance: 2},
		Purple:        models.ColorAdjustment{Hue: 7, Saturation: 14, Luminance: 3},
		Magenta:       models.ColorAdjustment{Hue: 2, Saturation: 9, Luminance: 1},
	}

	jsonBytes, err := ToJSON(recipe)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify key fields exist
	expectedFields := []string{
		"name", "sourceFormat", "exposure", "contrast", "highlights", "shadows",
		"whites", "blacks", "texture", "clarity", "dehaze", "vibrance", "saturation",
		"sharpness", "sharpnessRadius", "sharpnessDetail", "sharpnessMasking",
		"tint", "red", "orange", "yellow", "green", "aqua", "blue", "purple", "magenta",
	}

	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("Missing expected field '%s' in JSON output", field)
		}
	}

	// Verify nested color structures
	if red, ok := parsed["red"].(map[string]interface{}); ok {
		if _, hasHue := red["hue"]; !hasHue {
			t.Error("Red color adjustment missing 'hue' field")
		}
		if _, hasSat := red["saturation"]; !hasSat {
			t.Error("Red color adjustment missing 'saturation' field")
		}
		if _, hasLum := red["luminance"]; !hasLum {
			t.Error("Red color adjustment missing 'luminance' field")
		}
	} else {
		t.Error("Red color adjustment not serialized as object")
	}
}
