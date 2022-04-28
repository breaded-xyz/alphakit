package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Predicter = (*BreakoutPredicter)(nil)

type BreakoutPredicter struct {
	priceSelector ta.PriceSelector

	osc ta.Indicator
	sd  ta.Indicator
	mmi ta.Indicator

	prev float64
}

func NewBreakoutPredicter(osc, sd, mmi ta.Indicator) *BreakoutPredicter {
	return &BreakoutPredicter{
		priceSelector: ta.Close,
		osc:           osc,
		sd:            sd,
		mmi:           mmi,
	}
}

func (p *BreakoutPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.priceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.osc.Update(v); err != nil {
		return err
	}
	if err := p.sd.Update(p.osc.Value()); err != nil {
		return err
	}
	if err := p.mmi.Update(vDiff); err != nil {
		return err
	}

	return nil
}

func (p *BreakoutPredicter) Predict() float64 {

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
	case ta.CrossUp(p.osc.History(), lower):
		return score + 0.2
	case ta.CrossDown(p.osc.History(), upper):
		return -(score + 0.2)
	case ta.CrossDown(p.osc.History(), 0):
		return -(score + 0.6)
	case ta.CrossDown(p.osc.History(), lower):
		return -(score + 0.9)
	}

	return score
}

func (p *BreakoutPredicter) Valid() bool {
	return p.osc.Valid() && p.sd.Valid() && p.mmi.Valid()
}
