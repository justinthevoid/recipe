package converter

import (
	"math"
	"sort"

	"github.com/justin/recipe/internal/models"
)

// densifyPointCurve replaces recipe.PointCurve with 16 evenly-spaced samples of a
// Fritsch–Carlson monotonic cubic Hermite spline through the original control points,
// with (0,0) and (255,255) added as implicit endpoints when absent. Zero tangents are
// forced at the first and last points so the curve eases into the black and white
// anchors the way NX Studio renders tone curves — producing a "toe" region that
// Lightroom's linear-segment rendering can approximate.
//
// No-ops when len(PointCurve) < 2. The caller is responsible for setting
// recipe.SkipCurveSimplification so xmp.Generate skips its RDP pass.
func densifyPointCurve(recipe *models.UniversalRecipe) {
	if len(recipe.PointCurve) < 2 {
		return
	}

	pts := make([]models.ToneCurvePoint, len(recipe.PointCurve))
	copy(pts, recipe.PointCurve)
	sort.Slice(pts, func(i, j int) bool { return pts[i].Input < pts[j].Input })

	// Dedupe duplicate-input points (keep the last Y for each unique X). A malformed
	// NP3 with two points at the same input would otherwise produce h=0 intervals
	// that anchor the final densified output to the wrong Y.
	dedup := pts[:0:len(pts)]
	for i, p := range pts {
		if i > 0 && p.Input == dedup[len(dedup)-1].Input {
			dedup[len(dedup)-1] = p
			continue
		}
		dedup = append(dedup, p)
	}
	pts = dedup

	if pts[0].Input > 0 {
		pts = append([]models.ToneCurvePoint{{Input: 0, Output: 0}}, pts...)
	}
	if pts[len(pts)-1].Input < 255 {
		pts = append(pts, models.ToneCurvePoint{Input: 255, Output: 255})
	}

	n := len(pts)
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i, p := range pts {
		xs[i] = float64(p.Input)
		ys[i] = float64(p.Output)
	}

	d := make([]float64, n-1)
	for k := 0; k < n-1; k++ {
		dx := xs[k+1] - xs[k]
		if dx == 0 {
			d[k] = 0
		} else {
			d[k] = (ys[k+1] - ys[k]) / dx
		}
	}

	m := make([]float64, n)
	for k := 1; k < n-1; k++ {
		m[k] = (d[k-1] + d[k]) / 2
	}
	// Zero tangents at the anchor endpoints so the curve eases in/out of (0,0)
	// and (255,255) instead of tracing the straight-line secant.
	m[0] = 0
	m[n-1] = 0

	for k := 0; k < n-1; k++ {
		if d[k] == 0 {
			m[k] = 0
			m[k+1] = 0
			continue
		}
		alpha := m[k] / d[k]
		beta := m[k+1] / d[k]
		s := alpha*alpha + beta*beta
		if s > 9 {
			scale := 3 / math.Sqrt(s)
			m[k] *= scale
			m[k+1] *= scale
		}
	}

	const sampleCount = 16
	out := make([]models.ToneCurvePoint, sampleCount)
	k := 0
	for i := 0; i < sampleCount; i++ {
		x := float64(i) * 255 / float64(sampleCount-1)
		for k < n-2 && x > xs[k+1] {
			k++
		}
		h := xs[k+1] - xs[k]
		var y float64
		if h == 0 {
			y = ys[k]
		} else {
			t := (x - xs[k]) / h
			t2 := t * t
			t3 := t2 * t
			h00 := 2*t3 - 3*t2 + 1
			h10 := t3 - 2*t2 + t
			h01 := -2*t3 + 3*t2
			h11 := t3 - t2
			y = h00*ys[k] + h10*h*m[k] + h01*ys[k+1] + h11*h*m[k+1]
		}
		xi := int(math.Round(x))
		yi := int(math.Round(y))
		if yi < 0 {
			yi = 0
		} else if yi > 255 {
			yi = 255
		}
		out[i] = models.ToneCurvePoint{Input: xi, Output: yi}
	}

	recipe.PointCurve = out
}
