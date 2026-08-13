package user

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), user)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "invalid user data", http.StatusBadRequest)
			return
		}
		log.Printf("error while trying to register: %v", err)
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	log.Printf("successfully created user with id %d", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "getting id error", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "getting id error", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
