package server

import (
	"awesomeProject1/internal/entity"
	"encoding/json"
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
	ID               string      `json:"id"`
	UserID           string      `json:"user_id"`
	ProductIDs       []string    `json:"product_ids"`
	CreatedAt        AwesomeTime `json:"created_at"`
	UpdatedAt        AwesomeTime `json:"updated_at"`
	DeliveryDeadLine AwesomeTime `json:"delivery_dead_line"`
	Price            float64     `json:"price"`
	DeliveryType     string      `json:"delivery_type"`
	Address          string      `json:"address"`
	OrderStatus      string      `json:"order_status"`
}

type UpdateOrderStatusRequest struct {
	OrderID     string             `json:"1"`
	OrderStatus entity.OrderStatus `json:"order_status"`
}

type GetOrdersRequest struct {
	UserID entity.GetOrders `json:"user_id"`
}

type EditOrderRequest struct {
	OrderID  string   `json:"1"`
	Products []string `json:"product_ids"`
	Address  string   `json:"address"`
}

type AwesomeTime time.Time

func (t *AwesomeTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t).Format(time.RFC3339))
}
