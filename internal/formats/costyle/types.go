// Package costyle provides parsing and generation of Capture One .costyle preset files.
package costyle

import "encoding/xml"

// CaptureOneStyle represents the root structure of a .costyle XML file.
// Capture One style files use Adobe XMP-style XML with RDF/Description elements.
type CaptureOneStyle struct {
	XMLName xml.Name `xml:"xmpmeta"`
	RDF     RDF      `xml:"RDF"`
}

// RDF wraps the Description element containing adjustment parameters.
type RDF struct {
	Description Description `xml:"Description"`
}

// Description contains all adjustment parameters from the .costyle file.
// Parameter ranges and mappings:
//   - Exposure: -2.0 to +2.0 (direct map to UniversalRecipe.Exposure)
//   - Contrast: -100 to +100 (direct map to UniversalRecipe.Contrast)
//   - Saturation: -100 to +100 (direct map to UniversalRecipe.Saturation)
//   - Temperature: -100 to +100 (map to UniversalRecipe.Temperature)
//   - Tint: -100 to +100 (direct map to UniversalRecipe.Tint)
//   - Clarity: -100 to +100 (direct map to UniversalRecipe.Clarity)
//   - Color Balance: Hue/Saturation per tonal range (shadows, midtones, highlights)
type Description struct {
	// Core adjustments
	Exposure    float64 `xml:"Exposure,omitempty"`    // -2.0 to +2.0
	Contrast    int     `xml:"Contrast,omitempty"`    // -100 to +100
	Saturation  int     `xml:"Saturation,omitempty"`  // -100 to +100
	Temperature int     `xml:"Temperature,omitempty"` // -100 to +100
	Tint        int     `xml:"Tint,omitempty"`        // -100 to +100
	Clarity     int     `xml:"Clarity,omitempty"`     // -100 to +100

	// Color balance adjustments (per tonal range)
	ShadowsHue           int `xml:"ShadowsHue,omitempty"`
	ShadowsSaturation    int `xml:"ShadowsSaturation,omitempty"`
	MidtonesHue          int `xml:"MidtonesHue,omitempty"`
	MidtonesSaturation   int `xml:"MidtonesSaturation,omitempty"`
	HighlightsHue        int `xml:"HighlightsHue,omitempty"`
	HighlightsSaturation int `xml:"HighlightsSaturation,omitempty"`

	// Metadata fields
	Name        string `xml:"Name,omitempty"`
	Author      string `xml:"Author,omitempty"`
	Description string `xml:"Desc,omitempty"`
}
