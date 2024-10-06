package credits

import (
	"strings"

	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type controller struct {
}

func NewController() *controller {
	return &controller{}
}

func (c *controller) Init(ctx gscene.InitContext) {
	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Credits",
		Font: assets.Font2,
	}))

	panel := game.G.UI.NewPanel(eui.PanelConfig{
		MinWidth: 440,
	})
	root.AddChild(panel)

	panel.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: strings.Join([]string{
			"Thank you for playing NebuLeet Troopers!",
			"",
			"Made for a LDJAM 56 by @quasilyte in ~20 hours",
			"",
			"Written in Go, powered by Ebitengine",
		}, "\n"),
	}))

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text: "BACK",
		OnClick: func() {
			game.G.SceneManager.ChangeScene(game.G.NewMainMenuController())
		},
	}))

	game.G.UI.Build(ctx.Scene, root)
}

func (c *controller) Update(delta float64) {}
