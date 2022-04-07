package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
)

type Predicter struct {
	priceSelector ta.PriceSelector
	maFast        ta.Indicator
	maSlow        ta.Indicator
	osc           ta.Osc
	sd            ta.SD
	mmi           ta.MMI

	prev float64
}

func NewPredicter(maFast, maSlow ta.Indicator, sd ta.SD, mmi ta.MMI) *Predicter {
	return nil
}

func (p *Predicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.priceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.osc.Update(v); err != nil {
		return err
	}
	if err := p.sd.Update(v); err != nil {
		return err
	}
	if err := p.mmi.Update(vDiff); err != nil {
		return err
	}

	return nil
}

func (p *Predicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.mmi.History(), 1), ta.Lookback(p.mmi.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	threshold := p.sd.Value()
	upper := threshold
	lower := -threshold

	switch {
	case ta.CrossUp(p.osc.History(), upper):
		return score + 0.9
	case ta.CrossUp(p.osc.History(), 0):
		return score + 0.6
	case ta.CrossDown(p.osc.History(), 0):
		return -(score + 0.6)
	case ta.CrossDown(p.osc.History(), lower):
		return -(score + 0.9)
	}

	return score
}

func (p *Predicter) Valid() bool {
	return p.osc.Valid() && p.sd.Valid() && p.mmi.Valid()
}
