package errors

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

var _funcCaller = func(callDepth int) []string {
	return []string{_caller.FromDepth(callDepth), _caller.FromDepth(callDepth + 1)}
}

// Panic errors
func Panic(err interface{}) {
	if err == nil || _isNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) {
		return
	}

	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: _funcCaller(callDepth + 1),
	})
}

// PanicT errors
func PanicT(b bool, msg string, args ...interface{}) {
	if !b {
		return
	}

	panic(&Err{
		err:    fmt.Errorf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	})
}

// PanicTT errors
func PanicTT(b bool, fn func(err *Err)) {
	if !b {
		return
	}

	_err := &Err{caller: _funcCaller(callDepth + 1)}
	fn(_err)

	if _err.msg == "" {
		log.Fatalf("msg is null")
	}
	_err.err = errors.New(_err.msg)

	panic(_err)
}

// PanicM errors
func PanicM(err interface{}, msg string, args ...interface{}) {
	if err == nil || _isNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) {
		return
	}

	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	})
}

// PanicMM errors
func PanicMM(err interface{}, fn func(err *Err)) {
	if err == nil || _isNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) {
		return
	}

	_err := &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: _funcCaller(callDepth + 1),
	}
	fn(_err)
	if _err.msg == "" {
		log.Fatalf("msg is null")
	}
	panic(_err)
}

// Wrap errors
func Wrap(err interface{}, msg string, args ...interface{}) error {
	if err == nil || _isNone(err) {
		return nil
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) {
		return nil
	}

	return &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	}
}

// WrapM errors
func WrapM(err interface{}, fn func(err *Err)) error {
	if err == nil || _isNone(err) {
		return nil
	}

	_m := _handle(err)
	if _m == nil || _isNone(_m) {
		return nil
	}

	_err := &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: _funcCaller(callDepth + 1),
	}
	fn(_err)
	if _err.msg == "" {
		log.Fatalf("msg is null")
	}
	return _err
}

// AssertFn errors
func AssertFn(fn reflect.Value) error {
	if _isZero(fn) || fn.Kind() != reflect.Func {
		return fmt.Errorf("the func is nil[%#v] or not func type[%s]", fn, fn.Kind())
	}
	return nil
}
