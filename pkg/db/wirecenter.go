package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// WirecenterCreate sets up a brand new Wirecenter that sits on top of
// a line.
func (db *DB) WirecenterCreate(wc *types.Wirecenter) (uint, error) {
	res := db.d.Create(wc)
	return wc.ID, res.Error
}

// WirecenterList filters the list of Wirecenters by the provided
// instance.
func (db *DB) WirecenterList(filter *types.Wirecenter) ([]types.Wirecenter, error) {
	wcs := []types.Wirecenter{}
	res := db.d.Where(filter).Find(&wcs)
	return wcs, res.Error
}

// WirecenterGet returns detailed information on a single Wirecenter
// selected by the parameter in the filter instance.
func (db *DB) WirecenterGet(filter *types.Wirecenter) (types.Wirecenter, error) {
	wc := types.Wirecenter{}
	res := db.d.Where(filter).First(&wc)
	return wc, res.Error
}
