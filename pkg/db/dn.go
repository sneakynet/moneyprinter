package db

import (
	"gorm.io/gorm/clause"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// DNSave sets up a brand new DN that sits on top of a line.
func (db *DB) DNSave(dn *types.DN) (uint, error) {
	res := db.d.Save(dn)
	return dn.ID, res.Error
}

// DNList filters the list of DNs by the provided instance.
func (db *DB) DNList(filter *types.DN) ([]types.DN, error) {
	dns := []types.DN{}
	res := db.d.Where(filter).Preload(clause.Associations).Find(&dns)
	return dns, res.Error
}

// DNGet returns detailed information on a single DN selected by the
// parameter in the filter instance.
func (db *DB) DNGet(filter *types.DN) (types.DN, error) {
	dn := types.DN{}
	res := db.d.Where(filter).Preload(clause.Associations).First(&dn)
	return dn, res.Error
}

// DNDelete removes the specified DN from the database.
func (db *DB) DNDelete(dn *types.DN) error {
	return db.d.Delete(dn).Error
}
