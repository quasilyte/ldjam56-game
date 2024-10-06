package styles

import (
	graphics "github.com/quasilyte/ebitengine-graphics"
)

var (
	ColorBackground = graphics.RGB(0x6e579c).ScaleRGB(0.25)
	ColorDark       = graphics.RGB(0x19352c)
	ColorNormal     = graphics.RGB(0x579c6e)
	ColorBright     = graphics.RGB(0xa6fba6)
	ColorOrange     = graphics.RGB(0xdc5b24)
)

func Orange(s string) string {
	return BB(ColorOrange, s)
}

func Background(s string) string {
	return BB(ColorBackground, s)
}

func Dark(s string) string {
	return BB(ColorDark, s)
}

func Normal(s string) string {
	return BB(ColorNormal, s)
}

func Bright(s string) string {
	return BB(ColorBright, s)
}

func BB(clr graphics.ColorScale, s string) string {
	rgb := graphics.FormatRGB(clr.Color())
	return "[color=" + rgb + "]" + s + "[/color]"
}
