package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// AccountCreate creates a new account within the system.
func (db *DB) AccountCreate(a *types.Account) (uint, error) {
	res := db.d.Create(a)
	return a.ID, res.Error
}
