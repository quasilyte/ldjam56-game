package groundscape

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type Controller struct {
	scene *gscene.Scene

	level *gcombat.Level

	worldOffset gmath.Vec
}

type ControllerConfig struct {
	Level *gcombat.Level
}

func NewController(config ControllerConfig) *Controller {
	return &Controller{
		level: config.Level,
	}
}

func (c *Controller) Init(ctx gscene.InitContext) {
	c.scene = ctx.Scene

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	{
		bg := c.generateLevelImage()
		c.worldOffset.X = float64((int(game.G.WindowSize.X) - bg.Bounds().Dx()) / 2)
		c.worldOffset.Y = 32

		s := graphics.NewSprite()
		s.SetCentered(false)
		s.SetImage(bg)
		s.Pos.Offset = c.worldOffset
		ctx.Scene.AddGraphics(s, 0)
	}

	{
		s := game.G.NewSprite(assets.ImageUnitTank)
		s.Pos.Offset = (gmath.Vec{X: 32, Y: 32}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
	{
		s := game.G.NewSprite(assets.ImageUnitRifle)
		s.Pos.Offset = (gmath.Vec{X: 20, Y: 20}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
	{
		s := game.G.NewSprite(assets.ImageUnitLaser)
		s.Pos.Offset = (gmath.Vec{X: 40, Y: 20}).Add(c.worldOffset)
		ctx.Scene.AddGraphics(s, 0)
	}
}

func (c *Controller) Update(delta float64) {}

func (c *Controller) generateLevelImage() *ebiten.Image {
	mapHeight := len(c.level.Tiles) * 64
	mapWidth := len(c.level.Tiles[0]) * 64
	img := ebiten.NewImage(mapWidth, mapHeight)
	{
		vector.StrokeLine(img, 0, 1, float32(mapWidth), 0, 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, 0, float32(mapHeight), float32(mapWidth), float32(mapHeight), 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, 1, 0, 1, float32(mapHeight), 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, float32(mapWidth), 0, float32(mapWidth), float32(mapHeight), 1, styles.ColorBright.Color(), false)
	}
	for rowNum, rowTiles := range c.level.Tiles {
		for colNum, tileKind := range rowTiles {
			x := colNum * 64
			y := rowNum * 64
			var tileImg resource.ImageID
			switch tileKind {
			case gcombat.TilePlains:
				tileImg = assets.ImageTilePlains
			case gcombat.TileMountains:
				tileImg = assets.ImageTileMountains
			case gcombat.TileForest:
				tileImg = assets.ImageTileForest
			}
			var opts ebiten.DrawImageOptions
			bgSeed := game.G.Rand.IntRange(0, 2)

			opts.GeoM.Reset()
			if game.G.Rand.Bool() {
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(64, 0)
			}
			opts.GeoM.Translate(float64(x), float64(y))
			img.DrawImage(game.G.Loader.LoadImage(assets.ImageTileBg1+resource.ImageID(bgSeed)).Data, &opts)

			// opts.GeoM.Reset()
			// opts.GeoM.Translate(float64(x), float64(y))
			// img.DrawImage(game.G.Loader.LoadImage(assets.ImageTileGrid).Data, &opts)

			opts.GeoM.Reset()
			if game.G.Rand.Bool() {
				opts.GeoM.Scale(-1, 1)
				opts.GeoM.Translate(64, 0)
			}
			opts.GeoM.Translate(float64(x), float64(y))
			img.DrawImage(game.G.Loader.LoadImage(tileImg).Data, &opts)

		}
	}
	return img
}
