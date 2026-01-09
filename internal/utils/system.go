package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenFolder opens the specified folder in the default OS file explorer.
// It supports Windows (explorer), macOS (open), and Linux (xdg-open).
func OpenFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open folder: %w", err)
	}
	return nil
}
