package tradebot

import "github.com/colngroup/zero2algo/pricing"

type Bot interface {
	pricing.Receiver
	Close()
}

type ConfigurableBot interface {
	Bot
	Configure(config map[string]any) error
}
