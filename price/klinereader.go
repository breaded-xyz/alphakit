package price

type KlineReader interface {
	Read() (Kline, error)
}
