package scenes

import (
	"os"
	"runtime"

	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/credits"
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
		MinWidth: 440,
		OnClick: func() {
			game.G.State = &game.State{
				Level:   0,
				Credits: 0,
				Units: []gcombat.UnitKind{
					gcombat.UnitRifle,
					gcombat.UnitRifle,
					gcombat.UnitRifle,
					gcombat.UnitRifle,
				},
				UnitsUnlocked: [gcombat.NumUnitKinds]bool{
					gcombat.UnitRifle: true,
				},
				CardsUnlocked: make(map[gcombat.CardKind]struct{}),
			}
			game.G.State.EnterLevel()
			game.G.SceneManager.ChangeScene(NewLobbyController())
		},
	}))

	{
		settings := game.G.UI.NewButton(eui.ButtonConfig{
			Text: "SETTINGS",
			OnClick: func() {
				game.G.SceneManager.ChangeScene(NewSettingsController())
			},
		})
		root.AddChild(settings)
	}

	{
		settings := game.G.UI.NewButton(eui.ButtonConfig{
			Text: "CREDITS",
			OnClick: func() {
				game.G.SceneManager.ChangeScene(credits.NewController())
			},
		})
		root.AddChild(settings)
	}

	if runtime.GOARCH != "wasm" {
		root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

		root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
			Text: "EXIT",
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
