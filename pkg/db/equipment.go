package db

import (
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// EquipmentCreate sets up a brand new Equipment.
func (db *DB) EquipmentCreate(eqpmnt *types.Equipment) (uint, error) {
	res := db.d.Create(eqpmnt)
	return eqpmnt.ID, res.Error
}

// EquipmentList filters the list of Equipments by the provided
// instance.
func (db *DB) EquipmentList(filter *types.Equipment) ([]types.Equipment, error) {
	eqpmnts := []types.Equipment{}
	res := db.d.Where(filter).Find(&eqpmnts)
	return eqpmnts, res.Error
}

// EquipmentGet returns detailed information on a single Equipment
// selected by the parameter in the filter instance.
func (db *DB) EquipmentGet(filter *types.Equipment) (types.Equipment, error) {
	eqpmnt := types.Equipment{}
	res := db.d.Where(filter).First(&eqpmnt)
	return eqpmnt, res.Error
}
