package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// SwitchCreate sets up a brand new Switch.
func (db *DB) SwitchCreate(sw *types.Switch) (uint, error) {
	res := db.d.Create(sw)
	return sw.ID, res.Error
}

// SwitchList filters the list of Switchs by the provided instance.
func (db *DB) SwitchList(filter *types.Switch) ([]types.Switch, error) {
	sws := []types.Switch{}
	res := db.d.Where(filter).Find(&sws)
	return sws, res.Error
}

// SwitchGet returns detailed information on a single Switch selected
// by the parameter in the filter instance.
func (db *DB) SwitchGet(filter *types.Switch) (types.Switch, error) {
	sw := types.Switch{}
	res := db.d.Where(filter).First(&sw)
	return sw, res.Error
}
