package domain

import "time"

type Item struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	Quantity    uint64    `json:"quantity"`
}

type AddItemRequest struct {
	OwnerID     string  `json:"owner_id"`
	Name        string  `json:"name"`
	Description string  `json:"desc"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Quantity    uint64  `json:"quantity"`
}
