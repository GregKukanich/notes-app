package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting application on port 8080")

	db, err := sql.Open("sqlite3", "./testdb")
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
	http.HandleFunc("/notes", notesHandler(store))
	http.HandleFunc("/notes/", noteItemHandler(store))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
