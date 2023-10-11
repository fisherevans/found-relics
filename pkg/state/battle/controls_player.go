package battle

import (
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/moves"
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math/rand"
)

type PlayerController struct {
	selected int
}

func (c *PlayerController) Update(game state.Game, battle *combat.Battle, dt float64, elapsed combat.Time) {
	// selector player
	count := len(battle.PlayerTeam)
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		c.selected = (c.selected + 1) % count
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		c.selected = (count + c.selected - 1) % count
	}
	// queue moves
	char := battle.PlayerTeam[c.selected]
	var triggered []*combat.Move
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot1))
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot2))
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot3))
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		triggered = append(triggered, moves.Get(char.Details.Moves.Slot4))
	}
	for _, move := range triggered {
		if move == nil || char.MoveQueueTimeDepth+move.Duration.ToCombatTime() > char.Details.MaxMoveQueueDepth.ToCombatTime() {
			continue
		}
		targetId := rand.Intn(len(battle.OpponentTeam))
		battle.QueueMove(move, char, battle.OpponentTeam[targetId]) // TODO targeting
	}
}