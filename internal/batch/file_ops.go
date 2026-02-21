package batch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CopyFile copies file content from src to dst, preserving mode and modification time.
// It respects the provided context for cancellation.
func CopyFile(ctx context.Context, src, dst string) error {
	sFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer sFile.Close()

	info, err := sFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	dFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create dest: %w", err)
	}

	// We close explicitly later, but defer ensures cleanup on error
	defer dFile.Close()

	// Copy using buffer and checking context
	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := sFile.Read(buf)
		if n > 0 {
			if _, wErr := dFile.Write(buf[:n]); wErr != nil {
				return fmt.Errorf("failed to write: %w", wErr)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}
	}

	// Close before metadata operations to ensure flush and release
	if err := dFile.Close(); err != nil {
		return fmt.Errorf("failed to close dest file: %w", err)
	}

	// Preserve permissions
	if err := os.Chmod(dst, info.Mode()); err != nil {
		return fmt.Errorf("failed to chmod: %w", err)
	}

	// Preserve mod time
	if err := os.Chtimes(dst, info.ModTime(), info.ModTime()); err != nil {
		return fmt.Errorf("failed to chtimes: %w", err)
	}

	return nil
}

// CalculateFileHash calculates the SHA256 hash of a file.
func CalculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file content: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
