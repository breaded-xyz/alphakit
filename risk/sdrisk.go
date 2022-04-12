package risk

import (
	"context"

	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/shopspring/decimal"
)

var _ Risker = (*SDRisk)(nil)

type SDRisk struct {
	sd *ta.SD
}

func NewSDRisk(length int, factor float64) *SDRisk {
	return &SDRisk{
		sd: ta.NewSDWithFactor(length, factor),
	}
}

func (r *SDRisk) ReceivePrice(ctx context.Context, price market.Kline) error {
	return r.sd.Update(price.C.InexactFloat64())
}

func (r *SDRisk) Risk() decimal.Decimal {
	return dec.New(r.sd.Value())
}

func (r *SDRisk) Valid() bool {
	return r.sd.Valid()
}
