package game

import (
	"github.com/faiface/pixel/pixelgl"
)

type State interface {
	Init(game *Game, win *pixelgl.Window)
	Tick(game *Game, win *pixelgl.Window, dt float64)
}
