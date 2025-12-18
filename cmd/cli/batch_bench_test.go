package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// BenchmarkBatch100Files benchmarks batch conversion of 100 files (AC-3)
func BenchmarkBatch100Files(b *testing.B) {
	// Setup: Create 100 test files
	tmpDir := b.TempDir()
	files := make([]string, 100)

	for i := 0; i < 100; i++ {
		filename := fmt.Sprintf("file%d.xmp", i)
		path := filepath.Join(tmpDir, filename)
		createValidXMPFileBench(b, path)
		files[i] = path
	}

	flags := BatchFlags{
		To:        "np3",
		Parallel:  runtime.NumCPU(),
		Overwrite: true,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := processBatch(files, flags)
		if err != nil {
			b.Fatalf("processBatch failed: %v", err)
		}

		// Clean up output files for next iteration
		if i < b.N-1 {
			for _, file := range files {
				outputPath := generateOutputPath(file, "np3")
				os.Remove(outputPath)
			}
		}
	}
}

// BenchmarkBatch100FilesSequential benchmarks sequential processing for comparison
func BenchmarkBatch100FilesSequential(b *testing.B) {
	// Setup: Create 100 test files
	tmpDir := b.TempDir()
	files := make([]string, 100)

	for i := 0; i < 100; i++ {
		filename := fmt.Sprintf("file%d.xmp", i)
		path := filepath.Join(tmpDir, filename)
		createValidXMPFileBench(b, path)
		files[i] = path
	}

	flags := BatchFlags{
		To:        "np3",
		Parallel:  1, // Sequential: 1 worker
		Overwrite: true,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := processBatch(files, flags)
		if err != nil {
			b.Fatalf("processBatch failed: %v", err)
		}

		// Clean up output files for next iteration
		if i < b.N-1 {
			for _, file := range files {
				outputPath := generateOutputPath(file, "np3")
				os.Remove(outputPath)
			}
		}
	}
}

// BenchmarkBatchVariousWorkerCounts benchmarks different worker pool sizes
func BenchmarkBatchVariousWorkerCounts(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, runtime.NumCPU()}

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("%dworkers", workers), func(b *testing.B) {
			// Setup: Create 50 test files
			tmpDir := b.TempDir()
			files := make([]string, 50)

			for i := 0; i < 50; i++ {
				filename := fmt.Sprintf("file%d.xmp", i)
				path := filepath.Join(tmpDir, filename)
				createValidXMPFileBench(b, path)
				files[i] = path
			}

			flags := BatchFlags{
				To:        "np3",
				Parallel:  workers,
				Overwrite: true,
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := processBatch(files, flags)
				if err != nil {
					b.Fatalf("processBatch failed: %v", err)
				}

				// Clean up output files for next iteration
				if i < b.N-1 {
					for _, file := range files {
						outputPath := generateOutputPath(file, "np3")
						os.Remove(outputPath)
					}
				}
			}
		})
	}
}

// createValidXMPFileBench creates a valid XMP file for benchmarking
func createValidXMPFileBench(b *testing.B, path string) {
	b.Helper()
	// Minimal valid XMP content
	content := `<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description rdf:about="">
      <crs:Exposure2012>+0.50</crs:Exposure2012>
      <crs:Contrast2012>+10</crs:Contrast2012>
      <crs:Highlights2012>-20</crs:Highlights2012>
      <crs:Shadows2012>+15</crs:Shadows2012>
      <crs:Whites2012>+5</crs:Whites2012>
      <crs:Blacks2012>-5</crs:Blacks2012>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		b.Fatalf("failed to create test file %s: %v", path, err)
	}
}
