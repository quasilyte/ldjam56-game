package game

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/eui"
)

var G *GlobalContext

type GlobalContext struct {
	SceneManager *gscene.Manager

	State *State

	WindowSize gmath.Vec

	Loader *resource.Loader

	Rand gmath.Rand

	UI *eui.Builder
}

func (ctx *GlobalContext) NewSprite(id resource.ImageID) *graphics.Sprite {
	s := graphics.NewSprite()
	img := ctx.Loader.LoadImage(id)
	s.SetImage(img.Data)
	return s
}
