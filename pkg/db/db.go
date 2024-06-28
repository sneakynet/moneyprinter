package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// New returns a new database storage layer.
func New() (*DB, error) {
	return new(DB), nil
}

// Connect sets up the database connection
func (db *DB) Connect(file string) error {
	d, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return err
	}
	db.d = d
	return nil
}

// Migrate reconciles the database schema with the
func (db *DB) Migrate() error {
	if err := db.d.AutoMigrate(&types.Account{}); err != nil {
		return err
	}

	if err := db.d.AutoMigrate(&types.Circuit{}); err != nil {
		return err
	}

	if err := db.d.AutoMigrate(&types.DN{}); err != nil {
		return err
	}

	if err := db.d.AutoMigrate(&types.Line{}); err != nil {
		return err
	}
	return nil
}
