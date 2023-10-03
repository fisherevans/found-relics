package main

import (
	"fmt"
	"found-relics/driver/game"
	"found-relics/pkg/style"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(game.NewGame(&Test{
		A:    style.ColorHighlightBright,
		B:    style.ColorHighlightDark,
		Rate: 1,
	}).Run)
}

type Test struct {
	A    pixel.RGBA
	B    pixel.RGBA
	Rate float64
	Time float64
}

func (t *Test) Init(game *game.Game, win *pixelgl.Window) {
	fmt.Printf("A = r:%.3f, g:%.3f, b:%.3f\n", t.A.R, t.A.G, t.A.B)
	fmt.Printf("B = r:%.3f, g:%.3f, b:%.3f\n", t.B.R, t.B.G, t.B.B)
}

func (t *Test) Tick(game *game.Game, win *pixelgl.Window, dt float64) {
	t.Time += dt
	c := style.Flash(t.A, t.B, t.Time, t.Rate)
	win.Clear(c)
	if t.Time < 3 {
		fmt.Printf("%3.3f - ", t.Time)
		fmt.Printf("r:%.3f, g:%.3f, b:%.3f\n", c.R, c.G, c.B)
	}

}
