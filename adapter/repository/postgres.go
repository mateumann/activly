package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mateumann/activly/core/domain"
)

const dbInitNameLength = 120

type UserPostgresRepository struct {
	db *sql.DB
}

// NewUserPostgresRepository returns a new instance of UserPostgresRepository.
// It creates a connection string using environment variables for PostgreSQL connection details.
// Then, it opens a connection to the PostgreSQL database.
// If the connection fails, it panics with the error.
// Finally, it initialises and returns a UserPostgresRepository with the opened database connection.
func NewUserPostgresRepository() *UserPostgresRepository {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	return &UserPostgresRepository{db: db}
}

// Init initialises the UserPostgresRepository by creating the necessary database table.
// It starts a transaction using the database connection. If the transaction fails to
// begin, it returns an error indicating the failure. Otherwise, it executes a SQL
// statement to create the "users" table with the required columns. If the execution
// of SQL fails, it rolls back the transaction and returns an error. After successfully
// creating the table, it commits the transaction. If the commit fails, it returns an
// error. Otherwise, it returns nil to indicate successful initialization.
// Note: The length of the "name" column in the table is set based on the `dbInitNameLength`
// constant value, which is defined as 120.
//
// Example usage:
//
//	err := repo.Init()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Declaration:
//
//	func (r *UserPostgresRepository) Init() error {
//	    ...
//	}
func (r *UserPostgresRepository) Init() error {
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
	defer rows.Close()

	users := make([]*domain.User, 0)

	for rows.Next() {
		var u domain.User

		err := rows.Scan(&u.ID, &u.Name, &u.Settings, &u.LastSeen)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over user rows: %w", err)
	}

	return users, nil
}
