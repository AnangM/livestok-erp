package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnangM/livestok-erp/internal/service"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(s service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: s}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}
