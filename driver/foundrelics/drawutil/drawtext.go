package drawutil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

type AlignHorizontal int

const (
	AlignLeft AlignHorizontal = iota
	AlignCenter
	AlignRight
)

type AlignVertical int

const (
	AlignTop AlignVertical = iota
	AlignMiddle
	AlignBottom
)

type TextDrawer struct {
	Face      font.Face
	X         int
	Y         int
	TextColor color.Color

	shadow      bool
	shadowDx    int
	shadowDy    int
	shadowColor color.Color

	bounded         bool
	boundWidth      int
	boundHeight     int
	alignVertical   AlignVertical
	alignHorizontal AlignHorizontal
}

func NewTextDrawer(face font.Face, x, y int, textColor color.Color) *TextDrawer {
	return &TextDrawer{
		Face:      face,
		X:         x,
		Y:         y,
		TextColor: textColor,
	}
}

func (t *TextDrawer) Shadowed(dx, dy int, shadowColor color.Color) *TextDrawer {
	t.shadow = true
	t.shadowDx = dx
	t.shadowDy = dy
	t.shadowColor = shadowColor
	return t
}

func (t *TextDrawer) Move(dx, dy int) *TextDrawer {
	t.X += dx
	t.Y += dy
	return t
}

func (t *TextDrawer) Color(textColor color.Color) *TextDrawer {
	t.TextColor = textColor
	return t
}

func (t *TextDrawer) Bounded(width, height int, alignH AlignHorizontal, alignV AlignVertical) *TextDrawer {
	t.bounded = true
	t.boundWidth = width
	t.boundHeight = height
	t.alignHorizontal = alignH
	t.alignVertical = alignV
	return t
}

func (t *TextDrawer) GetBoundingBox() image.Rectangle {
	b := image.Rectangle{
		Min: image.Point{X: t.X, Y: t.Y},
		Max: image.Point{X: t.X, Y: t.Y},
	}
	if t.bounded {
		b.Max.X += t.boundWidth
		b.Max.Y += t.boundHeight
	}
	return b
}

func (t *TextDrawer) Draw(text string, target *ebiten.Image) image.Rectangle {
	bounds := t.BoundsOf(text)
	dx := 0
	switch t.alignHorizontal {
	case AlignLeft:
		break
	case AlignCenter:
		dx = (t.boundWidth - bounds.Dx()) / 2
		break
	case AlignRight:
		dx = t.boundWidth - bounds.Dx()
		break
	}
	dy := 0
	switch t.alignVertical {
	case AlignTop:
		break
	case AlignMiddle:
		dy = (t.boundHeight - bounds.Dy()) / 2
		break
	case AlignBottom:
		dy = t.boundHeight - bounds.Dy()
		break
	}
	t.drawText(text, dx, dy, target)
	return bounds.Add(image.Point{X: dx, Y: dy})
}

func (t *TextDrawer) drawText(textContent string, dx, dy int, dest *ebiten.Image) {
	x := t.X + dx
	y := t.Y + dy + t.Face.Metrics().CapHeight.Ceil()
	for _, line := range strings.Split(textContent, "\n") {
		if t.shadow {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(x+t.shadowDx), float64(y+t.shadowDy))
			opt.ColorScale.ScaleWithColor(t.shadowColor)
			text.DrawWithOptions(dest, line, t.Face, opt)
		}
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(x), float64(y))
		opt.ColorScale.ScaleWithColor(t.TextColor)
		text.DrawWithOptions(dest, line, t.Face, opt)
		y += t.Face.Metrics().Height.Round()
	}
}

func (t *TextDrawer) BoundsOf(txt string) image.Rectangle {
	if txt == "" {
		return image.Rectangle{Min: image.Point{X: t.X, Y: t.Y}, Max: image.Point{X: t.X, Y: t.Y}}
	}
	x := t.X
	y := t.Y
	var bounds image.Rectangle
	f := t.Face
	m := f.Metrics()
	h := m.Height
	lineHeight := h.Round()
	for n, line := range strings.Split(txt, "\n") {
		fb, _ := font.BoundString(t.Face, line)
		lb := image.Rect(x+fb.Min.X.Floor(), y+fb.Min.Y.Floor(), x+fb.Max.X.Ceil(), y+fb.Max.Y.Ceil())
		if n == 0 { // first line, text draws from bottom left, we want to draw from top left
			capHeight := t.Face.Metrics().CapHeight.Ceil()
			y += capHeight
			bounds = lb.Add(image.Point{X: 0, Y: capHeight})
		} else {
			bounds = bounds.Union(lb)
		}
		y += lineHeight
	}
	if t.shadow {
		bounds.Max.X += t.shadowDx
		bounds.Max.Y += t.shadowDy
	}
	return bounds
}

func DrawString(txt string, face font.Face, x, y int, color color.Color, dest *ebiten.Image) image.Rectangle {
	return NewTextDrawer(face, x, y, color).Draw(txt, dest)
}
