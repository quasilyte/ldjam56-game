package gcombat

type Level struct {
	EnemyTroops []UnitKind
	EnemyCards  []CardKind

	Tiles [][]TileKind

	CardPicks int

	Hint string
}

func LoadLevel(level *LevelDef) *Level {
	result := &Level{
		CardPicks:  level.CardPicks,
		Hint:       level.Hint,
		EnemyCards: level.EnemyCards,
	}

	for _, unitTag := range level.EnemyTroops {
		switch unitTag {
		case 'R':
			result.EnemyTroops = append(result.EnemyTroops, UnitRifle)
		case 'L':
			result.EnemyTroops = append(result.EnemyTroops, UnitLaser)
		}
	}

	result.Tiles = make([][]TileKind, len(level.Tiles))
	for i, rowTiles := range level.Tiles {
		result.Tiles[i] = make([]TileKind, len(rowTiles))
		for j, colTag := range rowTiles {
			var k TileKind
			switch colTag {
			case ' ':
				k = TilePlains
			case 'M':
				k = TileMountains
			case 'F':
				k = TileForest
			}
			result.Tiles[i][j] = k
		}
	}

	return result
}
