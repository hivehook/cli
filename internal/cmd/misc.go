package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(version)
			return nil
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show server status (version, queue depth, workers, DLQ size)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cl, err := client()
			if err != nil {
				return err
			}
			st, err := cl.Status.Get(cmd.Context())
			if err != nil {
				return err
			}
			return emit(st)
		},
	}
}
