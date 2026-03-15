package models

import (
	"encoding/json"
	"testing"
)

// BenchmarkUniversalRecipeJSONMarshal benchmarks JSON marshaling of UniversalRecipe
func BenchmarkUniversalRecipeJSONMarshal(b *testing.B) {
	recipe := UniversalRecipe{
		Name:       "Benchmark Test",
		Exposure:   1.5,
		Contrast:   50,
		Highlights: -25,
		Shadows:    30,
		Red: ColorAdjustment{
			Hue:        10,
			Saturation: 20,
			Luminance:  5,
		},
		PointCurve: []ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 128, Output: 140},
			{Input: 255, Output: 255},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUniversalRecipeJSONMarshalWithMetadata benchmarks JSON marshaling with Metadata field populated
func BenchmarkUniversalRecipeJSONMarshalWithMetadata(b *testing.B) {
	recipe := UniversalRecipe{
		Name:       "Benchmark Test",
		Exposure:   1.5,
		Contrast:   50,
		Highlights: -25,
		Shadows:    30,
		Red: ColorAdjustment{
			Hue:        10,
			Saturation: 20,
			Luminance:  5,
		},
		PointCurve: []ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 128, Output: 140},
			{Input: 255, Output: 255},
		},
		Metadata: map[string]interface{}{
			"xmp_tone_curve_pv2012": `[{"input":0,"output":0},{"input":128,"output":140}]`,
			"xmp_hsl_red_hue":       10,
			"format_split_toning": map[string]interface{}{
				"shadow_hue": 220,
				"balance":    10,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(recipe)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUniversalRecipeJSONUnmarshal benchmarks JSON unmarshaling
func BenchmarkUniversalRecipeJSONUnmarshal(b *testing.B) {
	recipe := UniversalRecipe{
		Name:       "Benchmark Test",
		Exposure:   1.5,
		Contrast:   50,
		Highlights: -25,
		Shadows:    30,
	}

	jsonData, _ := json.Marshal(recipe)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var r UniversalRecipe
		err := json.Unmarshal(jsonData, &r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMetadataAccess benchmarks direct access to Metadata field
func BenchmarkMetadataAccess(b *testing.B) {
	recipe := UniversalRecipe{
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = recipe.Metadata["key1"]
		_ = recipe.Metadata["key2"]
		_ = recipe.Metadata["key3"]
	}
}

// BenchmarkMetadataInsert benchmarks inserting values into Metadata
func BenchmarkMetadataInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recipe := UniversalRecipe{
			Metadata: make(map[string]interface{}),
		}
		recipe.Metadata["xmp_tone_curve"] = "test_data"
		recipe.Metadata["xmp_hsl_red"] = 10
		recipe.Metadata["format_split"] = map[string]interface{}{
			"hue": 220,
		}
	}
}
