package main

import (
	"os"
)

func main() {
	if err := Execute(); err != nil {
		// Cobra already prints error to stderr
		os.Exit(1)
	}
	os.Exit(0)
}
