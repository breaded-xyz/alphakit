// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package market

import "context"

// Receiver is a core interface implemented by many types that receive market data.
type Receiver interface {
	ReceivePrice(context.Context, Kline) error
}
