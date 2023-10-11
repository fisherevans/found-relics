package drawutil

import "math"

func InterpolateSin(in float64) float64 {
	return (math.Sin(((in*2.0)-1.0)*math.Pi/2.0) + 1.0) / 2.0
}
