package unitshop

import (
	"fmt"
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gslices"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type Controller struct {
	back gscene.Controller

	creditsCounter *widget.Text

	buttons []*buyButton
}

func NewController(back gscene.Controller) *Controller {
	return &Controller{
		back: back,
	}
}

type buyButton struct {
	btn  *widget.Button
	unit gcombat.UnitKind
}

func (c *Controller) Init(ctx gscene.InitContext) {
	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root := eui.NewTopLevelRows()

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Reinforcements Menu",
		Font: assets.Font2,
	}))

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{
			MinWidth: 440,
		})

		c.creditsCounter = game.G.UI.NewText(eui.TextConfig{
			AlignLeft:   true,
			ForceBBCode: true,
		})
		c.updateCreditsCounter()
		panel.AddChild(c.creditsCounter)

		root.AddChild(panel)
	}

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{
			MinWidth: 440,
		})

		table := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Spacing(12, 2),
			)),
		)
		panel.AddChild(table)

		for k, unlocked := range game.G.State.UnitsUnlocked {
			if !unlocked {
				continue
			}
			k := gcombat.UnitKind(k)
			stats := gcombat.GetUnitStats(k)
			btn := game.G.UI.NewTileButton(eui.TileButtonConfig{
				Text: k.String(),
				OnClick: func() {
					game.G.State.Credits -= stats.Cost
					game.G.State.Units = append(game.G.State.Units, k)
					c.updateCreditsCounter()
					for _, b := range c.buttons {
						stats := gcombat.GetUnitStats(b.unit)
						b.btn.GetWidget().Disabled = stats.Cost > game.G.State.Credits
					}
				},
				MinWidth: 108,
			})
			btn.GetWidget().Disabled = stats.Cost > game.G.State.Credits
			table.AddChild(btn)
			table.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:      styles.Normal(strconv.Itoa(stats.Cost)) + "$",
				Font:      assets.FontTiny,
				AlignLeft: true,
			}))
			table.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:      stats.Comment,
				Font:      assets.FontTiny,
				AlignLeft: true,
			}))
			c.buttons = append(c.buttons, &buyButton{
				btn:  btn,
				unit: k,
			})
		}

		root.AddChild(panel)
	}

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text:     "CONTINUE",
		MinWidth: 160,
		OnClick: func() {
			gslices.SortFunc(game.G.State.Units, func(u1, u2 gcombat.UnitKind) bool {
				return u1 < u2
			})
			game.G.SceneManager.ChangeScene(c.back)
		},
	}))

	game.G.UI.Build(ctx.Scene, root)
}

func (c *Controller) Update(delta float64) {

}

func (c *Controller) updateCreditsCounter() {
	s := fmt.Sprintf("Credits: %s$", styles.Normal(strconv.Itoa(game.G.State.Credits)))
	c.creditsCounter.Label = s
}
