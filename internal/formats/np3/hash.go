package np3

import (
	"crypto/sha256"
	"encoding/hex"
)

// CalculateMagicHash computes a SHA-256 hash of the provided NP3 file bytes.
// This is used to ensure data integrity, idempotency, and track file state without relying on file paths.
func CalculateMagicHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
