// Package verify provides verification tools for checking exported images against reference "golden" images.
//
// It includes an ImageComparator for pixel-by-pixel comparison with configurable thresholds,
// supporting JPEG and TIFF formats (via standard library interfaces).
// It implements robust comparison logic including automatic downsampling to handle
// compression artifacts and improve performance.
package verify
