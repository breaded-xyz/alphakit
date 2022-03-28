package backtest

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

func TestDealer_processOrder(t *testing.T) {
	order := broker.Order{
		OpenedAt:    time.Time{},
		FilledAt:    time.Time{},
		ClosedAt:    time.Time{},
		Asset:       market.Asset{},
		Side:        broker.Buy,
		Type:        broker.Market,
		Size:        decimal.Decimal{},
		ReduceOnly:  false,
		FilledPrice: decimal.Decimal{},
		FilledSize:  decimal.Decimal{},
	}

	dealer := NewDealer()

	order = dealer.processOrder(order)
}
