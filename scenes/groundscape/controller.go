package groundscape

import (
	"github.com/hajimehoshi/ebiten/v2"
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/dat"
	"github.com/quasilyte/ldjam56-game/game"
)

type Controller struct {
	scene *gscene.Scene

	level *dat.Level

	worldOffset gmath.Vec
}

type ControllerConfig struct {
	Level *dat.Level
}

func NewController(config ControllerConfig) *Controller {
	return &Controller{
		level: config.Level,
	}
}

func (c *Controller) Init(ctx gscene.InitContext) {
	c.scene = ctx.Scene

	{
		screenBg := ebiten.NewImage(1, 1)
		screenBg.Fill(graphics.RGB(0x111412).Color())
		s := graphics.NewSprite()
		s.SetImage(screenBg)
		s.SetScaleX(game.G.WindowSize.X)
		s.SetScaleY(game.G.WindowSize.Y)
		ctx.Scene.AddGraphics(s, 0)
	}

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
	for rowNum, rowTiles := range c.level.Tiles {
		for colNum, colTag := range rowTiles {
			x := colNum * 64
			y := rowNum * 64
			var tileImg resource.ImageID
			switch colTag {
			case ' ':
				tileImg = assets.ImageTilePlains
			case 'M':
				tileImg = assets.ImageTileMountains
			case 'F':
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
