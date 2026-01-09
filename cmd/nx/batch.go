package main

import "github.com/spf13/cobra"

func newBatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch",
		Short: "Batch process directories",
		Long:  `Batch processing commands for Nikon NX Studio sidecars.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is provided, show help
			return cmd.Help()
		},
	}
}
