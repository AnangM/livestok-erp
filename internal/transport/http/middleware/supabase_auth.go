package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const UserIdKey contextKey = "user_id"

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func SupabaseAuth(supabaseUrl, anonKey string) func(http.Handler) http.Handler {
	client := &http.Client{Timeout: 5 * time.Second}
	endpoint := fmt.Sprintf("%s/auth/v1/user", strings.TrimSuffix(supabaseUrl, "/"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing auth header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid auth header", http.StatusUnauthorized)
				return
			}

			req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint, nil)
			if err != nil {
				log.Printf("[Auth Middleware] Failed to create request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			req.Header.Set("Authorization", authHeader)
			req.Header.Set("apikey", anonKey)

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[Auth Middleware] Supabase request error: %v", err)
				http.Error(w, "Authentication service unavailable", http.StatusServiceUnavailable)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("[Auth Middleware] Supabase rejected token with status: %d", resp.StatusCode)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			var user SupabaseUser
			if err := json.NewDecoder(resp.Body).Decode(&user); err != nil || user.ID == "" {
				log.Printf("[Auth Middleware] Failed to parse user data: %v", err)
				http.Error(w, "Unable to retrieve user information", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), UserIdKey, user.ID))

			next.ServeHTTP(w, r)
		})
	}
}
