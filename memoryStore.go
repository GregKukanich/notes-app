package main

import (
	"fmt"
	"notesapp/internal/notes"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type inMemoryStore struct {
	mu    sync.Mutex
	notes []notes.Note
}

func (m *inMemoryStore) save(note notes.Note) (notes.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()
	note.ID = id
	m.notes = append(m.notes, note)

	return note, nil
}

func (m *inMemoryStore) delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, val := range m.notes {
		if val.ID == id {
			m.notes = slices.Delete(m.notes, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) update(id string, note notes.Note) (notes.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, val := range m.notes {
		if val.ID == id {
			originalNote := &m.notes[i]
			originalNote.Title = note.Title
			originalNote.Body = note.Body
			return m.notes[i], nil
		}
	}
	return notes.Note{}, fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) get(id string) (notes.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, note := range m.notes {
		if note.ID == id {
			return note, nil
		}
	}
	return notes.Note{}, fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) getAll() []notes.Note {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.notes
}
