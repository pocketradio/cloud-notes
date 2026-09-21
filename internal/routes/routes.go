package routes

import (
	"cloud-notes/internal/endpoints"
	"net/http"
)

func New(notesList *endpoints.Notelist) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", endpoints.HealthCheck)
	mux.HandleFunc("POST /notes", notesList.CreateNote)
	mux.HandleFunc("GET /notes", notesList.GetNotes)
	mux.HandleFunc("GET /notes/{id}", notesList.GetOneNote)
	return mux
}
