package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthService struct {
	supabaseUrl     string
	supabaseAnonKey string
	httpClient      *http.Client
}

func NewAuthService(url, anonKey string) AuthServiceInterface {
	return &AuthService{
		supabaseUrl:     url,
		supabaseAnonKey: anonKey,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Login implements [AuthServiceInterface].
func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(payload)

	endpoint := fmt.Sprintf("%s/auth/v1/token?grant_type=password", s.supabaseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", s.supabaseAnonKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Invalid email or password")
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
