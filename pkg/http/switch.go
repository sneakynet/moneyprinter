package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewSwitchList(w http.ResponseWriter, r *http.Request) {
	switches, err := s.d.SwitchList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "p2/views/switch_list.p2", pongo2.Context{"switches": switches})
}

func (s *Server) uiViewSwitchUpsert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	sw := types.Switch{
		ID:   s.strToUint(chi.URLParam(r, "id")),
		Name: r.FormValue("switch_name"),
		CLLI: r.FormValue("switch_clli"),
		Type: r.FormValue("switch_type"),
	}

	_, err := s.d.SwitchSave(&sw)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/switches", http.StatusSeeOther)
}

func (s *Server) uiViewSwitchFormCreate(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/switch_create.p2", nil)
}

func (s *Server) uiViewSwitchFormEdit(w http.ResponseWriter, r *http.Request) {
	sw, err := s.d.SwitchGet(&types.Switch{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/switch_create.p2", pongo2.Context{"switch": sw})
}

func (s *Server) uiViewSwitchDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.SwitchDelete(&types.Switch{ID: s.strToUint(chi.URLParam(r, "id"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/switches", http.StatusSeeOther)
}
