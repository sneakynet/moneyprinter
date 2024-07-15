package db

import (
	"gorm.io/gorm/clause"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// EquipmentSave sets up a brand new Equipment.
func (db *DB) EquipmentSave(eqpmnt *types.Equipment) (uint, error) {
	res := db.d.Save(eqpmnt)
	return eqpmnt.ID, res.Error
}

// EquipmentList filters the list of Equipments by the provided
// instance.
func (db *DB) EquipmentList(filter *types.Equipment) ([]types.Equipment, error) {
	eqpmnts := []types.Equipment{}
	res := db.d.Preload(clause.Associations).Where(filter).Find(&eqpmnts)
	return eqpmnts, res.Error
}

// EquipmentGet returns detailed information on a single Equipment
// selected by the parameter in the filter instance.
func (db *DB) EquipmentGet(filter *types.Equipment) (types.Equipment, error) {
	eqpmnt := types.Equipment{}
	res := db.d.Preload(clause.Associations).Where(filter).First(&eqpmnt)
	return eqpmnt, res.Error
}

// EquipmentDelete removes a specific equipment specified by its ID.
func (db *DB) EquipmentDelete(e *types.Equipment) error {
	return db.d.Delete(e).Error
}
