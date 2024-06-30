package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// DNCreate sets up a brand new DN that sits on top of a line.
func (db *DB) DNCreate(dn *types.DN) (uint, error) {
	res := db.d.Create(dn)
	return dn.ID, res.Error
}

// DNList filters the list of DNs by the provided instance.
func (db *DB) DNList(filter *types.DN) ([]types.DN, error) {
	dns := []types.DN{}
	res := db.d.Where(filter).Find(&dns)
	return dns, res.Error
}

// DNGet returns detailed information on a single DN selected by the
// parameter in the filter instance.
func (db *DB) DNGet(filter *types.DN) (types.DN, error) {
	dn := types.DN{}
	res := db.d.Where(filter).First(&dn)
	return dn, res.Error
}
