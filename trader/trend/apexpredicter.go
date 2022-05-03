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

	osc ta.Indicator
	sd  ta.Indicator
	mmi ta.Indicator

	prev float64
}

func NewApexPredicter(osc, sd, mmi ta.Indicator) *ApexPredicter {
	return &ApexPredicter{
		priceSelector: ta.Close,
		osc:           osc,
		sd:            sd,
		mmi:           mmi,
	}
}

func (p *ApexPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

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

func (p *ApexPredicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.mmi.History(), 1), ta.Lookback(p.mmi.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	//if p.osc.Value() > p.sd.Value() || p.osc.Value() < -p.sd.Value() {
	//	score += 0.05
	//}

	switch {
	case ta.Valley(p.osc.History()):
		score = (score + 0.9)
	case ta.Peak(p.osc.History()):
		score = -(score + 0.9)
	}

	/*sd05, sd15, sd2 := sd*0.5, sd*1.5, sd*2

	switch {
	case ta.Valley(p.osc.History()) && p.osc.Value() < 0-sd2:
		return score + 0.9
	case ta.Valley(p.osc.History()) && p.osc.Value() < 0-sd15:
		return score + 0.7
	case ta.Valley(p.osc.History()) && p.osc.Value() < 0-sd:
		return score + 0.5
	case ta.Valley(p.osc.History()) && p.osc.Value() < 0-sd05:
		return score + 0.3
	case ta.Valley(p.osc.History()) && p.osc.Value() < 0:
		return score + 0.1
	case ta.Peak(p.osc.History()) && p.osc.Value() > 0+sd2:
		return -(score + 0.9)
	case ta.Peak(p.osc.History()) && p.osc.Value() > 0+sd15:
		return -(score + 0.7)
	case ta.Peak(p.osc.History()) && p.osc.Value() > 0+sd:
		return -(score + 0.5)
	case ta.Peak(p.osc.History()) && p.osc.Value() > 0+sd05:
		return -(score + 0.3)
	case ta.Peak(p.osc.History()) && p.osc.Value() > 0:
		return -(score + 0.1)
	}*/

	return score
}

func (p *ApexPredicter) Valid() bool {
	return p.osc.Valid() && p.sd.Valid() && p.mmi.Valid()
}
