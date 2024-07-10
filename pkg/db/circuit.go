package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// CircuitCreate inserts a new Circuit into the database.
func (db *DB) CircuitCreate(l *types.Circuit) (uint, error) {
	res := db.d.Create(l)
	return l.ID, res.Error
}

// CircuitList just returns all the circuits that the system knows about.
func (db *DB) CircuitList(filter *types.Circuit) ([]types.Circuit, error) {
	circuits := []types.Circuit{}
	res := db.d.Where(filter).Preload("Lines").Find(&circuits)
	return circuits, res.Error
}

// CircuitGet retrieves a single circuit by its direct numeric ID.
func (db *DB) CircuitGet(filter *types.Circuit) (types.Circuit, error) {
	circuit := types.Circuit{}
	res := db.d.Where(filter).First(&circuit)
	return circuit, res.Error
}
