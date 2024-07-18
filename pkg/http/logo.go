package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewLogoForm(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "views/lec/logo.p2", nil)
}

func (s *Server) uiViewLogoSet(w http.ResponseWriter, r *http.Request) {
	lecID := s.strToUint(chi.URLParam(r, "id"))

	if err := r.ParseMultipartForm(5 * 1025 * 1024); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	logo := types.Logo{}

	buf := new(bytes.Buffer)
	f, _, err := r.FormFile("logo")
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	io.Copy(buf, f)
	f.Close()
	logo.Bytes = make([]byte, base64.StdEncoding.EncodedLen(len(buf.Bytes())))
	base64.StdEncoding.Encode(logo.Bytes, buf.Bytes())

	lec, err := s.d.LECGet(&types.LEC{ID: lecID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	lec.Logo = logo

	if _, err := s.d.LECSave(&lec); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/lecs/%d", lecID), http.StatusSeeOther)
}
