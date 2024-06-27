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
	Contact string

	Circuits []Circuit
}

// Circuit specifies a single connection that is paid for by an
// account.
type Circuit struct {
	gorm.Model

	ID        uint
	AccountID uint

	Location string
	Type     string
	DN       uint
}
