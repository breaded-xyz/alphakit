package backtest

import (
	"errors"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const _defaultTockInterval = time.Millisecond

var ErrInvalidOrderState = errors.New("order is not valid for processing")
var ErrRejectedOrder = errors.New("order rejected during processing")

type Simulator struct {
	clock       Clocker
	balance     broker.AccountBalance
	marketPrice market.Kline

	cost Coster

	orders    map[broker.DealID]broker.Order
	positions map[broker.DealID]broker.Position
	trades    map[broker.DealID]broker.Trade
	equity    broker.EquitySeries
}

func NewSimulator() *Simulator {
	return NewSimulatorWithCost(NewPerpCost())
}

func NewSimulatorWithCost(cost Coster) *Simulator {
	return &Simulator{
		clock:     NewClock(),
		cost:      cost,
		orders:    make(map[broker.DealID]broker.Order),
		positions: make(map[broker.DealID]broker.Position),
		trades:    make(map[broker.DealID]broker.Trade),
		equity:    make(broker.EquitySeries),
	}
}

func (s *Simulator) SetInitialCapital(amount decimal.Decimal) {
	s.balance.Trade = amount
}

func (s *Simulator) AddOrder(order broker.Order) (broker.Order, error) {
	var empty broker.Order
	if order.Side == 0 || order.Type == 0 || order.State() != broker.OrderPending || !order.Size.IsPositive() {
		return empty, ErrInvalidOrderState
	}
	return s.processOrder(order)
}

func (s *Simulator) Next(price market.Kline) error {

	// Init simulation clock the first time a price is received
	if s.clock.Peek().IsZero() {
		s.clock.Start(price.Start, _defaultTockInterval)
	}

	// Advance the clock epoch to the start time of the kline
	s.clock.Advance(price.Start)

	// Set the market price used in this epoch to the received price
	s.marketPrice = price

	// Deduct funding fees if an existing position is open
	s.balance.Trade = s.balance.Trade.Sub(s.cost.Funding(s.getPosition(), s.marketPrice.C, s.clock.Elapsed()))

	// Iterate open orders in the sequence they were placed (FIFO)
	// Go maps do not maintain insertion order so we must sort the keys in a slice first
	// The map key is a ULID seeded from a time and supports lexicographic sorting
	ks := maps.Keys(s.orders)
	slices.Sort(ks)
	for _, k := range ks {
		order := s.orders[k]
		if order.State() != broker.OrderOpen {
			continue
		}
		if _, err := s.processOrder(order); err != nil {
			return err
		}
	}

	// Add current portfolio equity to the history
	equity := s.markToMarket()
	s.equity[broker.Timestamp(s.clock.Peek().UnixMilli())] = equity
	s.balance.Equity = equity

	return nil
}

func (s *Simulator) CancelOrders() {
	for k := range s.orders {
		order := s.orders[k]
		if order.State() == broker.OrderOpen {
			order.ClosedAt = s.clock.Now()
		}
	}
}

func (s *Simulator) Orders() []broker.Order {
	return copySortMap(s.orders)
}

func (s *Simulator) Positions() []broker.Position {
	return copySortMap(s.positions)
}

func (s *Simulator) Trades() []broker.Trade {
	return copySortMap(s.trades)
}

func (s *Simulator) EquityHistory() broker.EquitySeries {
	copied := make(broker.EquitySeries, len(s.equity))
	maps.Copy(copied, s.equity)
	return copied
}

func (s *Simulator) Balance() broker.AccountBalance {
	return s.balance
}

func (s *Simulator) processOrder(order broker.Order) (broker.Order, error) {
	var err error

	switch order.State() {
	case broker.OrderPending:
		if order, err = s.processOrder(s.openOrder(order)); err != nil {
			return order, err
		}
	case broker.OrderOpen:

		// State transition condition:
		// Guard for temporal logic error whereby a past or future price is used to fill an order
		// Limit orders cannot be filled in the same epoch as the current price
		if order.Type == broker.Limit && equalClock(order.OpenedAt, s.clock.Peek()) {
			break
		}

		// State transition condition:
		// Order price must be matched to the available market price
		// Market type orders will always match the current close price
		matchedPrice := matchOrder(order, s.marketPrice)
		if !matchedPrice.IsPositive() {
			break
		}

		// Transition to filled state
		if order, err = s.processOrder(s.fillOrder(order, matchedPrice)); err != nil {
			return order, err
		}

	case broker.OrderFilled:
		if _, err = s.processPosition(s.getPosition(), order); err != nil {
			return order, err
		}
		s.balance.Trade = s.balance.Trade.Sub(s.cost.Transaction(order))
		order = s.closeOrder(order)
	}

	s.orders[order.ID] = order
	return order, nil
}

func (s *Simulator) openOrder(order broker.Order) broker.Order {
	order.ID = broker.NewIDWithTime(s.clock.Now())
	order.OpenedAt = s.clock.Now()
	return order
}

func (s *Simulator) fillOrder(order broker.Order, matchedPrice decimal.Decimal) broker.Order {
	order.FilledAt = s.clock.Now()
	var fillPrice decimal.Decimal

	switch order.Side {
	case broker.Buy:
		fillPrice = matchedPrice.Add(s.cost.Slippage(matchedPrice))
		fillPrice = fillPrice.Add(s.cost.Spread(fillPrice))
	case broker.Sell:
		fillPrice = matchedPrice.Sub(s.cost.Slippage(matchedPrice))
		fillPrice = fillPrice.Sub(s.cost.Spread(fillPrice))
	}

	order.FilledPrice = fillPrice
	order.FilledSize = order.Size

	return order
}

func (s *Simulator) closeOrder(order broker.Order) broker.Order {
	order.ClosedAt = s.clock.Now()
	return order
}

func (s *Simulator) getPosition() broker.Position {
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

func (s *Simulator) processPosition(position broker.Position, order broker.Order) (broker.Position, error) {
	var err error

	switch position.State() {
	case broker.PositionPending:

		// State transition condition:
		// Do not open a new position with a 'reduce-only' order
		// Reduce-only is typically used for stop loss orders and is only permitted to close a position
		if order.ReduceOnly {
			return position, ErrRejectedOrder
		}

		// Transition to open
		if position, err = s.processPosition(s.openPosition(order), order); err != nil {
			return position, err
		}

	case broker.PositionOpen:

		// State transition condition:
		// If processing the order that opened the position then do not attempt to close position
		if order.ID == position.ID {
			break
		}

		// State transition condition:
		// A new order can only close an opened position in full and never partially reduce or increase
		if !(order.Side == position.Side.Opposite() && order.FilledSize.Equal(position.Size)) {
			return position, ErrRejectedOrder
		}

		// Transition to closed
		if position, err = s.processPosition(s.closePosition(position, order), order); err != nil {
			return position, err
		}

	case broker.PositionClosed:

		// Create a trade for the closed position and update the account balance with the profit / loss
		trade := s.createTrade(position)
		s.balance.Trade = s.balance.Trade.Add(trade.Profit)
		s.trades[position.ID] = trade
	}

	s.positions[position.ID] = position
	return position, nil
}

func (s *Simulator) openPosition(order broker.Order) broker.Position {
	return broker.Position{
		ID:       order.ID,
		OpenedAt: order.FilledAt,
		Asset:    order.Asset,
		Side:     order.Side,
		Price:    order.FilledPrice,
		Size:     order.FilledSize,
	}
}

func (s *Simulator) closePosition(position broker.Position, order broker.Order) broker.Position {
	position.ClosedAt = order.FilledAt
	position.LiquidationPrice = order.FilledPrice
	return position
}

func (s *Simulator) createTrade(position broker.Position) broker.Trade {
	return broker.Trade{
		ID:         position.ID,
		CreatedAt:  position.ClosedAt,
		Asset:      position.Asset,
		Side:       position.Side,
		Size:       position.Size,
		Profit:     profit(position, position.LiquidationPrice),
		HoldPeriod: position.ClosedAt.Sub(position.OpenedAt),
	}
}

func (s *Simulator) markToMarket() decimal.Decimal {
	equity := s.balance.Trade
	if position := s.getPosition(); position.State() == broker.PositionOpen {
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

func profit(position broker.Position, price decimal.Decimal) decimal.Decimal {
	profit := price.Sub(position.Price).Mul(position.Size)
	if position.Side == broker.Sell {
		profit = profit.Neg()
	}
	return profit
}

func equalClock(t1, t2 time.Time) bool {
	t1H, t1M, t1S := t1.Clock()
	t2H, t2M, t2S := t2.Clock()
	return t1H == t2H && t1M == t2M && t1S == t2S
}

func copySortMap[K constraints.Ordered, V any](m map[K]V) []V {
	copied := make([]V, len(m))
	ks := maps.Keys(m)
	slices.Sort(ks)
	for i := range ks {
		copied[i] = m[ks[i]]
	}
	return copied
}
