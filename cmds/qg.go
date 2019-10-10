package cmds

import (
	"github.com/pubgo/g/gcmds"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd := gcmds.Default()
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gq",
		Short: "mycli gq",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}
