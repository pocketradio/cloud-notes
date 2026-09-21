package routes

import (
	"cloud-notes/internal/endpoints"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", endpoints.HealthCheck)
	return mux
}