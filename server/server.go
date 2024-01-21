package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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
	w.Write([]byte(fmt.Sprintf("Hello, Chi! request_id=%v", r.Context().Value(middleware.RequestIDKey))))
}

func auth0callback(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Callback has been called with body=%v, form=%v, post_form=%v, cookies=%v", r.Body, r.Form, r.PostForm, r.Cookies())))
}
