package backtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/dec"
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

	cost := &PerpCoster{
		SpreadPct:      dec.New(0.02),
		SlippagePct:    dec.New(0.01),
		TransactionPct: dec.New(0.1),
		FundingHourPct: dec.New(0.001),
	}
	dealer := NewDealerWithCost(cost)

	for i, price := range prices {
		assert.NoError(t, dealer.ReceivePrice(context.Background(), price))
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
	act := dealer.simulator.Balance().Trade
	assert.True(t, act.Equal(exp))
}

func TestShortTradeWithStop(t *testing.T) {
	start := time.Now()
	prices := []market.Kline{
		{Start: start.Add(0 * time.Hour), O: dec.New(10), H: dec.New(20), L: dec.New(5), C: dec.New(8)},
		{Start: start.Add(1 * time.Hour), O: dec.New(8), H: dec.New(15), L: dec.New(4), C: dec.New(10)},
		{Start: start.Add(2 * time.Hour), O: dec.New(10), H: dec.New(26), L: dec.New(10), C: dec.New(20)},
		{Start: start.Add(3 * time.Hour), O: dec.New(20), H: dec.New(25), L: dec.New(6), C: dec.New(10)},
		{Start: start.Add(4 * time.Hour), O: dec.New(10), H: dec.New(30), L: dec.New(3), C: dec.New(5)},
	}

	dealer := NewDealer()
	for i, price := range prices {
		assert.NoError(t, dealer.ReceivePrice(context.Background(), price))
		switch i {
		case 1:
			_, _, err := dealer.PlaceOrder(context.Background(), broker.Order{
				Side:       broker.Sell,
				Type:       broker.Limit,
				LimitPrice: dec.New(14),
				Size:       dec.New(2),
			})
			assert.NoError(t, err)
			_, _, err = dealer.PlaceOrder(context.Background(), broker.Order{
				Side:       broker.Buy,
				Type:       broker.Limit,
				LimitPrice: dec.New(28),
				Size:       dec.New(2),
				ReduceOnly: true,
			})
			assert.NoError(t, err)
		}
	}

	exp := dec.New(-28)
	act := dealer.simulator.Balance().Trade
	assert.True(t, act.Equal(exp))
}
