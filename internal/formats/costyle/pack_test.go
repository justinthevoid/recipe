package costyle

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/justin/recipe/internal/models"
)

func TestUnpack_ValidBundle(t *testing.T) {
	// Create a .costylepack with 3 valid .costyle files
	recipes := createTestRecipes(3)
	filenames := []string{"Portrait.costyle", "Landscape.costyle", "Product.costyle"}

	packData, err := Pack(recipes, filenames)
	if err != nil {
		t.Fatalf("Pack() failed: %v", err)
	}

	// Unpack the bundle
	extractedRecipes, err := Unpack(packData)
	if err != nil {
		t.Fatalf("Unpack() failed: %v", err)
	}

	// Verify all 3 recipes extracted
	if len(extractedRecipes) != 3 {
		t.Errorf("Expected 3 recipes, got %d", len(extractedRecipes))
	}

	// Verify filenames preserved in metadata
	expectedFilenames := []string{"Portrait.costyle", "Landscape.costyle", "Product.costyle"}
	for i, recipe := range extractedRecipes {
		if recipe.Metadata == nil {
			t.Errorf("Recipe %d has nil metadata", i)
			continue
		}
		filename, ok := recipe.Metadata["original_filename"].(string)
		if !ok {
			t.Errorf("Recipe %d missing original_filename in metadata", i)
			continue
		}
		if filename != expectedFilenames[i] {
			t.Errorf("Recipe %d: expected filename %s, got %s", i, expectedFilenames[i], filename)
		}
	}

	// Verify no errors returned
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

func TestUnpack_EmptyBundle(t *testing.T) {
	// Create valid empty ZIP with PK magic bytes but no files
	// We'll create a minimal valid ZIP structure manually
	emptyZIP := []byte{
		0x50, 0x4B, 0x05, 0x06, // End of central directory signature
		0x00, 0x00, 0x00, 0x00, // Number of this disk
		0x00, 0x00, 0x00, 0x00, // Disk where central directory starts
		0x00, 0x00, 0x00, 0x00, // Number of central directory records
		0x00, 0x00, 0x00, 0x00, // Size of central directory
		0x00, 0x00, 0x00, 0x00, // Offset of start of central directory
		0x00, 0x00,             // ZIP file comment length
	}

	// Unpack empty bundle
	recipes, err := Unpack(emptyZIP)

	// Verify empty slice returned
	if len(recipes) != 0 {
		t.Errorf("Expected empty slice, got %d recipes", len(recipes))
	}

	// Verify error message (empty ZIP should return error about no files)
	if err == nil {
		t.Error("Expected error for empty bundle, got nil")
	} else if !contains(err.Error(), "empty .costylepack") && !contains(err.Error(), "no files found") && !contains(err.Error(), "not a valid ZIP") {
		t.Errorf("Expected error about empty bundle or invalid ZIP, got: %v", err)
	}
}

func TestUnpack_CorruptZIP(t *testing.T) {
	// Provide truncated/corrupt ZIP bytes
	corruptData := []byte{0x50, 0x4B, 0x03, 0x04, 0x00, 0x00} // PK header but incomplete

	// Unpack corrupt ZIP
	recipes, err := Unpack(corruptData)

	// Verify error returned
	if err == nil {
		t.Error("Expected error for corrupt ZIP, got nil")
	}

	// Verify error contains "corrupt" or "failed"
	if err != nil && !contains(err.Error(), "failed to read .costylepack ZIP") {
		t.Errorf("Expected error about corrupt ZIP, got: %v", err)
	}

	// Verify no recipes returned
	if recipes != nil {
		t.Errorf("Expected nil recipes for corrupt ZIP, got %d recipes", len(recipes))
	}
}

func TestUnpack_InvalidMagicBytes(t *testing.T) {
	// Provide non-ZIP data
	invalidData := []byte("This is not a ZIP file")

	// Unpack invalid data
	recipes, err := Unpack(invalidData)

	// Verify error returned
	if err == nil {
		t.Error("Expected error for invalid magic bytes, got nil")
	}

	// Verify error contains "not a valid ZIP"
	if err != nil && !contains(err.Error(), "not a valid ZIP file") {
		t.Errorf("Expected error about invalid ZIP, got: %v", err)
	}

	// Verify no recipes returned
	if recipes != nil {
		t.Errorf("Expected nil recipes for invalid ZIP, got %d recipes", len(recipes))
	}
}

func TestUnpack_NonCostyleFiles(t *testing.T) {
	// Create .costylepack with mixed files: 2 .costyle, 2 .txt, 1 .json
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add 2 valid .costyle files
	recipes := createTestRecipes(2)
	for i, recipe := range recipes {
		xmlData, _ := Generate(recipe)
		filename := []string{"Style1.costyle", "Style2.costyle"}[i]
		fileWriter, _ := zipWriter.Create(filename)
		fileWriter.Write(xmlData)
	}

	// Add 2 .txt files
	txtWriter, _ := zipWriter.Create("readme.txt")
	txtWriter.Write([]byte("This is a readme"))

	txtWriter2, _ := zipWriter.Create("info.txt")
	txtWriter2.Write([]byte("More info"))

	// Add 1 .json file
	jsonWriter, _ := zipWriter.Create("metadata.json")
	jsonWriter.Write([]byte(`{"name": "test"}`))

	zipWriter.Close()

	// Unpack the bundle
	extractedRecipes, err := Unpack(buf.Bytes())

	// Verify only 2 recipes extracted (non-.costyle files skipped)
	if len(extractedRecipes) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(extractedRecipes))
	}

	// Verify no error (non-.costyle files skipped gracefully)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
}

func TestUnpack_PartialFailure(t *testing.T) {
	// Create .costylepack with 2 valid .costyle and 1 malformed .costyle
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add 2 valid .costyle files
	recipes := createTestRecipes(2)
	for i, recipe := range recipes {
		xmlData, _ := Generate(recipe)
		filename := []string{"Valid1.costyle", "Valid2.costyle"}[i]
		fileWriter, _ := zipWriter.Create(filename)
		fileWriter.Write(xmlData)
	}

	// Add 1 malformed .costyle
	malformedWriter, _ := zipWriter.Create("Malformed.costyle")
	malformedWriter.Write([]byte("<invalid>XML</not-closed>"))

	zipWriter.Close()

	// Unpack the bundle
	extractedRecipes, err := Unpack(buf.Bytes())

	// Verify 2 valid recipes extracted (malformed skipped)
	if len(extractedRecipes) != 2 {
		t.Errorf("Expected 2 recipes, got %d", len(extractedRecipes))
	}

	// Verify no fatal error (partial success allowed)
	if err != nil {
		t.Errorf("Expected nil error for partial success, got: %v", err)
	}
}

func TestPack_ValidRecipes(t *testing.T) {
	// Pack 3 recipes with explicit filenames
	recipes := createTestRecipes(3)
	filenames := []string{"Portrait.costyle", "Landscape.costyle", "Product.costyle"}

	packData, err := Pack(recipes, filenames)
	if err != nil {
		t.Fatalf("Pack() failed: %v", err)
	}

	// Verify ZIP created successfully
	if len(packData) == 0 {
		t.Error("Pack() returned empty data")
	}

	// Verify ZIP contains exactly 3 .costyle files
	zipReader, err := zip.NewReader(bytes.NewReader(packData), int64(len(packData)))
	if err != nil {
		t.Fatalf("Failed to read packed ZIP: %v", err)
	}

	if len(zipReader.File) != 3 {
		t.Errorf("Expected 3 files in ZIP, got %d", len(zipReader.File))
	}

	// Verify filenames match input
	for i, file := range zipReader.File {
		if file.Name != filenames[i] {
			t.Errorf("File %d: expected %s, got %s", i, filenames[i], file.Name)
		}
	}

	// Verify each .costyle XML is valid (can be parsed back)
	for i, file := range zipReader.File {
		rc, _ := file.Open()
		xmlData, _ := readAll(rc)
		rc.Close()

		_, err := Parse(xmlData)
		if err != nil {
			t.Errorf("File %d (%s) failed to parse: %v", i, file.Name, err)
		}
	}
}

func TestPack_AutoFilenames(t *testing.T) {
	// Pack 3 recipes without filenames (empty slice)
	recipes := createTestRecipes(3)
	filenames := []string{} // Empty = auto-generate

	packData, err := Pack(recipes, filenames)
	if err != nil {
		t.Fatalf("Pack() failed: %v", err)
	}

	// Verify filenames auto-generated: Style1.costyle, Style2.costyle, Style3.costyle
	zipReader, err := zip.NewReader(bytes.NewReader(packData), int64(len(packData)))
	if err != nil {
		t.Fatalf("Failed to read packed ZIP: %v", err)
	}

	expectedFilenames := []string{"Style1.costyle", "Style2.costyle", "Style3.costyle"}
	for i, file := range zipReader.File {
		if file.Name != expectedFilenames[i] {
			t.Errorf("File %d: expected %s, got %s", i, expectedFilenames[i], file.Name)
		}
	}
}

func TestPack_DuplicateFilenames(t *testing.T) {
	// Pack 3 recipes with duplicate filenames
	recipes := createTestRecipes(3)
	filenames := []string{"Preset.costyle", "Preset.costyle", "Other.costyle"}

	packData, err := Pack(recipes, filenames)
	if err != nil {
		t.Fatalf("Pack() failed: %v", err)
	}

	// Verify output filenames: Preset.costyle, Preset_1.costyle, Other.costyle
	zipReader, err := zip.NewReader(bytes.NewReader(packData), int64(len(packData)))
	if err != nil {
		t.Fatalf("Failed to read packed ZIP: %v", err)
	}

	expectedFilenames := []string{"Preset.costyle", "Preset_1.costyle", "Other.costyle"}
	for i, file := range zipReader.File {
		if file.Name != expectedFilenames[i] {
			t.Errorf("File %d: expected %s, got %s", i, expectedFilenames[i], file.Name)
		}
	}

	// Verify no filename collisions
	filenameSet := make(map[string]bool)
	for _, file := range zipReader.File {
		if filenameSet[file.Name] {
			t.Errorf("Duplicate filename in ZIP: %s", file.Name)
		}
		filenameSet[file.Name] = true
	}
}

func TestPack_EmptyRecipes(t *testing.T) {
	// Pack empty recipe slice
	recipes := []*models.UniversalRecipe{}
	filenames := []string{}

	_, err := Pack(recipes, filenames)

	// Verify error returned
	if err == nil {
		t.Error("Expected error for empty recipes, got nil")
	}

	// Verify error message
	if err != nil && !contains(err.Error(), "at least one recipe required") {
		t.Errorf("Expected error about empty recipes, got: %v", err)
	}
}

func TestPack_MismatchedFilenames(t *testing.T) {
	// Pack 3 recipes with 2 filenames (mismatch)
	recipes := createTestRecipes(3)
	filenames := []string{"File1.costyle", "File2.costyle"} // Only 2 filenames for 3 recipes

	_, err := Pack(recipes, filenames)

	// Verify error returned
	if err == nil {
		t.Error("Expected error for mismatched filenames, got nil")
	}

	// Verify error message
	if err != nil && !contains(err.Error(), "filenames count") {
		t.Errorf("Expected error about filename count mismatch, got: %v", err)
	}
}

func TestRoundTrip_Costylepack(t *testing.T) {
	// Pack 5 recipes → .costylepack bytes
	originalRecipes := createTestRecipes(5)
	filenames := []string{"R1.costyle", "R2.costyle", "R3.costyle", "R4.costyle", "R5.costyle"}

	packData1, err := Pack(originalRecipes, filenames)
	if err != nil {
		t.Fatalf("First Pack() failed: %v", err)
	}

	// Unpack bytes → recipes
	unpackedRecipes, err := Unpack(packData1)
	if err != nil {
		t.Fatalf("Unpack() failed: %v", err)
	}

	// Pack recipes again → .costylepack bytes
	packData2, err := Pack(unpackedRecipes, filenames)
	if err != nil {
		t.Fatalf("Second Pack() failed: %v", err)
	}

	// Unpack second time → recipes2
	finalRecipes, err := Unpack(packData2)
	if err != nil {
		t.Fatalf("Second Unpack() failed: %v", err)
	}

	// Compare original recipes with finalRecipes (verify 95%+ field match)
	if len(finalRecipes) != len(originalRecipes) {
		t.Errorf("Recipe count mismatch: original %d, final %d", len(originalRecipes), len(finalRecipes))
	}

	for i := 0; i < len(originalRecipes) && i < len(finalRecipes); i++ {
		orig := originalRecipes[i]
		final := finalRecipes[i]

		// Compare key fields (allow ±1 tolerance for integer rounding)
		if !floatEqual(orig.Exposure, final.Exposure, 0.01) {
			t.Errorf("Recipe %d: Exposure mismatch: %f vs %f", i, orig.Exposure, final.Exposure)
		}
		if !intEqual(orig.Contrast, final.Contrast, 1) {
			t.Errorf("Recipe %d: Contrast mismatch: %d vs %d", i, orig.Contrast, final.Contrast)
		}
		if !intEqual(orig.Saturation, final.Saturation, 1) {
			t.Errorf("Recipe %d: Saturation mismatch: %d vs %d", i, orig.Saturation, final.Saturation)
		}
		if !intEqual(orig.Clarity, final.Clarity, 1) {
			t.Errorf("Recipe %d: Clarity mismatch: %d vs %d", i, orig.Clarity, final.Clarity)
		}
	}
}

func TestUnpack_RealSampleFile(t *testing.T) {
	// Create a .costylepack from real sample files
	testdataPath := "testdata/costyle"
	samples := []string{"sample1-portrait.costyle", "sample2-minimal.costyle", "sample3-landscape.costyle"}

	// Read real .costyle files
	var recipes []*models.UniversalRecipe
	for _, filename := range samples {
		data, err := os.ReadFile(filepath.Join(testdataPath, filename))
		if err != nil {
			t.Skipf("Skipping test: sample file not found: %s", filename)
			return
		}
		recipe, err := Parse(data)
		if err != nil {
			t.Fatalf("Failed to parse sample %s: %v", filename, err)
		}
		recipes = append(recipes, recipe)
	}

	// Pack into .costylepack
	packData, err := Pack(recipes, samples)
	if err != nil {
		t.Fatalf("Pack() failed with real samples: %v", err)
	}

	// Unpack and verify
	extractedRecipes, err := Unpack(packData)
	if err != nil {
		t.Fatalf("Unpack() failed with real samples: %v", err)
	}

	if len(extractedRecipes) != len(samples) {
		t.Errorf("Expected %d recipes from real samples, got %d", len(samples), len(extractedRecipes))
	}
}

// Helper functions

func createTestRecipes(count int) []*models.UniversalRecipe {
	recipes := make([]*models.UniversalRecipe, count)
	for i := 0; i < count; i++ {
		recipes[i] = &models.UniversalRecipe{
			SourceFormat: "costyle",
			Name:         "Test Recipe",
			Exposure:     float64(i) * 0.5,
			Contrast:     i * 10,
			Saturation:   i * 5,
			Clarity:      i * 3,
			Metadata:     make(map[string]interface{}),
		}
	}
	return recipes
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func floatEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

func intEqual(a, b, tolerance int) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

func readAll(rc io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(rc)
	return buf.Bytes(), err
}

// Benchmarks

func BenchmarkPack50Files(b *testing.B) {
	recipes := createTestRecipes(50)
	filenames := make([]string, 50)
	for i := 0; i < 50; i++ {
		filenames[i] = "Style" + string(rune('0'+i)) + ".costyle"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Pack(recipes, filenames)
		if err != nil {
			b.Fatalf("Pack() failed: %v", err)
		}
	}
}

func BenchmarkUnpack50Files(b *testing.B) {
	// Create a 50-file .costylepack once
	recipes := createTestRecipes(50)
	filenames := make([]string, 50)
	for i := 0; i < 50; i++ {
		filenames[i] = "Style" + string(rune('0'+i)) + ".costyle"
	}
	packData, err := Pack(recipes, filenames)
	if err != nil {
		b.Fatalf("Pack() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Unpack(packData)
		if err != nil {
			b.Fatalf("Unpack() failed: %v", err)
		}
	}
}
