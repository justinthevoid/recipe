//go:build windows
// +build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// getAvailableDiskSpace returns available disk space in bytes (Windows)
func getAvailableDiskSpace() int64 {
	cwd, err := os.Getwd()
	if err != nil {
		return 0
	}

	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable, totalBytes, totalFreeBytes int64
	_, _, _ = c.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cwd))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)))

	return freeBytesAvailable
}
