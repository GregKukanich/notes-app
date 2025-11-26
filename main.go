package main

import (
	"log"
	"net/http"
)

func main() {
	log.Output(1, "Starting application on port 8080")
	store := &inMemoryStore{}
	http.HandleFunc("/notes", notesHandler(store))
	http.HandleFunc("/notes/", noteItemHandler(store))
	log.Fatal(http.ListenAndServe(":8080", nil))
}