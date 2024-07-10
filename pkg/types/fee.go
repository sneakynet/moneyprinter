package types

// A Fee is an individual line item that comprises a bill.  A bill is
// composed of fees as calulated for an account.  Fees can match
// against many different facets of an account and are evaluated
// within a FeeContext which includes details about an account's
// lines, circuits, services, and CDR information.
type Fee struct {
	ID uint

	Name       string
	AppliesTo  string
	Expression string
}
