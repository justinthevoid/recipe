package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "Recipe CLI v0.1.0"
	verbose bool
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:   "recipe",
	Short: "Convert photo presets between formats",
	Long: `Recipe - Universal Photo Preset Converter

Convert photo presets between Nikon NP3, Adobe Lightroom XMP, and lrtemplate formats.

All processing happens locally on your device - files are never uploaded to any server.
Your privacy is guaranteed by design.

Supported formats:
  - NP3 (Nikon Picture Control binary format)
  - XMP (Adobe Lightroom sidecar XML)
  - lrtemplate (Adobe Lightroom Lua preset)

Examples:
  recipe convert portrait.xmp --to np3
  recipe convert --batch *.xmp --to np3
  recipe --help

Documentation: https://github.com/justin/recipe`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is specified, show help
		cmd.Help()
	},
}

func init() {
	// Initialize logger with default (non-verbose) settings
	// This will be reconfigured in PersistentPreRun if --verbose flag is used
	logger = initLogger(false)

	// Global flags available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output in JSON format")

	// Configure version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", version))
}

// Execute runs the root command and initializes the logger based on flags.
// This is called by main.go and is the entry point for all CLI operations.
func Execute() error {
	// Reinitialize logger based on verbose flag before command execution
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger = initLogger(verbose)
	}

	return rootCmd.Execute()
}
