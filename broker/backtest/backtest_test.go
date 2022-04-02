package backtest

import (
	"context"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
)

func TestLongTradeWithCosts(t *testing.T) {

	start := time.Now()
	prices := []market.Kline{
		{Start: start.Add(0 * time.Hour), O: dec.New(10), H: dec.New(20), L: dec.New(5), C: dec.New(8)},
		{Start: start.Add(1 * time.Hour), O: dec.New(8), H: dec.New(15), L: dec.New(4), C: dec.New(10)},
		{Start: start.Add(2 * time.Hour), O: dec.New(10), H: dec.New(30), L: dec.New(10), C: dec.New(20)},
		{Start: start.Add(3 * time.Hour), O: dec.New(20), H: dec.New(25), L: dec.New(6), C: dec.New(10)},
		{Start: start.Add(4 * time.Hour), O: dec.New(10), H: dec.New(20), L: dec.New(3), C: dec.New(5)},
	}

	cost := &PerpCost{
		SpreadPct:      dec.New(0.01),
		SlippagePct:    dec.New(0.01),
		TransactionPct: dec.New(0.1),
		FundingHourPct: dec.New(0.001),
	}
	dealer := NewDealerWithCost(cost)

	for i, price := range prices {
		dealer.ReceivePrice(context.Background(), price)
		switch i {
		case 1:
			_, _, err := dealer.PlaceOrder(context.Background(), broker.Order{
				Side: broker.Buy,
				Type: broker.Market,
				Size: dec.New(2),
			})
			assert.NoError(t, err)
		case 4:
			_, _, err := dealer.PlaceOrder(context.Background(), broker.Order{
				Side: broker.Sell,
				Type: broker.Market,
				Size: dec.New(2),
			})
			assert.NoError(t, err)
		}
	}

	exp := dec.New(-13.6913)
	act := dealer.simulator.AccountBalance()
	assert.True(t, act.Equal(exp))
}
