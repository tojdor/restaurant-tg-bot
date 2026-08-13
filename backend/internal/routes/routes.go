package routes

import (
	menuitem "backend/internal/menu_item"
	"backend/internal/middleware"
	"backend/internal/order"
	orderitem "backend/internal/order_item"

	"backend/internal/table"
	"backend/internal/user"

	"net/http"
	"os"

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

	botSecret := os.Getenv("BOT_INTERNAL_SECRET")
	requireRoles := func(roles ...string) func(http.Handler) http.Handler {
		return middleware.RequireRoles(pool, botSecret, roles...)
	}
	mux.Handle("POST /dishes", requireRoles("admin")(http.HandlerFunc(menuItemHandler.Create)))
	mux.Handle("GET /dishes/category/{category}", requireRoles("admin", "waiter", "kitchen")(http.HandlerFunc(menuItemHandler.GetByCategory)))
	mux.Handle("GET /dishes/name/{name}", requireRoles("admin", "waiter", "kitchen")(http.HandlerFunc(menuItemHandler.GetByName)))
	mux.Handle("DELETE /dishes/{id}", requireRoles("admin")(http.HandlerFunc(menuItemHandler.Delete)))

	//order
	orderStorage := order.NewStorage(pool)
	orderService := order.NewService(orderStorage)
	orderHandler := order.NewOrderHandler(orderService)

	mux.Handle("POST /orders", requireRoles("waiter")(http.HandlerFunc(orderHandler.CreateOrder)))
	mux.Handle("GET /orders/waiter/{waiter_id}", requireRoles("admin", "waiter")(http.HandlerFunc(orderHandler.GetOrdersByWaiterID)))
	mux.Handle("GET /orders", requireRoles("admin", "kitchen")(http.HandlerFunc(orderHandler.GetAllOrders)))
	mux.Handle("GET /orders/table/{table_number}", requireRoles("admin", "waiter", "kitchen")(http.HandlerFunc(orderHandler.GetOrderByTableNumber)))
	mux.Handle("PATCH /orders/{id}", requireRoles("admin", "waiter")(http.HandlerFunc(orderHandler.UpdateOrder)))
	mux.Handle("DELETE /orders/{id}", requireRoles("admin", "waiter")(http.HandlerFunc(orderHandler.DeleteOrder)))

	// order item
	orderItemStorage := orderitem.NewStorage(pool)
	orderItemService := orderitem.NewService(orderItemStorage)
	orderItemHandler := orderitem.NewHandler(orderItemService)

	mux.Handle("POST /order-items", requireRoles("waiter")(http.HandlerFunc(orderItemHandler.AddItem)))
	mux.Handle("GET /order-items/{id}", requireRoles("admin", "waiter", "kitchen")(http.HandlerFunc(orderItemHandler.GetByOrderID)))
	mux.Handle("PATCH /order-items/{order_id}/{menu_item_id}/ready", requireRoles("kitchen")(http.HandlerFunc(orderItemHandler.SetReady)))
	mux.Handle("PATCH /order-items/{order_id}/{menu_item_id}/count", requireRoles("waiter")(http.HandlerFunc(orderItemHandler.UpdateCount)))
	mux.Handle("DELETE /order-items/{order_id}/{menu_item_id}", requireRoles("waiter")(http.HandlerFunc(orderItemHandler.DeleteItem)))

	// table
	tableStorage := table.NewStorage(pool)
	tableService := table.NewService(tableStorage)
	tableHandler := table.NewHandler(tableService)

	mux.Handle("POST /tables", requireRoles("admin")(http.HandlerFunc(tableHandler.Create)))
	mux.Handle("GET /tables", requireRoles("admin", "waiter", "kitchen")(http.HandlerFunc(tableHandler.GetAll)))
	mux.Handle("DELETE /tables/{number}", requireRoles("admin")(http.HandlerFunc(tableHandler.Delete)))

	// user
	userStorage := user.NewStorage(pool)
	userService := user.NewService(userStorage)
	userHandler := user.NewHandler(userService)

	mux.Handle("POST /users", requireRoles("admin")(http.HandlerFunc(userHandler.Register)))
	mux.Handle("GET /users", requireRoles("admin")(http.HandlerFunc(userHandler.GetAll)))
	mux.Handle("DELETE /users/{id}", requireRoles("admin")(http.HandlerFunc(userHandler.Delete)))
}
