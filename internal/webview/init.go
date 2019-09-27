package webview

import "runtime"

func init() {
	runtime.LockOSThread()
}
