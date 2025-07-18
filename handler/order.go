package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Gabbu98/orders-api/handler/response"
	"github.com/Gabbu98/orders-api/model"
	"github.com/Gabbu98/orders-api/repository/order"
	"github.com/Gabbu98/orders-api/util"
)


type Order struct {
	Repo	order.OrderRepository
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
	var body struct {
		CustomerID	string				`json:"customer_id"`
		LineItems	[]model.LineItem	`json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID: 	strconv.FormatUint(rand.Uint64(), 10),
		CustomerID: body.CustomerID,
		LineItems: 	body.LineItems,
		CreatedAt: &now,
	}

	if err := o.Repo.Insert(r.Context(), order); err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
	
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size: size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items	[]handler_responses.Order	`json:"items"`
		Next	uint64				`json:"next,omitempty"` // omits the value in the response if it isnt present
	}
	var orders []handler_responses.Order
	for i:=0; i<len(res.Orders); i++ {
		orders = append(orders, util.MapOrderToOrderResponse(res.Orders[i]))
	}

	response.Items = orders
	response.Next = res.Cursor

	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an order by ID")
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderFetched, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNoExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(orderFetched); err != nil {
		fmt.Println("failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")

	//decode param id
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//unmarshal request
	var body struct {
		Status	string	`json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fetch order
	orderFetched, err := o.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNoExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find him by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//update
	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
		case completedStatus: 
		{
			if orderFetched.CompletedAt != nil || orderFetched.ShippedAt == nil{
				w.WriteHeader(http.StatusConflict)
				return
			}
			orderFetched.CompletedAt = &now
		}
		case shippedStatus: 
		{
			if orderFetched.ShippedAt != nil {
				w.WriteHeader(http.StatusConflict)
				return
			}
			orderFetched.ShippedAt = &now
		}
		default: 
		{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if err := o.Repo.Update(r.Context(), orderFetched); err != nil {
		fmt.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(orderFetched)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//return updated response
	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")
	idParam := chi .URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	
	err = o.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNoExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}