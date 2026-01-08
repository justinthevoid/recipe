package nksc

import (
	"encoding/base64"

	"github.com/justin/recipe/internal/formats/np3"
)

// NewFromNP3 converts NP3 metadata into an NKSC structure.
// targetNEF is the filename of the NEF this sidecar belongs to (used in references).
func NewFromNP3(meta *np3.Metadata, targetNEF string) (*NKSC, error) {
	nksc := NewNKSC()

	// Set the About attribute to the target filename (or generic empty)
	// Sidecars often reference the original file.
	nksc.RDF.Description.About = targetNEF

	// Construct the payload string.
	// We wrap the raw NP3 data into the export data field.
	// This ensures the Picture Control is fully preserved.
	// We also format a human-readable summary if needed, but for now
	// we assume the opaque data is what NX Studio needs.
	// Note: In a real implementation we would format this string exactly as Nikon does,
	// which might be a comma-separated list or a hex blob.
	// Base64 is a safe container for binary data in XML.
	payload := base64.StdEncoding.EncodeToString(meta.RawBytes)

	// Also include a readable summary for debugging/inspection (optional, if field existed)

	// Create the Picture Control step
	step := Step{
		ParseType:                  "Resource",
		Name:                       "Picture Control",
		FilterParametersExportData: payload,
	}

	// If Label is present, maybe use it as StepName?
	// Usually StepName is the Type of edit (Picture Control),
	// and the specific preset name (Label) is inside the parameters.

	// Add to NKSC
	nksc.RDF.Description.NineEdits = NineEdits{
		Seq: Seq{
			Steps: []Step{step},
		},
	}

	return &nksc, nil
}
