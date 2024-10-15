package server

import (
	"apiGateway/pkg/logger"
	"context"
	"net/http"
)

// authMiddleware - промежуточное ПО для приватных запросов
func (route Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json, Authorization")

		logger.Info("Приватный запрос: %s [%s]", r.URL.String(), r.Method)

		ctx := context.WithValue(r.Context(), "admin", "test")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// publicMiddleware - промежуточное ПО для публичных запросов
func (route Router) publicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		logger.Info("Публичный запрос: %s [%s]", r.URL.String(), r.Method)

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
