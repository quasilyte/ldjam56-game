package gcombat

import "strings"

type Card struct {
	TeamIndex int
	Kind      CardKind
}

type CardInfo struct {
	Name string

	Category CardCategory

	Description string

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
	CardInfantryCharge
	CardTankRush
	CardStandGround
	CardTakeCover

	// Special action cards.
	CardSuppressiveFire
	CardFocusFire

	// Bonus cards.
	CardLuckyShot
	CardFirstAid

	// Modifier cards.
	CardBadCover
	CardIonStorm
)

var cardInfoTable = [...]CardInfo{
	CardInfantryCharge: {
		Name:     "Infantry Charge",
		Category: CardCategoryMovement,
		Duration: 1,
		Description: strings.Join([]string{
			"Makes infantry move faster",
			"also makes them immune to the",
			"effects of Suppressive Fire",
		}, "\n"),
	},
	CardStandGround: {
		Name:     "Stand Ground",
		Category: CardCategoryMovement,
		Duration: 1,
		Description: strings.Join([]string{
			"Order units to stop",
			"moving towards the enemy",
		}, "\n"),
	},
	CardTakeCover: {
		Name:     "Take Cover",
		Category: CardCategoryMovement,
		Duration: 1,
		Description: strings.Join([]string{
			"Order units to find the closest",
			"cover (e.g. a forest)",
		}, "\n"),
	},
	CardTankRush: {
		Name:     "Tank Rush",
		Category: CardCategoryMovement,
		Duration: 2,
		Description: strings.Join([]string{
			"Order your tanks to move forward",
		}, "\n"),
	},

	CardSuppressiveFire: {
		Name:     "Suppressive Fire",
		Category: CardCategorySpecial,
		Duration: 1,
		Description: strings.Join([]string{
			"Increases some units rate-of-fire",
			"while decreasing their accuracy.",
			"Routes enemy infantry for a while",
		}, "\n"),
	},
	CardFocusFire: {
		Name:     "Focus Fire",
		Category: CardCategorySpecial,
		Duration: 1,
		Description: strings.Join([]string{
			"Order units to concentrate",
			"their firepower on wounded targets",
		}, "\n"),
	},

	CardLuckyShot: {
		Name:     "Lucky Shot",
		Category: CardCategoryBonus,
		Duration: 1,
		Description: strings.Join([]string{
			"Significantly increases accuracy of",
			"your troops. Extra powerful against",
			"the targets under Bad Cover",
		}, "\n"),
	},
	CardFirstAid: {
		Name:     "First Aid",
		Category: CardCategoryBonus,
		Duration: 1,
		Description: strings.Join([]string{
			"Instantly heals most of the wounds of",
			"your infantry",
		}, "\n"),
	},

	CardBadCover: {
		Name:     "Bad Cover",
		Category: CardCategoryModifier,
		Duration: 2,
		Description: strings.Join([]string{
			"Forest stops being a good cover",
			"for the duration of this modifier.",
			"Affects everyone",
		}, "\n"),
	},
	CardIonStorm: {
		Name:     "Ion Storm",
		Category: CardCategoryModifier,
		Duration: 2,
		Description: strings.Join([]string{
			"Makes it impossible to fire laser weapons",
			"for the duration of this modifier.",
			"Affects everyone",
		}, "\n"),
	},
}
