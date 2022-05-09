package ta

// Lookback returns a value from series at n index ago.
// Series must be in chronological order, with the earliest value at slice index 0.
// n = 0 returns the latest value. n = 1 returns the value before the latest etc.
func Lookback(series []float64, n int) float64 {
	i := (len(series) - n) - 1
	if i < 0 {
		return 0
	}
	return series[i]
}

// Window returns a copied slice of series starting at n index ago.
// Semantics of n argument are the same as Lookback function.
func Window(series []float64, n int) []float64 {
	window := make([]float64, len(series))
	copy(window, series)

	ln := len(series)
	if ln <= n {
		return window
	}

	i := (ln - n) - 1
	return window[i:]
}

// WindowAppend appends a value to the end of the series and slices it to the window starting at n index ago.
// Semantics of n argument are the same as Window and Lookback functions.
func WindowAppend(series []float64, n int, v float64) []float64 {
	return Window(append(series, v), n)
}
