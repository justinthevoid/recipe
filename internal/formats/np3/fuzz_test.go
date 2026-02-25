package np3

import "testing"

// FuzzParse ensures the parser is resilient and will not panic even if provided
// with entirely random, corrupted, or fuzzed byte inputs.
func FuzzParse(f *testing.F) {
	// Provide seed corpora
	f.Add([]byte{})
	f.Add([]byte("NCP\x00\x00\x00\x00"))
	f.Add([]byte("NCP\x00\x00\x00\x0012345678901234567890"))

	f.Fuzz(func(t *testing.T, data []byte) {
		// The function must not panic under any circumstances
		_, _ = Parse(data)
	})
}
