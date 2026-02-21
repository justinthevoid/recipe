package models

import (
	"encoding/json"
	"testing"
)

// TestMetadataFieldExists verifies the Metadata field is accessible and has correct type
func TestMetadataFieldExists(t *testing.T) {
	recipe := UniversalRecipe{}

	// Verify field can be assigned
	recipe.Metadata = make(map[string]interface{})
	recipe.Metadata["test_key"] = "test_value"

	// Verify field can be read
	if val, ok := recipe.Metadata["test_key"]; !ok || val != "test_value" {
		t.Errorf("Metadata field assignment/retrieval failed")
	}
}

// TestMetadataJSONSerialization tests JSON serialization with various value types
func TestMetadataJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
	}{
		{
			name: "string values",
			metadata: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "numeric values",
			metadata: map[string]interface{}{
				"int_val":    42,
				"float_val":  3.14,
				"zero_val":   0,
				"negative":   -100,
			},
		},
		{
			name: "boolean values",
			metadata: map[string]interface{}{
				"bool_true":  true,
				"bool_false": false,
			},
		},
		{
			name: "array values",
			metadata: map[string]interface{}{
				"string_array": []interface{}{"a", "b", "c"},
				"number_array": []interface{}{1, 2, 3},
				"mixed_array":  []interface{}{"text", 42, true},
			},
		},
		{
			name: "nested map values",
			metadata: map[string]interface{}{
				"nested": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep_value",
					},
				},
			},
		},
		{
			name: "mixed complex structure",
			metadata: map[string]interface{}{
				"string":  "text",
				"number":  123,
				"boolean": true,
				"array":   []interface{}{1, 2, 3},
				"nested": map[string]interface{}{
					"key": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe := UniversalRecipe{
				Name:     "test_recipe",
				Metadata: tt.metadata,
			}

			// Serialize to JSON
			jsonData, err := json.Marshal(recipe)
			if err != nil {
				t.Fatalf("JSON serialization failed: %v", err)
			}

			// Verify metadata is in JSON
			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Fatalf("JSON unmarshaling failed: %v", err)
			}

			// Verify metadata field exists in JSON
			if _, ok := result["metadata"]; !ok {
				t.Error("metadata field missing from JSON output")
			}
		})
	}
}

// TestMetadataJSONDeserialization tests JSON deserialization into Metadata
func TestMetadataJSONDeserialization(t *testing.T) {
	tests := []struct {
		name         string
		jsonInput    string
		validateFunc func(t *testing.T, metadata map[string]interface{})
	}{
		{
			name:      "string values",
			jsonInput: `{"name":"test","metadata":{"key1":"value1","key2":"value2"}}`,
			validateFunc: func(t *testing.T, metadata map[string]interface{}) {
				if metadata["key1"] != "value1" {
					t.Error("String value not deserialized correctly")
				}
			},
		},
		{
			name:      "numeric values",
			jsonInput: `{"name":"test","metadata":{"int_val":42,"float_val":3.14}}`,
			validateFunc: func(t *testing.T, metadata map[string]interface{}) {
				// JSON numbers deserialize as float64
				if metadata["int_val"] != float64(42) {
					t.Errorf("Integer value not deserialized correctly: got %v", metadata["int_val"])
				}
				if metadata["float_val"] != 3.14 {
					t.Error("Float value not deserialized correctly")
				}
			},
		},
		{
			name:      "boolean values",
			jsonInput: `{"name":"test","metadata":{"bool_true":true,"bool_false":false}}`,
			validateFunc: func(t *testing.T, metadata map[string]interface{}) {
				if metadata["bool_true"] != true {
					t.Error("Boolean true not deserialized correctly")
				}
				if metadata["bool_false"] != false {
					t.Error("Boolean false not deserialized correctly")
				}
			},
		},
		{
			name:      "array values",
			jsonInput: `{"name":"test","metadata":{"array":[1,2,3]}}`,
			validateFunc: func(t *testing.T, metadata map[string]interface{}) {
				arr, ok := metadata["array"].([]interface{})
				if !ok {
					t.Fatal("Array not deserialized correctly")
				}
				if len(arr) != 3 {
					t.Errorf("Array length incorrect: got %d, want 3", len(arr))
				}
			},
		},
		{
			name:      "nested structures",
			jsonInput: `{"name":"test","metadata":{"nested":{"level2":{"level3":"deep"}}}}`,
			validateFunc: func(t *testing.T, metadata map[string]interface{}) {
				nested, ok := metadata["nested"].(map[string]interface{})
				if !ok {
					t.Fatal("Nested map not deserialized correctly")
				}
				level2, ok := nested["level2"].(map[string]interface{})
				if !ok {
					t.Fatal("Level 2 nested map not deserialized correctly")
				}
				if level2["level3"] != "deep" {
					t.Error("Level 3 value not deserialized correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var recipe UniversalRecipe
			if err := json.Unmarshal([]byte(tt.jsonInput), &recipe); err != nil {
				t.Fatalf("JSON deserialization failed: %v", err)
			}

			if recipe.Metadata == nil {
				t.Fatal("Metadata is nil after deserialization")
			}

			tt.validateFunc(t, recipe.Metadata)
		})
	}
}

// TestMetadataOmitempty tests that empty/nil maps are omitted from JSON output
func TestMetadataOmitempty(t *testing.T) {
	tests := []struct {
		name          string
		recipe        UniversalRecipe
		expectInJSON  bool
	}{
		{
			name: "nil metadata",
			recipe: UniversalRecipe{
				Name:     "test",
				Metadata: nil,
			},
			expectInJSON: false,
		},
		{
			name: "empty metadata",
			recipe: UniversalRecipe{
				Name:     "test",
				Metadata: make(map[string]interface{}),
			},
			expectInJSON: false,
		},
		{
			name: "populated metadata",
			recipe: UniversalRecipe{
				Name: "test",
				Metadata: map[string]interface{}{
					"key": "value",
				},
			},
			expectInJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.recipe)
			if err != nil {
				t.Fatalf("JSON serialization failed: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Fatalf("JSON unmarshaling failed: %v", err)
			}

			_, hasMetadata := result["metadata"]
			if hasMetadata != tt.expectInJSON {
				t.Errorf("metadata in JSON: got %v, want %v", hasMetadata, tt.expectInJSON)
			}
		})
	}
}

// TestMetadataNestedStructures tests complex nested maps and arrays
func TestMetadataNestedStructures(t *testing.T) {
	recipe := UniversalRecipe{
		Name: "test",
		Metadata: map[string]interface{}{
			"complex_structure": map[string]interface{}{
				"arrays": []interface{}{
					map[string]interface{}{
						"nested_array": []interface{}{1, 2, 3},
					},
					"string_value",
					42,
				},
				"nested_maps": map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{
							"a", "b", "c",
						},
					},
				},
			},
		},
	}

	// Serialize to JSON and back
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Verify structure is preserved
	if result.Metadata == nil {
		t.Fatal("Metadata is nil")
	}

	complex, ok := result.Metadata["complex_structure"].(map[string]interface{})
	if !ok {
		t.Fatal("Complex structure not preserved")
	}

	if _, ok := complex["arrays"]; !ok {
		t.Error("Arrays field not preserved")
	}

	if _, ok := complex["nested_maps"]; !ok {
		t.Error("Nested maps field not preserved")
	}
}

// TestMetadataStory1_8Example1 validates tone curve storage example from story 1-8
func TestMetadataStory1_8Example1(t *testing.T) {
	// Create recipe with point curve
	recipe := UniversalRecipe{
		Name: "tone_curve_test",
		PointCurve: []ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 128, Output: 140},
			{Input: 255, Output: 255},
		},
	}

	// Store tone curve when generating NP3 (simulated)
	recipe.Metadata = make(map[string]interface{})
	if len(recipe.PointCurve) > 0 {
		curveJSON, err := json.Marshal(recipe.PointCurve)
		if err != nil {
			t.Fatalf("Failed to marshal tone curve: %v", err)
		}
		recipe.Metadata["xmp_tone_curve_pv2012"] = string(curveJSON)
	}

	// Serialize and deserialize
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Retrieve tone curve when parsing back to XMP (simulated)
	if curveData, ok := result.Metadata["xmp_tone_curve_pv2012"]; ok {
		var curve []ToneCurvePoint
		if err := json.Unmarshal([]byte(curveData.(string)), &curve); err != nil {
			t.Fatalf("Failed to unmarshal stored tone curve: %v", err)
		}

		// Verify curve was preserved
		if len(curve) != 3 {
			t.Errorf("Curve length: got %d, want 3", len(curve))
		}
		if curve[1].Input != 128 || curve[1].Output != 140 {
			t.Errorf("Curve midpoint incorrect: got (%d,%d), want (128,140)",
				curve[1].Input, curve[1].Output)
		}
	} else {
		t.Error("Stored tone curve not found in metadata")
	}
}

// TestMetadataStory1_8Example2 validates HSL adjustments storage example from story 1-8
func TestMetadataStory1_8Example2(t *testing.T) {
	recipe := UniversalRecipe{
		Name: "hsl_test",
		Red: ColorAdjustment{
			Hue:        10,
			Saturation: 20,
			Luminance:  5,
		},
		Orange: ColorAdjustment{
			Hue:        -5,
			Saturation: 15,
			Luminance:  0,
		},
	}

	// Store HSL adjustments when generating NP3 (simulated)
	recipe.Metadata = make(map[string]interface{})
	hslData := map[string]interface{}{
		"red_hue":         recipe.Red.Hue,
		"red_saturation":  recipe.Red.Saturation,
		"red_luminance":   recipe.Red.Luminance,
		"orange_hue":      recipe.Orange.Hue,
		"orange_saturation": recipe.Orange.Saturation,
		"orange_luminance":  recipe.Orange.Luminance,
	}
	hslJSON, err := json.Marshal(hslData)
	if err != nil {
		t.Fatalf("Failed to marshal HSL data: %v", err)
	}
	recipe.Metadata["xmp_hsl_adjustments"] = string(hslJSON)

	// Serialize and deserialize
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Verify HSL data can be retrieved
	if hslStr, ok := result.Metadata["xmp_hsl_adjustments"]; ok {
		var hslResult map[string]interface{}
		if err := json.Unmarshal([]byte(hslStr.(string)), &hslResult); err != nil {
			t.Fatalf("Failed to unmarshal HSL data: %v", err)
		}

		// Verify values (JSON numbers are float64)
		if hslResult["red_hue"] != float64(10) {
			t.Errorf("Red hue: got %v, want 10", hslResult["red_hue"])
		}
		if hslResult["orange_saturation"] != float64(15) {
			t.Errorf("Orange saturation: got %v, want 15", hslResult["orange_saturation"])
		}
	} else {
		t.Error("Stored HSL data not found in metadata")
	}
}

// TestMetadataStory1_8Example3 validates split toning storage example from story 1-8
func TestMetadataStory1_8Example3(t *testing.T) {
	recipe := UniversalRecipe{
		Name:                    "split_toning_test",
		SplitShadowHue:          220,
		SplitShadowSaturation:   35,
		SplitHighlightHue:       45,
		SplitHighlightSaturation: 25,
		SplitBalance:            10,
	}

	// Store split toning when generating NP3 (simulated)
	recipe.Metadata = make(map[string]interface{})
	if recipe.SplitShadowHue != 0 || recipe.SplitHighlightHue != 0 {
		splitData := map[string]interface{}{
			"shadow_hue":        recipe.SplitShadowHue,
			"shadow_saturation": recipe.SplitShadowSaturation,
			"highlight_hue":     recipe.SplitHighlightHue,
			"highlight_saturation": recipe.SplitHighlightSaturation,
			"balance":           recipe.SplitBalance,
		}
		recipe.Metadata["lrtemplate_split_toning"] = splitData
	}

	// Serialize and deserialize
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Verify split toning data can be retrieved
	if splitData, ok := result.Metadata["lrtemplate_split_toning"]; ok {
		splitMap := splitData.(map[string]interface{})

		// Verify values (JSON numbers are float64)
		if splitMap["shadow_hue"] != float64(220) {
			t.Errorf("Shadow hue: got %v, want 220", splitMap["shadow_hue"])
		}
		if splitMap["highlight_saturation"] != float64(25) {
			t.Errorf("Highlight saturation: got %v, want 25", splitMap["highlight_saturation"])
		}
	} else {
		t.Error("Stored split toning data not found in metadata")
	}
}

// TestMetadataStory1_8Example4 validates grain effects storage example from story 1-8
func TestMetadataStory1_8Example4(t *testing.T) {
	recipe := UniversalRecipe{
		Name:          "grain_test",
		GrainAmount:   50,
		GrainSize:     35,
		GrainRoughness: 60,
	}

	// Store grain effects when generating NP3 (simulated)
	recipe.Metadata = make(map[string]interface{})
	if recipe.GrainAmount != 0 {
		grainData := map[string]interface{}{
			"amount":    recipe.GrainAmount,
			"size":      recipe.GrainSize,
			"roughness": recipe.GrainRoughness,
		}
		recipe.Metadata["effects_grain"] = grainData
	}

	// Serialize and deserialize
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Verify grain data can be retrieved
	if grainData, ok := result.Metadata["effects_grain"]; ok {
		grainMap := grainData.(map[string]interface{})

		// Verify values (JSON numbers are float64)
		if grainMap["amount"] != float64(50) {
			t.Errorf("Grain amount: got %v, want 50", grainMap["amount"])
		}
		if grainMap["size"] != float64(35) {
			t.Errorf("Grain size: got %v, want 35", grainMap["size"])
		}
		if grainMap["roughness"] != float64(60) {
			t.Errorf("Grain roughness: got %v, want 60", grainMap["roughness"])
		}
	} else {
		t.Error("Stored grain data not found in metadata")
	}
}

// TestMetadataKeyNamingConvention tests format_fieldname pattern
func TestMetadataKeyNamingConvention(t *testing.T) {
	recipe := UniversalRecipe{
		Name:     "key_naming_test",
		Metadata: make(map[string]interface{}),
	}

	// Test various key naming patterns from documentation
	testKeys := []string{
		"xmp_tone_curve_pv2012",
		"xmp_hsl_adjustments",
		"lrtemplate_split_toning",
		"np3_custom_data",
		"effects_grain",
		"format_specific_parameter",
	}

	for _, key := range testKeys {
		recipe.Metadata[key] = "test_value"
	}

	// Serialize and deserialize
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var result UniversalRecipe
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	// Verify all keys are preserved
	for _, key := range testKeys {
		if _, ok := result.Metadata[key]; !ok {
			t.Errorf("Key %q not preserved in metadata", key)
		}
	}
}
