package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

var (
	billableItems = map[string]string{
		"account": "Account",
		"circuit": "Circuit",
		"dn":      "Directory Number",
		"line":    "Provisioned Line",
		"cdr":     "Network Usage (CDR)",
	}
)

func (s *Server) uiViewFeeCreate(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/fee_create.p2", pongo2.Context{"BillableItems": billableItems})
}

func (s *Server) uiViewFeeEditForm(w http.ResponseWriter, r *http.Request) {
	fID := s.strToUint(chi.URLParam(r, "id"))

	fee, err := s.d.FeeGet(&types.Fee{ID: fID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"BillableItems": billableItems,
		"fee":           fee,
	}

	s.doTemplate(w, r, "p2/views/fee_create.p2", ctx)
}

func (s *Server) uiViewFeeUpsertPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	fID := s.strToUint(chi.URLParam(r, "id"))
	fName := r.FormValue("fee_name")
	fApply := r.FormValue("apply_to")
	fExpression := r.FormValue("fee_expression")

	fee := types.Fee{
		ID:         fID,
		Name:       fName,
		AppliesTo:  fApply,
		Expression: fExpression,
	}

	_, err := s.d.FeeSave(&fee)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/fees", http.StatusSeeOther)
}

func (s *Server) uiViewFeeList(w http.ResponseWriter, r *http.Request) {
	fees, err := s.d.FeeList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/fee_list.p2", pongo2.Context{"fees": fees})
}

func (s *Server) uiViewFeeDelete(w http.ResponseWriter, r *http.Request) {
	fID := s.strToUint(chi.URLParam(r, "id"))

	if err := s.d.FeeDelete(&types.Fee{ID: fID}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/fees", http.StatusSeeOther)
}
