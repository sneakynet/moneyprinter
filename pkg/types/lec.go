package types

import (
	"gorm.io/gorm"
)

// LEC or Local Exchange Company is an entity that actually provides
// service and receives the payout from the bill.
type LEC struct {
	gorm.Model

	ID      uint
	Name    string
	Byline  string
	Contact string
	Website string
	Logo    Logo
}

// Logo is an abuse of the database to store a picture in it.
type Logo struct {
	gorm.Model

	ID    uint
	LECID uint

	Bytes []byte
}
