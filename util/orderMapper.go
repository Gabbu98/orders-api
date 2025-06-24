package util

import (
	handler_responses "github.com/Gabbu98/orders-api/handler/response"
	"github.com/Gabbu98/orders-api/model"
)

func MapOrderToOrderResponse(order model.Order) handler_responses.Order {
	return handler_responses.Order{
		OrderID: order.OrderID,
		CustomerID: order.CustomerID,
		LineItems: mapLineItemsToLineItemsResponse(order.LineItems),
		CreatedAt: order.CreatedAt,
		ShippedAt: order.ShippedAt,
		CompletedAt: order.CompletedAt,
	}
}

func mapLineItemsToLineItemsResponse(lineItems []model.LineItem) []handler_responses.LineItem {
	var lineItemsResponse []handler_responses.LineItem
	for j:=0; j<len(lineItems); j++ {
		lineItemsResponse = append(lineItemsResponse, handler_responses.LineItem{
			ItemID: lineItems[j].ItemID,
			Quantity: lineItems[j].Quantity,
			Price: lineItems[j].Price,
		})
	}
	return lineItemsResponse
}