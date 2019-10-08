package cmds

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(func(cmd *cobra.Command) *cobra.Command {
		return cmd
	}(&cobra.Command{
		Use:   "dbless",
		Short: "mycli dbless server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		}},
	))
}
