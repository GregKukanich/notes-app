package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleGetNotes(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := store.getAll()
		if err != nil {
			http.Error(w, "failed to get notes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			http.Error(w, "failed to encode notes", http.StatusInternalServerError)
			return
		}
	}
}

func handleCreateNote(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			Body:  req.Body,
		}
		saved, err := store.save(note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(saved)

	}
}

func handleGetNote(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

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
}

func handleDeleteNote(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := store.delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleUpdateNote(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
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
			Body:  req.Body,
		}

		resp, err := store.update(id, note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode note", http.StatusInternalServerError)
			return
		}
	}
}
