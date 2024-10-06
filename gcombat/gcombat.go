package gcombat

import (
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/ldjam56-game/assets"
)

type Stage struct {
	Teams []*Team

	Level *Level

	Height float64
	Width  float64

	MapBg *ebiten.Image

	Time float64
}

type Projectile struct {
	Pos      gmath.Vec
	Rotation gmath.Rad
	Attacker *Unit
	Target   *Unit
	AimPos   gmath.Vec
	GoodAim  bool
	Disposed bool

	EventDisposed gsignal.Event[gsignal.Void]
}

type StageConfig struct {
	Level *Level
	Team1 *Team
	Team2 *Team
}

func CreateStage(config StageConfig) *Stage {
	stage := &Stage{
		Teams: []*Team{
			config.Team1,
			config.Team2,
		},
		Level:  config.Level,
		Height: float64(64 * len(config.Level.Tiles)),
		Width:  float64(64 * len(config.Level.Tiles[0])),
	}

	return stage
}

type Map struct {
	Tiles [][]Tile
}

type Tile struct {
	Kind TileKind
}

type TileKind int

const (
	TilePlains TileKind = iota
	TileMountains
	TileForest

	NumTileKinds
)

type Team struct {
	Index int

	Units []*Unit

	Casualties int

	Cards []Card
}

type Unit struct {
	Stats *UnitStats

	Team *Team

	Reload float64

	Pos      gmath.Vec
	SpawnPos gmath.Vec

	HP float64

	Waypoint gmath.Vec

	EventDisposed gsignal.Event[gsignal.Void]
}

func (u *Unit) IsDisposed() bool {
	return u.HP <= 0
}

func NewUnit(k UnitKind) *Unit {
	stats := &unitStatsTable[k]
	return &Unit{
		Stats: stats,
		HP:    stats.MaxHP,
	}
}

type UnitStats struct {
	Kind UnitKind

	Cost int

	Image           resource.ImageID
	ProjectileImage resource.ImageID
	FireSound       resource.AudioID

	ProjectileHitRadius float64
	Damage              float64

	MaxHP float64

	Speed          float64
	TerrainSpeed   [NumTileKinds]float64
	TerrainDefense [NumTileKinds]float64

	Reload       float64
	AccuracyDist float64
	BaseAccuracy float64

	Infantry bool
}

var unitStatsTable = [...]UnitStats{
	UnitRifle: {
		Kind:            UnitRifle,
		Cost:            10,
		Image:           assets.ImageUnitRifle,
		ProjectileImage: assets.ImageProjectileRifle,
		FireSound:       assets.AudioRifle1,
		MaxHP:           10,
		Speed:           18,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    1.0,
			TileMountains: 0.2,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.75,
			TileMountains: -0.2,
		},
		ProjectileHitRadius: 6,
		Damage:              2,
		Reload:              0.8,
		AccuracyDist:        64 * 3,
		BaseAccuracy:        0.4,
		Infantry:            true,
	},

	UnitLaser: {
		Kind:            UnitLaser,
		Cost:            15,
		Image:           assets.ImageUnitLaser,
		ProjectileImage: assets.ImageProjectileLaser,
		FireSound:       assets.AudioLaser1,
		MaxHP:           12,
		Speed:           12,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    0.8,
			TileMountains: 0.15,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.55,
			TileMountains: -0.2,
		},
		ProjectileHitRadius: 8,
		Damage:              4,
		Reload:              2.0,
		AccuracyDist:        64 * 5,
		BaseAccuracy:        0.6,
		Infantry:            true,
	},
}

type UnitKind int

const (
	UnitUnknown UnitKind = iota

	UnitRifle
	UnitLaser

	NumUnitKinds
)

func (k UnitKind) String() string {
	switch k {
	case UnitRifle:
		return "Rifle infantry"
	case UnitLaser:
		return "Laser infantry"
	default:
		return "?"
	}
}
