package services

import (
	"testing"

	"github.com/mateumann/activly/core/domain"
)

type mockUserRepository struct {
	err error
}

func (m *mockUserRepository) Save(_ domain.User) error {
	return m.err
}

func (m *mockUserRepository) ListUsers() ([]*domain.User, error) {
	return nil, m.err
}

func TestNewUserService(t *testing.T) {
	t.Run("UserServiceWithValidRepo", func(t *testing.T) {
		repo := &mockUserRepository{}
		userService := NewUserService(repo)

		if userService.repo != repo {
			t.Errorf("UserService.repo = %v, want %v", userService.repo, repo)
		}
	})

	t.Run("UserServiceWithNilRepo", func(t *testing.T) {
		userService := NewUserService(nil)

		if userService.repo != nil {
			t.Errorf("UserService.repo = %v, wanted nil", userService.repo)
		}
	})
}
