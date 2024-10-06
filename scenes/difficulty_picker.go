package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type difficultyPickerController struct {
}

func NewDifficultyPickerController() *difficultyPickerController {
	return &difficultyPickerController{}
}

func (c *difficultyPickerController) Init(ctx gscene.InitContext) {
	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Choose Difficulty",
		Font: assets.Font2,
	}))

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text:     "HARD",
		MinWidth: 440,
		OnClick: func() {
			c.start(false)
		},
	}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text:     "NORMAL",
		MinWidth: 440,
		OnClick: func() {
			c.start(true)
		},
	}))

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{})
		root.AddChild(panel)

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

		lines := []string{
			"Choose " + styles.Orange("normal") + " difficulty if you don't want to struggle",
		}
		for _, l := range lines {
			rows.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:       l,
				Font:       assets.FontTiny,
				LayoutData: widget.RowLayoutData{Stretch: true},
			}))
		}
		panel.AddChild(rows)
	}

	game.G.UI.Build(ctx.Scene, root)
}

func (c *difficultyPickerController) Update(delta float64) {}

func (c *difficultyPickerController) start(easyMode bool) {
	units := []gcombat.UnitKind{
		gcombat.UnitRifle,
		gcombat.UnitRifle,
		gcombat.UnitRifle,
		gcombat.UnitRifle,
	}
	if easyMode {
		units = append(units, gcombat.UnitRifle)
	}
	units = append(units, gcombat.UnitLaser)

	game.G.State = &game.State{
		Level:   0,
		Credits: 0,
		Easy:    easyMode,
		Units:   units,
		UnitsUnlocked: [gcombat.NumUnitKinds]bool{
			gcombat.UnitRifle: true,
		},
		CardsUnlocked: make(map[gcombat.CardKind]struct{}),
	}
	game.G.State.EnterLevel()
	game.G.SceneManager.ChangeScene(NewLobbyController())
}
