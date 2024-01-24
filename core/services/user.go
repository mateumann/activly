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

// NewUserService creates a new instance of UserService with the provided UserRepository.
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (*UserService) Ready() bool {
	return true // come up with something more useful, e.g. is database ready or something like that
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
