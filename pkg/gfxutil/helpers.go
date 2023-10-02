package gfxutil

import "github.com/faiface/pixel"

func Box(bottomLeft pixel.Vec, width, height float64) pixel.Rect {
	return pixel.Rect{
		Min: bottomLeft,
		Max: bottomLeft.Add(pixel.V(width, height)),
	}
}
