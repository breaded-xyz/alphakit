package day

import (
	"context"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"golang.org/x/exp/slices"
)

var _ trader.Bot = (*Bot)(nil)

type sessionRow struct {
	Day          time.Time
	YLow         float64
	YVAL         float64
	YPOC         float64
	YVAH         float64
	YHigh        float64
	HourClose    float64
	YLowDistPct  float64
	YVALDistPct  float64
	YPOCDistPct  float64
	YVAHDistPct  float64
	YHighDistPct float64
	CrossLow     bool
	CrossVAL     bool
	CrossPOC     bool
	CrossVAH     bool
	CrossHigh    bool
}

// Bot is a trader.Bot implementation for day trading.
type Bot struct {
	Levels   []VolumeLevel
	Profiles []*VolumeProfile
	Results  []sessionRow
}

// NewBot creates a new Bot instance.
func NewBot() *Bot {
	return &Bot{}
}

// SetDealer sets the broker used for placing orders.
func (b *Bot) SetDealer(dealer broker.Dealer) {

}

// Warmup seeds the Predicter and Risker with historical price data.
func (b *Bot) Warmup(ctx context.Context, prices []market.Kline) error {

	return nil
}

// ReceivePrice updates the algo with latest market price potentially triggering buy and/or sell orders.
func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {

	// Add new Level for kline using HL2
	levelsNow := VolumeLevel{
		Price:  util.RoundTo(ta.HL2(price), 1.0),
		Volume: util.RoundTo(price.Volume, 1.0),
	}
	b.Levels = append(b.Levels, levelsNow)

	// Initialize new day
	if h, m, s := price.Start.UTC().Clock(); h == 0 && m == 0 && s == 0 {
		if len(b.Levels) < 1440 {
			return nil
		}
		ySession := b.Levels[len(b.Levels)-1440:]
		slices.SortStableFunc(ySession, func(i, j VolumeLevel) bool {
			return i.Price < j.Price
		})
		vp := NewVolumeProfile(100, ySession)
		b.Profiles = append(b.Profiles, vp)
		b.Results = append(b.Results, sessionRow{
			Day:   price.Start.UTC(),
			YLow:  vp.Low,
			YVAL:  vp.VAL,
			YPOC:  vp.POC,
			YVAH:  vp.VAH,
			YHigh: vp.High,
		})
		return nil
	}

	// Initialize distance to key levels after first hour
	if h, m, s := price.Start.UTC().Clock(); h == 1 && m == 0 && s == 0 {
		if len(b.Results) == 0 {
			return nil
		}
		session := b.Results[len(b.Results)-1]
		session.HourClose = price.C.InexactFloat64()
		session.YLowDistPct = util.RoundTo((session.YLow-session.HourClose)/session.HourClose, 0.01)
		session.YVALDistPct = util.RoundTo((session.YVAL-session.HourClose)/session.HourClose, 0.01)
		session.YPOCDistPct = util.RoundTo((session.YPOC-session.HourClose)/session.HourClose, 0.01)
		session.YVAHDistPct = util.RoundTo((session.YVAH-session.HourClose)/session.HourClose, 0.01)
		session.YHighDistPct = util.RoundTo((session.YHigh-session.HourClose)/session.HourClose, 0.01)
		b.Results[len(b.Results)-1] = session
		return nil
	}

	if h, m, s := price.Start.UTC().Clock(); h > 0 && (h < 23 && m < 59 && s < 59) {

		if len(b.Results) == 0 || len(b.Levels) == 0 {
			return nil
		}

		session := b.Results[len(b.Results)-1]
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YLow) {
			session.CrossLow = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAL) {
			session.CrossVAL = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YPOC) {
			session.CrossPOC = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAH) {
			session.CrossVAH = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YHigh) {
			session.CrossHigh = true
		}
		b.Results[len(b.Results)-1] = session
	}

	return nil
}

func crossed(a, b, c float64) bool {
	return ta.CrossUp([]float64{a, b}, c) || ta.CrossDown([]float64{a, b}, c)
}

// Close exits all open positions at current market price.
func (b *Bot) Close(ctx context.Context) error {

	return nil
}
