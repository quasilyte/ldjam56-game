package groundscape

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/gsim"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type Controller struct {
	scene *gscene.Scene

	back gscene.Controller

	runner *gsim.Runner

	stage *gcombat.Stage

	victory        bool
	statusLabel    *widget.Text
	continueButton *widget.Button
}

type ControllerConfig struct {
	Stage *gcombat.Stage
	Back  gscene.Controller
}

func NewController(config ControllerConfig) *Controller {
	return &Controller{
		stage: config.Stage,
		back:  config.Back,
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
	c.runner.EventFinished.Connect(nil, func(winner *gcombat.Team) {
		c.continueButton.GetWidget().Visibility = widget.Visibility_Show
		c.victory = winner.Index == 0
		if c.victory {
			c.statusLabel.Label = "Status: victory"
		} else {
			c.statusLabel.Label = "Status: defeat"
		}
	})

	c.initUI()
}

func (c *Controller) initUI() {
	stage := c.stage

	rows := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(8, 8),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: widget.Insets{
					Top: 32 + stage.MapBg.Bounds().Dy() + 16,
				},
			}),
		),
	)

	c.statusLabel = game.G.UI.NewText(eui.TextConfig{
		Text: "Status: turn 1",
	})
	rows.AddChild(c.statusLabel)

	c.continueButton = game.G.UI.NewButton(eui.ButtonConfig{
		Text: "CONTINUE",
		OnClick: func() {
			if c.victory {
				game.G.State.Casualties += c.stage.Teams[0].Casualties
				game.G.State.Level++
				game.G.State.EnterLevel()
				game.G.State.Credits += c.stage.Level.Reward
				survivors := game.G.State.Units[:0]
				for _, u := range c.stage.Teams[0].Units {
					if u.HP > 0 {
						survivors = append(survivors, u.Stats.Kind)
					}
				}
				game.G.State.Units = survivors
			} else {
				game.G.State.Retries++
			}
			game.G.SceneManager.ChangeScene(c.back)
		},
		MinWidth: 300,
	})
	c.continueButton.GetWidget().Visibility = widget.Visibility_Hide_Blocking
	rows.AddChild(c.continueButton)

	game.G.UI.Build(c.scene, rows)
}

func (c *Controller) Update(delta float64) {
	c.runner.Update(delta)
}
