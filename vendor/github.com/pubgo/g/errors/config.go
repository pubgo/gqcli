package errors

import "time"

// Cfg config
var Cfg = struct {
	MaxObj        uint8
	MaxRetryDur   time.Duration
	Debug         bool
	Stack         bool
	SkipErrorFile bool
}{
	Debug:         false,
	Stack:         false,
	SkipErrorFile: false,
	MaxObj:        15,             // the max length of m keys, default=15
	MaxRetryDur:   time.Hour * 24, // the max retry duration, default=one day
}
