package models

import "time"

type MenuIem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

type Order struct {
	ID          int        `json:"id"`
	TableNumber int        `json:"table_number"`
	WaiterID    int        `json:"waiter_id"`
	IsServed    bool       `json:"is_served"`
	IsPayed     bool       `json:"is_payed"`
	CreatedAt   time.Time  `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}

type OrderItem struct {
	OrderID   int  `json:"order_id"`
	MenuIemID int  `json:"menu_item_id"`
	Count     int  `json:"count"`
	IsReady   bool `json:"is_ready"`
}

type Table struct {
	Number int    `json:"number"`
	Status string `json:"string"`
}

type User struct {
	ID            int    `json:"id"`
	Nickname      string `json:"nickname"`
	PhoneNumber string `json:"phone_number"`
	Role          string `json:"role"`
}
