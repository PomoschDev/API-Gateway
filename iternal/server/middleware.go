package server

import (
	"context"
	"net/http"
)

// middleware - промежуточное ПО
func (route Router) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json, Authorization")

		ctx := context.WithValue(r.Context(), "admin", "test")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
