package main

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type inMemoryStore struct {
	mu sync.Mutex
	notes []Note
}



func (m *inMemoryStore) save(note Note) (Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()
	note.ID = id
	m.notes = append(m.notes, note);


	return note, nil
}

func (m *inMemoryStore) delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i,val  := range m.notes {
		if val.ID == id {
			m.notes = slices.Delete(m.notes, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) update(id string, note Note) (Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i,val := range m.notes{
		if val.ID == id {
			originalNote := &m.notes[i]
			originalNote.Title = note.Title
			originalNote.Body = note.Body
			return m.notes[i], nil
		}
	}
	return Note{}, fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) get(id string) (Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _,note := range m.notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) getAll() []Note {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.notes
}