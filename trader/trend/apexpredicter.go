package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Predicter = (*ApexPredicter)(nil)

// ApexPredicter predicts price direction based on trend turning points
// with a market meaness index filter.
// Peak signals the start of a downward price trend.
// Valley signal the start of an upward price trend.
type ApexPredicter struct {
	// PriceSelector is the kline component to use for price. Close by default.
	PriceSelector ta.PriceSelector

	// MA is the smoothed price series to evaluate for a peak or valley.
	MA ta.Indicator

	// MMI is the trend filter.
	MMI ta.Indicator

	prev float64
}

// NewApexPredicter creates a new predicter with Close quote price selector.
func NewApexPredicter(ma, mmi ta.Indicator) *ApexPredicter {
	return &ApexPredicter{
		PriceSelector: ta.Close,
		MA:            ma,
		MMI:           mmi,
	}
}

// ReceivePrice updates the prediction algo with the next market price.
// Call Predict() to get the resulting score.
func (p *ApexPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.PriceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.MA.Update(v); err != nil {
		return err
	}
	if err := p.MMI.Update(vDiff); err != nil {
		return err
	}

	return nil
}

// Predict returns a score to indicate confidence of price direction.
//
// 1.0 = Valley with MMI confluence.
//
// 0.9 = Valley (no MMI confluence).
//
// -0.9 = Peak (no MMI confluence).
//
// -1.0 = Peak with MMI confluence.
//
// [0.0, 0.1] = Flat trend.
func (p *ApexPredicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.MMI.History(), 1), ta.Lookback(p.MMI.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	switch {
	case ta.Valley(p.MA.History()):
		score = (score + 0.9)
	case ta.Peak(p.MA.History()):
		score = -(score + 0.9)
	}

	return score
}

func (p *ApexPredicter) Valid() bool {
	return p.MA.Valid() && p.MMI.Valid()
}
