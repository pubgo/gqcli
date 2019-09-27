package cmds

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:     "server",
		Aliases: []string{"s"},
		Short:   "mycli server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}
