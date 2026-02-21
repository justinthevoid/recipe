package nksc

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/justin/recipe/internal/formats/np3"
)

// NKSCRecipe is a facade over the raw NP3 metadata, formatted for Nikon NX Studio.
type NKSCRecipe struct {
	metadata  *np3.Metadata // Source data
	targetNEF string        // The NEF file this sidecar applies to
	version   string        // NKSC schema version (default: "1.0")
}

// NewNKSCRecipe creates a new NKSC recipe wrapper for the given NP3 metadata.
func NewNKSCRecipe(metadata *np3.Metadata, targetNEF string) *NKSCRecipe {
	return &NKSCRecipe{
		metadata:  metadata,
		targetNEF: targetNEF,
		version:   "1.0",
	}
}

// MarshalXML serializes NKSCRecipe to NKSC XMP format.
func (r *NKSCRecipe) MarshalXML() ([]byte, error) {
	// convert internal metadata to the XML-mapped struct
	nksc, err := NewFromNP3(r.metadata, r.targetNEF)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}

	// serialize with indentation
	output, err := xml.MarshalIndent(nksc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("xml marshal failed: %w", err)
	}

	// encoding/xml does not add the XML header, but XMP is usually embedded or requires it?
	// The tests expect <x:xmpmeta at the root, which our struct provides.
	// We might want to prepend the standard XML header if strictly required for files,
	// checking standard lib behavior. For now, returning the XMP block is correct.

	return output, nil
}

// Write serializes the recipe to an NKSC sidecar file at the specified path.
// It uses atomic file writing (temp file + rename) to avoid partial writes.
func (r *NKSCRecipe) Write(path string) error {
	data, err := r.MarshalXML()
	if err != nil {
		return err
	}

	// Construct the full XMP packet
	// Header
	// const xmlHeader = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n"
	// XMP Packet Marker (with BOM)
	const xmpBegin = "<?xpacket begin=\"\xEF\xBB\xBF\" id=\"W5M0MpCehiHzreSzNTczkc9d\"?>\n"
	const xmpEnd = "\n<?xpacket end=\"w\"?>"

	// NX Studio output does not include the XML declaration, only the XMP packet wrapper.
	fullContent := []byte(xmpBegin + string(data) + xmpEnd)

	// Atomic write
	dir := filepath.Dir(path)
	// Ensure directory exists
	// Note: os.CreateTemp requires the directory to exist if specified.
	// But usually the caller handles the directory existence, or we can try to create it.
	// The story AC doesn't strictly say we must create the directory, but good practice.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		// Clean up on failure (if rename hasn't happened yet, file still exists)
		// os.Remove handles "file not found" poorly if we don't check errors, but
		// if we renamed it, it's gone from tmp path.
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(fullContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("failed to rename temp file to %s: %w", path, err)
	}

	return nil
}
