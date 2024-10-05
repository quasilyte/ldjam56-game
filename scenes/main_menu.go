package scenes

import (
	"os"
	"runtime"

	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type mainMenuController struct {
}

func NewMainMenuController() *mainMenuController {
	return &mainMenuController{}
}

func (c *mainMenuController) Init(ctx gscene.InitContext) {
	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "NebuLeet Troopers",
		Font: assets.Font3,
	}))

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text:     "START",
		MinWidth: 160,
	}))

	{
		settings := game.G.UI.NewButton(eui.ButtonConfig{
			Text:     "SETTINGS",
			MinWidth: 160,
		})
		settings.GetWidget().Disabled = true
		root.AddChild(settings)
	}

	{
		settings := game.G.UI.NewButton(eui.ButtonConfig{
			Text:     "CREDITS",
			MinWidth: 160,
		})
		settings.GetWidget().Disabled = true
		root.AddChild(settings)
	}

	if runtime.GOARCH != "wasm" {
		root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

		root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
			Text:     "EXIT",
			MinWidth: 160,
			OnClick: func() {
				os.Exit(0)
			},
		}))
	}

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "LDJAM 56 compo build 1",
		Font: assets.FontTiny,
	}))
	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Made with Ebitengine",
		Font: assets.FontTiny,
	}))

	game.G.UI.Build(ctx.Scene, root)
}

func (c *mainMenuController) Update(delta float64) {}
