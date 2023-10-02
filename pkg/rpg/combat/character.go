package combat

import (
	"fmt"
)

type CharacterDetails struct {
	Name              string
	Moves             AvailableMoves
	MaxMoveQueueDepth Beats
	MaxLife           int
}

type CharacterInstance struct {
	Details CharacterDetails
	Life    int

	MoveQueue          []*MoveInstance
	MoveQueueTimeDepth Time
}

func (c *CharacterInstance) Heal(amt int) {
	c.Life = c.Life + amt
	if c.Life > c.Details.MaxLife {
		c.Life = c.Details.MaxLife
	}
}

func (c *CharacterInstance) Damage(amt int) {
	c.Life = c.Life - amt
	if c.Life < 0 {
		c.Life = 0
	}
}

func (c *CharacterInstance) QueueMove(move *Move, target *CharacterInstance) {
	mi := &MoveInstance{
		Move:   move,
		Source: c,
		Target: target,
	}
	c.MoveQueue = append(c.MoveQueue, mi)
	c.MoveQueueTimeDepth += mi.Move.Duration.ToCombatTime()
	fmt.Printf("move '%s' queued for '%s'\n", move.Name, c.Details.Name)
}

func (c *CharacterInstance) Progress(toElapse Time, triggers *ElapsedTriggers) {
	var pastElapsed Time = 0
	for len(c.MoveQueue) > 0 && toElapse > 0 {
		m := c.MoveQueue[0]
		if m.ElapsedTime == 0 {
			triggers.Append(&TriggerNode{
				Order:    int(pastElapsed),
				Instance: m,
				Trigger:  m.Move.StartTrigger,
			})
		}
		start := m.ElapsedTime
		end := MinTime(m.Move.Duration.ToCombatTime(), m.ElapsedTime+toElapse)
		thisElapsed := end - start
		m.ElapsedTime += thisElapsed
		for _, tick := range m.Move.Ticks {
			if tick > start && tick <= end {
				triggers.Append(&TriggerNode{
					Order:    int(pastElapsed + tick - start),
					Instance: m,
					Trigger:  m.Move.TickTrigger,
				})
			}
		}
		if m.ElapsedTime >= m.Move.Duration.ToCombatTime() {
			triggers.Append(&TriggerNode{
				Order:    int(pastElapsed + m.Move.Duration.ToCombatTime()),
				Instance: m,
				Trigger:  m.Move.EndTrigger,
			})
			c.MoveQueue = c.MoveQueue[1:]
		}
		c.MoveQueueTimeDepth -= thisElapsed
		pastElapsed += thisElapsed
		toElapse -= thisElapsed
	}
}
