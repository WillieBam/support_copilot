package middlewares

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Preserve the auth token in context for downstream verification/authorization layers.
		ctx := context.WithValue(r.Context(), "authToken", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
