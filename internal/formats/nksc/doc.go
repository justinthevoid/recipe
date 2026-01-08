// Package nksc implements the Nikon Key Store Container (NKSC) sidecar format.
//
// NKSC is an XML-based metadata format used by Nikon NX Studio to store
// non-destructive edits. This package provides functionality to wrap NP3
// data into the NKSC structure without full conversion, treating the NP3
// data as an embedded opaque blob or referenced payload.
//
// Architecture:
//
//	internal/formats/nksc
//	   ↓ (imports)
//	internal/formats/np3 (for Metadata struct)
//
// This package is part of the recipe-nx binary and enforces the "facade" pattern
// over the raw NP3 metadata.
package nksc
