// Package inspect provides tools for inspecting and comparing photo editing recipes.
package inspect

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/justin/recipe/internal/models"
)

// DiffResult represents a single field difference between two UniversalRecipe instances.
type DiffResult struct {
	Field      string      `json:"field"`       // Field name
	OldValue   interface{} `json:"old_value"`   // Value in first recipe (nil if added)
	NewValue   interface{} `json:"new_value"`   // Value in second recipe (nil if removed)
	ChangeType string      `json:"change_type"` // "modified", "added", "removed", or "unchanged"
	Significant bool       `json:"significant"` // true if change is >5% or non-numeric change
}

// DiffOutput represents the complete diff output for JSON formatting.
type DiffOutput struct {
	File1              string       `json:"file1"`
	File2              string       `json:"file2"`
	Changes            int          `json:"changes"`
	SignificantChanges int          `json:"significant_changes"`
	FieldsCompared     int          `json:"fields_compared"`
	Differences        []DiffResult `json:"differences"`
	Unchanged          []string     `json:"unchanged,omitempty"`
}

// Diff compares two UniversalRecipe instances field-by-field and returns the differences.
// tolerance is the threshold for float comparison (default: 0.001).
func Diff(recipe1, recipe2 *models.UniversalRecipe, tolerance float64) ([]DiffResult, error) {
	if recipe1 == nil || recipe2 == nil {
		return nil, fmt.Errorf("cannot diff nil recipes")
	}

	var results []DiffResult

	// Get reflection values
	v1 := reflect.ValueOf(*recipe1)
	v2 := reflect.ValueOf(*recipe2)
	t := v1.Type()

	// Iterate over all struct fields
	for i := 0; i < v1.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		val1 := v1.Field(i)
		val2 := v2.Field(i)

		// Compare values
		diff := compareValues(field.Name, val1, val2, tolerance)
		if diff != nil {
			results = append(results, *diff)
		}
	}

	return results, nil
}

// compareValues compares two reflect.Value instances and returns a DiffResult if they differ.
func compareValues(fieldName string, val1, val2 reflect.Value, tolerance float64) *DiffResult {
	// Handle different types
	switch val1.Kind() {
	case reflect.Int:
		old := val1.Int()
		new := val2.Int()
		if old != new {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    old,
				NewValue:    new,
				ChangeType:  "modified",
				Significant: isSignificantChange(float64(old), float64(new)),
			}
		}

	case reflect.Float64:
		old := val1.Float()
		new := val2.Float()
		if math.Abs(old-new) > tolerance {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    old,
				NewValue:    new,
				ChangeType:  "modified",
				Significant: isSignificantChange(old, new),
			}
		}

	case reflect.String:
		old := val1.String()
		new := val2.String()
		if old != new {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    old,
				NewValue:    new,
				ChangeType:  "modified",
				Significant: true, // String changes always significant
			}
		}

	case reflect.Ptr:
		// Handle pointer fields (e.g., *int for Temperature)
		if val1.IsNil() && val2.IsNil() {
			return nil
		}
		if val1.IsNil() && !val2.IsNil() {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    nil,
				NewValue:    val2.Elem().Interface(),
				ChangeType:  "added",
				Significant: true,
			}
		}
		if !val1.IsNil() && val2.IsNil() {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    val1.Elem().Interface(),
				NewValue:    nil,
				ChangeType:  "removed",
				Significant: true,
			}
		}
		// Both non-nil, compare underlying values
		return compareValues(fieldName, val1.Elem(), val2.Elem(), tolerance)

	case reflect.Struct:
		// Compare nested structs (ColorAdjustment, CameraProfile, etc.)
		return compareStructs(fieldName, val1, val2, tolerance)

	case reflect.Slice:
		// Compare slices (ToneCurvePoint[], etc.)
		return compareSlices(fieldName, val1, val2, tolerance)

	case reflect.Map:
		// Compare maps (Metadata)
		return compareMaps(fieldName, val1, val2)
	}

	return nil
}

// compareStructs compares two struct values recursively.
func compareStructs(fieldName string, val1, val2 reflect.Value, tolerance float64) *DiffResult {
	// Check if structs are equal by comparing all fields
	t := val1.Type()
	hasChanges := false
	var changedFields []string

	for i := 0; i < val1.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}

		subField := t.Field(i)
		subVal1 := val1.Field(i)
		subVal2 := val2.Field(i)

		diff := compareValues(subField.Name, subVal1, subVal2, tolerance)
		if diff != nil {
			hasChanges = true
			changedFields = append(changedFields, subField.Name)
		}
	}

	if hasChanges {
		return &DiffResult{
			Field:       fieldName,
			OldValue:    val1.Interface(),
			NewValue:    val2.Interface(),
			ChangeType:  "modified",
			Significant: true, // Struct changes always significant
		}
	}

	return nil
}

// compareSlices compares two slice values.
func compareSlices(fieldName string, val1, val2 reflect.Value, tolerance float64) *DiffResult {
	len1 := val1.Len()
	len2 := val2.Len()

	// Both empty, no diff
	if len1 == 0 && len2 == 0 {
		return nil
	}

	// Length changed
	if len1 != len2 {
		return &DiffResult{
			Field:       fieldName,
			OldValue:    val1.Interface(),
			NewValue:    val2.Interface(),
			ChangeType:  "modified",
			Significant: true,
		}
	}

	// Compare element by element
	for i := 0; i < len1; i++ {
		elem1 := val1.Index(i)
		elem2 := val2.Index(i)

		// For structs, do deep comparison
		if elem1.Kind() == reflect.Struct {
			if diff := compareStructs(fmt.Sprintf("%s[%d]", fieldName, i), elem1, elem2, tolerance); diff != nil {
				return &DiffResult{
					Field:       fieldName,
					OldValue:    val1.Interface(),
					NewValue:    val2.Interface(),
					ChangeType:  "modified",
					Significant: true,
				}
			}
		} else {
			if !reflect.DeepEqual(elem1.Interface(), elem2.Interface()) {
				return &DiffResult{
					Field:       fieldName,
					OldValue:    val1.Interface(),
					NewValue:    val2.Interface(),
					ChangeType:  "modified",
					Significant: true,
				}
			}
		}
	}

	return nil
}

// compareMaps compares two map values.
func compareMaps(fieldName string, val1, val2 reflect.Value) *DiffResult {
	len1 := val1.Len()
	len2 := val2.Len()

	// Both empty, no diff
	if len1 == 0 && len2 == 0 {
		return nil
	}

	// Check all keys
	keys1 := val1.MapKeys()
	keys2 := val2.MapKeys()

	// Quick check: different number of keys
	if len(keys1) != len(keys2) {
		return &DiffResult{
			Field:       fieldName,
			OldValue:    val1.Interface(),
			NewValue:    val2.Interface(),
			ChangeType:  "modified",
			Significant: true,
		}
	}

	// Check each key-value pair
	for _, key := range keys1 {
		v1 := val1.MapIndex(key)
		v2 := val2.MapIndex(key)

		if !v2.IsValid() || !reflect.DeepEqual(v1.Interface(), v2.Interface()) {
			return &DiffResult{
				Field:       fieldName,
				OldValue:    val1.Interface(),
				NewValue:    val2.Interface(),
				ChangeType:  "modified",
				Significant: true,
			}
		}
	}

	return nil
}

// isSignificantChange determines if a numeric change is significant (>5%).
func isSignificantChange(old, new float64) bool {
	// If old is zero, any non-zero change is significant
	if old == 0 {
		return new != 0
	}

	// Calculate percentage change
	percentChange := math.Abs((new - old) / old)

	return percentChange > 0.05 // >5% is significant
}

// FormatDiff formats diff results as human-readable text.
// If unified is true, shows all fields including unchanged ones.
// If colorize is true, applies ANSI color codes.
func FormatDiff(results []DiffResult, unified, colorize bool) string {
	if len(results) == 0 {
		return "✓ No differences found"
	}

	var sb strings.Builder

	// Group by change type
	var modified, added, removed, unchanged []DiffResult
	for _, r := range results {
		switch r.ChangeType {
		case "modified":
			modified = append(modified, r)
		case "added":
			added = append(added, r)
		case "removed":
			removed = append(removed, r)
		case "unchanged":
			unchanged = append(unchanged, r)
		}
	}

	// Modified fields
	if len(modified) > 0 {
		sb.WriteString("MODIFIED:\n")
		for _, r := range modified {
			line := fmt.Sprintf("  %s: %v → %v", r.Field, formatValue(r.OldValue), formatValue(r.NewValue))
			if r.Significant {
				line += "                  *significant"
			}
			if colorize {
				line = colorizeOutput(line, "significant")
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	// Added fields
	if len(added) > 0 {
		sb.WriteString("ADDED (only in file 2):\n")
		for _, r := range added {
			line := fmt.Sprintf("  %s: %v", r.Field, formatValue(r.NewValue))
			if colorize {
				line = colorizeOutput(line, "added")
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	// Removed fields
	if len(removed) > 0 {
		sb.WriteString("REMOVED (only in file 1):\n")
		for _, r := range removed {
			line := fmt.Sprintf("  %s: %v", r.Field, formatValue(r.OldValue))
			if colorize {
				line = colorizeOutput(line, "removed")
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	// Unchanged fields (if unified mode)
	if unified && len(unchanged) > 0 {
		sb.WriteString("UNCHANGED:\n")
		for _, r := range unchanged {
			line := fmt.Sprintf("  %s: %v = %v", r.Field, formatValue(r.OldValue), formatValue(r.NewValue))
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	// Summary
	significantCount := 0
	for _, r := range modified {
		if r.Significant {
			significantCount++
		}
	}

	summary := fmt.Sprintf("Summary: %d modified", len(modified))
	if significantCount > 0 {
		summary += fmt.Sprintf(" (%d significant)", significantCount)
	}
	summary += fmt.Sprintf(", %d added, %d removed", len(added), len(removed))
	if unified {
		summary += fmt.Sprintf(", %d unchanged", len(unchanged))
	}

	sb.WriteString(summary)

	return sb.String()
}

// formatValue formats a value for display.
func formatValue(v interface{}) string {
	if v == nil {
		return "nil"
	}

	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", val)
	case int64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBold   = "\033[1m"
)

// colorizeOutput applies ANSI color codes to text based on the color type.
func colorizeOutput(text, colorType string) string {
	var color string
	switch colorType {
	case "added":
		color = ColorGreen
	case "removed":
		color = ColorRed
	case "significant":
		color = ColorYellow + ColorBold
	default:
		return text
	}

	return color + text + ColorReset
}
