package main
import "fmt"

func main() {
    // Chunk bytes: 05 00 00 00 00
    bytes := []byte{0x05, 0x00, 0x00, 0x00, 0x00}
    
    // Parser applies 128-offset normalization
    var sharpnessSum, sharpnessCount int
    for _, b := range bytes {
        if int(b) != 0 {  // Parser skips zeros
            adjusted := int(b) - 128
            sharpnessSum += adjusted
            sharpnessCount++
            fmt.Printf("Byte %02x → adjusted %d\n", b, adjusted)
        }
    }
    
    if sharpnessCount > 0 {
        avgSharpness := sharpnessSum / sharpnessCount
        sharpening := float64(avgSharpness) * 9.0 / 255.0
        fmt.Printf("Avg adjusted: %d, Sharpening: %.2f\n", avgSharpness, sharpening)
    }
}
