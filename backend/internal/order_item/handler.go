package orderitem

import (
	"backend/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type updateCountRequest struct {
	Count int `json:"count"`
}

type readyRequest struct {
	IsReady bool `json:"is_ready"`
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item models.OrderItem

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.AddItem(r.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetByOrderID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetByOrderID(
		r.Context(),
		orderID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SetReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID, err := strconv.Atoi(r.PathValue("order_id"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	menuItemID, err := strconv.Atoi(r.PathValue("menu_item_id"))
	if err != nil {
		http.Error(w, "invalid menu item id", http.StatusBadRequest)
		return
	}

	var req readyRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.SetReady(
		r.Context(),
		orderID,
		menuItemID,
		req.IsReady,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID, err := strconv.Atoi(r.PathValue("order_id"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	menuItemID, err := strconv.Atoi(r.PathValue("menu_item_id"))
	if err != nil {
		http.Error(w, "invalid menu item id", http.StatusBadRequest)
		return
	}

	var req updateCountRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateCount(
		r.Context(),
		orderID,
		menuItemID,
		req.Count,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID, err := strconv.Atoi(r.PathValue("order_id"))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	menuItemID, err := strconv.Atoi(r.PathValue("menu_item_id"))
	if err != nil {
		http.Error(w, "invalid menu item id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteItem(
		r.Context(),
		orderID,
		menuItemID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /order/items", h.AddItem)
	mux.HandleFunc("GET /orders/{id}/items", h.GetByOrderID)
	mux.HandleFunc("PATCH /orders/{order_id}/items/{menu_item_id}/ready", h.SetReady)
	mux.HandleFunc("PATCH /orders/{order_id}/items/{menu_item_id}/count", h.UpdateCount)
	mux.HandleFunc("DELETE /orders/{order_id}/items/{menu_item_id}", h.DeleteItem)
}
