package scenes

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/gslices"
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

	enemyCardsPool   []gcombat.CardKind
	enemyLabel       *widget.Text
	enemyCardsPicked []gcombat.CardKind

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

	c.enemyCardsPool = gslices.Clone(c.level.EnemyCards)
	gmath.Shuffle(&game.G.Rand, c.enemyCardsPool)

	root := eui.NewTopLevelRows()

	ctx.Scene.AddGraphics(sceneutil.NewBackgroundImage(), 0)

	root.AddChild(game.G.UI.NewText(eui.TextConfig{
		Text: "Tactic Cards",
		Font: assets.Font2,
	}))

	gridAnchor := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	grid := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(len(c.cardsSelection)),
			widget.GridLayoutOpts.Spacing(4, 4),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	gridAnchor.AddChild(grid)
	root.AddChild(gridAnchor)

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

	{
		panel := game.G.UI.NewPanel(eui.PanelConfig{})
		c.enemyLabel = game.G.UI.NewText(eui.TextConfig{
			MinWidth:    220,
			Font:        assets.FontTiny,
			ForceBBCode: true,
			LayoutData: widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			},
		})
		c.pickEnemyCard()
		panel.AddChild(c.enemyLabel)
		root.AddChild(panel)
	}

	c.startButton = game.G.UI.NewButton(eui.ButtonConfig{
		Text: "START",
		OnClick: func() {
			team1 := &gcombat.Team{Index: 0}
			units1 := make([]*gcombat.Unit, len(game.G.State.Units))
			for i, u := range game.G.State.Units {
				units1[i] = gcombat.NewUnit(u)
				units1[i].Team = team1
			}
			team1.Units = units1
			for _, k := range c.cardsPicked {
				team1.Cards = append(team1.Cards, gcombat.Card{
					TeamIndex: 0,
					Kind:      k,
				})
			}

			team2 := &gcombat.Team{Index: 1}
			units2 := make([]*gcombat.Unit, len(c.level.EnemyTroops))
			for i, u := range c.level.EnemyTroops {
				units2[i] = gcombat.NewUnit(u)
				units2[i].Team = team2
			}
			team2.Units = units2
			for _, k := range c.enemyCardsPicked {
				team2.Cards = append(team2.Cards, gcombat.Card{
					TeamIndex: 1,
					Kind:      k,
				})
			}

			stage := gcombat.CreateStage(gcombat.StageConfig{
				Level: c.level,
				Team1: team1,
				Team2: team2,
			})
			stage.MapBg = sceneutil.DrawCombatMap(c.level)
			game.G.State.CurrentStage = stage
			game.G.SceneManager.ChangeScene(NewTroopDelpoyController())
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

func (c *cardpickController) pickEnemyCard() {
	if len(c.enemyCardsPool) > 0 {
		enemySelectedCard := c.enemyCardsPool[len(c.enemyCardsPool)-1]
		c.enemyCardsPool = c.enemyCardsPool[:len(c.enemyCardsPool)-1]
		c.enemyCardsPicked = append(c.enemyCardsPicked, enemySelectedCard)
		c.enemyLabel.Label = fmt.Sprintf("Enemy picks %s", styles.Orange(enemySelectedCard.Info().Name))
	} else {
		c.enemyLabel.Label = styles.Orange("Tactics are scheduled!")
	}
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

	c.pickEnemyCard()

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
				Color: styles.ColorOrange.Color(),
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
