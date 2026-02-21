package testutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/testutil"
	"gopkg.in/yaml.v3"
)

func TestFixtureManifestSchema(t *testing.T) {
	// Read the sample manifest
	manifestPath := filepath.Join("..", "..", "testdata", "nx-fixtures", "manifest_schema_test.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("Failed to read test manifest: %v", err)
	}

	// This struct should be defined in internal/testutil/fixtures.go
	var manifest testutil.FixtureManifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	if len(manifest.Fixtures) != 1 {
		t.Fatalf("Expected 1 fixture, got %d", len(manifest.Fixtures))
	}

	f := manifest.Fixtures[0]
	if f.ID != "z6_base_iso" {
		t.Errorf("Expected ID 'z6_base_iso', got '%s'", f.ID)
	}
	if f.CameraModel != "Nikon Z 6" {
		t.Errorf("Expected CameraModel 'Nikon Z 6', got '%s'", f.CameraModel)
	}
	if f.URL != "https://example.com/fixtures/z6_base_iso.nef" {
		t.Errorf("Expected URL 'https://example.com/fixtures/z6_base_iso.nef', got '%s'", f.URL)
	}
	// Verify AC details
	if f.Lens == "" {
		t.Error("Lens field is missing or empty")
	}
	if f.NXStudioVersion == "" {
		t.Error("NXStudioVersion field is missing or empty")
	}
	if f.NP3Variant == "" {
		t.Error("NP3Variant field is missing or empty")
	}
	if f.SHA256 == "" {
		t.Error("SHA256 field is missing or empty")
	}
}
