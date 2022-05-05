package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Predicter = (*ApexPredicter)(nil)

type ApexPredicter struct {
	priceSelector ta.PriceSelector

	ma  ta.Indicator
	mmi ta.Indicator

	prev float64
}

func NewApexPredicter(ma, mmi ta.Indicator) *ApexPredicter {
	return &ApexPredicter{
		priceSelector: ta.Close,
		ma:            ma,
		mmi:           mmi,
	}
}

func (p *ApexPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.priceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.ma.Update(v); err != nil {
		return err
	}
	if err := p.mmi.Update(vDiff); err != nil {
		return err
	}

	return nil
}

func (p *ApexPredicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.mmi.History(), 1), ta.Lookback(p.mmi.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	switch {
	case ta.Valley(p.ma.History()):
		score = (score + 0.9)
	case ta.Peak(p.ma.History()):
		score = -(score + 0.9)
	}

	return score
}

func (p *ApexPredicter) Valid() bool {
	return p.ma.Valid() && p.mmi.Valid()
}
