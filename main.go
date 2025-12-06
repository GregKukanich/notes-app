package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting application on port 8080")

	db, err := sql.Open("sqlite3", "./notesapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement, err := os.ReadFile("db.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlStatement))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created table succesfully or it already exists")

	// store := &inMemoryStore{}
	store := &sqlStore{db: db}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/notes", func(r chi.Router) {
		// /notes
		r.Get("/", handleGetNotes(store))
		r.Post("/", handleCreateNote(store))

		// /notes/{id}
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handleGetNote(store))
			r.Put("/", handleUpdateNote(store))
			r.Delete("/", handleDeleteNote(store))
		})
	})

	// log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
	log.Fatal(http.ListenAndServe(":8080", r))
}
