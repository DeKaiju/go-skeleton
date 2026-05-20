package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dekaiju/go-skeleton/cmd/api"
	"github.com/dekaiju/go-skeleton/cmd/scripts"
)

var rootCmd = &cobra.Command{
	Use:          "go-skeleton",
	Short:        "A reusable Go service skeleton",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(api.Server)
	rootCmd.AddCommand(scripts.Register)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
