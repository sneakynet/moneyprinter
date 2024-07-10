package bill

import (
	"log/slog"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// NewProcessor returns a new processor that can generate multiple
// bills.
func NewProcessor(opts ...ProcessorOption) *Processor {
	p := new(Processor)

	for _, o := range opts {
		o(p)
	}

	return p
}

// Preload fetches a bunch of information from the database so that we
// can run multiple bills without needing to re-fetch this
// information.
func (p *Processor) Preload() error {
	srcFees, err := p.db.FeeList(nil)
	if err != nil {
		slog.Error("Error preloading billing information", "error", err)
		return err
	}

	p.fees = make(map[string][]Fee)

	for _, fee := range srcFees {
		newFee, err := NewDynamicFee(fee.Name, fee.Expression)
		if err != nil {
			continue
		}
		p.fees[fee.AppliesTo] = append(p.fees[fee.AppliesTo], newFee)
	}

	return nil
}

// BillAccount uses the preloaded fee information and applies that to
// an account.  WARNING: This fetches a lot of data to work out what
// services the account has and has consumed, and what it needs to be
// charged for.
func (p *Processor) BillAccount(ac types.Account) (Bill, error) {
	b := Bill{
		Account: ac,
	}

	fctx := FeeContext{Account: ac}

	// First bill anything for just having the account
	for _, fee := range p.fees["account"] {
		l := fee.Evaluate(fctx)
		l.Item = ac.BillText()
		if l.Cost == 0 {
			continue
		}
		b.Lines = append(b.Lines, l)
	}

	// Now we fetch all the lines for the account and bill for
	// those.
	lines, err := p.db.LineList(&types.Line{AccountID: ac.ID})
	if err != nil {
		slog.Error("Error retrieving lines to bill", "account", ac.ID, "error", err)
		return Bill{}, err
	}
	for _, line := range lines {
		for _, fee := range p.fees["line"] {
			fctx.Line = line
			l := fee.Evaluate(fctx)
			l.Item = line.BillText()
			if l.Cost == 0 {
				continue
			}
			b.Lines = append(b.Lines, l)
		}
	}

	// Now the circuits.  This covers things like dry loops and
	// use of our plant that doesn't have associated service on
	// it.
	circuits, err := p.db.CircuitList(&types.Circuit{AccountID: ac.ID})
	if err != nil {
		slog.Error("Error retrieving circuits to bill", "account", ac.ID, "error", err)
		return Bill{}, err
	}
	for _, circuit := range circuits {
		for _, fee := range p.fees["circuit"] {
			fctx.Circuit = circuit
			l := fee.Evaluate(fctx)
			l.Item = circuit.BillText()
			if l.Cost == 0 {
				continue
			}
			b.Lines = append(b.Lines, l)
		}
	}

	// Finally we bill for DNs that are receiving service.  We
	// also in this loop fetch all the CDRs which have the given
	// DN as the origin number and bill for all of those calls as
	// well.
	dns, err := p.db.DNList(&types.DN{AccountID: ac.ID})
	if err != nil {
		slog.Error("Error retreiving DNs to bill", "account", ac.ID, "error", err)
		return Bill{}, err
	}

	for _, dn := range dns {
		for _, fee := range p.fees["dn"] {
			fctx.DN = dn
			l := fee.Evaluate(fctx)
			l.Item = dn.BillText()
			if l.Cost == 0 {
				continue
			}
			b.Lines = append(b.Lines, l)
		}

		cdrs, err := p.db.CDRList(&types.CDR{CLID: dn.Number})
		if err != nil {
			slog.Error("Error retreiving CDRs to bill", "account", ac.ID, "dn", dn.Number, "error", err)
			return Bill{}, err
		}

		for _, cdr := range cdrs {
			for _, fee := range p.fees["cdr"] {
				fctx.CDR = cdr
				l := fee.Evaluate(fctx)
				l.Item = cdr.BillText()
				if l.Cost == 0 {
					continue
				}
				b.Lines = append(b.Lines, l)
			}
		}
	}

	return b, nil
}

// Cost is the total value of the entire Bill.
func (b Bill) Cost() int {
	total := 0
	for _, line := range b.Lines {
		total += line.Cost
	}
	return total
}
