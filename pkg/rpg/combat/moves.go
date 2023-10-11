package combat

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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

type MoveId string

type Move struct {
	Name     MoveId
	Sprite   *ebiten.Image
	Duration Beats

	Ticks []Time

	StartTrigger MoveTrigger
	TickTrigger  MoveTrigger
	EndTrigger   MoveTrigger
}

type AvailableMoves struct {
	Slot1 MoveId
	Slot2 MoveId
	Slot3 MoveId
	Slot4 MoveId
}

func (a AvailableMoves) AsSlice() []MoveId {
	return []MoveId{a.Slot1, a.Slot2, a.Slot3, a.Slot4}
}

type MoveInstance struct {
	Move        *Move
	Source      *BattleCharacter
	Target      *BattleCharacter
	ElapsedTime Time
}

type MoveTrigger func(battle *Battle, move *MoveInstance)
