package envs

import "github.com/pubgo/g/errors"

// Cfg default config
var Cfg = struct {
	Prefix string
	Env    string
}{
	Prefix: "",
	Env:    DevEnv.Dev,
}

func Init() {
	_env := GetEnv("env", "dev_env", "DevEnv")
	if _env != "" {
		errors.PanicT(_Match(_env), "dev env match error")
	}
	Cfg.Env = _env
}
