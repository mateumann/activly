package ports

import (
	"github.com/mateumann/activly/core/domain"
)

// User Ports are interfaces that define how the communication between an actor, and the core has to be done.

// UserService represents a service for managing user data.
// It provides methods for creating, retrieving, and listing users.
//
// ListUsers method returns a slice of pointers to domain.User objects and an error.
//
// For more information on the fields and methods of domain.User, refer to its documentation.
type UserService interface {
	/*
		GetUserByName(name string) (*domain.User, error)
		GetUserByUUID(uuid uuid.UUID) (*domain.User, error)
	*/
	CreateUser(name string, settings map[string]interface{}) error

	ListUsers() ([]*domain.User, error)

	Ready() bool
}

// UserRepository represents a repository for managing user data.
// It provides methods for creating, retrieving, and listing users.
//
// Method ListUsers returns a slice of pointers to domain.User objects and an error.
//
// For more information on the fields and methods of domain.User, refer to its documentation.
type UserRepository interface {
	Create(name string, settings map[string]interface{}) error
	Save(user domain.User) error
	ListUsers() ([]*domain.User, error)
}
