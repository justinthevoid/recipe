package inspect

import (
	"testing"

	"github.com/justin/recipe/internal/models"
)

func TestDiff_IdenticalRecipes(t *testing.T) {
	recipe1 := &models.UniversalRecipe{
		Name:       "Test Recipe",
		Exposure:   0.5,
		Contrast:   15,
		Saturation: -10,
	}

	recipe2 := &models.UniversalRecipe{
		Name:       "Test Recipe",
		Exposure:   0.5,
		Contrast:   15,
		Saturation: -10,
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected no differences, got %d", len(results))
	}
}

func TestDiff_ModifiedFields(t *testing.T) {
	recipe1 := &models.UniversalRecipe{
		Contrast:   0,
		Saturation: 0,
		Exposure:   0.5,
	}

	recipe2 := &models.UniversalRecipe{
		Contrast:   15,
		Saturation: -10,
		Exposure:   0.51,
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Should detect 3 modified fields
	modifiedCount := 0
	for _, r := range results {
		if r.ChangeType == "modified" {
			modifiedCount++
		}
	}

	if modifiedCount != 3 {
		t.Errorf("Expected 3 modified fields, got %d", modifiedCount)
	}

	// Check Contrast specifically
	var contrastResult *DiffResult
	for _, r := range results {
		if r.Field == "Contrast" {
			contrastResult = &r
			break
		}
	}

	if contrastResult == nil {
		t.Error("Contrast change not detected")
	} else {
		if contrastResult.OldValue != int64(0) {
			t.Errorf("Expected old Contrast value 0, got %v", contrastResult.OldValue)
		}
		if contrastResult.NewValue != int64(15) {
			t.Errorf("Expected new Contrast value 15, got %v", contrastResult.NewValue)
		}
	}
}

func TestDiff_AddedField(t *testing.T) {
	temp := 5500
	recipe1 := &models.UniversalRecipe{
		Temperature: nil,
	}

	recipe2 := &models.UniversalRecipe{
		Temperature: &temp,
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find Temperature result
	var tempResult *DiffResult
	for _, r := range results {
		if r.Field == "Temperature" {
			tempResult = &r
			break
		}
	}

	if tempResult == nil {
		t.Fatal("Temperature change not detected")
	}

	if tempResult.ChangeType != "added" {
		t.Errorf("Expected change type 'added', got '%s'", tempResult.ChangeType)
	}

	if tempResult.OldValue != nil {
		t.Errorf("Expected nil old value, got %v", tempResult.OldValue)
	}

	if tempResult.NewValue != 5500 {
		t.Errorf("Expected new value 5500, got %v", tempResult.NewValue)
	}
}

func TestDiff_RemovedField(t *testing.T) {
	temp := 5500
	recipe1 := &models.UniversalRecipe{
		Temperature: &temp,
	}

	recipe2 := &models.UniversalRecipe{
		Temperature: nil,
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find Temperature result
	var tempResult *DiffResult
	for _, r := range results {
		if r.Field == "Temperature" {
			tempResult = &r
			break
		}
	}

	if tempResult == nil {
		t.Fatal("Temperature change not detected")
	}

	if tempResult.ChangeType != "removed" {
		t.Errorf("Expected change type 'removed', got '%s'", tempResult.ChangeType)
	}

	if tempResult.OldValue != 5500 {
		t.Errorf("Expected old value 5500, got %v", tempResult.OldValue)
	}

	if tempResult.NewValue != nil {
		t.Errorf("Expected nil new value, got %v", tempResult.NewValue)
	}
}

func TestDiff_SignificantChanges(t *testing.T) {
	tests := []struct {
		name        string
		old         int
		new         int
		significant bool
	}{
		{"Zero to non-zero", 0, 15, true},              // Always significant
		{"Large change", 10, 20, true},                 // 100% change
		{"6% change", 100, 106, true},                  // >5%
		{"5% change", 100, 105, false},                 // Exactly 5%, not >5%
		{"Small change", 100, 102, false},              // 2% change
		{"Negative to positive", -50, -45, true},       // 10% change (5/50)
		{"Large percentage", 0, 1, true},               // Infinite % (zero base)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe1 := &models.UniversalRecipe{Contrast: tt.old}
			recipe2 := &models.UniversalRecipe{Contrast: tt.new}

			results, err := Diff(recipe1, recipe2, 0.001)
			if err != nil {
				t.Fatalf("Diff failed: %v", err)
			}

			// Find Contrast result
			var contrastResult *DiffResult
			for _, r := range results {
				if r.Field == "Contrast" {
					contrastResult = &r
					break
				}
			}

			if contrastResult == nil {
				t.Fatal("Contrast change not detected")
			}

			if contrastResult.Significant != tt.significant {
				t.Errorf("Expected significant=%v, got %v (old=%d, new=%d, change=%.1f%%)",
					tt.significant, contrastResult.Significant, tt.old, tt.new,
					float64(tt.new-tt.old)/float64(tt.old)*100)
			}
		})
	}
}

func TestDiff_ToleranceForFloats(t *testing.T) {
	tests := []struct {
		name      string
		old       float64
		new       float64
		tolerance float64
		hasDiff   bool
	}{
		{"Within tolerance", 0.500, 0.500, 0.001, false},
		{"Just within tolerance", 0.500, 0.5009, 0.001, false},
		{"Outside tolerance", 0.500, 0.502, 0.001, true},
		{"Large tolerance", 0.500, 0.510, 0.02, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe1 := &models.UniversalRecipe{Exposure: tt.old}
			recipe2 := &models.UniversalRecipe{Exposure: tt.new}

			results, err := Diff(recipe1, recipe2, tt.tolerance)
			if err != nil {
				t.Fatalf("Diff failed: %v", err)
			}

			// Check if Exposure changed
			exposureChanged := false
			for _, r := range results {
				if r.Field == "Exposure" {
					exposureChanged = true
					break
				}
			}

			if exposureChanged != tt.hasDiff {
				t.Errorf("Expected hasDiff=%v, got %v (old=%.4f, new=%.4f, tolerance=%.4f)",
					tt.hasDiff, exposureChanged, tt.old, tt.new, tt.tolerance)
			}
		})
	}
}

func TestDiff_StructComparison(t *testing.T) {
	recipe1 := &models.UniversalRecipe{
		Red: models.ColorAdjustment{Hue: 0, Saturation: 0, Luminance: 0},
	}

	recipe2 := &models.UniversalRecipe{
		Red: models.ColorAdjustment{Hue: 10, Saturation: 0, Luminance: 0},
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find Red result
	var redResult *DiffResult
	for _, r := range results {
		if r.Field == "Red" {
			redResult = &r
			break
		}
	}

	if redResult == nil {
		t.Fatal("Red ColorAdjustment change not detected")
	}

	if redResult.ChangeType != "modified" {
		t.Errorf("Expected change type 'modified', got '%s'", redResult.ChangeType)
	}

	if !redResult.Significant {
		t.Error("Expected struct change to be significant")
	}
}

func TestDiff_SliceComparison(t *testing.T) {
	recipe1 := &models.UniversalRecipe{
		PointCurve: []models.ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 255, Output: 255},
		},
	}

	recipe2 := &models.UniversalRecipe{
		PointCurve: []models.ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 128, Output: 140}, // Different point
			{Input: 255, Output: 255},
		},
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find PointCurve result
	var curveResult *DiffResult
	for _, r := range results {
		if r.Field == "PointCurve" {
			curveResult = &r
			break
		}
	}

	if curveResult == nil {
		t.Fatal("PointCurve change not detected")
	}

	if curveResult.ChangeType != "modified" {
		t.Errorf("Expected change type 'modified', got '%s'", curveResult.ChangeType)
	}
}

func TestDiff_MapComparison(t *testing.T) {
	recipe1 := &models.UniversalRecipe{
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	recipe2 := &models.UniversalRecipe{
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 43, // Changed
		},
	}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find Metadata result
	var metadataResult *DiffResult
	for _, r := range results {
		if r.Field == "Metadata" {
			metadataResult = &r
			break
		}
	}

	if metadataResult == nil {
		t.Fatal("Metadata change not detected")
	}

	if metadataResult.ChangeType != "modified" {
		t.Errorf("Expected change type 'modified', got '%s'", metadataResult.ChangeType)
	}
}

func TestDiff_StringComparison(t *testing.T) {
	recipe1 := &models.UniversalRecipe{Name: "Original"}
	recipe2 := &models.UniversalRecipe{Name: "Modified"}

	results, err := Diff(recipe1, recipe2, 0.001)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	// Find Name result
	var nameResult *DiffResult
	for _, r := range results {
		if r.Field == "Name" {
			nameResult = &r
			break
		}
	}

	if nameResult == nil {
		t.Fatal("Name change not detected")
	}

	if !nameResult.Significant {
		t.Error("String changes should always be significant")
	}
}

func TestDiff_NilRecipes(t *testing.T) {
	recipe := &models.UniversalRecipe{}

	_, err := Diff(nil, recipe, 0.001)
	if err == nil {
		t.Error("Expected error for nil recipe1")
	}

	_, err = Diff(recipe, nil, 0.001)
	if err == nil {
		t.Error("Expected error for nil recipe2")
	}
}

func TestFormatDiff_NoDifferences(t *testing.T) {
	results := []DiffResult{}

	output := FormatDiff(results, false, false)

	expected := "✓ No differences found"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFormatDiff_WithChanges(t *testing.T) {
	results := []DiffResult{
		{Field: "Contrast", OldValue: int64(0), NewValue: int64(15), ChangeType: "modified", Significant: true},
		{Field: "Vibrance", OldValue: nil, NewValue: int64(20), ChangeType: "added", Significant: false},
	}

	output := FormatDiff(results, false, false)

	// Check that output contains expected sections
	if !contains(output, "MODIFIED:") {
		t.Error("Output should contain MODIFIED section")
	}
	if !contains(output, "ADDED") {
		t.Error("Output should contain ADDED section")
	}
	if !contains(output, "Contrast: 0 → 15") {
		t.Error("Output should show Contrast change")
	}
	if !contains(output, "*significant") {
		t.Error("Output should mark significant change")
	}
	if !contains(output, "Summary: 1 modified (1 significant), 1 added, 0 removed") {
		t.Error("Summary should show correct counts")
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"nil", nil, "nil"},
		{"int64", int64(42), "42"},
		{"float64", 0.12345, "0.12"},
		{"string", "test", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.value)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
