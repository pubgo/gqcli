package main

import (
	"fmt"
	"github.com/zserge/webview"
)

func main() {
	// Open wikipedia in a 800x600 resizable window
	w := webview.New(webview.Settings{
		Width:     1200,
		Height:    600,
		Title:     "Simple canvas demo",
		URL:       "https://dd.diandao.pro",
		Resizable: true,
		Debug:     true,
		ExternalInvokeCallback: func(w webview.WebView, data string) {
			fmt.Println(w, data)
		},
	})

	w.Dispatch(func() {
		w.Eval(`external.invoke('close');`)
	})
	defer w.Exit()
	w.Run()
}
