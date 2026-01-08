package nksc

import (
	"strings"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
)

func TestNewNKSCRecipe(t *testing.T) {
	meta := &np3.Metadata{
		RawBytes: []byte("dummy payload"),
		Label:    "Test Label",
	}
	target := "image.nef"
	recipe := NewNKSCRecipe(meta, target)

	if recipe == nil {
		t.Fatal("NewNKSCRecipe returned nil")
	}
	if recipe.version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", recipe.version)
	}
	if recipe.targetNEF != target {
		t.Errorf("Expected targetNEF '%s', got '%s'", target, recipe.targetNEF)
	}
}

func TestNKSCRecipe_MarshalXML(t *testing.T) {
	meta := &np3.Metadata{
		RawBytes: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		Label:    "TestPictureControl",
	}
	targetName := "DSC_0001.NEF"
	recipe := NewNKSCRecipe(meta, targetName)

	xmlBytes, err := recipe.MarshalXML()
	if err != nil {
		t.Fatalf("MarshalXML failed: %v", err)
	}

	xmlStr := string(xmlBytes)

	// Verify essential components
	checks := []string{
		`<x:xmpmeta`, // Root
		`xmlns:nine="http://ns.nikon.com/nksc/1.0/nine/"`,                            // Namespace
		`rdf:about="DSC_0001.NEF"`,                                                   // Target NEF
		`nine:StepName>Picture Control</nine:StepName>`,                              // Step Name
		`nine:FilterParametersExportData>3q2+7w==</nine:FilterParametersExportData>`, // Base64(DEADBEEF)
	}

	for _, check := range checks {
		if !strings.Contains(xmlStr, check) {
			t.Errorf("XML output missing expected string: %s\nGot:\n%s", check, xmlStr)
		}
	}
}
