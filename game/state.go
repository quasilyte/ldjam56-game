package game

import (
	"github.com/quasilyte/ldjam56-game/gcombat"
)

type State struct {
	Level int

	Credits int

	Units []gcombat.UnitKind
}
