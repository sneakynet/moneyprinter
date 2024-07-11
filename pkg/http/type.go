package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/bill"
	"github.com/sneakynet/moneyprinter/pkg/types"
)

// Server handles the HTTP frontend of the application.
type Server struct {
	r  chi.Router
	n  *http.Server
	d  DB
	bp BillProcessor

	tpl *pongo2.TemplateSet
}

// Option configures the server in a composeable way.
type Option func(*Server)

// DB encapsulates all the logic that the webserver expects to be able
// to do.
type DB interface {
	AccountCreate(*types.Account) (uint, error)
	AccountList(*types.Account) ([]types.Account, error)
	AccountGet(*types.Account) (types.Account, error)

	LineCreate(*types.Line) (uint, error)
	LineList(*types.Line) ([]types.Line, error)
	LineGet(*types.Line) (types.Line, error)

	DNCreate(*types.DN) (uint, error)
	DNList(*types.DN) ([]types.DN, error)
	DNGet(*types.DN) (types.DN, error)

	SwitchSave(*types.Switch) (uint, error)
	SwitchList(*types.Switch) ([]types.Switch, error)
	SwitchGet(*types.Switch) (types.Switch, error)
	SwitchDelete(*types.Switch) error

	EquipmentCreate(*types.Equipment) (uint, error)
	EquipmentList(*types.Equipment) ([]types.Equipment, error)
	EquipmentGet(*types.Equipment) (types.Equipment, error)

	WirecenterSave(*types.Wirecenter) (uint, error)
	WirecenterList(*types.Wirecenter) ([]types.Wirecenter, error)
	WirecenterGet(*types.Wirecenter) (types.Wirecenter, error)
	WirecenterDelete(*types.Wirecenter) error

	CDRCreate(*types.CDR) (uint, error)
	CDRList(*types.CDR) ([]types.CDR, error)
	CDRGet(*types.CDR) (types.CDR, error)

	FeeSave(*types.Fee) (uint, error)
	FeeList(*types.Fee) ([]types.Fee, error)
	FeeGet(*types.Fee) (types.Fee, error)
	FeeDelete(*types.Fee) error
}

// BillProcessor actually generates and returns bills.
type BillProcessor interface {
	Preload() error
	BillAccount(types.Account) (bill.Bill, error)
}
