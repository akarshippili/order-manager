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

var receivedOrdersChannel = make(chan model.Order)
var validOrdersChannel = make(chan model.Order)
var inValidOrdersChannel = make(chan model.InvalidOrder)

func receiveOrders() {
	for _, rawOrder := range rawData {
		var order model.Order
		err := json.Unmarshal([]byte(rawOrder), &order)

		if err != nil {
			fmt.Println(err)
			continue
		}

		receivedOrdersChannel <- order
	}
}

func validateOrders() {
	order := <-receivedOrdersChannel
	if order.Quantity <= 0 {
		inValidOrdersChannel <- model.InvalidOrder{
			Order: order,
			Err:   errors.New("order quantity must be greater than 0"),
		}
	} else {
		validOrdersChannel <- order
	}
}

func main() {

	go receiveOrders()
	go validateOrders()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		fmt.Println("Waiting for Order")

		select {
		case validOrder := <-validOrdersChannel:
			fmt.Printf("valid order recived -> %v \n", validOrder.String())
		case invalidOrder := <-inValidOrdersChannel:
			fmt.Printf("invalid order recived -> %v \n", invalidOrder.Order.String())
		}
	}()

	wg.Wait()

}
