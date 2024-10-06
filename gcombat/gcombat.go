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
	Pos         gmath.Vec
	Rotation    gmath.Rad
	Attacker    *Unit
	Target      *Unit
	AimPos      gmath.Vec
	GoodAim     bool
	Disposed    bool
	FocusedFire bool

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

	Casualties   int
	CasualtyCost int

	Cards []Card
}

func (t *Team) EnemyIndex() int {
	if t.Index == 0 {
		return 1
	}
	return 0
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

func GetUnitStats(k UnitKind) *UnitStats {
	return &unitStatsTable[k]
}

type UnitStats struct {
	Kind UnitKind

	Cost int

	Image           resource.ImageID
	ProjectileImage resource.ImageID
	FireSound       resource.AudioID

	ProjectileHitRadius float64
	Damage              float64
	AntiArmorDamage     float64

	MaxHP float64

	Speed          float64
	TerrainSpeed   [NumTileKinds]float64
	TerrainDefense [NumTileKinds]float64

	Reload       float64
	AccuracyDist float64
	BaseAccuracy float64

	Comment string

	IonStorm       bool
	SplashDamage   bool
	Infantry       bool
	SuppressiveROF bool
}

var unitStatsTable = [...]UnitStats{
	UnitRifle: {
		Kind:            UnitRifle,
		Cost:            10,
		Comment:         "Rushes its lane, strong at short-range",
		Image:           assets.ImageUnitRifle,
		ProjectileImage: assets.ImageProjectileRifle,
		FireSound:       assets.AudioRifle1,
		MaxHP:           13,
		Speed:           18,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    1.0,
			TileMountains: 0.2,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.75,
			TileMountains: -0.3,
		},
		ProjectileHitRadius: 6,
		Damage:              2,
		AntiArmorDamage:     0.5,
		Reload:              0.8,
		AccuracyDist:        64 * 3,
		BaseAccuracy:        0.4,
		IonStorm:            false,
		Infantry:            true,
		SuppressiveROF:      true,
		SplashDamage:        false,
	},

	UnitLaser: {
		Kind:            UnitLaser,
		Cost:            15,
		Comment:         "A heavy laser guard, a situational mid-range support",
		Image:           assets.ImageUnitLaser,
		ProjectileImage: assets.ImageProjectileLaser,
		FireSound:       assets.AudioLaser1,
		MaxHP:           15,
		Speed:           12,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    0.8,
			TileMountains: 0.15,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.6,
			TileMountains: -0.3,
		},
		ProjectileHitRadius: 8,
		Damage:              4,
		AntiArmorDamage:     0.85,
		Reload:              1.9,
		AccuracyDist:        64 * 5,
		BaseAccuracy:        0.65,
		IonStorm:            true,
		Infantry:            true,
		SuppressiveROF:      true,
		SplashDamage:        false,
	},

	UnitMissile: {
		Kind:            UnitMissile,
		Cost:            30,
		Comment:         "Anti-vehicle trooper, good accuracy, splash",
		Image:           assets.ImageUnitMissile,
		ProjectileImage: assets.ImageProjectileMissile,
		FireSound:       assets.AudioMissile1,
		MaxHP:           12,
		Speed:           10,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    0.7,
			TileMountains: 0.15,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.3,
			TileMountains: -0.2,
		},
		ProjectileHitRadius: 10,
		Damage:              7,
		AntiArmorDamage:     1.5,
		Reload:              4.0,
		AccuracyDist:        64 * 7,
		BaseAccuracy:        0.85,
		IonStorm:            false,
		Infantry:            true,
		SuppressiveROF:      false,
		SplashDamage:        true,
	},

	UnitHunter: {
		Kind:            UnitHunter,
		Cost:            45,
		Comment:         "A laser crawler that doesn't fear mountains",
		Image:           assets.ImageUnitHunter,
		ProjectileImage: assets.ImageProjectileLaser,
		FireSound:       assets.AudioHunter1,
		MaxHP:           40,
		Speed:           22,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    0.8,
			TileMountains: 0.8,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    0.45,
			TileMountains: 0.55,
		},
		ProjectileHitRadius: 8,
		Damage:              4,
		AntiArmorDamage:     0.75,
		Reload:              0.85,
		AccuracyDist:        64 * 4,
		BaseAccuracy:        0.6,
		IonStorm:            true,
		Infantry:            false,
		SuppressiveROF:      true,
		SplashDamage:        false,
	},

	UnitTank: {
		Kind:            UnitTank,
		Cost:            55,
		Comment:         "Tanky, but low damage-per-second, splash",
		Image:           assets.ImageUnitTank,
		ProjectileImage: assets.ImageProjectileTank,
		FireSound:       assets.AudioTank1,
		MaxHP:           65,
		Speed:           20,
		TerrainSpeed: [NumTileKinds]float64{
			TilePlains:    1.0,
			TileForest:    0.4,
			TileMountains: 0.05,
		},
		TerrainDefense: [NumTileKinds]float64{
			TilePlains:    0.0,
			TileForest:    -0.3,
			TileMountains: -0.4,
		},
		ProjectileHitRadius: 8,
		Damage:              10,
		AntiArmorDamage:     1.0,
		Reload:              3.5,
		AccuracyDist:        64 * 5,
		BaseAccuracy:        0.7,
		IonStorm:            false,
		Infantry:            false,
		SuppressiveROF:      false,
		SplashDamage:        true,
	},
}

type UnitKind int

const (
	UnitUnknown UnitKind = iota

	UnitRifle
	UnitLaser
	UnitMissile
	UnitHunter
	UnitTank

	NumUnitKinds
)

func (k UnitKind) String() string {
	switch k {
	case UnitRifle:
		return "Rifle infantry"
	case UnitLaser:
		return "Laser infantry"
	case UnitMissile:
		return "Missile infantry"
	case UnitHunter:
		return "Hunter"
	case UnitTank:
		return "Tank"
	default:
		return "?"
	}
}
