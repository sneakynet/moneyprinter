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
func (db *DB) LineList(filter *types.Line) ([]types.Line, error) {
	lines := []types.Line{}
	res := db.d.Where(filter).Preload("DNs").Find(&lines)
	return lines, res.Error
}

// LineGet retrieves a single line by its direct numeric ID.
func (db *DB) LineGet(filter *types.Line) (types.Line, error) {
	line := types.Line{}
	res := db.d.Where(filter).First(&line)
	return line, res.Error
}
