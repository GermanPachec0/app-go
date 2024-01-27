package cmd

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(APICmd(ctx))
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
