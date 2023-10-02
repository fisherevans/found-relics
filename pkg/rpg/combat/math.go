package combat

func MaxTime(a, b Time) Time {
	if a > b {
		return a
	}
	return b
}

func MinTime(a, b Time) Time {
	if a < b {
		return a
	}
	return b
}
