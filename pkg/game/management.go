package game

import (
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/state"
	"found-relics/pkg/state/battle"
	"found-relics/pkg/state/load"
	"found-relics/pkg/state/selector"
)

func (g *FoundRelics) EnterBattle(b *combat.Battle) {
	g.currentState = load.LoadState(battle.NewBattle(b))
}

func (g *FoundRelics) EnterSelector() {
	g.currentState = selector.NewExampleBattles()
}

func (g *FoundRelics) EnterState(newState state.State) {
	g.currentState = newState
}

func (g *FoundRelics) Controller() state.Controller {
	return g.controller
}
