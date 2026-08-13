package table

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var table models.Table
	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		http.Error(w, "error decode body", http.StatusBadRequest)
		return
	}

	number, err := h.service.Create(r.Context(), table)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		Number int `json:"number"`
	}{
		Number: number,
	})
	if err != nil {
		http.Error(w, "error json encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tables, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tables)
	if err != nil {
		http.Error(w, "error json encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	numberStr := r.PathValue("number")

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "ivalid number", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), number)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "invalid number", http.StatusBadRequest)
			return
		}

		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
