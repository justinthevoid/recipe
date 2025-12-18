package main

import (
	"fmt"
	"math"
)

// Simplified structures and functions from the codebase

type ControlPoint struct {
	X, Y int
}

func PointsToCurveLUT(points []ControlPoint) []int {
	// Simple linear interpolation
	lut := make([]int, 256)
	if len(points) == 0 {
		return lut
	}
	for i := 0; i < len(points)-1; i++ {
		p1 := points[i]
		p2 := points[i+1]
		for x := p1.X; x <= p2.X; x++ {
			if x >= 0 && x < 256 {
				t := float64(x-p1.X) / float64(p2.X-p1.X)
				y := float64(p1.Y) + t*float64(p2.Y-p1.Y)
				lut[x] = int(math.Round(y))
			}
		}
	}
	return lut
}

func ApplyCurveToLUT(base, modifier []int) []int {
	res := make([]int, 256)
	for i := 0; i < 256; i++ {
		out := modifier[base[i]] // modifier(base(i))
		res[i] = out
	}
	return res
}

func main() {
	// Agfachrome Data
	// Point Curve (Approx from XMP)
	pointCurve := []ControlPoint{
		{0, 13}, {65, 65}, {190, 191}, {255, 255},
	}
	pointLUT := PointsToCurveLUT(pointCurve)

	// Parametric Curve (Approx Darks=-10)
	// If Darks=-10, assume it pulls down. 
	// Let's create a proxy Parametric LUT.
	// Adobe Darks affects region around 64. 
	// -10 roughly means -10% shift?
	// Let's assume ideal behavior: 64->58.
	parametricCurve := []ControlPoint{
		{0, 0}, {32, 32}, {64, 58}, {128, 128}, {255, 255}, 
	}
	parametricLUT := PointsToCurveLUT(parametricCurve)

	// Standard Base Curve v7 (Old)
	// 64 -> 45
	baseV7 := []ControlPoint{{0,0}, {32,15}, {64,45}, {128,128}, {192,215}, {255,255}}
	baseV7LUT := PointsToCurveLUT(baseV7)

	// Standard Base Curve v9 (Aggressive)
	// 64 -> 35
	baseV9 := []ControlPoint{{0,0}, {32,10}, {64,35}, {128,128}, {192,220}, {255,255}}
	baseV9LUT := PointsToCurveLUT(baseV9)

	// Calc v7 (Exclusive Point, ignored Parametric)
	// Final = Point(Base(x))
	v7LUT := ApplyCurveToLUT(baseV7LUT, pointLUT)

	// Calc v9 (Chained, Aggressive Base)
	// Final = Point(Parametric(Base(x)))
	// Inner = ApplyCurveToLUT(Parametric, Point) // Point(Parametric(x))
	chainedLUT := ApplyCurveToLUT(parametricLUT, pointLUT)
	v9LUT := ApplyCurveToLUT(baseV9LUT, chainedLUT)

	fmt.Printf("Input | v7 Output | v9 Output\n")
	fmt.Printf("------------------------------\n")
	inputs := []int{0, 10, 32, 64, 128}
	for _, x := range inputs {
		fmt.Printf("%3d   | %3d       | %3d\n", x, v7LUT[x], v9LUT[x])
	}
}
