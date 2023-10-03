package style

import (
	"fmt"
	"github.com/PerformLine/go-stockutil/colorutil"
	"github.com/faiface/pixel"
	"image/color"
	"math"
)

var ColorHealth = hex("#f94e6b")
var ColorHealthLag = hex("#ff96e6")
var ColorHealthBg = hex("#522c32")

var ColorHighlightBright = hsl(211, 0.9, 0.75, 1)
var ColorHighlightDark = hsl(211, 0.9, 0.65, 1)

var ColorBright1 = hex("#fafafa")
var ColorBright2 = hex("#d9d9d9")
var ColorBright3 = hex("#b3b3b3")
var ColorGray = hex("#808080")
var ColorDark3 = hex("#4d4d4d")
var ColorDark2 = hex("#262626")
var ColorDark1 = hex("#050505")

var ColorPurpleBright = hex("#ad4bde")
var ColorPurpleDark = hex("#9237bf")

var Transparent = pixel.RGBA{}

func Flash(a, b pixel.RGBA, time, rate float64) pixel.RGBA {
	bs := math.Sin(time*math.Pi*rate)/2.0 + 0.5
	as := 1 - bs
	color := pixel.RGBA{
		R: a.R*as + b.R*bs,
		G: a.G*as + b.G*bs,
		B: a.B*as + b.B*bs,
		A: a.A*as + b.A*bs,
	}
	return color
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func hex(s string) pixel.RGBA {
	c := color.RGBA{A: 0xff}
	switch len(s) {
	case 7:
		_, _ = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, _ = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	}
	return pixel.ToRGBA(c)
}

func hsl(hue, saturation, lightness, alpha float64) pixel.RGBA {
	r, g, b := colorutil.HslToRgb(hue, saturation, lightness)
	c := color.RGBA{R: r, G: g, B: b}
	pc := pixel.ToRGBA(c)
	pc.A = alpha
	return pc
}
