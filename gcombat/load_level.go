package gcombat

import (
	"github.com/quasilyte/ldjam56-game/dat"
)

type Level struct {
	EnemyTroops []UnitKind
}

func LoadLevel(level *dat.Level) *Level {
	result := &Level{}

	for _, unitTag := range level.EnemyTroops {
		switch unitTag {
		case 'R':
			result.EnemyTroops = append(result.EnemyTroops, UnitRifle)
		case 'L':
			result.EnemyTroops = append(result.EnemyTroops, UnitLaser)
		}
	}

	return result
}
