package webview

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pubgo/g/errors"
	"github.com/zserge/webview"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Webview struct {
	home string
	url  string
	w    webview.WebView
}

func (Webview) init() {

}

func (t *Webview) Back() {
	errors.Panic(t.w.Eval("window.history.back();"))
}

func (t *Webview) Forward() {
	errors.Panic(t.w.Eval("window.history.forward();"))
}

func (t *Webview) Home() {
	t.Open(t.home)
}

func (t *Webview) Open(url string) {
	t.url = url
	errors.Panic(t.w.Eval(fmt.Sprintf(`window.location.pathname = "%s";`, url)))
}

func (t *Webview) Reload() {
	errors.Panic(t.w.Eval(fmt.Sprintf("window.location.reload();")))
}

func (t *Webview) HandleRPC(w webview.WebView, data string) {
	args := strings.Split(data, ":")

	switch args[0] {
	case "back":
		t.Back()
	case "forwrd":
		t.Forward()
	case "reload":
		t.Reload()
	case "home":
		t.Home()
	case "open":
		t.Open(args[1])
	}

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
}

func md5sum(in string) string {
	hasher := md5.New()
	hasher.Write([]byte(in))
	return hex.EncodeToString(hasher.Sum(nil))
}

func slurpFile(fname string) string {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("Error in slurpFile: %s", err)
	}

	return string(content)
}

func slurpURL(url string) string {
	res, err := http.Head(url)
	if err != nil {
		log.Fatalf("Error executing http head: %s", err)
	}

	etag := res.Header.Get("Etag")
	tmpFname := fmt.Sprintf("/tmp/webview-cache-%s", md5sum(etag))

	if _, err := os.Stat(tmpFname); !os.IsNotExist(err) {
		log.Println("Reading response for %s from cache at %s", url, tmpFname)
		return slurpFile(tmpFname)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error executing http get: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading resp body: %s", err)
	}

	_, err = os.Create(tmpFname)
	if err != nil {
		log.Fatalf("Error creating cache file: %s", err)
	}

	err = ioutil.WriteFile(tmpFname, body, 0644)
	if err != nil {
		log.Fatalf("Error writing to cache file: %s", err)
	}

	return string(body)
}

func slurp(resource string) string {
	if _, err := os.Stat(resource); os.IsNotExist(err) {
		return slurpURL(resource)
	}

	return slurpFile(resource)
}

func setBodyHTML(content string) string {
	return fmt.Sprintf("(function(content){ document.body.innerHTML = content; }(`%s`))", content)
}
