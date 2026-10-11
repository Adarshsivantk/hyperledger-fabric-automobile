
package main

import "fmt"

func main() {
	// Use these functions to evaluate and submit transactions.

	// Create a car
	result := submitTxnFn(
		"org1",
		"autochannel",
		"KBA-Automobile",
		"CarContract",
		"invoke",
		make(map[string][]byte),
		"CreateCar",
		"Car-07",
		"Tata",
		"Harrier",
		"Black",
		"fac01",
		"25/10/2023",
	)

	// Create an order using private data
	// privateData := map[string][]byte{
	// 	"make":       []byte("Maruti"),
	// 	"model":      []byte("Alto"),
	// 	"color":      []byte("Red"),
	// 	"dealerName": []byte("Popular"),
	// }
	// result := submitTxnFn(
	// 	"org2", "autochannel", "KBA-Automobile",
	// 	"OrderContract", "private", privateData,
	// 	"CreateOrder", "ORD-03",
	// )

	// Read an order
	// result := submitTxnFn(
	// 	"org2", "autochannel", "KBA-Automobile",
	// 	"OrderContract", "query", make(map[string][]byte),
	// 	"ReadOrder", "ORD-05",
	// )

	// Get all cars
	// result := submitTxnFn(
	// 	"org1", "autochannel", "KBA-Automobile",
	// 	"CarContract", "query", make(map[string][]byte),
	// 	"GetAllCars",
	// )

	// Get all orders
	// result := submitTxnFn(
	// 	"org1", "autochannel", "KBA-Automobile",
	// 	"OrderContract", "query", make(map[string][]byte),
	// 	"GetAllOrders",
	// )

	// Get matching orders for a car
	// result := submitTxnFn(
	// 	"org1", "autochannel", "KBA-Automobile",
	// 	"CarContract", "query", make(map[string][]byte),
	// 	"GetMatchingOrders", "Car-06",
	// )

	// Match an order with a car
	// result := submitTxnFn(
	// 	"org1", "autochannel", "KBA-Automobile",
	// 	"CarContract", "invoke", make(map[string][]byte),
	// 	"MatchOrder", "Car-06", "ORD-01",
	// )

	// Register a car
	// result := submitTxnFn(
	// 	"org3", "autochannel", "KBA-Automobile",
	// 	"CarContract", "invoke", make(map[string][]byte),
	// 	"RegisterCar", "Car-06", "Dani", "KL-01-CD-01",
	// )

	fmt.Println(result)
}
