package scenes

import (
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
)

type settingsController struct {
}

func NewSettingsController() *settingsController {
	return &settingsController{}
}

func (c *settingsController) Init(ctx gscene.InitContext) {
	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Settings",
		Font: assets.Font2,
	}))

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	{
		cols := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Spacing(8, 0),
				widget.GridLayoutOpts.Stretch([]bool{false, true, false}, nil),
			)),
		)
		volumeLabel := game.G.UI.NewText(eui.TextConfig{
			Text: "Sound Level: " + strconv.Itoa(game.G.SoundVolume),
		})

		cols.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
			Text: "-",
			OnClick: func() {
				game.G.SoundVolume = gmath.ClampMin(game.G.SoundVolume-1, 0)
				volumeLabel.Label = "Sound Level: " + strconv.Itoa(game.G.SoundVolume)
				game.G.Audio.SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(game.G.SoundVolume))
			},
		}))
		cols.AddChild(volumeLabel)
		cols.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
			Text: "+",
			OnClick: func() {
				game.G.SoundVolume = gmath.ClampMax(game.G.SoundVolume+1, 5)
				volumeLabel.Label = "Sound Level: " + strconv.Itoa(game.G.SoundVolume)
				game.G.Audio.SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(game.G.SoundVolume))
			},
		}))

		root.AddChild(cols)
	}

	root.AddChild(game.G.UI.NewText(eui.TextConfig{Text: ""}))

	root.AddChild(game.G.UI.NewButton(eui.ButtonConfig{
		Text: "BACK",
		OnClick: func() {
			game.G.SceneManager.ChangeScene(NewMainMenuController())
		},
	}))

	game.G.UI.Build(ctx.Scene, root)
}

func (c *settingsController) Update(delta float64) {}

// rowContainer.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
// 	Resources:  c.state.UIResources,
// 	Input:      c.state.Input,
// 	Value:      &c.state.Settings.SoundLevel,
// 	Label:      "Effects volume",
// 	ValueNames: []string{"off", "1", "2", "3", "4", "5"},
// 	OnPressed: func() {
// 		if c.state.Settings.SoundLevel != 0 {
// 			scene.Audio().SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(c.state.Settings.SoundLevel))
// 			scene.Audio().PlaySound(assets.AudioPhotonCannon1)
// 		}
// 	},
// }))
