package zero2algo

// Enforce at compile time that the type implements the interface
var _ SimulatedDealer = (*BacktestDealer)(nil)

type BacktestDealer struct {
}

func NewBacktestDealer() *BacktestDealer {
	return nil
}

func (d *BacktestDealer) ListTrades() []Trade {
	return nil
}

func (d *BacktestDealer) EquityCurve() []Equity {
	return nil
}

func (d *BacktestDealer) ReceivePrice(price Kline) error {
	return nil
}
