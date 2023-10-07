package drawutil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

func DrawRect(rect image.Rectangle, color color.Color, target *ebiten.Image) {
	vector.DrawFilledRect(target,
		float32(rect.Min.X), float32(rect.Min.Y),
		float32(rect.Dx()), float32(rect.Dy()),
		color, false)
}
