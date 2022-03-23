package zero2algo

import (
	"time"

	"github.com/shopspring/decimal"
)

type Equity struct {
	At     time.Time
	Amount decimal.Decimal
}
