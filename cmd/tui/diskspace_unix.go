//go:build !windows
// +build !windows

package main

import (
	"os"
	"syscall"
)

// getAvailableDiskSpace returns available disk space in bytes (Unix)
func getAvailableDiskSpace() int64 {
	cwd, err := os.Getwd()
	if err != nil {
		return 0
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(cwd, &stat); err == nil {
		// Available blocks * block size
		return int64(stat.Bavail) * int64(stat.Bsize)
	}

	return 0
}
