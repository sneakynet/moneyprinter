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
	tables := []interface{}{
		&types.Account{},
		&types.Circuit{},
		&types.DN{},
		&types.Line{},
		&types.Switch{},
		&types.Equipment{},
		&types.Wirecenter{},
		&types.CDR{},
		&types.Fee{},
	}

	if err := db.d.AutoMigrate(tables...); err != nil {
		return err
	}

	return nil
}
