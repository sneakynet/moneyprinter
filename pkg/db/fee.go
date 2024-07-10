package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// FeeSave stores a new fee into the system.
func (db *DB) FeeSave(f *types.Fee) (uint, error) {
	res := db.d.Save(f)
	return f.ID, res.Error
}

// FeeList retrieves fees from the database.
func (db *DB) FeeList(f *types.Fee) ([]types.Fee, error) {
	fees := []types.Fee{}
	res := db.d.Where(f).Find(&fees)
	return fees, res.Error
}

// FeeGet returns a single fee matching the provided filter.
func (db *DB) FeeGet(f *types.Fee) (types.Fee, error) {
	res := db.d.Where(f).First(&f)
	return *f, res.Error
}

// FeeDelete removes a fee from the database.  Never run this because
// it will cause revenue to decrease.
func (db *DB) FeeDelete(f *types.Fee) error {
	res := db.d.Delete(f)
	return res.Error
}
