package nksc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
)

func TestWrite_Atomic(t *testing.T) {
	// Create a temporary directory for output
	tmpDir, err := os.MkdirTemp("", "nksc_write_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "output.nksc")

	// Create dummy NP3 metadata
	meta := &np3.Metadata{
		Label:    "Test Recipe",
		RawBytes: []byte("dummy bytes"),
	}

	recipe := NewNKSCRecipe(meta, "test.nef")

	// Call Write
	err = recipe.Write(outputPath)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("File was not created at %s", outputPath)
	}

	// Verify content is not empty
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Output file is empty")
	}
}

func TestWrite_Structure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nksc_structure_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "structure.nksc")
	meta := &np3.Metadata{
		Label:    "Structure Test",
		RawBytes: []byte{0x00, 0x01}, // minimal
	}
	recipe := NewNKSCRecipe(meta, "structure.nef")

	err = recipe.Write(outputPath)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	strContent := string(content)

	// Check Header
	expectedHeader := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
	if !strings.Contains(strContent, expectedHeader) {
		t.Errorf("Missing XML header")
	}

	// Check XMP Packet Wrapper
	// Note: The BOM might make exact string matching tricky, checking for parts
	if !strings.Contains(strContent, "<?xpacket begin=") {
		t.Errorf("Missing XMP packet begin")
	}
	if !strings.Contains(strContent, `id="W5M0MpCehiHzreSzNTczkc9d"`) {
		t.Errorf("Missing XMP packet ID")
	}
	if !strings.Contains(strContent, `<?xpacket end="w"?>`) {
		t.Errorf("Missing XMP packet end")
	}

	// Check Root Element
	if !strings.Contains(strContent, `<x:xmpmeta`) {
		t.Errorf("Missing x:xmpmeta root")
	}

	// Check Namespaces
	if !strings.Contains(strContent, `xmlns:ast="http://ns.nikon.com/nksc/1.0/asteria/"`) {
		t.Errorf("Missing ast namespace definition")
	}
	if !strings.Contains(strContent, `xmlns:nine="http://ns.nikon.com/nksc/1.0/nine/"`) {
		t.Errorf("Missing nine namespace definition")
	}

	// Double check validity by unmarshalling (skipping wrappers)
	// We need to strip the wrappers to parse with standard xml decoder if we want to validate structure
	// But for now, text search is enough for the wrapper checks.

	// Validate indentation (heuristic: look for newline and spaces)
	if !strings.Contains(strContent, "\n  <rdf:RDF") {
		t.Errorf("Output does not appear to be indented")
	}
}
