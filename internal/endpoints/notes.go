package endpoints

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Note struct {
	NoteID  uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type Notelist struct {
	notes []Note
}

func (N *Notelist) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(note.Title) == "" || strings.TrimSpace(note.Content) == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}
	note.NoteID = uuid.New()
	N.notes = append(N.notes, note)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func NewList() *Notelist {
	return &Notelist{
		notes: []Note{},
	}
}

func (N *Notelist) GetNotes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(N.notes); err != nil {
		http.Error(w, "failed to get notes", 500)
	}

}

func (N *Notelist) GetOneNote(w http.ResponseWriter, r *http.Request) {

	reqID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Please enter a valid ID", 400)
		return
	}

	flag := false

	var requestedNote Note
	for _, note := range N.notes {

		if note.NoteID == reqID {
			requestedNote = note
			flag = true
			break
		}
	}

	if !flag {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Note doesn't exist for the given ID",
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(requestedNote)
	}

}
