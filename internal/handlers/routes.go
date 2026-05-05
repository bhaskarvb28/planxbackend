package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	myMw "planx/internal/middleware"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(myMw.Auth)

		RegisterProfileRoutes(r)
		RegisterCreditsRoutes(r)
		RegisterFloorPlanRoutes(r)
		RegisterVendorRoutes(r)
		RegisterEngineerRoutes(r)
	})
	return r
}