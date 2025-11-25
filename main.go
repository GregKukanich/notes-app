package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)


func notesHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNotes(w,r, store)
		case http.MethodPost:
			handleCreateNote(w,r, store)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func noteItemHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
		
		if idStr == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleGetNote(w,r, idStr, store)
		case http.MethodDelete:
			handleDeleteNote(w,r, idStr, store)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetNotes(w http.ResponseWriter, r *http.Request, store Store){
	notes := store.getAll()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "failed to encode notes", http.StatusInternalServerError)
		return
	}
}
func handleCreateNote(w http.ResponseWriter, r *http.Request, store Store){
	var req CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Body == "" {
		http.Error(w, "title/body is required", http.StatusBadRequest)
		return
	}

	note := Note{
		Title: req.Title,
		Body: req.Body,
	}
	saved,err := store.save(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(saved)

}
func handleGetNote(w http.ResponseWriter, r *http.Request, id string, store Store){
	note, err := store.get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "failed to encode note", http.StatusInternalServerError)
		return
	}
}
func handleDeleteNote(w http.ResponseWriter, r *http.Request, id string, store Store){
	err := store.delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.Output(1, "Starting application")
	store := &inMemoryStore{}
	http.HandleFunc("/notes", notesHandler(store))
	http.HandleFunc("/notes/", noteItemHandler(store))
	log.Fatal(http.ListenAndServe(":8080", nil))
}