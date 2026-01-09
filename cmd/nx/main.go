package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var version = "dev"
var logLevel string

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "recipe-nx",
		Short:   "Recipe NX Studio Integration",
		Long:    `recipe-nx integrates with Nikon NX Studio to apply recipes to NEF files.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return setupLogger(logLevel)
		},
	}

	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

	batchCmd := newBatchCmd()
	batchCmd.AddCommand(newApplyCmd())
	cmd.AddCommand(batchCmd)

	return cmd
}

func setupLogger(level string) error {
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	var handler slog.Handler
	if term.IsTerminal(int(os.Stderr.Fd())) {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
