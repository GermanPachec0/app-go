package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(APICmd(ctx))
	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
