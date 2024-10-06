package groundscape

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/game"
)

type effectNode struct {
	sprite        *graphics.Sprite
	pos           gmath.Vec
	anim          Animation
	rotationSpeed gmath.Rad
	fades         bool
	noFlip        bool

	rotation gmath.Rad

	EventCompleted gsignal.Event[gsignal.Void]
}

type effectNodeConfig struct {
	Sprite *graphics.Sprite
	Pos    gmath.Vec
}

func newEffectNode(config effectNodeConfig) *effectNode {
	return &effectNode{
		sprite: config.Sprite,
		pos:    config.Pos,
	}
}

func (e *effectNode) Init(s *gscene.Scene) {
	e.sprite.Rotation = &e.rotation
	e.sprite.Pos.Base = &e.pos
	e.sprite.SetFrameWidth(e.sprite.ImageHeight())
	if !e.noFlip {
		e.sprite.SetHorizontalFlip(game.G.Rand.Bool())
		e.sprite.SetVerticalFlip(game.G.Rand.Bool())
	}
	if e.rotationSpeed != 0 {
		e.rotation = game.G.Rand.Rad()
	}
	s.AddGraphics(e.sprite, 0)

	e.anim.SetSprite(e.sprite, -1)
	e.anim.SetFPS(15)
}

func (e *effectNode) IsDisposed() bool {
	return e.sprite.IsDisposed()
}

func (e *effectNode) Dispose() {
	e.sprite.Dispose()
}

func (e *effectNode) Update(delta float64) {
	if e.fades && e.anim.IsLastFrame() {
		alphaDecrease := 8 * float32(delta)
		if e.sprite.GetAlpha() > alphaDecrease {
			e.sprite.SetAlpha(e.sprite.GetAlpha() - alphaDecrease)
		}
	}

	if e.anim.Tick(delta) {
		e.EventCompleted.Emit(gsignal.Void{})
		e.Dispose()
		return
	}

	if e.rotationSpeed != 0 {
		e.rotation += gmath.Rad(delta * float64(e.rotationSpeed))
	}
}
