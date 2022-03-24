package zero2algo

type KlineReader interface {
	Read() (Kline, error)
}
