package optimize

// ObjectiveRanker is used by an Optimizer to sort the results of backtest trials
// and select the best performing ParamSet. ObjectiveRanker is equivalent to a 'less' comparison function.
type ObjectiveRanker func(a, b PhaseReport) bool

// SharpeRanker ranks by Sharpe ratio.
func SharpeRanker(a, b PhaseReport) bool {
	return a.Sharpe < b.Sharpe
}

// PRRRanker ranks by Pessimistic Return Ratio.
func PRRRanker(a, b PhaseReport) bool {
	return a.PRR < b.PRR
}
