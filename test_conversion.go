package main

import (
	"fmt"
	"os"
	"github.com/justin/recipe/internal/converter"
)

func main() {
	data, err := os.ReadFile("testdata/xmp/Fuji Industrial 400.xmp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		os.Exit(1)
	}

	np3Data, err := converter.Convert(data, "xmp", "np3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Convert error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converted successfully: %d bytes\n", len(np3Data))
	
	// Write to output for inspection
	os.WriteFile("test_debug.np3", np3Data, 0644)
}
