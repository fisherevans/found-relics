package game

import (
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type FoundRelics struct {
	currentState state.State
	lastTime     time.Time
	controller   state.Controller
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
