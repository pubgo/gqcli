package config

import "github.com/pubgo/g/logs"

func (t *config) InitLog() {
	logs.Default().Init()
}
