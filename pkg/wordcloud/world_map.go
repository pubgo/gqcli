package wordcloud

import (
	"fmt"
	"strconv"
)

type WorldMap struct {
	Width           int
	Height          int
	CollisionMap    []int
	RealImageWidth  int
	RealImageHeight int
}

func (w *WorldMap) PrintMap() {
	for y := 0; y < w.Height; y++ {
		str := ""
		for x := 0; x < w.Width; x++ {
			idx := y*w.Width + x
			str = str + strconv.Itoa(w.CollisionMap[idx])
		}
		fmt.Println(str)
	}
}
