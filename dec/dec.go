package dec

import (
	"github.com/shopspring/decimal"
)

type number interface {
	~int | ~int64 | ~float64
}

func New[T number](v T) decimal.Decimal {
	var dec decimal.Decimal

	switch (any)(v).(type) {
	case int:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(int64(v))
	case float64:
		return decimal.NewFromFloat(float64(v))
	}

	return dec
}

func Between(v, lower, upper decimal.Decimal) bool {
	return v.GreaterThanOrEqual(lower) && v.LessThanOrEqual(upper)
}
