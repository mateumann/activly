package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mateumann/activly/core/services"
)

// HTTPHandler is a type that handles HTTP requests related to user functionality.
type HTTPHandler struct {
	userService *services.UserService
}

// NewHTTPHandler creates a new instance of the HTTPHandler struct using the provided UserService.
// It returns a pointer to the HTTPHandler struct.
func NewHTTPHandler(userService *services.UserService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
}

// ListUsers handles the HTTP request for listing users. It writes the list of users to the http.ResponseWriter
// in a formatted string.
// It uses the UserService's ListUsers method to retrieve the list of users from the repository.
// If there is an error while writing the response or retrieving the list of users, it returns
// an http.StatusInternalServerError.
// The requestID is obtained from the request's context.
//
// Usage Example:
//
//	func initRoutes(userService *services.UserService) *chi.Mux {
//	  r := chi.NewRouter()
//	  h := handler.NewHTTPHandler(userService)
//
//	  r.Get("/users", h.ListUsers)
//
//	  return r
//	}
//
// Method Signature:
//
//	func (h *HTTPHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
//	  ...
//	}
func (h *HTTPHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.RequestIDKey)

	_, err := fmt.Fprintf(w, "listing users, request_id=%s\n", requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	users, err := h.userService.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	for _, u := range users {
		_, err := fmt.Fprintf(w, "  · %s\n", u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}
