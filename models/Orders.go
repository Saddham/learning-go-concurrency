package models

import (
	"fmt"
	"math"
	"sync"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatus_New       OrderStatus = "New"
	OrderStatus_Rejected  OrderStatus = "Rejected"
	OrderStatus_Completed OrderStatus = "Completed"
)

type Order struct {
	ID      string
	Product *Product
	Total   *float64
	Status  OrderStatus
}

type OrderDB struct {
	orders sync.Map
}

func NewOrder(product *Product) Order {
	total := (math.Round(float64(product.Quantity)*10.5*100) / 100)

	return Order{
		ID:      uuid.New().String(),
		Status:  OrderStatus_New,
		Product: product,
		Total:   &total,
	}
}

func (o *Order) Complete() {
	o.Status = OrderStatus_Completed
}

func (o *Order) Reject() {
	o.Status = OrderStatus_Rejected
}

func NewOrderDB() *OrderDB {
	return &OrderDB{}
}

func (odb *OrderDB) Find(id string) (Order, error) {
	order, ok := odb.orders.Load(id)
	if !ok {
		return Order{}, fmt.Errorf("No product found for product id %s", id)
	}

	return order.(Order), nil
}

func (odb *OrderDB) Upsert(order Order) {
	odb.orders.Store(order.ID, order)
}
