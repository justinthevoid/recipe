package batch

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies file content from src to dst, preserving mode and modification time.
func CopyFile(src, dst string) error {
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
	// Note: dFile.Close() usually returns error if already closed? Go docs say "returns an error, if any".
	// But calling it twice is generally safe logic-wise if we don't check second error or if os.File handles it.
	// Actually os.File.Close() returns error if already closed (PathError).
	// But defer ignores return value.
	defer dFile.Close()

	if _, err := io.Copy(dFile, sFile); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Close before metadata operations to ensure flush and release (important for Windows)
	// We assign the error check here.
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
