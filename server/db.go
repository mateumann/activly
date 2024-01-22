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

	return sql.Open("postgres", psqlInfo)
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
		return err
	}

	stmt, err := tx.Prepare(
		"CREATE TABLE IF NOT EXISTS public.users (" +
			"id uuid NOT NULL UNIQUE, " +
			fmt.Sprintf("name varchar(%d) NOT NULL UNIQUE, ", dbInitNameLenght) +
			"settings jsonb, " +
			"last_seen timestamp without time zone DEFAULT now()" +
			");",
	)
	if err != nil {
		_ = tx.Rollback()

		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func dbGetUsers(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT id, name, settings, last_seen FROM public.users ORDER BY name;")
}
