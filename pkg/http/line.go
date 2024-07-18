package http

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewSwitchLineList(w http.ResponseWriter, r *http.Request) {
	swID := s.strToUint(chi.URLParam(r, "id"))
	s.uiViewSwitchLineListFilter(w, r, &types.Line{SwitchID: swID})
}

func (s *Server) uiViewSwitchLineListAll(w http.ResponseWriter, r *http.Request) {
	s.uiViewSwitchLineListFilter(w, r, nil)
}

func (s *Server) uiViewSwitchLineListFilter(w http.ResponseWriter, r *http.Request, f *types.Line) {
	lines, err := s.d.LineList(f)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/line/list.p2", pongo2.Context{"lines": lines, "switch_id": chi.URLParam(r, "id")})
}

func (s *Server) uiViewSwitchLineDetail(w http.ResponseWriter, r *http.Request) {
	id := s.strToUint(chi.URLParam(r, "lid"))

	line, err := s.d.LineGet(&types.Line{ID: id})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/line/detail.p2", pongo2.Context{"line": line})
}

func (s *Server) uiViewSwitchLineFormCreate(w http.ResponseWriter, r *http.Request) {
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

	equipment, err := s.d.EquipmentList(&types.Equipment{SwitchID: swID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"switch":    sw,
		"accounts":  accounts,
		"equipment": equipment,
		"linetypes": []string{"FXS-LOOP-START"},
	}

	s.doTemplate(w, r, "views/line/create.p2", ctx)
}

func (s *Server) uiViewSwitchLineFormEdit(w http.ResponseWriter, r *http.Request) {
	equipment, err := s.d.EquipmentList(&types.Equipment{SwitchID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	line, err := s.d.LineGet(&types.Line{ID: s.strToUint(chi.URLParam(r, "lid"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"equipment": equipment,
		"accounts":  []types.Account{line.Account},
		"line":      line,
	}

	s.doTemplate(w, r, "views/line/create.p2", ctx)
}

func (s *Server) uiViewSwitchLineUpsert(w http.ResponseWriter, r *http.Request) {
	swID := s.strToUint(chi.URLParam(r, "id"))
	lID := s.strToUint(chi.URLParam(r, "lid"))

	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	l := types.Line{
		ID:          lID,
		AccountID:   s.strToUint(r.FormValue("account_id")),
		SwitchID:    swID,
		EquipmentID: s.strToUint(r.FormValue("equipment_id")),
	}

	lID, err := s.d.LineSave(&l)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/switches/%d/lines/%d", swID, lID), http.StatusSeeOther)

}

func (s *Server) uiViewSwitchLineDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.LineDelete(&types.Line{ID: s.strToUint(chi.URLParam(r, "lid"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/lines/", http.StatusSeeOther)
}
