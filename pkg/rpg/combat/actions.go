package combat

import (
	"fmt"
)

type EventType int

const Damage EventType = 0
const Heal EventType = 1

type Event struct {
	Type   EventType
	Health int
	Source *BattleCharacter
	At     Time
}

func (b *Battle) Heal(amt int, source, target *BattleCharacter) {
	amt = minInt(amt, target.Details.MaxLife-target.Life)
	target.Life = target.Life + amt
	target.Events = append(target.Events, Event{
		Type:   Heal,
		Health: amt,
		Source: source,
		At:     b.ElapsedTime,
	})
}

func (b *Battle) Damage(amt int, source, target *BattleCharacter) {
	amt = minInt(amt, target.Life)
	target.Life = target.Life - amt
	target.Events = append(target.Events, Event{
		Type:   Damage,
		Health: amt,
		Source: source,
		At:     b.ElapsedTime,
	})
}

func (b *Battle) QueueMove(move *Move, source, target *BattleCharacter) {
	mi := &MoveInstance{
		Move:   move,
		Source: source,
		Target: target,
	}
	source.MoveQueue = append(source.MoveQueue, mi)
	source.MoveQueueTimeDepth += mi.Move.Duration.ToCombatTime()
	fmt.Printf("move '%s' queued for '%s'\n", move.Name, source.Details.Name)
}

func (b *Battle) DequeueLastMove(char *BattleCharacter) bool {
	l := len(char.MoveQueue)
	if l < 2 {
		return false
	}
	move := char.MoveQueue[l-1]
	char.MoveQueue = char.MoveQueue[:l-1]
	char.MoveQueueTimeDepth -= move.Move.Duration.ToCombatTime()
	return true
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func absInt(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
