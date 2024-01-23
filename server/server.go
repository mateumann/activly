package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const RequestTimeout = 3

type Server struct {
	db *sql.DB
}

type user struct {
	id       uuid.UUID
	name     string
	settings []byte
	lastSeen *time.Time
}

func (u user) String() string {
	result := fmt.Sprintf("user %s (%s)", u.name, u.id)
	if u.lastSeen != nil {
		result += fmt.Sprintf(", last seen at %s", u.lastSeen)
	}

	if u.settings != nil {
		var settings map[string]any
		_ = json.Unmarshal(u.settings, &settings)
		result += fmt.Sprintf(", settings %v", settings)
	}

	return result
}

func (s Server) root(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.RequestIDKey)

	_, err := fmt.Fprintf(w, "hello, request_id=%s, db=%s", requestID, dbInfo(s.db))
	if err != nil {
		panic(err)
	}
}

func (s Server) auth0callback(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.RequestIDKey)

	_, err := fmt.Fprintf(w, "callback has been called, request_id=%s, db=%s", requestID, dbInfo(s.db))
	if err != nil {
		panic(err)
	}
}

func (s Server) listUsers(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.RequestIDKey)

	_, err := fmt.Fprintf(w, "listing users, request_id=%s\n", requestID)
	if err != nil {
		panic(err)
	}

	users, err := dbGetUsers(s.db)
	if err != nil {
		panic(err)
	}
	defer users.Close()

	for users.Next() {
		var u user
		if err = users.Scan(&u.id, &u.name, &u.settings, &u.lastSeen); err != nil {
			panic(err)
		}

		_, err := fmt.Fprintf(w, "  · %s\n", u)
		if err != nil {
			panic(err)
		}
	}
}

func Start() {
	db, err := dbConnect()
	if err != nil {
		panic(err)
	}

	defer dbClose(db)

	err = dbInit(db)
	if err != nil {
		panic(err)
	}

	handler := Server{db: db}

	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handler.root)
	r.Get("/users", handler.listUsers)
	r.Get("/callback", handler.auth0callback)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: RequestTimeout * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
