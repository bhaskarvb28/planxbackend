package handlers

import (
	"encoding/json"
	"net/http"

	"planx/internal/middleware"
	"planx/internal/services"

	"github.com/go-chi/chi/v5"
)

type CreateUserProfileRequest struct {
	Name string `json:"name"`
}

func RegisterProfileRoutes(r chi.Router) {
	r.Get("/profile/me", GetMyProfile)
	r.Post("/profile", CreateProfile)
}

func GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	profile, err := services.GetMyProfile(userID)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"profile": profile,
	})
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req CreateUserProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		http.Error(w, "Invalid request", 400)
		return
	}

	profile, status, err := services.CreateProfile(userID, req.Name)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"profile": profile,
	})
}