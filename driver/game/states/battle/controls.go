package battle

import (
	"found-relics/pkg/rpg/combat"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

type PlayerControls struct {
	selected int
}

func (p *PlayerControls) Update(win *pixelgl.Window, battle *combat.Battle) {
	// select player
	count := len(battle.PlayerTeam)
	if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyS) {
		p.selected = (p.selected + 1) % count
	}
	if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyW) {
		p.selected = (count + p.selected - 1) % count
	}
	// queue moves
	char := battle.PlayerTeam[p.selected]
	var moves []*combat.Move
	if win.JustPressed(pixelgl.Key1) {
		moves = append(moves, char.Details.Moves.Slot1)
	}
	if win.JustPressed(pixelgl.Key2) {
		moves = append(moves, char.Details.Moves.Slot2)
	}
	if win.JustPressed(pixelgl.Key3) {
		moves = append(moves, char.Details.Moves.Slot3)
	}
	if win.JustPressed(pixelgl.Key4) {
		moves = append(moves, char.Details.Moves.Slot4)
	}
	for _, move := range moves {
		if move == nil || char.MoveQueueTimeDepth+move.Duration.ToCombatTime() > char.Details.MaxMoveQueueDepth.ToCombatTime() {
			continue
		}
		targetId := rand.Intn(len(battle.OpponentTeam))
		battle.QueueMove(move, char, battle.OpponentTeam[targetId]) // TODO targeting
	}
}

type OpponentControls struct {
}

func (p *OpponentControls) Update(win *pixelgl.Window, battle *combat.Battle) {
	for battle.OpponentTeam[0].MoveQueueTimeDepth < combat.Beats(6).ToCombatTime() {
		var move *combat.Move
		//switch rand.Intn(3) {
		switch 1 {
		case 0:
			move = battle.OpponentTeam[0].Details.Moves.Slot1
			break
		case 1:
			move = battle.OpponentTeam[0].Details.Moves.Slot2
			break
		default:
			move = battle.OpponentTeam[0].Details.Moves.Slot3
		}
		targetId := rand.Intn(len(battle.PlayerTeam))
		battle.QueueMove(move, battle.OpponentTeam[0], battle.PlayerTeam[targetId]) // todo actual targeting
	}
}
