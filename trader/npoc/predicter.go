package npoc

import (
	"context"
	"time"

	"github.com/thecolngroup/alphakit/internal/util"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/ta"
	"github.com/thecolngroup/alphakit/trader"
	"golang.org/x/exp/slices"
)

var _ trader.Predicter = (*Predicter)(nil)

// Predicter predicts the reclaim of the previous session's Point of Control.
type Predicter struct {
	// PriceSelector is the kline component to use for price. Close by default.
	PriceSelector ta.PriceSelector

	// PrevSessionProfile is the volume profile for yesterdays session.
	PrevSessionProfile *ta.VolumeProfile

	SessionStart     time.Time
	SessionOpenClose time.Time

	KlineCountPerDay int

	prices          []market.Kline
	insideVA        bool
	nakedPOC        bool
	nakedPOCDistPct float64
	score           float64
}

// NewPredicter creates a new predicter with Close quote price selector.
func NewPredicter(osc, mmi ta.Indicator[float64]) *Predicter {
	return &Predicter{
		PriceSelector:    ta.Close,
		SessionStart:     time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
		SessionOpenClose: time.Date(0, 0, 0, 1, 0, 0, 0, time.UTC),
		KlineCountPerDay: 1440,
	}
}

// ReceivePrice updates the prediction algo with the next market price.
// Call Predict() to get the resulting score.
func (p *Predicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	p.score = 0

	p.prices = append(p.prices, price)
	h, m, s := price.Start.Clock()

	if p.PrevSessionProfile != nil {
		h, l := price.H.InexactFloat64(), price.L.InexactFloat64()
		p.nakedPOC = !util.Between(p.PrevSessionProfile.POC, l, h)
	}

	switch {
	case h == p.SessionStart.Hour(), m == p.SessionStart.Minute(), s == 0:

		if len(p.prices) < p.KlineCountPerDay {
			return nil
		}
		p.nakedPOC = true
		p.insideVA = false

		p.prices = p.prices[len(p.prices)-p.KlineCountPerDay:]
		vls := make([]ta.VolumeLevel, len(p.prices))
		for i := range p.prices {
			vls[i] = ta.VolumeLevel{
				Price:  ta.HLC3(p.prices[i]),
				Volume: p.prices[i].Volume,
			}
		}
		slices.SortStableFunc(vls, func(i, j ta.VolumeLevel) bool {
			return i.Price < j.Price
		})
		p.PrevSessionProfile = ta.NewVolumeProfile(25, vls)
		p.insideVA = util.Between(price.O.InexactFloat64(), p.PrevSessionProfile.Low, p.PrevSessionProfile.High)

	case h == p.SessionOpenClose.Hour(), m == p.SessionOpenClose.Minute(), s == 0:
		priceClose := price.C.InexactFloat64()
		p.nakedPOCDistPct = (p.PrevSessionProfile.POC - priceClose) / priceClose

		if p.nakedPOC && p.insideVA {
			switch {
			case p.nakedPOCDistPct >= 1.0 && p.nakedPOCDistPct <= 2.0:
				p.score = 1.0
			case p.nakedPOCDistPct <= -1.0 && p.nakedPOCDistPct >= -2.0:
				p.score = -1.0
			}
		}

	}

	return nil
}

// Predict returns a score to indicate confidence of price direction.
//
// 1.0 = Long: reclaim PoC at higher price level.
//
// -1.0 = Short: reclaim PoC at lower price level.
//
// 0.0 = No prediction
func (p *Predicter) Predict() float64 {

	var score float64

	return score
}

// Valid returns true if the prediction is valid.
func (p *Predicter) Valid() bool {
	return true
}
