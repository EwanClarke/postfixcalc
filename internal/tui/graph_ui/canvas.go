package graphui

import (
	"fmt"
	"math"
	"strings"

	"github.com/EwanClarke/postfixcalc/internal/engine"
)

type Size struct {
	width  int
	height int
}

type Bounds struct {
	xMin float64
	xMax float64
	yMin float64
	yMax float64
}

type Canvas struct {
	canvasSize Size
	bounds     Bounds
	engine     *engine.Engine
	postfix    []engine.Token
}

func NewCanvas() *Canvas {
	return &Canvas{
		canvasSize: Size{
			width:  1,
			height: 1,
		},
		bounds: Bounds{
			xMin: 0,
			xMax: 0,
			yMin: 0,
			yMax: 0,
		},
		engine:  engine.NewEngine(),
		postfix: []engine.Token{},
	}
}

func (c *Canvas) Resize(width int, height int) error {
	if width <= 0 {
		return fmt.Errorf("width %d", width)
	}
	if height <= 0 {
		return fmt.Errorf("height %d is less than minimum height 1", height)
	}

	c.canvasSize.width = width
	c.canvasSize.height = height

	return nil
}

func (c *Canvas) Size() (int, int) {
	return c.canvasSize.width, c.canvasSize.height
}

func (c *Canvas) SetBounds(xMin, xMax, yMin, yMax float64) {
	c.bounds = Bounds{
		xMin: xMin,
		xMax: xMax,
		yMin: yMin,
		yMax: yMax,
	}
}

func (c *Canvas) ChangeExpression(expression string) error {
	postfixTokens, err := c.engine.Parse(expression)
	if err != nil {
		return err
	}

	c.postfix = postfixTokens
	return nil
}

func (c *Canvas) Render() (string, error) {
	if len(c.postfix) == 0 {
		return "", nil
	}

	scaleX, scaleY, dotScaleY := c.calculateScales()
	plotPoints := c.evaluatePoints(scaleX, scaleY, dotScaleY)

	return c.toBraille(plotPoints), nil
}

func (c *Canvas) calculateScales() (scaleX, scaleY, dotScaleY float64) {
	scaleX = (c.bounds.xMax - c.bounds.xMin) / float64(c.canvasSize.width*2)
	scaleY = (c.bounds.yMax - c.bounds.yMin) / float64(c.canvasSize.height*4)
	dotScaleY = scaleY / 4
	return
}

func (c *Canvas) evaluatePoints(scaleX, scaleY, dotScaleY float64) map[int]map[int]bool {
	plotPoints := make(map[int]map[int]bool)

	var prevX, prevYIdx int
	var hasPrev bool

	for charX := 0; charX < c.canvasSize.width*2; charX++ {
		xBase := c.bounds.xMin + float64(charX)*scaleX

		y, err := c.engine.EvaluateAt(c.postfix, xBase)
		if err != nil {
			continue
		}

		yIdx := c.canvasSize.height*4 - 1 - int(math.Round(y*scaleY/dotScaleY))

		if hasPrev {
			c.plotLine(prevX, prevYIdx, charX, yIdx, plotPoints)
		}

		prevX, prevYIdx = charX, yIdx
		hasPrev = true
	}

	return plotPoints
}

func (c *Canvas) toBraille(plotPoints map[int]map[int]bool) string {
	var result strings.Builder

	for row := 0; row < c.canvasSize.height; row++ {
		result.WriteString(c.renderRow(row, plotPoints))
		if row < c.canvasSize.height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

func (c *Canvas) renderRow(row int, plotPoints map[int]map[int]bool) string {
	var rowResult strings.Builder
	baseY := row * 4

	for charX := 0; charX < c.canvasSize.width; charX++ {
		dots := c.dotsForCell(charX, baseY, plotPoints)
		rowResult.WriteRune(dotsToBraille(dots))
	}

	return rowResult.String()
}

func (c *Canvas) dotsForCell(charX, baseY int, plotPoints map[int]map[int]bool) []bool {
	dots := make([]bool, 8)

	leftCol := charX * 2
	rightCol := leftCol + 1

	if col, ok := plotPoints[leftCol]; ok {
		dots[0] = col[baseY]
		dots[1] = col[baseY+1]
		dots[2] = col[baseY+2]
		dots[6] = col[baseY+3]
	}
	if col, ok := plotPoints[rightCol]; ok {
		dots[3] = col[baseY]
		dots[4] = col[baseY+1]
		dots[5] = col[baseY+2]
		dots[7] = col[baseY+3]
	}

	return dots
}

func dotsToBraille(dots []bool) rune {
	var result rune = 0x2800

	for i, on := range dots {
		if on {
			result |= 1 << i
		}
	}
	return result
}

func (c *Canvas) plotLine(x1, y1, x2, y2 int, plotPoints map[int]map[int]bool) {
	Bresenham(x1, y1, x2, y2, func(px, py int) {
		if plotPoints[px] == nil {
			plotPoints[px] = make(map[int]bool)
		}
		plotPoints[px][py] = true
	})
}

type PlotFunc func(x, y int)

func Bresenham(x1, y1, x2, y2 int, plot PlotFunc) {
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx + dy

	for {
		plot(x1, y1)
		if x1 == x2 && y1 == y2 {
			return
		}
		e2 := 2 * err
		if e2 > dy {
			err += dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
