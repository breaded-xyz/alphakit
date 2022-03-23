package zero2algo

type PriceReceiver interface {
	ReceivePrice(Kline) error
}
