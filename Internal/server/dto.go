package server

import (
	"awesomeProject1/Internal/entity"
	"time"
)

type CreateOrderRequest struct {
	UserID       string   `json:"user_id"`
	Products     []string `json:"products"`
	Price        float64  `json:"price"`
	DeliveryType string   `json:"delivery_type"`
	AddressID    string   `json:"address_id"`
}

type Order struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	ProductIDs       []string  `json:"product_ids"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeliveryDeadLine time.Time `json:"delivery_dead_line"`
	Price            float64   `json:"price"`
	DeliveryType     string    `json:"delivery_type"`
	Address          string    `json:"address"`
	OrderStatus      string    `json:"order_status"`
}

type UpdateOrderStatusRequest struct {
	OrderID     string             `json:"1"`
	OrderStatus entity.OrderStatus `json:"order_status"`
}

type GetOrdersRequest struct {
	UserID string `json:"user_id"`
}
