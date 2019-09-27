package logs

import (
	"encoding/json"
	"fmt"
	"github.com/pubgo/g/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// P log
func P(s string, d ...interface{}) {
	fmt.Print(s)
	for _, i := range d {
		if i == nil || _isNone(i) {
			continue
		}

		dt, err := json.MarshalIndent(i, "", "\t")
		errors.PanicM(err, "P json MarshalIndent error")
		fmt.Println(string(dt))
	}
}

// LevelMatch 匹配log level
func Match(e string) bool {
	return zerolog.DebugLevel.String() == e ||
		zerolog.ErrorLevel.String() == e ||
		zerolog.WarnLevel.String() == e ||
		zerolog.FatalLevel.String() == e ||
		zerolog.InfoLevel.String() == e ||
		zerolog.PanicLevel.String() == e
}

// InitDebugLog init log
func InitDebugLog() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Timestamp().Logger()
}
