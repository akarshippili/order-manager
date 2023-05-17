package model

import "fmt"

type Order struct {
	ProductCode int
	Quantity    float64
	Status      OrderStatus
}

func (order *Order) String() string {
	return fmt.Sprintf("Procuct Code: %v, Quantity: %v, Satus: %v\n", order.ProductCode, order.Quantity, order.Status)
}
