package model

type OrderStatus int

const (
	none OrderStatus = iota
	new
	received
	reserved
	filled
)
