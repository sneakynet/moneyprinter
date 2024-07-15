package types

import (
	"gorm.io/gorm"
)

// Account provides the top level context to which all circuits are
// billed.
type Account struct {
	gorm.Model

	ID      uint
	Name    string
	Alias   string
	Contact string

	DNs      []DN
	Lines    []Line
	Circuits []Circuit
}

// A DN is associated with a Line and has an identifying name with it.
// This allows the system to provision caller ID records.
type DN struct {
	ID      uint
	Number  uint
	Display string

	LineID    uint
	AccountID uint
}

// A Line has one or more DNs and is carried by a circuit.  The line
// is the basic billable unit of access, which may have one or more
// DNs associated with it.
type Line struct {
	ID uint

	CircuitID   uint
	Account     Account
	AccountID   uint
	Switch      Switch
	SwitchID    uint
	Equipment   Equipment
	EquipmentID uint
	Type        string
	DNs         []DN
}

// Wirecenter represents a single location to which wire comes back
// to.  A central location if you will.
type Wirecenter struct {
	ID uint

	Name      string
	Equipment []Equipment
}

// Switch represents a switch with some amount of capacity on it.
type Switch struct {
	ID uint

	Name string
	CLLI string
	Type string

	Lines     []Line
	Equipment []Equipment
}

// Equipment is used to serve lines and typically represents a port on
// a switch
type Equipment struct {
	ID uint

	Switch       Switch
	SwitchID     uint
	Wirecenter   Wirecenter
	WirecenterID uint
	Name         string
	Description  string
	Port         string
}

// Circuit specifies a single connection that is paid for by an
// account.
type Circuit struct {
	gorm.Model

	ID        uint
	AccountID uint

	Location string
	Type     string
	Lines    []Line
}
