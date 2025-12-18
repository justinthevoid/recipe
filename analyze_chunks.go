package main
import "fmt"

func main() {
    // Chunk #2 (offset 66-75): 05 00 00 00 00 02 ff 01 00 00
    // Chunk #3 (offset 76-85): 06 00 00 00 00 02 74 04 00 00
    // Chunk #4 (offset 86-95): 07 00 00 00 00 02 7e 04 00 00
    
    chunk2Value := []byte{0xff, 0x01}
    chunk3Value := []byte{0x74, 0x04}
    chunk4Value := []byte{0x7e, 0x04}
    
    fmt.Printf("Chunk #2 (id=0x05) value: 0x%02x%02x\n", chunk2Value[0], chunk2Value[1])
    fmt.Printf("Chunk #3 (id=0x06) value: 0x%02x%02x = sharpness byte 0x%02x\n", chunk3Value[0], chunk3Value[1], chunk3Value[0])
    fmt.Printf("Chunk #4 (id=0x07) value: 0x%02x%02x = clarity byte 0x%02x\n", chunk4Value[0], chunk4Value[1], chunk4Value[0])
    
    // Decode sharpness (offset 82 = chunk #3 value byte 0)
    sharpnessRaw := chunk3Value[0] // 0x74 = 116 decimal
    sharpnessAdjusted := int(sharpnessRaw) - 128 // 116 - 128 = -12
    sharpening := float64(sharpnessAdjusted) * 9.0 / 255.0 // -12 * 9/255 = -0.42
    fmt.Printf("Sharpening from chunk: %.2f (expected 0.0)\n", sharpening)
    
    // Decode clarity (offset 92 = chunk #4 value byte 0, but offset 92 is NOT in chunk #4!)
    // Chunk #4 is at 86-95, so offset 92 is at position 6 within the chunk
    // That's the value field! So chunk value should contain clarity
    clarityRaw := chunk4Value[0] // 0x7e = 126 decimal
    clarityAdjusted := int(clarityRaw) - 128 // 126 - 128 = -2
    clarity := float64(clarityAdjusted) * 5.0 / 127.0 // -2 * 5/127 = -0.08
    fmt.Printf("Clarity from chunk: %.2f (expected -10.0)\n", clarity)
}
