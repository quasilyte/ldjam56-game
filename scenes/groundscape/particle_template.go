package groundscape

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/ebitengine-graphics/particle"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ldjam56-game/styles"
)

var (
	particleTemplateMissileTrail *particle.Template
)

func InitParticleTemplates() {
	offsets := [...]gmath.Vec{
		{X: -1, Y: -1},
		{X: 0, Y: -1},
		{X: +1, Y: -1},

		{X: -1, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: +1, Y: 0},

		{X: -1, Y: +1},
		{X: 0, Y: +1},
		{X: +1, Y: +1},
	}

	frame3x2 := ebiten.NewImage(3, 2)
	frame3x2.Fill(color.White)

	frame1x1 := ebiten.NewImage(1, 1)
	frame1x1.Fill(color.White)

	{
		tmpl := particle.NewTemplate()

		clrFrom := graphics.ColorScale{R: 1, G: 1, B: 1, A: 1.0}
		clrTo := graphics.ColorScale{R: 1, G: 1, B: 1, A: 0}

		tmpl.SetImage(frame1x1)
		tmpl.SetUpdateColorScaleFunc(func(ctx particle.UpdateContext) graphics.ColorScale {
			t := ctx.Time()
			return clrFrom.Lerp(clrTo, t)
		})
		tmpl.SetSpawnOffsetFunc(func(ctx particle.SpawnContext) gmath.Vec {
			rnd := ctx.RandUint()
			return offsets[rnd%uint64(len(offsets))]
		})

		tmpl.SetParticleSpeedRange(10, 40)
		tmpl.SetEmitInterval(0.015)
		tmpl.SetEmitBurst(1, 2)
		tmpl.SetParticleLifetimeRange(0.6, 0.9)
		tmpl.SetParticleDirection(math.Pi, 0.09)

		particleTemplateMissileTrail = tmpl.Clone()
		particleTemplateMissileTrail.SetPalette([]graphics.ColorScale{
			styles.ColorOrange,
		})
	}
}
