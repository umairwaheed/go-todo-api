package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates the todos table
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create the todos table if it doesn't exist
	if err := createTodosTable(db); err != nil {
		return nil, fmt.Errorf("failed to create todos table: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// createTodosTable creates the todos table with the required schema
func createTodosTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
} 