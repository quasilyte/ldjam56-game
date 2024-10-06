package gcombat

type Card struct {
	TeamIndex int
	Kind      CardKind
}

type CardInfo struct {
	Name string

	Category CardCategory

	Duration int
}

type CardCategory int

const (
	CardCategoryMovement CardCategory = iota
	CardCategorySpecial
	CardCategoryBonus
	CardCategoryModifier
)

func (c CardCategory) String() string {
	switch c {
	case CardCategoryMovement:
		return "Movement"
	case CardCategorySpecial:
		return "Special"
	case CardCategoryBonus:
		return "Bonus"
	case CardCategoryModifier:
		return "Mod"
	default:
		return "?"
	}
}

type CardKind int

func (k CardKind) Info() *CardInfo {
	return &cardInfoTable[k]
}

const (
	CardUnknown CardKind = iota

	// Movement command cards.
	CardInfatryCharge
	CardStandGround
	CardTakeCover

	// Special action cards.
	CardSuppressiveFire
	CardArtilleryFire

	// Bonus cards.
	CardLuckyShot

	// Modifier cards.
	CardIonStorm
)

var cardInfoTable = [...]CardInfo{
	CardInfatryCharge: {
		Name:     "Infantry Charge",
		Category: CardCategoryMovement,
		Duration: 1,
	},
	CardStandGround: {
		Name:     "Stand Ground",
		Category: CardCategoryMovement,
		Duration: 1,
	},
	CardTakeCover: {
		Name:     "Take Cover",
		Category: CardCategoryMovement,
		Duration: 1,
	},

	CardSuppressiveFire: {
		Name:     "Suppressive Fire",
		Category: CardCategorySpecial,
		Duration: 1,
	},
	CardArtilleryFire: {
		Name:     "Artillery Fire",
		Category: CardCategorySpecial,
		Duration: 1,
	},

	CardLuckyShot: {
		Name:     "Lucky Shot",
		Category: CardCategoryBonus,
		Duration: 1,
	},

	CardIonStorm: {
		Name:     "Ion Storm",
		Category: CardCategoryModifier,
		Duration: 1,
	},
}
