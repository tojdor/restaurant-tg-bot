package menuitem

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type menuItemRequest struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dishRequest menuItemRequest
	err := json.NewDecoder(r.Body).Decode(&dishRequest)
	if err != nil {
		http.Error(w, "error decode body", http.StatusBadRequest)
		return
	}

	dish := models.MenuItem{
		Name:     dishRequest.Name,
		Price:    dishRequest.Price,
		Category: dishRequest.Category,
	}
	id, err := h.service.Create(r.Context(), dish)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
	if err != nil {
		http.Error(w, "encode id error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category := r.PathValue("category")

	dishes, err := h.service.GetByCategory(category, r.Context())
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(dishes)
	if err != nil {
		http.Error(w, "error encode body", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.PathValue("name")

	dish, err := h.service.GetByName(name, r.Context())
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "invalid name", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(dish)
	if err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /dish", h.Create)
	mux.HandleFunc("GET /dish/{category}", h.GetByCategory)
	mux.HandleFunc("GET /dish/{name}", h.GetByName)
	mux.HandleFunc("DELETE /dish/{id}", h.Delete)
}
