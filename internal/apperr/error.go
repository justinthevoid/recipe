// Package apperr provides domain-specific error types and structured error wrapping
// for the Recipe application.
//
// It defines a standard Error struct that captures the operation phase,
// file context, and underlying cause, ensuring consistent logging and user feedback.
// It also provides sentinel errors for common domain failures.
package apperr

import (
	"errors"
	"fmt"
)

// Sentinel errors for known domain failure modes.
var (
	ErrInvalidNP3 = errors.New("invalid NP3 data")
	ErrCorruptNEF = errors.New("corrupt NEF file")
)

// Error represents a structured domain error with context.
type Error struct {
	Operation string            // The phase or operation (e.g., "parse", "process")
	File      string            // The file being processed (base name)
	Cause     error             // The underlying error
	Context   map[string]string // Additional key-value context
}

// New creates a new structured Error.
func New(op, file string, cause error) *Error {
	return &Error{
		Operation: op,
		File:      file,
		Cause:     cause,
		Context:   make(map[string]string),
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.File != "" {
		return fmt.Sprintf("%s %s: %v", e.Operation, e.File, e.Cause)
	}
	return fmt.Sprintf("%s: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying cause of the error.
func (e *Error) Unwrap() error {
	return e.Cause
}

// With adds a key-value pair to the error's context.
// It returns the error itself for chaining.
func (e *Error) With(key, value string) *Error {
	e.Context[key] = value
	return e
}
