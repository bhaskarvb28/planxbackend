package handlers

import (
	"planx/internal/middleware"
	"planx/internal/services"
	"net/http"
	"encoding/json"

	"github.com/go-chi/chi/v5"
)

func RegisterVendorRoutes(r chi.Router) {
	r.Route("/vendor", func(r chi.Router) {

		// Public (or just auth)
		r.Post("/apply", ApplyVendor)

		// Routes that need role
		r.With(middleware.WithRole).Group(func(r chi.Router) {
			r.Get("/applications", GetVendorApplications)
			r.Patch("/applications/{id}/approve", ApproveVendor)
			r.Patch("/applications/{id}/reject", RejectVendor)
		})
	})
}

func ApplyVendor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		ShopName string `json:"shop_name"`
		GSTIN    string `json:"gstin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	if body.ShopName == "" || body.GSTIN == "" {
		http.Error(w, "missing fields", 400)
		return
	}

	err := services.CreateVendorApplication(userID, body.ShopName, body.GSTIN)
	if err != nil {
		http.Error(w, "failed to apply", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "application submitted",
	})
}

func GetVendorApplications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	role := middleware.GetRole(r.Context())

	status := r.URL.Query().Get("status")

	var apps interface{}
	var err error

	if role == "admin" {
		// admin → see all
		apps, err = services.ListVendorApplications(status)
	} else {
		// user → see only their own
		apps, err = services.ListUserVendorApplications(userID, status)
	}

	if err != nil {
		http.Error(w, "failed to fetch applications", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func ApproveVendor(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	appID := chi.URLParam(r, "id")
	if appID == "" {
		http.Error(w, "missing application id", 400)
		return
	}

	err := services.ApproveVendorApplication(appID)
	if err != nil {
		http.Error(w, "failed to approve", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "vendor approved",
	})
}

func RejectVendor(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	appID := chi.URLParam(r, "id")
	if appID == "" {
		http.Error(w, "missing application id", 400)
		return
	}

	err := services.RejectVendorApplication(appID)
	if err != nil {
		http.Error(w, "failed to reject", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "application rejected",
	})
}