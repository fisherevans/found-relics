package battle

import (
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/moves"
	"found-relics/pkg/state"
	"math/rand"
)

type PlayerController struct {
	selected int
}

func (c *PlayerController) Update(game state.Game, battle *combat.Battle, dt float64, elapsed combat.Time) {
	// selector player
	count := len(battle.PlayerTeam)
	if game.Controller().JustPressed(state.InputDown) {
		c.selected = (c.selected + 1) % count
	}
	if game.Controller().JustPressed(state.InputUp) {
		c.selected = (count + c.selected - 1) % count
	}
	// queue moves
	char := battle.PlayerTeam[c.selected]
	var triggered []*combat.Move
	if game.Controller().JustPressed(state.InputOpt1) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot1))
	}
	if game.Controller().JustPressed(state.InputOpt2) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot2))
	}
	if game.Controller().JustPressed(state.InputOpt3) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot3))
	}
	if game.Controller().JustPressed(state.InputOpt4) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot4))
	}
	for _, move := range triggered {
		if move == nil || char.MoveQueueTimeDepth+move.Duration.ToCombatTime() > char.Details.MaxMoveQueueDepth.ToCombatTime() {
			continue
		}
		targetId := rand.Intn(len(battle.OpponentTeam))
		battle.QueueMove(move, char, battle.OpponentTeam[targetId]) // TODO targeting
	}
	if game.Controller().JustPressed(state.InputBack) {
		battle.DequeueLastMove(char)
	}
}
