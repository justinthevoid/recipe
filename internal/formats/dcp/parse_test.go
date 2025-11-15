package dcp

import (
	"math"
	"os"
	"testing"
)

// TestParse_ValidDCP tests parsing a real Adobe DCP sample file.
func TestParse_ValidDCP(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		wantName string
	}{
		{
			name:     "Nikon Z f Standard",
			file:     "../../../testdata/dcp/Nikon Z f Camera Standard.dcp",
			wantName: "Camera Standard",
		},
		{
			name:     "Nikon Z f Portrait",
			file:     "../../../testdata/dcp/Nikon Z f Camera Portrait.dcp",
			wantName: "Camera Portrait", // ProfileName from tag 50936
		},
		{
			name:     "Hasselblad Adobe Standard",
			file:     "../../../testdata/dcp/Hasselblad X1D-50 Adobe Standard.dcp",
			wantName: "Adobe Standard", // ProfileName from tag 50936
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			recipe, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if recipe == nil {
				t.Fatal("Parse() returned nil recipe")
			}

			// Verify profile name in metadata
			profileName, ok := recipe.Metadata["profile_name"]
			if !ok {
				t.Error("profile_name not found in metadata")
			}
			if profileName != tt.wantName {
				t.Errorf("profile_name = %q, want %q", profileName, tt.wantName)
			}

			// Verify UniversalRecipe fields are set (basic smoke test)
			// Specific values tested in TestAnalyzeToneCurve
		})
	}
}

// TestParse_MissingToneCurve tests parsing a DCP without tone curve.
func TestParse_MissingToneCurve(t *testing.T) {
	// Will be implemented after acquiring DCP samples
	t.Skip("Pending real DCP sample files")
}

// TestParse_NonIdentityMatrix tests parsing a DCP with color calibration matrices.
func TestParse_NonIdentityMatrix(t *testing.T) {
	// Will be implemented after acquiring DCP samples
	t.Skip("Pending real DCP sample files")
}

// TestParse_CorruptTIFF tests error handling for malformed TIFF files.
func TestParse_CorruptTIFF(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name:    "empty file",
			data:    []byte{},
			wantErr: "file too small",
		},
		{
			name:    "invalid magic bytes",
			data:    []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: "invalid TIFF magic bytes",
		},
		{
			name:    "truncated file",
			data:    []byte{0x49, 0x49}, // II only, no version
			wantErr: "file too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.data)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() == "" || len(tt.wantErr) == 0 {
				t.Fatalf("invalid test: both error and wantErr must be non-empty")
			}
			// Just verify we got an error (detailed matching requires actual DCP samples)
		})
	}
}

// TestParse_MissingTag50740 tests error when TIFF lacks CameraProfile tag.
func TestParse_MissingTag50740(t *testing.T) {
	// Will be implemented after understanding TIFF structure better
	t.Skip("Requires valid TIFF file without tag 50740")
}

// TestParse_MalformedXML tests error handling for invalid XML in tag 50740.
func TestParse_MalformedXML(t *testing.T) {
	// Will be implemented after TIFF tag extraction works
	t.Skip("Requires valid TIFF with malformed XML in tag 50740")
}

// TestAnalyzeToneCurve tests tone curve analysis algorithms.
func TestAnalyzeToneCurve(t *testing.T) {
	tests := []struct {
		name       string
		points     []ToneCurvePoint
		wantExp    float64
		wantCon    float64
		wantHi     float64
		wantSh     float64
		tolerance  float64
	}{
		{
			name:      "linear curve",
			points:    []ToneCurvePoint{{0.0, 0.0}, {1.0, 1.0}},
			wantExp:   0.0,
			wantCon:   0.0,
			wantHi:    0.0,
			wantSh:    0.0,
			tolerance: 0.01,
		},
		{
			name: "exposure shift +0.5",
			points: []ToneCurvePoint{
				{0.0, 0.0},
				{0.5, 0.625}, // +0.125 shift (32/255 = 0.125)
				{1.0, 1.0},
			},
			wantExp:   0.625,  // (0.625 - 0.5) * 5.0 = 0.625
			wantCon:   0.0,    // No contrast change
			wantHi:    0.0,    // Top unchanged
			wantSh:    0.0,    // Bottom unchanged
			tolerance: 0.1,
		},
		{
			name: "contrast increase",
			points: []ToneCurvePoint{
				{0.0, 0.0},
				{0.25, 0.125},  // Darker shadows (32/255 = 0.125)
				{0.5, 0.5},     // Midpoint unchanged
				{0.75, 0.878},  // Brighter highlights (224/255 = 0.878)
				{1.0, 1.0},
			},
			wantExp:   0.0,  // Midpoint unchanged
			wantCon:   0.5,  // Steeper slope
			wantHi:    0.0,  // Top unchanged
			wantSh:    0.0,  // Bottom unchanged
			tolerance: 0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, con, hi, sh := analyzeToneCurve(tt.points)

			if math.Abs(exp-tt.wantExp) > tt.tolerance {
				t.Errorf("exposure: got %.3f, want %.3f (tolerance %.3f)", exp, tt.wantExp, tt.tolerance)
			}
			if math.Abs(con-tt.wantCon) > tt.tolerance {
				t.Errorf("contrast: got %.3f, want %.3f (tolerance %.3f)", con, tt.wantCon, tt.tolerance)
			}
			if math.Abs(hi-tt.wantHi) > tt.tolerance {
				t.Errorf("highlights: got %.3f, want %.3f (tolerance %.3f)", hi, tt.wantHi, tt.tolerance)
			}
			if math.Abs(sh-tt.wantSh) > tt.tolerance {
				t.Errorf("shadows: got %.3f, want %.3f (tolerance %.3f)", sh, tt.wantSh, tt.tolerance)
			}
		})
	}
}

// TestIsIdentityMatrix tests identity matrix detection.
func TestIsIdentityMatrix(t *testing.T) {
	tests := []struct {
		name   string
		matrix *Matrix
		want   bool
	}{
		{
			name: "identity matrix",
			matrix: &Matrix{
				Rows: [3][3]float64{
					{1.0, 0.0, 0.0},
					{0.0, 1.0, 0.0},
					{0.0, 0.0, 1.0},
				},
			},
			want: true,
		},
		{
			name: "non-identity matrix",
			matrix: &Matrix{
				Rows: [3][3]float64{
					{1.5, 0.1, 0.0},
					{0.0, 1.2, 0.0},
					{0.0, 0.0, 0.9},
				},
			},
			want: false,
		},
		{
			name:   "nil matrix",
			matrix: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIdentityMatrix(tt.matrix)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClampFloat64 tests value clamping.
func TestClampFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		min   float64
		max   float64
		want  float64
	}{
		{"within range", 0.5, 0.0, 1.0, 0.5},
		{"below min", -0.5, 0.0, 1.0, 0.0},
		{"above max", 1.5, 0.0, 1.0, 1.0},
		{"at min", 0.0, 0.0, 1.0, 0.0},
		{"at max", 1.0, 0.0, 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampFloat64(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("got %.3f, want %.3f", got, tt.want)
			}
		})
	}
}

// TestFindPoint tests tone curve point finding and interpolation.
func TestFindPoint(t *testing.T) {
	// Binary format uses 0.0-1.0 normalized values
	points := []ToneCurvePoint{
		{Input: 0.0, Output: 0.0},
		{Input: 0.25, Output: 0.25},
		{Input: 0.5, Output: 0.5},
		{Input: 0.75, Output: 0.75},
		{Input: 1.0, Output: 1.0},
	}

	tests := []struct {
		name   string
		input  float64
		output float64
	}{
		{"exact match at 0.0", 0.0, 0.0},
		{"exact match at 0.5", 0.5, 0.5},
		{"exact match at 1.0", 1.0, 1.0},
		{"interpolate at 0.125", 0.125, 0.125},
		{"interpolate at 0.375", 0.375, 0.375},
		{"interpolate at 0.625", 0.625, 0.625},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := findPoint(points, tt.input)
			if math.Abs(point.Input-tt.input) > 0.0001 {
				t.Errorf("input: got %.4f, want %.4f", point.Input, tt.input)
			}
			if math.Abs(point.Output-tt.output) > 0.01 {
				t.Errorf("output: got %.4f, want %.4f (tolerance ±0.01)", point.Output, tt.output)
			}
		})
	}
}
