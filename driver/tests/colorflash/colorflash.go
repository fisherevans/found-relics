package main

import (
	"fmt"
	"found-relics/pkg/style"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lucasb-eyer/go-colorful"
	"log"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Color Flash")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	t := &Test{
		A:    style.ColorHighlightBright,
		B:    style.ColorHighlightDark,
		Rate: 1,
	}
	fmt.Printf("A = r:%.3f, g:%.3f, b:%.3f\n", t.A.R, t.A.G, t.A.B)
	fmt.Printf("B = r:%.3f, g:%.3f, b:%.3f\n", t.B.R, t.B.G, t.B.B)
	if err := ebiten.RunGame(t); err != nil {
		log.Fatal(err)
	}
}

type Test struct {
	A    colorful.Color
	B    colorful.Color
	Rate float64
	Time float64
}

func (t *Test) Update() error {
	t.Time += 1.0 / ebiten.ActualTPS()
	return nil
}

func (t *Test) Draw(screen *ebiten.Image) {
	c := style.Flash(t.A, t.B, t.Time, t.Rate)
	ebitenutil.DrawRect(screen, 10, 10, 200, 200, c)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("time:%3.3f\nr:%.3f, g:%.3f, b:%.3f\n", t.Time, c.R, c.G, c.B))
}

func (t *Test) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
