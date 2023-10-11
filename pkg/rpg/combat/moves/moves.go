package moves

import (
	"found-relics/pkg/rpg/combat"
)

var All map[combat.MoveId]*combat.Move = map[combat.MoveId]*combat.Move{}

const (
	None        combat.MoveId = ""
	Sit         combat.MoveId = "Sit"
	Poke        combat.MoveId = "Poke"
	Eat         combat.MoveId = "Eat"
	Heal        combat.MoveId = "Heal"
	BloodRitual combat.MoveId = "Blood Ritual"
	Smite       combat.MoveId = "Smite"
	HeavySwing  combat.MoveId = "Heavy Swing"
	Flourish    combat.MoveId = "Flourish"
)

func Get(moveId combat.MoveId) *combat.Move {
	m, ok := All[moveId]
	if !ok {
		return nil
	}
	return m
}
