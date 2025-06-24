package model

import (
	"time"

)

type Order struct {
	OrderID		string		`bson:"orderId"`
	CustomerID	string		`bson:"customerId"`
	LineItems	[]LineItem	`bson:"lineItems"`
	CreatedAt	*time.Time	`bson:"createdAt"`
	ShippedAt 	*time.Time	`bson:"shippedAt"`
	CompletedAt *time.Time	`bson:"completedAt"`
}

type LineItem struct {
	ItemID 		string		`bson:"itemId"`
	Quantity 	uint		`bson:"quantity"`
	Price		uint		`bson:"price"`
}