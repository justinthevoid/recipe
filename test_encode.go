package main
import "fmt"

func EncodeSigned8(value int) byte {
	adjusted := value + 128
	if adjusted > 255 {
		adjusted = 255
	} else if adjusted < 0 {
		adjusted = 0
	}
	return byte(adjusted)
}

func main() {
	fmt.Printf("EncodeSigned8(50) = 0x%02x (%d)\n", EncodeSigned8(50), EncodeSigned8(50))
	fmt.Printf("EncodeSigned8(0) = 0x%02x (%d)\n", EncodeSigned8(0), EncodeSigned8(0))
}
