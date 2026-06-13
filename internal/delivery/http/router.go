package http

import (
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(healthHandler *HealthHandler, donateHandler *DonateHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.Check)
	mux.HandleFunc("GET /public/donate/{id}", donateHandler.GetByID)

	return enableCORS(mux)
}
