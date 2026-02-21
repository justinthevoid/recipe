package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// FixtureFile represents a file to be downloaded (NEF or Sidecar).
type FixtureFile struct {
	Filename string `yaml:"filename"`
	URL      string `yaml:"url"`
	SHA256   string `yaml:"sha256"`
}

// FixtureEntry represents a singe test fixture file and its metadata.
type FixtureEntry struct {
	ID              string        `yaml:"id"`
	Filename        string        `yaml:"filename"` // Primary file (NEF)
	URL             string        `yaml:"url"`      // Primary URL
	SHA256          string        `yaml:"sha256"`   // Primary SHA256
	CameraModel     string        `yaml:"camera_model"`
	Lens            string        `yaml:"lens"`
	NXStudioVersion string        `yaml:"nx_studio_version"`
	NP3Variant      string        `yaml:"np3_variant"`
	Sidecars        []FixtureFile `yaml:"sidecars,omitempty"`
}

// FixtureManifest represents the collection of test fixtures.
type FixtureManifest struct {
	Fixtures []FixtureEntry `yaml:"fixtures"`
}

// EnsureFixtures checks if fixtures are available.
// If fixtures are missing, it invokes t.Skip().
func EnsureFixtures(t *testing.T) {
	t.Helper()

	root, err := findProjectRoot()
	if err != nil {
		// Can't find root, assume we can't check fixtures.
		// In a real scenario this might be a fail, but for now log and return.
		t.Logf("Could not find project root: %v", err)
		return
	}

	manifestPath := filepath.Join(root, "testdata", "nx-fixtures", "manifest.yaml")
	fixturesDir := filepath.Join(root, "testdata", "nx-fixtures")

	data, err := os.ReadFile(manifestPath)
	if os.IsNotExist(err) {
		t.Skip("Manifest not found at " + manifestPath)
		return
	}
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}

	var manifest FixtureManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("Failed to parse manifest: %v", err)
	}

	if len(manifest.Fixtures) == 0 {
		return
	}

	missing := 0
	for _, f := range manifest.Fixtures {
		path := filepath.Join(fixturesDir, f.Filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing++
		}
	}

	if missing > 0 {
		t.Logf("%d fixtures missing. Run 'make setup-fixtures' to download.", missing)
		t.Skip("Skipping due to missing fixtures")
	}
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
