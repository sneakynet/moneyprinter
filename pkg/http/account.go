package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiAccountCreateForm(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/account_create.p2", nil)
}

func (s *Server) uiAccountCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.internalErrorHandler(w, r, err)
		return
	}

	aName := r.FormValue("account_name")
	aContact := r.FormValue("account_contact")

	slog.Debug("Want to create account", "name", aName, "contact", aContact)

	id, err := s.d.AccountCreate(&types.Account{Name: aName, Contact: aContact})

	if err != nil {
		s.internalErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/account/%d", id), http.StatusSeeOther)
}

func (s *Server) uiAccountsList(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.d.AccountList()
	if err != nil {
		s.internalErrorHandler(w, r, err)
		return
	}

	s.doTemplate(w, r, "p2/views/account_list.p2", pongo2.Context{"accounts": accounts})
}
