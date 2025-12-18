package models

import "testing"

func TestWarnLevel_String(t *testing.T) {
	tests := []struct {
		level WarnLevel
		want  string
	}{
		{WarnInfo, "Info"},
		{WarnAdvisory, "Advisory"},
		{WarnCritical, "Critical"},
		{WarnLevel(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.level.String()
		if got != tt.want {
			t.Errorf("WarnLevel(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestConversionResult_AddWarning(t *testing.T) {
	result := &ConversionResult{}

	result.AddWarning(WarnAdvisory, "Clarity", "50", "Test message", "Use something else")

	if len(result.Warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d", len(result.Warnings))
	}

	w := result.Warnings[0]
	if w.Level != WarnAdvisory {
		t.Errorf("Level = %v, want WarnAdvisory", w.Level)
	}
	if w.Parameter != "Clarity" {
		t.Errorf("Parameter = %q, want %q", w.Parameter, "Clarity")
	}
	if w.Value != "50" {
		t.Errorf("Value = %q, want %q", w.Value, "50")
	}
}

func TestConversionResult_HasCritical(t *testing.T) {
	result := &ConversionResult{}

	// No warnings
	if result.HasCritical() {
		t.Error("HasCritical() = true, want false (no warnings)")
	}

	// Add advisory warning
	result.AddWarning(WarnAdvisory, "Test", "", "Test", "")
	if result.HasCritical() {
		t.Error("HasCritical() = true, want false (only advisory)")
	}

	// Add critical warning
	result.AddWarning(WarnCritical, "Test2", "", "Test", "")
	if !result.HasCritical() {
		t.Error("HasCritical() = false, want true (has critical)")
	}
}

func TestConversionResult_HasWarnings(t *testing.T) {
	result := &ConversionResult{}

	if result.HasWarnings() {
		t.Error("HasWarnings() = true, want false (no warnings)")
	}

	result.AddWarning(WarnInfo, "Test", "", "Test", "")
	if !result.HasWarnings() {
		t.Error("HasWarnings() = false, want true")
	}
}

func TestConversionResult_ByLevel(t *testing.T) {
	result := &ConversionResult{}
	result.AddWarning(WarnInfo, "Info1", "", "", "")
	result.AddWarning(WarnAdvisory, "Advisory1", "", "", "")
	result.AddWarning(WarnCritical, "Critical1", "", "", "")
	result.AddWarning(WarnAdvisory, "Advisory2", "", "", "")

	// Get all warnings at Advisory or higher
	advisoryPlus := result.ByLevel(WarnAdvisory)
	if len(advisoryPlus) != 3 {
		t.Errorf("ByLevel(Advisory) returned %d warnings, want 3", len(advisoryPlus))
	}

	// Get only Critical
	critical := result.ByLevel(WarnCritical)
	if len(critical) != 1 {
		t.Errorf("ByLevel(Critical) returned %d warnings, want 1", len(critical))
	}
}

func TestUnsupportedXMPParameters(t *testing.T) {
	// Check that key parameters are defined
	params := []string{"Clarity", "Texture", "Dehaze", "ToneCurvePV2012Red"}

	for _, param := range params {
		if _, exists := UnsupportedXMPParameters[param]; !exists {
			t.Errorf("Expected %q in UnsupportedXMPParameters", param)
		}
	}

	// Check Clarity is Advisory level
	if UnsupportedXMPParameters["Clarity"].Level != WarnAdvisory {
		t.Error("Clarity should be WarnAdvisory level")
	}

	// Check RGB curves are Critical level
	if UnsupportedXMPParameters["ToneCurvePV2012Red"].Level != WarnCritical {
		t.Error("ToneCurvePV2012Red should be WarnCritical level")
	}
}
