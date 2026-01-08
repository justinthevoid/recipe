// Package batch implements the parallel worker pool for processing NEF files.
//
// It orchestrates the reading of input files, generation of NKSC sidecars,
// and handling of concurrency limits. The package ensures thread-safe
// operation and efficient resource utilization (Goroutine pool).
//
// Architecture:
//
//	internal/batch
//	   ↓
//	internal/formats/nksc (for recipe generation)
//
// This package is the core engine for the recipe-nx "batch" command.
package batch
