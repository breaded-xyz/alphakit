package market

// KlineReader is an interface for reading candlesticks.
type KlineReader interface {
	Read() (Kline, error)
	ReadAll() ([]Kline, error)
}
