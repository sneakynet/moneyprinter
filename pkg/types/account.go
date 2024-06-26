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
}

// Circuit specifies a single connection that is paid for by an
// account.
type Circuit struct {
	gorm.Model

	ID       uint
	DN       uint
	Location string
}
