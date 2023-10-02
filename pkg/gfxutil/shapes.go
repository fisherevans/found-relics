package gfxutil

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type StrokePosition int

const StrokeInner StrokePosition = 0
const StrokeCenter StrokePosition = 1
const StrokeOuter StrokePosition = 2

type RectShape struct {
	Bounds pixel.Rect

	Fill      bool
	FillColor pixel.RGBA

	StrokeWidth    float64
	StrokeColor    pixel.RGBA
	StrokeShape    imdraw.EndShape
	StrokePosition StrokePosition
}

func (r RectShape) Draw(imd *imdraw.IMDraw) {
	bl := r.Bounds.Min
	w := r.Bounds.W()
	h := r.Bounds.H()
	hsw := r.StrokeWidth / 2
	if r.StrokePosition == StrokeInner {
		bl = bl.Add(pixel.V(hsw, hsw))
		w -= r.StrokeWidth
		h -= r.StrokeWidth
	} else if r.StrokePosition == StrokeOuter {
		bl = bl.Sub(pixel.V(hsw, hsw))
		w += r.StrokeWidth
		h += r.StrokeWidth
	}
	ps := []pixel.Vec{
		pixel.V(bl.X, bl.Y),
		pixel.V(bl.X+w, bl.Y),
		pixel.V(bl.X+w, bl.Y+h),
		pixel.V(bl.X, bl.Y+h)}
	if r.Fill {
		imd.Color = r.FillColor
		imd.Push(ps...)
		imd.Polygon(0)
	}
	if r.StrokeWidth > 0 {
		ps = append(ps, ps[0])
		imd.EndShape = r.StrokeShape
		imd.Color = r.StrokeColor
		imd.Push(ps...)
		imd.Polygon(r.StrokeWidth)
	}
}
