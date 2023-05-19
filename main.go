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

func main() {

	receivedOrdersChannel := receiveOrders()
	validOrdersChannel, inValidOrdersChannel := validateOrders(receivedOrdersChannel)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(validOrdersChannel <-chan model.Order, invalidOrdersChannel <-chan model.InvalidOrder) {
		defer wg.Done()
		fmt.Println("Waiting for Orders")

	loop:
		for {
			select {
			case validOrder, ok := <-validOrdersChannel:
				if ok {
					fmt.Printf("valid order recived -> %v \n", validOrder.String())
				} else {
					break loop
				}
			case invalidOrder, ok := <-inValidOrdersChannel:
				if ok {
					fmt.Printf("invalid order recived -> %v \n", invalidOrder.Order.String())
				} else {
					break loop
				}
			}
		}

	}(validOrdersChannel, inValidOrdersChannel)

	wg.Wait()
}
