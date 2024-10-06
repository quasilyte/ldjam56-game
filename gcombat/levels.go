package gcombat

type LevelDef struct {
	Tiles [][]rune

	EnemyTroops []rune
	EnemyDeploy [][]rune

	CardPicks int

	DeployWidth int

	Hint   string
	Reward int

	NewCards []CardKind

	EnemyCards []CardKind
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
	//
	// Deployment:
	// * 'i'=infantry

	{
		// Level 1:
		// Enemy has range advantage, the player is
		// expected to rush with their gatling troops.
		// Using forests for cover is encouraged.
		// Choosing a mountain-heavy lane is discouraged.
		Hint:   "Rifles beat lasers when close enough",
		Reward: 45,
		Tiles: [][]rune{
			{'F', ' ', 'M', ' ', ' ', ' '},
			{' ', ' ', 'M', 'M', 'M', ' '},
			{' ', ' ', ' ', ' ', ' ', ' '},
			{'F', 'F', 'F', 'F', ' ', ' '},
			{' ', 'M', 'M', 'F', ' ', ' '},
		},
		EnemyDeploy: [][]rune{
			{'F', ' ', 'M', ' ', ' ', ' '},
			{' ', ' ', 'M', 'M', 'M', ' '},
			{' ', ' ', ' ', ' ', ' ', ' '},
			{'F', 'F', 'F', 'F', ' ', 'i'},
			{' ', 'M', 'M', 'F', ' ', 'i'},
		},
		EnemyTroops: []rune{'L', 'L', 'L', 'L'},
		CardPicks:   2,
		DeployWidth: 1,
		NewCards: []CardKind{
			CardInfatryCharge,
			CardTakeCover,

			CardSuppressiveFire,

			CardLuckyShot,
		},
		EnemyCards: []CardKind{
			CardLuckyShot,
			CardSuppressiveFire,
		},
	},

	{
		// Level 2:
		// Basically a level 1 reversed.
		// The player is expected to get some laser troops.
		Hint: "TODO",
		Tiles: [][]rune{
			{' ', 'M', ' ', 'M', ' ', ' ', 'F'},
			{' ', 'M', ' ', ' ', ' ', ' ', 'F'},
			{' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'F', ' ', ' ', 'F', ' '},
			{' ', 'M', 'F', 'F', 'F', ' ', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', 'M', ' ', 'M', ' ', ' ', 'i'},
			{' ', 'M', ' ', ' ', ' ', ' ', 'i'},
			{' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'F', ' ', ' ', 'F', ' '},
			{' ', 'M', 'F', 'F', 'F', ' ', 'i'},
		},
		EnemyTroops: []rune{'R', 'R', 'R', 'R', 'R'},
		CardPicks:   3,
		DeployWidth: 2,
		NewCards:    []CardKind{
			// CardInfatryCharge,
			// CardTakeCover,

			// CardSuppressiveFire,

			// CardLuckyShot,
		},
		EnemyCards: []CardKind{
			CardLuckyShot,
			CardSuppressiveFire,
			CardSuppressiveFire,
		},
	},
}
