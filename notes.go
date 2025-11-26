package main

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type Note struct {
	ID string
	Title string
	Body string
}

type CreateNoteRequest struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

type Store interface {
	save(note Note) (Note, error)
	delete(id  string) error
	get(id string) (Note, error)
	getAll() []Note
	update(id string, note Note) (Note, error)
}

type inMemoryStore struct {
	notes []Note
}



func (m *inMemoryStore) save(note Note) (Note, error) {
	id := uuid.NewString()
	note.ID = id
	m.notes = append(m.notes, note);
	return note, nil
}

func (m *inMemoryStore) delete(id string) error {
	for i,val  := range m.notes {
		if val.ID == id {
			m.notes = slices.Delete(m.notes, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) update(id string, note Note) (Note, error) {
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
	for _,note := range m.notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, fmt.Errorf("note with id %s not found", id)
}

func (m *inMemoryStore) getAll() []Note {
	return m.notes
}