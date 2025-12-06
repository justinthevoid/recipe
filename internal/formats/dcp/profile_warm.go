package dcp

// generateColorMatrix2Warm generates a warmer ColorMatrix2 for D65 illuminant.
//
// This is a custom variant of generateColorMatrix2() with modified coefficients
// to match Nikon's warmer color rendering as documented in FINAL_CONCLUSIONS.md.
//
// Key changes from Adobe's original matrix:
//   Blue→Red coefficient: -0.0977 → +0.08 (warmth boost: +0.1777)
//   Red diagonal: 1.1607 → 1.25 (more red passthrough)
//   Blue diagonal: 0.7616 → 0.85 (higher blue sensitivity)
//   Green→Red: -0.4491 → -0.35 (less green correction)
//
// Warm ColorMatrix2 (D65 illuminant):
//   1.25   -0.35    0.08
//  -0.40    1.20    0.15
//  -0.02    0.10    0.85
//
// Returns array of 9 SRational values ready for tag 50722 (ColorMatrix2).
func generateColorMatrix2Warm() []SRational {
	return []SRational{
		// Row 1: [1.25, -0.35, 0.08] ← WARMER (positive blue→red)
		{Numerator: 12500, Denominator: 10000}, {Numerator: -3500, Denominator: 10000}, {Numerator: 800, Denominator: 10000},
		// Row 2: [-0.40, 1.20, 0.15]
		{Numerator: -4000, Denominator: 10000}, {Numerator: 12000, Denominator: 10000}, {Numerator: 1500, Denominator: 10000},
		// Row 3: [-0.02, 0.10, 0.85]
		{Numerator: -200, Denominator: 10000}, {Numerator: 1000, Denominator: 10000}, {Numerator: 8500, Denominator: 10000},
	}
}
