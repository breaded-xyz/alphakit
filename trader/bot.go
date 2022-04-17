package trader

import (
	"context"
	"errors"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
)

var ErrInvalidConfig = errors.New("invalid bot config")

type Bot interface {
	market.Receiver
	SetDealer(broker.Dealer)
	Close(context.Context) error
}

type ConfigurableBot interface {
	Bot
	Configure(config map[string]any) error
}

type BotMakerFunc func() ConfigurableBot
