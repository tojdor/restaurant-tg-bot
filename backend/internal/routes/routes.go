package routes

import (
	menuitem "backend/internal/menu_item"
	orderitem "backend/internal/order_item"

	"backend/internal/table"
	"backend/internal/user"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(
	mux *http.ServeMux,
	pool *pgxpool.Pool,
) {
	// menu item
	menuItemStorage := menuitem.NewStorage(pool)
	menuItemService := menuitem.NewService(menuItemStorage)
	menuItemHandler := menuitem.NewHandler(menuItemService)

	mux.HandleFunc("POST /dish", menuItemHandler.Create)
	mux.HandleFunc("GET /dish/{category}", menuItemHandler.GetByCategory)
	mux.HandleFunc("GET /dish/{name}", menuItemHandler.GetByName)
	mux.HandleFunc("DELETE /dish/{id}", menuItemHandler.Delete)

	// order item
	orderItemStorage := orderitem.NewStorage(pool)
	orderItemService := orderitem.NewService(orderItemStorage)
	orderItemHandler := orderitem.NewHandler(orderItemService)

	mux.HandleFunc("POST /order/items", orderItemHandler.AddItem)
	mux.HandleFunc("GET /orders/{id}/items", orderItemHandler.GetByOrderID)
	mux.HandleFunc("PATCH /orders/{order_id}/items/{menu_item_id}/ready", orderItemHandler.SetReady)
	mux.HandleFunc("PATCH /orders/{order_id}/items/{menu_item_id}/count", orderItemHandler.UpdateCount)
	mux.HandleFunc("DELETE /orders/{order_id}/items/{menu_item_id}", orderItemHandler.DeleteItem)

	// table
	tableStorage := table.NewStorage(pool)
	tableService := table.NewService(tableStorage)
	tableHandler := table.NewHandler(tableService)

	mux.HandleFunc("POST /table", tableHandler.Create)
	mux.HandleFunc("GET /table", tableHandler.GetAll)
	mux.HandleFunc("DELETE /table/{number}", tableHandler.Delete)

	// user
	userStorage := user.NewStorage(pool)
	userService := user.NewService(userStorage)
	userHandler := user.NewHandler(userService)

	mux.HandleFunc("POST /user", userHandler.Register)
	mux.HandleFunc("GET /user", userHandler.GetAll)
	mux.HandleFunc("GET /user/login", userHandler.IsRegistered)
	mux.HandleFunc("DELETE /user/{id}", userHandler.Delete)
}
