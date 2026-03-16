package main

import (
	"encoding/json"
	"testing"
)

// TestConversionResultJSONMarshaling tests JSON marshaling of ConversionResult
func TestConversionResultJSONMarshaling(t *testing.T) {
	tests := []struct {
		name   string
		result ConversionResult
		verify func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "successful conversion with all fields",
			result: ConversionResult{
				Input:         "test.xmp",
				Output:        "test.np3",
				SourceFormat:  "xmp",
				TargetFormat:  "np3",
				Success:       true,
				DurationMs:    15,
				FileSizeBytes: 1234,
				Warnings:      []string{"Warning 1", "Warning 2"},
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				if data["input"] != "test.xmp" {
					t.Errorf("input = %v, want test.xmp", data["input"])
				}
				if data["output"] != "test.np3" {
					t.Errorf("output = %v, want test.np3", data["output"])
				}
				if data["source_format"] != "xmp" {
					t.Errorf("source_format = %v, want xmp", data["source_format"])
				}
				if data["target_format"] != "np3" {
					t.Errorf("target_format = %v, want np3", data["target_format"])
				}
				if data["success"] != true {
					t.Errorf("success = %v, want true", data["success"])
				}
				if data["duration_ms"] != float64(15) {
					t.Errorf("duration_ms = %v, want 15", data["duration_ms"])
				}
				if data["file_size_bytes"] != float64(1234) {
					t.Errorf("file_size_bytes = %v, want 1234", data["file_size_bytes"])
				}
				warnings := data["warnings"].([]interface{})
				if len(warnings) != 2 {
					t.Errorf("len(warnings) = %v, want 2", len(warnings))
				}
			},
		},
		{
			name: "failed conversion with error",
			result: ConversionResult{
				Input:        "corrupted.xmp",
				Output:       "",
				SourceFormat: "xmp",
				TargetFormat: "np3",
				Success:      false,
				DurationMs:   5,
				Error:        "parse error: invalid XML",
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				if data["success"] != false {
					t.Errorf("success = %v, want false", data["success"])
				}
				if data["error"] != "parse error: invalid XML" {
					t.Errorf("error = %v, want 'parse error: invalid XML'", data["error"])
				}
				// file_size_bytes should not be present (omitempty)
				if _, exists := data["file_size_bytes"]; exists {
					t.Errorf("file_size_bytes should not be present for failed conversion")
				}
			},
		},
		{
			name: "successful conversion without optional fields",
			result: ConversionResult{
				Input:        "test.xmp",
				Output:       "test.np3",
				SourceFormat: "xmp",
				TargetFormat: "np3",
				Success:      true,
				DurationMs:   10,
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				if data["success"] != true {
					t.Errorf("success = %v, want true", data["success"])
				}
				// Optional fields should not be present (omitempty)
				if _, exists := data["file_size_bytes"]; exists {
					t.Errorf("file_size_bytes should not be present (omitempty)")
				}
				if _, exists := data["warnings"]; exists {
					t.Errorf("warnings should not be present (omitempty)")
				}
				if _, exists := data["error"]; exists {
					t.Errorf("error should not be present (omitempty)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Unmarshal to map for verification
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify fields
			tt.verify(t, parsed)
		})
	}
}

// TestBatchResultJSONMarshaling tests JSON marshaling of BatchResult
func TestBatchResultJSONMarshaling(t *testing.T) {
	result := BatchResult{
		Batch:        true,
		Total:        3,
		SuccessCount: 2,
		ErrorCount:   1,
		DurationMs:   45,
		Results: []ConversionResult{
			{
				Input:         "file1.xmp",
				Output:        "file1.np3",
				SourceFormat:  "xmp",
				TargetFormat:  "np3",
				Success:       true,
				DurationMs:    15,
				FileSizeBytes: 1234,
			},
			{
				Input:        "file2.xmp",
				Output:       "file2.np3",
				SourceFormat: "xmp",
				TargetFormat: "np3",
				Success:      true,
				DurationMs:   18,
			},
			{
				Input:        "file3.xmp",
				Output:       "",
				SourceFormat: "xmp",
				TargetFormat: "np3",
				Success:      false,
				DurationMs:   12,
				Error:        "parse error",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Unmarshal to map for verification
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify batch fields
	if parsed["batch"] != true {
		t.Errorf("batch = %v, want true", parsed["batch"])
	}
	if parsed["total"] != float64(3) {
		t.Errorf("total = %v, want 3", parsed["total"])
	}
	if parsed["success_count"] != float64(2) {
		t.Errorf("success_count = %v, want 2", parsed["success_count"])
	}
	if parsed["error_count"] != float64(1) {
		t.Errorf("error_count = %v, want 1", parsed["error_count"])
	}
	if parsed["duration_ms"] != float64(45) {
		t.Errorf("duration_ms = %v, want 45", parsed["duration_ms"])
	}

	// Verify results array
	results := parsed["results"].([]interface{})
	if len(results) != 3 {
		t.Errorf("len(results) = %v, want 3", len(results))
	}

	// Verify each result has required fields
	for i, r := range results {
		result := r.(map[string]interface{})
		if result["input"] == nil {
			t.Errorf("results[%d].input is nil", i)
		}
		if result["success"] == nil {
			t.Errorf("results[%d].success is nil", i)
		}
	}
}

// TestJSONFieldNaming verifies snake_case field names
func TestJSONFieldNaming(t *testing.T) {
	result := ConversionResult{
		Input:         "test.xmp",
		Output:        "test.np3",
		SourceFormat:  "xmp",
		TargetFormat:  "np3",
		Success:       true,
		DurationMs:    15,
		FileSizeBytes: 1234,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// Verify snake_case field names (not camelCase)
	expectedFields := []string{
		"\"input\"",
		"\"output\"",
		"\"source_format\"",
		"\"target_format\"",
		"\"success\"",
		"\"duration_ms\"",
		"\"file_size_bytes\"",
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON does not contain field %s (expected snake_case)", field)
		}
	}

	// Verify no camelCase fields
	unwantedFields := []string{
		"\"sourceFormat\"",
		"\"targetFormat\"",
		"\"durationMs\"",
		"\"fileSizeBytes\"",
	}

	for _, field := range unwantedFields {
		if contains(jsonStr, field) {
			t.Errorf("JSON contains camelCase field %s (should be snake_case)", field)
		}
	}
}

// TestFormatMilliseconds tests millisecond formatting
func TestFormatMilliseconds(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{5, "5ms"},
		{15, "15ms"},
		{999, "999ms"},
		{1000, "1.00s"},
		{1500, "1.50s"},
		{2345, "2.34s"}, // Note: %.2f truncates, doesn't round up
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatMilliseconds(tt.ms)
			if got != tt.want {
				t.Errorf("formatMilliseconds(%d) = %v, want %v", tt.ms, got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkJSONMarshalConversionResult benchmarks JSON marshaling of ConversionResult
// Tech spec requirement: <10ms overhead for JSON output
func BenchmarkJSONMarshalConversionResult(b *testing.B) {
	result := ConversionResult{
		Input:         "test.xmp",
		Output:        "test.np3",
		SourceFormat:  "xmp",
		TargetFormat:  "np3",
		Success:       true,
		DurationMs:    15,
		FileSizeBytes: 1234,
		Warnings:      []string{"Warning 1", "Warning 2"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(result)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

// BenchmarkJSONMarshalIndentConversionResult benchmarks pretty-printed JSON marshaling
// This is what the actual CLI uses
func BenchmarkJSONMarshalIndentConversionResult(b *testing.B) {
	result := ConversionResult{
		Input:         "test.xmp",
		Output:        "test.np3",
		SourceFormat:  "xmp",
		TargetFormat:  "np3",
		Success:       true,
		DurationMs:    15,
		FileSizeBytes: 1234,
		Warnings:      []string{"Warning 1", "Warning 2"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			b.Fatalf("json.MarshalIndent() error = %v", err)
		}
	}
}

// BenchmarkJSONMarshalBatchResult benchmarks JSON marshaling of BatchResult with 100 results
// Validates performance for typical batch operations
func BenchmarkJSONMarshalBatchResult(b *testing.B) {
	// Create batch result with 100 conversions (typical batch size)
	results := make([]ConversionResult, 100)
	for i := 0; i < 100; i++ {
		results[i] = ConversionResult{
			Input:         "file.xmp",
			Output:        "file.np3",
			SourceFormat:  "xmp",
			TargetFormat:  "np3",
			Success:       true,
			DurationMs:    15,
			FileSizeBytes: 1234,
		}
	}

	batchResult := BatchResult{
		Batch:        true,
		Total:        100,
		SuccessCount: 100,
		ErrorCount:   0,
		DurationMs:   1500,
		Results:      results,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.MarshalIndent(batchResult, "", "  ")
		if err != nil {
			b.Fatalf("json.MarshalIndent() error = %v", err)
		}
	}
}
