package np3

import (
	"testing"
)

func TestCameraProfileHeuristic(t *testing.T) {
	tests := []struct {
		name            string
		presetName      string
		expectedProfile string
	}{
		{"Standard", "My Preset", "Camera Standard"},
		{"Neutral", "Nikon Neutral 01", "Camera Neutral"},
		{"Flat", "Flat Profile", "Camera Flat"},
		{"Monochrome", "BW Mono", "Camera Monochrome"},
		{"Portrait", "Portrait Soft", "Camera Portrait"},
		{"Landscape", "Landscape Vivid", "Camera Landscape"},
		{"Vivid", "Super Vivid", "Camera Vivid"},
		{"Case Insensitive", "nEuTrAl", "Camera Neutral"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal valid NP3 file structure
			// Magic bytes (3) + Version (4) + Padding (13) + Name Header (4) + Name (40) + Padding (2) + Chunks...
			// Total size must be >= 300
			data := make([]byte, 480)
			copy(data[0:], []byte{'N', 'C', 'P'}) // Magic

			// Set name
			nameBytes := []byte(tt.presetName)
			copy(data[20:], nameBytes)

			// Parse
			recipe, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if recipe.CameraProfileName != tt.expectedProfile {
				t.Errorf("expected profile %q, got %q", tt.expectedProfile, recipe.CameraProfileName)
			}
		})
	}
}
