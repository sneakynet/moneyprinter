package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewLECList(w http.ResponseWriter, r *http.Request) {
	lecs, err := s.d.LECList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/lec/list.p2", pongo2.Context{"lecs": lecs})
}

func (s *Server) uiViewLECDetail(w http.ResponseWriter, r *http.Request) {
	lec, err := s.d.LECGet(&types.LEC{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/lec/detail.p2", pongo2.Context{"lec": lec})
}

func (s *Server) uiViewLECUpsert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	sw := types.LEC{
		ID:      s.strToUint(chi.URLParam(r, "id")),
		Name:    r.FormValue("lec_name"),
		Byline:  r.FormValue("lec_byline"),
		Contact: r.FormValue("lec_contact"),
		Website: r.FormValue("lec_website"),
	}

	_, err := s.d.LECSave(&sw)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/lecs", http.StatusSeeOther)
}

func (s *Server) uiViewLECFormCreate(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "views/lec/create.p2", nil)
}

func (s *Server) uiViewLECFormEdit(w http.ResponseWriter, r *http.Request) {
	sw, err := s.d.LECGet(&types.LEC{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/lec/create.p2", pongo2.Context{"lec": sw})
}

func (s *Server) uiViewLECDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.LECDelete(&types.LEC{ID: s.strToUint(chi.URLParam(r, "id"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/lecs", http.StatusSeeOther)
}
