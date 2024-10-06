package game

import (
	"github.com/quasilyte/ldjam56-game/gcombat"
)

type State struct {
	Level int

	Retries    int
	Casualties int

	Credits int

	Easy bool

	CurrentStage *gcombat.Stage

	UnitsUnlocked [gcombat.NumUnitKinds]bool
	CardsUnlocked map[gcombat.CardKind]struct{}

	Units []gcombat.UnitKind
}

func (s *State) EnterLevel() {
	level := gcombat.LevelList[s.Level]

	for _, k := range level.NewCards {
		s.CardsUnlocked[k] = struct{}{}
	}
	for _, k := range level.NewUnits {
		s.UnitsUnlocked[k] = true
	}
}
