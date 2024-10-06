package groundscape

import (
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/gsim"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type Controller struct {
	scene *gscene.Scene

	runner *gsim.Runner

	stage *gcombat.Stage
}

type ControllerConfig struct {
	Stage *gcombat.Stage
}

func NewController(config ControllerConfig) *Controller {
	return &Controller{
		stage: config.Stage,
	}
}

func (c *Controller) Init(ctx gscene.InitContext) {
	c.scene = ctx.Scene

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)
	ctx.Scene.AddGraphics(sceneutil.CombatMapSprite(game.G.State.CurrentStage.MapBg), 0)

	for _, team := range c.stage.Teams {
		for _, u := range team.Units {
			n := newUnitNode(u)
			c.scene.AddObject(n)
		}
	}

	c.runner = gsim.NewRunner(c.stage)
	c.runner.EventProjectileCreated.Connect(nil, func(p *gcombat.Projectile) {
		n := newProjectileNode(p)
		c.scene.AddObject(n)
	})
}

func (c *Controller) Update(delta float64) {
	c.runner.Update(delta)
}
