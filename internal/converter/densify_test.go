package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justin/recipe/internal/formats/np3"
	"github.com/justin/recipe/internal/models"
)

const densifySampleCount = 16

func TestDensifyPointCurve_NoOpSmallInput(t *testing.T) {
	tests := []struct {
		name   string
		points []models.ToneCurvePoint
	}{
		{name: "empty", points: nil},
		{name: "single point", points: []models.ToneCurvePoint{{Input: 128, Output: 128}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recipe := &models.UniversalRecipe{PointCurve: tc.points}
			densifyPointCurve(recipe)
			if len(recipe.PointCurve) != len(tc.points) {
				t.Fatalf("expected %d points untouched, got %d", len(tc.points), len(recipe.PointCurve))
			}
		})
	}
}

func TestDensifyPointCurve_LinearIdentity(t *testing.T) {
	recipe := &models.UniversalRecipe{
		PointCurve: []models.ToneCurvePoint{
			{Input: 0, Output: 0},
			{Input: 255, Output: 255},
		},
	}
	densifyPointCurve(recipe)
	if len(recipe.PointCurve) != densifySampleCount {
		t.Fatalf("expected %d points, got %d", densifySampleCount, len(recipe.PointCurve))
	}
	if recipe.PointCurve[0] != (models.ToneCurvePoint{Input: 0, Output: 0}) {
		t.Errorf("expected first point (0,0), got %+v", recipe.PointCurve[0])
	}
	if recipe.PointCurve[densifySampleCount-1] != (models.ToneCurvePoint{Input: 255, Output: 255}) {
		t.Errorf("expected last point (255,255), got %+v", recipe.PointCurve[densifySampleCount-1])
	}
	// With zero tangents at both endpoints and a straight (0,0)–(255,255) line, the cubic
	// Hermite degenerates to an S-curve that eases in at both ends. At midpoint (x≈136),
	// the output must still sit close to the diagonal — a tolerance of ±12 catches gross
	// pathologies without over-constraining the eased shape.
	midpoint := recipe.PointCurve[densifySampleCount/2]
	diff := midpoint.Output - midpoint.Input
	if diff < -12 || diff > 12 {
		t.Errorf("midpoint output too far from linear: %+v", midpoint)
	}
}

func TestDensifyPointCurve_ScottLeiter(t *testing.T) {
	path := filepath.Join("..", "formats", "np3", "testdata", "densify-curves-input.np3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	recipe, err := np3.Parse(data)
	if err != nil {
		t.Fatalf("parse NP3: %v", err)
	}
	if len(recipe.PointCurve) < 2 {
		t.Fatalf("fixture has no point curve: %+v", recipe.PointCurve)
	}
	originals := append([]models.ToneCurvePoint(nil), recipe.PointCurve...)

	densifyPointCurve(recipe)

	if len(recipe.PointCurve) != densifySampleCount {
		t.Fatalf("expected %d samples, got %d", densifySampleCount, len(recipe.PointCurve))
	}
	if recipe.PointCurve[0] != (models.ToneCurvePoint{Input: 0, Output: 0}) {
		t.Errorf("first sample must be (0,0), got %+v", recipe.PointCurve[0])
	}
	if recipe.PointCurve[densifySampleCount-1] != (models.ToneCurvePoint{Input: 255, Output: 255}) {
		t.Errorf("last sample must be (255,255), got %+v", recipe.PointCurve[densifySampleCount-1])
	}
	// Strict monotonicity of sample inputs.
	for i := 1; i < len(recipe.PointCurve); i++ {
		if recipe.PointCurve[i].Input <= recipe.PointCurve[i-1].Input {
			t.Errorf("sample inputs not strictly increasing at %d: %+v", i, recipe.PointCurve)
			break
		}
	}
	// Monotonic outputs (given input-monotonic Scott's Leiter data, F–C preserves it).
	for i := 1; i < len(recipe.PointCurve); i++ {
		if recipe.PointCurve[i].Output < recipe.PointCurve[i-1].Output {
			t.Errorf("output not monotonic at %d: %+v", i, recipe.PointCurve)
			break
		}
	}
	// Shape assertion: with the first original control point at (55,41), the densified
	// curve must bow BELOW the straight (0,0)→(55,41) diagonal — at x=17 the diagonal
	// gives y=12.7, so an eased curve must sit below that.
	sampleAt17 := recipe.PointCurve[1]
	if sampleAt17.Input != 17 {
		t.Fatalf("expected sample index 1 to have input=17, got %+v", sampleAt17)
	}
	if sampleAt17.Output >= 13 {
		t.Errorf("curve does not bow below diagonal at x=17: got y=%d (expected < 13)", sampleAt17.Output)
	}
	// Every original control point, when read from the densified curve via linear
	// interpolation (Lightroom-style), must land within ±2 of its original output.
	// ±2 absorbs cumulative rounding through sample + linear interpolation.
	for _, orig := range originals {
		y := linearInterp(recipe.PointCurve, orig.Input)
		diff := y - orig.Output
		if diff < -2 || diff > 2 {
			t.Errorf("reproduction of %+v via linear interp gave y=%d (diff %d)", orig, y, diff)
		}
	}
}

func TestDensifyPointCurve_FlattenWins(t *testing.T) {
	path := filepath.Join("..", "formats", "np3", "testdata", "densify-curves-input.np3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	out, err := ConvertWithOptions(data, FormatNP3, FormatXMP, ConvertOptions{
		FlattenCurves: true,
		DensifyCurves: true,
	})
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if strings.Contains(string(out), "ToneCurvePV2012") {
		t.Errorf("flatten+densify should remove ToneCurvePV2012, but it was present in output")
	}
}

func TestDensifyPointCurve_DefaultOffUsesRDP(t *testing.T) {
	path := filepath.Join("..", "formats", "np3", "testdata", "densify-curves-input.np3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	out, err := Convert(data, FormatNP3, FormatXMP)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	s := string(out)
	i := strings.Index(s, "<crs:ToneCurvePV2012>")
	if i < 0 {
		t.Fatalf("expected ToneCurvePV2012 in baseline output")
	}
	end := strings.Index(s[i:], "</crs:ToneCurvePV2012>")
	if end < 0 {
		t.Fatalf("unterminated ToneCurvePV2012")
	}
	count := strings.Count(s[i:i+end], "<rdf:li>")
	// Baseline (RDP simplified) must be well under the densified 16-point output.
	// The existing RDP targets ≤6 points; anything >10 indicates densification leaked.
	if count >= 10 {
		t.Errorf("default-off output should be RDP-simplified (≤~6 entries), got %d — densification appears to have leaked", count)
	}
}

func TestDensifyPointCurve_EndToEndWriteCount(t *testing.T) {
	path := filepath.Join("..", "formats", "np3", "testdata", "densify-curves-input.np3")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	out, err := ConvertWithOptions(data, FormatNP3, FormatXMP, ConvertOptions{DensifyCurves: true})
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	s := string(out)
	i := strings.Index(s, "<crs:ToneCurvePV2012>")
	if i < 0 {
		t.Fatalf("expected ToneCurvePV2012 in output, got:\n%s", firstN(s, 500))
	}
	end := strings.Index(s[i:], "</crs:ToneCurvePV2012>")
	if end < 0 {
		t.Fatalf("unterminated ToneCurvePV2012 element")
	}
	block := s[i : i+end]
	count := strings.Count(block, "<rdf:li>")
	if count != densifySampleCount {
		t.Errorf("expected %d <rdf:li> entries in master curve, got %d\nblock:\n%s", densifySampleCount, count, block)
	}
}

func linearInterp(pts []models.ToneCurvePoint, x int) int {
	if len(pts) == 0 {
		return 0
	}
	if x <= pts[0].Input {
		return pts[0].Output
	}
	if x >= pts[len(pts)-1].Input {
		return pts[len(pts)-1].Output
	}
	for i := 1; i < len(pts); i++ {
		if x <= pts[i].Input {
			x0, y0 := pts[i-1].Input, pts[i-1].Output
			x1, y1 := pts[i].Input, pts[i].Output
			if x1 == x0 {
				return y0
			}
			return y0 + (y1-y0)*(x-x0)/(x1-x0)
		}
	}
	return pts[len(pts)-1].Output
}

func firstN(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}
