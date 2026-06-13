package http

import (
	"net/http"
)

func NewRouter(healthHandler *HealthHandler, donateHandler *DonateHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.Check)
	mux.HandleFunc("GET /public/donate/{id}", donateHandler.GetByID)

	return mux
}
