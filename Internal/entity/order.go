package entity

import (
	"time"
)

type DType string

const (
	Courier DType = "courier"
	Drone   DType = "Drone"
	Myself  DType = "Myself"
)

type OrderStatus string

const (
	Created   OrderStatus = "created"
	Paid      OrderStatus = "paid"
	Collect   OrderStatus = "collect"
	Collected OrderStatus = "collected"
	Delivery  OrderStatus = "delivery"
	Done      OrderStatus = "done"
	Cancelled OrderStatus = "cancelled"
)

type Order struct {
	ID               string
	UserID           string
	ProductIDs       []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeliveryDeadLine time.Time
	Price            float64
	DeliveryType     DType
	Address          string
	OrderStatus      OrderStatus
}

type CreateOrderRequest struct {
	UserID       string
	Products     []string
	Price        float64
	DeliveryType DType
	AddressID    string
}

type GetOrders struct {
	UserID string
	Limit  uint
	Page   uint
	Asc    bool
}
