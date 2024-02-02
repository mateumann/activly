package services

import (
	"fmt"

	"github.com/mateumann/activly/core/domain"
	"github.com/mateumann/activly/core/ports"
)

// UserService represents a service responsible for managing users.
type UserService struct {
	repo ports.UserRepository
}

// CreateUser creates a user with the given name and settings in the UserRepository.
// It returns an error if creating the user in the repository fails.
func (s *UserService) CreateUser(name string, settings map[string]interface{}) error {
	err := s.repo.Create(name, settings)
	if err != nil {
		return fmt.Errorf("failed creating a user in the repository: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users from the UserRepository.
// It returns a slice of domain.User objects and any error encountered during the process.
func (s *UserService) ListUsers() ([]*domain.User, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("failed listing users from the repository: %w", err)
	}

	return users, nil
}

func (s *UserService) Ready() bool {
	return s.repo != nil
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
