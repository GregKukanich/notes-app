package main

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
	getAll() ([]Note, error)
	update(id string, note Note) (Note, error)
}

