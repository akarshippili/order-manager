package model

type InvalidOrder struct {
	Order Order
	Err   error
}
