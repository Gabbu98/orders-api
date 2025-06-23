package model

import (
	"time"

)

type Order struct {
	OrderID		string		`json:"order_id"`
	CustomerID	string		`json:"customer_id"`
	LineItems	[]LineItem	`json:"line_items"`
	CreatedAt	*time.Time	`json:"created_at"`
	ShippedAt 	*time.Time	`json:"shipped_at"`
	CompletedAt *time.Time	`json:"completed_at"`
}

type LineItem struct {
	ItemID 		string		`json:"item_id"`
	Quantity 	uint		`json:"quantity"`
	Price		uint		`json:"price"`
}