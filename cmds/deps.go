package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "deps",
		Short: "deps info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gen struct from toml, https://xuri.me/toml-to-go")
		},
	})
}
