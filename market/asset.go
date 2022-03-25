package market

type Asset struct {
	Symbol string
}

func NewAsset(symbol string) Asset {
	return Asset{
		Symbol: symbol,
	}
}
