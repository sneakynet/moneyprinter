package types

// FeeType identifies whether a fee is static or dynamic.
type FeeType uint8

const (
	_ FeeType = iota

	// StaticFee is a fee that is a fixed amount and is applied
	// across all matches.
	StaticFee

	// DynamicFee is a more complex fee that is evaluated inside a
	// FeeContext which contains detailed information about an
	// account and its services.
	DynamicFee
)

// A Fee is an individual line item that comprises a bill.  A bill is
// composed of fees as calulated for an account.  Fees can be either
// static (meaning the charge amount is fixed) or dynamic (the charge
// amount is calculated dynamically).  Fees can match against many
// different facets of an account and are evaluated within a
// FeeContext which includes details about an account's lines,
// circuits, services, and CDR information.
type Fee struct {
	ID uint

	Name      string
	Type      FeeType
	AppliesTo string

	StaticCost  uint
	DynamicExpr string
}
