package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mateumann/activly/core/domain"
)

var errBang = errors.New("bang")

type mockUserRepository struct {
	users []*domain.User
	err   error
}

func (m *mockUserRepository) Create(_ string, _ map[string]interface{}) error {
	return m.err
}

func (m *mockUserRepository) Save(_ domain.User) error {
	return m.err
}

func (m *mockUserRepository) ListUsers() ([]*domain.User, error) {
	return m.users, m.err
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

func TestUserService_ListUsers(t *testing.T) {
	t.Run("NilUsersAnError", func(t *testing.T) {
		repo := &mockUserRepository{nil, errBang}
		userService := NewUserService(repo)

		users, err := userService.ListUsers()
		if err == nil {
			t.Errorf("UserService.ListUsers() error = nil, wanted an %v", errBang)
		}

		if users != nil {
			t.Errorf("UserService.ListUsers() = %v, wanted nil", users)
		}
	})

	t.Run("SomeUsersAnError", func(t *testing.T) {
		repoUsers := []*domain.User{{Name: "Alice"}, {Name: "Bob"}}
		repo := &mockUserRepository{repoUsers, errBang}
		userService := NewUserService(repo)

		repoUsers, err := userService.ListUsers()
		if err == nil {
			t.Errorf("UserService.ListUsers() error = nil, wanted an %v", errBang)
		}

		if repoUsers != nil {
			t.Errorf("UserService.ListUsers() = %v, wanted nil", repoUsers)
		}
	})

	t.Run("SomeUsersNoError", func(t *testing.T) {
		repoUsers := []*domain.User{{Name: "Alice"}, {Name: "Bob"}}
		repo := &mockUserRepository{repoUsers, nil}
		userService := NewUserService(repo)

		users, err := userService.ListUsers()
		if err != nil {
			t.Errorf("UserService.ListUsers() error = %v, wanted nil", err)
		}

		if !reflect.DeepEqual(users, repoUsers) {
			t.Errorf("UserService.ListUsers() = %v, wanted %v", users, repoUsers)
		}
	})

	t.Run("NoUsersNoError", func(t *testing.T) {
		var repoUsers []*domain.User
		repo := &mockUserRepository{repoUsers, nil}
		userService := NewUserService(repo)

		users, err := userService.ListUsers()
		if err != nil {
			t.Errorf("UserService.ListUsers() error = %v, wanted nil", err)
		}

		if !reflect.DeepEqual(users, repoUsers) {
			t.Errorf("UserService.ListUsers() = %v, wanted %v", users, repoUsers)
		}
	})
}

func TestUserService_Ready(t *testing.T) {
	t.Run("AlwaysTrue", func(t *testing.T) {
		userService := NewUserService(nil)

		ready := userService.Ready()
		if ready != true {
			t.Errorf("UserService.Ready() = %v, want true", ready)
		}
	})
}
