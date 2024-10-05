package gcombat

type Card struct {
	Kind CardKind
}

type CardInfo struct {
	Name string

	Category CardCategory
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
)

var cardInfoTable = [...]CardInfo{
	CardInfatryCharge: {
		Name:     "Infantry Charge",
		Category: CardCategoryMovement,
	},
	CardStandGround: {
		Name:     "Stand Ground",
		Category: CardCategoryMovement,
	},
	CardTakeCover: {
		Name:     "Take Cover",
		Category: CardCategoryMovement,
	},

	CardSuppressiveFire: {
		Name:     "Suppressive Fire",
		Category: CardCategorySpecial,
	},
	CardArtilleryFire: {
		Name:     "Artillery Fire",
		Category: CardCategorySpecial,
	},

	CardLuckyShot: {
		Name:     "Lucky Shot",
		Category: CardCategoryBonus,
	},
}
