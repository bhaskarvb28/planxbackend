package middleware

import (
	"context"
	"net/http"

	"planx/internal/services"
)

const roleKey contextKey = "role"

func GetRole(ctx context.Context) string {
	if val, ok := ctx.Value(roleKey).(string); ok {
		return val
	}
	return ""
}

func WithRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())

		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		role, err := services.GetUserRole(userID)
		if err != nil {
			http.Error(w, "failed to fetch role", 500)
			return
		}

		ctx := context.WithValue(r.Context(), roleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}