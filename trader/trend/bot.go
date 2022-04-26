package trend

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/money"
	"github.com/colngroup/zero2algo/risk"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"github.com/shopspring/decimal"
)

var _ trader.ConfigurableBot = (*Bot)(nil)

type Bot struct {
	EnterLong float64
	ExitLong  float64

	EnterShort float64
	ExitShort  float64

	Predicter *Predicter
	Risker    risk.Risker
	Sizer     money.Sizer

	asset  market.Asset
	dealer broker.Dealer
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) SetDealer(dealer broker.Dealer) {
	b.dealer = dealer
}

func (b *Bot) Warmup(ctx context.Context, prices []market.Kline) error {
	for i := range prices {
		if err := b.updateIndicators(ctx, prices[i]); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {

	if err := b.updateIndicators(ctx, price); err != nil {
		return err
	}
	if !(b.Predicter.Valid() && b.Risker.Valid()) {
		return nil
	}

	enterSide, exitSide := b.signal(b.Predicter.Predict())

	if err := b.exit(ctx, exitSide); err != nil {
		return err
	}

	if err := b.enter(ctx, enterSide, price.C); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Close(ctx context.Context) error {

	if err := b.exit(ctx, broker.Buy); err != nil {
		return err
	}
	if err := b.exit(ctx, broker.Sell); err != nil {
		return err
	}

	return nil
}

func (b *Bot) updateIndicators(ctx context.Context, price market.Kline) error {
	if err := b.Risker.ReceivePrice(ctx, price); err != nil {
		return err
	}
	if err := b.Predicter.ReceivePrice(ctx, price); err != nil {
		return err
	}
	return nil
}

func (b *Bot) signal(prediction float64) (enter, exit broker.OrderSide) {

	switch {
	case prediction == 0:
		return
	case prediction >= b.EnterLong:
		return broker.Buy, broker.Sell
	case prediction >= b.ExitShort:
		return 0, broker.Sell
	case prediction <= b.EnterShort:
		return broker.Sell, broker.Buy
	case prediction <= b.ExitLong:
		return 0, broker.Buy
	}

	return enter, exit
}

func (b *Bot) getOpenPosition(ctx context.Context, side broker.OrderSide) (broker.Position, error) {
	var empty broker.Position

	positions, _, err := b.dealer.ListPositions(ctx, nil)
	if err != nil {
		return empty, err
	}
	opens := filterPositions(positions, b.asset, side, broker.PositionOpen)
	if len(opens) == 0 {
		return empty, err
	}

	return opens[len(opens)-1], nil
}

func (b *Bot) exit(ctx context.Context, exitSide broker.OrderSide) error {
	if exitSide == 0 {
		return nil
	}

	position, err := b.getOpenPosition(ctx, exitSide)
	if err != nil {
		return err
	}
	if position.State() == broker.PositionOpen {
		if _, err := b.executeExitOrder(ctx, exitSide, position.Size); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) executeExitOrder(ctx context.Context, side broker.OrderSide, size decimal.Decimal) (broker.Order, error) {
	var empty broker.Order

	if _, err := b.dealer.CancelOrders(ctx); err != nil {
		return empty, err
	}

	order := broker.Order{
		Asset:      b.asset,
		Type:       broker.Market,
		Side:       side.Opposite(),
		Size:       size,
		ReduceOnly: true,
	}
	placed, _, err := b.dealer.PlaceOrder(ctx, order)
	if err != nil {
		return empty, err
	}
	if placed == nil {
		return empty, nil
	}

	return *placed, err
}

func (b *Bot) enter(ctx context.Context, enterSide broker.OrderSide, price decimal.Decimal) error {
	if enterSide == 0 {
		return nil
	}

	position, err := b.getOpenPosition(ctx, enterSide)
	if err != nil {
		return err
	}
	if position.State() == broker.PositionOpen {
		return nil
	}

	balance, _, err := b.dealer.GetBalance(ctx)
	if err != nil {
		return err
	}
	capital := balance.Trade
	unitaryRisk := b.Risker.Risk()
	size := b.Sizer.Size(price, capital, unitaryRisk)
	_, err = b.executeEnterOrder(ctx, enterSide, price, size, unitaryRisk)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) executeEnterOrder(ctx context.Context, side broker.OrderSide, price, size, risk decimal.Decimal) (broker.BracketOrder, error) {
	var bracket, empty broker.BracketOrder

	if _, err := b.dealer.CancelOrders(ctx); err != nil {
		return empty, err
	}

	order := broker.Order{
		Asset: b.asset,
		Type:  broker.Market,
		Side:  side,
		Size:  size,
	}
	enterPlaced, _, err := b.dealer.PlaceOrder(ctx, order)
	if err != nil {
		return empty, err
	}
	bracket = broker.BracketOrder{Enter: *enterPlaced}

	stop := broker.Order{
		Asset:      b.asset,
		Type:       broker.Limit,
		Side:       side.Opposite(),
		Size:       size,
		ReduceOnly: true,
	}
	if side == broker.Buy {
		stop.LimitPrice = price.Sub(risk)
	} else {
		stop.LimitPrice = price.Add(risk)
	}
	if !stop.LimitPrice.IsPositive() {
		return bracket, nil
	}
	stopPlaced, _, err := b.dealer.PlaceOrder(ctx, stop)
	if err != nil {
		return bracket, err
	}
	bracket.Stop = *stopPlaced

	return bracket, nil
}

func (b *Bot) Configure(config map[string]any) error {

	b.asset = market.NewAsset(config["asset"].(string))

	b.EnterLong = config["enterlong"].(float64)
	b.EnterShort = config["entershort"].(float64)
	b.ExitLong = config["exitlong"].(float64)
	b.ExitShort = config["exitshort"].(float64)

	maFastLength := util.ToInt(config["mafastlength"])
	maSlowLength := util.ToInt(config["maslowlength"])
	if maFastLength >= maSlowLength {
		return trader.ErrInvalidConfig
	}
	maOsc := ta.NewOsc(ta.NewALMA(maFastLength), ta.NewALMA(maSlowLength))
	maSDFilter := ta.NewSDWithFactor(util.ToInt(config["masdfilterlength"]), config["masdfilterfactor"].(float64))
	mmi := ta.NewMMIWithSmoother(util.ToInt(config["mmilength"]), ta.NewALMA(util.ToInt(config["mmismootherlength"])))
	b.Predicter = NewPredicter(maOsc, maSDFilter, mmi)

	riskSDLength := util.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		b.Risker = risk.NewSDRisk(util.ToInt(config["riskersdlength"]), config["riskersdfactor"].(float64))
	} else {
		b.Risker = risk.NewFullRisk()
	}

	initialCapital := dec.New(config["initialcapital"].(float64))
	sizerF := config["sizerf"].(float64)
	if sizerF > 0 {
		b.Sizer = money.NewSafeFSizer(initialCapital, sizerF, config["sizerscalef"].(float64))
	} else {
		b.Sizer = money.NewFixedSizer(initialCapital)
	}

	return nil
}

func filterPositions(positions []broker.Position, asset market.Asset, side broker.OrderSide, state broker.PositionState) []broker.Position {

	filtered := make([]broker.Position, 0, len(positions))
	for i := range positions {
		p := positions[i]
		if p.Asset.Equal(asset) && p.Side == side && p.State() == state {
			filtered = append(filtered, p)
		}
	}

	return filtered
}
