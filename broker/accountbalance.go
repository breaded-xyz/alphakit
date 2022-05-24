package broker

import "github.com/shopspring/decimal"

// AccountBalance is a representation of a broker's account balance.
type AccountBalance struct {

	// Trade is the balance amount available for trading.
	Trade decimal.Decimal

	// Equity is the total notional account value including unrealized gains on open positions.
	Equity decimal.Decimal
}
