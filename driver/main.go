package main

import (
	"combat/driver/game"
	"combat/driver/game/states/battle"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	g := game.NewGame(battle.NewExampleBattles())
	pixelgl.Run(g.Run)
}
