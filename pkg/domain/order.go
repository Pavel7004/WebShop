package domain

import (
	"time"
)

type StatusID string

var (
	CREATED   StatusID = "created"
	PAID      StatusID = "paid"
	DELIVERED StatusID = "delivered"
	CANCELED  StatusID = "canceled"
	REFUNDED  StatusID = "refunded"
)

type OrderItem struct {
	ID       string `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

type Order struct {
	ID         string      `json:"id"`
	Total      float64     `json:"total"`
	Items      []OrderItem `json:"items"`
	Status     StatusID    `json:"status"`
	CustomerID string      `json:"customer_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	Items      []OrderItem `json:"items"`
	CustomerID string      `json:"customer_id"`
}

type UpdateOrderRequest struct {
	Items  *[]OrderItem `json:"items"`
	Status *StatusID    `json:"status"`
}
