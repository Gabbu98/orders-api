package handler_responses

import (
	"time"

)

type Order struct {
	OrderID		string		`json:"orderId"`
	CustomerID	string		`json:"customerId"`
	LineItems	[]LineItem	`json:"lineItems"`
	CreatedAt	*time.Time	`json:"createdAt"`
	ShippedAt 	*time.Time	`json:"shippedAt"`
	CompletedAt *time.Time	`json:"completedAt"`
}

type LineItem struct {
	ItemID 		string		`json:"itemId"`
	Quantity 	uint		`json:"quantity"`
	Price		uint		`json:"price"`
}