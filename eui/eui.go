package eui

import (
	"image/color"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	sound "github.com/quasilyte/ebitengine-sound"
	"github.com/quasilyte/gscene"
	"github.com/quasilyte/ldjam56-game/assets"
	"github.com/quasilyte/ldjam56-game/styles"
	"golang.org/x/image/font"
)

type Widget = widget.PreferredSizeLocateableWidget

type Builder struct {
	button     *buttonDefaults
	tileButton *buttonDefaults
	panel      *panelDefaults

	currentObject *SceneObject

	loader *resource.Loader

	audio *sound.System
}

type buttonDefaults struct {
	image      *widget.ButtonImage
	padding    widget.Insets
	textColors *widget.ButtonTextColor
}

type panelDefaults struct {
	image   *image.NineSlice
	padding widget.Insets
}

type Config struct {
	Loader *resource.Loader

	Audio *sound.System
}

func NewBuilder(config Config) *Builder {
	b := &Builder{
		loader: config.Loader,
		audio:  config.Audio,
	}
	return b
}

func (b *Builder) Init() {
	l := b.loader

	{
		disabled := loadNineSliced(l, assets.ImageUIButtonDisabled, 26, 20)
		idle := loadNineSliced(l, assets.ImageUIButtonIdle, 26, 20)
		hover := loadNineSliced(l, assets.ImageUIButtonHover, 26, 20)
		pressed := loadNineSliced(l, assets.ImageUIButtonPressed, 26, 20)
		buttonPadding := widget.Insets{
			Left:   12,
			Right:  12,
			Top:    6,
			Bottom: 6,
		}
		b.button = &buttonDefaults{
			image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  pressed,
				Disabled: disabled,
			},
			padding: buttonPadding,
			textColors: &widget.ButtonTextColor{
				Idle:     styles.ColorBright.Color(),
				Disabled: styles.ColorNormal.Color(),
			},
		}
	}

	{
		disabled := loadNineSliced(l, assets.ImageUITileButtonDisabled, 26, 20)
		idle := loadNineSliced(l, assets.ImageUITileButtonIdle, 26, 20)
		hover := loadNineSliced(l, assets.ImageUITileButtonHover, 26, 20)
		pressed := loadNineSliced(l, assets.ImageUITileButtonPressed, 26, 20)
		buttonPadding := widget.Insets{
			Left:   0,
			Right:  0,
			Top:    0,
			Bottom: 0,
		}
		b.tileButton = &buttonDefaults{
			image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  pressed,
				Disabled: disabled,
			},
			padding: buttonPadding,
			textColors: &widget.ButtonTextColor{
				Idle:     styles.ColorBright.Color(),
				Disabled: styles.ColorNormal.Color(),
			},
		}
	}

	{
		normal := loadNineSliced(l, assets.ImageUIPanel, 18, 18)
		b.panel = &panelDefaults{
			image: normal,
			padding: widget.Insets{
				Left:   12,
				Right:  12,
				Top:    12,
				Bottom: 12,
			},
		}
	}
}

type ButtonConfig struct {
	Text         string
	OnClick      func()
	OnMouseEnter func()
	OnMouseExit  func()
	MinWidth     int
	MinHeight    int
	Font         font.Face
	LayoutData   any
}

func (b *Builder) NewButton(config ButtonConfig) *widget.Button {
	ff := config.Font
	if ff == nil {
		ff = assets.Font1
	}

	defaults := b.button

	colors := b.button.textColors
	padding := defaults.padding
	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(defaults.image),
		widget.ButtonOpts.Text(config.Text, ff, colors),
		widget.ButtonOpts.TextPadding(padding),
	}

	if strings.Contains(config.Text, "[color=") {
		options = append(options, widget.ButtonOpts.TextProcessBBCode(true))
	}

	options = append(options, widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
		if config.OnClick != nil {
			config.OnClick()
		}
		b.audio.PlaySound(assets.AudioButtonClick)
	}))

	if config.MinWidth != 0 || config.MinHeight != 0 {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, config.MinHeight)))
	}

	buttonWidget := widget.NewButton(options...)
	return buttonWidget
}

type TileButtonConfig struct {
	Text       string
	MinWidth   int
	MinHeight  int
	OnClick    func()
	LayoutData any
}

func (b *Builder) NewTileButton(config TileButtonConfig) *widget.Button {
	defaults := b.tileButton

	padding := defaults.padding
	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(defaults.image),
	}

	options = append(options, widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
		if config.OnClick != nil {
			config.OnClick()
		}
		b.audio.PlaySound(assets.AudioButtonClickSoft)
	}))
	if config.Text != "" {
		options = append(options,
			widget.ButtonOpts.Text(config.Text, assets.FontTiny, b.button.textColors),
			widget.ButtonOpts.TextPadding(padding))
		if strings.Contains(config.Text, "[color=") {
			options = append(options, widget.ButtonOpts.TextProcessBBCode(true))
		}
	}

	if config.MinWidth != 0 || config.MinHeight != 0 {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, config.MinHeight)))
	}

	buttonWidget := widget.NewButton(options...)
	return buttonWidget
}

type TextConfig struct {
	Text     string
	Font     font.Face
	Color    color.Color
	MinWidth int
	MaxWidth int

	LayoutData any

	AlignLeft   bool
	ForceBBCode bool
	AlignRight  bool
	AlignTop    bool
}

func (b *Builder) NewText(config TextConfig) *widget.Text {
	var clr color.Color = styles.ColorBright.Color()
	if config.Color != nil {
		clr = config.Color
	}
	ff := assets.Font1
	if config.Font != nil {
		ff = config.Font
	}

	verticalPos := widget.TextPositionCenter
	if config.AlignTop {
		verticalPos = widget.TextPositionStart
	}

	opts := []widget.TextOpt{
		widget.TextOpts.Text(config.Text, ff, clr),
	}
	if config.LayoutData != nil {
		opts = append(opts, widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(config.LayoutData)))
	}
	if config.MinWidth != 0 {
		opts = append(opts, widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}
	if config.MaxWidth != 0 {
		opts = append(opts, widget.TextOpts.MaxWidth(float64(config.MaxWidth)))
	}
	switch {
	case config.AlignLeft:
		opts = append(opts, widget.TextOpts.Position(widget.TextPositionStart, verticalPos))
	case config.AlignRight:
		opts = append(opts, widget.TextOpts.Position(widget.TextPositionEnd, verticalPos))
	default:
		opts = append(opts, widget.TextOpts.Position(widget.TextPositionCenter, verticalPos))
	}
	if config.ForceBBCode || strings.Contains(config.Text, "[color=") {
		opts = append(opts, widget.TextOpts.ProcessBBCode(true))
	}
	return widget.NewText(opts...)
}

type PanelConfig struct {
	MinWidth   int
	MinHeight  int
	Padding    *widget.Insets
	LayoutData any
}

func (b *Builder) NewPanel(config PanelConfig) *widget.Container {
	defaults := b.panel

	padding := defaults.padding
	if config.Padding != nil {
		padding = *config.Padding
	}

	var ld any
	if config.LayoutData != nil {
		ld = config.LayoutData
	} else {
		ld = widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		}
	}

	opts := []widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(defaults.image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(padding),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(ld),
			widget.WidgetOpts.MinSize(config.MinWidth, config.MinHeight),
		),
	}

	panel := widget.NewContainer(opts...)

	return panel
}

func (b *Builder) Build(scene *gscene.Scene, root *widget.Container) *ebitenui.UI {
	anchor := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	anchor.AddChild(root)

	uiObject := b.newSceneObject(anchor)
	scene.AddGraphics(uiObject, 0)
	scene.AddObject(uiObject)

	return uiObject.ui
}

func loadNineSliced(l *resource.Loader, id resource.ImageID, offsetX, offsetY int) *image.NineSlice {
	i := l.LoadImage(id).Data
	return nineSliceImage(i, offsetX, offsetY)
}

func nineSliceImage(i *ebiten.Image, offsetX, offsetY int) *image.NineSlice {
	size := i.Bounds().Size()
	w := size.X
	h := size.Y
	return image.NewNineSlice(i,
		[3]int{offsetX, w - 2*offsetX, offsetX},
		[3]int{offsetY, h - 2*offsetY, offsetY},
	)
}
