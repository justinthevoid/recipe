// +build ignore

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/tiff"
)

func main() {
	data, err := os.ReadFile("testdata/dcp/Nikon Z f Camera Standard.dcp")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Convert DNG to TIFF
	if len(data) >= 4 && data[0] == 0x49 && data[1] == 0x49 && data[2] == 0x52 && data[3] == 0x43 {
		data[2] = 0x2A
		data[3] = 0x00
	}

	reader := tiff.NewReadAtReadSeeker(bytes.NewReader(data))
	tiffFile, err := tiff.Parse(reader, tiff.DefaultTagSpace, tiff.DefaultFieldTypeSpace)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	ifds := tiffFile.IFDs()
	fmt.Printf("Number of IFDs: %d\n", len(ifds))

	for i, ifd := range ifds {
		fmt.Printf("\nIFD %d:\n", i)
		fields := ifd.Fields()
		fmt.Printf("  Number of fields: %d\n", len(fields))

		// List all tags
		for j, field := range fields {
			if j < 30 {
				fmt.Printf("  Tag %d (0x%04x): type=%v\n", field.Tag(), field.Tag(), field.Type())
			}
		}

		// Check for tag 50740 specifically
		if ifd.HasField(50740) {
			fmt.Printf("\n  *** Found tag 50740! ***\n")
			field := ifd.GetField(50740)
			if field != nil {
				fmt.Printf("  Tag 50740 type: %v\n", field.Type())
				data := field.Value().Bytes()
				fmt.Printf("  Tag 50740 data length: %d bytes\n", len(data))
				if len(data) > 0 {
					fmt.Printf("  First 100 bytes: %s\n", string(data[:min(100, len(data))]))
				}
			}
		} else {
			fmt.Printf("\n  *** Tag 50740 NOT found ***\n")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
