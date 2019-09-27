package envs

import (
	"github.com/pubgo/g/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func SetEnv(e, v string) error {
	return errors.Wrap(os.Setenv(_EnvKey(e), v), "set env error")
}

func GetEnv(e ...string) string {
	for _, _e := range e {
		if _v := os.Getenv(_EnvKey(_e)); _v != "" {
			return _v
		}
	}
	return ""
}

// Env env
func _EnvKey(e string) string {
	e = strings.ToUpper(e)
	if Cfg.Prefix == "" {
		return e
	}

	return Cfg.Prefix + "_" + e
}

// LoadFile 加载.env文件并添加前缀
func LoadFile(envFiles ...string) (err error) {
	defer errors.RespErr(&err)

	for _, filename := range filenamesOrDefault(envFiles) {
		_envPath, err := filepath.EvalSymlinks(filename)
		errors.PanicM(err, "%s EvalSymlinks error", filename)

		f, err := ioutil.ReadFile(_envPath)
		errors.PanicM(err, "%s ReadFile error", filename)

		for _, env := range strings.Split(string(f), "\n") {
			if _envs := strings.SplitN(strings.TrimSpace(env), "=", 2); len(_envs) == 2 {
				_env := _EnvKey(_envs[0])
				errors.Panic(os.Unsetenv(_env))
				errors.PanicM(os.Setenv(_env, doubleQuoteEscape(_envs[1])), "set env(%s=%s) error", _envs[0], _envs[1])
			}
		}
	}

	return
}

// Parse 解析envs map[string]string
func Parse() map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		e = strings.TrimSpace(e)

		if !strings.HasPrefix(Cfg.Prefix, e) {
			continue
		}

		if ev := strings.Split(e, "="); len(ev) == 2 && ev[0] != "" {
			envs[strings.ToUpper(ev[0])] = doubleQuoteEscape(ev[1])
		}
	}
	return envs
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}
