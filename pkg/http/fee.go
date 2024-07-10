package http

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewFeeCreate(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{
		"BillableItems": map[string]string{
			"account": "Account",
			"circuit": "Circuit",
			"dn":      "Directory Number",
			"line":    "Provisioned Line",
			"cdr":     "Network Usage (CDR)",
		},
	}

	s.doTemplate(w, r, "p2/views/fee_create.p2", ctx)
}

func (s *Server) uiViewFeeCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	fName := r.FormValue("fee_name")
	fType := r.FormValue("fee_type")
	fApply := r.FormValue("apply_to")
	fParam := r.FormValue("param")

	fee := types.Fee{
		Name:      fName,
		AppliesTo: fApply,
	}

	switch fType {
	case "static":
		fee.Type = types.StaticFee
		fee.StaticCost = s.strToUint(fParam)
	case "dynamic":
		fee.Type = types.DynamicFee
		fee.DynamicExpr = fParam
	}

	id, err := s.d.FeeCreate(&fee)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/fees/%d", id), http.StatusSeeOther)
}

func (s *Server) uiViewFeeList(w http.ResponseWriter, r *http.Request) {
	fees, err := s.d.FeeList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/fee_list.p2", pongo2.Context{"fees": fees})
}
