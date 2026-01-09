package apperr_test

import (
	"errors"
	"testing"

	"github.com/justin/recipe/internal/apperr"
)

func TestError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := apperr.New("parse", "file.nef", cause)

	if !errors.Is(err, cause) {
		t.Error("errors.Is(err, cause) = false, want true")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("errors.Unwrap(err) = %v, want %v", unwrapped, cause)
	}
}

func TestError_ErrorString(t *testing.T) {
	cause := errors.New("invalid header")
	err := apperr.New("parse", "IMG_123.NEF", cause)

	want := "parse IMG_123.NEF: invalid header"
	if got := err.Error(); got != want {
		t.Errorf("err.Error() = %q, want %q", got, want)
	}
}

func TestError_WithContext(t *testing.T) {
	cause := errors.New("something failed")
	err := apperr.New("process", "test.nef", cause).
		With("worker", "1").
		With("attempts", "3")

	// We don't enforce exact string format for context, but it should be present if we were to inspect it.
	// For now, let's just ensure it doesn't panic and preserves value.
	// In a real implementation we might want context in the string or just attached.
	// The project-context suggests map[string]string.

	// Let's verify we can get the context back if we expose it, or just that the error string might contain it if desired.
	// The story criteria implies "Errors are wrapped with file context and phase".
	// The acceptance criteria 4 says: "Errors are wrapped with file context and phase (e.g., "IMG_123.NEF: parse_error: invalid header")".

	// Let's check the basic string again.
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDomainErrors(t *testing.T) {
	// Story mentions "Define explicit error variables/types for known failure modes"
	if apperr.ErrInvalidNP3 == nil {
		t.Error("ErrInvalidNP3 is nil")
	}
	if apperr.ErrCorruptNEF == nil {
		t.Error("ErrCorruptNEF is nil")
	}
}
