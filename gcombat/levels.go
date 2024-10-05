package gcombat

type LevelDef struct {
	Tiles [][]rune

	EnemyTroops []rune

	CardPicks int

	Hint string

	NewCards []CardKind
}

var LevelList = []*LevelDef{
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
		Hint: "Rifles beat lasers when close enough",
		Tiles: [][]rune{
			{'F', ' ', 'M', ' ', ' ', ' '},
			{' ', ' ', 'M', 'M', 'M', ' '},
			{' ', ' ', ' ', ' ', ' ', ' '},
			{'F', 'F', 'F', 'F', ' ', ' '},
			{' ', 'M', 'M', 'F', ' ', ' '},
		},
		EnemyTroops: []rune{'L', 'L', 'L', 'L'},
		CardPicks:   2,
		NewCards: []CardKind{
			CardInfatryCharge,
			CardTakeCover,

			CardSuppressiveFire,

			CardLuckyShot,
		},
	},
}
