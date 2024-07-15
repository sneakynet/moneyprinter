package db

import (
	"gorm.io/gorm/clause"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// LineSave persists a line into the database.
func (db *DB) LineSave(l *types.Line) (uint, error) {
	res := db.d.Save(l)
	return l.ID, res.Error
}

// LineList just returns all the lines that the system knows about.
func (db *DB) LineList(filter *types.Line) ([]types.Line, error) {
	lines := []types.Line{}
	res := db.d.Where(filter).Preload(clause.Associations).Find(&lines)
	return lines, res.Error
}

// LineGet retrieves a single line by its direct numeric ID.
func (db *DB) LineGet(filter *types.Line) (types.Line, error) {
	line := types.Line{}
	res := db.d.Where(filter).Preload(clause.Associations).First(&line)
	return line, res.Error
}

// LineDelete removes a line from the database.
func (db *DB) LineDelete(l *types.Line) error {
	return db.d.Delete(l).Error
}
