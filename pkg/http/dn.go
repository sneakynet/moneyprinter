package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewDNList(w http.ResponseWriter, r *http.Request) {
	dns, err := s.d.DNList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/dn/list.p2", pongo2.Context{"dns": dns})
}

func (s *Server) uiViewDNDetail(w http.ResponseWriter, r *http.Request) {
	id := s.strToUint(chi.URLParam(r, "id"))

	dn, err := s.d.DNGet(&types.DN{ID: id})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/dn/detail.p2", pongo2.Context{"dn": dn})
}

func (s *Server) uiViewDNFormCreate(w http.ResponseWriter, r *http.Request) {
	swID := s.strToUint(chi.URLParam(r, "id"))

	sw, err := s.d.SwitchGet(&types.Switch{ID: swID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	accounts, err := s.d.AccountList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	lines, err := s.d.LineList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"switch":   sw,
		"accounts": accounts,
		"lines":    lines,
		"dntypes":  []string{"FXS-LOOP-START"},
	}

	s.doTemplate(w, r, "views/dn/create.p2", ctx)
}

func (s *Server) uiViewDNFormEdit(w http.ResponseWriter, r *http.Request) {
	equipment, err := s.d.EquipmentList(&types.Equipment{SwitchID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	dn, err := s.d.DNGet(&types.DN{ID: s.strToUint(chi.URLParam(r, "lid"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"equipment": equipment,
		"dn":        dn,
		"dntypes":   []string{"FXS-LOOP-START"},
	}

	s.doTemplate(w, r, "views/dn/create.p2", ctx)
}

func (s *Server) uiViewDNUpsert(w http.ResponseWriter, r *http.Request) {
	ID := s.strToUint(chi.URLParam(r, "id"))

	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	l := types.DN{
		ID:        ID,
		AccountID: s.strToUint(r.FormValue("account_id")),
		LineID:    s.strToUint(r.FormValue("line_id")),
		Number:    s.strToUint(r.FormValue("dn/number")),
		Display:   r.FormValue("dn_display"),
	}

	_, err := s.d.DNSave(&l)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/dn", http.StatusSeeOther)

}

func (s *Server) uiViewDNDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.DNDelete(&types.DN{ID: s.strToUint(chi.URLParam(r, "id"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/dn/", http.StatusSeeOther)
}
