package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

type chunkData struct {
	id     uint32
	length uint16
	value  []byte
}

func parseChunks(data []byte, startOffset int) ([]chunkData, error) {
	var chunks []chunkData
	pos := startOffset

	for pos+12 <= len(data) {
		chunkID := binary.LittleEndian.Uint32(data[pos : pos+4])
		_ = binary.LittleEndian.Uint32(data[pos+4 : pos+8])
		valueLen := binary.LittleEndian.Uint16(data[pos+8 : pos+10])

		if pos+10+int(valueLen) > len(data) {
			break
		}

		valueStart := pos + 10
		valueEnd := valueStart + int(valueLen)
		value := make([]byte, valueLen)
		copy(value, data[valueStart:valueEnd])

		chunks = append(chunks, chunkData{
			id:     chunkID,
			length: valueLen,
			value:  value,
		})

		pos = valueEnd

		if len(chunks) > 100 {
			break
		}
	}

	return chunks, nil
}

func main() {
	files, _ := filepath.Glob("examples/np3/**/*.np3")

	for i, f := range files[:3] {
		fmt.Printf("\n=== File %d: %s ===\n", i+1, filepath.Base(f))
		data, _ := os.ReadFile(f)

		chunks, _ := parseChunks(data, 0x2C)
		fmt.Printf("Total chunks: %d\n", len(chunks))

		// Show first 15 chunk IDs and values if length=2
		for j, chunk := range chunks {
			if j >= 15 {
				break
			}
			if chunk.length == 2 {
				val := binary.LittleEndian.Uint16(chunk.value)
				fmt.Printf("Chunk %d: ID=%d, Value=0x%04X (%d)\n", j, chunk.id, val, val)
			} else {
				fmt.Printf("Chunk %d: ID=%d, Length=%d\n", j, chunk.id, chunk.length)
			}
		}
	}
}
