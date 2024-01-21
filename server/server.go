package server

import (
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Serve() {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", root)
	r.Get("/callback", auth0callback)

	_ = http.ListenAndServe(":8080", r)
}

func root(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("hello, request_id=%v", r.Context().Value(middleware.RequestIDKey))))
	if err != nil {
		panic(err)
	}
}

func auth0callback(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(fmt.Sprintf("callback has been called, request_id=%v",
		r.Context().Value(middleware.RequestIDKey))))
	if err != nil {
		panic(err)
	}
}
