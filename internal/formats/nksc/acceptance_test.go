//go:build acceptance

package nksc

import (
	"os"
	"testing"
)

// TestNXStudioIntegration verifies that the environment is set up correctly for acceptance tests.
// This requires the 'acceptance' build tag.
func TestNXStudioIntegration(t *testing.T) {
	nxPath := os.Getenv("NX_STUDIO_PATH")
	if nxPath == "" {
		// IN CI, this MUST be set. Locally, we skip.
		if os.Getenv("CI") != "" {
			t.Fatal("NX_STUDIO_PATH not set in CI environment")
		}
		t.Skip("NX_STUDIO_PATH not set, skipping integration test")
	}

	// Verify the executable exists
	info, err := os.Stat(nxPath)
	if err != nil {
		t.Fatalf("NX Studio executable not found at %s: %v", nxPath, err)
	}
	if info.IsDir() {
		t.Fatalf("NX Studio path %s is a directory, not a file", nxPath)
	}

	// In a real scenario, we might try to run it with a flag, but it's a GUI app.
	// We'll just verify file permissions/existence as a proxy for "installed".
	t.Logf("NX Studio found at %s", nxPath)
}
