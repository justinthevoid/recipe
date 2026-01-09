// Package apperr provides standardized error handling types for the application.
//
// It encapsulates the domain error strategy, offering structured error types associated
// with specific operations and files. This package ensures that errors can be wrapped
// with context (Operation, File, Cause) and unwrapped for causal analysis, compatible
// with the standard library `errors` package.
//
// This package is intended to be used by `internal/batch`, `internal/formats`, and
// `cmd/nx` to maintain consistent error reporting and logging.
package apperr
