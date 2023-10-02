package moves

import (
	"combat/pkg/rpg/combat"
	"math/rand"
)

var SitThere = &combat.Move{
	Name:     "Sit",
	Duration: 1,
}

var Poke = &combat.Move{
	Name:     "Poke",
	Duration: 2,
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Damage(50)
	},
}

var Eat = &combat.Move{
	Name:     "Eat",
	Duration: 6,
	Ticks:    GenTicks(0, 1, 5),
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Source.Heal(25)
	},
}

var HealSome = &combat.Move{
	Name:     "Heal",
	Duration: 3,
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Heal(rand.Intn(300) + 100)
	},
}

var BloodRitual = &combat.Move{
	Name:     "Blood Ritual",
	Duration: 6,
	Ticks:    GenTicks(0, 1, 11),
	StartTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Source.Damage(500)
	},
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Heal(100)
	},
}

var HeavySwing = &combat.Move{
	Name:     "Heavy Swing",
	Duration: 8,

	Ticks: GenTicks(0, 6, 1),
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Damage(rand.Intn(1000) + 500)
	},
}

var Flourish = &combat.Move{
	Name:     "Flourish",
	Duration: 4,
	Ticks:    GenTicks(1, 1, 3),
	StartTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Damage(100)
	},
	TickTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Damage(50)
	},
	EndTrigger: func(battle *combat.Battle, move *combat.MoveInstance) {
		move.Target.Damage(200)
	},
}

func GenTicks(additionalDelay, duration combat.Time, count int) []combat.Time {
	var ticks []combat.Time
	for i := 1; i <= count; i++ {
		ticks = append(ticks, additionalDelay+combat.Time(i)*duration)
	}
	//fmt.Printf("(%d, %d, %d) = %s\n", additionalDelay, duration, count, ticks)
	return ticks
}
