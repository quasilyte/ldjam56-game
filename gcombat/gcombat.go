package gcombat

import (
	"github.com/quasilyte/gmath"
)

type Stage struct {
	Teams []*Team

	Time float64
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

	Pos gmath.Vec

	HP float64
}

type UnitStats struct {
	Kind UnitKind

	MaxHP float64

	Speed float64

	Infantry bool
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
