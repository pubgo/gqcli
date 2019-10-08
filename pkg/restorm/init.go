package restorm

import (
	"github.com/pubgo/g/gotry"
	"github.com/pubgo/g/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _isNone = pkg.IsNone
var _retryAt = gotry.RetryAt


var logger zerolog.Logger
func init() {
	logger = log.With().Str("pkg", "restorm").Logger()
}
