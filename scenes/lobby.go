package scenes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/dat"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type lobbyController struct {
	level *gcombat.Level
}

func NewLobbyController() *lobbyController {
	return &lobbyController{}
}

func (c *lobbyController) Init(ctx gscene.InitContext) {
	root := eui.NewTopLevelRows()

	c.level = gcombat.LoadLevel(dat.LevelList[game.G.State.Level])

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: fmt.Sprintf("Level %d", game.G.State.Level+1),
		Font: assets.Font2,
	}))

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{
			MinWidth: 440,
		})

		rows := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(12),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal: true,
				}),
			),
		)
		panel.AddChild(rows)

		rows.AddChild(game.G.UI.NewText(eui.TextConfig{
			Text: fmt.Sprintf("Credits: %s\n", styles.Normal(strconv.Itoa(game.G.State.Credits))),
		}))

		cols := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Stretch: true,
				}),
			),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(12, 4),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, nil),
			)),
		)
		rows.AddChild(cols)

		{
			var unitSet [gcombat.NumUnitKinds]int
			for _, u := range game.G.State.Units {
				unitSet[u]++
			}
			var textLines []string
			textLines = append(textLines, styles.Normal("[Your troops]"))
			for i, count := range unitSet {
				kind := gcombat.UnitKind(i)
				if count == 0 {
					continue
				}
				textLines = append(textLines, fmt.Sprintf("%s: %s", kind.String(), styles.Normal(strconv.Itoa(count))))
			}

			cols.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:      strings.Join(textLines, "\n"),
				AlignLeft: true,
			}))
		}
		{
			var unitSet [gcombat.NumUnitKinds]int
			for _, u := range c.level.EnemyTroops {
				unitSet[u]++
			}
			var textLines []string
			textLines = append(textLines, styles.Normal("[Enemy troops]"))
			for i, count := range unitSet {
				kind := gcombat.UnitKind(i)
				if count == 0 {
					continue
				}
				textLines = append(textLines, fmt.Sprintf("%s: %s", kind.String(), styles.Normal(strconv.Itoa(count))))
			}

			cols.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:      strings.Join(textLines, "\n"),
				AlignLeft: true,
			}))
		}

		root.AddChild(panel)
	}

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text:     "READY",
		MinWidth: 160,
	}))

	game.G.UI.Build(ctx.Scene, root)
}

func (c *lobbyController) Update(delta float64) {}
