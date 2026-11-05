package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ReynoldArun09/school-procurement-management/internal/services"
	"github.com/ReynoldArun09/school-procurement-management/internal/utils"
)

const sessionCookieName = "session_token"

type AdminAuthHandler struct {
	service *services.AdminAuthService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAdminAuthHandler(service *services.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{
		service: service,
	}
}

func extractToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if errors.Is(err, services.ErrInactiveAdmin) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    result.SessionID,
		Path:     "/",
		Expires:  result.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(utils.ApiResponse{
		Message:    "Login successful",
		StatusCode: http.StatusOK,
		Data: map[string]any{
			"admin":        result.Admin,
			"sessionToken": result.SessionID,
			"expiresAt":    result.ExpiresAt,
		},
	})
}

func (h *AdminAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)

	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	admin, err := h.service.Authenticate(r.Context(), token)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrInvalidSession) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(utils.ApiResponse{
		Message:    "Admin profile fetched successfully",
		StatusCode: http.StatusOK,
		Data:       admin,
	})
}

func (h *AdminAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)

	if token != "" {
		if err := h.service.Logout(r.Context(), token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(utils.ApiResponse{
		Message:    "Logged out successfully",
		StatusCode: http.StatusOK,
	})
}
