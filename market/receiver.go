package market

import "context"

type Receiver interface {
	ReceivePrice(context.Context, Kline) error
}
