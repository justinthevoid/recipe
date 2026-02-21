package verify

import (
	"fmt"
	"image"
	"math"
)

// Comparator handles image comparison logic.
type Comparator struct {
	targetSize int
}

// NewComparator creates a new image comparator.
func NewComparator() *Comparator {
	return &Comparator{
		targetSize: 256, // Fixed small size as per requirements
	}
}

// Compare compares two images and returns a difference score (RMSE).
// Returns 0.0 for identical images, and ranges up to 1.0 for maximum difference.
// It returns an error if the images cannot be compared (e.g. nil images).
func (c *Comparator) Compare(img1, img2 image.Image) (float64, error) {
	if img1 == nil || img2 == nil {
		return 0, fmt.Errorf("cannot compare nil images")
	}

	// Downsample both images to fixed size
	small1 := c.downsample(img1)
	small2 := c.downsample(img2)

	return c.calculateRMSE(small1, small2), nil
}

func (c *Comparator) downsample(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Create target buffer
	out := image.NewRGBA(image.Rect(0, 0, c.targetSize, c.targetSize))

	// Nearest neighbor scaling
	xRatio := float64(w) / float64(c.targetSize)
	yRatio := float64(h) / float64(c.targetSize)

	for y := 0; y < c.targetSize; y++ {
		srcY := bounds.Min.Y + int(float64(y)*yRatio)
		// Clamp to bounds to be safe due to floating point precision
		if srcY >= bounds.Max.Y {
			srcY = bounds.Max.Y - 1
		}

		for x := 0; x < c.targetSize; x++ {
			srcX := bounds.Min.X + int(float64(x)*xRatio)
			if srcX >= bounds.Max.X {
				srcX = bounds.Max.X - 1
			}

			// Get color directly
			out.Set(x, y, img.At(srcX, srcY))
		}
	}
	return out
}

func (c *Comparator) calculateRMSE(img1, img2 *image.RGBA) float64 {
	bounds := img1.Bounds() // Should be same size 256x256
	w, h := bounds.Dx(), bounds.Dy()

	var sumSqDiff float64

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Convert to 0-1 float range logic
			// RGBA() returns 0-0xffff (16-bit)
			const maxVal = 65535.0

			dr := (float64(r1) - float64(r2)) / maxVal
			dg := (float64(g1) - float64(g2)) / maxVal
			db := (float64(b1) - float64(b2)) / maxVal

			sumSqDiff += dr*dr + dg*dg + db*db
		}
	}

	meanSq := sumSqDiff / float64(w*h*3)
	return math.Sqrt(meanSq)
}
