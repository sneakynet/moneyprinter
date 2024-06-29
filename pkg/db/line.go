package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// LineCreate inserts a new Line into the database.
func (db *DB) LineCreate(l *types.Line) (uint, error) {
	res := db.d.Create(l)
	return l.ID, res.Error
}

// LineList just returns all the lines that the system knows about.
func (db *DB) LineList() ([]types.Line, error) {
	lines := []types.Line{}
	res := db.d.Find(&lines)
	return lines, res.Error
}

// LineListByAccountID filters the lines table by the account number
// and returns only lines that are active for a given account.
func (db *DB) LineListByAccountID(id uint) ([]types.Line, error) {
	lines := []types.Line{}
	res := db.d.Where(&types.Line{AccountID: id}).Preload("DNs").Find(&lines)
	return lines, res.Error
}

// LineGet retrieves a single line by its direct numeric ID.
func (db *DB) LineGet(id uint) (types.Line, error) {
	line := types.Line{}
	res := db.d.First(&line, id)
	return line, res.Error
}
