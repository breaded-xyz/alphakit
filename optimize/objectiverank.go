package optimize

type ObjectiveRanker func(a, b Report) bool

func SharpeRanker(a, b Report) bool {
	return a.Sharpe < b.Sharpe
}

func PRRRanker(a, b Report) bool {
	return a.PRR < b.PRR
}
