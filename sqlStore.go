package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type sqlStore struct {
	db *sql.DB
}

func (s *sqlStore) save(note Note) (Note, error) {
	id := uuid.NewString()
	note.ID = id

	stmt, err := s.db.Prepare("INSERT INTO notes (id, title, body) VALUES (?, ?, ?)")
	if err != nil {
		return Note{}, fmt.Errorf("save note: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(note.ID, note.Title, note.Body)
	if err != nil {
		return Note{}, fmt.Errorf("save note: %w", err)
	}

	return note, nil
}

func (s *sqlStore) delete(id string) error {

	stmt, err := s.db.Prepare("DELETE FROM notes WHERE id = ?;")
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("delete note %s: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete note %s: getting rows affected: %w", id, err)
	}

	if rows == 0 {
		return errors.New("note not found")
	}

	return nil
}

func (s *sqlStore) update(id string, note Note) (Note, error) {
	_, err := s.db.Exec("UPDATE notes SET title = ?, body = ? WHERE id = ?;", note.Title, note.Body, id)
	if err != nil {
		return Note{}, fmt.Errorf("update note: %s: %w", note.ID, err)
	}

	var n Note
	row := s.db.QueryRow("SELECT id, title, body FROM notes WHERE id = ?;", id)
	err2 := row.Scan(&n.ID, &n.Title, &n.Body)
	if errors.Is(err2, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err2 != nil {
		return Note{}, fmt.Errorf("get note %s: %w", id, err2)
	}
	return n, nil
}

func (s *sqlStore) get(id string) (Note, error) {
	var n Note

	row := s.db.QueryRow("SELECT id, title, body FROM notes WHERE id = ?;", id)

	err := row.Scan(&n.ID, &n.Title, &n.Body)
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, errors.New("note not found")
	}
	if err != nil {
		return Note{}, fmt.Errorf("get note %s: %w", id, err)
	}

	return n, nil
}

func (s *sqlStore) getAll() ([]Note, error) {
	rows, err := s.db.Query("SELECT id, title, body FROM notes;")
	if err != nil {
		return nil, fmt.Errorf("get all notes: %w", err)
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return notes, nil
}
