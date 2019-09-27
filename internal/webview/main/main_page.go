package main

import (
	"fmt"
	"github.com/zserge/webview"
	"net/url"
)

func main() {

	w := webview.New(webview.Settings{
		Width:  1200,
		Height: 600,
		Title:  "Loaded: Injected via JavaScript",
		URL: `data:text/html,` + url.PathEscape(`
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>

<body>
<h1>命运与波动</h1>
<p>我还是想说交流电的事情: 50Hz 220V.</p>
<p>波粒二象性: </p>
<ol start="1">
  <li>
    <p>当我们以超过 20ms 的观察手段去观察的时候, 每次都能拿到一根正弦波. 平均值是 220V, 我们认为这是稳定可靠的结论.</p>
  </li>
</ol>
<ol start="2">
  <li>
    <p>随着观察尺度越来越小, 当我们以 1ms 手段去观察的时候, 每次得到的都是一个 0 - 311V 的随机数.</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>观察手段本身包含着时间, 空间尺度. 这个观察方式本身会与事物本质发生影响.</p>
  </li>
</ol>
<ol start="2" data-lake-indent="1">
  <li>
    <p>于是我们通过统计学, 观察 10000 次, 能得出一个概率波函数.</p>
  </li>
</ol>
<ol start="3" data-lake-indent="1">
  <li>
    <p><img src="https://cdn.nlark.com/yuque/0/2018/png/85323/1539908273857-31d64fc0-7417-4992-9007-de0e593686d4.png#width=300" style="max-width: 600px; width: 300px;" /></p>
  </li>
</ol>
<p><br /></p>
<p>量子共振: 我们虽然无法确定一个量子的位置, 但是我们可以通过一定<span style="color: #333333;"><span style="background-color: #FFFFFF;">频率的电磁波</span></span>, 去提高量子出现在某一位置的概率.</p>
<p><br /></p>
<p>命运论: </p>
<ol start="1">
  <li>
    <p>每个人的未来有无数可能, 可以类似总结出一个概率波函数.</p>
  </li>
</ol>
<ol start="2">
  <li>
    <p>我们虽然无法确定未来会出现在哪个位置. 但是我们可以通过一个波形, 去提高未来出现在某个位置的概率.</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>这个引导波: 就是你的意识.</p>
  </li>
</ol>
<ol start="3">
  <li>
    <p>so: 你所相信的, 就是你的命运!</p>
  </li>
</ol>
<p><br /></p>
<p>命运叠加平衡:</p>
<p>虽然有无数种可能. 但是在几个大因素下维持相对平衡:</p>
<ol start="1">
  <li>
    <p>环境</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>社会环境, 家庭环境, 职业环境</p>
  </li>
</ol>
<ol start="2">
  <li>
    <p>性格</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>喜欢什么, 讨厌什么.</p>
  </li>
</ol>
<ol start="3">
  <li>
    <p>意志</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>相信什么, 不相信什么.</p>
  </li>
</ol>
<ol start="4">
  <li>
    <p>选择</p>
  </li>
</ol>
<ol start="1" data-lake-indent="1">
  <li>
    <p>靠近什么, 远离什么.</p>
  </li>
</ol>
<p><br /></p>
<p>方法论:</p>
<ol start="1">
  <li>
    <p>当你的观察尺度小于事物本身规律的时候. 得出的结论是片面的, 需要多观察几次. 学会用统计学的方式获取真正的本质规律.</p>
  </li>
</ol>
<ol start="2">
  <li>
    <p>选择你相信的, 坚信的, 靠近它, 与之共振. 否则人生处处都是干涉波.</p>
  </li>
</ol>
<p><br /></p>
<p><br /></p>
<p><br /></p>
</body>
`),
		Resizable: true,
		Debug:     true,
		ExternalInvokeCallback: func(w webview.WebView, data string) {
			fmt.Println(w, data)
		},
	})
	w.Dispatch(func() {
		//w.Eval(`document.body.innerHTML = "<h1>Hello, world</h1>";`)
	})

	defer w.Exit()
	w.Run()
}
