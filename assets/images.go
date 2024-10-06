package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"

	_ "image/png"
)

func registerImageResources(loader *resource.Loader) {
	resources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonIdle:         {Path: "image/ui/button_idle.png"},
		ImageUIButtonHover:        {Path: "image/ui/button_hover.png"},
		ImageUIButtonPressed:      {Path: "image/ui/button_pressed.png"},
		ImageUIButtonDisabled:     {Path: "image/ui/button_disabled.png"},
		ImageUITileButtonIdle:     {Path: "image/ui/tilebutton_idle.png"},
		ImageUITileButtonHover:    {Path: "image/ui/tilebutton_hover.png"},
		ImageUITileButtonPressed:  {Path: "image/ui/tilebutton_pressed.png"},
		ImageUITileButtonDisabled: {Path: "image/ui/tilebutton_disabled.png"},
		ImageUIPanel:              {Path: "image/ui/panel.png"},
		ImageUITooltip:            {Path: "image/ui/tooltip.png"},

		ImageTilePlains:    {Path: "image/tile_plains.png"},
		ImageTileMountains: {Path: "image/tile_mountains.png"},
		ImageTileForest:    {Path: "image/tile_forest.png"},
		ImageTileBg1:       {Path: "image/tile_bg1.png"},
		ImageTileBg2:       {Path: "image/tile_bg2.png"},
		ImageTileBg3:       {Path: "image/tile_bg3.png"},
		ImageTileGrid:      {Path: "image/tile_grid.png"},
		ImageTileSelector:  {Path: "image/tile_selector.png"},

		ImageUnitRifle:   {Path: "image/unit_rifle.png"},
		ImageUnitLaser:   {Path: "image/unit_laser.png"},
		ImageUnitMissile: {Path: "image/unit_missile.png"},
		ImageUnitHunter:  {Path: "image/unit_hunter.png"},
		ImageUnitTank:    {Path: "image/unit_tank.png"},

		ImageProjectileRifle:   {Path: "image/projectile_rifle.png"},
		ImageProjectileLaser:   {Path: "image/projectile_laser.png"},
		ImageProjectileMissile: {Path: "image/projectile_missile.png"},
		ImageProjectileTank:    {Path: "image/projectile_tank.png"},

		ImageExplosion: {Path: "image/explosion.png", FrameWidth: 12},
	}

	for id, info := range resources {
		loader.ImageRegistry.Set(id, info)
		loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUIButtonDisabled
	ImageUITileButtonIdle
	ImageUITileButtonHover
	ImageUITileButtonPressed
	ImageUITileButtonDisabled
	ImageUIPanel
	ImageUITooltip

	ImageTilePlains
	ImageTileMountains
	ImageTileForest
	ImageTileBg1
	ImageTileBg2
	ImageTileBg3
	ImageTileGrid
	ImageTileSelector

	ImageUnitRifle
	ImageUnitLaser
	ImageUnitMissile
	ImageUnitHunter
	ImageUnitTank

	ImageProjectileRifle
	ImageProjectileLaser
	ImageProjectileMissile
	ImageProjectileTank

	ImageExplosion
)
