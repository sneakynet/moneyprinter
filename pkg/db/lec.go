package db

import (
	"gorm.io/gorm/clause"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// LECSave persists a LEC
func (db *DB) LECSave(l *types.LEC) (uint, error) {
	return l.ID, db.d.Save(l).Error
}

// LECList retrieves all LECs that are known.
func (db *DB) LECList(l *types.LEC) ([]types.LEC, error) {
	lecs := []types.LEC{}
	return lecs, db.d.Preload(clause.Associations).Where(l).Find(&lecs).Error
}

// LECGet fetches a single LEC from the database.
func (db *DB) LECGet(l *types.LEC) (types.LEC, error) {
	lec := types.LEC{}
	res := db.d.Preload(clause.Associations).Where(l).First(&lec)
	return lec, res.Error
}

// LECDelete removes a LEC from the database.  Use with caution as
// this can corrupt references from accounts.
func (db *DB) LECDelete(l *types.LEC) error {
	return db.d.Delete(l).Error
}

// LogoSave persists a logo into the database.
func (db *DB) LogoSave(l *types.Logo) (uint, error) {
	return l.ID, db.d.Save(l).Error
}
