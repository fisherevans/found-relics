package combat

import "found-relics/pkg/assets"

const BeatToTimeScale = 1000

type Beats int

func (b Beats) ToCombatTime() Time {
	return Time(b * BeatToTimeScale)
}

type Time int

func (t Time) ToBeatRoundedDown() Beats {
	return Beats(t / BeatToTimeScale)
}

func (t Time) ToBeatRoundedUp() Beats {
	quot := t.ToBeatRoundedDown()
	if t%BeatToTimeScale != 0 {
		quot += 1
	}
	return quot
}

type Move struct {
	Name     string
	Sprite   assets.SpriteReference
	Duration Beats

	Ticks []Time

	StartTrigger MoveTrigger
	TickTrigger  MoveTrigger
	EndTrigger   MoveTrigger
}

type AvailableMoves struct {
	Slot1 *Move
	Slot2 *Move
	Slot3 *Move
	Slot4 *Move
}

func (a AvailableMoves) AsSlice() []*Move {
	return []*Move{a.Slot1, a.Slot2, a.Slot3, a.Slot4}
}

type MoveInstance struct {
	Move        *Move
	Source      *BattleCharacter
	Target      *BattleCharacter
	ElapsedTime Time
}

type MoveTrigger func(battle *Battle, move *MoveInstance)
