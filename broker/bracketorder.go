package broker

type BracketOrder struct {
	Primary Order
	Stop    Order
}
