package types

import (
	"fmt"
)

// BillText formats the text that will be displayed or this item on a
// Bill.
func (a Account) BillText() string {
	return fmt.Sprintf("Account #%d", a.ID)
}

// BillText formats the text that will be displayed or this item on a
// Bill.
func (l Line) BillText() string {
	return fmt.Sprintf("Line #%d (%s)", l.ID, l.Equipment.Type)
}

// BillText formats the text that will be displayed or this item on a
// Bill.
func (c Circuit) BillText() string {
	return fmt.Sprintf("Circuit #%d", c.ID)
}

// BillText formats the text that will be displayed or this item on a
// Bill.
func (dn DN) BillText() string {
	return fmt.Sprintf("Number: %d (%s)", dn.Number, dn.Display)
}

// BillText formats the text that will be displayed or this item on a
// Bill.
func (cdr CDR) BillText() string {
	return fmt.Sprintf("Call to %d (%s)", cdr.DNIS, cdr.End.Sub(cdr.Start))
}
