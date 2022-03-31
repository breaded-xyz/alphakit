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
	clock          Clocker
	accountBalance decimal.Decimal
	marketPrice    market.Kline

	orders    map[broker.DealID]broker.Order
	positions map[broker.DealID]broker.Position
	trades    map[broker.DealID]broker.Trade
	equity    broker.EquitySeries
}

func NewSimulator() *Simulator {
	return &Simulator{
		clock:     NewClock(),
		orders:    make(map[broker.DealID]broker.Order),
		positions: make(map[broker.DealID]broker.Position),
		trades:    make(map[broker.DealID]broker.Trade),
		equity:    make(broker.EquitySeries),
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
	s.clock.NextEpoch(closeTime(s.marketPrice.Start, price.Start))
	s.marketPrice = price

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

	s.equity[broker.Timestamp(s.clock.Epoch().Unix())] = s.equityNow()

	return nil
}

func (s *Simulator) Orders() []broker.Order {
	return maps.Values(s.orders)
}

func (s *Simulator) Positions() []broker.Position {
	return maps.Values(s.positions)
}

func (s *Simulator) Trades() []broker.Trade {
	return maps.Values(s.trades)
}

func (s *Simulator) Equity() broker.EquitySeries {
	var copied broker.EquitySeries
	maps.Copy(copied, s.equity)
	return copied
}

func (s *Simulator) processOrder(order broker.Order) broker.Order {
	switch order.State() {
	case broker.OrderPending:
		order = s.processOrder(s.openOrder(order))
	case broker.OrderOpen:
		if matchedPrice := matchOrder(order, s.marketPrice); matchedPrice.IsPositive() {
			order = s.processOrder(s.fillOrder(order, matchedPrice))
		}
	case broker.OrderFilled:
		s.processPosition(s.getLatestOrNewPosition(), order)
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

func (s *Simulator) processPosition(position broker.Position, order broker.Order) broker.Position {
	switch position.State() {
	case broker.PositionPending:
		position = s.processPosition(s.openPosition(order), order)
	case broker.PositionOpen:
		if order.Side == position.Side.Opposite() {
			position = s.processPosition(s.closePosition(position, order), order)
		}
	case broker.PositionClosed:
		s.trades[position.ID] = s.createTrade(position)
	}
	s.positions[position.ID] = position
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

func (s *Simulator) createTrade(position broker.Position) broker.Trade {
	return broker.Trade{
		ID:        position.ID,
		CreatedAt: s.clock.Now(),
		Asset:     position.Asset,
		Side:      position.Side,
		Size:      position.Size,
		Profit:    profit(position, position.LiquidationPrice),
	}
}

func (s *Simulator) equityNow() decimal.Decimal {
	equity := s.accountBalance
	if position := s.getLatestOrNewPosition(); position.State() == broker.PositionOpen {
		equity = equity.Add(profit(position, s.marketPrice.C))
	}
	return equity
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

func profit(position broker.Position, price decimal.Decimal) decimal.Decimal {
	profit := price.Sub(position.Price).Mul(position.Size)
	if position.Side == broker.Sell {
		profit = profit.Neg()
	}
	return profit
}
