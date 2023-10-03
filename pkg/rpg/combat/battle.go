package combat

import "fmt"

type Battle struct {
	ElapsedTime  Time
	PlayerTeam   []*BattleCharacter
	OpponentTeam []*BattleCharacter
}

func (b *Battle) Progress(elapse Time) Time {
	// determine how much time to elapse
	for _, c := range b.PlayerTeam {
		elapse = MinTime(elapse, c.MoveQueueTimeDepth)
	}
	for _, c := range b.OpponentTeam {
		elapse = MinTime(elapse, c.MoveQueueTimeDepth)
	}
	// collect all elapsed move triggers
	triggers := &ElapsedTriggers{}
	for _, c := range b.PlayerTeam {
		c.Progress(elapse, triggers)
	}
	for _, c := range b.OpponentTeam {
		c.Progress(elapse, triggers)
	}
	// execute triggers in order
	startElapse := b.ElapsedTime
	endElapse := b.ElapsedTime + elapse
	for triggers.First != nil {
		b.ElapsedTime = startElapse + triggers.First.RelativeAt
		if t := triggers.First; t.Trigger != nil {
			fmt.Printf("running move trigger for '%s' during '%s'\n", triggers.First.Instance.Source.Details.Name, triggers.First.Instance.Move.Name)
			t.Trigger(b, t.Instance)
		}
		triggers.First = triggers.First.Next
	}
	// record progress time
	b.ElapsedTime = endElapse
	return elapse
}

type ElapsedTriggers struct {
	First *TriggerNode
}

type TriggerNode struct {
	RelativeAt Time
	Instance   *MoveInstance
	Trigger    MoveTrigger
	Next       *TriggerNode
}

func (p *ElapsedTriggers) Append(new *TriggerNode) {
	if p.First == nil { // if nil, start the list
		p.First = new
		return
	}
	if new.RelativeAt < p.First.RelativeAt { // if we need to replace the first one
		new.Next = p.First
		p.First = new
		return
	}
	t := p.First
	for t.Next != nil { // for each node, if new goes betwen t and next, insert and return
		if t.Next.RelativeAt > new.RelativeAt {
			new.Next = t.Next
			t.Next = new
			return
		}
		t = t.Next
	}
	t.Next = new // else, append to the end
}
