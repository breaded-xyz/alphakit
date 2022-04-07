package ta

func Slope(t1, t2 float64) int {
	s := t2 - t1
	switch {
	case s < 0:
		return -1
	case s > 0:
		return 1
	}
	return 0
}
