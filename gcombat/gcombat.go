package gcombat

import (
	"github.com/quasilyte/gmath"
)

type Stage struct {
	Teams []*Team

	Time float64
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
)

type Team struct {
	Units []*Unit

	Cards []Card
}

type Unit struct {
	Stats *UnitStats

	Team *Team

	Pos gmath.Vec

	HP float64
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

	MaxHP float64

	Speed float64

	Infantry bool
}

var unitStatsTable = [...]UnitStats{
	UnitRifle: {
		Kind:     UnitRifle,
		MaxHP:    10,
		Speed:    32,
		Infantry: true,
	},

	UnitLaser: {
		Kind:     UnitLaser,
		MaxHP:    12,
		Speed:    20,
		Infantry: true,
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
