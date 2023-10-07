package battle

import (
	"found-relics/driver/foundrelics/state"
	"found-relics/pkg/rpg/combat"
	"math/rand"
)

type OpponentController struct {
}

func (c *OpponentController) Update(game state.Game, battle *combat.Battle, dt float64, elapsed combat.Time) {
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
