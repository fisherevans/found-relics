package combat

import (
	"math"
)

const healthLagPctPerBeat = 0.05

type BattleCharacter struct {
	Details Character
	Life    int
	LagLife int

	MoveQueue          []*MoveInstance
	MoveQueueTimeDepth Time

	Events []Event
}

func (c *BattleCharacter) Progress(toElapse Time, triggers *ElapsedTriggers) {
	// lag health
	if c.Life != c.LagLife {
		dt := float64(toElapse) / float64(BeatToTimeScale)
		lagStep := int(math.Ceil(healthLagPctPerBeat * dt * float64(c.Details.MaxLife)))
		if c.LagLife > c.Life {
			lagStep = -1 * lagStep
		}
		//fmt.Printf("life: %4d, lag:%4d, dt: %.4f, step:%4d\n", c.Life, c.LagLife, dt, lagStep)
		if absInt(lagStep) > absInt(c.LagLife-c.Life) {
			c.LagLife = c.Life
		} else {
			c.LagLife += lagStep
		}
	}

	// queue move triggers
	var pastElapsed Time = 0
	for len(c.MoveQueue) > 0 && toElapse > 0 {
		m := c.MoveQueue[0]
		if m.ElapsedTime == 0 {
			triggers.Append(&TriggerNode{
				RelativeAt: pastElapsed,
				Instance:   m,
				Trigger:    m.Move.StartTrigger,
			})
		}
		start := m.ElapsedTime
		end := MinTime(m.Move.Duration.ToCombatTime(), m.ElapsedTime+toElapse)
		thisElapsed := end - start
		m.ElapsedTime += thisElapsed
		for _, tick := range m.Move.Ticks {
			if tick > start && tick <= end {
				triggers.Append(&TriggerNode{
					RelativeAt: pastElapsed + tick - start,
					Instance:   m,
					Trigger:    m.Move.TickTrigger,
				})
			}
		}
		if m.ElapsedTime >= m.Move.Duration.ToCombatTime() {
			triggers.Append(&TriggerNode{
				RelativeAt: pastElapsed + m.Move.Duration.ToCombatTime(),
				Instance:   m,
				Trigger:    m.Move.EndTrigger,
			})
			c.MoveQueue = c.MoveQueue[1:]
		}
		c.MoveQueueTimeDepth -= thisElapsed
		pastElapsed += thisElapsed
		toElapse -= thisElapsed
	}
}

type Character struct {
	Name              string
	Moves             AvailableMoves
	MaxMoveQueueDepth Beats
	MaxLife           int
}

func (c Character) NewBattleCharacter() *BattleCharacter {
	return &BattleCharacter{
		Details: c,
		Life:    c.MaxLife,
		LagLife: c.MaxLife,
	}
}
