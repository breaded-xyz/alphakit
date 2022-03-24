package pricing

type KlineReader interface {
	Read() (Kline, error)
	ReadAll() ([]Kline, error)
}
