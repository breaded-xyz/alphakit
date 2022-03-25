package tradebot

import (
	"errors"

	"github.com/colngroup/zero2algo/pricing"
)

var ErrInvalidConfig = errors.New("invalid bot config")

type Bot interface {
	pricing.Receiver
	Close() error
}

type ConfigurableBot interface {
	Bot
	Configure(config map[string]any) error
}
