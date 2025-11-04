package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/justin/recipe/internal/formats/np3"
)

func main() {
	files, _ := filepath.Glob("examples/np3/**/*.np3")
	seen := make(map[string]int)

	for i, f := range files {
		if i >= 15 { break } // Check first 15 files

		data, _ := os.ReadFile(f)
		r, err := np3.Parse(data)
		if err == nil {
			key := fmt.Sprintf("S%d-C%d-Sa%d", r.Sharpness, r.Contrast, r.Saturation)
			seen[key]++
			fmt.Printf("%-30s: %s\n", filepath.Base(f), key)
		}
	}

	fmt.Printf("\n=== DIVERSITY CHECK ===\n")
	fmt.Printf("Unique combinations: %d / %d files\n", len(seen), min(len(files), 15))
	if len(seen) == 1 {
		fmt.Printf("⚠️  WARNING: ALL files produce IDENTICAL parameters!\n")
	} else {
		fmt.Printf("✓ Parameters vary across files\n")
	}
}
