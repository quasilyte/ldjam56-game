package groundscape

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type Controller struct {
	scene *gscene.Scene

	level *gcombat.Level

	worldOffset gmath.Vec
}

type ControllerConfig struct {
	Level *gcombat.Level
}

func NewController(config ControllerConfig) *Controller {
	return &Controller{
		level: config.Level,
	}
}

func (c *Controller) Init(ctx gscene.InitContext) {
	c.scene = ctx.Scene

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)
	ctx.Scene.AddGraphics(sceneutil.CombatMapSprite(game.G.State.CurrentStage.MapBg), 0)

	{
		s := game.G.NewSprite(assets.ImageUnitTank)
		s.Pos.Offset = (gmath.Vec{X: 32, Y: 32}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
	{
		s := game.G.NewSprite(assets.ImageUnitRifle)
		s.Pos.Offset = (gmath.Vec{X: 20, Y: 20}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
	{
		s := game.G.NewSprite(assets.ImageUnitLaser)
		s.Pos.Offset = (gmath.Vec{X: 40, Y: 20}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
}

func (c *Controller) Update(delta float64) {}
