package dat

type Level struct {
	Tiles [][]rune

	EnemyTroops []rune
}

var LevelList = []*Level{
	// Tiles:
	// * ' '=plains
	// * 'F'=forest
	// * 'M'=mountains
	//
	// Units:
	// * 'R'=rifle
	// * 'L'=laser

	{
		// Level 1:
		// Enemy has range advantage, the player is
		// expected to rush with their gatling troops.
		// Using forests for cover is encouraged.
		// Choosing a mountain-heavy lane is discouraged.
		Tiles: [][]rune{
			{'F', ' ', 'M', ' ', ' ', ' '},
			{' ', ' ', 'M', 'M', 'M', ' '},
			{' ', ' ', ' ', ' ', ' ', ' '},
			{'F', 'F', 'F', 'F', ' ', ' '},
			{' ', 'M', 'M', 'F', ' ', ' '},
		},
		EnemyTroops: []rune{'L', 'L', 'L', 'L'},
	},
}
