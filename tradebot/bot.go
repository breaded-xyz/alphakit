package tradebot

import "github.com/colngroup/zero2algo/price"

type Bot interface {
	price.Receiver
	Close()
}
