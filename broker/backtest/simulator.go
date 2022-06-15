package backtest

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/gou/dec"
	"golang.org/x/exp/maps"
)

const _defaultTockInterval = time.Millisecond

const _defaultInitialCapital = 1000

// ErrInvalidOrderState is returned when an order is not in a valid state for the simulator to open it.
var ErrInvalidOrderState = errors.New("order is not valid for processing")

// ErrRejectedOrder is returned when an order is rejected during processing due to an exceptional condition.
var ErrRejectedOrder = errors.New("order rejected during processing")

// Simulator is a backtest simulator that simulates the execution of orders against a market.
// Only a single position can be opened at a time, and must be closed in full before another can be opened.
// Partial fills are not supported.
// Account balance can go negative and trading will continue.
// Inspect the equity curve to understand equity change over time and the capital requirements for the algo.
// To advance the simulation call Next() with the next market price.
// Market and Limit orders are supported. Market orders execute immediately with the last available close price.
// To place a stop loss or take profit style order use a limit order with 'ReduceOnly' set to true.
type Simulator struct {
	clock       Clocker
	balance     broker.AccountBalance
	marketPrice market.Kline

	cost Coster

	orders     []broker.Order
	positions  []broker.Position
	roundturns []broker.RoundTurn
	equity     broker.EquitySeries
}

// NewSimulator create a new backtest simulator with zero cost model.
func NewSimulator() *Simulator {
	return NewSimulatorWithCost(NewPerpCoster())
}

// NewSimulatorWithCost creates a new backtest simulator with the given cost model.
func NewSimulatorWithCost(cost Coster) *Simulator {
	return &Simulator{
		balance: broker.AccountBalance{},
		clock:   NewClock(),
		cost:    cost,
		equity:  make(broker.EquitySeries),
	}
}

// SetInitialCapital sets the initial trade balance.
func (s *Simulator) SetInitialCapital(amount decimal.Decimal) {
	s.balance.Trade = amount
}

// AddOrder adds an order to the simulator and returns the processed order or an error.
func (s *Simulator) AddOrder(order broker.Order) (broker.Order, error) {
	var empty broker.Order
	if order.Side == 0 || order.Type == 0 || order.State() != broker.OrderPending || !order.Size.IsPositive() {
		return empty, ErrInvalidOrderState
	}
	order, err := s.processOrder(order)
	if err != nil {
		return empty, err
	}

	s.orders = append(s.orders, order)

	return order, nil
}

// Next advances the simulation by one kline.
func (s *Simulator) Next(price market.Kline) error {

	// Init simulation clock the first time a price is received
	if s.clock.Peek().IsZero() {
		s.clock.Start(price.Start, _defaultTockInterval)
	}

	// Advance the clock epoch to the start time of the kline
	s.clock.Advance(price.Start)

	// Set the market price used in this epoch to the received price
	s.marketPrice = price

	for i := range s.orders {
		order := s.orders[i]
		if order.State() != broker.OrderOpen {
			continue
		}
		order, err := s.processOrder(order)
		if err != nil {
			return err
		}
		s.orders[i] = order
	}

	// Init equity balance with trade (realized cash) balance
	equity := s.balance.Trade

	// Mark open position to market and add unrealized PNL to equity
	if position := s.getPosition(); position.State() == broker.OrderOpen {
		// Deduct funding fees from position PNL
		position.Cost = position.Cost.Add(s.cost.Funding(position, s.marketPrice.C, s.clock.Elapsed()))
		// Mark position PNL to latest price
		position = markPositionToMarket(position, s.marketPrice.C)
		s.upsertPosition(position)

		equity = equity.Add(position.PNL)
	}

	// Update equity balance
	s.equity[broker.Timestamp(s.clock.Peek().UnixMilli())] = equity
	s.balance.Equity = equity

	return nil
}

// CancelOrders cancels all open orders and returns the cancelled orders.
func (s *Simulator) CancelOrders() []broker.Order {
	cancelled := make([]broker.Order, 0, len(s.orders))
	for i := range s.orders {
		order := s.orders[i]
		if order.State() == broker.OrderOpen {
			order.ClosedAt = s.clock.Now()
			cancelled = append(cancelled, order)
			s.orders[i] = order
		}
	}
	return cancelled
}

// Orders returns a copy of all historical and open orders.
func (s *Simulator) Orders() []broker.Order {
	copied := make([]broker.Order, len(s.orders))
	copy(copied, s.orders)
	return s.orders
}

// Positions returns a copy of all historical and open positions.
func (s *Simulator) Positions() []broker.Position {
	copied := make([]broker.Position, len(s.positions))
	copy(copied, s.positions)
	return copied
}

// RoundTurns returns a copy of all historical roundturns.
func (s *Simulator) RoundTurns() []broker.RoundTurn {
	copied := make([]broker.RoundTurn, len(s.roundturns))
	copy(copied, s.roundturns)
	return copied
}

// EquityHistory returns a copy of the equity curve.
func (s *Simulator) EquityHistory() broker.EquitySeries {
	copied := make(broker.EquitySeries, len(s.equity))
	maps.Copy(copied, s.equity)
	return copied
}

// Balance returns the current account balance.
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
		position, err := s.processPosition(s.getPosition(), order)
		if err != nil {
			return order, err
		}
		s.upsertPosition(position)
		order = s.closeOrder(order)
	}
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
	order.Fee = s.cost.Transaction(order)

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
	position = s.positions[len(s.positions)-1]
	if position.State() == broker.PositionClosed {
		return empty
	}
	return position
}

func (s *Simulator) upsertPosition(position broker.Position) {
	if len(s.positions) == 0 {
		s.positions = append(s.positions, position)
		return
	}
	last := s.positions[len(s.positions)-1]
	if position.ID == last.ID {
		s.positions[len(s.positions)-1] = position
		return
	}
	s.positions = append(s.positions, position)
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
		position = s.openPosition(order)
		if position, err = s.processPosition(position, order); err != nil {
			return position, err
		}

	case broker.PositionOpen:

		// State transition condition:
		// If processing the order that opened the position then do not attempt to close position
		if order.ID == position.ID {
			break
		}

		// State transition condition:
		// A new order can only adjust down an opened position to zero, it cannot be forced negative
		if order.Side == position.Side.Opposite() && order.FilledSize.GreaterThan(position.Size) {
			return position, ErrRejectedOrder
		}

		position = s.adjustPosition(position, order)

		// Closed position
		if position.Size.IsZero() {
			position = s.closePosition(position, order)
			if position, err = s.processPosition(position, order); err != nil {
				return position, err
			}
		}

	case broker.PositionClosed:
		// Create a round-turn for the closed position
		// Realize the position PNL to the account balance
		// Mark price is the fill price of the order that closed the position,
		// note this is a specific edge case for closing a position in order to correctly handle limit orders
		position = markPositionToMarket(position, order.FilledPrice)
		roundturn := s.createRoundTurn(position)
		s.balance.Trade = s.balance.Trade.Add(roundturn.Profit)
		s.roundturns = append(s.roundturns, roundturn)
	}

	return position, nil
}

func (s *Simulator) openPosition(order broker.Order) broker.Position {
	position := broker.Position{
		ID:       order.ID,
		OpenedAt: order.FilledAt,
		Asset:    order.Asset,
		Side:     order.Side,
	}

	return s.adjustPosition(position, order)
}

func (s *Simulator) adjustPosition(position broker.Position, order broker.Order) broker.Position {

	position.TradeCount++

	orderCost := order.FilledSize.Mul(order.FilledPrice)

	switch position.Side {
	case order.Side:
		position.Cost = position.Cost.Add(orderCost).Add(order.Fee)
		position.Size = position.Size.Add(order.FilledSize)
	case order.Side.Opposite():
		position.Cost = position.Cost.Sub(orderCost).Add(order.Fee)
		position.Size = position.Size.Sub(order.FilledSize)
	}

	position.EntryPrice = position.Cost.Div(dec.NZ(position.Size, dec.New(1))).Abs()

	return position
}

func (s *Simulator) closePosition(position broker.Position, order broker.Order) broker.Position {
	position.ClosedAt = order.FilledAt
	position.ExitPrice = order.FilledPrice
	return position
}

func (s *Simulator) createRoundTurn(position broker.Position) broker.RoundTurn {
	return broker.RoundTurn{
		ID:         position.ID,
		CreatedAt:  position.ClosedAt,
		Asset:      position.Asset,
		Side:       position.Side,
		Profit:     position.PNL,
		HoldPeriod: position.ClosedAt.Sub(position.OpenedAt),
		TradeCount: position.TradeCount,
	}
}

func markPositionToMarket(position broker.Position, markPrice decimal.Decimal) broker.Position {
	position.MarkPrice = markPrice
	position.PNL = position.Size.Mul(position.MarkPrice).Sub(position.Cost)
	if position.Side == broker.Sell {
		position.PNL = position.PNL.Mul(dec.New(-1))
	}
	return position
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

func equalClock(t1, t2 time.Time) bool {
	t1H, t1M, t1S := t1.Clock()
	t2H, t2M, t2S := t2.Clock()
	return t1H == t2H && t1M == t2M && t1S == t2S
}
