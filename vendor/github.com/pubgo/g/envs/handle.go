package envs

import (
	"github.com/pubgo/g/errors"
)

// SetDebug 设置debug模式
func SetDebug() error {
	return errors.Wrap(SetEnv("debug", "true"), "set debug env error")
}

// SetSkipErrorFile 设置跳过错误
func SetSkipErrorFile() error {
	return errors.Wrap(SetEnv("skip_error_file", "true"), "set skip_error_file env error")
}

// IsDebug 是否是debug模式
func IsDebug() bool {
	debug := GetEnv("debug")
	return debug == "true" || debug == "t" || debug == "1" || debug == "ok"
}

// IsSkipErrorFile 是否跳过错误文件
func IsSkipErrorFile() bool {
	skipErrorFile := GetEnv("skip_error_file")
	return skipErrorFile == "true" || skipErrorFile == "t" || skipErrorFile == "1" || skipErrorFile == "ok"
}

// IsDev 是否是dev模式
func IsDev() bool {
	return Cfg.Env == DevEnv.Dev
}

// IsStag 是否是stag模式
func IsStag() bool {
	return Cfg.Env == DevEnv.Stag
}

// IsProd 是否是prod模式
func IsProd() bool {
	return Cfg.Env == DevEnv.Prod
}
