package main

import (
    "fmt"
    "github.com/justin/recipe/internal/formats/np3"
)

func main() {
    // 1. Create Identity LUT (Standard workflow assumption)
    lut := make([]int, 256)
    for i := range lut {
        lut[i] = i
    }

    // 2. Apply Compensation
    fmt.Println("Applying Flexible Color Baseline Compensation...")
    compensated := np3.ApplyFlexibleColorBaselineCompensation(lut)

    // 3. Print Analysis
    fmt.Println("\nCurve Analysis (Input -> Output):")
    fmt.Println("Idx | Original | Compensated | Diff")
    fmt.Println("----|----------|-------------|-----")
    
    // Sample points
    points := []int{0, 10, 32, 64, 96, 128, 160, 192, 224, 250, 255}
    
    for _, i := range points {
        orig := lut[i]
        comp := compensated[i]
        diff := comp - orig
        fmt.Printf("%3d | %3d      | %3d         | %3d\n", i, orig, comp, diff)
    }

    // Check monotonicity
    for i := 0; i < 255; i++ {
        if compensated[i] > compensated[i+1] {
            fmt.Printf("WARNING: Non-monotonic at index %d: %d > %d\n", i, compensated[i], compensated[i+1])
        }
    }
}
