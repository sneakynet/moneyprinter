package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewWirecenterList(w http.ResponseWriter, r *http.Request) {
	wirecenteres, err := s.d.WirecenterList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "p2/views/wirecenter_list.p2", pongo2.Context{"wirecenters": wirecenteres})
}

func (s *Server) uiViewWirecenterUpsert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	sw := types.Wirecenter{
		ID:   s.strToUint(chi.URLParam(r, "id")),
		Name: r.FormValue("wirecenter_name"),
	}

	_, err := s.d.WirecenterSave(&sw)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/wirecenters", http.StatusSeeOther)
}

func (s *Server) uiViewWirecenterFormCreate(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/wirecenter_create.p2", nil)
}

func (s *Server) uiViewWirecenterFormEdit(w http.ResponseWriter, r *http.Request) {
	sw, err := s.d.WirecenterGet(&types.Wirecenter{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/wirecenter_create.p2", pongo2.Context{"wirecenter": sw})
}

func (s *Server) uiViewWirecenterDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.WirecenterDelete(&types.Wirecenter{ID: s.strToUint(chi.URLParam(r, "id"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/wirecenters", http.StatusSeeOther)
}
