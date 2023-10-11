package game

import (
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type FoundRelics struct {
	currentState state.State
	lastTime     time.Time
}

func (g *FoundRelics) Update() error {
	if g.lastTime.IsZero() {
		g.lastTime = time.Now()
	}
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now
	g.currentState.Update(g, dt)
	return nil
}

func (g *FoundRelics) Draw(screen *ebiten.Image) {
	g.currentState.Draw(g, screen)
}

func (g *FoundRelics) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.DeviceScaleFactor()
	s = 1
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
