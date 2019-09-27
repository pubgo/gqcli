package envs

// Env env配置
type _Env struct {
	Dev  string
	Stag string
	Prod string
}

var DevEnv = _Env{
	Dev:  "dev",
	Stag: "stag",
	Prod: "prod",
}

// Match 匹配环境
func _Match(envs ...string) bool {
	for _, e := range envs {
		if DevEnv.Dev == e || DevEnv.Stag == e || DevEnv.Prod == e {
			return true
		}
	}
	return false
}
