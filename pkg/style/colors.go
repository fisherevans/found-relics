package style

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"log"
	"math"
)

var ColorHealth = hex("#f94e6b")
var ColorHealthLagDown = hex("#a23346")
var ColorHealthLagUp = hex("#ff8196")
var ColorHealthBg = hex("#522c32")

var ColorHighlightBright = colorful.Hsl(211, 0.9, 0.75)
var ColorHighlightDark = colorful.Hsl(211, 0.9, 0.65)

var ColorBright1 = hex("#fafafa")
var ColorBright2 = hex("#d9d9d9")
var ColorBright3 = hex("#b3b3b3")
var ColorGray = hex("#808080")
var ColorDark3 = hex("#4d4d4d")
var ColorDark2 = hex("#262626")
var ColorDark1 = hex("#050505")

var ColorPurpleBright = hex("#ad4bde")
var ColorPurpleDark = hex("#9237bf")

var FullyTransparent = color.RGBA{0, 0, 0, 0}

func Flash(a, b colorful.Color, time, rate float64) colorful.Color {
	at := math.Sin(time*math.Pi*rate)/2.0 + 0.5
	return a.BlendLab(b, at)
}

func Transparent(c color.Color, alpha float64) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(alpha * 255),
	}
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func hex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		log.Fatal("invalid color "+s, err)
	}
	return c
}
