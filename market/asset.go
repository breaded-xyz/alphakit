package market

// Asset represents a tradeable asset identified by its symbol, e.g. BTCUSD.
type Asset struct {
	Symbol string
}

// NewAsset creates a new Asset with the given symbol.
func NewAsset(symbol string) Asset {
	return Asset{
		Symbol: symbol,
	}
}

// Equal asserts equality based on the symbol.
func (a *Asset) Equal(other Asset) bool {
	return a.Symbol == other.Symbol
}
