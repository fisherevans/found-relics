package state

import (
	"found-relics/pkg/rpg/combat"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game interface {
	EnterBattle(battle *combat.Battle)
	EnterSelector()
	EnterState(newState State)
}

type State interface {
	Update(game Game, dt float64)
	Draw(game Game, screen *ebiten.Image)
}

type LoadableState interface {
	State
	Load()
}
