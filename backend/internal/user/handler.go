package user

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

type IsRegisteredBody struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone_number"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), user)
	if err != nil {
		log.Fatal("Error while trying to register")
		return
	}

	fmt.Printf("Succesfully created user by id %d", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) IsRegisteredHandler(w http.ResponseWriter, r *http.Request) {

	var userinfo IsRegisteredBody

	if err := json.NewDecoder(r.Body).Decode(&userinfo); err != nil {
		http.Error(w, "Bad requst", http.StatusBadRequest)
		return
	}

	role, err := h.service.IsRegistered(r.Context(), userinfo.Nickname, userinfo.Phone)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)

}
