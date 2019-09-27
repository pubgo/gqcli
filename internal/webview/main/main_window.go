package main

import (
	"fmt"
	"github.com/zserge/webview"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {

	var ok = false
	var ok1 = false

	w := webview.New(webview.Settings{
		Width:  1200,
		Height: 600,
		//Title:  "Loaded: Injected via JavaScript",
		URL: "https://www.yuque.com/kai.fangk/wave-balance/lyyup6",
		Title: `data:text/html,` + url.PathEscape(`
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
		<button onclick="external.invoke('close')">Close</button>
		<button onclick="external.invoke('fullscreen')">Fullscreen</button>
		<button onclick="external.invoke('unfullscreen')">Unfullscreen</button>
		<button onclick="external.invoke('open')">Open</button>
		<button onclick="external.invoke('opendir')">Open directory</button>
		<button onclick="external.invoke('save')">Save</button>
		<button onclick="external.invoke('message')">Message</button>
		<button onclick="external.invoke('info')">Info</button>
		<button onclick="external.invoke('warning')">Warning</button>
		<button onclick="external.invoke('error')">Error</button>
		<button onclick="external.invoke('changeTitle:'+document.getElementById('new-title').value)">
			Change title
		</button>
		<input id="new-title" type="text" />
		<button onclick="external.invoke('changeColor:'+document.getElementById('new-color').value)">
			Change color
		</button>
		<input id="new-color" value="#e91e63" type="color" />
	</body>
</html>
`),
		Resizable: true,
		Debug:     true,
		ExternalInvokeCallback: func(w webview.WebView, data string) {
			switch {
			case data == "close":
				w.Terminate()
			case data == "fullscreen":
				w.SetFullscreen(true)
			case data == "unfullscreen":
				w.SetFullscreen(false)
			case data == "open":
				log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
			case data == "opendir":
				log.Println("open", w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", ""))
			case data == "save":
				log.Println("save", w.Dialog(webview.DialogTypeSave, 0, "Save file", ""))
			case data == "message":
				w.Dialog(webview.DialogTypeAlert, 0, "Hello", "Hello, world!")
			case data == "info":
				w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Hello", "Hello, info!")
			case data == "warning":
				w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "Hello", "Hello, warning!")
			case data == "error":
				w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Hello", "Hello, error!")
			case strings.HasPrefix(data, "changeTitle:"):
				w.SetTitle(strings.TrimPrefix(data, "changeTitle:"))
			case strings.HasPrefix(data, "changeColor:"):
				hex := strings.TrimPrefix(strings.TrimPrefix(data, "changeColor:"), "#")
				num := len(hex) / 2
				if !(num == 3 || num == 4) {
					log.Println("Color must be RRGGBB or RRGGBBAA")
					return
				}
				i, err := strconv.ParseUint(hex, 16, 64)
				if err != nil {
					log.Println(err)
					return
				}
				if num == 3 {
					r := uint8((i >> 16) & 0xFF)
					g := uint8((i >> 8) & 0xFF)
					b := uint8(i & 0xFF)
					w.SetColor(r, g, b, 255)
					return
				}
				if num == 4 {
					r := uint8((i >> 24) & 0xFF)
					g := uint8((i >> 16) & 0xFF)
					b := uint8((i >> 8) & 0xFF)
					a := uint8(i & 0xFF)
					w.SetColor(r, g, b, a)
					return
				}
			}

			ok = data == "true"
			fmt.Println(w, data)
		},
	})
	w.SetColor(255, 255, 255, 255)

	fmt.Println(w.Loop(true))

	w.Dispatch(func() {
		w.Eval(`document.addEventListener("DOMContentLoaded",function () {
    console.log("DOMContentLoaded"+new Date())
});
        document.addEventListener("readystatechange",function () {
            console.log("B_____"+new Date());
            console.log(document.readyState)
//            switch (document.readyState){
//                case "loading":
//                    console.log("LOADING"+new Date());
//                    break;
//                case "loaded":
//                    console.log("loaded"+new Date());
//                    break;
//                case "interactive":
//                    console.log("interactive"+new Date());
//                    break;
//                case "complete":
//                    console.log("complete"+new Date());
//                    break;
//            }
        });
`)

	})

	go func() {

		for _, url := range []string{
			"https://mp.weixin.qq.com/s?__biz=MzA4NTIyOTI3NQ==&mid=2247483820&idx=1&sn=91505ad643049fa2d624d9efb2006923&chksm=9fda6a18a8ade30ea3e88d72de97cba5c521e87b28803dc510c512b86b32c164d4c2e37e726b&scene=21#wechat_redirect",
			"https://mp.weixin.qq.com/s/oE7Q-9wPrngHfr02gbBlLw",
			"https://mp.weixin.qq.com/s/dprkCOvPZHr6fi_qC91dVw",
			"https://mp.weixin.qq.com/s/CG8hu_1tgJwRl4jma24Avw",
			"https://www.yunyingpai.com/",
			"https://www.yunyingpai.com/data/560862.html",
			"https://www.yunyingpai.com/data/560863.html",
			"https://github.com/google/jsonapi",
		} {
			for i := 0; ; i += 100 {

				if ok {
					ok = false
					ok1 = false
					break
				}

				//if !ok1 {

				w.Dispatch(func() {
					fmt.Println(url, w.Eval(fmt.Sprintf(`
console.log(document.readyState);
if (document.readyState=="complete"){
	document.location.reload();

	// console.log(document.documentElement.outerHTML);
	// console.log(document.documentElement.outerHTML);	
	external.invoke(document.readyState);
	external.invoke(window.location.href);
	external.invoke(document.querySelector("p").text);
	external.invoke("true");
	window.location.href = "%s";
}`, url)))
				})
				time.Sleep(time.Millisecond * time.Duration(100+i))
			}
		}

	}()

	defer w.Exit()
	w.Run()
}
