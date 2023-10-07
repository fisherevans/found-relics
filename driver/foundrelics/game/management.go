package game

import (
	"found-relics/driver/foundrelics/state"
	"found-relics/driver/foundrelics/state/battle"
	"found-relics/driver/foundrelics/state/load"
	"found-relics/driver/foundrelics/state/selector"
	"found-relics/pkg/rpg/combat"
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
