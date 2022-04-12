package risk

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type FullRisk struct {
	price decimal.Decimal
}

func NewFullRisk() *FullRisk {
	return &FullRisk{}
}

func (r *FullRisk) ReceivePrice(ctx context.Context, price market.Kline) error {
	r.price = price.C
	return nil
}

func (r *FullRisk) Risk() decimal.Decimal {
	return r.price
}

func (r *FullRisk) Valid() bool {
	return true
}
