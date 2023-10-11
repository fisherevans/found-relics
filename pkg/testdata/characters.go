package testdata

import (
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/moves"
)

var HankTheTank combat.Character
var ReluctantHealer combat.Character
var BadGuy combat.Character

func loadCharacters() {
	HankTheTank = combat.Character{
		Name: "Hank the Tank",
		Moves: combat.AvailableMoves{
			Slot1: moves.HeavySwing,
			Slot2: moves.Flourish,
			Slot3: moves.Sit,
		},
		MaxLife:           5000,
		MaxMoveQueueDepth: 15,
	}

	ReluctantHealer = combat.Character{
		Name: "Reluctant Healer",
		Moves: combat.AvailableMoves{
			Slot1: moves.Heal,
			Slot2: moves.BloodRitual,
			Slot3: moves.Sit,
			Slot4: moves.Smite,
		},
		MaxLife:           3000,
		MaxMoveQueueDepth: 10,
	}

	BadGuy = combat.Character{
		Name: "Bad Guy",
		Moves: combat.AvailableMoves{
			Slot1: moves.Eat,
			Slot2: moves.Poke,
			Slot3: moves.Sit,
		},
		MaxLife: 2000,
	}
}

func Battle2v1() *combat.Battle {
	return &combat.Battle{
		PlayerTeam: []*combat.BattleCharacter{
			HankTheTank.NewBattleCharacter(),
			ReluctantHealer.NewBattleCharacter(),
		},
		OpponentTeam: []*combat.BattleCharacter{
			BadGuy.NewBattleCharacter(),
		},
	}
}
