package drawutil

import (
	"github.com/fogleman/gg"
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

type VAlignMode int

const (
	VAlignCapHeight VAlignMode = iota
	VAlignLineHeight
)

type HAlignMode int

const (
	HAlignBlock HAlignMode = iota
	HAlignLine
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
	boundShadow bool

	bounded             bool
	boundWidth          int
	boundHeight         int
	AlignVertical       AlignVertical
	AlignHorizontal     AlignHorizontal
	AlignModeHorizontal HAlignMode
	AlignModeVertical   VAlignMode
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

func (t *TextDrawer) BoundShadow(do bool) *TextDrawer {
	t.boundShadow = do
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

func (t *TextDrawer) FontFace(face font.Face) *TextDrawer {
	t.Face = face
	return t
}

func (t *TextDrawer) Bounded(width, height int, alignH AlignHorizontal, alignV AlignVertical) *TextDrawer {
	t.bounded = true
	t.boundWidth = width
	t.boundHeight = height
	t.AlignHorizontal = alignH
	t.AlignVertical = alignV
	return t
}

func (t *TextDrawer) Unbounded() *TextDrawer {
	t.bounded = false
	return t
}

func (t *TextDrawer) WithAlignModeHorizontal(mode HAlignMode) *TextDrawer {
	t.AlignModeHorizontal = mode
	return t
}

func (t *TextDrawer) WithAlignModeVertical(mode VAlignMode) *TextDrawer {
	t.AlignModeVertical = mode
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
	dx, dy, bounds := t.calcPosition(text)
	t.drawText(text, dx, dy, bounds.Dx(), target)
	return bounds
}

func (t *TextDrawer) DrawToCanvas(text string, canvas *gg.Context) image.Rectangle {
	dx, dy, bounds := t.calcPosition(text)
	t.drawTextCanvas(text, dx, dy, bounds.Dx(), canvas)
	return bounds
}

func (t *TextDrawer) calcPosition(text string) (int, int, image.Rectangle) {
	bounds := t.BoundsOf(text)
	dx := 0
	switch t.AlignHorizontal {
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
	switch t.AlignVertical {
	case AlignTop:
		break
	case AlignMiddle:
		dy = (t.boundHeight - bounds.Dy()) / 2
		break
	case AlignBottom:
		dy = t.boundHeight - bounds.Dy()
		break
	}
	return dx, dy, bounds.Add(image.Point{X: dx, Y: dy})
}

func (t *TextDrawer) drawText(textContent string, dx, dy, width int, dest *ebiten.Image) {
	blockX := t.X + dx
	y := t.Y + dy
	if t.AlignModeVertical == VAlignCapHeight {
		y += t.Face.Metrics().CapHeight.Ceil()
	} else {
		y += t.Face.Metrics().Height.Ceil() - t.Face.Metrics().Descent.Ceil()
	}
	for _, line := range strings.Split(textContent, "\n") {
		lineX := blockX
		if t.AlignModeHorizontal == HAlignLine {
			fb, _ := font.BoundString(t.Face, line)
			switch t.AlignHorizontal {
			case AlignRight:
				lineX += width - (fb.Max.X.Round() - fb.Min.X.Round())
				break
			case AlignCenter:
				lineX += (width - (fb.Max.X.Round() - fb.Min.X.Round())) / 2
				break
			}
		}
		if t.shadow {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(lineX+t.shadowDx), float64(y+t.shadowDy))
			opt.ColorScale.ScaleWithColor(t.shadowColor)
			text.DrawWithOptions(dest, line, t.Face, opt)
		}
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(lineX), float64(y))
		opt.ColorScale.ScaleWithColor(t.TextColor)
		text.DrawWithOptions(dest, line, t.Face, opt)
		y += t.Face.Metrics().Height.Round()
	}
}

func (t *TextDrawer) drawTextCanvas(textContent string, dx, dy, width int, dest *gg.Context) {
	blockX := t.X + dx
	y := t.Y + dy
	if t.AlignModeVertical == VAlignCapHeight {
		y += t.Face.Metrics().CapHeight.Ceil()
	} else {
		y += t.Face.Metrics().Height.Ceil()
	}
	for _, line := range strings.Split(textContent, "\n") {
		lineX := blockX
		if t.AlignModeHorizontal == HAlignLine {
			fb, _ := font.BoundString(t.Face, line)
			switch t.AlignHorizontal {
			case AlignRight:
				lineX += width - (fb.Max.X.Round() - fb.Min.X.Round())
				break
			case AlignCenter:
				lineX += (width - (fb.Max.X.Round() - fb.Min.X.Round())) / 2
				break
			}
		}
		dest.SetFontFace(t.Face)
		if t.shadow {
			dest.SetColor(t.shadowColor)
			dest.DrawString(line, float64(lineX+t.shadowDx), float64(y+t.shadowDy))
		}
		dest.SetColor(t.TextColor)
		dest.DrawString(line, float64(lineX), float64(y))
		y += t.Face.Metrics().Height.Round()
	}
}

func (t *TextDrawer) BoundsOf(txt string) image.Rectangle {
	if txt == "" {
		return image.Rectangle{Min: image.Point{X: t.X, Y: t.Y}, Max: image.Point{X: t.X, Y: t.Y}}
	}

	width := 0
	lines := strings.Split(txt, "\n")
	for _, line := range lines {
		fb, _ := font.BoundString(t.Face, line)
		dx := fb.Max.X.Round() - fb.Min.X.Round()
		if dx > width {
			width = dx
		}
	}

	var height int
	lineH := t.Face.Metrics().Height.Round()
	if t.AlignModeVertical == VAlignCapHeight {
		capH := t.Face.Metrics().CapHeight.Ceil()
		lineGap := lineH - capH
		height = (len(lines) * capH) + ((len(lines) - 1) * lineGap)
	} else {
		height = len(lines) * lineH
	}

	bounds := NewSizedRect(t.X, t.Y, width, height)
	if t.boundShadow {
		bounds.Max.X += t.shadowDx
		bounds.Max.Y += t.shadowDy
	}
	return bounds
}

func DrawString(txt string, face font.Face, x, y int, color color.Color, dest *ebiten.Image) image.Rectangle {
	return NewTextDrawer(face, x, y, color).Draw(txt, dest)
}
