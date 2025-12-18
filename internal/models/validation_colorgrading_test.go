package models

import (
	"testing"
)

func TestValidateMidRangeSharpening(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"valid min", -5.0, false},
		{"valid max", 5.0, false},
		{"valid zero", 0.0, false},
		{"valid mid positive", 2.5, false},
		{"valid mid negative", -2.5, false},
		{"invalid too low", -5.1, true},
		{"invalid too high", 5.1, true},
		{"invalid way too low", -10.0, true},
		{"invalid way too high", 10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMidRangeSharpening(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMidRangeSharpening(%f) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColorGradingHue(t *testing.T) {
	tests := []struct {
		name    string
		hue     int
		wantErr bool
	}{
		{"valid min", 0, false},
		{"valid max", 360, false},
		{"valid mid", 180, false},
		{"valid low", 45, false},
		{"valid high", 270, false},
		{"invalid negative", -1, true},
		{"invalid too high", 361, true},
		{"invalid way negative", -90, true},
		{"invalid way too high", 720, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColorGradingHue(tt.hue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColorGradingHue(%d) error = %v, wantErr %v", tt.hue, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColorGradingChroma(t *testing.T) {
	tests := []struct {
		name    string
		chroma  int
		wantErr bool
	}{
		{"valid min", -100, false},
		{"valid max", 100, false},
		{"valid zero", 0, false},
		{"valid mid positive", 50, false},
		{"valid mid negative", -50, false},
		{"invalid too low", -101, true},
		{"invalid too high", 101, true},
		{"invalid way too low", -200, true},
		{"invalid way too high", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColorGradingChroma(tt.chroma)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColorGradingChroma(%d) error = %v, wantErr %v", tt.chroma, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColorGradingBrightness(t *testing.T) {
	tests := []struct {
		name       string
		brightness int
		wantErr    bool
	}{
		{"valid min", -100, false},
		{"valid max", 100, false},
		{"valid zero", 0, false},
		{"valid mid positive", 50, false},
		{"valid mid negative", -50, false},
		{"invalid too low", -101, true},
		{"invalid too high", 101, true},
		{"invalid way too low", -200, true},
		{"invalid way too high", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColorGradingBrightness(tt.brightness)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColorGradingBrightness(%d) error = %v, wantErr %v", tt.brightness, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColorGradingBlending(t *testing.T) {
	tests := []struct {
		name     string
		blending int
		wantErr  bool
	}{
		{"valid min", 0, false},
		{"valid max", 100, false},
		{"valid mid", 50, false},
		{"valid low", 25, false},
		{"valid high", 75, false},
		{"invalid negative", -1, true},
		{"invalid too high", 101, true},
		{"invalid way negative", -50, true},
		{"invalid way too high", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColorGradingBlending(tt.blending)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColorGradingBlending(%d) error = %v, wantErr %v", tt.blending, err, tt.wantErr)
			}
		})
	}
}

func TestValidateColorGradingBalance(t *testing.T) {
	tests := []struct {
		name    string
		balance int
		wantErr bool
	}{
		{"valid min", -100, false},
		{"valid max", 100, false},
		{"valid zero", 0, false},
		{"valid mid positive", 50, false},
		{"valid mid negative", -50, false},
		{"invalid too low", -101, true},
		{"invalid too high", 101, true},
		{"invalid way too low", -200, true},
		{"invalid way too high", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColorGradingBalance(tt.balance)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColorGradingBalance(%d) error = %v, wantErr %v", tt.balance, err, tt.wantErr)
			}
		})
	}
}

// Test comprehensive ColorGrading validation
func TestColorGradingFullValidation(t *testing.T) {
	// Valid ColorGrading structure
	validGrading := &ColorGrading{
		Highlights: ColorGradingZone{
			Hue:        180,
			Chroma:     20,
			Brightness: -10,
		},
		Midtone: ColorGradingZone{
			Hue:        90,
			Chroma:     -15,
			Brightness: 5,
		},
		Shadows: ColorGradingZone{
			Hue:        270,
			Chroma:     30,
			Brightness: -20,
		},
		Blending: 50,
		Balance:  10,
	}

	// Validate all fields
	if err := ValidateColorGradingHue(validGrading.Highlights.Hue); err != nil {
		t.Errorf("Highlights Hue validation failed: %v", err)
	}
	if err := ValidateColorGradingChroma(validGrading.Highlights.Chroma); err != nil {
		t.Errorf("Highlights Chroma validation failed: %v", err)
	}
	if err := ValidateColorGradingBrightness(validGrading.Highlights.Brightness); err != nil {
		t.Errorf("Highlights Brightness validation failed: %v", err)
	}

	if err := ValidateColorGradingHue(validGrading.Midtone.Hue); err != nil {
		t.Errorf("Midtone Hue validation failed: %v", err)
	}
	if err := ValidateColorGradingChroma(validGrading.Midtone.Chroma); err != nil {
		t.Errorf("Midtone Chroma validation failed: %v", err)
	}
	if err := ValidateColorGradingBrightness(validGrading.Midtone.Brightness); err != nil {
		t.Errorf("Midtone Brightness validation failed: %v", err)
	}

	if err := ValidateColorGradingHue(validGrading.Shadows.Hue); err != nil {
		t.Errorf("Shadows Hue validation failed: %v", err)
	}
	if err := ValidateColorGradingChroma(validGrading.Shadows.Chroma); err != nil {
		t.Errorf("Shadows Chroma validation failed: %v", err)
	}
	if err := ValidateColorGradingBrightness(validGrading.Shadows.Brightness); err != nil {
		t.Errorf("Shadows Brightness validation failed: %v", err)
	}

	if err := ValidateColorGradingBlending(validGrading.Blending); err != nil {
		t.Errorf("Blending validation failed: %v", err)
	}
	if err := ValidateColorGradingBalance(validGrading.Balance); err != nil {
		t.Errorf("Balance validation failed: %v", err)
	}
}
