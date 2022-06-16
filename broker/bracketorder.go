// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package broker

// BracketOrder groups together a set of dependent orders to open and manage a new position.
type BracketOrder struct {
	Enter Order
	Stop  Order
}
