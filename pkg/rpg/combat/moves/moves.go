package moves

import (
	"found-relics/pkg/assets"
	"found-relics/pkg/rpg/combat"
	"math/rand"
)

var SitThere = &combat.Move{
	Name:     "Sit",
	Sprite:   assets.BattleIconReference(6, 1, 6),
	Duration: 1,
}

var Poke = &combat.Move{
	Name:     "Poke",
	Sprite:   assets.BattleIconReference(4, 6, 8),
	Duration: 2,
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Damage(500, move.Source, move.Target)
	},
}

var Eat = &combat.Move{
	Name:     "Eat",
	Sprite:   assets.BattleIconReference(4, 5, 10),
	Duration: 6,
	Ticks:    []combat.Time{1, 2, 3, 4},
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Heal(25, move.Source, move.Target)
	},
}

var Heal = &combat.Move{
	Name:     "Heal",
	Sprite:   assets.BattleIconReference(3, 6, 9),
	Duration: 3,
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Heal(rand.Intn(300)+100, move.Source, move.Target)
	},
}

var BloodRitual = &combat.Move{
	Name:     "Blood Ritual",
	Sprite:   assets.BattleIconReference(8, 1, 9),
	Duration: 6,
	Ticks:    []combat.Time{1, 2, 3, 4, 5},
	StartTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Damage(500, move.Source, move.Target)
	},
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Heal(100, move.Source, move.Target)
	},
}

var Smite = &combat.Move{
	Name:     "Smite",
	Sprite:   assets.BattleIconReference(4, 4, 1),
	Duration: 4,
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Damage(rand.Intn(500)+100, move.Source, move.Target)
	},
}

var HeavySwing = &combat.Move{
	Name:     "Heavy Swing",
	Sprite:   assets.BattleIconReference(3, 4, 10),
	Duration: 8,
	Ticks:    []combat.Time{5},
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		battle.Damage(rand.Intn(1000)+500, move.Source, move.Target)
	},
}

var Flourish = &combat.Move{
	Name:     "Flourish",
	Sprite:   assets.BattleIconReference(8, 4, 1),
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
}
