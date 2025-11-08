// test_temperature.go - Check temperature/tint in recipe
package main

import (
	"fmt"
	"os"

	"github.com/justin/recipe/internal/formats/np3"
)

func main() {
	data, err := os.ReadFile("testdata/FIlmstill's Nostalgic Negative.NP3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	r, err := np3.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Temperature: %v\n", r.Temperature)
	fmt.Printf("Tint: %d\n", r.Tint)
}
