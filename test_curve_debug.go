package main

import (
	"fmt"
	"os"
	"github.com/justin/recipe/internal/formats/xmp"
)

func main() {
	data, _ := os.ReadFile("testdata/xmp/Fuji Industrial 400.xmp")
	recipe, err := xmp.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Preset name: %s\n", recipe.Name)
	fmt.Printf("PointCurve length: %d\n", len(recipe.PointCurve))
	if len(recipe.PointCurve) > 0 {
		fmt.Println("First 5 points:")
		for i := 0; i < 5 && i < len(recipe.PointCurve); i++ {
			fmt.Printf("  [%d] Input=%d, Output=%d\n", i, recipe.PointCurve[i].Input, recipe.PointCurve[i].Output)
		}
	} else {
		fmt.Println("PointCurve is EMPTY or NIL")
	}
}
