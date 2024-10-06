package groundscape

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/gsim"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/scenes/unitshop"
)

type Controller struct {
	scene *gscene.Scene

	back gscene.Controller

	runner *gsim.Runner

	stage *gcombat.Stage

	victory        bool
	statusLabel    *widget.Text
	team1cards     *widget.Container
	team2cards     *widget.Container
	continueButton *widget.Button
	cardsPanel     *widget.Container
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
		c.cardsPanel.GetWidget().Visibility = widget.Visibility_Hide
		c.continueButton.GetWidget().Visibility = widget.Visibility_Show
		c.victory = winner.Index == 0
		if c.victory {
			c.statusLabel.Label = "Status: victory"
		} else {
			c.statusLabel.Label = "Status: defeat"
		}
	})
	c.runner.EventUpdateCards.Connect(nil, c.updateCards)

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

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{
			MinHeight: 36,
		})
		c.cardsPanel = panel

		cols := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(8, 0),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, nil),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal: true,
				}),
			),
		)

		c.team1cards = widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
			)),
		)
		c.team2cards = widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
			)),
		)
		cols.AddChild(c.team1cards)
		cols.AddChild(c.team2cards)

		panel.AddChild(cols)

		rows.AddChild(panel)
	}

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
				game.G.SceneManager.ChangeScene(unitshop.NewController(c.back))
			} else {
				game.G.State.Retries++
				game.G.SceneManager.ChangeScene(c.back)
			}
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

func (c *Controller) updateCards(cards []gcombat.Card) {
	c.team1cards.RemoveChildren()
	c.team2cards.RemoveChildren()

	for _, card := range cards {
		info := card.Kind.Info()
		switch {
		case info.Category == gcombat.CardCategoryModifier:
			c.team1cards.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text: card.Kind.Info().Name,
				Font: assets.FontTiny,
			}))
			c.team2cards.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text: card.Kind.Info().Name,
				Font: assets.FontTiny,
			}))
		case card.TeamIndex == 0:
			c.team1cards.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text: card.Kind.Info().Name,
				Font: assets.FontTiny,
			}))
		case card.TeamIndex == 1:
			c.team2cards.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text: card.Kind.Info().Name,
				Font: assets.FontTiny,
			}))
		}
	}
}
