package state

import (
	"found-relics/pkg/rpg/combat"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game interface {
	EnterBattle(battle *combat.Battle)
	EnterSelector()
	EnterState(newState State)
	Controller() Controller
}

type State interface {
	Update(game Game, dt float64)
	Draw(game Game, screen *ebiten.Image)
}

type LoadableState interface {
	State
	Load()
}

type InputKey int

const (
	InputUp = iota
	InputDown
	InputLeft
	InputRight

	InputOpt1 // A
	InputOpt2 // B
	InputOpt3 // X
	InputOpt4 // Y

	InputAlt  // L Bumper
	InputBack // R Bumper

	InputInspect // Select
	InputMenu    // Start/Esc

	InputAliasSelect = InputOpt1
)

type Controller interface {
	JustPressed(k InputKey) bool
	IsPressed(k InputKey) bool
}
