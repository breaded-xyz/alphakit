package backtest

import (
	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/price"
)

// Enforce at compile time that the type implements the interface
var _ broker.SimulatedDealer = (*Dealer)(nil)

type Dealer struct {
}

func NewDealer() *Dealer {
	return nil
}

func (d *Dealer) ListTrades() []broker.Trade {
	return nil
}

func (d *Dealer) EquityCurve() []broker.Equity {
	return nil
}

func (d *Dealer) ReceivePrice(price price.Kline) error {
	return nil
}
