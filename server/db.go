package server

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	dbInitNameLenght = 120
)

func dbConnect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection to database: %w", err)
	}

	return db, nil
}

func dbClose(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func dbInfo(db *sql.DB) string {
	return fmt.Sprintf("%v", db)
}

func dbInit(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin a db transaction: %w", err)
	}

	_, err = tx.Exec(
		"CREATE TABLE IF NOT EXISTS public.users " + "(" +
			"id uuid PRIMARY KEY DEFAULT gen_random_uuid(), " +
			fmt.Sprintf("name varchar(%d) NOT NULL UNIQUE, ", dbInitNameLenght) +
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

func dbGetUsers(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, name, settings, last_seen FROM public.users ORDER BY last_seen;")
	if err != nil {
		return nil, fmt.Errorf("failed querying users: %w", err)
	}

	return rows, nil
}
