package gcombat

type Level struct {
	EnemyTroops []UnitKind
	EnemyCards  []CardKind

	EnemyInfantrySpots [][2]int

	Tiles [][]TileKind

	Reward      int
	CardPicks   int
	DeployWidth int

	Hint string
}

func LoadLevel(level *LevelDef) *Level {
	result := &Level{
		CardPicks:   level.CardPicks,
		Hint:        level.Hint,
		EnemyCards:  level.EnemyCards,
		DeployWidth: level.DeployWidth,
		Reward:      level.Reward,
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

	for i, rowTiles := range level.EnemyDeploy {
		for j, colTag := range rowTiles {
			switch colTag {
			case 'i':
				result.EnemyInfantrySpots = append(result.EnemyInfantrySpots, [2]int{i, j})
			}
		}
	}

	return result
}
