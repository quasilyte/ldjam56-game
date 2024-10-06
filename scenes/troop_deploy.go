package scenes

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	graphics "github.com/quasilyte/ebitengine-graphics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/groundscape"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type troopDeployController struct {
	deployHint            *widget.Text
	troopsDeployedCounter *widget.Text
	deployed              int

	bgSprite *graphics.Sprite

	startButton *widget.Button
	tileButtons []*widget.Button

	scene *gscene.Scene
}

func NewTroopDelpoyController() *troopDeployController {
	return &troopDeployController{}
}

func (c *troopDeployController) Init(ctx gscene.InitContext) {
	c.scene = ctx.Scene

	stage := game.G.State.CurrentStage

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	c.bgSprite = sceneutil.CombatMapSprite(stage.MapBg)
	ctx.Scene.AddGraphics(c.bgSprite, 0)

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
					// Top: 32 + stage.MapBg.Bounds().Dy() + 16,
					Top: 32,
				},
			}),
		),
	)

	tileButtonsGrid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(len(stage.Level.Tiles[0])),
		)),
	)
	for row := 0; row < len(stage.Level.Tiles); row++ {
		for col := 0; col < len(stage.Level.Tiles[0]); col++ {
			btn := game.G.UI.NewTileButton(eui.TileButtonConfig{
				MinWidth:  64,
				MinHeight: 64,
				OnClick: func() {
					c.onUnitDeployed(row, col)
				},
			})
			terrain := stage.Level.Tiles[row][col]
			btn.GetWidget().Disabled = col >= stage.Level.DeployWidth ||
				terrain == gcombat.TileMountains
			if !btn.GetWidget().Disabled {
				c.tileButtons = append(c.tileButtons, btn)
			}
			tileButtonsGrid.AddChild(btn)
		}
	}
	rows.AddChild(tileButtonsGrid)

	c.troopsDeployedCounter = game.G.UI.NewText(eui.TextConfig{
		Font: assets.FontTiny,
	})
	c.updateTroopsDeployedCounter()
	rows.AddChild(c.troopsDeployedCounter)

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{})
		c.deployHint = game.G.UI.NewText(eui.TextConfig{
			MinWidth:    220,
			Font:        assets.FontTiny,
			ForceBBCode: true,
			LayoutData: widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			},
		})
		c.updateCurrentlyDeploying()
		panel.AddChild(c.deployHint)
		rows.AddChild(panel)
	}

	c.startButton = game.G.UI.NewButton(eui.ButtonConfig{
		Text: "FIGHT",
		OnClick: func() {
			game.G.SceneManager.ChangeScene(groundscape.NewController(groundscape.ControllerConfig{
				Stage: game.G.State.CurrentStage,
				Back:  NewLobbyController(),
			}))
		},
	})
	c.startButton.GetWidget().Disabled = true
	rows.AddChild(c.startButton)

	c.deployEnemyTroops()

	game.G.UI.Build(ctx.Scene, rows)
}

func (c *troopDeployController) Update(delta float64) {}

func (c *troopDeployController) deployEnemyTroops() {
	stage := game.G.State.CurrentStage
	for _, u := range stage.Teams[1].Units {
		var rowcol [2]int
		switch {
		case u.Stats.Infantry:
			rowcol = gmath.RandElem(&game.G.Rand, stage.Level.EnemyInfantrySpots)
		default:
			rowcol = gmath.RandElem(&game.G.Rand, stage.Level.EnemyVehicleSpots)
		}
		c.deployUnit(u, rowcol[0], rowcol[1])
	}
}

func (c *troopDeployController) deployUnit(u *gcombat.Unit, row, col int) {
	pos := gmath.Vec{
		X: float64(col*64) + 32,
		Y: float64(row*64) + 32,
	}
	pos = pos.Add(game.G.Rand.Offset(-16, 16))
	u.Pos = pos
	u.SpawnPos = pos

	sprite := game.G.NewSprite(u.Stats.Image)
	sprite.Pos.Offset = c.bgSprite.Pos.Offset.Add(u.Pos)
	if u.Team != game.G.State.CurrentStage.Teams[0] {
		sprite.SetHorizontalFlip(true)
	}
	c.scene.AddGraphics(sprite, 0)
}

func (c *troopDeployController) onUnitDeployed(row, col int) {
	u := game.G.State.CurrentStage.Teams[0].Units[c.deployed]

	c.deployUnit(u, row, col)

	c.deployed++
	c.updateTroopsDeployedCounter()

	if c.deployed == len(game.G.State.Units) {
		c.startButton.GetWidget().Disabled = false
		for _, b := range c.tileButtons {
			b.GetWidget().Disabled = true
		}
		c.deployHint.Label = styles.Normal("All troops are deployed!")
	} else {
		c.updateCurrentlyDeploying()
	}
}

func (c *troopDeployController) updateTroopsDeployedCounter() {
	c.troopsDeployedCounter.Label = fmt.Sprintf("Troops deployed: %d/%d", c.deployed, len(game.G.State.Units))
}

func (c *troopDeployController) updateCurrentlyDeploying() {
	u := game.G.State.Units[c.deployed]
	c.deployHint.Label = fmt.Sprintf("Currently deploying: %s", styles.Orange(u.String()))
}
