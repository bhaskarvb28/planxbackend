package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"planx/internal/services"
	"planx/internal/middleware"
)

func RegisterFloorPlanRoutes(r chi.Router) {
	r.Post("/floor-plan/generate", GenerateFloorPlan)
}

type GenerateFloorPlanRequest struct {
	Prompt string `json:"prompt"`
}

func GenerateFloorPlan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req GenerateFloorPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}

	// 1. Check credits
	credits, err := services.GetCredits(userID)
	if err != nil {
		http.Error(w, "failed to fetch credits", http.StatusInternalServerError)
		return
	}

	if credits == nil || credits.Credits < 1 {
		http.Error(w, "not enough credits", http.StatusForbidden)
		return
	}

	// 2. Generate image
	imgBytes, err := services.GenerateFloorPlanImage(req.Prompt)
	if err != nil {
		http.Error(w, "failed to generate floor plan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Deduct credits AFTER success
	err = services.DeductCredits(userID, 1)
	if err != nil {
		http.Error(w, "failed to deduct credits", http.StatusInternalServerError)
		return
	}

	// 4. Return image
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=\"floorplan.png\"")
	w.WriteHeader(http.StatusOK)
	w.Write(imgBytes)
}