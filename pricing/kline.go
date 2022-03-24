package pricing
import (
	"time"

	"github.com/shopspring/decimal"
)

type Kline struct {
	Start  time.Time
	O      decimal.Decimal
	H      decimal.Decimal
	L      decimal.Decimal
	C      decimal.Decimal
	Volume float64
}
