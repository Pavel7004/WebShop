package domain

import "time"

type Item struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type AddItemRequest struct {
	Name        string  `json:"name"`
	OwnerID     string  `json:"owner_id"`
	Description string  `json:"desc"`
	Price       float64 `json:"price"`
}
