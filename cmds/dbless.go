package cmds

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/pubgo/g/gcmds"
	"github.com/pubgo/mycli/internal/dbless/rest"
	"github.com/pubgo/mycli/pkg/logs"
	"github.com/pubgo/mycli/pkg/restorm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd := gcmds.Default()
	rootCmd.AddCommand(func(cmd *cobra.Command) *cobra.Command {
		return cmd
	}(&cobra.Command{
		Use:   "dbless",
		Short: "mycli dbless server",
		RunE: func(cmd *cobra.Command, args []string) error {
			rtm := restorm.Default()
			rtm.DbAdd("mydb", &restorm.Config{
				Enable:  true,
				Driver:  "sqlite3",
				Dsn:     "file:test.db?cache=shared",
				ShowSQL: true,
			})

			logs.P("rtm.DbStats()", rtm.DbStats())

			return rest.App().Run(":8080")
		}},
	))
}
