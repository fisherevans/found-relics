package game

import (
	"found-relics/pkg/assets"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

type Game struct {
	Resources assets.Resources

	currentState            State
	currentStateInitialized bool
}

func NewGame(initialState State) *Game {
	return &Game{
		currentState: initialState,
	}
}

func (g *Game) Run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Game Prototype",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	targetFps := time.Tick(time.Second / 120)
	last := time.Now()
	for !win.Closed() {
		if !g.currentStateInitialized {
			g.currentState.Init(g, win)
			g.currentStateInitialized = true
		}
		now := time.Now()
		dt := now.Sub(last).Seconds()
		last = now
		g.currentState.Tick(g, win, dt)
		win.Update()
		<-targetFps
	}
}

func (g *Game) SwapState(state State) {
	g.currentState = state
	g.currentStateInitialized = false
	// TODO other events, like leave and close
}
