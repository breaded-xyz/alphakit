package pricing

type Receiver interface {
	ReceivePrice(Kline) error
}
