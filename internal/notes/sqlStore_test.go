package notes

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

func setup(t *testing.T) *NotesStore {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = os.Remove("./test.db")
		if err != nil {
			log.Fatalf("Error deleting file: %v", err)
		}
		db.Close()
	})

	sqlStatement, err := os.ReadFile("db.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(sqlStatement))
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Created table succesfully or it already exists")

	// store := &inMemoryStore{}
	store := &NotesStore{db: db}
	return store
}

func TestSave(t *testing.T) {
	store := setup(t)
	testNote := Note{
		Title: "test title",
		Body:  "test body",
	}
	note, err := store.save(testNote)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	arr, err := store.getAll()
	if err != nil {
		t.Fatal("Failed to get notes from db")
	}

	if len(arr) != 1 {
		t.Fatal("More then 1 row in DB")
	}

	checkNote := arr[0]
	if checkNote != note {
		t.Fatalf("Note: %v did not match Note: %v", checkNote, note)
	}
}

func TestUpdate(t *testing.T) {
	store := setup(t)
	testNote := Note{
		Title: "test title",
		Body:  "test body",
	}
	note, err := store.save(testNote)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	note.Body = "updated body"
	_, err = store.update(note.ID, note)
	if err != nil {
		t.Fatalf("failed to update note: %v \n error: \n %v", note, err)
	}
	dbNote, err := store.get(note.ID)
	if err != nil {
		t.Fatalf("failed to get note with id: %s", note.ID)
	}
	if note.Body != dbNote.Body {
		t.Fatalf("Update call failed note: \n %v \n does not match note: \n %v", note, dbNote)
	}
}

func TestGet(t *testing.T) {
	store := setup(t)
	testNote := Note{
		Title: "test title",
		Body:  "test body",
	}
	note, err := store.save(testNote)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	dbNote, err := store.get(note.ID)
	if err != nil {
		t.Fatalf("Failed to get note from db with id: %s", note.ID)
	}

	if dbNote != note {
		t.Fatal("Note in DB does not match")
	}
}

func TestGetAll(t *testing.T) {
	store := setup(t)
	testNote := Note{
		Title: "test title",
		Body:  "test body",
	}
	//Loop through and create 5 notes
	for i := 0; i < 5; i++ {
		_, err := store.save(testNote)
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}
	}

	notesArray, err := store.getAll()
	if err != nil {
		t.Fatal("Failed to get all notes from DB")
	}

	if len(notesArray) != 5 {
		t.Fatalf("Incorrect number of notes in array.\n Expected : 5 \n Got: %d", len(notesArray))
	}
}
