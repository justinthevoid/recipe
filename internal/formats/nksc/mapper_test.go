package nksc

import (
	"encoding/base64"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
)

func TestNewFromNP3(t *testing.T) {
	// Create sample metadata
	raw := []byte("fake_np3_binary_data_12345")
	// Note: We don't strictly need valid NP3 for this unit test if we just check mapping of fields
	meta := &np3.Metadata{
		RawBytes:   raw,
		Label:      "Test Style",
		Sharpening: 3.0,
		Contrast:   10,
	}

	nkscStruct, err := NewFromNP3(meta, "test.nef")
	if err != nil {
		t.Fatalf("NewFromNP3 failed: %v", err)
	}

	// Verify Structure
	if nkscStruct.RDF.Description.About != "test.nef" {
		t.Errorf("Expected About='test.nef', got '%s'", nkscStruct.RDF.Description.About)
	}

	// Verify NineEdits
	steps := nkscStruct.RDF.Description.NineEdits.Seq.Steps
	if len(steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(steps))
	}

	step := steps[0]
	if step.Name != "Picture Control" {
		t.Errorf("Expected Step Name 'Picture Control', got '%s'", step.Name)
	}

	// Verify Payload
	expectedPayload := base64.StdEncoding.EncodeToString(raw)
	if step.FilterParametersExportData != expectedPayload {
		t.Errorf("Payload mismatch. Got %s, want %s", step.FilterParametersExportData, expectedPayload)
	}
}
