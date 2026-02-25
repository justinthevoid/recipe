package np3

import "errors"

var (
	// ErrFileTooSmall indicates the provided byte slice is shorter than the minimum expected NP3 file size.
	ErrFileTooSmall = errors.New("np3: file too small or truncated")

	// ErrFileTooLarge indicates the provided byte slice exceeds the maximum permitted capacity (1MB).
	ErrFileTooLarge = errors.New("np3: file exceeds 1MB capacity limit")

	// ErrInvalidMagic indicates the file header does not match the expected "NCP" signature.
	ErrInvalidMagic = errors.New("np3: invalid magic bytes or corrupted header")

	// ErrChecksumMismatch indicates the calculated checksum does not match the embedded checksum.
	ErrChecksumMismatch = errors.New("np3: checksum mismatch")
)
