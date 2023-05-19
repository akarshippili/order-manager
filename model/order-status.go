package model

type OrderStatus int

const (
	None OrderStatus = iota
	New
	Received
	Reserved
	Filled
)
