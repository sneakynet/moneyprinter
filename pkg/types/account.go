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

	CircuitID uint
	AccountID uint
	Type      string
	DNs       []DN
	Equipment string
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
