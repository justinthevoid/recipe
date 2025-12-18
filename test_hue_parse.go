package main
import "fmt"

func main() {
    // Offsets 76-79 in known-good file: 06 00 00 00
    // Offsets 76-79 in our file: 06 00 00 00 (same!)
    bytes := []byte{0x06, 0x00, 0x00, 0x00}
    
    // Parser reads these with 128-offset normalization
    var hueSum, hueCount int
    for _, b := range bytes {
        if b != 0 {  // Skip zeros
            adjusted := int(b) - 128
            hueSum += adjusted
            hueCount++
            fmt.Printf("Byte 0x%02x → adjusted %d\n", b, adjusted)
        }
    }
    
    if hueCount > 0 {
        avgHue := hueSum / hueCount
        hue := avgHue * 9 / 128
        fmt.Printf("Avg adjusted: %d, Hue: %d\n", avgHue, hue)
    } else {
        fmt.Println("No non-zero bytes, parser will use default hue")
    }
}
