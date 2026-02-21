package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyCmd_Flags(t *testing.T) {
	cmd := newVerifyCmd()

	tests := []struct {
		hasFlag string
	}{
		{"input"},
		{"reference"},
		{"threshold"},
	}

	for _, tt := range tests {
		if cmd.Flag(tt.hasFlag) == nil {
			t.Errorf("verify command missing flag: %s", tt.hasFlag)
		}
	}
}

func createSolidJPEG(path string, width, height int, c color.Color) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, nil)
}

func TestVerifyCmd_Run_FailDiff(t *testing.T) {
	// Setup generic "different" test case
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	refDir := filepath.Join(tmpDir, "ref")
	os.Mkdir(inputDir, 0755)
	os.Mkdir(refDir, 0755)

	// Create different images
	createSolidJPEG(filepath.Join(inputDir, "test.jpg"), 100, 100, color.White)
	createSolidJPEG(filepath.Join(refDir, "test.jpg"), 100, 100, color.Black)

	cmd := newVerifyCmd()
	// Set flags
	cmd.SetArgs([]string{"--input", inputDir, "--reference", refDir})

	// Execute should fail because images differ
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for differing images, got nil")
	}
}

func TestVerifyCmd_Run_PassSame(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	refDir := filepath.Join(tmpDir, "ref")
	os.Mkdir(inputDir, 0755)
	os.Mkdir(refDir, 0755)

	// Create same images
	createSolidJPEG(filepath.Join(inputDir, "test.jpg"), 100, 100, color.Black)
	createSolidJPEG(filepath.Join(refDir, "test.jpg"), 100, 100, color.Black)

	cmd := newVerifyCmd()
	cmd.SetArgs([]string{"--input", inputDir, "--reference", refDir})

	// Execute should pass
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected success for same images, got error: %v", err)
	}
}
