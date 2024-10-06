package sceneutil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	graphics "github.com/quasilyte/ebitengine-graphics"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/styles"
)

func NewBackgroundImage() *graphics.Sprite {
	screenBg := ebiten.NewImage(1, 1)
	screenBg.Fill(styles.ColorBackground.ScaleRGB(0.75).Color())
	// screenBg.Fill(graphics.RGB(0x57839C).ScaleRGB(0.4).Color())
	// screenBg.Fill(graphics.RGB(0x70579C).ScaleRGB(0.4).Color())
	s := graphics.NewSprite()
	s.SetImage(screenBg)
	s.SetScaleX(game.G.WindowSize.X)
	s.SetScaleY(game.G.WindowSize.Y)
	return s
}

func CombatMapOffset(img *ebiten.Image) gmath.Vec {
	var worldOffset gmath.Vec
	worldOffset.X = float64((int(game.G.WindowSize.X) - img.Bounds().Dx()) / 2)
	worldOffset.Y = 32
	return worldOffset
}

func CombatMapSprite(img *ebiten.Image) *graphics.Sprite {
	worldOffset := CombatMapOffset(img)

	s := graphics.NewSprite()
	s.SetCentered(false)
	s.SetImage(img)
	s.Pos.Offset = worldOffset
	return s
}

func DrawCombatMap(level *gcombat.Level) *ebiten.Image {
	mapHeight := len(level.Tiles) * 64
	mapWidth := len(level.Tiles[0]) * 64
	img := ebiten.NewImage(mapWidth, mapHeight)
	{
		vector.StrokeLine(img, 0, 1, float32(mapWidth), 0, 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, 0, float32(mapHeight), float32(mapWidth), float32(mapHeight), 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, 1, 0, 1, float32(mapHeight), 1, styles.ColorBright.Color(), false)
		vector.StrokeLine(img, float32(mapWidth), 0, float32(mapWidth), float32(mapHeight), 1, styles.ColorBright.Color(), false)
	}
	for rowNum, rowTiles := range level.Tiles {
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
