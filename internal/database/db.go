package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	} // test connection
	return db, nil
}

func RunMigrations(db *sql.DB) error { // go migrations are manual, unlike ef
	content, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	} // thats new, %v has different formats as well, neat
	_, err = db.Exec(string(content))
	return err
}
