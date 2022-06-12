package domain

import "time"

type StatusID string

var (
	CREATED StatusID = "created"
	PAID    StatusID = "paid"
)

type Order struct {
	ID         string    `json:"id"`
	Total      float64   `json:"total"`
	ItemIDs    []string  `json:"item_ids"`
	CreatedAt  time.Time `json:"created_at"`
	Status     StatusID  `json:"status"`
	CustomerID string    `json:"customer_id"`
}

type CreateOrderRequest struct {
	ItemIDs    []string `json:"item_ids"`
	Total      float64  `json:"-"`
	CustomerID string   `json:"customer_id"`
}

type UpdateOrderRequest struct {
	ItemIDs *[]string `json:"item_ids"`
	Status  *StatusID `json:"status"`
}
