package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Predicter = (*CrossPredicter)(nil)

type CrossPredicter struct {
	priceSelector ta.PriceSelector

	osc ta.Indicator
	mmi ta.Indicator

	prev float64
}

func NewCrossPredicter(osc, mmi ta.Indicator) *CrossPredicter {
	return &CrossPredicter{
		priceSelector: ta.HL2,
		osc:           osc,
		mmi:           mmi,
	}
}

func (p *CrossPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.priceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.osc.Update(v); err != nil {
		return err
	}
	if err := p.mmi.Update(vDiff); err != nil {
		return err
	}

	return nil
}

func (p *CrossPredicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.mmi.History(), 1), ta.Lookback(p.mmi.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	switch {
	case ta.CrossUp(p.osc.History(), 0):
		score += 0.9
	case ta.CrossDown(p.osc.History(), 0):
		score = -(score + 0.9)
	}

	return score
}

func (p *CrossPredicter) Valid() bool {
	return p.osc.Valid() && p.mmi.Valid()
}
