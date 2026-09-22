package endpoints

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	db *sql.DB
}

func NewList(db *sql.DB) *Notelist {
	return &Notelist{db: db}
}

func (N *Notelist) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(note.Title) == "" || strings.TrimSpace(note.Content) == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	note.NoteID = uuid.New()
	_, err := N.db.ExecContext(r.Context(),
		"INSERT INTO notes (id, title, content) VALUES (?, ?, ?)",
		note.NoteID.String(), note.Title, note.Content,
	)
	if err != nil {
		http.Error(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(note)
}

func (N *Notelist) GetNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := N.db.QueryContext(r.Context(), "SELECT id, title, content FROM notes ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "failed to get notes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var note Note
		var id string
		if err := rows.Scan(&id, &note.Title, &note.Content); err != nil {
			http.Error(w, "failed to read note", http.StatusInternalServerError)
			return
		}
		note.NoteID, err = uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid note id in database", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "failed to get notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

func (N *Notelist) GetOneNote(w http.ResponseWriter, r *http.Request) {
	reqID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Please enter a valid ID", http.StatusBadRequest)
		return
	}

	var note Note
	var id string
	err = N.db.QueryRowContext(r.Context(),
		"SELECT id, title, content FROM notes WHERE id = ?", reqID.String(),
	).Scan(&id, &note.Title, &note.Content)
	if errors.Is(err, sql.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Note doesn't exist for the given ID"})
		return
	}
	if err != nil {
		http.Error(w, "failed to get note", http.StatusInternalServerError)
		return
	}

	note.NoteID, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid note id in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

func (N *Notelist) DeleteNote(w http.ResponseWriter, r *http.Request) {
	deleteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Please enter a valid ID", http.StatusBadRequest)
		return
	}

	result, err := N.db.ExecContext(r.Context(), "DELETE FROM notes WHERE id = ?", deleteID.String())
	if err != nil {
		http.Error(w, "failed to delete note", http.StatusInternalServerError)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to delete note", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Note doesn't exist for the given ID"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "note deleted successfully"})
}

func (N *Notelist) UpdateNote(w http.ResponseWriter, r *http.Request) {
	updateID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Please enter a valid ID", http.StatusBadRequest)
		return
	}

	var updatedNote Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, "Invalid json payload", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(updatedNote.Title) == "" || strings.TrimSpace(updatedNote.Content) == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	result, err := N.db.ExecContext(r.Context(),
		"UPDATE notes SET title = ?, content = ? WHERE id = ?",
		updatedNote.Title, updatedNote.Content, updateID.String(),
	)
	if err != nil {
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Note doesn't exist for the given ID"})
		return
	}

	updatedNote.NoteID = updateID
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedNote)
}
