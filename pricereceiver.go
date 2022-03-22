package zero2algo

type PriceReceiver interface {
	Receive(Kline) error
}
