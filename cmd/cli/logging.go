package main

import (
	"log/slog"
	"os"
)

// logger is the global logger instance used by all CLI commands
var logger *slog.Logger

// initLogger initializes the slog logger with appropriate level based on verbose flag.
// Returns configured logger writing to stderr.
//
// Parameters:
//   - verbose: If true, sets log level to Debug (all logs). If false, sets to Error (minimal output).
//
// The logger uses TextHandler for human-readable output and writes all logs to stderr
// to keep stdout clean for piping and JSON output.
func initLogger(verbose bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelError, // Default: errors only
	}

	if verbose {
		opts.Level = slog.LevelDebug // Verbose: all levels (Debug, Info, Warn, Error)
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler)
}
