package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/akarshippili/order-manager/model"
)

var rawData = []string{
	`{"productCode": 1111, "quantity": -5, "status": 1}`,
	`{"productCode": 2222, "quantity": 10, "status": 1}`,
	`{"productCode": 22, "quantity": -10, "status": 1}`,
	`{"productCode": 3333, "quantity": 100.0002, "status": 1}`,
	`{"productCode": 4444, "quantity": 0.000099, "status": 1}`,
	`{"productCode": 4, "quantity": -0.000099, "status": 1}`,
	`{"productCode": 5555, "quantity": 1000.000001, "status": 1}`,
}

// directional channel which makes the intent of the channels more clear
func receiveOrders() <-chan model.Order {
	receivedOrdersChannel := make(chan model.Order)

	go func() {
		for _, rawOrder := range rawData {
			var order model.Order
			err := json.Unmarshal([]byte(rawOrder), &order)

			if err != nil {
				fmt.Println(err)
				continue
			}
			receivedOrdersChannel <- order
		}
		close(receivedOrdersChannel)
	}()

	return receivedOrdersChannel
}

func validateOrders(receivedOrdersChannel <-chan model.Order) (<-chan model.Order, <-chan model.InvalidOrder) {

	validOrdersChannel := make(chan model.Order)
	inValidOrdersChannel := make(chan model.InvalidOrder, 1)

	go func() {
		for order := range receivedOrdersChannel {
			if order.Quantity <= 0 {
				inValidOrdersChannel <- model.InvalidOrder{
					Order: order,
					Err:   errors.New("order quantity must be greater than 0"),
				}
			} else {
				validOrdersChannel <- order
			}
		}
		close(validOrdersChannel)
		close(inValidOrdersChannel)
	}()

	return validOrdersChannel, inValidOrdersChannel
}

func reserverInventory(in <-chan model.Order) <-chan model.Order {

	out := make(chan model.Order)

	go func() {
		for order := range in {
			order.Status = model.Reserved
			out <- order
		}
		close(out)
	}()

	return out
}

func main() {

	receivedOrdersChannel := receiveOrders()
	validOrdersChannel, inValidOrdersChannel := validateOrders(receivedOrdersChannel)
	reservedOrdersChannel := reserverInventory(validOrdersChannel)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(reservedOrderChannel <-chan model.Order) {
		defer wg.Done()

		for reservedOrder := range reservedOrdersChannel {
			fmt.Printf("Inventory recived for order -> %v \n", reservedOrder.String())
		}
	}(receivedOrdersChannel)

	go func(invalidOrdersChannel <-chan model.InvalidOrder) {
		defer wg.Done()

		for invalidOrder := range inValidOrdersChannel {
			fmt.Printf("invalid order recived -> %v \n", invalidOrder.Order.String())
		}
	}(inValidOrdersChannel)

	wg.Wait()
	fmt.Println("Done processing orders.")
}
