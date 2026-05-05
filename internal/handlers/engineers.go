package handlers

import (
	"encoding/json"
	"net/http"
	"planx/internal/services"

	"github.com/go-chi/chi/v5"

	"fmt"

	"planx/internal/middleware"
)

func RegisterEngineerRoutes(r chi.Router) {
	r.Route("/engineers", func(r chi.Router) {
		r.Get("/", GetEngineers)

		r.With(middleware.WithRole).
			Post("/", CreateEngineer)
	})
}

func GetEngineers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := 10
	offset := 0

	if l := query.Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := query.Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	specialization := query.Get("specialization")
	city := query.Get("city")

	engineers, err := services.ListEngineers(limit, offset, specialization, city)
	if err != nil {
		http.Error(w, "failed to fetch engineers", 500)
		return
	}

	json.NewEncoder(w).Encode(engineers)
}

func CreateEngineer(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var body struct {
		Name           string `json:"name"`
		Phone          string `json:"phone"`
		Email          string `json:"email"`
		Specialization string `json:"specialization"`
		City           string `json:"city"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	if body.Name == "" || body.Phone == "" {
		http.Error(w, "missing required fields", 400)
		return
	}

	err := services.CreateEngineer(
		body.Name,
		body.Phone,
		body.Email,
		body.Specialization,
		body.City,
	)
	if err != nil {
		http.Error(w, "failed to create engineer", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "engineer created",
	})
}