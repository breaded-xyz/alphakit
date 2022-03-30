package backtest

import (
	"errors"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var ErrInvalidOrderState = errors.New("order is not valid for processing")

type Simulator struct {
	clock Clock
	price market.Kline

	orders    map[broker.DealID]broker.Order
	positions map[broker.DealID]broker.Position
	trades    map[broker.DealID]broker.Trade
}

func NewSimulator() *Simulator {
	return &Simulator{
		clock:     NewClock(),
		orders:    make(map[broker.DealID]broker.Order),
		positions: make(map[broker.DealID]broker.Position),
		trades:    make(map[broker.DealID]broker.Trade),
	}
}

func (s *Simulator) AddOrder(order broker.Order) (broker.Order, error) {
	var empty broker.Order
	if order.Side == 0 || order.Type == 0 || order.State() != broker.OrderPending || !order.Size.IsPositive() {
		return empty, ErrInvalidOrderState
	}
	return s.processOrder(order), nil
}

func (s *Simulator) Next(price market.Kline) error {
	s.clock.NextEpoch(closeTime(s.price.Start, price.Start))
	s.price = price

	// Iterate open orders in the order they were placed with the dealer
	// Go maps do not maintain insertion order so we must sort the keys in a slice first
	// The map key is a ULID seeded from a time and supports lexicographic sorting
	ks := maps.Keys(s.orders)
	slices.Sort(ks)
	for _, k := range ks {
		order := s.orders[k]
		if order.State() == broker.OrderOpen {
			order = s.processOrder(order)
		}
	}

	return nil
}

func (s *Simulator) Orders() []broker.Order {
	return nil
}

func (s *Simulator) Positions() []broker.Position {
	return nil
}

func (s *Simulator) Trades() []broker.Trade {
	return nil
}

func (s *Simulator) EquityHistory() []broker.Equity {
	return nil
}

func (s *Simulator) processOrder(order broker.Order) broker.Order {
	switch order.State() {
	case broker.OrderPending:
		order = s.processOrder(s.openOrder(order))
	case broker.OrderOpen:
		if matchedPrice := matchOrder(order, s.price); matchedPrice.IsPositive() {
			order = s.processOrder(s.fillOrder(order, matchedPrice))
		}
	case broker.OrderFilled:
		position := s.updatePosition(s.getLatestOrNewPosition(), order)
		if position.State() == broker.OrderClosed {
			s.trades[position.ID] = s.newTrade(position)
		}
		s.positions[position.ID] = position
		order = s.closeOrder(order)
	}

	s.orders[order.ID] = order
	return order
}

func (s *Simulator) openOrder(order broker.Order) broker.Order {
	order.ID = broker.NewIDWithTime(s.clock.Now())
	order.OpenedAt = s.clock.Now()
	return order
}

func (s *Simulator) fillOrder(order broker.Order, matchedPrice decimal.Decimal) broker.Order {
	order.FilledAt = s.clock.Now()
	order.FilledPrice = matchedPrice
	order.FilledSize = order.Size
	return order
}

func (s *Simulator) closeOrder(order broker.Order) broker.Order {
	order.ClosedAt = s.clock.Now()
	return order
}

func (s *Simulator) getLatestOrNewPosition() broker.Position {
	var empty, position broker.Position

	if len(s.positions) == 0 {
		return empty
	}

	ks := maps.Keys(s.positions)
	slices.Sort(ks)
	position = s.positions[ks[len(ks)-1]]

	if position.State() == broker.PositionClosed {
		return empty
	}
	return position
}

func (s *Simulator) updatePosition(position broker.Position, order broker.Order) broker.Position {

	switch position.State() {
	case broker.PositionPending:
		position = s.openPosition(order)
	case broker.PositionOpen:
		if order.Side == position.Side.Opposite() {
			position = s.closePosition(position, order)
		}
	}
	return position
}

func (s *Simulator) openPosition(order broker.Order) broker.Position {
	return broker.Position{
		ID:       order.ID,
		OpenedAt: s.clock.Now(),
		Asset:    order.Asset,
		Side:     order.Side,
		Price:    order.FilledPrice,
		Size:     order.FilledSize,
	}
}

func (s *Simulator) closePosition(position broker.Position, order broker.Order) broker.Position {
	position.ClosedAt = s.clock.Now()
	position.LiquidationPrice = order.FilledPrice
	return position
}

func (s *Simulator) newTrade(position broker.Position) broker.Trade {
	return broker.Trade{
		ID:        position.ID,
		CreatedAt: s.clock.Now(),
		Asset:     position.Asset,
		Side:      position.Side,
		Size:      position.Size,
		Profit:    profit(position),
	}
}

func matchOrder(order broker.Order, quote market.Kline) decimal.Decimal {
	var matchedPrice decimal.Decimal

	switch order.Type {
	case broker.Limit:
		if dec.Between(order.LimitPrice, quote.L, quote.H) {
			matchedPrice = order.LimitPrice
		}
	case broker.Market:
		matchedPrice = quote.C
	}

	return matchedPrice
}

func closeTime(start1, start2 time.Time) time.Time {
	if start1.IsZero() {
		return start2
	}
	interval := start2.Sub(start1)
	return start2.Add(interval)
}

func profit(position broker.Position) decimal.Decimal {
	profit := position.LiquidationPrice.Sub(position.Price).Mul(position.Size)
	if position.Side == broker.Sell {
		profit = profit.Neg()
	}
	return profit
}
