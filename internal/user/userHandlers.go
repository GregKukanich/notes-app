package user

import (
	"encoding/json"
	"log"
	"net/http"
	"notesapp/internal/session"
	"time"

	"github.com/go-chi/chi/v5"
)

func HandleCreateUser(store UserStoreIntf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req UserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.UserName == "" || req.Password == "" {
			http.Error(w, "username/password is required", http.StatusBadRequest)
			return
		}

		hash, err := HashPassword(req.Password)
		if err != nil {
			log.Printf("Failed hashing password: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		user := User{
			UserName: req.UserName,
			Password: hash,
		}
		saved, err := store.CreateUser(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(saved)
	}
}

func HandleGetUserById(store UserStoreIntf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		ur, err := store.GetByID(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ur); err != nil {
			http.Error(w, "failed to encode note", http.StatusInternalServerError)
			return
		}
	}
}

func HandleLogin(store UserStoreIntf, sessionStore session.SessionStoreIntf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req UserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.UserName == "" || req.Password == "" {
			http.Error(w, "username/password is required", http.StatusBadRequest)
			return
		}

		hash, err := store.GetUserPasswordHash(ctx, req.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		err = VerifyPassword(hash, req.Password)
		if err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		user, err := store.GetUserByUserName(ctx, req.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Success
		session, err := sessionStore.CreateSession(ctx, user.ID)
		if err != nil {
			log.Printf("Failed to create session for user: %s", user.ID)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set cookie with session ID
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,                 // true in real HTTPS setups
			SameSite: http.SameSiteLaxMode, // or Strict/None depending on use
			Expires:  time.Unix(session.Expiration, 0),
		})

		// Optionally write a JSON body or just a 204/200
		w.WriteHeader(http.StatusOK)

	}
}
