package main

import (
	"database/sql"
	"log"
	"net/http"
	"notesapp/internal/notes"
	"notesapp/internal/session"
	"notesapp/internal/user"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting application on port 8080")

	db, err := sql.Open("sqlite3", "notesapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement, err := os.ReadFile("internal/db/db.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlStatement))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created table succesfully or it already exists")

	store := notes.NewNotesStore(db)
	userStore := user.NewUserStore(db)
	sessionStore := session.NewSessionStore(db)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(sessionStore))
		r.Route("/notes", func(r chi.Router) {
			// /notes
			r.Get("/", notes.HandleGetNotes(store))
			r.Post("/", notes.HandleCreateNote(store))

			// /notes/{id}
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", notes.HandleGetNote(store))
				r.Put("/", notes.HandleUpdateNote(store))
				r.Delete("/", notes.HandleDeleteNote(store))
			})
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/", user.HandleCreateUser(userStore))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", user.HandleGetUserById(userStore))
		})
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", user.HandleLogin(userStore, sessionStore))
	})

	// log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
	log.Fatal(http.ListenAndServe(":8080", r))
}

// AuthMiddleware takes a Store and returns an http.HandlerFunc
// Don't use *session.SessionStore unless that's literally a struct type;
// if it's an interface, use session.SessionStore (no *).

func AuthMiddleware(s *session.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// better to use cookie since you set session_id as a cookie:
			c, err := r.Cookie("session_id")
			if err != nil || c.Value == "" {
				log.Printf("%v", err.Error())
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			sessionID := c.Value

			authenticated, err := s.CheckSession(ctx, sessionID)
			if err != nil || !authenticated {
				log.Printf("%v", err.Error())
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
