package bill

// StaticFee encapsulates all the logic of a fee that has a fixed
// cost.
type StaticFee struct {
	cost uint
	name string
}

// NewStaticFee sets up a StaticFee with a given cost and name as
// defined in the database.  The fee doesn't need to know what it
// charges for since static fees have the same price everywhere.
func NewStaticFee(cost uint, name string) Fee {
	return StaticFee{cost: cost, name: name}
}

// Evaluate resolves the fee to a line item in a specific FeeContext.
// For static fees, this just copies the FeeContext's item reference
// into the LineItem
func (sf StaticFee) Evaluate(fc FeeContext) LineItem {
	return LineItem{Fee: sf.name, Cost: int(sf.cost)}
}
