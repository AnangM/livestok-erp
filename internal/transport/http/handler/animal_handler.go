package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/AnangM/livestok-erp/internal/service"
	"github.com/AnangM/livestok-erp/internal/transport/http/middleware"
)

type AnimalHandler struct {
	service service.AnimalServiceInterface
}

func NewAnimalHandler(s service.AnimalServiceInterface) *AnimalHandler {
	return &AnimalHandler{service: s}
}

func (h *AnimalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var animal domain.Animal

	err := json.NewDecoder(r.Body).Decode(&animal)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get farm/user ID from the authenticated Supabase context
	farmID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || farmID == "" {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

	animal.FarmID = farmID

	err = h.service.RegisterNewAnimal(r.Context(), &animal)
	if err != nil {
		log.Printf("[AnimalHandler] Create failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}

func (h *AnimalHandler) List(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || farmID == "" {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

	animals, err := h.service.ListAnimals(r.Context(), farmID)
	if err != nil {
		log.Printf("[AnimalHandler] List failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}
