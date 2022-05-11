package day

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"golang.org/x/exp/slices"
)

var _ trader.Bot = (*Bot)(nil)

func NewTable() dataframe.DataFrame {
	return dataframe.New(
		series.New(nil, series.Int, "session"),
		series.New(nil, series.Float, "yLow"),
		series.New(nil, series.Float, "yVAL"),
		series.New(nil, series.Float, "yPOC"),
		series.New(nil, series.Float, "yVAH"),
		series.New(nil, series.Float, "yHigh"),
		series.New(nil, series.Float, "hourClose"),
		series.New(nil, series.Bool, "crossLow"),
		series.New(nil, series.Bool, "crossVAL"),
		series.New(nil, series.Bool, "crossPOC"),
		series.New(nil, series.Bool, "crossVAH"),
		series.New(nil, series.Bool, "crossHigh"),
	)
}

// Bot is a trader.Bot implementation for day trading.
type Bot struct {
	Levels   []VolumeLevel
	Profiles []*VolumeProfile
	Results  dataframe.DataFrame
}

// NewBot creates a new Bot instance.
func NewBot() *Bot {
	return &Bot{
		Results: NewTable(),
	}
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
	b.Levels = append(b.Levels, VolumeLevel{
		Price:  util.RoundTo(ta.HL2(price), 1.0),
		Volume: util.RoundTo(price.Volume, 1.0),
	})

	// If hour, minute, second is 0
	if h, m, s := price.Start.UTC().Clock(); h == 0 && m == 0 && s == 0 {
		if len(b.Levels) >= 1440 {
			ySession := b.Levels[len(b.Levels)-1440:]
			slices.SortStableFunc(ySession, func(i, j VolumeLevel) bool {
				return i.Price < j.Price
			})
			b.Profiles = append(b.Profiles, NewVolumeProfile(100, ySession))
		}
	}

	// Create new table row for today with yPOC, yVAH, yVAL, yHigh, yLow
	//b.Table.Set(nil)
	//b.Table.

	// If hour is 1 and minute, second is 0
	// Record close price in table for opening hour

	// Record percentage change from hour close price to yPOC, yVAH, yVAL, yHigh, yLow

	// If clock > 1:00:00 and clock < 23:59:59
	// Check if any key volume profile levels crossed and record in table time
	// crossed(levels, yPOC)

	return nil
}

// Close exits all open positions at current market price.
func (b *Bot) Close(ctx context.Context) error {

	return nil
}
