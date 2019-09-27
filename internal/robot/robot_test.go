package robot

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
	"github.com/pubgo/errors"
	"log"
	"testing"
	"time"
)

// 组件对象， 对象的坐标
// 获取组件的对象和坐标
// 关联鼠标事件和键盘事件等
// 下拉截止判断 通过获取当前鼠标的横竖区域的颜色，然后得到一个hash值，如果该hash值在一个范围内不变，则最后了

type Point struct {
	X float64
	Y float64
}

func point(x, y float64) Point {
	return Point{X: x, Y: y}
}

const activeName = "WeChat"

// 屏幕大小
var screenPoint = point(1440, 900)

// 最大化, 最小化
// robotgo.KeyTap("f", "lctrl", "lcmd")

// 公众号回退按钮
var gzhBackPoint = point(25, 50)

// 公众号前进按钮
var gzhGoPoint = point(65, 50)

// 公众号文章列表间隔
const gzhListStep = 105

// 微信公众号第一篇文章的位置
var gzhFirstList = point(715, 385)

// 聊天列表 187,95
var contactListPoint = point(187, 95)

// 聚焦通信列表 35 167
var contactListFocusPoint = point(35, 167)

// 公众号关闭 175 189
var gzhgbanPoint = point(175, 189)

// 公众号 第一篇 位置 597 293 竖列+10
var gzhgbanPoint1 = point(175, 189)

// 文章回退 185 223
// 文章刷新 253 222
// 复制网址 1006 223
// 文章正上点击 621 291
// data_ok
// url_ok
// 公众号返回 391 28
// 公众号列表1409 33 --  1362 290

func move() {
	robotgo.Move(100, 200)

	// move the mouse to 100, 200
	robotgo.MoveMouse(100, 200)

	robotgo.Drag(10, 10)
	robotgo.Drag(20, 20, "right")
	//
	robotgo.DragSmooth(10, 10)
	robotgo.DragSmooth(100, 200, 1.0, 100.0)

	// smooth move the mouse to 100, 200
	robotgo.MoveSmooth(100, 200)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)

	errors.Panic(clipboard.WriteAll("日本語"))
	text, err := clipboard.ReadAll()
	if err != nil {
		log.Println("clipboard read all error: ", err)
	} else {
		if text != "" {
			log.Println("text is: ", text)
		}
	}

	for i := 0; i < 1080; i += 1000 {
		fmt.Println(i)
		robotgo.MoveMouse(800, i)
	}
}

func click() {

	// click the left mouse button
	robotgo.Click()

	// click the right mouse button
	robotgo.Click("right", false)

	// double click the left mouse button
	robotgo.MouseClick("left", true)
}

func toggleAndScroll() {
	// scrolls the mouse either up
	robotgo.ScrollMouse(100, "down")

	robotgo.Scroll(100, 200)

	// toggles right mouse button
	//robotgo.MouseToggle("down", "right")

	//robotgo.MouseToggle("up")
}

func get() {
	// gets the mouse coordinates
	x, y := robotgo.GetMousePos()
	fmt.Println("pos:", x, y)
	if x == 456 && y == 586 {
		fmt.Println("mouse...", "586")
	}

	robotgo.MoveMouse(x, y)
}

func toggleAndScroll1() {
	// scrolls the mouse either up
	robotgo.ScrollMouse(gzhListStep, "up")
	//robotgo.Scroll(100, 200)

	// toggles right mouse button
	robotgo.MouseToggle("down", "right")

	//robotgo.MouseToggle("up")
}

func mouse() {
	fmt.Println("start")
	move()
	fmt.Println("mouse")

	click()

	fmt.Println("click")
	get()

	fmt.Println("get")

	toggleAndScroll()
	fmt.Println("toggleAndScroll")
}

func TestA1(t *testing.T) {
	//nps, err := robotgo.Process()
	//errors.Panic(err)
	//for _, p := range nps {
	//	fmt.Println(p)
	//}

	title := robotgo.GetTitle()
	fmt.Println("title@@@ ", title)

	errors.Panic(robotgo.ActiveName(activeName))

	//errors.Panic(clipboard.WriteAll("日本語"))
	//text, err := clipboard.ReadAll()
	//if err != nil {
	//	log.Println("clipboard read all error: ", err)
	//} else {
	//	if text != "" {
	//		log.Println("text is: ", text)
	//	}
	//}

	//mouse()

	//_t := point(537, 411)

	//fmt.Println(robotgo.KeyTap("f", "lctrl", "lcmd"))

	robotgo.MoveMouseSmooth(int(gzhFirstList.X), int(gzhFirstList.Y))

	_ii := 0

	for {

		if _ii == 0 {
			_ii = 4
			//x, y := robotgo.GetMousePos()
			//robotgo.MoveMouseSmooth(x, y+_ii*105)
		}

		//robotgo.ScrollMouse(10, "down")
		//robotgo.SetMouseDelay(1000)
		//robotgo.Scroll(0, gzhListStep)

		// 获取当前的坐标
		x, y := robotgo.GetMousePos()
		fmt.Println(x, y)

		//robotgo.MoveMouse(int(gzhFirstList.X), int(gzhFirstList.Y)+105)
		//robotgo.MouseToggle("down")
		//robotgo.DragMouse(int(gzhFirstList.X), int(gzhFirstList.Y))
		//robotgo.MouseToggle("up")

		toggleAndScroll()

		time.Sleep(time.Second)
		//fmt.Println(robotgo.IsValid())
		//fmt.Println(robotgo.GetActive())
		//fmt.Println(robotgo.FindNames())
		_ii -= 1
	}

	//s := robotgo.Start()
	//defer robotgo.End()
	//for e := range s {
	//errors.Panic(clipboard.WriteAll(e.String()))
	//text, err := clipboard.ReadAll()
	//if err != nil {
	//	log.Println("clipboard read all error: ", err)
	//} else {
	//	if text != "" {
	//cs, _ := charsetutil.GuessString(text)
	//text = charsetutil.MustDecodeString(text, cs.Charset())
	//fmt.Println("text is: ", text)
	//}
	//}

	// color := robotgo.GetMouseColor()
	// fmt.Println("color---- ", color)

	//

	// bitmap := robotgo.CaptureScreen(10, 20, 1500, 30)
	// defer robotgo.FreeBitmap(bitmap)

	// s := sha1.New()
	// s.Write([]byte(robotgo.TostringBitmap(bitmap)))
	// fmt.Println(hex.EncodeToString(s.Sum(nil)))

	//robotgo.CountColor()
	//robotgo.GoCaptureScreen()

	//robotgo.ShowAlert("dd","ggg")

	//robotgo.ScrollMouse(10, "down")

	//fmt.Println("hook: ", e.String())

	//}
}
