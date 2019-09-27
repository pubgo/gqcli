package errors

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
)

// ErrUnknownType error
var ErrUnknownType = errors.New("unknown type error")

type errF1 func() (err error)
type errF2 func(...interface{}) (err error)
type errF3 func(...interface{}) func(...interface{}) error
type errF4 func(...interface{}) func(...interface{}) func(...interface{}) error

func _handle(err interface{}) *Err {
	if err == nil || _isNone(err) {
		return nil
	}

	switch _e := err.(type) {
	case errF1:
		err = _e()
	case errF2:
		err = _e()
	case errF3:
		err = _e()()
	case errF4:
		err = _e()()()
	}

	m := &Err{}
	switch _e := err.(type) {
	case *Err:
		m = _e
	case error:
		m.err = _e
		m.msg = m.err.Error()
	case string:
		m.msg = _e
		m.err = errors.New(m.msg)
	default:
		m.msg = fmt.Sprintf("unknown type error, input: %#v", _e)
		m.err = ErrUnknownType
		m.tag = ErrTag.UnknownTypeCode
	}
	return m
}

// Assert errors
func Assert() {
	ErrHandle(recover(), func(err *Err) {
		if Cfg.Debug {
			fmt.Println(err.P())
		}
		if Cfg.Stack {
			debug.PrintStack()
		}
		os.Exit(1)
	})
}

// Resp errors
func Resp(fn func(err *Err)) {
	ErrHandle(recover(), func(err *Err) {
		err.Caller(_caller.FromFunc(reflect.ValueOf(fn)))
		fn(err)
	})
}

// RespErr errors
func RespErr(err *error) {
	ErrHandle(recover(), func(_err *Err) {
		*err = _err
	})
}

// ErrLog errors
func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.Caller(_caller.FromDepth(callDepth))
		fmt.Println(err.P())
	})
}

// Debug errors
func Debug() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
		debug.PrintStack()
	})
}

// ErrHandle errors
func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if err == nil || _isNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) || _m.err == nil {
		return
	}

	if len(fn) == 0 {
		return
	}

	PanicM(AssertFn(reflect.ValueOf(fn[0])), "func error")
	fn[0](_m)
}
