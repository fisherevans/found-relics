package combat

import "fmt"

type Battle struct {
	ElapsedTime  Time
	PlayerTeam   []*CharacterInstance
	OpponentTeam []*CharacterInstance
}

func (b *Battle) Progress(elapse Time) Time {
	for _, c := range b.PlayerTeam {
		elapse = MinTime(elapse, c.MoveQueueTimeDepth)
	}
	for _, c := range b.OpponentTeam {
		elapse = MinTime(elapse, c.MoveQueueTimeDepth)
	}
	triggers := &ElapsedTriggers{}
	for _, c := range b.PlayerTeam {
		c.Progress(elapse, triggers)
	}
	for _, c := range b.OpponentTeam {
		c.Progress(elapse, triggers)
	}
	for triggers.First != nil {
		if triggers.First.Trigger != nil {
			fmt.Printf("running move trigger for '%s' during '%s'\n", triggers.First.Instance.Source.Details.Name, triggers.First.Instance.Move.Name)
			triggers.First.Trigger(b, triggers.First.Instance)
		}
		triggers.First = triggers.First.Next
	}
	b.ElapsedTime += elapse
	return elapse
}

type ElapsedTriggers struct {
	First *TriggerNode
}

type TriggerNode struct {
	Order    int
	Instance *MoveInstance
	Trigger  MoveTrigger
	Next     *TriggerNode
}

func (p *ElapsedTriggers) Append(new *TriggerNode) {
	if p.First == nil { // if nil, start the list
		p.First = new
		return
	}
	if new.Order < p.First.Order { // if we need to replace the first one
		new.Next = p.First
		p.First = new
		return
	}
	t := p.First
	for t.Next != nil { // for each node, if new goes betwen t and next, insert and return
		if t.Next.Order > new.Order {
			new.Next = t.Next
			t.Next = new
			return
		}
		t = t.Next
	}
	t.Next = new // else, append to the end
}
