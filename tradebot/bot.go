package tradebot

import "github.com/colngroup/zero2algo/price"

type Bot interface {
	price.Receiver
	Close()
}

type ConfigurableBot interface {
	Bot
	Configure(config map[string]any) error
}
