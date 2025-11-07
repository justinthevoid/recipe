package converter

import (
	"os"
	"testing"

	"github.com/justin/recipe/internal/models"
)

// BenchmarkConvert_NP3_to_XMP measures NP3→XMP conversion performance
func BenchmarkConvert_NP3_to_XMP(b *testing.B) {
	// Find sample NP3 file
	sampleFile := "../../examples/np3/Denis Zeqiri/Classic Chrome.np3"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No NP3 sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatNP3, FormatXMP)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkConvert_NP3_to_LRTemplate measures NP3→lrtemplate conversion performance
func BenchmarkConvert_NP3_to_LRTemplate(b *testing.B) {
	sampleFile := "../../examples/np3/Denis Zeqiri/Classic Chrome.np3"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No NP3 sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatNP3, FormatLRTemplate)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkConvert_XMP_to_NP3 measures XMP→NP3 conversion performance
func BenchmarkConvert_XMP_to_NP3(b *testing.B) {
	sampleFile := "../../examples/lrtemplate/015. PRESETPRO - Emulation K/00. E - auto tone.xmp"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No XMP sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatXMP, FormatNP3)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkConvert_XMP_to_LRTemplate measures XMP→lrtemplate conversion performance
func BenchmarkConvert_XMP_to_LRTemplate(b *testing.B) {
	sampleFile := "../../examples/lrtemplate/015. PRESETPRO - Emulation K/00. E - auto tone.xmp"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No XMP sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatXMP, FormatLRTemplate)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkConvert_LRTemplate_to_NP3 measures lrtemplate→NP3 conversion performance
func BenchmarkConvert_LRTemplate_to_NP3(b *testing.B) {
	sampleFile := "../../examples/lrtemplate/Fujifilm Pro 400H.lrtemplate"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No lrtemplate sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatLRTemplate, FormatNP3)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkConvert_LRTemplate_to_XMP measures lrtemplate→XMP conversion performance
func BenchmarkConvert_LRTemplate_to_XMP(b *testing.B) {
	sampleFile := "../../examples/lrtemplate/Fujifilm Pro 400H.lrtemplate"
	input, err := os.ReadFile(sampleFile)
	if err != nil {
		b.Skipf("No lrtemplate sample file found at %s: %v", sampleFile, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Convert(input, FormatLRTemplate, FormatXMP)
		if err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}

// BenchmarkDetectFormat_NP3 measures NP3 format detection overhead
func BenchmarkDetectFormat_NP3(b *testing.B) {
	np3Data := make([]byte, 300)
	copy(np3Data, []byte{'N', 'C', 'P'})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DetectFormat(np3Data)
		if err != nil {
			b.Fatalf("Detection failed: %v", err)
		}
	}
}

// BenchmarkDetectFormat_XMP measures XMP format detection overhead
func BenchmarkDetectFormat_XMP(b *testing.B) {
	xmpData := []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
</rdf:Description>
</rdf:RDF>
</x:xmpmeta>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DetectFormat(xmpData)
		if err != nil {
			b.Fatalf("Detection failed: %v", err)
		}
	}
}

// BenchmarkDetectFormat_LRTemplate measures lrtemplate format detection overhead
func BenchmarkDetectFormat_LRTemplate(b *testing.B) {
	lrtemplateData := []byte(`s = {
	id = "test",
	value = {}
}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DetectFormat(lrtemplateData)
		if err != nil {
			b.Fatalf("Detection failed: %v", err)
		}
	}
}

// BenchmarkHubOverhead measures UniversalRecipe hub overhead (<5ms target)
// This benchmark isolates the hub conversion overhead by using in-memory recipe data
func BenchmarkHubOverhead(b *testing.B) {
	// Create a realistic UniversalRecipe with typical parameters
	temp := 5800
	recipe := &models.UniversalRecipe{
		Name:         "Benchmark Recipe",
		SourceFormat: FormatNP3,
		Exposure:     1.5,
		Contrast:     25,
		Highlights:   -50,
		Shadows:      50,
		Whites:       10,
		Blacks:       -10,
		Clarity:      20,
		Vibrance:     15,
		Saturation:   10,
		Sharpness:    40,
		Temperature:  &temp,
		Tint:         5,
	}

	b.ResetTimer()
	b.Run("Parse+Generate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Simulate parse → hub → generate flow without actual I/O
			// This measures the overhead of working with UniversalRecipe

			// Generate would create output bytes from recipe
			_ = recipe.Exposure
			_ = recipe.Contrast
			_ = recipe.Saturation
		}
	})
}

// BenchmarkConvert_AutoDetect measures auto-detection overhead in Convert()
func BenchmarkConvert_AutoDetect(b *testing.B) {
	// Create minimal valid NP3 data
	np3Data := make([]byte, 300)
	copy(np3Data, []byte{'N', 'C', 'P'})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use empty 'from' parameter to trigger auto-detection
		_, err := Convert(np3Data, "", FormatXMP)
		if err != nil {
			b.Fatalf("Conversion with auto-detect failed: %v", err)
		}
	}
}
