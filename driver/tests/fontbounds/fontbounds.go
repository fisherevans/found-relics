package main

import (
	"fmt"
	assets2 "found-relics/assets"
	"golang.org/x/image/font"
)

var f font.Face

func main() {
	assets2.Initialize()
	f = assets2.Fonts.TextLarge.Regular

	fmt.Printf("height: %d, caph: %d, xh: %d, des: %d, asc: %d\n",
		f.Metrics().Height.Round(),
		f.Metrics().CapHeight.Round(),
		f.Metrics().XHeight.Round(),
		f.Metrics().Descent.Round(),
		f.Metrics().Ascent.Round())

	info("o")
	info("oo")
	info("p")
	info("P")
	info("Pp")
	info("oooooopP")
}

func info(s string) {
	fb, a := font.BoundString(f, s)
	dx := fb.Max.X.Ceil() - fb.Min.X.Floor()
	dy := fb.Max.Y.Ceil() - fb.Min.Y.Floor()
	fmt.Printf("%-10s - min(%3d, %3d) max(%3d, %3d) delta(%3d, %3d) adv(%3d)\n", s, fb.Min.X.Floor(), fb.Min.Y.Floor(), fb.Max.X.Ceil(), fb.Max.Y.Ceil(), dx, dy, a.Round())
}
