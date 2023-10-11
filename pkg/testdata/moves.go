package testdata

import (
	"found-relics/assets"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/moves"
	"math/rand"
)

func loadMoves() {
	add := func(m *combat.Move) {
		moves.All[m.Name] = m
	}
	add(&combat.Move{
		Name:     moves.Sit,
		Sprite:   assets.Images.MoveIconSheets[5].Sprite(1, 6),
		Duration: 1,
	})
	add(&combat.Move{
		Name:     moves.Poke,
		Sprite:   assets.Images.MoveIconSheets[3].Sprite(6, 8),
		Duration: 2,
		EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(500, move.Source, move.Target)
		},
	})
	add(&combat.Move{
		Name:     moves.Eat,
		Sprite:   assets.Images.MoveIconSheets[3].Sprite(5, 10),
		Duration: 6,
		Ticks:    []combat.Time{1, 2, 3, 4},
		TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Heal(25, move.Source, move.Target)
		},
	})
	add(&combat.Move{
		Name:     moves.Heal,
		Sprite:   assets.Images.MoveIconSheets[2].Sprite(6, 9),
		Duration: 3,
		EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Heal(rand.Intn(300)+100, move.Source, move.Source) // TODO target
		},
	})
	add(&combat.Move{
		Name:     moves.BloodRitual,
		Sprite:   assets.Images.MoveIconSheets[7].Sprite(1, 9),
		Duration: 6,
		Ticks:    []combat.Time{1, 2, 3, 4, 5},
		StartTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(500, move.Source, move.Target)
		},
		TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Heal(100, move.Source, move.Target)
		},
	})
	add(&combat.Move{
		Name:     moves.Smite,
		Sprite:   assets.Images.MoveIconSheets[3].Sprite(4, 1),
		Duration: 4,
		EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(rand.Intn(500)+100, move.Source, move.Target)
		},
	})
	add(&combat.Move{
		Name:     moves.HeavySwing,
		Sprite:   assets.Images.MoveIconSheets[2].Sprite(4, 10),
		Duration: 8,
		Ticks:    []combat.Time{5},
		TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(rand.Intn(1000)+500, move.Source, move.Target)
		},
	})
	add(&combat.Move{
		Name:     moves.Flourish,
		Sprite:   assets.Images.MoveIconSheets[7].Sprite(4, 1),
		Duration: 4,
		Ticks:    []combat.Time{1, 2},
		StartTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(100, move.Source, move.Target)
		},
		TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(50, move.Source, move.Target)
		},
		EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
			battle.Damage(200, move.Source, move.Target)
		},
	})

}
