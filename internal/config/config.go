package config

import (
	"github.com/pubgo/g/envs"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/g/logs"
	"github.com/pubgo/mycli/pdd/cfg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"sync"
)

type config struct {
	Cfg      *cfg.Cfg
	Debug    bool
	Env      string
	LogLevel string
}

func (t *config) InitCfg() {
	log.Debug().Msg("init cfg")

	t.Debug = envs.IsDebug()
	t.Env = envs.Cfg.Env

	t.LogLevel = t.Cfg.Log.Level
	if _l, ok := os.LookupEnv("log_level"); ok {
		t.LogLevel = _l
		errors.PanicT(!logs.Match(t.LogLevel), "the env value is not match,value(%s)", _l)
	}
	if _d := viper.Get("ll"); _d != nil {
		t.LogLevel = viper.GetString("ll")
		errors.PanicT(!logs.Match(t.LogLevel), "the env value is not match,value(%s)", t.LogLevel)
	}
	log.Debug().Msg("init cfg ok")
}

func (t *config) Init() {

	// cfg 初始化
	t.InitCfg()

	// log初始化
	t.InitLog()

}

func (t *config) Parse() {
	log.Debug().Msg("parse cfg")
	errors.PanicM(viper.Unmarshal(t.Cfg), "cfg parse error")
	log.Debug().Msg("parse cfg ok")
}

var _cfg *config
var once sync.Once

func Default() *config {
	once.Do(func() {
		_cfg = &config{Cfg: new(cfg.Cfg)}
		_cfg.LogLevel = zerolog.DebugLevel.String()
		_cfg.Debug = true
	})
	return _cfg
}
