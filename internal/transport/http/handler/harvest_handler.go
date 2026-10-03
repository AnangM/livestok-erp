package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/AnangM/livestok-erp/internal/service"
	"github.com/AnangM/livestok-erp/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

type HarvestHandler struct {
	service service.HarvestServiceInterface
}

func NewHarvestHandler(s service.HarvestServiceInterface) *HarvestHandler {
	return &HarvestHandler{service: s}
}

func (h *HarvestHandler) Create(w http.ResponseWriter, r *http.Request) {

	var harvest domain.Harvest

	err := json.NewDecoder(r.Body).Decode(&harvest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	farmID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || farmID == "" {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

	harvest.FarmID = farmID

	err = h.service.RegisterNewHarvest(r.Context(), &harvest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(harvest)
}

func (h *HarvestHandler) List(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || farmID == "" {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

	harvests, err := h.service.ListHarvests(r.Context(), farmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(harvests)
}

func (h *HarvestHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Harvest ID is required", http.StatusBadRequest)
		return
	}

	harvest, err := h.service.GetHarvest(r.Context(), id)
	if err != nil {
		http.Error(w, "Harvest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(harvest)
}

func (h *HarvestHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Harvest ID is required", http.StatusBadRequest)
		return
	}

	var harvest domain.Harvest
	err := json.NewDecoder(r.Body).Decode(&harvest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	farmID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok || farmID == "" {
		http.Error(w, "Unauthorized: missing user identity", http.StatusUnauthorized)
		return
	}

	harvest.FarmID = farmID

	updatedHarvest, err := h.service.UpdateHarvest(r.Context(), id, &harvest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedHarvest)
}

func (h *HarvestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Harvest ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteHarvest(r.Context(), id)
	if err != nil {
		http.Error(w, "Harvest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
