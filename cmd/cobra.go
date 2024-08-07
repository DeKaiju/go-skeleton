package cmd

import (
	"errors"
	"os"

	"github.com/dekaiju/go-skeleton/cmd/api"
	"github.com/dekaiju/go-skeleton/cmd/scripts"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "entry",
	SilenceUsage: false,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run:               func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(api.Server)
	rootCmd.AddCommand(scripts.Register)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
