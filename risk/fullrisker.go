package risk

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type FullRisker struct {
	price decimal.Decimal
}

func NewFullRisker() *FullRisker {
	return &FullRisker{}
}

func (r *FullRisker) ReceivePrice(ctx context.Context, price market.Kline) error {
	r.price = price.C
	return nil
}

func (r *FullRisker) Risk() decimal.Decimal {
	return r.price
}

func (r *FullRisker) Valid() bool {
	return true
}
