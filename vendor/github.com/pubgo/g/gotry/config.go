package gotry

import "github.com/pubgo/g/envs"

// Cfg default config
var Cfg = struct {
	Debug bool
}{
	Debug: false,
}

func init() {
	Cfg.Debug = envs.IsDebug()
}
