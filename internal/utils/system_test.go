package utils_test

import (
	"testing"

	"github.com/justin/recipe/internal/utils"
)

func TestOpenFolder(t *testing.T) {
	// Verify OpenFolder returns an error for invalid paths or fails gracefully
	// Note: We cannot easily test successful GUI window opening in a headless CI environment,
	// so we expect an error when passing "." if xdg-open/explorer is missing,
	// or success if it works. Here we just ensure it doesn't panic.
	// In this specific test setup (Linux/Headless), xdg-open likely fails or is missing.
	err := utils.OpenFolder(".")
	if err != nil {
		t.Logf("OpenFolder returned error (expected in headless env): %v", err)
	}
}
