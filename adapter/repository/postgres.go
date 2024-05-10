package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/mateumann/activly/core/domain"
)

const dbInitNameLength = 120

type UserPostgresRepository struct {
	db *sql.DB
}

func (r *UserPostgresRepository) init() error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin a db transaction: %w", err)
	}

	_, err = tx.Exec(
		"CREATE TABLE IF NOT EXISTS public.users " + "(" +
			"id uuid PRIMARY KEY DEFAULT gen_random_uuid(), " +
			fmt.Sprintf("name varchar(%d) NOT NULL UNIQUE, ", dbInitNameLength) +
			"settings jsonb, " +
			"last_seen timestamp without time zone DEFAULT now()" +
			");",
	)
	if err != nil {
		_ = tx.Rollback()

		return fmt.Errorf("failed to create users table: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit a db transaction: %w", err)
	}

	return nil
}

// Close closes the database connection used by the UserPostgresRepository.
// If an error occurs while closing the connection, it panics to indicate the failure.
// It is recommended to call this method when the repository is no longer needed to properly
// release the resources associated with the connection.
func (r *UserPostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		panic(err)
	}
}

var errPostgresRepositorySaveNotImplemented = errors.New("UserPostgresRepository.Save not yet implemented")

func (r *UserPostgresRepository) Create(name string, settings map[string]interface{}) error {
	stmt, err := r.db.Prepare("INSERT INTO public.users (name, settings) VALUES ($1, $2);")
	if err != nil {
		return fmt.Errorf("failed to prepare an insert user statement: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	_, err = stmt.Exec(name, settings)
	if err != nil {
		return fmt.Errorf("failed to insert a user: %w", err)
	}

	return nil
}

func (*UserPostgresRepository) Save(_ domain.User) error {
	return errPostgresRepositorySaveNotImplemented
}

// ListUsers retrieves a list of users from the database.
// It executes a SQL query to select the user ID, name, settings, and last seen timestamp
// from the "users" table in ascending order by last seen timestamp.
// If the query execution fails, it returns an error.
// It then iterates over the rows of the result set and scans each row into a User struct.
// If scanning fails for any row, it returns an error.
// Otherwise, it appends each user to a slice of User pointers.
// If there is an error during iteration, such as a broken connection, it returns an error.
// Finally, it returns the slice of users and nil to indicate success.
//
// Example usage:
//
//	users, err := repo.ListUsers()
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, user := range users {
//		fmt.Println(user)
//	}
//
// Declaration:
//
//	func (r *UserPostgresRepository) ListUsers() ([]*domain.User, error) {
//		...
//	}
func (r *UserPostgresRepository) ListUsers() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, settings, last_seen FROM public.users ORDER BY last_seen")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	users := make([]*domain.User, 0)

	for rows.Next() {
		var u domain.User

		var settings []byte

		err := rows.Scan(&u.ID, &u.Name, &settings, &u.LastSeen)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		if len(settings) > 0 {
			err = json.Unmarshal(settings, &u.Settings)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal user settings: %w", err)
			}
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over user rows: %w", err)
	}

	return users, nil
}

// NewUserPostgresRepository returns a new instance of UserPostgresRepository.
// It creates a connection string using environment variables for PostgreSQL connection details: POSTGRES_HOST,
// POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER and POSTGRES_PASSWORD.  Then, it opens a connection to the PostgreSQL
// database.  If the connection fails, it panics with the error.  Finally, it initialises and returns
// a UserPostgresRepository with the opened database connection.
func NewUserPostgresRepository() *UserPostgresRepository {
	conn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	repo := &UserPostgresRepository{db: db}
	if err := repo.init(); err != nil {
		panic(err)
	}

	return repo
}
