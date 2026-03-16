package main

import (
	"os"
	"testing"
)

// BenchmarkDetectFormatExtension benchmarks extension-based detection
// Target: <1μs (should be nanoseconds)
func BenchmarkDetectFormatExtension(b *testing.B) {
	testPath := "portrait.xmp"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectFormat(testPath)
	}
}

// BenchmarkDetectFormatExtensionNP3 benchmarks NP3 extension detection
func BenchmarkDetectFormatExtensionNP3(b *testing.B) {
	testPath := "sample.np3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectFormat(testPath)
	}
}

// BenchmarkDetectFormatFromBytes benchmarks content-based detection
// Target: <5ms for typical files
func BenchmarkDetectFormatFromBytes(b *testing.B) {
	// Load a real XMP file if available
	data, err := os.ReadFile("testdata/xmp/AFGA APX 100.xmp")
	if err != nil {
		// Use synthetic data if real file not available
		data = []byte(`<?xml version="1.0"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
    <rdf:RDF xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/">
    </rdf:RDF>
</x:xmpmeta>`)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectFormatFromBytes(data)
	}
}

// BenchmarkDetectFormatFromBytesNP3 benchmarks NP3 content detection
func BenchmarkDetectFormatFromBytesNP3(b *testing.B) {
	// Load a real NP3 file if available
	data, err := os.ReadFile("testdata/np3/Classic Chrome.np3")
	if err != nil {
		// Use synthetic NP3 data if real file not available
		data = make([]byte, 300)
		copy(data, []byte{'N', 'C', 'P'})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectFormatFromBytes(data)
	}
}

// BenchmarkValidateFormat benchmarks format validation
func BenchmarkValidateFormat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateFormat("xmp")
	}
}
