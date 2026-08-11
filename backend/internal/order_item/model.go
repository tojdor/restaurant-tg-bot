package models

type OrderItem struct {
	OrderID   int  `json:"order_id"`
	MenuIemID int  `json:"menu_item_id"`
	Count     int  `json:"count"`
	IsReady   bool `json:"is_ready"`
}
