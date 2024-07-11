package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// SwitchSave persists a switch into the database.
func (db *DB) SwitchSave(sw *types.Switch) (uint, error) {
	res := db.d.Save(sw)
	return sw.ID, res.Error
}

// SwitchList filters the list of Switchs by the provided instance.
func (db *DB) SwitchList(filter *types.Switch) ([]types.Switch, error) {
	sws := []types.Switch{}
	res := db.d.Preload("Lines").Preload("Equipment").Where(filter).Find(&sws)
	return sws, res.Error
}

// SwitchGet returns detailed information on a single Switch selected
// by the parameter in the filter instance.
func (db *DB) SwitchGet(filter *types.Switch) (types.Switch, error) {
	sw := types.Switch{}
	res := db.d.Where(filter).First(&sw)
	return sw, res.Error
}

// SwitchDelete removes a switch from the database.
func (db *DB) SwitchDelete(sw *types.Switch) error {
	res := db.d.Delete(sw)
	return res.Error
}
