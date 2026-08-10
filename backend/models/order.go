package models

import "time"

type Order struct {
	ID          int        `json:"id"`
	TableNumber int        `json:"table_number"`
	WaiterID    int        `json:"waiter_id"`
	IsServed    bool       `json:"is_served"`
	IsPayed     bool       `json:"is_payed"`
	CreatedAt   time.Time  `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}
