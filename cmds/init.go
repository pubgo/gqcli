package cmds

import (
	"github.com/pubgo/g/gcmds"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

func init() {
	rootCmd := gcmds.Default()
	rootCmd.AddCommand(&cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
