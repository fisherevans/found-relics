package main

import (
	"found-relics/driver/game"
	"found-relics/driver/game/states/initialize"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	g := game.NewGame(initialize.New())
	pixelgl.Run(g.Run)
}
