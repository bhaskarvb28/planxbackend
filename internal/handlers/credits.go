package handlers

import (
	"encoding/json"
	"net/http"

	"planx/internal/middleware"
	"planx/internal/services"

	"github.com/go-chi/chi/v5"
)

func RegisterCreditsRoutes(r chi.Router) {
	r.Get("/credits", GetCredits)
	r.Post("/credits/add", AddCreditsHandler) 
}

func GetCredits(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	credits, err := services.GetCredits(userID)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	if credits == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"credits": 0,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"credits": credits.Credits,
	})
}

func AddCreditsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if body.Amount <= 0 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	err := services.AddCredits(userID, body.Amount)
	if err != nil {
		http.Error(w, "failed to add credits", http.StatusInternalServerError)
		return
	}

	// 👇 Fetch updated credits
	updated, err := services.GetCredits(userID)
	if err != nil {
		http.Error(w, "failed to fetch updated credits", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "credits added successfully",
		"credits": updated.Credits,
	})
}