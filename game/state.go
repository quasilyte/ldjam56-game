package game

import (
	"github.com/quasilyte/ldjam56-game/gcombat"
)

type State struct {
	Level int

	Retries    int
	Casualties int

	Credits int

	CurrentStage *gcombat.Stage

	CardsUnlocked map[gcombat.CardKind]struct{}

	Units []gcombat.UnitKind
}

func (s *State) EnterLevel() {
	level := gcombat.LevelList[s.Level]

	for _, k := range level.NewCards {
		s.CardsUnlocked[k] = struct{}{}
	}
}
