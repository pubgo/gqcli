package cmds

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pubgo/g/envs"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/mycli/internal/config"
	"github.com/pubgo/mycli/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const homeFlag = "home"
const debugFlag = "debug"

var rootCmd = &cobra.Command{
	Use:     "mycli",
	Short:   "mycli app",
	Version: version.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		defer errors.Assert()

		for _, name := range []string{"version"} {
			if name == cmd.Name() {
				return
			}
		}

		errors.PanicM(viper.BindPFlags(cmd.Flags()), "Flags Error")

		viper.SetConfigType("toml")
		viper.SetConfigName("config")

		homeDir := viper.GetString(homeFlag)
		viper.AddConfigPath(homeDir)                          // search root directory
		viper.AddConfigPath(filepath.Join(homeDir, "config")) // search root directory /kdata

		home, err := homedir.Dir()
		errors.PanicM(err, "home dir error")
		viper.AddConfigPath(home)

		viper.AddConfigPath("/app/config")
		viper.AddConfigPath("config")
		viper.AddConfigPath("pdd/config")
		viper.AddConfigPath("pdd/cfg")

		// load data
		errors.PanicM(viper.ReadInConfig(), "check kata error")

		config.Default().Parse()
	},
}

func Execute(envPrefix, defaultHome string) {
	defer errors.Assert()

	cobra.OnInitialize(func() { initEnv(envPrefix) })
	rootCmd.PersistentFlags().StringP(homeFlag, "", defaultHome, "directory for data")
	rootCmd.PersistentFlags().BoolP(debugFlag, "d", true, "debug mode")
	rootCmd.PersistentFlags().StringP("ll", "l", "debug", "log level")
	rootCmd.PersistentFlags().StringP("env", "e", "dev", "project environment")
	errors.Panic(rootCmd.Execute())
}

// initEnv sets to use ENV variables if set.
func initEnv(prefix string) {
	// set global env prefix
	envs.Cfg.Prefix = prefix
	envs.Init()

	// env variables with TM prefix (eg. TM_ROOT)
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}
