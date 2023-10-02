package gfxutil

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type VAlignment int

const Top VAlignment = 1
const Middle VAlignment = 0
const Bottom VAlignment = -1

type HAlignment int

const Left HAlignment = -1
const Center HAlignment = 0
const Right HAlignment = 1

type Text struct {
	String string
	Color  pixel.RGBA
	Bounds pixel.Rect
	VAlign VAlignment
	HAlign HAlignment

	ShadowColor pixel.RGBA
	ShadowDepth pixel.Vec
}

func (t Text) Draw(font *text.Text, delta pixel.Vec) {
	sBounds := font.BoundsOf(t.String)
	dx, dy := 0.0, 0.0
	hDiff := t.Bounds.H() - sBounds.H()
	switch t.VAlign {
	case Top:
		dy = hDiff
		break
	case Middle:
		dy = hDiff / 2.0
		break
	}
	wDiff := t.Bounds.W() - sBounds.W()
	switch t.HAlign {
	case Center:
		dx = wDiff / 2.0
		break
	case Right:
		dx = wDiff
		break
	}
	dot := t.Bounds.Min.Add(delta).Add(pixel.V(dx, dy))
	if t.ShadowDepth.X != 0 && t.ShadowDepth.Y != 0 {
		font.Color = t.ShadowColor
		font.Dot = dot.Add(t.ShadowDepth)
		font.WriteString(t.String)
	}
	font.Color = t.Color
	font.Dot = dot
	font.WriteString(t.String)
}
