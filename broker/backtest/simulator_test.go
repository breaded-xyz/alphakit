package backtest

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var _fixed time.Time = time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)

// newSimulatorForTest sets the simulation clock to a fixed time
func newSimulatorForTest() *Simulator {
	sim := NewSimulator()
	sim.clock = &StubClock{Fixed: _fixed}
	return sim
}

func TestSimulatorAddOrder(t *testing.T) {
	tests := []struct {
		name string
		give broker.Order
		want error
	}{
		{
			name: "invalid side",
			give: broker.Order{
				Side: 0,
				Type: broker.Market,
				Size: dec.New(1),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "invalid type",
			give: broker.Order{
				Side: broker.Buy,
				Type: 0,
				Size: dec.New(1),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "invalid size",
			give: broker.Order{
				Side: broker.Buy,
				Type: broker.Market,
				Size: dec.New(0),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "no pending state",
			give: broker.Order{
				OpenedAt: time.Now(),
			},
			want: ErrInvalidOrderState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sim Simulator
			_, err := sim.AddOrder(tt.give)
			assert.ErrorIs(t, err, tt.want)
		})
	}
}

func TestSimulatorProcessOrder(t *testing.T) {
	tests := []struct {
		name      string
		give      broker.Order
		wantOrder broker.Order
		wantState broker.OrderState
	}{
		{
			name: "market order filled",
			give: broker.Order{
				Type: broker.Market,
				Side: broker.Buy,
				Size: dec.New(1),
			},
			wantOrder: broker.Order{
				FilledPrice: dec.New(10),
				FilledSize:  dec.New(1),
			},
			wantState: broker.OrderClosed,
		},
		{
			name: "limit order filled: open time before now",
			give: broker.Order{
				OpenedAt:   _fixed.Add(time.Hour),
				Side:       broker.Buy,
				Type:       broker.Limit,
				LimitPrice: dec.New(8),
				Size:       dec.New(1),
			},
			wantOrder: broker.Order{
				FilledPrice: dec.New(8),
				FilledSize:  dec.New(1),
			},
			wantState: broker.OrderClosed,
		},
		{
			name: "limit order not filled: open time equals now",
			give: broker.Order{
				Side:       broker.Buy,
				Type:       broker.Limit,
				LimitPrice: dec.New(8),
				Size:       dec.New(1),
			},
			wantOrder: broker.Order{},
			wantState: broker.OrderOpen,
		},
		{
			name: "limit order opened but no price match",
			give: broker.Order{
				OpenedAt:   _fixed.Add(time.Hour),
				Type:       broker.Limit,
				LimitPrice: dec.New(100),
				Size:       dec.New(1),
			},
			wantOrder: broker.Order{},
			wantState: broker.OrderOpen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sim := newSimulatorForTest()
			sim.marketPrice = market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)}
			act := sim.processOrder(tt.give)

			assert.Equal(t, tt.wantOrder.FilledSize, act.FilledSize)
			assert.Equal(t, tt.wantOrder.FilledPrice, act.FilledPrice)
			assert.Equal(t, tt.wantState, act.State())
		})
	}
}

func TestSimulatorOpenOrder(t *testing.T) {
	sim := newSimulatorForTest()
	exp := broker.Order{
		ID:       broker.NewIDWithTime(_fixed),
		OpenedAt: _fixed,
	}
	act := sim.openOrder(broker.Order{})
	assert.Equal(t, exp, act)
}

func TestSimulatorFillOrder(t *testing.T) {
	sim := newSimulatorForTest()
	exp := broker.Order{
		Side:        broker.Buy,
		FilledAt:    _fixed,
		Size:        dec.New(1),
		FilledPrice: dec.New(100),
		FilledSize:  dec.New(1),
	}
	act := sim.fillOrder(broker.Order{Side: broker.Buy, Size: exp.Size}, exp.FilledPrice)
	assert.Equal(t, exp, act)
}

func TestSimulatorCloseOrder(t *testing.T) {
	sim := newSimulatorForTest()
	exp := broker.Order{
		ClosedAt: _fixed,
	}
	act := sim.closeOrder(broker.Order{})
	assert.Equal(t, exp, act)
}

func TestSimulatorGetPosition(t *testing.T) {
	tests := []struct {
		name string
		give map[broker.DealID]broker.Position
		want broker.PositionState
	}{
		{
			name: "no positions",
			give: map[broker.DealID]broker.Position{},
			want: broker.PositionPending,
		},
		{
			name: "latest position is closed",
			give: map[broker.DealID]broker.Position{
				"1": {OpenedAt: _fixed},
				"2": {ClosedAt: _fixed},
			},
			want: broker.PositionPending,
		},
		{
			name: "latest position is open",
			give: map[broker.DealID]broker.Position{"1": {ID: "1", OpenedAt: _fixed}},
			want: broker.PositionOpen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := newSimulatorForTest()
			sim.positions = tt.give
			act := sim.getPosition()
			assert.Equal(t, tt.want, act.State())
		})
	}
}

func TestSimulatorProcessPosition(t *testing.T) {
	tests := []struct {
		name         string
		giveOrder    broker.Order
		givePosition broker.Position
		wantPosition broker.Position
		wantState    broker.PositionState
	}{
		{
			name:         "open new position",
			giveOrder:    broker.Order{Side: broker.Buy, FilledAt: _fixed, FilledPrice: dec.New(10), FilledSize: dec.New(1)},
			givePosition: broker.Position{},
			wantPosition: broker.Position{
				OpenedAt: _fixed,
				Side:     broker.Buy,
				Price:    dec.New(10),
				Size:     dec.New(1),
			},
			wantState: broker.PositionOpen,
		},
		{
			name:         "close existing position",
			giveOrder:    broker.Order{FilledAt: _fixed, Side: broker.Sell, FilledPrice: dec.New(20), FilledSize: dec.New(1)},
			givePosition: broker.Position{OpenedAt: _fixed, Side: broker.Buy, Price: dec.New(10), Size: dec.New(1)},
			wantPosition: broker.Position{
				OpenedAt:         _fixed,
				ClosedAt:         _fixed,
				Side:             broker.Buy,
				Price:            dec.New(10),
				Size:             dec.New(1),
				LiquidationPrice: dec.New(20),
			},
			wantState: broker.PositionClosed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := newSimulatorForTest()
			act := sim.processPosition(tt.givePosition, tt.giveOrder)
			assert.Equal(t, tt.wantPosition, act)
			assert.Equal(t, tt.wantState, act.State())
		})
	}
}

func TestSimulatorOpenPosition(t *testing.T) {

	sim := newSimulatorForTest()

	exp := broker.Position{
		ID:       "1",
		OpenedAt: _fixed,
		Asset:    market.NewAsset("BTCUSD"),
		Side:     broker.Buy,
		Price:    dec.New(10),
		Size:     dec.New(1),
	}

	act := sim.openPosition(broker.Order{
		ID:          exp.ID,
		FilledAt:    _fixed,
		Asset:       exp.Asset,
		Side:        exp.Side,
		FilledPrice: exp.Price,
		FilledSize:  exp.Size,
	})

	assert.Equal(t, exp, act)
}

func TestSimulatorReceivePrice(t *testing.T) {

	sim := NewSimulator()
	sim.clock.Start(time.Now(), time.Millisecond)

	k1 := broker.NewIDWithTime(sim.clock.Now())
	sim.orders[k1] = broker.Order{ID: k1, Type: broker.Limit, LimitPrice: dec.New(15), OpenedAt: sim.clock.Now()}

	k2 := broker.NewIDWithTime(sim.clock.Now())
	sim.orders[k2] = broker.Order{ID: k2, Type: broker.Limit, LimitPrice: dec.New(15), OpenedAt: sim.clock.Now()}

	k3 := broker.NewIDWithTime(sim.clock.Now())
	sim.orders[k3] = broker.Order{ID: k3, Type: broker.Limit, LimitPrice: dec.New(10), OpenedAt: sim.clock.Now()}

	price := market.Kline{
		Start: sim.clock.Now().Add(time.Hour * 1),
		O:     dec.New(8),
		H:     dec.New(15),
		L:     dec.New(5),
		C:     dec.New(10)}

	sim.Next(price)

	t.Run("all open orders are closed", func(t *testing.T) {
		for _, v := range sim.orders {
			if v.State() != broker.OrderClosed {
				assert.Fail(t, "expect all orders to be closed")
			}
		}
	})

	t.Run("orders are processed in sequence they were created", func(t *testing.T) {
		assert.True(t, sim.orders[k1].ClosedAt.Before(sim.orders[k2].ClosedAt))
		assert.True(t, sim.orders[k2].ClosedAt.Before(sim.orders[k3].ClosedAt))
	})
}

func TestMatchOrder(t *testing.T) {
	tests := []struct {
		name      string
		giveOrder broker.Order
		giveQuote market.Kline
		want      decimal.Decimal
	}{
		{
			name: "match limit",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(12),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(12),
		},
		{
			name: "match limit lower bound inclusive",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(5),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(5),
		},
		{
			name: "match limit upper bound inclusive",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(15),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(15),
		},
		{
			name: "no match limit below lower bound",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(2),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      decimal.Decimal{},
		},
		{
			name: "no match limit above upper bound",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(100),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      decimal.Decimal{},
		},
		{
			name: "always match market on latest close price",
			giveOrder: broker.Order{
				Type: broker.Market,
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := matchOrder(tt.giveOrder, tt.giveQuote)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestCloseTime(t *testing.T) {
	interval := time.Hour * 4
	start1 := time.Now()
	start2 := start1.Add(interval)

	t.Run("valid start times", func(t *testing.T) {
		exp := start2.Add(interval)
		act := closeTime(start1, start2)
		assert.EqualValues(t, exp, act)
	})

	t.Run("start 1 is zero", func(t *testing.T) {
		act := closeTime(time.Time{}, start2)
		assert.EqualValues(t, start2, act)
	})
}

func TestProfit(t *testing.T) {
	tests := []struct {
		name string
		give broker.Position
		want decimal.Decimal
	}{
		{
			name: "buy side profit",
			give: broker.Position{
				Side:             broker.Buy,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(20),
			},
			want: dec.New(20),
		},
		{
			name: "sell side profit",
			give: broker.Position{
				Side:             broker.Sell,
				Price:            dec.New(100),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(50),
			},
			want: dec.New(100),
		},
		{
			name: "buy side loss",
			give: broker.Position{
				Side:             broker.Buy,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(5),
			},
			want: dec.New(-10),
		},
		{
			name: "sell side loss",
			give: broker.Position{
				Side:             broker.Sell,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(20),
			},
			want: dec.New(-20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := profit(tt.give, tt.give.LiquidationPrice)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestCreateTrade(t *testing.T) {

	sim := newSimulatorForTest()

	give := broker.Position{
		ID:               "1",
		ClosedAt:         _fixed,
		Asset:            market.NewAsset("BTCUSD"),
		Side:             broker.Sell,
		Price:            dec.New(10),
		Size:             dec.New(2),
		LiquidationPrice: dec.New(20),
	}

	want := broker.Trade{
		ID:        give.ID,
		CreatedAt: give.ClosedAt,
		Asset:     give.Asset,
		Side:      give.Side,
		Size:      give.Size,
		Profit:    dec.New(-20),
	}

	act := sim.createTrade(give)
	assert.Equal(t, want, act)
}

func TestEquityNow(t *testing.T) {
	sim := newSimulatorForTest()
	sim.accountBalance = dec.New(10)
	sim.marketPrice = market.Kline{C: dec.New(20)}

	t.Run("open position - unrealized profit", func(t *testing.T) {
		sim.positions = map[broker.DealID]broker.Position{
			"1": {OpenedAt: _fixed, Side: broker.Sell, Price: dec.New(10), Size: dec.New(2)},
		}
		exp := dec.New(-10)
		act := sim.equityNow()
		assert.Equal(t, exp, act)
	})

	t.Run("closed position - just account balance", func(t *testing.T) {
		sim.positions = map[broker.DealID]broker.Position{
			"1": {ClosedAt: _fixed, Side: broker.Sell, Price: dec.New(10), Size: dec.New(2)},
		}
		exp := dec.New(10)
		act := sim.equityNow()
		assert.Equal(t, exp, act)
	})

}

func TestEqualClock(t *testing.T) {
	t1 := time.Date(0, 0, 0, 1, 1, 1, 5, time.Local)
	t2 := time.Date(0, 0, 0, 1, 1, 1, 8, time.Local)
	assert.True(t, equalClock(t1, t2))

	t1 = time.Date(0, 0, 0, 1, 2, 1, 5, time.Local)
	t2 = time.Date(0, 0, 0, 1, 1, 1, 8, time.Local)
	assert.False(t, equalClock(t1, t2))
}
