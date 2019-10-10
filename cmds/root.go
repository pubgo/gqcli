package cmds

import (
	"github.com/pubgo/g/gcmds"
	"github.com/pubgo/mycli/version"
	"github.com/spf13/cobra"
)

var Execute = gcmds.Execute
var _ = gcmds.Default(func(cmd *cobra.Command) {
	cmd.Use = "mycli"
	cmd.Version = version.Version
})
