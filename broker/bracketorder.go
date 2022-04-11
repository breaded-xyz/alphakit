package broker

type BracketOrder struct {
	Enter Order
	Stop  Order
}
