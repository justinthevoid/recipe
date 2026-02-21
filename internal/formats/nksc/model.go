package nksc

import "encoding/xml"

// Namespaces
const (
	NsXmpMeta  = "adobe:ns:meta/"
	NsRDF      = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	NsSDC      = "http://ns.nikon.com/sdc/1.0/"
	NsAsteroid = "http://ns.nikon.com/asteroid/1.0/"
	NsAstType  = "http://ns.nikon.com/asteroid/Types/1.0/"
	NsNine     = "http://ns.nikon.com/nine/1.0/"
)

// NKSC represents the root XMP wrapper for Nikon Key Store Container files.
type NKSC struct {
	XMLName xml.Name `xml:"x:xmpmeta"`
	X       string   `xml:"xmlns:x,attr"` // xmlns:x="adobe:ns:meta/"
	XmpTk   string   `xml:"x:xmptk,attr"` // x:xmptk="XMP Core 5.5.0"
	RDF     RDF      `xml:"rdf:RDF"`
}

// RDF container
type RDF struct {
	RDFNS       string      `xml:"xmlns:rdf,attr"` // xmlns:rdf="..."
	Description Description `xml:"rdf:Description"`
}

// Description holds the Nikon specific namespaces
type Description struct {
	About     string `xml:"rdf:about,attr"`
	SdcNS     string `xml:"xmlns:sdc,attr"`
	AstNS     string `xml:"xmlns:ast,attr"`
	AstTypeNS string `xml:"xmlns:astype,attr"`
	NineNS    string `xml:"xmlns:nine,attr"`

	// SDC fields
	SdcAbout      string `xml:"sdc:about"`
	SdcVersion    string `xml:"sdc:version"`
	SdcAppVersion string `xml:"sdc:appversion"`
	SdcAppName    string `xml:"sdc:appname"`

	// Asteroid fields
	AstAbout        string          `xml:"ast:about"`
	AstVersion      string          `xml:"ast:version"`
	AstXMLPackets   *ParsedResource `xml:"ast:XMLPackets,omitempty"`
	AstGPSVersionID *ParsedResource `xml:"ast:GPSVersionID,omitempty"`

	// Nine fields
	NineAbout   string     `xml:"nine:about"`
	NineVersion string     `xml:"nine:version"`
	NineEdits   *NineEdits `xml:"nine:NineEdits"` // Pointer to allow nil if needed
}

// ParsedResource generic placeholder for rdf:parseType="Resource" fields
type ParsedResource struct {
	ParseType string `xml:"rdf:parseType,attr"`
}

// NineEdits represents the editing steps in the NKSC file.
type NineEdits struct {
	Seq Seq `xml:"rdf:Seq"`
}

// Seq represents the RDF Sequence container
type Seq struct {
	Steps []Step `xml:"rdf:li"`
}

// Step represents a single editing step (Filter) in the NineEdits list.
type Step struct {
	ParseType string `xml:"rdf:parseType,attr,omitempty"` // usually "Resource"

	Name    string `xml:"nine:StepName"`
	Version string `xml:"nine:StepVersion,omitempty"`

	// The payload is often stored in these opaque fields.
	// We map our NP3 data into these.
	FilterParametersExportData string `xml:"nine:FilterParametersExportData,omitempty"`
}

// NewNKSC creates a basic initialized NKSC struct structure
func NewNKSC() NKSC {
	return NKSC{
		X:     NsXmpMeta,
		XmpTk: "XMP Core 5.5.0",
		RDF: RDF{
			RDFNS: NsRDF,
			Description: Description{
				About:         "core-sidear-tags/1.0", // Matches observed reference, likely typo in Nikon SW
				SdcNS:         NsSDC,
				AstNS:         NsAsteroid,
				AstTypeNS:     NsAstType,
				NineNS:        NsNine,
				SdcAbout:      "nikon sidecar/1.0",
				SdcVersion:    "1.1.0",
				SdcAppVersion: "1.10 W",
				SdcAppName:    "NX Studio",
				AstAbout:      "core-asteroid-tags",
				AstVersion:    "11.0.0.3000",
				NineAbout:     "nine-tags",
				NineVersion:   "2.0.0",
			},
		},
	}
}
