package drawutil

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
)

type StrokePosition int

const (
	StrokeCenter StrokePosition = iota
	StrokeOutside
	StrokeInside
)

type StrokedRectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64

	RoundedCornerSize float64

	StrokeThickness float64
	StrokePosition  StrokePosition
}

func NewStrokedRectangle(x, y, width, height float64) *StrokedRectangle {
	return &StrokedRectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r *StrokedRectangle) Rounded(size float64) *StrokedRectangle {
	r.RoundedCornerSize = size
	return r
}

func (r *StrokedRectangle) Stroked(thickness float64, position StrokePosition) *StrokedRectangle {
	r.StrokeThickness = thickness
	r.StrokePosition = position
	return r
}

func (r *StrokedRectangle) Draw(canvas *gg.Context) {
	x, y, width, height := r.X, r.Y, r.Width, r.Height
	t := r.StrokeThickness
	if t > 0 {
		if r.StrokePosition == StrokeInside {
			t = math.Min(t, math.Min(width/2.0, height/2.0))
		}
		canvas.SetLineWidth(t)
		ht := t / 2.0
		switch r.StrokePosition {
		case StrokeInside:
			x += ht
			y += ht
			width -= t
			height -= t
			break
		case StrokeOutside:
			x -= ht
			y -= ht
			width += t
			height += t
			break
		}
	}
	if r.RoundedCornerSize > 0 {
		canvas.DrawRoundedRectangle(x, y, width, height, r.RoundedCornerSize)
	} else {
		canvas.DrawRectangle(x, y, width, height)
	}
}

func NewSizedRect(x, y, width, height int) image.Rectangle {
	return image.Rect(x, y, x+width, y+height)
}

func NewSizedRectF(x, y, width, height float64) image.Rectangle {
	return NewSizedRect(int(x), int(y), int(width), int(height))
}

func DrawRect(rect image.Rectangle, color color.Color, target *ebiten.Image) {
	vector.DrawFilledRect(target,
		float32(rect.Min.X), float32(rect.Min.Y),
		float32(rect.Dx()), float32(rect.Dy()),
		color, false)
}
