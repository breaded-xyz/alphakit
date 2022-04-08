package broker

import "github.com/shopspring/decimal"

type AccountBalance struct {
	Trade  decimal.Decimal
	Equity decimal.Decimal
}
