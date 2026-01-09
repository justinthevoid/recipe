// Package main is the entry point for the recipe-nx binary.
//
// recipe-nx is a standalone tool designed to integrate with Nikon NX Studio
// via the "Open with..." functionality. It receives a list of NEF files,
// parses a sidecar "NKSC" recipe (XML format), and applies the edits to the
// target NEF files by generating or updating ".nksc" sidecar files.
//
// Architecture:
//
//	cmd/nx (main)
//	   ↓
//	internal/batch (Orchestrator)
//	   ↓
//	internal/formats/nksc (Facade)
//	   ↓
//	internal/formats/np3 (Core Data)
//
// This binary is distinct from the primary "recipe" tool to ensure clean
// separation of concerns and independent release cycles.
package main
