package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// CDRCreate inserts a new CDR into the database.  Care should be
// taken to avoid inserting duplicate CDRs.
func (db *DB) CDRCreate(c *types.CDR) (uint, error) {
	res := db.d.Create(c)
	return c.ID, res.Error
}

// CDRList returns a list of CDRs matching the given filter.
func (db *DB) CDRList(filter *types.CDR) ([]types.CDR, error) {
	cdrs := []types.CDR{}
	res := db.d.Where(filter).Find(&cdrs)
	return cdrs, res.Error
}

// CDRGet returns the first CDR that matches the filter, which will be
// the oldest.
func (db *DB) CDRGet(filter *types.CDR) (types.CDR, error) {
	cdr := types.CDR{}
	res := db.d.Where(filter).First(&cdr)
	return cdr, res.Error
}
