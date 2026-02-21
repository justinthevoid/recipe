package batch_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/batch"
	"github.com/justin/recipe/internal/formats/np3"
)

func BenchmarkProcessBatch(b *testing.B) {
	// Setup inputs
	tmpDir := b.TempDir()
	outDir := b.TempDir()
	fileCount := 50 // Reduced for benchmark speed, scale up if needed

	// Generate dummy NEF files
	data := []byte("dummy nef content")
	for i := 0; i < fileCount; i++ {
		os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("bench_%d.nef", i)), data, 0644)
	}

	scenarios := []struct {
		name    string
		workers int
	}{
		{"Serial", 1},
		{"Parallel", 4}, // Fixed parallel count for consistent benchmark
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				iterOut := filepath.Join(outDir, fmt.Sprintf("%s_%d", sc.name, n))

				orch := batch.NewOrchestrator(batch.Config{
					InputPath:  tmpDir,
					OutputPath: iterOut,
					Recipe:     &np3.Metadata{Label: "Bench"},
					Workers:    sc.workers,
				})
				b.StartTimer()

				if _, err := orch.ProcessBatch(context.Background()); err != nil {
					b.Fatalf("Run failed: %v", err)
				}
			}
		})
	}
}
