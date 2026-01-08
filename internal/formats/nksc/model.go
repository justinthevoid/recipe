package nksc

import "encoding/xml"

// Namespaces
const (
	NsXmpMeta = "adobe:ns:meta/"
	NsRDF     = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	NsAsteria = "http://ns.nikon.com/nksc/1.0/asteria/"
	NsNine    = "http://ns.nikon.com/nksc/1.0/nine/"
)

// NKSC represents the root XMP wrapper for Nikon Key Store Container files.
type NKSC struct {
	XMLName xml.Name `xml:"x:xmpmeta"`
	X       string   `xml:"xmlns:x,attr"` // xmlns:x="adobe:ns:meta/"
	RDF     RDF      `xml:"rdf:RDF"`
}

// RDF container
type RDF struct {
	RDFNS       string      `xml:"xmlns:rdf,attr"` // xmlns:rdf="..."
	Description Description `xml:"rdf:Description"`
}

// Description holds the Nikon specific namespaces
type Description struct {
	About     string    `xml:"rdf:about,attr"`
	AstNS     string    `xml:"xmlns:ast,attr"`
	NineNS    string    `xml:"xmlns:nine,attr"`
	NineEdits NineEdits `xml:"nine:NineEdits"`
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
		X: NsXmpMeta,
		RDF: RDF{
			RDFNS: NsRDF,
			Description: Description{
				About:  "",
				AstNS:  NsAsteria,
				NineNS: NsNine,
			},
		},
	}
}
