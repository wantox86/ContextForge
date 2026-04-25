// Package storage handles database connection and migrations.
package storage

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps the GORM database connection.
type DB struct {
	conn *gorm.DB
}

// New opens a SQLite database at the given path and runs auto-migrations.
func New(dbPath string) (*DB, error) {
	conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("storage: open database: %w", err)
	}

	if err := conn.AutoMigrate(&Template{}, &Rule{}, &Project{}); err != nil {
		return nil, fmt.Errorf("storage: auto migrate: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Conn returns the underlying GORM connection.
func (db *DB) Conn() *gorm.DB {
	return db.conn
}
