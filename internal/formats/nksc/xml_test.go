package nksc

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestNKSC_Serialization(t *testing.T) {
	// Setup a basic NKSC structure
	// Note: We are testing the raw structs from model.go, not the NKSCRecipe facade yet.
	// We want to ensure the XML tags generate the correct Nikon namespaces.

	nksc := NewNKSC()

	// This is our "Golden Master" expectation for the namespaces
	// We expect proper XMP wrapping.
	output, err := xml.MarshalIndent(nksc, "", "  ")

	if err != nil {
		t.Fatalf("Failed to marshal NKSC: %v", err)
	}

	xmlStr := string(output)

	// Check for standard XMP wrapper
	if !strings.Contains(xmlStr, "<x:xmpmeta") {
		t.Errorf("Missing x:xmpmeta root tag. Got:\n%s", xmlStr)
	}

	// These checks will FAIL currently because implementation is empty
	expectedNamespaces := []string{
		`xmlns:x="adobe:ns:meta/"`,
		`xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"`,
		`xmlns:ast="http://ns.nikon.com/nksc/1.0/asteria/"`,
		`xmlns:nine="http://ns.nikon.com/nksc/1.0/nine/"`,
	}

	for _, ns := range expectedNamespaces {
		if !strings.Contains(xmlStr, ns) {
			t.Errorf("Missing namespace definition: %s", ns)
		}
	}
}
