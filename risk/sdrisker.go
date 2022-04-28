package risk

import (
	"context"

	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/shopspring/decimal"
)

var _ Risker = (*SDRisker)(nil)

type SDRisker struct {
	sd *ta.SD
}

func NewSDRisker(length int, factor float64) *SDRisker {
	return &SDRisker{
		sd: ta.NewSDWithFactor(length, factor),
	}
}

func (r *SDRisker) ReceivePrice(ctx context.Context, price market.Kline) error {
	return r.sd.Update(price.C.InexactFloat64())
}

func (r *SDRisker) Risk() decimal.Decimal {
	return dec.New(r.sd.Value())
}

func (r *SDRisker) Valid() bool {
	return r.sd.Valid()
}
