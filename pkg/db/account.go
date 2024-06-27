package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// AccountCreate creates a new account within the system.
func (db *DB) AccountCreate(a *types.Account) (uint, error) {
	res := db.d.Create(a)
	return a.ID, res.Error
}

// AccountList provides a listing of all accounts in the system.  This
// is not paginated and is one of the limiting points in the system.
func (db *DB) AccountList() ([]types.Account, error) {
	accounts := []types.Account{}
	res := db.d.Find(&accounts)
	return accounts, res.Error
}

// AccountGet returns a single account identified by its specific ID
func (db *DB) AccountGet(id uint) (types.Account, error) {
	acct := types.Account{}
	res := db.d.First(&acct, id)
	return acct, res.Error
}
