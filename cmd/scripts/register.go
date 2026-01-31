package scripts

import (
	"github.com/spf13/cobra"

	"github.com/dekaiju/go-skeleton/cmd/scripts/demo"
)

var (
	Register = &cobra.Command{
		Use: "cmd",
	}
)

func init() {
	Register.AddCommand(demo.Register)
}
