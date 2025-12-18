package main
import "fmt"
func main() {
    data := make([]byte, 1050)
    fmt.Printf("Bytes 66-70: %02x %02x %02x %02x %02x\n", data[66], data[67], data[68], data[69], data[70])
}
