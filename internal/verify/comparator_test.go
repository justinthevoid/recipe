package verify

import (
	"image"
	"image/color"
	"testing"
)

// createSolidImage creates a uniform color image of given size.
func createSolidImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestComparator_Compare(t *testing.T) {
	c := NewComparator()

	t.Run("Identical Images", func(t *testing.T) {
		img1 := createSolidImage(100, 100, color.RGBA{255, 0, 0, 255})
		img2 := createSolidImage(100, 100, color.RGBA{255, 0, 0, 255})

		diff, err := c.Compare(img1, img2)
		if err != nil {
			t.Fatalf("Compare failed: %v", err)
		}
		if diff != 0 {
			t.Errorf("Expected 0 difference for identical images, got %f", diff)
		}
	})

	t.Run("Different Images (Black vs White)", func(t *testing.T) {
		img1 := createSolidImage(100, 100, color.Black)
		img2 := createSolidImage(100, 100, color.White)

		diff, err := c.Compare(img1, img2)
		if err != nil {
			t.Fatalf("Compare failed: %v", err)
		}
		if diff <= 0 {
			t.Errorf("Expected positive difference for different images, got %f", diff)
		}
	})

	t.Run("Resize Verification", func(t *testing.T) {
		// Test that internal resizing logic works by passing huge images
		// We can't easily peek inside without exporting, but we can ensure it runs
		img1 := createSolidImage(2000, 2000, color.Black)
		img2 := createSolidImage(2000, 2000, color.Black)

		diff, err := c.Compare(img1, img2)
		if err != nil {
			t.Fatalf("Compare failed on large images: %v", err)
		}
		if diff != 0 {
			t.Errorf("Expected 0 difference, got %f", diff)
		}
	})

	t.Run("Similar Images (Minor Noise)", func(t *testing.T) {
		img1 := createSolidImage(100, 100, color.Gray{Y: 128})
		// Create img2 with slightly different color
		img2 := createSolidImage(100, 100, color.Gray{Y: 129}) // Very small difference

		diff, err := c.Compare(img1, img2)
		if err != nil {
			t.Fatalf("Compare failed: %v", err)
		}
		if diff == 0 {
			t.Error("Expected non-zero difference for similar images")
		}
		if diff > 0.1 {
			t.Errorf("Expected small difference, got %f", diff)
		}
	})

	t.Run("Error Handling", func(t *testing.T) {
		img := createSolidImage(100, 100, color.White)
		_, err := c.Compare(nil, img)
		if err == nil {
			t.Error("Expected error when first image is nil")
		}

		_, err = c.Compare(img, nil)
		if err == nil {
			t.Error("Expected error when second image is nil")
		}
	})
}
