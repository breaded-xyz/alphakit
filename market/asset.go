package market

type Asset struct {
	Symbol string
}

func NewAsset(symbol string) Asset {
	return Asset{
		Symbol: symbol,
	}
}

func (a *Asset) Equal(other Asset) bool {
	return a.Symbol == other.Symbol
}
