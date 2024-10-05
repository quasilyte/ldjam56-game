package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/eui"
	"github.com/quasilyte/ldjam56-game/game"
	"github.com/quasilyte/ldjam56-game/gcombat"
	"github.com/quasilyte/ldjam56-game/scenes/sceneutil"
	"github.com/quasilyte/ldjam56-game/styles"
)

type cardpickController struct {
	level *gcombat.Level

	counterLabel *widget.Text

	cardsSelection [][]gcombat.CardKind
	cardsPicked    []gcombat.CardKind
	cardButtons    [][]*cardpickButton

	startButton *widget.Button
}

type cardpickButton struct {
	selected bool
	btn      *widget.Button
	contents *widget.Container
}

func NewCardpickController() *cardpickController {
	return &cardpickController{}
}

func (c *cardpickController) Init(ctx gscene.InitContext) {
	c.level = gcombat.LoadLevel(gcombat.LevelList[game.G.State.Level])

	c.cardsSelection = c.generateCards()

	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Tactics",
		Font: assets.Font2,
	}))

	grid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(len(c.cardsSelection)),
			widget.GridLayoutOpts.Spacing(4, 4),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
			}),
		),
	)
	root.AddChild(grid)

	c.cardButtons = make([][]*cardpickButton, len(c.cardsSelection))
	for i, rowCards := range c.cardsSelection {
		c.cardButtons[i] = make([]*cardpickButton, len(rowCards))
		for j, k := range rowCards {
			stacked := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewStackedLayout()),
			)
			btn := game.G.UI.NewButton(eui.ButtonConfig{
				Font:      assets.FontTiny,
				MinWidth:  96,
				MinHeight: 60,
				OnClick: func() {
					c.onCardPicked(i, j)
				},
			})
			stacked.AddChild(btn)

			contents := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(2),
				)),
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
						HorizontalPosition: widget.AnchorLayoutPositionCenter,
						VerticalPosition:   widget.AnchorLayoutPositionCenter,
					}),
				),
			)
			contentsAnchor := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			)
			contentsAnchor.AddChild(contents)
			contents.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text: k.Info().Category.String(),
				Font: assets.Font1,
				LayoutData: widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			}))
			stacked.AddChild(contentsAnchor)

			grid.AddChild(stacked)

			c.cardButtons[i][j] = &cardpickButton{
				btn:      btn,
				contents: contents,
			}
		}
	}

	c.counterLabel = game.G.UI.NewText(eui.TextConfig{
		Font: assets.FontTiny,
	})
	c.updateCounterLabel()
	root.AddChild(c.counterLabel)

	c.startButton = game.G.UI.NewButton(eui.ButtonConfig{
		Text: "START",
		OnClick: func() {
			team := &gcombat.Team{}
			units := make([]*gcombat.Unit, len(game.G.State.Units))
			for i, u := range game.G.State.Units {
				units[i] = gcombat.NewUnit(u)
				units[i].Team = team
			}
			team.Units = units
			stage := gcombat.CreateStage(gcombat.StageConfig{
				Level: c.level,
				Team1: team,
			})
			game.G.State.CurrentStage = stage
		},
	})
	c.startButton.GetWidget().Disabled = true
	root.AddChild(c.startButton)

	game.G.UI.Build(ctx.Scene, root)
}

func (c *cardpickController) Update(delta float64) {}

func (c *cardpickController) updateCounterLabel() {
	c.counterLabel.Label = fmt.Sprintf("Cards picked: %d/%d",
		len(c.cardsPicked), c.level.CardPicks)
}

func (c *cardpickController) onCardPicked(row, col int) {
	selectedCard := c.cardsSelection[row][col]
	c.cardsPicked = append(c.cardsPicked, selectedCard)

	if len(c.cardsPicked) == 1 {
		// The first move.
		// Disable all buttons.
		for _, buttons := range c.cardButtons {
			for _, b := range buttons {
				b.btn.GetWidget().Disabled = true
				label := b.contents.Children()[0].(*widget.Text)
				label.Color = styles.ColorNormal.Color()
			}
		}
	}

	selectedButton := c.cardButtons[row][col]
	selectedButton.selected = true
	selectedButton.btn.GetWidget().Disabled = true
	{
		selectedButton.contents.RemoveChildren()
		words := strings.Fields(selectedCard.Info().Name)
		for _, word := range words {
			selectedButton.contents.AddChild(game.G.UI.NewText(eui.TextConfig{
				Text:  word,
				Font:  assets.FontTiny,
				Color: styles.ColorBackground.Color(),
				LayoutData: widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			}))
		}
		selectedButton.contents.AddChild(game.G.UI.NewText(eui.TextConfig{
			Text: fmt.Sprintf("[ %d ]", len(c.cardsPicked)),
			Font: assets.FontTiny,
			LayoutData: widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			},
		}))
	}

	// Unlock the neighbours.
	probes := [][2]int{
		{+1, 0},
		{0, +1},
		{-1, 0},
		{0, -1},
	}
	for _, pair := range probes {
		newRow := row + pair[0]
		newCol := col + pair[1]
		if newRow < 0 || newRow >= len(c.cardButtons) {
			continue
		}
		if newCol < 0 || newCol >= len(c.cardButtons) {
			continue
		}

		b := c.cardButtons[newRow][newCol]
		if b.selected {
			continue
		}
		if b.btn.GetWidget().Disabled {
			b.btn.GetWidget().Disabled = false
			b.contents.RemoveChildren()
			card := c.cardsSelection[newRow][newCol]
			words := strings.Fields(card.Info().Name)
			for _, word := range words {
				b.contents.AddChild(game.G.UI.NewText(eui.TextConfig{
					Text:  word,
					Font:  assets.FontTiny,
					Color: styles.ColorBright.Color(),
					LayoutData: widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
						Stretch:  true,
					},
				}))
			}
		}
	}

	if len(c.cardsPicked) == c.level.CardPicks {
		// All cards are picked.
		for _, buttons := range c.cardButtons {
			for _, b := range buttons {
				if b.selected {
					continue
				}
				b.btn.GetWidget().Disabled = true
				for _, label := range b.contents.Children() {
					label := label.(*widget.Text)
					label.Color = styles.ColorNormal.Color()
				}
			}
		}
		c.startButton.GetWidget().Disabled = false
	}

	c.updateCounterLabel()
}

func (c *cardpickController) generateCards() [][]gcombat.CardKind {
	result := make([][]gcombat.CardKind, c.level.CardPicks)

	var cardsPool []gcombat.CardKind
	nextCard := func() gcombat.CardKind {
		if len(cardsPool) == 0 {
			for k := range game.G.State.CardsUnlocked {
				cardsPool = append(cardsPool, k)
			}
			gmath.Shuffle(&game.G.Rand, cardsPool)

		}
		c := cardsPool[len(cardsPool)-1]
		cardsPool = cardsPool[:len(cardsPool)-1]
		return c
	}

	for i := 0; i < c.level.CardPicks; i++ {
		result[i] = make([]gcombat.CardKind, c.level.CardPicks)
		for j := 0; j < c.level.CardPicks; j++ {
			k := nextCard()
			result[i][j] = k
		}
	}

	return result
}
