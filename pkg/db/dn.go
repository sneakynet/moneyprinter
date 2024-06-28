package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// DNCreate sets up a brand new DN that sits on top of a line.
func (db *DB) DNCreate(dn *types.DN) (uint, error) {
	res := db.d.Create(dn)
	return dn.ID, res.Error
}

// DNListByAccountID filters the list of DNs by account.
func (db *DB) DNListByAccountID(acct uint) ([]types.DN, error) {
	dns := []types.DN{}
	res := db.d.Where(&types.DN{AccountID: acct}).Find(&dns)
	return dns, res.Error
}

// DNGet returns detailed information on a single DN by its numeric
// ID.  The numeric ID is **not** the directory number.
func (db *DB) DNGet(id uint) (types.DN, error) {
	dn := types.DN{}
	res := db.d.First(&dn, id)
	return dn, res.Error
}

