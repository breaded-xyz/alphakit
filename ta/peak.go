package ta

func Peak(t1, t2, t3 float64) bool {
	slo1 := Slope(t1, t2)
	slo2 := Slope(t2, t3)
	return (slo1 == 1 || slo1 == 0) && slo2 == -1
}

func Valley(t1, t2, t3 float64) bool {
	slo1 := Slope(t1, t2)
	slo2 := Slope(t2, t3)
	return (slo1 == -1 || slo1 == 0) && slo2 == 1
}
