package server

import (
	"fmt"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const RequestTimeout = 3

func Serve() {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", root)
	r.Get("/callback", auth0callback)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: RequestTimeout * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello, request_id=%s", r.Context().Value(middleware.RequestIDKey))
	if err != nil {
		panic(err)
	}
}

func auth0callback(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.RequestIDKey)
	_, err := fmt.Fprintf(w, "callback has been called, request_id=%v", requestID)
	if err != nil {
		panic(err)
	}
}
