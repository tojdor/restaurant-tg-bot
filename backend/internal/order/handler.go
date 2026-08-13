package order

import (
	"backend/internal/middleware"
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service *Service
}

func NewOrderHandler(service *Service) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.Order

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "authenticated user is missing", http.StatusInternalServerError)
		return
	}
	req.WaiterID = principal.UserID

	id, err := h.service.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "Invalid order data: waiter_id and table_number are required", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderByTableNumber(w http.ResponseWriter, r *http.Request) {

	numberStr := r.PathValue("table_number")
	number, err := strconv.Atoi(numberStr)
	if err != nil || number <= 0 {
		http.Error(w, "Invalid table number", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByTableNumber(r.Context(), number)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "Invalid table number", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "Order not found for this table", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	if principal, ok := middleware.PrincipalFromContext(r.Context()); ok && principal.Role == "waiter" && principal.UserID != order.WaiterID {
		http.Error(w, "waiters can view only their own orders", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrdersByWaiterID(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("waiter_id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid waiter ID", http.StatusBadRequest)
		return
	}
	if principal, ok := middleware.PrincipalFromContext(r.Context()); ok && principal.Role == "waiter" && principal.UserID != id {
		http.Error(w, "waiters can view only their own orders", http.StatusForbidden)
		return
	}

	orders, err := h.service.GetByWaiterId(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "Invalid waiter ID", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req OrderRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.ID = id
	if principal, ok := middleware.PrincipalFromContext(r.Context()); ok && principal.Role == "waiter" {
		req.WaiterID = principal.UserID
	}
	if !h.canManageOrder(w, r, id) {
		return
	}

	if err := h.service.Update(r.Context(), req); err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "Invalid update data", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order updated successfully"})
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	if !h.canManageOrder(w, r, id) {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) canManageOrder(w http.ResponseWriter, r *http.Request, orderID int) bool {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok || principal.Role != "waiter" {
		return true
	}

	owned, err := h.service.IsOwnedByWaiter(r.Context(), orderID, principal.UserID)
	if err != nil {
		http.Error(w, "failed to authorize order access", http.StatusInternalServerError)
		return false
	}
	if !owned {
		http.Error(w, "waiters can manage only their own orders", http.StatusForbidden)
		return false
	}
	return true
}
