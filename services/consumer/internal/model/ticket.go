package model

import "time"

type Ticket struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
