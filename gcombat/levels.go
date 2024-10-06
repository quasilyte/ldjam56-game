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
	NewUnits []UnitKind

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
	// * 'M'=missile
	// * 'H'=hunter
	// * 'T'=tank
	//
	// Deployment:
	// * 'i'=infantry
	// * 'v'=normal vehicle (not artillery)

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
			{' ', ' ', ' ', ' ', ' ', 'v'},
			{'F', 'F', 'F', 'F', ' ', 'i'},
			{' ', 'M', 'M', 'F', ' ', 'i'},
		},
		EnemyTroops: []rune{'L', 'L', 'L', 'L'},
		CardPicks:   2,
		DeployWidth: 1,
		NewCards: []CardKind{
			CardInfantryCharge,

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
		Hint:   "Traversing mountains reduces defense",
		Reward: 60,
		Tiles: [][]rune{
			{' ', 'F', ' ', 'M', 'M', ' ', 'F'},
			{' ', 'M', ' ', ' ', ' ', ' ', 'F'},
			{' ', ' ', ' ', ' ', ' ', 'M', ' '},
			{' ', ' ', 'F', ' ', ' ', 'F', ' '},
			{' ', 'M', ' ', 'F', 'F', ' ', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', 'M', ' ', 'M', ' ', ' ', 'i'},
			{' ', 'M', ' ', ' ', ' ', ' ', 'i'},
			{' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', 'F', ' ', ' ', 'F', ' '},
			{' ', 'M', ' ', 'F', 'F', ' ', 'i'},
		},
		EnemyTroops: []rune{'R', 'R', 'R', 'R', 'R', 'R'},
		// EnemyTroops: []rune{'R'},
		CardPicks:   3,
		DeployWidth: 2,
		NewCards: []CardKind{
			CardTakeCover,
			CardFirstAid,
		},
		NewUnits: []UnitKind{
			UnitLaser,
		},
		EnemyCards: []CardKind{
			CardLuckyShot,
			CardTakeCover,
			CardInfantryCharge,
		},
	},

	{
		// Level 3:
		// The first level with vehicles.
		// The player is expected to get some anti-armor.
		Hint:   "You get a 20% refund for casualties",
		Reward: 85,
		Tiles: [][]rune{
			{' ', 'F', ' ', ' ', 'M', ' '},
			{'F', ' ', ' ', ' ', ' ', ' '},
			{' ', 'F', 'M', 'M', ' ', ' '},
			{' ', ' ', 'M', ' ', ' ', ' '},
			{' ', 'F', 'F', ' ', ' ', 'F'},
		},
		EnemyDeploy: [][]rune{
			{' ', 'F', ' ', ' ', 'M', ' '},
			{'F', ' ', ' ', ' ', 'v', 'v'},
			{' ', 'F', 'M', 'M', ' ', 'i'},
			{' ', ' ', 'M', ' ', ' ', 'i'},
			{' ', 'F', 'F', ' ', ' ', 'F'},
		},
		EnemyTroops: []rune{'T', 'R', 'R', 'L', 'L'},
		CardPicks:   4,
		DeployWidth: 2,
		NewCards: []CardKind{
			CardStandGround,
			CardBadCover,
		},
		NewUnits: []UnitKind{
			UnitMissile,
		},
		EnemyCards: []CardKind{
			CardBadCover,
			CardTankRush,
			CardFirstAid,
			CardTakeCover,
		},
	},

	{
		// Level 4:
		Hint:   "Infantry can't survive the treads of a tank",
		Reward: 95,
		Tiles: [][]rune{
			{' ', ' ', 'F', 'F', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', 'F', ' ', ' ', ' ', ' '},
			{'M', ' ', ' ', ' ', 'F', ' ', ' ', 'M'},
			{' ', 'M', 'M', ' ', 'F', ' ', ' ', ' '},
			{' ', 'M', 'M', 'F', 'F', ' ', ' ', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', ' ', 'F', 'F', ' ', ' ', ' ', 'v'},
			{' ', ' ', ' ', 'F', ' ', ' ', ' ', 'v'},
			{'M', ' ', ' ', ' ', 'F', ' ', ' ', 'M'},
			{' ', 'M', 'M', ' ', 'F', ' ', ' ', 'i'},
			{' ', 'M', 'M', 'F', 'F', ' ', ' ', 'i'},
		},
		EnemyTroops: []rune{'T', 'T', 'M', 'M', 'L', 'R'},
		CardPicks:   4,
		DeployWidth: 1,
		NewCards: []CardKind{
			CardFocusFire,
		},
		NewUnits: []UnitKind{
			UnitHunter,
		},
		EnemyCards: []CardKind{
			CardTankRush,
			CardTankRush,
			CardFirstAid,
			CardTakeCover,
		},
	},

	{
		// Level 5:
		Hint:   "Missiles and tanks deal splash damage",
		Reward: 75,
		Tiles: [][]rune{
			{' ', ' ', 'F', ' ', 'F', 'F'},
			{' ', 'F', 'M', ' ', 'F', 'M'},
			{' ', ' ', ' ', ' ', 'F', ' '},
			{' ', ' ', 'M', 'F', 'F', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', ' ', 'F', ' ', 'F', 'i'},
			{' ', 'F', ' ', ' ', 'F', 'M'},
			{' ', ' ', ' ', ' ', 'F', 'i'},
			{' ', ' ', 'M', 'F', 'F', 'i'},
		},
		EnemyTroops: []rune{
			'R', 'R', 'R', 'R', 'R',
			'R', 'R', 'R', 'R', 'R',
			'R', 'R', 'R', 'R', 'R',
		},
		CardPicks:   4,
		DeployWidth: 1,
		NewCards: []CardKind{
			CardTankRush,
		},
		NewUnits: []UnitKind{
			UnitTank,
		},
		EnemyCards: []CardKind{
			CardInfantryCharge,
			CardSuppressiveFire,
			CardFocusFire,
			CardIonStorm,
		},
	},

	{
		// Level 6:
		Hint:   "The final battle is near",
		Reward: 100,
		Tiles: [][]rune{
			{'F', ' ', 'M', 'M', ' ', 'F', 'F', ' ', 'F', 'F'},
			{' ', ' ', ' ', ' ', 'M', 'F', 'F', ' ', ' ', 'F'},
			{' ', 'F', ' ', 'M', ' ', 'M', ' ', ' ', ' ', ' '},
			{' ', 'F', ' ', 'M', 'F', ' ', 'M', ' ', 'F', ' '},
			{' ', 'F', 'M', ' ', 'F', ' ', ' ', 'M', 'F', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', ' ', 'M', 'M', ' ', 'F', 'F', ' ', 'F', 'i'},
			{' ', ' ', ' ', ' ', 'M', 'F', 'F', ' ', ' ', 'i'},
			{' ', 'F', ' ', 'M', ' ', 'M', ' ', ' ', ' ', 'v'},
			{' ', 'F', ' ', 'M', 'F', ' ', 'M', ' ', 'F', ' '},
			{' ', 'F', 'M', ' ', 'F', ' ', ' ', 'M', 'F', ' '},
		},
		EnemyTroops: []rune{'T', 'M', 'M', 'M', 'M', 'L'},
		CardPicks:   5,
		DeployWidth: 3,
		NewCards: []CardKind{
			CardIonStorm,
		},
		NewUnits: []UnitKind{},
		EnemyCards: []CardKind{
			CardIonStorm,
			CardFirstAid,
			CardLuckyShot,
			CardLuckyShot,
			CardStandGround,
		},
	},

	{
		// Level 7:
		Hint:   "Turn their strength into weakness",
		Reward: 100,
		Tiles: [][]rune{
			{' ', ' ', 'M', 'F', 'F', ' ', 'M', 'F', ' ', 'F'},
			{' ', ' ', 'M', 'M', 'F', 'F', 'F', ' ', 'F', 'F'},
			{'F', 'M', ' ', 'F', 'M', ' ', 'F', 'F', 'F', ' '},
		},
		EnemyDeploy: [][]rune{
			{' ', ' ', 'M', 'F', 'F', ' ', 'M', 'v', ' ', 'i'},
			{' ', ' ', 'M', 'M', 'F', 'F', 'F', ' ', 'v', 'i'},
			{'F', 'M', ' ', 'F', 'M', ' ', 'F', 'v', 'i', ' '},
		},
		EnemyTroops: []rune{
			'H', 'H', 'H', 'H',
			'L', 'L', 'L', 'L', 'L',
			'R', 'R', 'R', 'R', 'R',
			'R', 'R',
		},
		CardPicks:   5,
		DeployWidth: 3,
		NewCards:    []CardKind{},
		NewUnits:    []UnitKind{},
		EnemyCards: []CardKind{
			CardInfantryCharge,
			CardStandGround,
			CardTakeCover,
			CardTakeCover,
			CardInfantryCharge,
		},
	},
}
