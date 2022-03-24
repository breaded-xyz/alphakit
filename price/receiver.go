package price

type Receiver interface {
	ReceivePrice(Kline) error
}
