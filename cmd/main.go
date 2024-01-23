package main

import (
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/mateumann/activly/adapter/handler"
	"github.com/mateumann/activly/adapter/repository"
	"github.com/mateumann/activly/core/services"
)

const RequestTimeoutSeconds = 3

func main() {
	store := repository.NewUserPostgresRepository()
	service := services.NewUserService(store)

	r := initRoutes(service)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: RequestTimeoutSeconds * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initRoutes(svc *services.UserService) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewHTTPHandler(svc)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/users", h.ListUsers)

	return r
}
