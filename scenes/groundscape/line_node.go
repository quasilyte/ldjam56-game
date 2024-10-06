package groundscape

import (
	"image/color"

	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
)

type lineNode struct {
	from gmath.Vec
	to   gmath.Vec

	color      color.NRGBA
	width      float64
	line       *graphics.Line
	opaqueTime float64
}

func newLineNode(from, to gmath.Vec, clr color.NRGBA) *lineNode {
	return &lineNode{
		from:  from,
		to:    to,
		width: 1,
		color: clr,
	}
}

func (l *lineNode) Init(scene *gscene.Scene) {
	l.line = graphics.NewLine(gmath.Pos{Offset: l.from}, gmath.Pos{Offset: l.to})
	var c graphics.ColorScale
	c.SetColor(l.color)
	l.line.SetColorScale(c)
	l.line.SetWidth(1)
	scene.AddGraphics(l.line, 0)
}

func (l *lineNode) IsDisposed() bool {
	return l.line.IsDisposed()
}

func (l *lineNode) Update(delta float64) {
	if l.opaqueTime > 0 {
		l.opaqueTime -= delta
		return
	}

	if l.line.GetAlpha() < 0.1 {
		l.line.Dispose()
		return
	}
	l.line.SetAlpha(l.line.GetAlpha() - float32(delta*4))
}
