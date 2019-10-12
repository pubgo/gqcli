package wordcloud

const (
	IS_NOT_FIT = 1
	IS_FIT     = 2
	OUT_INDEX  = 3
	DEGREE_360 = 360
	DEGREE_180 = 180
	IS_EMPTY   = 0
	XUNIT      = 2
	YUNIT      = 2
)

type Position struct {
	Xpos   int
	Ypos   int
	Value  int
	XLeiji int
	YLeiji int
}

func NewPosition(xpos, ypos, value, xleiji, yleiji int) *Position {
	pos := &Position{
		Xpos:   xpos,
		Ypos:   ypos,
		Value:  value,
		XLeiji: xleiji,
		YLeiji: yleiji,
	}
	return pos
}

type Grid struct {
	Width     int
	Height    int
	positions []*Position
	XScale    int
	YScale    int
}

func (g *Grid) IsFit(xIncrement, yIncrement, width, height int, gridIntArray []int) int {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			index := i*g.Width + j
			position := g.positions[index]
			if position.Value != IS_EMPTY {
				position.Xpos = position.XLeiji + xIncrement
				position.Ypos = position.YLeiji + yIncrement
				if position.Xpos < 0 || position.Ypos < 0 || position.Xpos >= width || position.Ypos >= height {
					return OUT_INDEX
				}
				index = position.Ypos*width + position.Xpos
				if position.Value != 0 && gridIntArray[index] == position.Value {
					return IS_NOT_FIT
				}
			}
		}
	}
	return IS_FIT
}

func (g *Grid) SetCollisionMap(collisionMap []int, width, height int) {
	g.Width = width
	g.Height = height
	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := collisionMap[index]
			position := NewPosition(x, y, value, 0, 0)
			g.positions = append(g.positions, position)
			index++
		}
	}
}

func (g *Grid) Fill(gridIntArrayWidth, gridIntArrayHeight int, gridIntArray []int) {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			index := y*g.Width + x
			position := g.positions[index]
			index = position.Ypos*gridIntArrayWidth + position.Xpos
			if position.Value != IS_EMPTY {
				gridIntArray[index] = position.Value
			}
		}
	}
}
